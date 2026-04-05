package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/realtong/logictl-cli/internal/ipc"
	platformmacos "github.com/realtong/logictl-cli/internal/platform/macos"
)

type fakePlatformDoctor struct {
	activeBundleID string
	report         platformmacos.PermissionReport
}

func (f fakePlatformDoctor) ActiveBundleID(context.Context) (string, error) {
	return f.activeBundleID, nil
}

func (f fakePlatformDoctor) Permissions(context.Context) (platformmacos.PermissionReport, error) {
	return f.report, nil
}

type fakeStatusReporter struct {
	status ipc.Status
}

func (f fakeStatusReporter) Status() (ipc.Status, error) {
	return f.status, nil
}

func TestDoctorReportsMissingAccessibilityPermission(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newDoctorCmd(
		fakePlatformDoctor{
			activeBundleID: "com.google.Chrome",
			report: platformmacos.PermissionReport{
				AccessibilityGranted:   false,
				InputMonitoringGranted: true,
			},
		},
		fakeStatusReporter{status: ipc.Status{Message: "stopped"}},
		"../../testdata/config/valid.toml",
	)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "Accessibility: missing") {
		t.Fatalf("output = %q, want missing accessibility status", got)
	}
	if !strings.Contains(got, "Input Monitoring: granted") {
		t.Fatalf("output = %q, want granted input monitoring status", got)
	}
	if !strings.Contains(got, "Frontmost App: com.google.Chrome") {
		t.Fatalf("output = %q, want active bundle ID", got)
	}
	if !strings.Contains(got, "Config Path: ../../testdata/config/valid.toml") {
		t.Fatalf("output = %q, want config path", got)
	}
	if !strings.Contains(got, "Config: valid") {
		t.Fatalf("output = %q, want valid config status", got)
	}
	if !strings.Contains(got, "Daemon: stopped") {
		t.Fatalf("output = %q, want stopped daemon status", got)
	}
}

func TestDoctorReportsMissingConfig(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newDoctorCmd(
		fakePlatformDoctor{
			activeBundleID: "",
			report: platformmacos.PermissionReport{
				AccessibilityGranted:   true,
				InputMonitoringGranted: false,
			},
		},
		fakeStatusReporter{status: ipc.Status{Running: true, Message: "running"}},
		"/tmp/definitely-missing-logictl-config.toml",
	)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "Frontmost App: unavailable") {
		t.Fatalf("output = %q, want unavailable active app", got)
	}
	if !strings.Contains(got, "Config Path: /tmp/definitely-missing-logictl-config.toml") {
		t.Fatalf("output = %q, want config path", got)
	}
	if !strings.Contains(got, "Config: missing") {
		t.Fatalf("output = %q, want missing config status", got)
	}
	if !strings.Contains(got, "Daemon: running") {
		t.Fatalf("output = %q, want running daemon status", got)
	}
}
