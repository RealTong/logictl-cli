package cli

import (
	"github.com/realtong/logi-cli/internal/hidapi"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	return newRootCmd(hidapi.NewClient())
}

func newRootCmd(hidClient hidapi.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logi",
		Short: "Configure Logitech device behavior on macOS",
	}

	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.AddCommand(newDevicesCmd(hidClient))
	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newValidateCmd())
	cmd.AddCommand(newVersionCmd())
	return cmd
}
