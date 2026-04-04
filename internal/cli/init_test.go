package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/realtong/logictl-cli/internal/app"
)

func TestStarterConfigContentIsNonEmpty(t *testing.T) {
	if got := starterConfigContent(); len(got) == 0 {
		t.Fatal("starterConfigContent returned empty content")
	} else if !strings.Contains(string(got), "[[devices]]") {
		t.Fatalf("starterConfigContent = %q, want embedded example config", string(got))
	}
}

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

func TestInitCmdRefusesToOverwriteExistingConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	paths := app.DefaultPaths()
	if err := os.MkdirAll(paths.ConfigDir, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}

	const original = "custom config\n"
	if err := os.WriteFile(paths.ConfigFile, []byte(original), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"init"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil, want existing config error")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("Execute error = %v, want existing config error", err)
	}

	data, readErr := os.ReadFile(paths.ConfigFile)
	if readErr != nil {
		t.Fatalf("ReadFile returned error: %v", readErr)
	}
	if got := string(data); got != original {
		t.Fatalf("config file was modified: got %q, want %q", got, original)
	}
}
