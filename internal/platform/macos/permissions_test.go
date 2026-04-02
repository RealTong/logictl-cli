package macos

import (
	"context"
	"errors"
	"testing"
)

func TestEnvironmentPermissionsUsesInjectedLookup(t *testing.T) {
	want := PermissionReport{
		AccessibilityGranted:   true,
		InputMonitoringGranted: false,
	}
	env := Environment{
		lookupActiveBundleID: func(context.Context) (string, error) {
			return "", nil
		},
		lookupPermissions: func(context.Context) (PermissionReport, error) {
			return want, nil
		},
	}

	got, err := env.Permissions(context.Background())
	if err != nil {
		t.Fatalf("Permissions() returned error: %v", err)
	}
	if got != want {
		t.Fatalf("Permissions() = %#v, want %#v", got, want)
	}
}

func TestEnvironmentPermissionsPropagatesErrors(t *testing.T) {
	wantErr := errors.New("permissions failed")
	env := Environment{
		lookupActiveBundleID: func(context.Context) (string, error) {
			return "", nil
		},
		lookupPermissions: func(context.Context) (PermissionReport, error) {
			return PermissionReport{}, wantErr
		},
	}

	_, err := env.Permissions(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("Permissions() error = %v, want %v", err, wantErr)
	}
}
