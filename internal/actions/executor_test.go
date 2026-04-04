package actions

import (
	"context"
	"errors"
	"testing"

	"github.com/realtong/logictl-cli/internal/config"
)

type fakeShortcutEmitter struct {
	events [][]string
	err    error
}

func (f *fakeShortcutEmitter) Emit(_ context.Context, keys []string) error {
	if f.err != nil {
		return f.err
	}
	f.events = append(f.events, append([]string(nil), keys...))
	return nil
}

type fakeSystemRunner struct {
	actions []string
	err     error
}

func (f *fakeSystemRunner) Run(_ context.Context, action string) error {
	if f.err != nil {
		return f.err
	}
	f.actions = append(f.actions, action)
	return nil
}

type fakeScriptRunner struct {
	paths []string
	err   error
}

func (f *fakeScriptRunner) Run(_ context.Context, path string) error {
	if f.err != nil {
		return f.err
	}
	f.paths = append(f.paths, path)
	return nil
}

func TestExecutorRunsShortcutAction(t *testing.T) {
	emitter := &fakeShortcutEmitter{}
	exec := Executor{ShortcutEmitter: emitter}

	err := exec.Execute(context.Background(), config.Action{ID: "close_tab", Type: "shortcut", Keys: []string{"cmd", "w"}})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(emitter.events) != 1 {
		t.Fatalf("len(emitter.events) = %d, want 1", len(emitter.events))
	}
	if got := emitter.events[0]; len(got) != 2 || got[0] != "cmd" || got[1] != "w" {
		t.Fatalf("emitter.events[0] = %#v, want [cmd w]", got)
	}
}

func TestExecutorRunsSystemAction(t *testing.T) {
	runner := &fakeSystemRunner{}
	exec := Executor{SystemRunner: runner}

	err := exec.Execute(context.Background(), config.Action{ID: "mission_control", Type: "system", System: "mission_control"})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(runner.actions) != 1 {
		t.Fatalf("len(runner.actions) = %d, want 1", len(runner.actions))
	}
	if got := runner.actions[0]; got != "mission_control" {
		t.Fatalf("runner.actions[0] = %q, want mission_control", got)
	}
}

func TestExecutorRunsScriptAction(t *testing.T) {
	runner := &fakeScriptRunner{}
	exec := Executor{ScriptRunner: runner}

	err := exec.Execute(context.Background(), config.Action{ID: "run_script", Type: "script", Script: "/tmp/test.sh"})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(runner.paths) != 1 {
		t.Fatalf("len(runner.paths) = %d, want 1", len(runner.paths))
	}
	if got := runner.paths[0]; got != "/tmp/test.sh" {
		t.Fatalf("runner.paths[0] = %q, want /tmp/test.sh", got)
	}
}

func TestExecutorRejectsUnsupportedActionType(t *testing.T) {
	err := (Executor{}).Execute(context.Background(), config.Action{ID: "noop", Type: "noop"})
	if err == nil {
		t.Fatal("Execute returned nil, want unsupported action type error")
	}
}

func TestExecutorPropagatesEmitterErrors(t *testing.T) {
	wantErr := errors.New("emit failed")
	exec := Executor{
		ShortcutEmitter: &fakeShortcutEmitter{err: wantErr},
	}

	err := exec.Execute(context.Background(), config.Action{ID: "close_tab", Type: "shortcut", Keys: []string{"cmd", "w"}})
	if !errors.Is(err, wantErr) {
		t.Fatalf("Execute error = %v, want %v", err, wantErr)
	}
}
