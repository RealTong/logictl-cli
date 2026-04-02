package cli

import (
	"context"
	"os"

	appcore "github.com/realtong/logi-cli/internal/app"
	"github.com/realtong/logi-cli/internal/daemon"
	platformmacos "github.com/realtong/logi-cli/internal/platform/macos"
	"github.com/spf13/cobra"
)

type daemonServiceManager interface {
	Start(context.Context) error
	Stop(context.Context) error
	Restart(context.Context) error
}

type daemonPreflighter interface {
	Preflight() error
}

type launchAgentServiceManager struct {
	paths      appcore.Paths
	executable func() (string, error)
}

func newDaemonCmd(app *daemon.App) *cobra.Command {
	return newDaemonCmdWithServiceManager(app, launchAgentServiceManager{
		paths:      appcore.DefaultPaths(),
		executable: os.Executable,
	})
}

func newDaemonCmdWithServiceManager(app *daemon.App, manager daemonServiceManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Control the logi-cli daemon",
	}

	cmd.AddCommand(newDaemonRunCmd(app))
	cmd.AddCommand(newDaemonStartCmd(app, manager))
	cmd.AddCommand(newDaemonStatusCmd(app))
	cmd.AddCommand(newDaemonStopCmd(manager))
	cmd.AddCommand(newDaemonRestartCmd(app, manager))
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

func newDaemonStartCmd(app daemonPreflighter, manager daemonServiceManager) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Install and start the LaunchAgent-managed daemon",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.Preflight(); err != nil {
				return err
			}
			if err := manager.Start(cmd.Context()); err != nil {
				return err
			}
			cmd.Println("started")
			return nil
		},
	}
}

func newDaemonStopCmd(manager daemonServiceManager) *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the LaunchAgent-managed daemon",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := manager.Stop(cmd.Context()); err != nil {
				return err
			}
			cmd.Println("stopped")
			return nil
		},
	}
}

func newDaemonRestartCmd(app daemonPreflighter, manager daemonServiceManager) *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Restart the LaunchAgent-managed daemon",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.Preflight(); err != nil {
				return err
			}
			if err := manager.Restart(cmd.Context()); err != nil {
				return err
			}
			cmd.Println("restarted")
			return nil
		},
	}
}

func (m launchAgentServiceManager) Start(ctx context.Context) error {
	binary, err := m.executable()
	if err != nil {
		return err
	}
	return platformmacos.StartLaunchAgent(ctx, m.paths, binary)
}

func (m launchAgentServiceManager) Stop(ctx context.Context) error {
	return platformmacos.StopLaunchAgent(ctx, m.paths)
}

func (m launchAgentServiceManager) Restart(ctx context.Context) error {
	binary, err := m.executable()
	if err != nil {
		return err
	}
	return platformmacos.RestartLaunchAgent(ctx, m.paths, binary)
}
