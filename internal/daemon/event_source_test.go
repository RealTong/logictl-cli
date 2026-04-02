package daemon

import (
	"strings"
	"testing"

	"github.com/realtong/logi-cli/internal/hidapi"
)

func TestMXMaster4EventSourceResolvePathRejectsSharedPrimaryPointerPath(t *testing.T) {
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
	}}

	_, err := source.resolvePath()
	if err == nil {
		t.Fatal("resolvePath() returned nil, want unsafe primary pointer error")
	}
	if !strings.Contains(err.Error(), "unsafe") {
		t.Fatalf("resolvePath() error = %v, want unsafe-path guidance", err)
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
