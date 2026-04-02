package daemon

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/realtong/logi-cli/internal/devices/mxmaster4"
	"github.com/realtong/logi-cli/internal/events"
	"github.com/realtong/logi-cli/internal/hidapi"
)

const mxMaster4DeviceID = "mx-master-4"
const (
	usagePageGenericDesktop = 0x0001
	usageGenericDesktop     = 0x0001
	usageMouse              = 0x0002
)

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

func (s mxMaster4EventSource) Validate() error {
	_, err := s.resolvePath()
	return err
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
	groups := make(map[string][]hidapi.DeviceInfo)
	for _, device := range devices {
		if !adapter.Matches(device) || device.Path == "" {
			continue
		}
		groups[device.Path] = append(groups[device.Path], device)
	}

	if len(groups) == 0 {
		return "", fmt.Errorf("no supported MX Master 4 HID device available")
	}

	var safePaths []string
	var unsafePaths []string
	for path, group := range groups {
		switch {
		case groupHasPrimaryPointerUsage(group):
			unsafePaths = append(unsafePaths, path)
		case groupHasVendorSpecificUsage(group):
			safePaths = append(safePaths, path)
		}
	}

	sort.Strings(safePaths)
	sort.Strings(unsafePaths)

	switch len(safePaths) {
	case 1:
		return safePaths[0], nil
	case 0:
		if len(unsafePaths) > 0 {
			return "", fmt.Errorf("unsafe MX Master 4 HID path layout: refusing to open primary pointer path(s) %v because macOS/BLE would seize the mouse", unsafePaths)
		}
		return "", fmt.Errorf("no safe MX Master 4 vendor-specific HID path available")
	default:
		return "", fmt.Errorf("multiple safe MX Master 4 HID paths found: %v", safePaths)
	}
}

func groupHasPrimaryPointerUsage(group []hidapi.DeviceInfo) bool {
	for _, device := range group {
		if device.UsagePage != usagePageGenericDesktop {
			continue
		}
		if device.Usage == usageGenericDesktop || device.Usage == usageMouse {
			return true
		}
	}
	return false
}

func groupHasVendorSpecificUsage(group []hidapi.DeviceInfo) bool {
	for _, device := range group {
		if device.UsagePage >= 0xff00 {
			return true
		}
	}
	return false
}
