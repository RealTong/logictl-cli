package macos

import (
	"context"
	"sync"
	"time"

	"github.com/realtong/logictl-cli/internal/config"
)

const (
	defaultScrollMatchWindow = 50 * time.Millisecond
	defaultSmoothSteps       = 4
	defaultPixelStepScale    = 12
)

type ScrollRewriter interface {
	Start(context.Context) error
	Record(deviceID, gesture string, settings config.ScrollConfig, at time.Time)
}

type scrollAxis int

const (
	scrollAxisUnknown scrollAxis = iota
	scrollAxisVertical
	scrollAxisHorizontal
)

type scrollUnit int

const (
	scrollUnitLine scrollUnit = iota
	scrollUnitPixel
)

type nativeScrollEvent struct {
	VerticalLine    int
	HorizontalLine  int
	VerticalPoint   int
	HorizontalPoint int
	At              time.Time
}

type scrollEmission struct {
	Unit       scrollUnit
	Vertical   int
	Horizontal int
}

type scrollRewritePlan struct {
	Settings  config.ScrollConfig
	Emissions []scrollEmission
}

type pendingScrollMatch struct {
	axis     scrollAxis
	gesture  string
	at       time.Time
	settings config.ScrollConfig
}

type axisSuppression struct {
	axis  scrollAxis
	until time.Time
}

type scrollMatcher struct {
	window       time.Duration
	mu           sync.Mutex
	pending      []pendingScrollMatch
	suppressions []axisSuppression
}

func newScrollMatcher(window time.Duration) *scrollMatcher {
	if window <= 0 {
		window = defaultScrollMatchWindow
	}
	return &scrollMatcher{window: window}
}

func (m *scrollMatcher) Record(_ string, gesture string, settings config.ScrollConfig, at time.Time) {
	axis := axisForGesture(gesture)
	if axis == scrollAxisUnknown || !shouldRewriteScroll(settings) {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.pending = append(m.pending, pendingScrollMatch{
		axis:     axis,
		gesture:  gesture,
		at:       at,
		settings: normalizeScrollConfig(settings),
	})
}

func (m *scrollMatcher) Match(event nativeScrollEvent) (scrollRewritePlan, bool) {
	axis := axisForNativeScroll(event)
	if axis == scrollAxisUnknown {
		return scrollRewritePlan{}, false
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.purgeExpiredLocked(event.At)
	if candidate, ok := m.popPendingLocked(axis); ok {
		m.suppressions = append(m.suppressions, axisSuppression{
			axis:  axis,
			until: event.At.Add(m.window),
		})
		plan := buildScrollRewritePlan(candidate.gesture, candidate.settings)
		return plan, true
	}

	if m.hasSuppressionLocked(axis, event.At) {
		return scrollRewritePlan{}, true
	}

	return scrollRewritePlan{}, false
}

func (m *scrollMatcher) purgeExpiredLocked(now time.Time) {
	if now.IsZero() {
		return
	}
	keep := m.pending[:0]
	for _, candidate := range m.pending {
		if now.Sub(candidate.at) <= m.window {
			keep = append(keep, candidate)
		}
	}
	m.pending = keep

	active := m.suppressions[:0]
	for _, suppression := range m.suppressions {
		if !now.After(suppression.until) {
			active = append(active, suppression)
		}
	}
	m.suppressions = active
}

func (m *scrollMatcher) popPendingLocked(axis scrollAxis) (pendingScrollMatch, bool) {
	for index, candidate := range m.pending {
		if candidate.axis != axis {
			continue
		}
		m.pending = append(m.pending[:index], m.pending[index+1:]...)
		return candidate, true
	}
	return pendingScrollMatch{}, false
}

func (m *scrollMatcher) hasSuppressionLocked(axis scrollAxis, now time.Time) bool {
	for _, suppression := range m.suppressions {
		if suppression.axis != axis {
			continue
		}
		if now.After(suppression.until) {
			continue
		}
		return true
	}
	return false
}

func buildScrollRewritePlan(gesture string, settings config.ScrollConfig) scrollRewritePlan {
	settings = normalizeScrollConfig(settings)
	if settings.SmoothScroll {
		return scrollRewritePlan{
			Settings:  settings,
			Emissions: splitPixelScroll(gesture, settings),
		}
	}
	vertical, horizontal := gestureLineDelta(gesture, settings)
	return scrollRewritePlan{
		Settings: settings,
		Emissions: []scrollEmission{{
			Unit:       scrollUnitLine,
			Vertical:   vertical,
			Horizontal: horizontal,
		}},
	}
}

func splitPixelScroll(gesture string, settings config.ScrollConfig) []scrollEmission {
	lineVertical, lineHorizontal := gestureLineDelta(gesture, settings)
	totalVertical := lineVertical * defaultPixelStepScale
	totalHorizontal := lineHorizontal * defaultPixelStepScale

	if totalVertical == 0 && totalHorizontal == 0 {
		return nil
	}

	vParts := splitSteps(totalVertical, defaultSmoothSteps)
	hParts := splitSteps(totalHorizontal, defaultSmoothSteps)
	steps := max(len(vParts), len(hParts))
	out := make([]scrollEmission, 0, steps)
	for i := 0; i < steps; i++ {
		out = append(out, scrollEmission{
			Unit:       scrollUnitPixel,
			Vertical:   valueAt(vParts, i),
			Horizontal: valueAt(hParts, i),
		})
	}
	return out
}

func shouldRewriteScroll(settings config.ScrollConfig) bool {
	return settings.SmoothScroll || settings.Direction == "standard"
}

func normalizeScrollConfig(settings config.ScrollConfig) config.ScrollConfig {
	if settings.Direction == "" {
		settings.Direction = "natural"
	}
	return settings
}

func axisForGesture(gesture string) scrollAxis {
	switch gesture {
	case "wheel_up", "wheel_down":
		return scrollAxisVertical
	case "thumb_wheel_left", "thumb_wheel_right":
		return scrollAxisHorizontal
	default:
		return scrollAxisUnknown
	}
}

func axisForNativeScroll(event nativeScrollEvent) scrollAxis {
	if absInt(event.VerticalLine) > 0 || absInt(event.VerticalPoint) > 0 {
		if absInt(event.VerticalLine)+absInt(event.VerticalPoint) >= absInt(event.HorizontalLine)+absInt(event.HorizontalPoint) {
			return scrollAxisVertical
		}
	}
	if absInt(event.HorizontalLine) > 0 || absInt(event.HorizontalPoint) > 0 {
		return scrollAxisHorizontal
	}
	return scrollAxisUnknown
}

func applyScrollDirection(value int, settings config.ScrollConfig) int {
	if settings.Direction == "standard" {
		return -value
	}
	return value
}

func gestureLineDelta(gesture string, settings config.ScrollConfig) (int, int) {
	vertical, horizontal := naturalGestureDelta(gesture)
	return applyScrollDirection(vertical, settings), applyScrollDirection(horizontal, settings)
}

func naturalGestureDelta(gesture string) (int, int) {
	switch gesture {
	case "wheel_up":
		return 1, 0
	case "wheel_down":
		return -1, 0
	case "thumb_wheel_left":
		return 0, 1
	case "thumb_wheel_right":
		return 0, -1
	default:
		return 0, 0
	}
}

func splitSteps(total, steps int) []int {
	if total == 0 || steps <= 0 {
		return nil
	}

	sign := 1
	if total < 0 {
		sign = -1
		total = -total
	}

	out := make([]int, 0, steps)
	remaining := total
	remainingSteps := steps
	for remaining > 0 && remainingSteps > 0 {
		part := remaining / remainingSteps
		if remaining%remainingSteps != 0 {
			part++
		}
		if part == 0 {
			part = 1
		}
		if part > remaining {
			part = remaining
		}
		out = append(out, sign*part)
		remaining -= part
		remainingSteps--
	}
	return out
}

func valueAt(values []int, index int) int {
	if index < len(values) {
		return values[index]
	}
	return 0
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
