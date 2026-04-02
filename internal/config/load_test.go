package config

import (
	"path/filepath"
	"testing"
)

func TestLoadFileValidFixture(t *testing.T) {
	cfg, err := LoadFile(filepath.Join("..", "..", "testdata", "config", "valid.toml"))
	if err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}

	if got, want := len(cfg.Actions), 1; got != want {
		t.Fatalf("len(cfg.Actions) = %d, want %d", got, want)
	}
}
