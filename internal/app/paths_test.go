package app

import (
	"os/user"
	"path/filepath"
	"testing"
)

func TestDefaultPathsFromHome(t *testing.T) {
	t.Setenv("HOME", "/tmp/logi-home")

	paths := DefaultPaths()

	if got, want := paths.ConfigDir, "/tmp/logi-home/.config/logi-cli"; got != want {
		t.Fatalf("ConfigDir = %q, want %q", got, want)
	}
	if got, want := paths.ConfigFile, "/tmp/logi-home/.config/logi-cli/config.toml"; got != want {
		t.Fatalf("ConfigFile = %q, want %q", got, want)
	}
	if got, want := paths.StateDir, "/tmp/logi-home/.config/logi-cli/state"; got != want {
		t.Fatalf("StateDir = %q, want %q", got, want)
	}
	if got, want := paths.LogDir, "/tmp/logi-home/.config/logi-cli/logs"; got != want {
		t.Fatalf("LogDir = %q, want %q", got, want)
	}
	if got, want := paths.SocketFile, "/tmp/logi-home/.config/logi-cli/state/daemon.sock"; got != want {
		t.Fatalf("SocketFile = %q, want %q", got, want)
	}
	if got, want := paths.PlistFile, "/tmp/logi-home/Library/LaunchAgents/io.realtong.logi-cli.plist"; got != want {
		t.Fatalf("PlistFile = %q, want %q", got, want)
	}
}

func TestDefaultPathsFallsBackWhenHomeUnset(t *testing.T) {
	t.Setenv("HOME", "")

	current, err := user.Current()
	if err != nil {
		t.Fatalf("user.Current returned error: %v", err)
	}
	if current.HomeDir == "" {
		t.Fatal("user.Current returned empty HomeDir")
	}

	paths := DefaultPaths()

	if got, want := paths.ConfigDir, filepath.Join(current.HomeDir, ".config", "logi-cli"); got != want {
		t.Fatalf("ConfigDir = %q, want %q", got, want)
	}
	if got, want := paths.ConfigFile, filepath.Join(current.HomeDir, ".config", "logi-cli", "config.toml"); got != want {
		t.Fatalf("ConfigFile = %q, want %q", got, want)
	}
	if got, want := paths.StateDir, filepath.Join(current.HomeDir, ".config", "logi-cli", "state"); got != want {
		t.Fatalf("StateDir = %q, want %q", got, want)
	}
	if got, want := paths.LogDir, filepath.Join(current.HomeDir, ".config", "logi-cli", "logs"); got != want {
		t.Fatalf("LogDir = %q, want %q", got, want)
	}
	if got, want := paths.SocketFile, filepath.Join(current.HomeDir, ".config", "logi-cli", "state", "daemon.sock"); got != want {
		t.Fatalf("SocketFile = %q, want %q", got, want)
	}
	if got, want := paths.PlistFile, filepath.Join(current.HomeDir, "Library", "LaunchAgents", "io.realtong.logi-cli.plist"); got != want {
		t.Fatalf("PlistFile = %q, want %q", got, want)
	}
}
