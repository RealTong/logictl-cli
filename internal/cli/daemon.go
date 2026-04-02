package cli

import (
	"github.com/realtong/logi-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func newDaemonCmd(app *daemon.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Control the logi-cli daemon",
	}

	cmd.AddCommand(newDaemonRunCmd(app))
	cmd.AddCommand(newDaemonStatusCmd(app))
	return cmd
}

func newDaemonRunCmd(app *daemon.App) *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run the daemon in the foreground",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Run(cmd.Context())
		},
	}
}

func newDaemonStatusCmd(app *daemon.App) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Query daemon status over the local control socket",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := app.Status()
			if err != nil {
				return err
			}

			if status.Message != "" {
				cmd.Println(status.Message)
				return nil
			}
			if status.Running {
				cmd.Println("running")
				return nil
			}

			cmd.Println("stopped")
			return nil
		},
	}
}
