package app

import (
	"os/user"
	"path/filepath"
	"testing"
)

func TestDefaultPathsFromHome(t *testing.T) {
	t.Setenv("HOME", "/tmp/logictl-home")

	paths := DefaultPaths()

	if got, want := paths.ConfigDir, "/tmp/logictl-home/.config/logictl"; got != want {
		t.Fatalf("ConfigDir = %q, want %q", got, want)
	}
	if got, want := paths.ConfigFile, "/tmp/logictl-home/.config/logictl/config.toml"; got != want {
		t.Fatalf("ConfigFile = %q, want %q", got, want)
	}
	if got, want := paths.StateDir, "/tmp/logictl-home/.config/logictl/state"; got != want {
		t.Fatalf("StateDir = %q, want %q", got, want)
	}
	if got, want := paths.LogDir, "/tmp/logictl-home/.config/logictl/logs"; got != want {
		t.Fatalf("LogDir = %q, want %q", got, want)
	}
	if got, want := paths.SocketFile, "/tmp/logictl-home/.config/logictl/state/daemon.sock"; got != want {
		t.Fatalf("SocketFile = %q, want %q", got, want)
	}
	if got, want := paths.PlistFile, "/tmp/logictl-home/Library/LaunchAgents/io.realtong.logictl.plist"; got != want {
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

	if got, want := paths.ConfigDir, filepath.Join(current.HomeDir, ".config", "logictl"); got != want {
		t.Fatalf("ConfigDir = %q, want %q", got, want)
	}
	if got, want := paths.ConfigFile, filepath.Join(current.HomeDir, ".config", "logictl", "config.toml"); got != want {
		t.Fatalf("ConfigFile = %q, want %q", got, want)
	}
	if got, want := paths.StateDir, filepath.Join(current.HomeDir, ".config", "logictl", "state"); got != want {
		t.Fatalf("StateDir = %q, want %q", got, want)
	}
	if got, want := paths.LogDir, filepath.Join(current.HomeDir, ".config", "logictl", "logs"); got != want {
		t.Fatalf("LogDir = %q, want %q", got, want)
	}
	if got, want := paths.SocketFile, filepath.Join(current.HomeDir, ".config", "logictl", "state", "daemon.sock"); got != want {
		t.Fatalf("SocketFile = %q, want %q", got, want)
	}
	if got, want := paths.PlistFile, filepath.Join(current.HomeDir, "Library", "LaunchAgents", "io.realtong.logictl.plist"); got != want {
		t.Fatalf("PlistFile = %q, want %q", got, want)
	}
}
