package mxmaster4

import (
	"strings"

	"github.com/realtong/logictl-cli/internal/events"
	"github.com/realtong/logictl-cli/internal/hidapi"
)

const productID = 0xb042

type Adapter struct {
	pressedButtons byte
}

func (Adapter) Matches(info hidapi.DeviceInfo) bool {
	if info.VendorID != 0x046d {
		return false
	}
	if info.ProductID == productID {
		return true
	}
	return strings.Contains(strings.ToLower(info.Product), "mx master 4")
}

func (a *Adapter) Decode(report events.RawReport) ([]events.DeviceEvent, error) {
	decoded, err := decodeReport(report)
	if err != nil {
		return nil, err
	}

	switch decoded.kind {
	case standardReport:
		return a.decodeStandardReport(report, decoded), nil
	case modeReport:
		return []events.DeviceEvent{
			triggerEvent(report, "mode_shift_button_press"),
			triggerEvent(report, modeGesture(decoded.modeFreeSpin)),
		}, nil
	default:
		return nil, ErrUnsupportedReport
	}
}

func (a *Adapter) decodeStandardReport(report events.RawReport, decoded decodedReport) []events.DeviceEvent {
	out := make([]events.DeviceEvent, 0, len(buttonSpecs)+4)

	changed := a.pressedButtons ^ decoded.buttons
	if changed&buttonMaskHaptic != 0 && decoded.buttons&buttonMaskHaptic != 0 {
		out = append(out, triggerEvent(report, "haptic_panel_press"))
	}
	for _, spec := range buttonSpecs {
		if changed&spec.mask == 0 {
			continue
		}
		kind := events.ButtonUp
		if decoded.buttons&spec.mask != 0 {
			kind = events.ButtonDown
		}
		out = append(out, buttonEvent(report, spec.control, kind))
	}
	a.pressedButtons = decoded.buttons

	if decoded.buttons&buttonMaskGesture != 0 && (decoded.deltaX != 0 || decoded.deltaY != 0) {
		out = append(out, pointerMoveEvent(report, decoded.deltaX, decoded.deltaY))
	}
	out = append(out, emitTicks(report, decoded.wheel, "wheel_down", "wheel_up")...)
	out = append(out, emitTicks(report, decoded.thumbWheel, "thumb_wheel_right", "thumb_wheel_left")...)
	return out
}

func emitTicks(report events.RawReport, delta int, positive, negative string) []events.DeviceEvent {
	if delta == 0 {
		return nil
	}

	gesture := positive
	if delta < 0 {
		gesture = negative
		delta = -delta
	}

	out := make([]events.DeviceEvent, 0, delta)
	for i := 0; i < delta; i++ {
		out = append(out, triggerEvent(report, gesture))
	}
	return out
}

func modeGesture(freeSpin bool) string {
	if freeSpin {
		return "wheel_mode_free_spin"
	}
	return "wheel_mode_ratchet"
}
