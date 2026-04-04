package ipc

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestQueryStatusMissingSocket(t *testing.T) {
	if _, err := QueryStatus("/tmp/does-not-exist.sock"); err == nil {
		t.Fatal("QueryStatus returned nil error, want socket error")
	}
}

func TestListenRejectsActiveSocket(t *testing.T) {
	socketPath := filepath.Join(os.TempDir(), fmt.Sprintf("logictl-ipc-%d.sock", time.Now().UnixNano()))
	t.Cleanup(func() {
		_ = os.Remove(socketPath)
	})

	listener, err := Listen(socketPath)
	if err != nil {
		t.Fatalf("Listen returned error: %v", err)
	}
	defer listener.Close()

	second, err := Listen(socketPath)
	if err == nil {
		second.Close()
		t.Fatal("second Listen returned nil error, want active socket error")
	}
}
