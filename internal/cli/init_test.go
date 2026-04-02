package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/realtong/logi-cli/internal/app"
)

func TestInitCmdCreatesStarterConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"init"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	paths := app.DefaultPaths()
	if _, err := os.Stat(paths.ConfigFile); err != nil {
		t.Fatalf("starter config missing: %v", err)
	}

	data, err := os.ReadFile(paths.ConfigFile)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("starter config is empty")
	}

	if _, err := os.Stat(filepath.Join(paths.ConfigDir)); err != nil {
		t.Fatalf("config dir missing: %v", err)
	}
}
