package mxmaster4

import (
	"strings"

	"github.com/realtong/logi-cli/internal/events"
	"github.com/realtong/logi-cli/internal/hidapi"
)

const productID = 0xb042

type Adapter struct {
	thumbButtonHeld bool
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
			return buttonEvent(report, events.ButtonDown), nil
		}
		if deltaX != 0 || deltaY != 0 {
			return pointerMoveEvent(report, deltaX, deltaY), nil
		}
		return events.DeviceEvent{}, nil
	}

	if a.thumbButtonHeld {
		a.thumbButtonHeld = false
		return buttonEvent(report, events.ButtonUp), nil
	}

	if state == 0x00 && deltaX == 0 && deltaY == 0 {
		return events.DeviceEvent{}, nil
	}

	if state == 0x00 && (deltaX != 0 || deltaY != 0) {
		return events.DeviceEvent{}, unsupportedReportError(state)
	}

	return events.DeviceEvent{}, unsupportedReportError(state)
}
