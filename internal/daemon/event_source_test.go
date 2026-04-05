package daemon

import (
	"context"
	"testing"
	"time"

	"github.com/realtong/logictl-cli/internal/events"
	"github.com/realtong/logictl-cli/internal/hidapi"
)

type fakeNativeSourceFactory struct {
	validateSpec nativeMatchSpec
	validateErr  error
	openSpec     nativeMatchSpec
	reports      []events.RawReport
	err          error
}

func (f *fakeNativeSourceFactory) Validate(spec nativeMatchSpec) error {
	f.validateSpec = spec
	return f.validateErr
}

func (f *fakeNativeSourceFactory) Open(spec nativeMatchSpec) events.Source {
	f.openSpec = spec
	return fakeRawEventSource{reports: f.reports, err: f.err}
}

type fakeRawEventSource struct {
	reports []events.RawReport
	err     error
}

func (s fakeRawEventSource) Stream(context.Context) (<-chan events.RawReport, <-chan error) {
	reportsCh := make(chan events.RawReport, len(s.reports))
	for _, report := range s.reports {
		reportsCh <- report
	}
	close(reportsCh)

	errCh := make(chan error, 1)
	if s.err != nil {
		errCh <- s.err
	}
	close(errCh)
	return reportsCh, errCh
}

func TestMXMaster4EventSourceValidateFallsBackToNativePassiveCaptureForSharedPrimaryPointerPath(t *testing.T) {
	nativeFactory := &fakeNativeSourceFactory{}
	source := mxMaster4EventSource{hidClient: hidapi.FakeClient{
		Devices: []hidapi.DeviceInfo{
			{
				Path:      "ble-shared",
				VendorID:  0x046d,
				ProductID: 0xb042,
				UsagePage: 0x0001,
				Usage:     0x0002,
				Product:   "MX Master 4",
			},
			{
				Path:      "ble-shared",
				VendorID:  0x046d,
				ProductID: 0xb042,
				UsagePage: 0x0001,
				Usage:     0x0001,
				Product:   "MX Master 4",
			},
			{
				Path:      "ble-shared",
				VendorID:  0x046d,
				ProductID: 0xb042,
				UsagePage: 0xff43,
				Usage:     0x0202,
				Product:   "MX Master 4",
			},
		},
	}, openNative: nativeFactory}

	if err := source.Validate(); err != nil {
		t.Fatalf("Validate() returned error: %v", err)
	}

	wantSpec := nativeMatchSpec{
		VendorID:  0x046d,
		ProductID: 0xb042,
		UsagePage: 0x0001,
		Usage:     0x0002,
		Product:   "MX Master 4",
	}
	if nativeFactory.validateSpec != wantSpec {
		t.Fatalf("Validate() native spec = %#v, want %#v", nativeFactory.validateSpec, wantSpec)
	}
}

func TestMXMaster4EventSourceResolvePathPrefersDedicatedVendorSpecificPath(t *testing.T) {
	source := mxMaster4EventSource{hidClient: hidapi.FakeClient{
		Devices: []hidapi.DeviceInfo{
			{
				Path:      "mouse-path",
				VendorID:  0x046d,
				ProductID: 0xb042,
				UsagePage: 0x0001,
				Usage:     0x0002,
				Product:   "MX Master 4",
			},
			{
				Path:      "vendor-path",
				VendorID:  0x046d,
				ProductID: 0xb042,
				UsagePage: 0xff43,
				Usage:     0x0202,
				Product:   "MX Master 4",
			},
		},
	}}

	got, err := source.resolvePath()
	if err != nil {
		t.Fatalf("resolvePath() returned error: %v", err)
	}
	if got != "vendor-path" {
		t.Fatalf("resolvePath() = %q, want %q", got, "vendor-path")
	}
}

func TestMXMaster4EventSourceStreamUsesNativePassiveCaptureForSharedPrimaryPointerPath(t *testing.T) {
	nativeFactory := &fakeNativeSourceFactory{
		reports: []events.RawReport{
			{
				At:    time.Unix(1, 0),
				Bytes: []byte{0x02, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			},
			{
				At:    time.Unix(2, 0),
				Bytes: []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			},
		},
	}
	source := mxMaster4EventSource{hidClient: hidapi.FakeClient{
		Devices: []hidapi.DeviceInfo{
			{
				Path:      "ble-shared",
				VendorID:  0x046d,
				ProductID: 0xb042,
				UsagePage: 0x0001,
				Usage:     0x0002,
				Product:   "MX Master 4",
			},
			{
				Path:      "ble-shared",
				VendorID:  0x046d,
				ProductID: 0xb042,
				UsagePage: 0x0001,
				Usage:     0x0001,
				Product:   "MX Master 4",
			},
			{
				Path:      "ble-shared",
				VendorID:  0x046d,
				ProductID: 0xb042,
				UsagePage: 0xff43,
				Usage:     0x0202,
				Product:   "MX Master 4",
			},
		},
	}, openNative: nativeFactory}

	eventsCh, errs := source.Stream(context.Background())

	var got []events.DeviceEvent
	for event := range eventsCh {
		got = append(got, event)
	}

	for err := range errs {
		if err != nil {
			t.Fatalf("Stream() reported error: %v", err)
		}
	}

	if len(got) != 2 {
		t.Fatalf("Stream() produced %d events, want 2", len(got))
	}
	if got[0].Control != "gesture_button" || got[0].Kind != events.ButtonDown {
		t.Fatalf("first event = %#v, want gesture button down", got[0])
	}
	if got[1].Control != "gesture_button" || got[1].Kind != events.ButtonUp {
		t.Fatalf("second event = %#v, want gesture button up", got[1])
	}
	if nativeFactory.openSpec.UsagePage != 0x0001 || nativeFactory.openSpec.Usage != 0x0002 {
		t.Fatalf("Open() spec = %#v, want primary mouse usage", nativeFactory.openSpec)
	}
}

func TestMXMaster4EventSourcePreservesGestureMotionFromReleaseReport(t *testing.T) {
	nativeFactory := &fakeNativeSourceFactory{
		reports: []events.RawReport{
			{
				At:    time.Unix(1, 0),
				Bytes: []byte{0x02, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			},
			{
				At:    time.Unix(2, 0),
				Bytes: []byte{0x02, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x00},
			},
		},
	}
	source := mxMaster4EventSource{hidClient: hidapi.FakeClient{
		Devices: []hidapi.DeviceInfo{
			{
				Path:      "ble-shared",
				VendorID:  0x046d,
				ProductID: 0xb042,
				UsagePage: 0x0001,
				Usage:     0x0002,
				Product:   "MX Master 4",
			},
			{
				Path:      "ble-shared",
				VendorID:  0x046d,
				ProductID: 0xb042,
				UsagePage: 0x0001,
				Usage:     0x0001,
				Product:   "MX Master 4",
			},
			{
				Path:      "ble-shared",
				VendorID:  0x046d,
				ProductID: 0xb042,
				UsagePage: 0xff43,
				Usage:     0x0202,
				Product:   "MX Master 4",
			},
		},
	}, openNative: nativeFactory}

	eventsCh, errs := source.Stream(context.Background())

	var got []events.DeviceEvent
	for event := range eventsCh {
		got = append(got, event)
	}

	for err := range errs {
		if err != nil {
			t.Fatalf("Stream() reported error: %v", err)
		}
	}

	if !containsDeviceEvent(got, func(event events.DeviceEvent) bool {
		return event.Kind == events.ButtonHold && event.Control == "gesture_button"
	}) {
		t.Fatalf("stream = %#v, want gesture_button hold before release", got)
	}
	if !containsDeviceEvent(got, func(event events.DeviceEvent) bool {
		return event.Gesture == "hold(gesture_button)+move(right)"
	}) {
		t.Fatalf("stream = %#v, want hold(gesture_button)+move(right) from release report motion", got)
	}
	if !containsDeviceEvent(got, func(event events.DeviceEvent) bool {
		return event.Kind == events.ButtonUp && event.Control == "gesture_button"
	}) {
		t.Fatalf("stream = %#v, want gesture_button up", got)
	}
}

func containsDeviceEvent(eventsList []events.DeviceEvent, predicate func(events.DeviceEvent) bool) bool {
	for _, event := range eventsList {
		if predicate(event) {
			return true
		}
	}
	return false
}
