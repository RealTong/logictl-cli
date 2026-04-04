package mxmaster4

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/realtong/logictl-cli/internal/events"
)

const (
	reportIDStandard = 0x02
	reportIDMode     = 0x11

	buttonMaskLeft    = 0x01
	buttonMaskRight   = 0x02
	buttonMaskMiddle  = 0x04
	buttonMaskBack    = 0x08
	buttonMaskForward = 0x10
	buttonMaskGesture = 0x20
	knownButtonMask   = buttonMaskLeft | buttonMaskRight | buttonMaskMiddle | buttonMaskBack | buttonMaskForward | buttonMaskGesture
)

var ErrUnsupportedReport = errors.New("unsupported mx master 4 report")

type reportKind int

const (
	standardReport reportKind = iota
	modeReport
)

type decodedReport struct {
	kind         reportKind
	buttons      byte
	deltaX       int
	deltaY       int
	wheel        int
	thumbWheel   int
	modeFreeSpin bool
}

type buttonSpec struct {
	mask    byte
	control string
}

var buttonSpecs = []buttonSpec{
	{mask: buttonMaskLeft, control: "left_button"},
	{mask: buttonMaskRight, control: "right_button"},
	{mask: buttonMaskMiddle, control: "middle_button"},
	{mask: buttonMaskBack, control: "back_button"},
	{mask: buttonMaskForward, control: "forward_button"},
	{mask: buttonMaskGesture, control: "gesture_button"},
}

func Decode(report events.RawReport) ([]events.DeviceEvent, error) {
	adapter := Adapter{}
	return adapter.Decode(report)
}

func decodeReport(report events.RawReport) (decodedReport, error) {
	if len(report.Bytes) == 0 {
		return decodedReport{}, fmt.Errorf("mx master 4 report too short: %d", len(report.Bytes))
	}

	switch report.Bytes[0] {
	case reportIDStandard:
		if len(report.Bytes) < 8 {
			return decodedReport{}, fmt.Errorf("mx master 4 report too short: %d", len(report.Bytes))
		}

		buttons := report.Bytes[1]
		if buttons&^knownButtonMask != 0 {
			return decodedReport{}, unsupportedReportError(buttons)
		}

		return decodedReport{
			kind:       standardReport,
			buttons:    buttons,
			deltaX:     int(int8(report.Bytes[3])),
			deltaY:     int(int16(binary.LittleEndian.Uint16(report.Bytes[4:6]))),
			wheel:      int(int8(report.Bytes[6])),
			thumbWheel: int(int8(report.Bytes[7])),
		}, nil
	case reportIDMode:
		if len(report.Bytes) < 5 {
			return decodedReport{}, fmt.Errorf("mx master 4 mode report too short: %d", len(report.Bytes))
		}
		if len(report.Bytes) < 4 || report.Bytes[1] != 0xff || report.Bytes[2] != 0x12 || report.Bytes[3] != 0x10 {
			return decodedReport{}, fmt.Errorf("%w: mode payload=% x", ErrUnsupportedReport, report.Bytes[1:])
		}

		return decodedReport{
			kind:         modeReport,
			modeFreeSpin: report.Bytes[4] == 0x01,
		}, nil
	default:
		return decodedReport{}, fmt.Errorf("mx master 4 unsupported report id: 0x%02x", report.Bytes[0])
	}
}

func buttonEvent(report events.RawReport, control string, kind events.EventKind) events.DeviceEvent {
	return events.DeviceEvent{
		DeviceID: report.DeviceID,
		At:       report.At,
		Control:  control,
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

func triggerEvent(report events.RawReport, gesture string) events.DeviceEvent {
	return events.DeviceEvent{
		DeviceID: report.DeviceID,
		At:       report.At,
		Kind:     events.Gesture,
		Gesture:  gesture,
	}
}

func unsupportedReportError(state byte) error {
	return fmt.Errorf("%w: state=0x%02x", ErrUnsupportedReport, state)
}
