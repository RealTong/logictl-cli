package actions

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRunScriptTimesOut(t *testing.T) {
	scriptPath := writeExecutableScript(t, "#!/bin/sh\nsleep 1\n")

	err := runScript(context.Background(), ScriptAction{
		Path:    scriptPath,
		Timeout: 10 * time.Millisecond,
	})
	if err == nil {
		t.Fatal("runScript returned nil, want timeout error")
	}
}

func TestRunScriptExecutesExecutable(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "out.txt")
	scriptPath := writeExecutableScript(t, "#!/bin/sh\necho ok > \""+outputPath+"\"\n")

	err := runScript(context.Background(), ScriptAction{
		Path:    scriptPath,
		Timeout: time.Second,
	})
	if err != nil {
		t.Fatalf("runScript returned error: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) returned error: %v", outputPath, err)
	}
	if string(data) != "ok\n" {
		t.Fatalf("script output = %q, want %q", string(data), "ok\n")
	}
}

func writeExecutableScript(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "script.sh")
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("WriteFile(%q) returned error: %v", path, err)
	}
	return path
}
