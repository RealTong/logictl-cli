package cli

import (
	"bytes"
	"context"
	"errors"
	"testing"
)

type fakeDaemonPreflight struct {
	err error
}

func (f fakeDaemonPreflight) Preflight() error {
	return f.err
}

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
	cmd := newDaemonStartCmd(fakeDaemonPreflight{}, manager)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(manager.calls) != 1 || manager.calls[0] != "start" {
		t.Fatalf("manager.calls = %#v, want [start]", manager.calls)
	}
}

func TestDaemonRestartCmdInvokesServiceManager(t *testing.T) {
	manager := &fakeDaemonServiceManager{}
	cmd := newDaemonRestartCmd(fakeDaemonPreflight{}, manager)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(manager.calls) != 1 || manager.calls[0] != "restart" {
		t.Fatalf("manager.calls = %#v, want [restart]", manager.calls)
	}
}

func TestDaemonStartCmdRejectsPreflightFailures(t *testing.T) {
	manager := &fakeDaemonServiceManager{}
	cmd := newDaemonStartCmd(fakeDaemonPreflight{err: errors.New("unsafe path")}, manager)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil, want preflight error")
	}
	if len(manager.calls) != 0 {
		t.Fatalf("manager.calls = %#v, want no service-manager calls on preflight failure", manager.calls)
	}
}
