package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/realtong/logictl-cli/internal/hidapi"
)

type fakeSource struct {
	events []RawReport
	err    error
}

func (s fakeSource) Stream(context.Context) (<-chan RawReport, <-chan error) {
	eventsCh := make(chan RawReport, len(s.events))
	for _, event := range s.events {
		eventsCh <- event
	}
	close(eventsCh)

	errCh := make(chan error, 1)
	if s.err != nil {
		errCh <- s.err
	}
	close(errCh)
	return eventsCh, errCh
}

func TestTestEventPrintsRawBytes(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newTestEventCmd(fakeSource{
		events: []RawReport{{DeviceID: "mx-master-4", Bytes: []byte{0x10, 0x01, 0xff}}},
	}, buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if !strings.Contains(buf.String(), "10 01 ff") {
		t.Fatalf("output = %q, want hex bytes", buf.String())
	}
}

func TestStreamRawReportsWritesCaptureFile(t *testing.T) {
	buf := new(bytes.Buffer)
	outputPath := filepath.Join(t.TempDir(), "capture.txt")

	err := streamRawReports(context.Background(), fakeSource{
		events: []RawReport{{DeviceID: "mx-master-4", Bytes: []byte{0x10, 0x01, 0xff}}},
	}, buf, outputPath)
	if err != nil {
		t.Fatalf("streamRawReports() returned error: %v", err)
	}

	gotFile, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) returned error: %v", outputPath, err)
	}

	if got := string(gotFile); !strings.Contains(got, "10 01 ff") {
		t.Fatalf("capture file = %q, want hex bytes", got)
	}
	if got := buf.String(); !strings.Contains(got, "10 01 ff") {
		t.Fatalf("stdout = %q, want hex bytes", got)
	}
}

func TestResolveEventDevicePathRejectsMultipleDevicesWithoutPath(t *testing.T) {
	_, err := resolveEventDevicePath(hidapi.FakeClient{
		Devices: []hidapi.DeviceInfo{
			{Path: "first", VendorID: 0x046d, Product: "MX Master 3"},
			{Path: "second", VendorID: 0x046d, Product: "MX Master 4"},
			{Path: "other", VendorID: 0x1234},
		},
	}, "")
	if err == nil {
		t.Fatal("resolveEventDevicePath() returned nil, want multiple-device error")
	}
	if !strings.Contains(err.Error(), "--path") {
		t.Fatalf("resolveEventDevicePath() error = %v, want --path guidance", err)
	}
}

func TestResolveEventDevicePathSelectsSingleSupportedMXMasterCandidate(t *testing.T) {
	got, err := resolveEventDevicePath(hidapi.FakeClient{
		Devices: []hidapi.DeviceInfo{
			{Path: "keyboard", VendorID: 0x046d, Product: "MX Keys"},
			{Path: "mx-master-4", VendorID: 0x046d, Product: "MX Master 4"},
			{Path: "trackpad", VendorID: 0x05ac},
		},
	}, "")
	if err != nil {
		t.Fatalf("resolveEventDevicePath() returned error: %v", err)
	}
	if got != "mx-master-4" {
		t.Fatalf("resolveEventDevicePath() = %q, want %q", got, "mx-master-4")
	}
}

func TestResolveEventDevicePathRejectsUnsupportedLogitechDevice(t *testing.T) {
	_, err := resolveEventDevicePath(hidapi.FakeClient{
		Devices: []hidapi.DeviceInfo{
			{Path: "keyboard", VendorID: 0x046d, Product: "MX Keys"},
			{Path: "trackpad", VendorID: 0x05ac, Product: "Magic Trackpad"},
		},
	}, "")
	if err == nil {
		t.Fatal("resolveEventDevicePath() returned nil, want unsupported-device error")
	}
	if !strings.Contains(err.Error(), "--path") {
		t.Fatalf("resolveEventDevicePath() error = %v, want --path guidance", err)
	}
}

func TestResolveEventDevicePathCollapsesSupportedInterfacesForSamePhysicalDevice(t *testing.T) {
	got, err := resolveEventDevicePath(hidapi.FakeClient{
		Devices: []hidapi.DeviceInfo{
			{
				Path:            "mx-master-4-iface-1",
				VendorID:        0x046d,
				ProductID:       0xc548,
				Product:         "MX Master 4",
				Manufacturer:    "Logitech",
				Transport:       "Bluetooth",
				SerialNumber:    "ABC123",
				InterfaceNumber: 1,
			},
			{
				Path:            "mx-master-4-iface-2",
				VendorID:        0x046d,
				ProductID:       0xc548,
				Product:         "MX Master 4",
				Manufacturer:    "Logitech",
				Transport:       "Bluetooth",
				SerialNumber:    "ABC123",
				InterfaceNumber: 2,
			},
		},
	}, "")
	if err != nil {
		t.Fatalf("resolveEventDevicePath() returned error: %v", err)
	}
	if got != "mx-master-4-iface-1" {
		t.Fatalf("resolveEventDevicePath() = %q, want %q", got, "mx-master-4-iface-1")
	}
}

func TestTestEventDeviceCmdDefaultsToSemanticOutput(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newTestEventDeviceCmd(
		hidapi.FakeClient{
			Devices: []hidapi.DeviceInfo{
				{Path: "mx-master-4", VendorID: 0x046d, Product: "MX Master 4"},
			},
		},
		func(path string) rawSource {
			if path != "mx-master-4" {
				t.Fatalf("openSource path = %q, want %q", path, "mx-master-4")
			}
			return fakeSource{
				events: []RawReport{
					{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
					{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x20, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00}},
					{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
				},
			}
		},
	)
	cmd.SetOut(buf)
	cmd.SetErr(io.Discard)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "gesture_button_down") {
		t.Fatalf("output = %q, want gesture_button_down", got)
	}
	if !strings.Contains(got, "gesture_button_hold") {
		t.Fatalf("output = %q, want gesture_button_hold", got)
	}
	if !strings.Contains(got, "hold(gesture_button)+move(down)") {
		t.Fatalf("output = %q, want hold(gesture_button)+move(down)", got)
	}
	if strings.Contains(got, "02 20 00 00 00 00 00 00") {
		t.Fatalf("output = %q, want semantic events instead of raw bytes", got)
	}
}

func TestTestEventDeviceCmdRawFlagPreservesRawOutput(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newTestEventDeviceCmd(
		hidapi.FakeClient{
			Devices: []hidapi.DeviceInfo{
				{Path: "mx-master-4", VendorID: 0x046d, Product: "MX Master 4"},
			},
		},
		func(path string) rawSource {
			return fakeSource{
				events: []RawReport{
					{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
				},
			}
		},
	)
	cmd.SetOut(buf)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"--raw"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "02 20 00 00 00 00 00 00") {
		t.Fatalf("output = %q, want raw bytes", got)
	}
	if strings.Contains(got, "gesture_button_down") {
		t.Fatalf("output = %q, want raw output when --raw is set", got)
	}
}

func TestTestEventDeviceCmdSemanticModeRejectsUnsupportedAutoSelectedDevice(t *testing.T) {
	cmd := newTestEventDeviceCmd(
		hidapi.FakeClient{
			Devices: []hidapi.DeviceInfo{
				{Path: "mx-master-3", VendorID: 0x046d, ProductID: 0xb023, Product: "MX Master 3"},
			},
		},
		func(path string) rawSource {
			t.Fatalf("openSource called for unsupported semantic device %q", path)
			return fakeSource{}
		},
	)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil, want semantic-device error")
	}
	if !strings.Contains(err.Error(), "MX Master 4") {
		t.Fatalf("Execute error = %v, want MX Master 4 guidance", err)
	}
}

func TestTestEventDeviceCmdSemanticModeRejectsUnsupportedExplicitPath(t *testing.T) {
	cmd := newTestEventDeviceCmd(
		hidapi.FakeClient{
			Devices: []hidapi.DeviceInfo{
				{Path: "mx-master-3", VendorID: 0x046d, ProductID: 0xb023, Product: "MX Master 3"},
			},
		},
		func(path string) rawSource {
			t.Fatalf("openSource called for unsupported semantic device %q", path)
			return fakeSource{}
		},
	)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"--path", "mx-master-3"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil, want unsupported explicit-path error")
	}
	if !strings.Contains(err.Error(), "unsupported semantic capture") {
		t.Fatalf("Execute error = %v, want unsupported semantic capture error", err)
	}
}

func TestResolveSemanticEventDevicePathRejectsMissingExplicitPath(t *testing.T) {
	_, err := resolveSemanticEventDevicePath(hidapi.FakeClient{
		Devices: []hidapi.DeviceInfo{
			{Path: "mx-master-4", VendorID: 0x046d, ProductID: 0xb042, Product: "MX Master 4"},
		},
	}, "stale-path")
	if err == nil {
		t.Fatal("resolveSemanticEventDevicePath() returned nil, want missing-path error")
	}
	if !strings.Contains(err.Error(), "not currently available") {
		t.Fatalf("resolveSemanticEventDevicePath() error = %v, want not-currently-available guidance", err)
	}
}

func TestStreamSemanticEventsWritesUnsupportedReportVisibility(t *testing.T) {
	buf := new(bytes.Buffer)

	err := streamSemanticEvents(context.Background(), fakeSource{
		events: []RawReport{
			{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		},
	}, buf, "")
	if err != nil {
		t.Fatalf("streamSemanticEvents() returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "unsupported_report") {
		t.Fatalf("output = %q, want unsupported_report visibility", got)
	}
	if !strings.Contains(got, "02 80 00 00 00 00 00 00") {
		t.Fatalf("output = %q, want raw bytes for unsupported report", got)
	}
}

func TestStreamSemanticEventsContinuesAfterDecodeError(t *testing.T) {
	buf := new(bytes.Buffer)

	err := streamSemanticEvents(context.Background(), fakeSource{
		events: []RawReport{
			{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
			{DeviceID: "mx-master-4", Bytes: []byte{0x01, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
			{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x20, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00}},
			{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		},
	}, buf, "")
	if err != nil {
		t.Fatalf("streamSemanticEvents() returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "gesture_button_down") {
		t.Fatalf("output = %q, want gesture_button_down before decode error", got)
	}
	if !strings.Contains(got, "ignored_report") {
		t.Fatalf("output = %q, want ignored_report visibility", got)
	}
	if !strings.Contains(got, "unsupported report id") {
		t.Fatalf("output = %q, want decode error details", got)
	}
	if !strings.Contains(got, "hold(gesture_button)+move(down)") {
		t.Fatalf("output = %q, want later semantic gesture after decode error", got)
	}
	if gotCount := strings.Count(got, "ignored_report"); gotCount != 1 {
		t.Fatalf("output = %q, ignored_report count = %d, want 1", got, gotCount)
	}
}

func TestStreamSemanticEventsSuppressesKnownReleaseTailNoise(t *testing.T) {
	buf := new(bytes.Buffer)

	err := streamSemanticEvents(context.Background(), fakeSource{
		events: []RawReport{
			{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
			{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x20, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00}},
			{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x20, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00}},
			{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		},
	}, buf, "")
	if err != nil {
		t.Fatalf("streamSemanticEvents() returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "gesture_button_down") {
		t.Fatalf("output = %q, want gesture_button_down", got)
	}
	if !strings.Contains(got, "gesture_button_hold") {
		t.Fatalf("output = %q, want gesture_button_hold", got)
	}
	if !strings.Contains(got, "hold(gesture_button)+move(down)") {
		t.Fatalf("output = %q, want hold(gesture_button)+move(down)", got)
	}
	if strings.Contains(got, "unsupported_report") {
		t.Fatalf("output = %q, want known release-tail noise suppressed", got)
	}
	if strings.Contains(got, "ignored_report") {
		t.Fatalf("output = %q, want no decode-error noise for known fixture", got)
	}
}

func TestStreamSemanticEventsSkipsKnownNonGestureStates(t *testing.T) {
	buf := new(bytes.Buffer)

	err := streamSemanticEvents(context.Background(), fakeSource{
		events: []RawReport{
			{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x01, 0x00, 0x00, 0xe0, 0xff, 0x00, 0x00}},
			{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}},
			{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
			{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		},
	}, buf, "")
	if err != nil {
		t.Fatalf("streamSemanticEvents() returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "gesture_button_down") {
		t.Fatalf("output = %q, want gesture_button_down", got)
	}
	if !strings.Contains(got, "gesture_button_up") {
		t.Fatalf("output = %q, want gesture_button_up", got)
	}
	if strings.Contains(got, "unsupported_report") {
		t.Fatalf("output = %q, want known non-thumb states suppressed", got)
	}
	if strings.Contains(got, "ignored_report") {
		t.Fatalf("output = %q, want known non-thumb states suppressed", got)
	}
}

func TestStreamSemanticEventsDeduplicatesRepeatedDecodeErrors(t *testing.T) {
	buf := new(bytes.Buffer)

	err := streamSemanticEvents(context.Background(), fakeSource{
		events: []RawReport{
			{DeviceID: "mx-master-4", Bytes: []byte{0x01, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
			{DeviceID: "mx-master-4", Bytes: []byte{0x01, 0x40, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00}},
			{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		},
	}, buf, "")
	if err != nil {
		t.Fatalf("streamSemanticEvents() returned error: %v", err)
	}

	got := buf.String()
	if gotCount := strings.Count(got, "ignored_report"); gotCount != 1 {
		t.Fatalf("output = %q, ignored_report count = %d, want 1", got, gotCount)
	}
	if !strings.Contains(got, "gesture_button_down") {
		t.Fatalf("output = %q, want later semantic event after repeated decode errors", got)
	}
}

func mustLoadFixtureReports(t *testing.T, name string) []RawReport {
	t.Helper()

	path := filepath.Join("..", "..", "testdata", "mxmaster4", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) returned error: %v", path, err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	reports := make([]RawReport, 0, len(lines))
	for _, line := range lines {
		report, err := parseFixtureReportLine(line)
		if err != nil {
			t.Fatalf("parseFixtureReportLine(%q) returned error: %v", line, err)
		}
		reports = append(reports, report)
	}

	return reports
}

func parseFixtureReportLine(line string) (RawReport, error) {
	fields := strings.Fields(line)
	if len(fields) != 9 {
		return RawReport{}, strconv.ErrSyntax
	}

	at, err := time.Parse(time.RFC3339Nano, fields[0])
	if err != nil {
		return RawReport{}, err
	}

	bytes := make([]byte, 0, len(fields)-1)
	for _, field := range fields[1:] {
		value, err := strconv.ParseUint(field, 16, 8)
		if err != nil {
			return RawReport{}, err
		}
		bytes = append(bytes, byte(value))
	}

	return RawReport{
		DeviceID: "mx-master-4",
		At:       at,
		Bytes:    bytes,
	}, nil
}
