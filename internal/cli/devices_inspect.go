package cli

import (
	"fmt"

	"github.com/realtong/logictl-cli/internal/hidapi"
	"github.com/spf13/cobra"
)

func newDevicesInspectCmd(client hidapi.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "inspect <path>",
		Short: "Print details for a HID device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			devices, err := client.ListDevices()
			if err != nil {
				return err
			}

			for _, device := range devices {
				if device.Path == args[0] {
					printDevice(cmd, device)
					return nil
				}
			}

			return fmt.Errorf("device not found: %s", args[0])
		},
	}
}
