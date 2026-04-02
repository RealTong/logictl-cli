package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/realtong/logi-cli/internal/hidapi"
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
					{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
					{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x40, 0x00, 0x00, 0x20, 0x01, 0x00, 0x00}},
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
	if !strings.Contains(got, "thumb_button_down") {
		t.Fatalf("output = %q, want thumb_button_down", got)
	}
	if !strings.Contains(got, "thumb_button_hold") {
		t.Fatalf("output = %q, want thumb_button_hold", got)
	}
	if !strings.Contains(got, "hold(thumb_button)+move(down)") {
		t.Fatalf("output = %q, want hold(thumb_button)+move(down)", got)
	}
	if strings.Contains(got, "02 40 00 00 00 00 00 00") {
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
					{DeviceID: "mx-master-4", Bytes: []byte{0x02, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
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
	if !strings.Contains(got, "02 40 00 00 00 00 00 00") {
		t.Fatalf("output = %q, want raw bytes", got)
	}
	if strings.Contains(got, "thumb_button_down") {
		t.Fatalf("output = %q, want raw output when --raw is set", got)
	}
}
