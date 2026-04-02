package actions

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const defaultScriptTimeout = 5 * time.Second

type CommandScriptRunner struct {
	Timeout time.Duration
}

func (r CommandScriptRunner) Run(ctx context.Context, path string) error {
	timeout := r.Timeout
	if timeout <= 0 {
		timeout = defaultScriptTimeout
	}

	return runScript(ctx, ScriptAction{
		Path:    path,
		Timeout: timeout,
	})
}

func runScript(parent context.Context, action ScriptAction) error {
	if action.Path == "" {
		return errors.New("script path is required")
	}

	timeout := action.Timeout
	if timeout <= 0 {
		timeout = defaultScriptTimeout
	}

	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	path, err := expandHome(action.Path)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, path)
	if err := cmd.Run(); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("script timed out after %s", timeout)
		}
		return err
	}

	return nil
}

func expandHome(path string) (string, error) {
	if !strings.HasPrefix(path, "~/") {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, strings.TrimPrefix(path, "~/")), nil
}
