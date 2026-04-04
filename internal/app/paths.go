package app

import (
	"os"
	"os/user"
	"path/filepath"
)

type Paths struct {
	ConfigDir  string
	ConfigFile string
	StateDir   string
	LogDir     string
	SocketFile string
	PlistFile  string
}

func DefaultPaths() Paths {
	home := homeDir()
	base := filepath.Join(home, ".config", "logictl")

	return Paths{
		ConfigDir:  base,
		ConfigFile: filepath.Join(base, "config.toml"),
		StateDir:   filepath.Join(base, "state"),
		LogDir:     filepath.Join(base, "logs"),
		SocketFile: filepath.Join(base, "state", "daemon.sock"),
		PlistFile:  filepath.Join(home, "Library", "LaunchAgents", "io.realtong.logictl.plist"),
	}
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	current, err := user.Current()
	if err == nil && current.HomeDir != "" {
		return current.HomeDir
	}

	return string(filepath.Separator)
}
