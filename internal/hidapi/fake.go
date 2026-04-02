package hidapi

type FakeClient struct {
	Devices []DeviceInfo
	Err     error
}

func (c FakeClient) ListDevices() ([]DeviceInfo, error) {
	if c.Err != nil {
		return nil, c.Err
	}

	devices := make([]DeviceInfo, len(c.Devices))
	copy(devices, c.Devices)
	return devices, nil
}
