package macos

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	appcore "github.com/realtong/logi-cli/internal/app"
)

func TestInstallLaunchAgentWritesExpectedPlist(t *testing.T) {
	paths := launchAgentTestPaths(t)
	binary := "/tmp/logi-cli"

	if err := InstallLaunchAgent(paths, binary); err != nil {
		t.Fatalf("InstallLaunchAgent returned error: %v", err)
	}

	data, err := os.ReadFile(paths.PlistFile)
	if err != nil {
		t.Fatalf("ReadFile(%q) returned error: %v", paths.PlistFile, err)
	}

	got := string(data)
	for _, want := range []string{
		binary,
		"<string>daemon</string>",
		"<string>run</string>",
		paths.ConfigFile,
		paths.SocketFile,
		filepath.Join(paths.LogDir, "daemon.stdout.log"),
		filepath.Join(paths.LogDir, "daemon.stderr.log"),
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("plist = %q, want %q", got, want)
		}
	}
}

func TestRestartLaunchAgentRunsBootoutThenBootstrap(t *testing.T) {
	paths := launchAgentTestPaths(t)
	binary := "/tmp/logi-cli"

	previousRun := runLaunchctl
	previousUID := currentLaunchctlUID
	t.Cleanup(func() {
		runLaunchctl = previousRun
		currentLaunchctlUID = previousUID
	})

	var calls []string
	runLaunchctl = func(_ context.Context, args ...string) error {
		calls = append(calls, strings.Join(args, " "))
		return nil
	}
	currentLaunchctlUID = func() (string, error) {
		return "501", nil
	}

	if err := RestartLaunchAgent(context.Background(), paths, binary); err != nil {
		t.Fatalf("RestartLaunchAgent returned error: %v", err)
	}

	if len(calls) != 2 {
		t.Fatalf("len(calls) = %d, want 2", len(calls))
	}
	if got, want := calls[0], "bootout gui/501 "+paths.PlistFile; got != want {
		t.Fatalf("calls[0] = %q, want %q", got, want)
	}
	if got, want := calls[1], "bootstrap gui/501 "+paths.PlistFile; got != want {
		t.Fatalf("calls[1] = %q, want %q", got, want)
	}
}

func launchAgentTestPaths(t *testing.T) appcore.Paths {
	t.Helper()

	base := t.TempDir()
	return appcore.Paths{
		ConfigDir:  filepath.Join(base, "config"),
		ConfigFile: filepath.Join(base, "config", "config.toml"),
		StateDir:   filepath.Join(base, "state"),
		LogDir:     filepath.Join(base, "logs"),
		SocketFile: filepath.Join(base, "state", "daemon.sock"),
		PlistFile:  filepath.Join(base, "LaunchAgents", "io.realtong.logi-cli.plist"),
	}
}
