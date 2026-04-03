package config

import (
	"fmt"
)

type bindingKey struct {
	app     string
	device  string
	trigger string
}

func Validate(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	actionIDs := make(map[string]Action, len(cfg.Actions))
	for _, action := range cfg.Actions {
		if action.ID == "" {
			return fmt.Errorf("action id is required")
		}
		if _, exists := actionIDs[action.ID]; exists {
			return fmt.Errorf("duplicate action id %q", action.ID)
		}
		if err := validateAction(action); err != nil {
			return err
		}
		actionIDs[action.ID] = action
	}

	deviceIDs := make(map[string]struct{}, len(cfg.Devices))
	for _, device := range cfg.Devices {
		if device.ID == "" {
			return fmt.Errorf("device id is required")
		}
		if _, exists := deviceIDs[device.ID]; exists {
			return fmt.Errorf("duplicate device id %q", device.ID)
		}
		if err := validateScrollConfig(device.ID, device.Scroll); err != nil {
			return err
		}
		deviceIDs[device.ID] = struct{}{}
	}

	seen := map[bindingKey][]Binding{}
	for _, profile := range cfg.Profiles {
		for _, binding := range profile.Bindings {
			if binding.Trigger == "" {
				return fmt.Errorf("profile %q has binding with empty trigger", profile.ID)
			}
			if binding.Action == "" {
				return fmt.Errorf("profile %q binding %q has empty action", profile.ID, binding.Trigger)
			}
			if _, ok := actionIDs[binding.Action]; !ok {
				return fmt.Errorf("binding %q references unknown action %q", binding.Trigger, binding.Action)
			}
			if binding.Device != "" {
				if _, ok := deviceIDs[binding.Device]; !ok {
					return fmt.Errorf("binding %q references unknown device %q", binding.Trigger, binding.Device)
				}
			}

			key := bindingKey{
				app:     profile.AppBundleID,
				device:  binding.Device,
				trigger: binding.Trigger,
			}

			if existingBindings, ok := seen[key]; ok {
				if hasAmbiguousCollision(existingBindings, binding) {
					return fmt.Errorf("ambiguous binding for app %q device %q trigger %q", profile.AppBundleID, binding.Device, binding.Trigger)
				}
			}
			seen[key] = append(seen[key], binding)
		}
	}

	return nil
}

func validateScrollConfig(deviceID string, scroll ScrollConfig) error {
	switch scroll.Direction {
	case "", "natural", "standard":
		return nil
	default:
		return fmt.Errorf("device %q has unsupported scroll direction %q", deviceID, scroll.Direction)
	}
}

func validateAction(action Action) error {
	switch action.Type {
	case "shortcut":
		if len(action.Keys) == 0 {
			return fmt.Errorf("action %q type shortcut requires keys", action.ID)
		}
	case "system":
		if action.System == "" {
			return fmt.Errorf("action %q type system requires system", action.ID)
		}
	case "script":
		if action.Script == "" {
			return fmt.Errorf("action %q type script requires script", action.ID)
		}
	default:
		return fmt.Errorf("action %q has unknown type %q", action.ID, action.Type)
	}

	return nil
}

func hasAmbiguousCollision(existing []Binding, current Binding) bool {
	for _, candidate := range existing {
		if candidate.Priority == nil || current.Priority == nil {
			return true
		}
		if *candidate.Priority == *current.Priority {
			return true
		}
	}
	return false
}
