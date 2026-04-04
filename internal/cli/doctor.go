package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/realtong/logictl-cli/internal/app"
	"github.com/realtong/logictl-cli/internal/config"
	"github.com/realtong/logictl-cli/internal/ipc"
	platformmacos "github.com/realtong/logictl-cli/internal/platform/macos"
	"github.com/spf13/cobra"
)

type platformDoctor interface {
	ActiveBundleID(context.Context) (string, error)
	Permissions(context.Context) (platformmacos.PermissionReport, error)
}

type daemonStatusReporter interface {
	Status() (ipc.Status, error)
}

func newDoctorCmd(doctor platformDoctor, daemonReporter daemonStatusReporter, defaultConfigPath string) *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Inspect macOS permissions, config health, and daemon status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			permissions, err := doctor.Permissions(cmd.Context())
			if err != nil {
				return err
			}

			activeBundleID, err := doctor.ActiveBundleID(cmd.Context())
			if err != nil {
				return err
			}

			status, err := daemonReporter.Status()
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Accessibility: %s\n", permissionLabel(permissions.AccessibilityGranted))
			fmt.Fprintf(cmd.OutOrStdout(), "Input Monitoring: %s\n", permissionLabel(permissions.InputMonitoringGranted))
			if activeBundleID == "" {
				activeBundleID = "unavailable"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Frontmost App: %s\n", activeBundleID)
			fmt.Fprintf(cmd.OutOrStdout(), "Config: %s\n", configStatus(configPath))
			fmt.Fprintf(cmd.OutOrStdout(), "Daemon: %s\n", daemonStatusLabel(status))
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", defaultConfigPath, "config file to inspect")
	return cmd
}

func permissionLabel(granted bool) string {
	if granted {
		return "granted"
	}
	return "missing"
}

func configStatus(path string) string {
	cfg, err := config.LoadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "missing"
		}
		return fmt.Sprintf("invalid (%v)", err)
	}
	if err := config.Validate(cfg); err != nil {
		return fmt.Sprintf("invalid (%v)", err)
	}
	return "valid"
}

func daemonStatusLabel(status ipc.Status) string {
	if status.Message != "" {
		return status.Message
	}
	if status.Running {
		return "running"
	}
	return "stopped"
}

func defaultDoctorCmd(daemonReporter daemonStatusReporter) *cobra.Command {
	return newDoctorCmd(platformmacos.NewEnvironment(), daemonReporter, app.DefaultPaths().ConfigFile)
}
