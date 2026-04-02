package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/realtong/logi-cli/internal/daemon"
)

type fakeDaemonServiceManager struct {
	calls []string
}

func (m *fakeDaemonServiceManager) Start(context.Context) error {
	m.calls = append(m.calls, "start")
	return nil
}

func (m *fakeDaemonServiceManager) Stop(context.Context) error {
	m.calls = append(m.calls, "stop")
	return nil
}

func (m *fakeDaemonServiceManager) Restart(context.Context) error {
	m.calls = append(m.calls, "restart")
	return nil
}

func TestDaemonStartCmdInvokesServiceManager(t *testing.T) {
	manager := &fakeDaemonServiceManager{}
	cmd := newDaemonCmdWithServiceManager(daemon.NewApp(testPaths(t)), manager)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"start"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(manager.calls) != 1 || manager.calls[0] != "start" {
		t.Fatalf("manager.calls = %#v, want [start]", manager.calls)
	}
}

func TestDaemonRestartCmdInvokesServiceManager(t *testing.T) {
	manager := &fakeDaemonServiceManager{}
	cmd := newDaemonCmdWithServiceManager(daemon.NewApp(testPaths(t)), manager)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"restart"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(manager.calls) != 1 || manager.calls[0] != "restart" {
		t.Fatalf("manager.calls = %#v, want [restart]", manager.calls)
	}
}
