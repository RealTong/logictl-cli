package app

import "testing"

func TestDefaultPathsFromHome(t *testing.T) {
	t.Setenv("HOME", "/tmp/logi-home")

	paths := DefaultPaths()

	if got, want := paths.ConfigFile, "/tmp/logi-home/.config/logi-cli/config.toml"; got != want {
		t.Fatalf("ConfigFile = %q, want %q", got, want)
	}
}
