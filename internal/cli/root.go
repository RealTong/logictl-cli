package cli

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logi",
		Short: "Configure Logitech device behavior on macOS",
	}

	cmd.AddCommand(newVersionCmd())
	return cmd
}
