package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateRejectsAmbiguousBindingsFixture(t *testing.T) {
	cfg, err := LoadFile(filepath.Join("..", "..", "testdata", "config", "duplicate_binding.toml"))
	if err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}

	err = Validate(cfg)
	if err == nil {
		t.Fatal("Validate returned nil, want ambiguous binding error")
	}
	if !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("Validate error = %v, want ambiguous binding error", err)
	}
}

func TestValidateRejectsMissingActionFixture(t *testing.T) {
	cfg, err := LoadFile(filepath.Join("..", "..", "testdata", "config", "missing_action.toml"))
	if err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}

	if err := Validate(cfg); err == nil {
		t.Fatal("Validate returned nil, want missing action error")
	}
}
