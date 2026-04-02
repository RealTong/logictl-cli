package cli

import (
	"bytes"
	"context"
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
