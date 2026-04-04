package macos

import "github.com/realtong/logictl-cli/internal/events"

type HIDReportMatch struct {
	VendorID     uint16
	ProductID    uint16
	UsagePage    uint16
	Usage        uint16
	SerialNumber string
	Product      string
}

type HIDReportSourceFactory interface {
	Validate(HIDReportMatch) error
	Open(HIDReportMatch) events.Source
}
