package macos

import (
	"context"
	"sync"
	"time"

	"github.com/realtong/logi-cli/internal/config"
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
	at       time.Time
	settings config.ScrollConfig
}

type scrollMatcher struct {
	window  time.Duration
	mu      sync.Mutex
	pending []pendingScrollMatch
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
	for index, candidate := range m.pending {
		if candidate.axis != axis {
			continue
		}

		m.pending = append(m.pending[:index], m.pending[index+1:]...)
		plan := buildScrollRewritePlan(event, candidate.settings)
		if len(plan.Emissions) == 0 {
			return scrollRewritePlan{}, false
		}
		return plan, true
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
}

func buildScrollRewritePlan(event nativeScrollEvent, settings config.ScrollConfig) scrollRewritePlan {
	settings = normalizeScrollConfig(settings)
	if settings.SmoothScroll {
		return scrollRewritePlan{
			Settings:  settings,
			Emissions: splitPixelScroll(event, settings),
		}
	}
	return scrollRewritePlan{
		Settings: settings,
		Emissions: []scrollEmission{{
			Unit:       scrollUnitLine,
			Vertical:   applyScrollDirection(nonZeroOrFallback(event.VerticalLine, event.VerticalPoint), settings),
			Horizontal: applyScrollDirection(nonZeroOrFallback(event.HorizontalLine, event.HorizontalPoint), settings),
		}},
	}
}

func splitPixelScroll(event nativeScrollEvent, settings config.ScrollConfig) []scrollEmission {
	totalVertical := nonZeroOrFallback(event.VerticalPoint, event.VerticalLine*defaultPixelStepScale)
	totalHorizontal := nonZeroOrFallback(event.HorizontalPoint, event.HorizontalLine*defaultPixelStepScale)
	totalVertical = applyScrollDirection(totalVertical, settings)
	totalHorizontal = applyScrollDirection(totalHorizontal, settings)

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

func nonZeroOrFallback(primary, fallback int) int {
	if primary != 0 {
		return primary
	}
	return fallback
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
