package cli

import (
	"strings"
	"testing"
)

func TestValidateCmdAcceptsValidFixture(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{
		"validate",
		"--config",
		"../../testdata/config/valid.toml",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
}

func TestValidateCmdRejectsDuplicateBindingFixture(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{
		"validate",
		"--config",
		"../../testdata/config/duplicate_binding.toml",
	})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil, want validation error")
	}
	if !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("Execute error = %v, want ambiguous binding error", err)
	}
}
