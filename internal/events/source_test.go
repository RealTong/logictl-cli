package events

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestFormatRawReportPrintsHexBytes(t *testing.T) {
	got := FormatRawReport(RawReport{
		DeviceID: "mx-master-4",
		Bytes:    []byte{0x10, 0x01, 0xff},
	})

	if !strings.Contains(got, "10 01 ff") {
		t.Fatalf("FormatRawReport() = %q, want hex bytes", got)
	}
}

func TestHIDSourceRawReportOmitsOpaquePathDeviceID(t *testing.T) {
	now := time.Date(2026, time.April, 2, 12, 0, 0, 0, time.UTC)
	source := hidSource{
		path: "IOService:/opaque/device/path",
		now:  func() time.Time { return now },
	}

	payload := []byte{0x10, 0x01, 0xff}
	got := source.rawReport(payload)

	if got.DeviceID != "" {
		t.Fatalf("rawReport().DeviceID = %q, want empty", got.DeviceID)
	}
	if !got.At.Equal(now) {
		t.Fatalf("rawReport().At = %v, want %v", got.At, now)
	}
	if !bytes.Equal(got.Bytes, payload) {
		t.Fatalf("rawReport().Bytes = %v, want %v", got.Bytes, payload)
	}

	payload[0] = 0x99
	if got.Bytes[0] != 0x10 {
		t.Fatalf("rawReport() did not copy payload bytes: %v", got.Bytes)
	}
}
