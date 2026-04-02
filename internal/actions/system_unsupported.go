//go:build !darwin

package actions

import (
	"context"
	"errors"
)

type MacOSSystemRunner struct {
	ShortcutEmitter ShortcutEmitter
}

func (MacOSSystemRunner) Run(context.Context, string) error {
	return errors.New("system actions are only supported on macOS")
}
