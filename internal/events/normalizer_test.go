package events_test

import (
	"testing"
	"time"

	"github.com/realtong/logi-cli/internal/events"
)

func TestNormalizeHoldMoveDown(t *testing.T) {
	stream := []events.DeviceEvent{
		deviceEvent("gesture_button", events.ButtonDown, 0, 0, ""),
		deviceEvent("pointer", events.PointerMove, 0, 20, ""),
		deviceEvent("pointer", events.PointerMove, 0, 20, ""),
		deviceEvent("gesture_button", events.ButtonUp, 0, 0, ""),
	}

	got := events.Normalize(stream, events.NormalizeConfig{GestureDistance: 32})

	if !containsEvent(got, func(event events.DeviceEvent) bool {
		return event.Kind == events.ButtonDown && event.Control == "gesture_button"
	}) {
		t.Fatalf("normalized stream = %#v, want gesture_button down", got)
	}
	if !containsEvent(got, func(event events.DeviceEvent) bool {
		return event.Kind == events.ButtonHold && event.Control == "gesture_button"
	}) {
		t.Fatalf("normalized stream = %#v, want gesture_button hold", got)
	}
	if !containsEvent(got, func(event events.DeviceEvent) bool {
		return event.Gesture == "hold(gesture_button)+move(down)"
	}) {
		t.Fatalf("normalized stream = %#v, want hold(gesture_button)+move(down)", got)
	}
}

func TestNormalizeHoldMoveLeft(t *testing.T) {
	stream := []events.DeviceEvent{
		deviceEvent("gesture_button", events.ButtonDown, 0, 0, ""),
		deviceEvent("pointer", events.PointerMove, 2, -16, ""),
		deviceEvent("pointer", events.PointerMove, 2, -16, ""),
		deviceEvent("gesture_button", events.ButtonUp, 0, 0, ""),
	}

	got := events.Normalize(stream, events.NormalizeConfig{GestureDistance: 32})

	if !containsEvent(got, func(event events.DeviceEvent) bool {
		return event.Gesture == "hold(gesture_button)+move(left)"
	}) {
		t.Fatalf("normalized stream = %#v, want hold(gesture_button)+move(left)", got)
	}
}

func deviceEvent(control string, kind events.EventKind, deltaX, deltaY int, gesture string) events.DeviceEvent {
	return events.DeviceEvent{
		DeviceID: "mx-master-4",
		Control:  control,
		Kind:     kind,
		DeltaX:   deltaX,
		DeltaY:   deltaY,
		Gesture:  gesture,
		At:       time.Unix(1, 0),
	}
}

func containsEvent(eventsList []events.DeviceEvent, predicate func(events.DeviceEvent) bool) bool {
	for _, event := range eventsList {
		if predicate(event) {
			return true
		}
	}
	return false
}
