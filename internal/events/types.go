package events

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type RawReport struct {
	DeviceID string
	Bytes    []byte
	At       time.Time
}

type Source interface {
	Stream(ctx context.Context) (<-chan RawReport, <-chan error)
}

type EventKind string

const (
	ButtonDown  EventKind = "button_down"
	ButtonHold  EventKind = "button_hold"
	ButtonUp    EventKind = "button_up"
	PointerMove EventKind = "pointer_move"
	Gesture     EventKind = "gesture"
)

type DeviceEvent struct {
	DeviceID string
	Control  string
	Kind     EventKind
	DeltaX   int
	DeltaY   int
	Gesture  string
	At       time.Time
}

func FormatDeviceEvent(event DeviceEvent) string {
	parts := make([]string, 0, 3)
	if !event.At.IsZero() {
		parts = append(parts, event.At.Format(time.RFC3339Nano))
	}
	if event.DeviceID != "" {
		parts = append(parts, event.DeviceID)
	}

	switch {
	case event.Gesture != "":
		parts = append(parts, event.Gesture)
	case event.Control != "" && event.Kind != "":
		switch event.Kind {
		case ButtonDown:
			parts = append(parts, fmt.Sprintf("%s_down", event.Control))
		case ButtonHold:
			parts = append(parts, fmt.Sprintf("%s_hold", event.Control))
		case ButtonUp:
			parts = append(parts, fmt.Sprintf("%s_up", event.Control))
		case PointerMove:
			parts = append(parts, fmt.Sprintf("%s_move dx=%d dy=%d", event.Control, event.DeltaX, event.DeltaY))
		default:
			parts = append(parts, strings.TrimSpace(fmt.Sprintf("%s %s", event.Control, event.Kind)))
		}
	default:
		parts = append(parts, "<empty>")
	}

	return strings.Join(parts, " ")
}
