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

	appcore "github.com/realtong/logictl-cli/internal/app"
	"github.com/realtong/logictl-cli/internal/daemon"
	"github.com/realtong/logictl-cli/internal/ipc"
)

type fakeDaemonPreflight struct {
	err error
}

func (f fakeDaemonPreflight) Preflight() error {
	return f.err
}

type fakeDaemonServiceManager struct {
	calls []string
}

func (m *fakeDaemonServiceManager) Install(context.Context) (string, error) {
	m.calls = append(m.calls, "install")
	return "/tmp/logictl-daemon", nil
}

func (m *fakeDaemonServiceManager) Start(context.Context) error {
	m.calls = append(m.calls, "start")
	return nil
}

func (m *fakeDaemonServiceManager) Stop(context.Context) error {
	m.calls = append(m.calls, "stop")
	return nil
}

func (m *fakeDaemonServiceManager) Restart(context.Context) error {
	m.calls = append(m.calls, "restart")
	return nil
}

func TestDaemonStartCmdInvokesServiceManager(t *testing.T) {
	manager := &fakeDaemonServiceManager{}
	cmd := newDaemonStartCmd(fakeDaemonPreflight{}, manager)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(manager.calls) != 1 || manager.calls[0] != "start" {
		t.Fatalf("manager.calls = %#v, want [start]", manager.calls)
	}
}

func TestDaemonInstallCmdInvokesServiceManager(t *testing.T) {
	manager := &fakeDaemonServiceManager{}
	cmd := newDaemonInstallCmd(manager)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(manager.calls) != 1 || manager.calls[0] != "install" {
		t.Fatalf("manager.calls = %#v, want [install]", manager.calls)
	}
	if got := buf.String(); !strings.Contains(got, "/tmp/logictl-daemon") {
		t.Fatalf("output = %q, want installed path", got)
	}
}

func TestDaemonRestartCmdInvokesServiceManager(t *testing.T) {
	manager := &fakeDaemonServiceManager{}
	cmd := newDaemonRestartCmd(fakeDaemonPreflight{}, manager)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(manager.calls) != 1 || manager.calls[0] != "restart" {
		t.Fatalf("manager.calls = %#v, want [restart]", manager.calls)
	}
}

func TestDaemonStartCmdRejectsPreflightFailures(t *testing.T) {
	manager := &fakeDaemonServiceManager{}
	cmd := newDaemonStartCmd(fakeDaemonPreflight{err: errors.New("unsafe path")}, manager)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil, want preflight error")
	}
	if len(manager.calls) != 0 {
		t.Fatalf("manager.calls = %#v, want no service-manager calls on preflight failure", manager.calls)
	}
}

func TestStageLaunchAgentBinaryCopiesExecutableIntoStableStatePath(t *testing.T) {
	root := t.TempDir()
	sourceDir := filepath.Join(root, "go-build123", "b001", "exe")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) returned error: %v", sourceDir, err)
	}

	source := filepath.Join(sourceDir, "logictl")
	if err := os.WriteFile(source, []byte("binary"), 0o755); err != nil {
		t.Fatalf("WriteFile(%q) returned error: %v", source, err)
	}

	paths := appcore.Paths{
		StateDir: filepath.Join(root, "state"),
	}

	installed, err := installLaunchAgentBinary(paths, source)
	if err == nil {
		t.Fatalf("installLaunchAgentBinary() error = nil, want go run rejection for %q", source)
	}
	if !strings.Contains(err.Error(), "build a stable binary") {
		t.Fatalf("installLaunchAgentBinary() error = %q, want stable-binary guidance", err)
	}
	if installed != "" {
		t.Fatalf("installLaunchAgentBinary() = %q, want empty installed path on rejection", installed)
	}
}

func TestInstallLaunchAgentBinaryCopiesStableExecutableIntoInstalledPath(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "bin", "logictl")
	if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) returned error: %v", filepath.Dir(source), err)
	}
	if err := os.WriteFile(source, []byte("binary"), 0o755); err != nil {
		t.Fatalf("WriteFile(%q) returned error: %v", source, err)
	}

	paths := appcore.Paths{
		StateDir: filepath.Join(root, "state"),
	}

	installed, err := installLaunchAgentBinary(paths, source)
	if err != nil {
		t.Fatalf("installLaunchAgentBinary() returned error: %v", err)
	}
	if installed == source {
		t.Fatalf("installLaunchAgentBinary() = %q, want copied installed path distinct from source", installed)
	}

	got, err := os.ReadFile(installed)
	if err != nil {
		t.Fatalf("ReadFile(%q) returned error: %v", installed, err)
	}
	if string(got) != "binary" {
		t.Fatalf("installed file contents = %q, want %q", string(got), "binary")
	}
}

func TestResolveInstalledLaunchAgentBinaryReturnsInstalledBinaryPath(t *testing.T) {
	root := t.TempDir()
	paths := appcore.Paths{
		StateDir: filepath.Join(root, "state"),
	}
	installed := filepath.Join(paths.StateDir, "logictl-daemon")
	if err := os.MkdirAll(paths.StateDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) returned error: %v", paths.StateDir, err)
	}
	if err := os.WriteFile(installed, []byte("binary"), 0o755); err != nil {
		t.Fatalf("WriteFile(%q) returned error: %v", installed, err)
	}

	got, err := resolveInstalledLaunchAgentBinary(paths)
	if err != nil {
		t.Fatalf("resolveInstalledLaunchAgentBinary() returned error: %v", err)
	}
	if got != installed {
		t.Fatalf("resolveInstalledLaunchAgentBinary() = %q, want %q", got, installed)
	}
}

func TestResolveInstalledLaunchAgentBinaryRequiresInstallStep(t *testing.T) {
	root := t.TempDir()
	paths := appcore.Paths{
		StateDir: filepath.Join(root, "state"),
	}

	got, err := resolveInstalledLaunchAgentBinary(paths)
	if err == nil {
		t.Fatal("resolveInstalledLaunchAgentBinary() error = nil, want install guidance")
	}
	if !strings.Contains(err.Error(), "daemon install") {
		t.Fatalf("resolveInstalledLaunchAgentBinary() error = %q, want install guidance", err)
	}
	if got != "" {
		t.Fatalf("resolveInstalledLaunchAgentBinary() = %q, want empty path on missing install", got)
	}
}

func TestWaitForDaemonReadyReportsRunningSocket(t *testing.T) {
	socketPath := filepath.Join(os.TempDir(), fmt.Sprintf("logictl-ready-%d.sock", time.Now().UnixNano()))
	t.Cleanup(func() {
		_ = os.Remove(socketPath)
	})
	server := daemon.NewServer(socketPath, ipc.Status{Running: true, Message: "running"})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = server.Run(ctx)
	}()

	if err := waitForDaemonReady(socketPath, 2*time.Second); err != nil {
		t.Fatalf("waitForDaemonReady() returned error: %v", err)
	}
}

func TestLaunchAgentStartFailureUsesPermissionGuidanceFromDaemonLog(t *testing.T) {
	root := t.TempDir()
	paths := appcore.Paths{
		LogDir:   filepath.Join(root, "logs"),
		StateDir: filepath.Join(root, "state"),
	}
	if err := os.MkdirAll(paths.LogDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) returned error: %v", paths.LogDir, err)
	}

	stagedBinary := filepath.Join(paths.StateDir, "logictl-daemon")
	if err := os.MkdirAll(paths.StateDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) returned error: %v", paths.StateDir, err)
	}
	if err := os.WriteFile(filepath.Join(paths.LogDir, "daemon.stderr.log"), []byte("Error: IOHIDManagerOpen failed for MX Master 4: 0xe00002e2\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(stderr log) returned error: %v", err)
	}

	err := launchAgentStartFailure(paths, stagedBinary, fmt.Errorf("daemon socket never became ready"))
	if err == nil {
		t.Fatal("launchAgentStartFailure() returned nil, want permission guidance")
	}
	if got := err.Error(); !strings.Contains(got, stagedBinary) || !strings.Contains(got, "Input Monitoring") {
		t.Fatalf("launchAgentStartFailure() = %q, want Input Monitoring guidance for %q", got, stagedBinary)
	}
}

func TestLaunchAgentStartFailureFindsPermissionErrorBeforeCobraUsageTail(t *testing.T) {
	root := t.TempDir()
	paths := appcore.Paths{
		LogDir:   filepath.Join(root, "logs"),
		StateDir: filepath.Join(root, "state"),
	}
	if err := os.MkdirAll(paths.LogDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) returned error: %v", paths.LogDir, err)
	}
	if err := os.MkdirAll(paths.StateDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) returned error: %v", paths.StateDir, err)
	}

	stagedBinary := filepath.Join(paths.StateDir, "logictl-daemon")
	logBody := "Error: IOHIDManagerOpen failed for MX Master 4: 0xe00002e2\nUsage:\n  logictl daemon run [flags]\n\nFlags:\n  -h, --help   help for run\n"
	if err := os.WriteFile(filepath.Join(paths.LogDir, "daemon.stderr.log"), []byte(logBody), 0o644); err != nil {
		t.Fatalf("WriteFile(stderr log) returned error: %v", err)
	}

	err := launchAgentStartFailure(paths, stagedBinary, fmt.Errorf("daemon socket never became ready"))
	if err == nil {
		t.Fatal("launchAgentStartFailure() returned nil, want permission guidance")
	}
	if got := err.Error(); !strings.Contains(got, "Input Monitoring") {
		t.Fatalf("launchAgentStartFailure() = %q, want Input Monitoring guidance", got)
	}
}
