package actions

import (
	"context"
	"fmt"

	"github.com/realtong/logi-cli/internal/config"
)

type Executor struct {
	ShortcutEmitter ShortcutEmitter
	SystemRunner    SystemRunner
	ScriptRunner    ScriptRunner
}

func (e Executor) Execute(ctx context.Context, action config.Action) error {
	switch action.Type {
	case "shortcut":
		if e.ShortcutEmitter == nil {
			return fmt.Errorf("shortcut emitter is not configured")
		}
		return e.ShortcutEmitter.Emit(ctx, action.Keys)
	case "system":
		if e.SystemRunner == nil {
			return fmt.Errorf("system runner is not configured")
		}
		return e.SystemRunner.Run(ctx, action.System)
	case "script":
		if e.ScriptRunner == nil {
			return fmt.Errorf("script runner is not configured")
		}
		return e.ScriptRunner.Run(ctx, action.Script)
	default:
		return fmt.Errorf("unsupported action type %q", action.Type)
	}
}
