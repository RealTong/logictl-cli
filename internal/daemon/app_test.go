package daemon

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/realtong/logictl-cli/internal/ipc"
)

func TestServerStatusEndpoint(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "daemon.sock")
	server := NewServer(socketPath, ipc.Status{Running: true})

	go func() {
		_ = server.Run(context.Background())
	}()

	deadline := time.Now().Add(2 * time.Second)
	for {
		status, err := ipc.QueryStatus(socketPath)
		if err == nil {
			if !status.Running {
				t.Fatal("status.Running = false, want true")
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("QueryStatus returned error: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestReloadDoesNotChangeSteadyStateStatus(t *testing.T) {
	socketPath := shortSocketPath(t)
	server := newServer(socketPath, NewRuntime())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = server.Run(ctx)
	}()

	waitForStatus(t, socketPath)

	reloadStatus, err := ipc.RequestReload(socketPath)
	if err != nil {
		t.Fatalf("RequestReload returned error: %v", err)
	}
	if got, want := reloadStatus.Message, "reload requested"; got != want {
		t.Fatalf("reload message = %q, want %q", got, want)
	}

	status, err := ipc.QueryStatus(socketPath)
	if err != nil {
		t.Fatalf("QueryStatus returned error: %v", err)
	}
	if got, want := status.Message, "running"; got != want {
		t.Fatalf("status.Message = %q, want %q", got, want)
	}
}

func waitForStatus(t *testing.T, socketPath string) {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for {
		if _, err := ipc.QueryStatus(socketPath); err == nil {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("socket %q did not become ready", socketPath)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func shortSocketPath(t *testing.T) string {
	t.Helper()

	socketPath := filepath.Join(os.TempDir(), fmt.Sprintf("logictl-daemon-%d.sock", time.Now().UnixNano()))
	t.Cleanup(func() {
		_ = os.Remove(socketPath)
	})
	return socketPath
}
