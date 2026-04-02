package hidapi

import (
	"errors"

	hid "github.com/sstallion/go-hid"
)

type Client interface {
	ListDevices() ([]DeviceInfo, error)
}

type DeviceInfo struct {
	Path            string
	VendorID        uint16
	ProductID       uint16
	ReleaseNumber   uint16
	InterfaceNumber int
	UsagePage       uint16
	Usage           uint16
	SerialNumber    string
	Manufacturer    string
	Product         string
	Transport       string
}

type client struct{}

func NewClient() Client {
	return client{}
}

func (client) ListDevices() ([]DeviceInfo, error) {
	if err := hid.Init(); err != nil {
		return nil, err
	}

	var devices []DeviceInfo
	enumErr := hid.Enumerate(hid.VendorIDAny, hid.ProductIDAny, func(info *hid.DeviceInfo) error {
		devices = append(devices, DeviceInfo{
			Path:            info.Path,
			VendorID:        info.VendorID,
			ProductID:       info.ProductID,
			ReleaseNumber:   info.ReleaseNbr,
			InterfaceNumber: info.InterfaceNbr,
			UsagePage:       info.UsagePage,
			Usage:           info.Usage,
			SerialNumber:    info.SerialNbr,
			Manufacturer:    info.MfrStr,
			Product:         info.ProductStr,
			Transport:       info.BusType.String(),
		})
		return nil
	})

	exitErr := hid.Exit()
	if enumErr != nil || exitErr != nil {
		return nil, errors.Join(enumErr, exitErr)
	}

	return devices, nil
}
