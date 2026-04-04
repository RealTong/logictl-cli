package rules

import (
	"errors"
	"fmt"
	"sort"

	"github.com/realtong/logictl-cli/internal/config"
	"github.com/realtong/logictl-cli/internal/events"
)

var (
	ErrNoBinding        = errors.New("no binding matched event")
	ErrAmbiguousBinding = errors.New("ambiguous binding match")
)

type Context struct {
	AppBundleID string
}

type Engine struct {
	actions  map[string]config.Action
	profiles []config.Profile
}

type candidate struct {
	profile config.Profile
	binding config.Binding
	action  config.Action
	rank    int
}

func NewEngine(cfg *config.Config) *Engine {
	engine := &Engine{
		actions: make(map[string]config.Action),
	}
	if cfg == nil {
		return engine
	}

	engine.profiles = append(engine.profiles, cfg.Profiles...)
	for _, action := range cfg.Actions {
		engine.actions[action.ID] = action
	}

	return engine
}

func (e *Engine) Match(ctx Context, event events.DeviceEvent) (config.Action, error) {
	trigger, ok := triggerForEvent(event)
	if !ok {
		return config.Action{}, fmt.Errorf("%w for event kind %q", ErrNoBinding, event.Kind)
	}

	candidates, err := e.candidatesFor(ctx, event, trigger)
	if err != nil {
		return config.Action{}, err
	}
	if len(candidates) == 0 {
		return config.Action{}, fmt.Errorf("%w for app %q device %q trigger %q", ErrNoBinding, ctx.AppBundleID, event.DeviceID, trigger)
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].rank != candidates[j].rank {
			return candidates[i].rank > candidates[j].rank
		}
		if bindingPriority(candidates[i].binding) != bindingPriority(candidates[j].binding) {
			return bindingPriority(candidates[i].binding) > bindingPriority(candidates[j].binding)
		}
		if candidates[i].profile.AppBundleID != candidates[j].profile.AppBundleID {
			return candidates[i].profile.AppBundleID < candidates[j].profile.AppBundleID
		}
		if candidates[i].binding.Device != candidates[j].binding.Device {
			return candidates[i].binding.Device < candidates[j].binding.Device
		}
		if candidates[i].binding.Action != candidates[j].binding.Action {
			return candidates[i].binding.Action < candidates[j].binding.Action
		}
		return candidates[i].profile.ID < candidates[j].profile.ID
	})

	if len(candidates) > 1 && sameMatchScore(candidates[0], candidates[1]) {
		return config.Action{}, fmt.Errorf("%w for app %q device %q trigger %q", ErrAmbiguousBinding, ctx.AppBundleID, event.DeviceID, trigger)
	}

	return candidates[0].action, nil
}

func (e *Engine) candidatesFor(ctx Context, event events.DeviceEvent, trigger string) ([]candidate, error) {
	candidates := make([]candidate, 0)
	for _, profile := range e.profiles {
		if !matchesProfile(profile, ctx) {
			continue
		}
		for _, binding := range profile.Bindings {
			if !matchesBinding(binding, event, trigger) {
				continue
			}

			action, ok := e.actions[binding.Action]
			if !ok {
				return nil, fmt.Errorf("binding %q references unknown action %q", binding.Trigger, binding.Action)
			}

			candidates = append(candidates, candidate{
				profile: profile,
				binding: binding,
				action:  action,
				rank:    matchRank(profile, binding),
			})
		}
	}
	return candidates, nil
}

func sameMatchScore(left, right candidate) bool {
	return left.rank == right.rank && bindingPriority(left.binding) == bindingPriority(right.binding)
}
