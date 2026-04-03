package rules

import (
	"testing"

	"github.com/realtong/logi-cli/internal/events"
)

func TestTriggerForEventPrefersGesture(t *testing.T) {
	event := events.DeviceEvent{
		Control: "gesture_button",
		Kind:    events.ButtonDown,
		Gesture: "hold(gesture_button)+move(down)",
	}

	got, ok := triggerForEvent(event)
	if !ok {
		t.Fatal("triggerForEvent() reported no trigger, want gesture trigger")
	}
	if got != "hold(gesture_button)+move(down)" {
		t.Fatalf("triggerForEvent() = %q, want hold(gesture_button)+move(down)", got)
	}
}

func TestTriggerForEventFormatsButtonKinds(t *testing.T) {
	tests := []struct {
		name  string
		event events.DeviceEvent
		want  string
	}{
		{
			name: "button down",
			event: events.DeviceEvent{
				Control: "gesture_button",
				Kind:    events.ButtonDown,
			},
			want: "gesture_button_down",
		},
		{
			name: "button hold",
			event: events.DeviceEvent{
				Control: "gesture_button",
				Kind:    events.ButtonHold,
			},
			want: "gesture_button_hold",
		},
		{
			name: "button up",
			event: events.DeviceEvent{
				Control: "gesture_button",
				Kind:    events.ButtonUp,
			},
			want: "gesture_button_up",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := triggerForEvent(tc.event)
			if !ok {
				t.Fatal("triggerForEvent() reported no trigger")
			}
			if got != tc.want {
				t.Fatalf("triggerForEvent() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTriggerForEventRejectsPointerMotionWithoutGesture(t *testing.T) {
	got, ok := triggerForEvent(events.DeviceEvent{
		Control: "pointer",
		Kind:    events.PointerMove,
		DeltaY:  42,
	})
	if ok {
		t.Fatalf("triggerForEvent() = %q, want no trigger for raw pointer motion", got)
	}
}
