package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewRootCmdHelpHidesCompletion(t *testing.T) {
	cmd := NewRootCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "completion") {
		t.Fatalf("help output exposes completion command: %s", out)
	}
	if !strings.Contains(out, "version") {
		t.Fatalf("help output missing version command: %s", out)
	}
}
