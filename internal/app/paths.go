package app

import (
	"os"
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
	home := os.Getenv("HOME")
	base := filepath.Join(home, ".config", "logi-cli")

	return Paths{
		ConfigDir:  base,
		ConfigFile: filepath.Join(base, "config.toml"),
		StateDir:   filepath.Join(base, "state"),
		LogDir:     filepath.Join(base, "logs"),
		SocketFile: filepath.Join(base, "state", "daemon.sock"),
		PlistFile:  filepath.Join(home, "Library", "LaunchAgents", "io.realtong.logi-cli.plist"),
	}
}
