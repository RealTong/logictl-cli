//go:build !darwin

package actions

import (
	"context"
	"errors"
)

type AppleScriptShortcutEmitter struct{}

func (AppleScriptShortcutEmitter) Emit(context.Context, []string) error {
	return errors.New("shortcut emission is only supported on macOS")
}
