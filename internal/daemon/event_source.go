package daemon

import (
	"context"
	"errors"
	"fmt"

	"github.com/realtong/logi-cli/internal/devices/mxmaster4"
	"github.com/realtong/logi-cli/internal/events"
	"github.com/realtong/logi-cli/internal/hidapi"
)

const mxMaster4DeviceID = "mx-master-4"

type rawSourceFactory func(path string) events.Source

type mxMaster4EventSource struct {
	hidClient hidapi.Client
	openRaw   rawSourceFactory
}

func newMXMaster4EventSource(hidClient hidapi.Client, openRaw rawSourceFactory) eventSource {
	return mxMaster4EventSource{
		hidClient: hidClient,
		openRaw:   openRaw,
	}
}

func (s mxMaster4EventSource) Stream(ctx context.Context) (<-chan events.DeviceEvent, <-chan error) {
	out := make(chan events.DeviceEvent)
	errs := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errs)

		path, err := s.resolvePath()
		if err != nil {
			errs <- err
			return
		}

		rawReports, rawErrs := s.openRaw(path).Stream(ctx)
		adapter := mxmaster4.Adapter{}
		normalizer := events.NewNormalizer(events.NormalizeConfig{})

		for rawReports != nil || rawErrs != nil {
			select {
			case report, ok := <-rawReports:
				if !ok {
					rawReports = nil
					continue
				}

				report.DeviceID = mxMaster4DeviceID
				event, err := adapter.Decode(report)
				if err != nil {
					if errors.Is(err, mxmaster4.ErrUnsupportedReport) {
						continue
					}
					continue
				}
				if event == (events.DeviceEvent{}) {
					continue
				}

				for _, normalized := range normalizer.Push(event) {
					select {
					case out <- normalized:
					case <-ctx.Done():
						return
					}
				}
			case err, ok := <-rawErrs:
				if !ok {
					rawErrs = nil
					continue
				}
				if err != nil {
					errs <- err
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, errs
}

func (s mxMaster4EventSource) resolvePath() (string, error) {
	devices, err := s.hidClient.ListDevices()
	if err != nil {
		return "", err
	}

	adapter := mxmaster4.Adapter{}
	for _, device := range devices {
		if !adapter.Matches(device) || device.Path == "" {
			continue
		}
		return device.Path, nil
	}

	return "", fmt.Errorf("no supported MX Master 4 HID device available")
}
