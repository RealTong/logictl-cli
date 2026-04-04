package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	appcore "github.com/realtong/logictl-cli/internal/app"
	"github.com/realtong/logictl-cli/internal/daemon"
	"github.com/realtong/logictl-cli/internal/ipc"
	platformmacos "github.com/realtong/logictl-cli/internal/platform/macos"
	"github.com/spf13/cobra"
)

type daemonServiceManager interface {
	Install(context.Context) (string, error)
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
		Short: "Control the logictl daemon",
	}

	cmd.AddCommand(newDaemonInstallCmd(manager))
	cmd.AddCommand(newDaemonRunCmd(app))
	cmd.AddCommand(newDaemonStartCmd(app, manager))
	cmd.AddCommand(newDaemonStatusCmd(app))
	cmd.AddCommand(newDaemonStopCmd(manager))
	cmd.AddCommand(newDaemonRestartCmd(app, manager))
	return cmd
}

func newDaemonInstallCmd(manager daemonServiceManager) *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Install the stable LaunchAgent daemon binary",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			installed, err := manager.Install(cmd.Context())
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "installed %s\n", installed)
			return nil
		},
	}
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

func (m launchAgentServiceManager) Install(ctx context.Context) (string, error) {
	_ = ctx
	binary, err := m.executable()
	if err != nil {
		return "", err
	}
	return installLaunchAgentBinary(m.paths, binary)
}

func (m launchAgentServiceManager) Start(ctx context.Context) error {
	installedBinary, err := resolveInstalledLaunchAgentBinary(m.paths)
	if err != nil {
		return err
	}
	if err := platformmacos.StartLaunchAgent(ctx, m.paths, installedBinary); err != nil {
		return err
	}
	if err := waitForDaemonReady(m.paths.SocketFile, 2*time.Second); err != nil {
		return launchAgentStartFailure(m.paths, installedBinary, err)
	}
	return nil
}

func (m launchAgentServiceManager) Stop(ctx context.Context) error {
	return platformmacos.StopLaunchAgent(ctx, m.paths)
}

func (m launchAgentServiceManager) Restart(ctx context.Context) error {
	installedBinary, err := resolveInstalledLaunchAgentBinary(m.paths)
	if err != nil {
		return err
	}
	if err := platformmacos.RestartLaunchAgent(ctx, m.paths, installedBinary); err != nil {
		return err
	}
	if err := waitForDaemonReady(m.paths.SocketFile, 2*time.Second); err != nil {
		return launchAgentStartFailure(m.paths, installedBinary, err)
	}
	return nil
}

func installLaunchAgentBinary(paths appcore.Paths, binary string) (string, error) {
	if requiresLaunchAgentStaging(binary) {
		return "", fmt.Errorf("current executable %q looks like a go run build cache binary; build a stable binary such as ./bin/logictl with `go build -o ./bin/logictl ./cmd/logictl` and run `./bin/logictl daemon install`", binary)
	}

	if err := os.MkdirAll(paths.StateDir, 0o755); err != nil {
		return "", err
	}

	installedPath := launchAgentBinaryPath(paths)
	if filepath.Clean(binary) == filepath.Clean(installedPath) {
		return installedPath, nil
	}
	return installedPath, copyExecutable(binary, installedPath)
}

func resolveInstalledLaunchAgentBinary(paths appcore.Paths) (string, error) {
	installedPath := launchAgentBinaryPath(paths)
	info, err := os.Stat(installedPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("no installed daemon binary at %s; build a stable binary with `go build -o ./bin/logictl ./cmd/logictl`, then run `./bin/logictl daemon install`", installedPath)
		}
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("installed daemon binary path %q is a directory; rebuild `./bin/logictl` and rerun `./bin/logictl daemon install`", installedPath)
	}
	if info.Mode()&0o111 == 0 {
		return "", fmt.Errorf("installed daemon binary %q is not executable; rebuild `./bin/logictl` and rerun `./bin/logictl daemon install`", installedPath)
	}
	return installedPath, nil
}

func launchAgentBinaryPath(paths appcore.Paths) string {
	return filepath.Join(paths.StateDir, "logictl-daemon")
}

func requiresLaunchAgentStaging(binary string) bool {
	clean := filepath.Clean(binary)
	return strings.Contains(clean, string(filepath.Separator)+"go-build")
}

func copyExecutable(source string, target string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()

	info, err := input.Stat()
	if err != nil {
		return err
	}

	output, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode()&0o777)
	if err != nil {
		return err
	}

	if _, err := io.Copy(output, input); err != nil {
		_ = output.Close()
		return err
	}
	if err := output.Close(); err != nil {
		return err
	}

	return nil
}

func waitForDaemonReady(socketPath string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		status, err := ipc.QueryStatus(socketPath)
		if err == nil && (status.Running || status.Message != "") {
			return nil
		}
		lastErr = err
		time.Sleep(50 * time.Millisecond)
	}
	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("daemon socket %q did not become ready", socketPath)
}

func launchAgentStartFailure(paths appcore.Paths, stagedBinary string, cause error) error {
	logBody := readDaemonStderrLog(filepath.Join(paths.LogDir, "daemon.stderr.log"))
	lowerLog := strings.ToLower(logBody)
	if strings.Contains(logBody, "0xe00002e2") || strings.Contains(lowerLog, "not permitted") {
		return fmt.Errorf("launch agent started but exited immediately because %s lacks Input Monitoring permission; grant Input Monitoring to %s and retry", stagedBinary, stagedBinary)
	}
	if summary := summarizeDaemonStderrLog(logBody); summary != "" {
		return fmt.Errorf("launch agent started but daemon did not become ready: %w; recent stderr: %s", cause, summary)
	}
	if cause != nil && !errors.Is(cause, os.ErrNotExist) {
		return fmt.Errorf("launch agent started but daemon did not become ready: %w", cause)
	}
	return fmt.Errorf("launch agent started but daemon did not become ready")
}

func readDaemonStderrLog(path string) string {
	data, err := os.ReadFile(path)
	if err != nil || len(data) == 0 {
		return ""
	}
	return string(data)
}

func summarizeDaemonStderrLog(body string) string {
	lines := strings.Split(strings.TrimSpace(body), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" && !strings.HasPrefix(line, "Usage:") && !strings.HasPrefix(line, "Flags:") && !strings.Contains(line, "help for") {
			return line
		}
	}
	return ""
}
