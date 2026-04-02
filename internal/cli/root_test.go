package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/realtong/logi-cli/internal/app"
	"github.com/realtong/logi-cli/internal/daemon"
	"github.com/realtong/logi-cli/internal/hidapi"
	"github.com/realtong/logi-cli/internal/ipc"
)

func TestNewRootCmdHelpHidesCompletion(t *testing.T) {
	cmd := NewRootCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "completion") {
		t.Fatalf("help output exposes completion command: %s", out)
	}
	if !strings.Contains(out, "version") {
		t.Fatalf("help output missing version command: %s", out)
	}
}

func TestDaemonStatusCmdReportsRunningDaemon(t *testing.T) {
	daemonApp := daemon.NewApp(testPaths(t))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = daemonApp.Run(ctx)
	}()

	waitForSocket(t, daemonApp.SocketPath())

	cmd := newRootCmdWithDaemon(hidapi.FakeClient{}, daemonApp)
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"daemon", "status"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "running") {
		t.Fatalf("status output = %q, want running", out)
	}
}

func TestDaemonStatusCmdReportsStoppedWhenSocketMissing(t *testing.T) {
	daemonApp := daemon.NewApp(testPaths(t))

	cmd := newRootCmdWithDaemon(hidapi.FakeClient{}, daemonApp)
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"daemon", "status"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "stopped") {
		t.Fatalf("status output = %q, want stopped", out)
	}
}

func TestReloadCmdRequestsReload(t *testing.T) {
	daemonApp := daemon.NewApp(testPaths(t))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = daemonApp.Run(ctx)
	}()

	waitForSocket(t, daemonApp.SocketPath())

	cmd := newRootCmdWithDaemon(hidapi.FakeClient{}, daemonApp)
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"reload"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "reload requested") {
		t.Fatalf("reload output = %q, want reload requested", out)
	}
}

func testPaths(t *testing.T) app.Paths {
	t.Helper()

	base := t.TempDir()
	socketPath := filepath.Join(os.TempDir(), fmt.Sprintf("logi-cli-%d.sock", time.Now().UnixNano()))
	t.Cleanup(func() {
		_ = os.Remove(socketPath)
	})

	return app.Paths{
		ConfigDir:  filepath.Join(base, "config"),
		ConfigFile: filepath.Join(base, "config", "config.toml"),
		StateDir:   filepath.Join(base, "state"),
		LogDir:     filepath.Join(base, "logs"),
		SocketFile: socketPath,
		PlistFile:  filepath.Join(base, "LaunchAgents", "io.realtong.logi-cli.plist"),
	}
}

func waitForSocket(t *testing.T, socketPath string) {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := ipc.QueryStatus(socketPath); err == nil {
			return
		} else if !errors.Is(err, os.ErrNotExist) {
			var opErr *os.PathError
			if !errors.As(err, &opErr) {
				time.Sleep(10 * time.Millisecond)
				continue
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("socket %q did not become ready", socketPath)
}
