package daemon

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/realtong/logictl-cli/internal/devices/mxmaster4"
	"github.com/realtong/logictl-cli/internal/events"
	"github.com/realtong/logictl-cli/internal/hidapi"
)

const mxMaster4DeviceID = "mx-master-4"
const (
	usagePageGenericDesktop = 0x0001
	usageGenericDesktop     = 0x0001
	usageMouse              = 0x0002
	usagePageVendorSpecific = 0xff43
	usageVendorSpecific     = 0x0202
)

type rawSourceFactory func(path string) events.Source
type nativeSourcePlanFactory interface {
	Validate(nativeMatchSpec) error
	Open(nativeMatchSpec) events.Source
}

type nativeMatchSpec struct {
	VendorID     uint16
	ProductID    uint16
	UsagePage    uint16
	Usage        uint16
	SerialNumber string
	Product      string
}

type eventSourcePlan struct {
	path        string
	nativeMatch *nativeMatchSpec
}

type mxMaster4EventSource struct {
	hidClient  hidapi.Client
	openRaw    rawSourceFactory
	openNative nativeSourcePlanFactory
}

func newMXMaster4EventSource(hidClient hidapi.Client, openRaw rawSourceFactory, openNative nativeSourcePlanFactory) eventSource {
	return mxMaster4EventSource{
		hidClient:  hidClient,
		openRaw:    openRaw,
		openNative: openNative,
	}
}

func (s mxMaster4EventSource) Validate() error {
	plan, err := s.resolvePlan()
	if err != nil {
		return err
	}
	if plan.nativeMatch == nil {
		return nil
	}
	if s.openNative == nil {
		return fmt.Errorf("native MX Master 4 HID capture is not available for vendor-specific match %+v", *plan.nativeMatch)
	}
	return s.openNative.Validate(*plan.nativeMatch)
}

func (s mxMaster4EventSource) Stream(ctx context.Context) (<-chan events.DeviceEvent, <-chan error) {
	out := make(chan events.DeviceEvent)
	errs := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errs)

		source, err := s.resolveRawSource()
		if err != nil {
			errs <- err
			return
		}

		rawReports, rawErrs := source.Stream(ctx)
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
				decoded, err := adapter.Decode(report)
				if err != nil {
					if errors.Is(err, mxmaster4.ErrUnsupportedReport) {
						continue
					}
					continue
				}

				for _, event := range decoded {
					for _, normalized := range normalizer.Push(event) {
						select {
						case out <- normalized:
						case <-ctx.Done():
							return
						}
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

func (s mxMaster4EventSource) resolveRawSource() (events.Source, error) {
	plan, err := s.resolvePlan()
	if err != nil {
		return nil, err
	}
	if plan.nativeMatch != nil {
		if s.openNative == nil {
			return nil, fmt.Errorf("native MX Master 4 HID capture is not available for vendor-specific match %+v", *plan.nativeMatch)
		}
		return s.openNative.Open(*plan.nativeMatch), nil
	}
	return s.openRaw(plan.path), nil
}

func (s mxMaster4EventSource) resolvePlan() (eventSourcePlan, error) {
	devices, err := s.hidClient.ListDevices()
	if err != nil {
		return eventSourcePlan{}, err
	}

	return resolveEventSourcePlan(devices, s.openNative != nil)
}

func (s mxMaster4EventSource) resolvePath() (string, error) {
	plan, err := s.resolvePlan()
	if err != nil {
		return "", err
	}
	if plan.nativeMatch != nil {
		return "", fmt.Errorf("MX Master 4 shared BLE HID layout requires native vendor-specific capture")
	}
	return plan.path, nil
}

func resolveEventSourcePlan(devices []hidapi.DeviceInfo, nativeAvailable bool) (eventSourcePlan, error) {
	adapter := mxmaster4.Adapter{}
	groups := make(map[string][]hidapi.DeviceInfo)
	for _, device := range devices {
		if !adapter.Matches(device) || device.Path == "" {
			continue
		}
		groups[device.Path] = append(groups[device.Path], device)
	}

	if len(groups) == 0 {
		return eventSourcePlan{}, fmt.Errorf("no supported MX Master 4 HID device available")
	}

	var safePaths []string
	var unsafePaths []string
	var nativeCandidates []nativeMatchSpec
	for path, group := range groups {
		switch {
		case groupHasPrimaryPointerUsage(group):
			unsafePaths = append(unsafePaths, path)
			if match, ok := nativePassiveMatch(group); ok {
				nativeCandidates = append(nativeCandidates, match)
			}
		case groupHasVendorSpecificUsage(group):
			safePaths = append(safePaths, path)
		}
	}

	sort.Strings(safePaths)
	sort.Strings(unsafePaths)

	switch len(safePaths) {
	case 1:
		return eventSourcePlan{path: safePaths[0]}, nil
	case 0:
		if nativeAvailable {
			switch len(nativeCandidates) {
			case 1:
				spec := nativeCandidates[0]
				return eventSourcePlan{nativeMatch: &spec}, nil
			case 0:
			default:
				return eventSourcePlan{}, fmt.Errorf("multiple MX Master 4 native HID capture candidates found")
			}
		}
		if len(unsafePaths) > 0 {
			return eventSourcePlan{}, fmt.Errorf("unsafe MX Master 4 HID path layout: refusing to open primary pointer path(s) %v because macOS/BLE would seize the mouse", unsafePaths)
		}
		return eventSourcePlan{}, fmt.Errorf("no safe MX Master 4 vendor-specific HID path available")
	default:
		return eventSourcePlan{}, fmt.Errorf("multiple safe MX Master 4 HID paths found: %v", safePaths)
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
		if device.UsagePage == usagePageVendorSpecific && device.Usage == usageVendorSpecific {
			return true
		}
	}
	return false
}

func vendorSpecificMatches(group []hidapi.DeviceInfo) []nativeMatchSpec {
	matches := make([]nativeMatchSpec, 0, 1)
	for _, device := range group {
		if device.UsagePage != usagePageVendorSpecific || device.Usage != usageVendorSpecific {
			continue
		}
		matches = append(matches, nativeMatchSpec{
			VendorID:     device.VendorID,
			ProductID:    device.ProductID,
			UsagePage:    device.UsagePage,
			Usage:        device.Usage,
			SerialNumber: device.SerialNumber,
			Product:      device.Product,
		})
	}
	return matches
}

func nativePassiveMatch(group []hidapi.DeviceInfo) (nativeMatchSpec, bool) {
	for _, device := range group {
		if device.UsagePage != usagePageGenericDesktop || device.Usage != usageMouse {
			continue
		}
		return nativeMatchSpec{
			VendorID:     device.VendorID,
			ProductID:    device.ProductID,
			UsagePage:    device.UsagePage,
			Usage:        device.Usage,
			SerialNumber: device.SerialNumber,
			Product:      device.Product,
		}, true
	}

	for _, device := range group {
		if device.UsagePage != usagePageGenericDesktop || device.Usage != usageGenericDesktop {
			continue
		}
		return nativeMatchSpec{
			VendorID:     device.VendorID,
			ProductID:    device.ProductID,
			UsagePage:    device.UsagePage,
			Usage:        device.Usage,
			SerialNumber: device.SerialNumber,
			Product:      device.Product,
		}, true
	}

	return nativeMatchSpec{}, false
}
