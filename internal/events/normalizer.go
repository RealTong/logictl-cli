package events

import "fmt"

const defaultGestureDistance = 256

type NormalizeConfig struct {
	GestureDistance int
}

type Normalizer struct {
	cfg            NormalizeConfig
	heldControl    string
	holdEmitted    bool
	gestureEmitted bool
	accumulatedX   int
	accumulatedY   int
}

func NewNormalizer(cfg NormalizeConfig) *Normalizer {
	if cfg.GestureDistance <= 0 {
		cfg.GestureDistance = defaultGestureDistance
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
		n.gestureEmitted = false
		n.accumulatedX = 0
		n.accumulatedY = 0
		return []DeviceEvent{event}
	case ButtonUp:
		n.heldControl = ""
		n.holdEmitted = false
		n.gestureEmitted = false
		n.accumulatedX = 0
		n.accumulatedY = 0
		return []DeviceEvent{event}
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

		if !n.gestureEmitted &&
			abs(n.accumulatedY) >= n.cfg.GestureDistance &&
			abs(n.accumulatedY) >= abs(n.accumulatedX) {
			n.gestureEmitted = true
			direction := "down"
			if n.accumulatedY < 0 {
				direction = "up"
			}
			out = append(out, DeviceEvent{
				DeviceID: event.DeviceID,
				At:       event.At,
				Control:  n.heldControl,
				Kind:     Gesture,
				Gesture:  fmt.Sprintf("hold(%s)+move(%s)", n.heldControl, direction),
			})
		}

		return out
	default:
		if event.Gesture != "" {
			return []DeviceEvent{event}
		}
		return nil
	}
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
