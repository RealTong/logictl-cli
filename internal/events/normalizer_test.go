package events_test

import (
	"testing"
	"time"

	"github.com/realtong/logictl-cli/internal/events"
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
		deviceEvent("pointer", events.PointerMove, -20, -2, ""),
		deviceEvent("pointer", events.PointerMove, -20, -2, ""),
		deviceEvent("gesture_button", events.ButtonUp, 0, 0, ""),
	}

	got := events.Normalize(stream, events.NormalizeConfig{GestureDistance: 32})

	if !containsEvent(got, func(event events.DeviceEvent) bool {
		return event.Gesture == "hold(gesture_button)+move(left)"
	}) {
		t.Fatalf("normalized stream = %#v, want hold(gesture_button)+move(left)", got)
	}
}

func TestNormalizeHoldMoveRightWithUpwardDrift(t *testing.T) {
	stream := []events.DeviceEvent{
		deviceEvent("gesture_button", events.ButtonDown, 0, 0, ""),
		deviceEvent("pointer", events.PointerMove, 20, -2, ""),
		deviceEvent("pointer", events.PointerMove, 20, -2, ""),
		deviceEvent("gesture_button", events.ButtonUp, 0, 0, ""),
	}

	got := events.Normalize(stream, events.NormalizeConfig{GestureDistance: 32})

	if !containsEvent(got, func(event events.DeviceEvent) bool {
		return event.Gesture == "hold(gesture_button)+move(right)"
	}) {
		t.Fatalf("normalized stream = %#v, want hold(gesture_button)+move(right)", got)
	}
}

func TestNormalizeHoldMoveUp(t *testing.T) {
	stream := []events.DeviceEvent{
		deviceEvent("gesture_button", events.ButtonDown, 0, 0, ""),
		deviceEvent("pointer", events.PointerMove, -2, -20, ""),
		deviceEvent("pointer", events.PointerMove, -2, -20, ""),
		deviceEvent("gesture_button", events.ButtonUp, 0, 0, ""),
	}

	got := events.Normalize(stream, events.NormalizeConfig{GestureDistance: 32})

	if !containsEvent(got, func(event events.DeviceEvent) bool {
		return event.Gesture == "hold(gesture_button)+move(up)"
	}) {
		t.Fatalf("normalized stream = %#v, want hold(gesture_button)+move(up)", got)
	}
}

func TestNormalizeHoldMoveLeftWithVerticalDrift(t *testing.T) {
	stream := []events.DeviceEvent{
		deviceEvent("gesture_button", events.ButtonDown, 0, 0, ""),
		deviceEvent("pointer", events.PointerMove, -15, -17, ""),
		deviceEvent("pointer", events.PointerMove, -15, -17, ""),
		deviceEvent("gesture_button", events.ButtonUp, 0, 0, ""),
	}

	got := events.Normalize(stream, events.NormalizeConfig{GestureDistance: 32})

	if !containsEvent(got, func(event events.DeviceEvent) bool {
		return event.Gesture == "hold(gesture_button)+move(left)"
	}) {
		t.Fatalf("normalized stream = %#v, want hold(gesture_button)+move(left)", got)
	}
}

func TestNormalizeHoldMoveLeftUsesNormalizedMagnitude(t *testing.T) {
	stream := []events.DeviceEvent{
		deviceEvent("gesture_button", events.ButtonDown, 0, 0, ""),
		deviceEvent("pointer", events.PointerMove, -14, -4, ""),
		deviceEvent("pointer", events.PointerMove, -14, -4, ""),
		deviceEvent("gesture_button", events.ButtonUp, 0, 0, ""),
	}

	got := events.Normalize(stream, events.NormalizeConfig{GestureDistance: 32})

	if !containsEvent(got, func(event events.DeviceEvent) bool {
		return event.Gesture == "hold(gesture_button)+move(left)"
	}) {
		t.Fatalf("normalized stream = %#v, want hold(gesture_button)+move(left) after normalization", got)
	}
}

func TestNormalizerDefersGestureUntilButtonRelease(t *testing.T) {
	normalizer := events.NewNormalizer(events.NormalizeConfig{GestureDistance: 32})

	downEvents := normalizer.Push(deviceEvent("gesture_button", events.ButtonDown, 0, 0, ""))
	moveEvents := normalizer.Push(deviceEvent("pointer", events.PointerMove, -40, -2, ""))

	if containsEvent(downEvents, func(event events.DeviceEvent) bool { return event.Kind == events.Gesture }) {
		t.Fatalf("downEvents = %#v, want no gesture before release", downEvents)
	}
	if containsEvent(moveEvents, func(event events.DeviceEvent) bool { return event.Kind == events.Gesture }) {
		t.Fatalf("moveEvents = %#v, want no gesture before release", moveEvents)
	}

	upEvents := normalizer.Push(deviceEvent("gesture_button", events.ButtonUp, 0, 0, ""))
	if len(upEvents) != 2 {
		t.Fatalf("len(upEvents) = %d, want 2", len(upEvents))
	}
	if upEvents[0].Gesture != "hold(gesture_button)+move(left)" {
		t.Fatalf("upEvents[0] = %#v, want left gesture on release", upEvents[0])
	}
	if upEvents[1].Kind != events.ButtonUp || upEvents[1].Control != "gesture_button" {
		t.Fatalf("upEvents[1] = %#v, want gesture_button up after gesture", upEvents[1])
	}
}

func TestNormalizerUsesFinalAccumulatedDirectionOnRelease(t *testing.T) {
	stream := []events.DeviceEvent{
		deviceEvent("gesture_button", events.ButtonDown, 0, 0, ""),
		deviceEvent("pointer", events.PointerMove, 18, -2, ""),
		deviceEvent("pointer", events.PointerMove, -28, 0, ""),
		deviceEvent("pointer", events.PointerMove, -28, 0, ""),
		deviceEvent("gesture_button", events.ButtonUp, 0, 0, ""),
	}

	got := events.Normalize(stream, events.NormalizeConfig{GestureDistance: 32})

	if containsEvent(got, func(event events.DeviceEvent) bool {
		return event.Gesture == "hold(gesture_button)+move(right)"
	}) {
		t.Fatalf("normalized stream = %#v, want no stale right gesture", got)
	}
	if !containsEvent(got, func(event events.DeviceEvent) bool {
		return event.Gesture == "hold(gesture_button)+move(left)"
	}) {
		t.Fatalf("normalized stream = %#v, want final left gesture", got)
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
