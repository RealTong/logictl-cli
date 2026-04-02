package events

import (
	"strings"
	"testing"
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
