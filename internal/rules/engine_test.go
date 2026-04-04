package rules

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/realtong/logictl-cli/internal/config"
	"github.com/realtong/logictl-cli/internal/events"
)

func TestMatchPrefersDeviceAndAppSpecificBinding(t *testing.T) {
	engine := NewEngine(sampleConfig())
	event := events.DeviceEvent{
		DeviceID: "mx-master-4",
		Gesture:  "hold(gesture_button)+move(down)",
	}

	action, err := engine.Match(Context{AppBundleID: "com.google.Chrome"}, event)
	if err != nil {
		t.Fatalf("Match returned error: %v", err)
	}
	if action.ID != "close_tab" {
		t.Fatalf("action.ID = %q, want close_tab", action.ID)
	}
}

func TestMatchPrefersDeviceGlobalBindingOverAnyDeviceAppBinding(t *testing.T) {
	cfg := sampleConfig()
	cfg.Profiles = []config.Profile{
		{
			ID: "global",
			Bindings: []config.Binding{
				{
					Device:  "mx-master-4",
					Trigger: "hold(gesture_button)+move(down)",
					Action:  "device_global",
				},
			},
		},
		{
			ID:          "chrome",
			AppBundleID: "com.google.Chrome",
			Bindings: []config.Binding{
				{
					Trigger: "hold(gesture_button)+move(down)",
					Action:  "app_specific",
				},
			},
		},
	}
	cfg.Actions = []config.Action{
		{ID: "device_global", Type: "system", System: "mission_control"},
		{ID: "app_specific", Type: "system", System: "launchpad"},
	}

	engine := NewEngine(cfg)
	action, err := engine.Match(Context{AppBundleID: "com.google.Chrome"}, events.DeviceEvent{
		DeviceID: "mx-master-4",
		Gesture:  "hold(gesture_button)+move(down)",
	})
	if err != nil {
		t.Fatalf("Match returned error: %v", err)
	}
	if action.ID != "device_global" {
		t.Fatalf("action.ID = %q, want device_global", action.ID)
	}
}

func TestMatchUsesHigherPriorityWithinSameSpecificity(t *testing.T) {
	low := 10
	high := 50
	cfg := sampleConfig()
	cfg.Profiles = []config.Profile{
		{
			ID:          "chrome",
			AppBundleID: "com.google.Chrome",
			Bindings: []config.Binding{
				{
					Device:   "mx-master-4",
					Trigger:  "hold(gesture_button)+move(down)",
					Action:   "close_tab",
					Priority: &low,
				},
				{
					Device:   "mx-master-4",
					Trigger:  "hold(gesture_button)+move(down)",
					Action:   "mission_control",
					Priority: &high,
				},
			},
		},
	}

	engine := NewEngine(cfg)
	action, err := engine.Match(Context{AppBundleID: "com.google.Chrome"}, events.DeviceEvent{
		DeviceID: "mx-master-4",
		Gesture:  "hold(gesture_button)+move(down)",
	})
	if err != nil {
		t.Fatalf("Match returned error: %v", err)
	}
	if action.ID != "mission_control" {
		t.Fatalf("action.ID = %q, want mission_control", action.ID)
	}
}

func TestMatchRejectsAmbiguousTopRank(t *testing.T) {
	priority := 25
	cfg := sampleConfig()
	cfg.Profiles = []config.Profile{
		{
			ID:          "chrome",
			AppBundleID: "com.google.Chrome",
			Bindings: []config.Binding{
				{
					Device:   "mx-master-4",
					Trigger:  "hold(gesture_button)+move(down)",
					Action:   "close_tab",
					Priority: &priority,
				},
				{
					Device:   "mx-master-4",
					Trigger:  "hold(gesture_button)+move(down)",
					Action:   "mission_control",
					Priority: &priority,
				},
			},
		},
	}

	engine := NewEngine(cfg)
	_, err := engine.Match(Context{AppBundleID: "com.google.Chrome"}, events.DeviceEvent{
		DeviceID: "mx-master-4",
		Gesture:  "hold(gesture_button)+move(down)",
	})
	if err == nil {
		t.Fatal("Match returned nil error, want ambiguity error")
	}
	if !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("Match error = %v, want ambiguity error", err)
	}
}

func TestMatchReturnsNoBindingForUnknownTrigger(t *testing.T) {
	engine := NewEngine(sampleConfig())
	_, err := engine.Match(Context{AppBundleID: "com.google.Chrome"}, events.DeviceEvent{
		DeviceID: "mx-master-4",
		Gesture:  "hold(gesture_button)+move(left)",
	})
	if err == nil {
		t.Fatal("Match returned nil error, want no-match error")
	}
	if !strings.Contains(err.Error(), "no binding") {
		t.Fatalf("Match error = %v, want no binding error", err)
	}
}

func TestMatchLoadsValidFixtureConfig(t *testing.T) {
	cfg, err := config.LoadFile(filepath.Join("..", "..", "testdata", "config", "valid.toml"))
	if err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}
	if err := config.Validate(cfg); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}

	engine := NewEngine(cfg)
	action, err := engine.Match(Context{AppBundleID: "com.google.Chrome"}, events.DeviceEvent{
		DeviceID: "mx-master-4",
		Gesture:  "hold(gesture_button)+move(down)",
	})
	if err != nil {
		t.Fatalf("Match returned error: %v", err)
	}
	if action.ID != "close_tab" {
		t.Fatalf("action.ID = %q, want close_tab", action.ID)
	}
}

func sampleConfig() *config.Config {
	return &config.Config{
		Devices: []config.Device{
			{
				ID:             "mx-master-4",
				MatchVendorID:  1133,
				MatchProductID: 0xb042,
			},
		},
		Actions: []config.Action{
			{ID: "close_tab", Type: "shortcut", Keys: []string{"cmd", "w"}},
			{ID: "mission_control", Type: "system", System: "mission_control"},
			{ID: "launchpad", Type: "system", System: "launchpad"},
			{ID: "device_global", Type: "system", System: "mission_control"},
			{ID: "app_specific", Type: "system", System: "launchpad"},
		},
		Profiles: []config.Profile{
			{
				ID: "global",
				Bindings: []config.Binding{
					{
						Device:  "mx-master-4",
						Trigger: "hold(gesture_button)+move(down)",
						Action:  "mission_control",
					},
					{
						Trigger: "hold(gesture_button)+move(down)",
						Action:  "launchpad",
					},
				},
			},
			{
				ID:          "chrome",
				AppBundleID: "com.google.Chrome",
				Bindings: []config.Binding{
					{
						Device:  "mx-master-4",
						Trigger: "hold(gesture_button)+move(down)",
						Action:  "close_tab",
					},
				},
			},
		},
	}
}
