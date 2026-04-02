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

	if got, want := len(cfg.Devices), 1; got != want {
		t.Fatalf("len(cfg.Devices) = %d, want %d", got, want)
	}

	device := cfg.Devices[0]
	if got, want := device.ID, "mx-master-4"; got != want {
		t.Fatalf("device.ID = %q, want %q", got, want)
	}
	if got, want := device.MatchVendorID, 1133; got != want {
		t.Fatalf("device.MatchVendorID = %d, want %d", got, want)
	}
	if got, want := device.Capabilities["thumb_button"], "button_5"; got != want {
		t.Fatalf("device.Capabilities[thumb_button] = %q, want %q", got, want)
	}

	if got, want := len(cfg.Actions), 3; got != want {
		t.Fatalf("len(cfg.Actions) = %d, want %d", got, want)
	}

	if got, want := cfg.Profiles[0].ID, "chrome"; got != want {
		t.Fatalf("profiles[0].ID = %q, want %q", got, want)
	}
	if got, want := cfg.Profiles[0].Bindings[0].Device, "mx-master-4"; got != want {
		t.Fatalf("bindings[0].Device = %q, want %q", got, want)
	}
}

func TestLoadFileRejectsUnknownKeys(t *testing.T) {
	if _, err := LoadFile(filepath.Join("..", "..", "testdata", "config", "unknown_key.toml")); err == nil {
		t.Fatal("LoadFile returned nil, want unknown key error")
	}
}
