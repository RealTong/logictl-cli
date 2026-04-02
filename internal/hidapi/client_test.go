package hidapi

import "testing"

func TestFakeClientListDevicesReturnsConfiguredDevices(t *testing.T) {
	client := FakeClient{
		Devices: []DeviceInfo{
			{Path: "one", Product: "MX Master 3"},
			{Path: "two", Product: "MX Keys"},
		},
	}

	devices, err := client.ListDevices()
	if err != nil {
		t.Fatalf("ListDevices returned error: %v", err)
	}
	if len(devices) != 2 {
		t.Fatalf("len(devices) = %d, want 2", len(devices))
	}
	if devices[0].Path != "one" || devices[1].Path != "two" {
		t.Fatalf("devices = %#v, want configured order", devices)
	}

	devices[0].Path = "changed"
	if client.Devices[0].Path != "one" {
		t.Fatalf("client devices were modified through returned slice: %#v", client.Devices)
	}
}
