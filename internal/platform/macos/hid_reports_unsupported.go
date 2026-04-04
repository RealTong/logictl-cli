//go:build !darwin

package macos

import (
	"context"
	"fmt"

	"github.com/realtong/logictl-cli/internal/events"
)

type unsupportedHIDReportSourceFactory struct{}

func NewHIDReportSourceFactory() HIDReportSourceFactory {
	return unsupportedHIDReportSourceFactory{}
}

func (unsupportedHIDReportSourceFactory) Validate(match HIDReportMatch) error {
	return fmt.Errorf("native HID report capture is not supported on this platform for %s", match.Product)
}

func (unsupportedHIDReportSourceFactory) Open(HIDReportMatch) events.Source {
	return unsupportedSource{}
}

type unsupportedSource struct{}

func (unsupportedSource) Stream(ctx context.Context) (<-chan events.RawReport, <-chan error) {
	reports := make(chan events.RawReport)
	close(reports)
	errs := make(chan error, 1)
	errs <- fmt.Errorf("native HID report capture is not supported on this platform")
	close(errs)
	return reports, errs
}
