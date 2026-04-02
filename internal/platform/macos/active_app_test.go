package macos

import (
	"context"
	"errors"
	"testing"
)

func TestEnvironmentActiveBundleIDUsesInjectedLookup(t *testing.T) {
	env := Environment{
		lookupActiveBundleID: func(context.Context) (string, error) {
			return "com.google.Chrome", nil
		},
		lookupPermissions: func(context.Context) (PermissionReport, error) {
			return PermissionReport{}, nil
		},
	}

	got, err := env.ActiveBundleID(context.Background())
	if err != nil {
		t.Fatalf("ActiveBundleID() returned error: %v", err)
	}
	if got != "com.google.Chrome" {
		t.Fatalf("ActiveBundleID() = %q, want com.google.Chrome", got)
	}
}

func TestEnvironmentActiveBundleIDPropagatesErrors(t *testing.T) {
	wantErr := errors.New("lookup failed")
	env := Environment{
		lookupActiveBundleID: func(context.Context) (string, error) {
			return "", wantErr
		},
		lookupPermissions: func(context.Context) (PermissionReport, error) {
			return PermissionReport{}, nil
		},
	}

	_, err := env.ActiveBundleID(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("ActiveBundleID() error = %v, want %v", err, wantErr)
	}
}
