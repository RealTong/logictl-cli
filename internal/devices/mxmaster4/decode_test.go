package mxmaster4

import (
	"errors"
	"testing"
	"time"

	"github.com/realtong/logictl-cli/internal/events"
)

func TestDecodeGestureButtonDown(t *testing.T) {
	got, err := Decode(rawReport(0x02, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
	if got[0].Kind != events.ButtonDown || got[0].Control != "gesture_button" {
		t.Fatalf("got[0] = %#v, want gesture_button down", got[0])
	}
}

func TestAdapterDecodeGestureButtonHoldMoveDownAndRelease(t *testing.T) {
	adapter := Adapter{}
	stream := decodeAll(t, &adapter,
		rawReport(0x02, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
		rawReport(0x02, 0x20, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00),
		rawReport(0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
	)

	if !containsEvent(stream, func(event events.DeviceEvent) bool {
		return event.Kind == events.ButtonDown && event.Control == "gesture_button"
	}) {
		t.Fatalf("stream = %#v, want gesture_button down", stream)
	}
	if !containsEvent(stream, func(event events.DeviceEvent) bool {
		return event.Kind == events.PointerMove && event.Control == "pointer" && event.DeltaY > 0
	}) {
		t.Fatalf("stream = %#v, want pointer move down while holding gesture_button", stream)
	}
	if !containsEvent(stream, func(event events.DeviceEvent) bool {
		return event.Kind == events.ButtonUp && event.Control == "gesture_button"
	}) {
		t.Fatalf("stream = %#v, want gesture_button up", stream)
	}
}

func TestAdapterDecodeStandardButtons(t *testing.T) {
	tests := []struct {
		name    string
		control string
		mask    byte
	}{
		{name: "left", control: "left_button", mask: 0x01},
		{name: "right", control: "right_button", mask: 0x02},
		{name: "middle", control: "middle_button", mask: 0x04},
		{name: "back", control: "back_button", mask: 0x08},
		{name: "forward", control: "forward_button", mask: 0x10},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			adapter := Adapter{}
			stream := decodeAll(t, &adapter,
				rawReport(0x02, tc.mask, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
				rawReport(0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
			)

			if len(stream) != 2 {
				t.Fatalf("len(stream) = %d, want 2", len(stream))
			}
			if stream[0].Control != tc.control || stream[0].Kind != events.ButtonDown {
				t.Fatalf("stream[0] = %#v, want %s down", stream[0], tc.control)
			}
			if stream[1].Control != tc.control || stream[1].Kind != events.ButtonUp {
				t.Fatalf("stream[1] = %#v, want %s up", stream[1], tc.control)
			}
		})
	}
}

func TestAdapterDecodeVerticalWheelTicks(t *testing.T) {
	adapter := Adapter{}
	stream := decodeAll(t, &adapter,
		rawReport(0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00),
		rawReport(0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00),
	)

	if len(stream) != 2 {
		t.Fatalf("len(stream) = %d, want 2", len(stream))
	}
	if stream[0].Gesture != "wheel_up" {
		t.Fatalf("stream[0] = %#v, want wheel_up", stream[0])
	}
	if stream[1].Gesture != "wheel_down" {
		t.Fatalf("stream[1] = %#v, want wheel_down", stream[1])
	}
}

func TestAdapterDecodeThumbWheelTicks(t *testing.T) {
	adapter := Adapter{}
	stream := decodeAll(t, &adapter,
		rawReport(0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01),
		rawReport(0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff),
	)

	if len(stream) != 2 {
		t.Fatalf("len(stream) = %d, want 2", len(stream))
	}
	if stream[0].Gesture != "thumb_wheel_right" {
		t.Fatalf("stream[0] = %#v, want thumb_wheel_right", stream[0])
	}
	if stream[1].Gesture != "thumb_wheel_left" {
		t.Fatalf("stream[1] = %#v, want thumb_wheel_left", stream[1])
	}
}

func TestAdapterDecodeModeShiftReport(t *testing.T) {
	adapter := Adapter{}
	stream := decodeAll(t, &adapter,
		rawReport(0x11, 0xff, 0x12, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
		rawReport(0x11, 0xff, 0x12, 0x10, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
	)

	if len(stream) != 4 {
		t.Fatalf("len(stream) = %d, want 4", len(stream))
	}
	if stream[0].Gesture != "mode_shift_button_press" || stream[1].Gesture != "wheel_mode_ratchet" {
		t.Fatalf("ratchet events = %#v %#v, want press + wheel_mode_ratchet", stream[0], stream[1])
	}
	if stream[2].Gesture != "mode_shift_button_press" || stream[3].Gesture != "wheel_mode_free_spin" {
		t.Fatalf("free-spin events = %#v %#v, want press + wheel_mode_free_spin", stream[2], stream[3])
	}
}

func TestAdapterDecodeHapticPanelPressDeduplicatesRepeatedReports(t *testing.T) {
	adapter := Adapter{}
	stream := decodeAll(t, &adapter,
		rawReport(0x02, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00),
		rawReport(0x02, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00),
		rawReport(0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
		rawReport(0x02, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00),
	)

	var count int
	for _, event := range stream {
		if event.Gesture == "haptic_panel_press" {
			count++
		}
	}
	if count != 2 {
		t.Fatalf("haptic_panel_press count = %d, want 2", count)
	}
}

func TestDecodeRejectsUnsupportedReportID(t *testing.T) {
	_, err := Decode(rawReport(0x01, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00))
	if err == nil {
		t.Fatal("Decode returned nil, want unsupported report error")
	}
}

func TestAdapterDecodeRejectsUnknownButtonBits(t *testing.T) {
	adapter := Adapter{}
	_, err := adapter.Decode(rawReport(0x02, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00))
	if !errors.Is(err, ErrUnsupportedReport) {
		t.Fatalf("Decode error = %v, want ErrUnsupportedReport", err)
	}
}

func rawReport(bytes ...byte) events.RawReport {
	return events.RawReport{
		DeviceID: "mx-master-4",
		At:       time.Unix(1, 0),
		Bytes:    append([]byte(nil), bytes...),
	}
}

func decodeAll(t *testing.T, adapter *Adapter, reports ...events.RawReport) []events.DeviceEvent {
	t.Helper()

	var out []events.DeviceEvent
	for _, report := range reports {
		decoded, err := adapter.Decode(report)
		if err != nil {
			t.Fatalf("Decode(% x) returned error: %v", report.Bytes, err)
		}
		out = append(out, decoded...)
	}
	return out
}

func containsEvent(stream []events.DeviceEvent, predicate func(events.DeviceEvent) bool) bool {
	for _, event := range stream {
		if predicate(event) {
			return true
		}
	}
	return false
}
