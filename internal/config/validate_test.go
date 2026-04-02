package config

import (
	"path/filepath"
	"testing"
)

func TestValidateRejectsDuplicateBindingFixture(t *testing.T) {
	cfg, err := LoadFile(filepath.Join("..", "..", "testdata", "config", "duplicate_binding.toml"))
	if err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}

	if err := Validate(cfg); err == nil {
		t.Fatal("Validate returned nil, want duplicate binding error")
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
