package events

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	hid "github.com/sstallion/go-hid"
)

const defaultReportBufferSize = 64

type hidSource struct {
	path        string
	readTimeout time.Duration
	now         func() time.Time
}

func NewHIDSource(path string) Source {
	return hidSource{
		path:        path,
		readTimeout: 250 * time.Millisecond,
		now:         time.Now,
	}
}

func FormatRawReport(report RawReport) string {
	parts := make([]string, 0, 3)
	if !report.At.IsZero() {
		parts = append(parts, report.At.Format(time.RFC3339Nano))
	}
	if report.DeviceID != "" {
		parts = append(parts, report.DeviceID)
	}
	if len(report.Bytes) == 0 {
		parts = append(parts, "<empty>")
	} else {
		encoded := make([]string, 0, len(report.Bytes))
		for _, b := range report.Bytes {
			encoded = append(encoded, fmt.Sprintf("%02x", b))
		}
		parts = append(parts, strings.Join(encoded, " "))
	}
	return strings.Join(parts, " ")
}

func (s hidSource) Stream(ctx context.Context) (<-chan RawReport, <-chan error) {
	reports := make(chan RawReport)
	errs := make(chan error, 1)

	go func() {
		defer close(reports)
		defer close(errs)

		reportErr := func(err error) {
			if err == nil {
				return
			}
			select {
			case errs <- err:
			default:
			}
		}

		if err := hid.Init(); err != nil {
			reportErr(err)
			return
		}
		defer func() {
			reportErr(hid.Exit())
		}()

		device, err := hid.OpenPath(s.path)
		if err != nil {
			reportErr(err)
			return
		}
		defer func() {
			reportErr(device.Close())
		}()

		buffer := make([]byte, defaultReportBufferSize)
		for {
			if ctx.Err() != nil {
				return
			}

			n, err := device.ReadWithTimeout(buffer, s.readTimeout)
			if err != nil {
				if errors.Is(err, hid.ErrTimeout) && ctx.Err() == nil {
					continue
				}
				if ctx.Err() != nil {
					return
				}
				reportErr(err)
				return
			}
			if n == 0 {
				continue
			}

			report := RawReport{
				DeviceID: s.path,
				Bytes:    append([]byte(nil), buffer[:n]...),
				At:       s.now(),
			}

			select {
			case reports <- report:
			case <-ctx.Done():
				return
			}
		}
	}()

	return reports, errs
}
