package rules

import (
	"fmt"

	"github.com/realtong/logictl-cli/internal/config"
	"github.com/realtong/logictl-cli/internal/events"
)

func triggerForEvent(event events.DeviceEvent) (string, bool) {
	if event.Gesture != "" {
		return event.Gesture, true
	}

	switch event.Kind {
	case events.ButtonDown:
		return fmt.Sprintf("%s_down", event.Control), event.Control != ""
	case events.ButtonHold:
		return fmt.Sprintf("%s_hold", event.Control), event.Control != ""
	case events.ButtonUp:
		return fmt.Sprintf("%s_up", event.Control), event.Control != ""
	default:
		return "", false
	}
}

func matchesProfile(profile config.Profile, ctx Context) bool {
	return profile.AppBundleID == "" || profile.AppBundleID == ctx.AppBundleID
}

func matchesBinding(binding config.Binding, event events.DeviceEvent, trigger string) bool {
	if binding.Trigger != trigger {
		return false
	}
	return binding.Device == "" || binding.Device == event.DeviceID
}

func matchRank(profile config.Profile, binding config.Binding) int {
	rank := 0
	if binding.Device != "" {
		rank += 2
	}
	if profile.AppBundleID != "" {
		rank += 1
	}
	return rank
}

func bindingPriority(binding config.Binding) int {
	if binding.Priority == nil {
		return 0
	}
	return *binding.Priority
}
