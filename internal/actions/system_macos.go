//go:build darwin

package actions

import (
	"context"
	"fmt"
)

type MacOSSystemRunner struct {
	ShortcutEmitter ShortcutEmitter
}

func (r MacOSSystemRunner) Run(ctx context.Context, action string) error {
	if r.ShortcutEmitter == nil {
		return fmt.Errorf("shortcut emitter is not configured")
	}

	switch action {
	case "mission_control":
		return r.ShortcutEmitter.Emit(ctx, []string{"ctrl", "up"})
	case "app_expose":
		return r.ShortcutEmitter.Emit(ctx, []string{"ctrl", "down"})
	case "launchpad":
		return r.ShortcutEmitter.Emit(ctx, []string{"f4"})
	default:
		return fmt.Errorf("unsupported system action %q", action)
	}
}
