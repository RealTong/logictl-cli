package actions

import (
	"context"
	"time"
)

type ShortcutEmitter interface {
	Emit(ctx context.Context, keys []string) error
}

type SystemRunner interface {
	Run(ctx context.Context, action string) error
}

type ScriptRunner interface {
	Run(ctx context.Context, path string) error
}

type ScriptAction struct {
	Path    string
	Timeout time.Duration
}
