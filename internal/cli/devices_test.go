package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/realtong/logi-cli/internal/hidapi"
)

func TestDevicesListPrintsSummaryLines(t *testing.T) {
	cmd := newRootCmd(hidapi.FakeClient{
		Devices: []hidapi.DeviceInfo{
			{
				Path:            "IOService:/AppleACPIPlatformExpert/PCI0@0",
				VendorID:        0x046d,
				ProductID:       0xc548,
				ReleaseNumber:   0x0111,
				InterfaceNumber: 1,
				UsagePage:       0x0001,
				Usage:           0x0002,
				SerialNumber:    "ABC123",
				Manufacturer:    "Logitech",
				Product:         "MX Master 3",
				Transport:       "USB",
			},
			{
				VendorID:  0x046d,
				ProductID: 0xb023,
				Product:   "MX Keys",
			},
		},
	})

	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"devices", "list"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	out := buf.String()
	want := "046d:c548 MX Master 3 IOService:/AppleACPIPlatformExpert/PCI0@0\n046d:b023 MX Keys\n"
	if out != want {
		t.Fatalf("output = %q, want %q", out, want)
	}
	for _, unwanted := range []string{"Path:", "Release:", "Transport:", "Serial:", "[", "]"} {
		if strings.Contains(out, unwanted) {
			t.Fatalf("output unexpectedly included %q:\n%s", unwanted, out)
		}
	}
}

func TestDevicesInspectPrintsMatchingDevice(t *testing.T) {
	cmd := newRootCmd(hidapi.FakeClient{
		Devices: []hidapi.DeviceInfo{
			{Path: "first", Product: "Ignored"},
			{
				Path:            "second",
				VendorID:        0x046d,
				ProductID:       0xb023,
				Manufacturer:    "Logitech",
				Product:         "MX Keys",
				SerialNumber:    "XYZ789",
				Transport:       "Bluetooth",
				InterfaceNumber: 2,
			},
		},
	})

	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"devices", "inspect", "second"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{
		"Path: second",
		"VID:PID: 046d:b023",
		"Product: MX Keys",
		"Manufacturer: Logitech",
		"Serial: XYZ789",
		"Transport: Bluetooth",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "Path: first") {
		t.Fatalf("inspect output included non-matching device:\n%s", out)
	}
}

func TestDevicesListRejectsUnexpectedArgs(t *testing.T) {
	cmd := newRootCmd(hidapi.FakeClient{
		Devices: []hidapi.DeviceInfo{
			{Path: "ignored", Product: "Should Not List"},
		},
	})

	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"devices", "list", "foo"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil, want argument validation error")
	}
	if !strings.Contains(err.Error(), "foo") {
		t.Fatalf("Execute error = %v, want failure mentioning unexpected arg", err)
	}
	if strings.Contains(buf.String(), "Should Not List") {
		t.Fatalf("unexpected device output when args are invalid:\n%s", buf.String())
	}
}
