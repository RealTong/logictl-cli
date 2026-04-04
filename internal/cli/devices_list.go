package cli

import (
	"fmt"

	"github.com/realtong/logictl-cli/internal/hidapi"
	"github.com/spf13/cobra"
)

func newDevicesCmd(client hidapi.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devices",
		Short: "List and inspect attached HID devices",
	}

	cmd.AddCommand(newDevicesListCmd(client))
	cmd.AddCommand(newDevicesInspectCmd(client))
	return cmd
}

func newDevicesListCmd(client hidapi.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List attached HID devices",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			devices, err := client.ListDevices()
			if err != nil {
				return err
			}

			for _, device := range devices {
				printDeviceSummary(cmd, device)
			}
			return nil
		},
	}
}

func printDeviceSummary(cmd *cobra.Command, device hidapi.DeviceInfo) {
	if device.Path == "" {
		fmt.Fprintf(cmd.OutOrStdout(), "%04x:%04x %s\n", device.VendorID, device.ProductID, device.Product)
		return
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%04x:%04x %s %s\n", device.VendorID, device.ProductID, device.Product, device.Path)
}

func printDevice(cmd *cobra.Command, device hidapi.DeviceInfo) {
	out := cmd.OutOrStdout()

	fmt.Fprintf(out, "Path: %s\n", device.Path)
	fmt.Fprintf(out, "VID:PID: %04x:%04x\n", device.VendorID, device.ProductID)
	fmt.Fprintf(out, "Release: %04x\n", device.ReleaseNumber)
	fmt.Fprintf(out, "Interface: %d\n", device.InterfaceNumber)
	fmt.Fprintf(out, "Usage Page: %04x\n", device.UsagePage)
	fmt.Fprintf(out, "Usage: %04x\n", device.Usage)
	fmt.Fprintf(out, "Transport: %s\n", device.Transport)
	fmt.Fprintf(out, "Manufacturer: %s\n", device.Manufacturer)
	fmt.Fprintf(out, "Product: %s\n", device.Product)
	fmt.Fprintf(out, "Serial: %s\n", device.SerialNumber)
}
