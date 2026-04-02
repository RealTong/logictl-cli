package ipc

import "testing"

func TestQueryStatusMissingSocket(t *testing.T) {
	if _, err := QueryStatus("/tmp/does-not-exist.sock"); err == nil {
		t.Fatal("QueryStatus returned nil error, want socket error")
	}
}
