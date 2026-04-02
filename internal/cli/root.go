package cli

import (
	appcore "github.com/realtong/logi-cli/internal/app"
	"github.com/realtong/logi-cli/internal/daemon"
	"github.com/realtong/logi-cli/internal/hidapi"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	return newRootCmd(hidapi.NewClient())
}

func newRootCmd(hidClient hidapi.Client) *cobra.Command {
	return newRootCmdWithDaemon(hidClient, daemon.NewApp(appcore.DefaultPaths()))
}

func newRootCmdWithDaemon(hidClient hidapi.Client, daemonApp *daemon.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logi",
		Short: "Configure Logitech device behavior on macOS",
	}

	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.AddCommand(newDaemonCmd(daemonApp))
	cmd.AddCommand(newDevicesCmd(hidClient))
	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newReloadCmd(daemonApp))
	cmd.AddCommand(newValidateCmd())
	cmd.AddCommand(newVersionCmd())
	return cmd
}
