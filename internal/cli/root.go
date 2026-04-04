package cli

import (
	appcore "github.com/realtong/logictl-cli/internal/app"
	"github.com/realtong/logictl-cli/internal/daemon"
	"github.com/realtong/logictl-cli/internal/events"
	"github.com/realtong/logictl-cli/internal/hidapi"
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
		Use:   "logictl",
		Short: "Configure Logitech device behavior on macOS",
	}

	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.AddCommand(newDaemonCmd(daemonApp))
	cmd.AddCommand(defaultDoctorCmd(daemonApp))
	cmd.AddCommand(newDevicesCmd(hidClient))
	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newReloadCmd(daemonApp))
	cmd.AddCommand(newTestCmd(hidClient, func(path string) rawSource {
		return events.NewHIDSource(path)
	}))
	cmd.AddCommand(newValidateCmd())
	cmd.AddCommand(newVersionCmd())
	return cmd
}
