//go:build darwin

package actions

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"unicode/utf8"
)

type AppleScriptShortcutEmitter struct{}

func (AppleScriptShortcutEmitter) Emit(ctx context.Context, keys []string) error {
	script, err := shortcutAppleScript(keys)
	if err != nil {
		return err
	}
	return exec.CommandContext(ctx, "osascript", "-e", script).Run()
}

func shortcutAppleScript(keys []string) (string, error) {
	primary, modifiers, err := parseShortcut(keys)
	if err != nil {
		return "", err
	}

	if code, ok := specialKeyCode(primary); ok {
		if len(modifiers) == 0 {
			return fmt.Sprintf(`tell application "System Events" to key code %d`, code), nil
		}
		return fmt.Sprintf(`tell application "System Events" to key code %d using {%s}`, code, strings.Join(modifiers, ", ")), nil
	}

	if utf8.RuneCountInString(primary) != 1 {
		return "", fmt.Errorf("unsupported shortcut key %q", primary)
	}

	if len(modifiers) == 0 {
		return fmt.Sprintf(`tell application "System Events" to keystroke %q`, primary), nil
	}
	return fmt.Sprintf(`tell application "System Events" to keystroke %q using {%s}`, primary, strings.Join(modifiers, ", ")), nil
}

func parseShortcut(keys []string) (string, []string, error) {
	if len(keys) == 0 {
		return "", nil, fmt.Errorf("shortcut requires at least one key")
	}

	modifiers := make([]string, 0, len(keys))
	var primary string
	for _, key := range keys {
		switch strings.ToLower(strings.TrimSpace(key)) {
		case "cmd", "command":
			modifiers = append(modifiers, "command down")
		case "ctrl", "control":
			modifiers = append(modifiers, "control down")
		case "alt", "option":
			modifiers = append(modifiers, "option down")
		case "shift":
			modifiers = append(modifiers, "shift down")
		default:
			if primary != "" {
				return "", nil, fmt.Errorf("shortcut must have exactly one non-modifier key")
			}
			primary = strings.ToLower(strings.TrimSpace(key))
		}
	}

	if primary == "" {
		return "", nil, fmt.Errorf("shortcut requires one non-modifier key")
	}

	return primary, modifiers, nil
}

func specialKeyCode(key string) (int, bool) {
	switch key {
	case "up":
		return 126, true
	case "down":
		return 125, true
	case "left":
		return 123, true
	case "right":
		return 124, true
	case "return":
		return 36, true
	case "enter":
		return 76, true
	case "tab":
		return 48, true
	case "space":
		return 49, true
	case "escape", "esc":
		return 53, true
	case "delete":
		return 51, true
	case "forward_delete":
		return 117, true
	case "f1":
		return 122, true
	case "f2":
		return 120, true
	case "f3":
		return 99, true
	case "f4":
		return 118, true
	case "f5":
		return 96, true
	case "f6":
		return 97, true
	case "f7":
		return 98, true
	case "f8":
		return 100, true
	case "f9":
		return 101, true
	case "f10":
		return 109, true
	case "f11":
		return 103, true
	case "f12":
		return 111, true
	default:
		return 0, false
	}
}
