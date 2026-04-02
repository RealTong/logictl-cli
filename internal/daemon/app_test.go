package daemon

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/realtong/logi-cli/internal/ipc"
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
