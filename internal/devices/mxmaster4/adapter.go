package mxmaster4

import (
	"strings"

	"github.com/realtong/logi-cli/internal/events"
	"github.com/realtong/logi-cli/internal/hidapi"
)

const productID = 0xb042
const postReleaseState = 0x20

type Adapter struct {
	thumbButtonHeld  bool
	thumbButtonMoved bool
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

func (a *Adapter) Decode(report events.RawReport) (events.DeviceEvent, error) {
	state, deltaX, deltaY, err := decodeReport(report)
	if err != nil {
		return events.DeviceEvent{}, err
	}

	if state == thumbButtonPressed {
		if !a.thumbButtonHeld {
			a.thumbButtonHeld = true
			a.thumbButtonMoved = false
			return buttonEvent(report, events.ButtonDown), nil
		}
		if deltaX != 0 || deltaY != 0 {
			a.thumbButtonMoved = true
			return pointerMoveEvent(report, deltaX, deltaY), nil
		}
		return events.DeviceEvent{}, nil
	}

	if a.thumbButtonHeld {
		a.thumbButtonHeld = false
		moved := a.thumbButtonMoved
		a.thumbButtonMoved = false
		if !moved && isIdleState(state, deltaX, deltaY) {
			return buttonEvent(report, events.ButtonUp), nil
		}
	}

	if isIgnoredState(state) {
		return events.DeviceEvent{}, nil
	}

	return events.DeviceEvent{}, unsupportedReportError(state)
}

func isIdleState(state byte, deltaX, deltaY int) bool {
	return isIgnoredState(state) && deltaX == 0 && deltaY == 0
}

func isIgnoredState(state byte) bool {
	return state == 0x00 || state == 0x01 || state == postReleaseState
}
