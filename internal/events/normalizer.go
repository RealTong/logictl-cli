package events

import "fmt"

const (
	defaultGestureDistance        = 32
	defaultHorizontalAxisWeight   = 1.25
	defaultVerticalDominanceRatio = 1.1
)

type NormalizeConfig struct {
	GestureDistance        int
	HorizontalAxisWeight   float64
	VerticalDominanceRatio float64
}

type Normalizer struct {
	cfg          NormalizeConfig
	heldControl  string
	holdEmitted  bool
	accumulatedX int
	accumulatedY int
}

func NewNormalizer(cfg NormalizeConfig) *Normalizer {
	if cfg.GestureDistance <= 0 {
		cfg.GestureDistance = defaultGestureDistance
	}
	if cfg.HorizontalAxisWeight <= 0 {
		cfg.HorizontalAxisWeight = defaultHorizontalAxisWeight
	}
	if cfg.VerticalDominanceRatio <= 0 {
		cfg.VerticalDominanceRatio = defaultVerticalDominanceRatio
	}
	return &Normalizer{cfg: cfg}
}

func Normalize(stream []DeviceEvent, cfg NormalizeConfig) []DeviceEvent {
	normalizer := NewNormalizer(cfg)
	normalized := make([]DeviceEvent, 0, len(stream))
	for _, event := range stream {
		normalized = append(normalized, normalizer.Push(event)...)
	}
	return normalized
}

func (n *Normalizer) Push(event DeviceEvent) []DeviceEvent {
	switch event.Kind {
	case ButtonDown:
		n.heldControl = event.Control
		n.holdEmitted = false
		n.accumulatedX = 0
		n.accumulatedY = 0
		return []DeviceEvent{event}
	case ButtonUp:
		out := make([]DeviceEvent, 0, 2)
		if n.heldControl != "" && n.gestureMagnitude(n.accumulatedX, n.accumulatedY) >= float64(n.cfg.GestureDistance) {
			direction := n.gestureDirection(n.accumulatedX, n.accumulatedY)
			if direction != "" {
				out = append(out, DeviceEvent{
					DeviceID: event.DeviceID,
					At:       event.At,
					Control:  n.heldControl,
					Kind:     Gesture,
					Gesture:  fmt.Sprintf("hold(%s)+move(%s)", n.heldControl, direction),
				})
			}
		}
		out = append(out, event)

		n.heldControl = ""
		n.holdEmitted = false
		n.accumulatedX = 0
		n.accumulatedY = 0
		return out
	case PointerMove:
		if n.heldControl == "" {
			return nil
		}

		out := make([]DeviceEvent, 0, 2)
		if !n.holdEmitted {
			n.holdEmitted = true
			out = append(out, DeviceEvent{
				DeviceID: event.DeviceID,
				At:       event.At,
				Control:  n.heldControl,
				Kind:     ButtonHold,
			})
		}

		n.accumulatedX += event.DeltaX
		n.accumulatedY += event.DeltaY

		return out
	default:
		if event.Gesture != "" {
			return []DeviceEvent{event}
		}
		return nil
	}
}

func (n *Normalizer) gestureMagnitude(x, y int) float64 {
	horizontal := float64(abs(x)) * n.cfg.HorizontalAxisWeight
	vertical := float64(abs(y))
	if horizontal > vertical {
		return horizontal
	}
	return vertical
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func maxAbs(x, y int) int {
	if abs(x) > abs(y) {
		return abs(x)
	}
	return abs(y)
}

func (n *Normalizer) gestureDirection(x, y int) string {
	horizontal := float64(abs(x)) * n.cfg.HorizontalAxisWeight
	vertical := float64(abs(y))

	switch {
	case y > 0 && vertical >= horizontal*n.cfg.VerticalDominanceRatio:
		return "down"
	case y < 0 && vertical >= horizontal*n.cfg.VerticalDominanceRatio:
		return "up"
	case x > 0:
		return "right"
	case x < 0:
		return "left"
	case y > 0:
		return "down"
	case y < 0:
		return "up"
	default:
		return ""
	}
}
