package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func validConfigForValidationTests() *Config {
	return &Config{
		Devices: []Device{
			{
				ID:             "mx-master-4",
				MatchVendorID:  1133,
				MatchProductID: 1234,
				Capabilities: map[string]string{
					"thumb_button": "button_5",
				},
			},
		},
		Actions: []Action{
			{
				ID:   "close_tab",
				Type: "shortcut",
				Keys: []string{"cmd", "w"},
			},
		},
		Profiles: []Profile{
			{
				ID:          "chrome",
				AppBundleID: "com.google.Chrome",
				Bindings: []Binding{
					{
						Device:  "mx-master-4",
						Trigger: "hold(thumb_button)+move(down)",
						Action:  "close_tab",
					},
				},
			},
		},
	}
}

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

func TestValidateRejectsDuplicateActionIDs(t *testing.T) {
	cfg := validConfigForValidationTests()
	cfg.Actions = append(cfg.Actions, Action{
		ID:     "close_tab",
		Type:   "system",
		System: "mission_control",
	})

	if err := Validate(cfg); err == nil {
		t.Fatal("Validate returned nil, want duplicate action id error")
	}
}

func TestValidateRejectsDuplicateDeviceIDs(t *testing.T) {
	cfg := validConfigForValidationTests()
	cfg.Devices = append(cfg.Devices, Device{
		ID:             "mx-master-4",
		MatchVendorID:  1133,
		MatchProductID: 1234,
	})

	if err := Validate(cfg); err == nil {
		t.Fatal("Validate returned nil, want duplicate device id error")
	}
}

func TestValidateRejectsUnknownDeviceReferences(t *testing.T) {
	cfg := validConfigForValidationTests()
	cfg.Profiles[0].Bindings[0].Device = "missing-device"

	if err := Validate(cfg); err == nil {
		t.Fatal("Validate returned nil, want unknown device reference error")
	}
}

func TestValidateRejectsInvalidActionPayloads(t *testing.T) {
	tests := []struct {
		name   string
		action Action
	}{
		{
			name: "shortcut without keys",
			action: Action{
				ID:   "close_tab",
				Type: "shortcut",
			},
		},
		{
			name: "system without system action",
			action: Action{
				ID:   "close_tab",
				Type: "system",
			},
		},
		{
			name: "script without script path",
			action: Action{
				ID:   "close_tab",
				Type: "script",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := validConfigForValidationTests()
			cfg.Actions[0] = tc.action

			if err := Validate(cfg); err == nil {
				t.Fatal("Validate returned nil, want invalid action payload error")
			}
		})
	}
}
