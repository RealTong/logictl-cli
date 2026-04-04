package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/realtong/logictl-cli/internal/config"
	"github.com/realtong/logictl-cli/internal/daemon"
)

func TestSmokeConfigParsesAndBuildsRuntime(t *testing.T) {
	cfg, err := config.LoadFile(filepath.Join("..", "..", "examples", "config.toml"))
	if err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}
	if _, err := daemon.NewFromConfig(cfg); err != nil {
		t.Fatalf("NewFromConfig returned error: %v", err)
	}
}
