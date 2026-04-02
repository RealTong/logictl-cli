package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/realtong/logi-cli/internal/app"
)

func TestInitCmdCreatesStarterConfigWithoutRepoLookup(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cwd := t.TempDir()
	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	defer func() { _ = os.Chdir(original) }()
	if err := os.Chdir(cwd); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"init"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	paths := app.DefaultPaths()
	data, err := os.ReadFile(paths.ConfigFile)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("starter config is empty")
	}

	if _, err := os.Stat(filepath.Dir(paths.ConfigFile)); err != nil {
		t.Fatalf("config dir missing: %v", err)
	}
}
