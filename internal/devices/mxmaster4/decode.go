package mxmaster4

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/realtong/logi-cli/internal/events"
)

const (
	reportID           = 0x02
	thumbButtonPressed = 0x40
)

var ErrUnsupportedReport = errors.New("unsupported mx master 4 report")

func Decode(report events.RawReport) (events.DeviceEvent, error) {
	state, deltaX, deltaY, err := decodeReport(report)
	if err != nil {
		return events.DeviceEvent{}, err
	}
	if state != thumbButtonPressed {
		return events.DeviceEvent{}, unsupportedReportError(state)
	}
	if deltaX == 0 && deltaY == 0 {
		return buttonEvent(report, events.ButtonDown), nil
	}
	return pointerMoveEvent(report, deltaX, deltaY), nil
}

func decodeReport(report events.RawReport) (byte, int, int, error) {
	if len(report.Bytes) < 6 {
		return 0, 0, 0, fmt.Errorf("mx master 4 report too short: %d", len(report.Bytes))
	}
	if report.Bytes[0] != reportID {
		return 0, 0, 0, fmt.Errorf("mx master 4 unsupported report id: 0x%02x", report.Bytes[0])
	}

	state := report.Bytes[1]
	deltaX := int(int8(report.Bytes[3]))
	deltaY := int(int16(binary.LittleEndian.Uint16(report.Bytes[4:6])))
	return state, deltaX, deltaY, nil
}

func buttonEvent(report events.RawReport, kind events.EventKind) events.DeviceEvent {
	return events.DeviceEvent{
		DeviceID: report.DeviceID,
		At:       report.At,
		Control:  "thumb_button",
		Kind:     kind,
	}
}

func pointerMoveEvent(report events.RawReport, deltaX, deltaY int) events.DeviceEvent {
	return events.DeviceEvent{
		DeviceID: report.DeviceID,
		At:       report.At,
		Control:  "pointer",
		Kind:     events.PointerMove,
		DeltaX:   deltaX,
		DeltaY:   deltaY,
	}
}

func unsupportedReportError(state byte) error {
	return fmt.Errorf("%w: state=0x%02x", ErrUnsupportedReport, state)
}
