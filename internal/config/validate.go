package config

import "fmt"

func Validate(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	actionIDs := make(map[string]struct{}, len(cfg.Actions))
	for _, action := range cfg.Actions {
		if action.ID == "" {
			return fmt.Errorf("action id is required")
		}
		if _, exists := actionIDs[action.ID]; exists {
			return fmt.Errorf("duplicate action id %q", action.ID)
		}
		actionIDs[action.ID] = struct{}{}
	}

	for _, profile := range cfg.Profiles {
		seenBindings := map[string]struct{}{}
		for _, binding := range profile.Bindings {
			if binding.Trigger == "" {
				return fmt.Errorf("profile %q has binding with empty trigger", profile.Name)
			}
			if binding.Action == "" {
				return fmt.Errorf("profile %q binding %q has empty action", profile.Name, binding.Trigger)
			}
			if _, ok := actionIDs[binding.Action]; !ok {
				return fmt.Errorf("binding %q references unknown action %q", binding.Trigger, binding.Action)
			}
			if _, exists := seenBindings[binding.Trigger]; exists {
				return fmt.Errorf("profile %q has duplicate binding for %q", profile.Name, binding.Trigger)
			}
			seenBindings[binding.Trigger] = struct{}{}
		}
	}

	return nil
}
