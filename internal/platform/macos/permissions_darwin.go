//go:build darwin

package macos

/*
#cgo LDFLAGS: -framework ApplicationServices -framework IOKit
#include <ApplicationServices/ApplicationServices.h>
#include <IOKit/hidsystem/IOHIDLib.h>

static int logi_accessibility_granted(void) {
	return AXIsProcessTrusted() ? 1 : 0;
}

static int logi_input_monitoring_granted(void) {
	return IOHIDCheckAccess(kIOHIDRequestTypeListenEvent) == kIOHIDAccessTypeGranted ? 1 : 0;
}
*/
import "C"

import "context"

func defaultPermissionReport(ctx context.Context) (PermissionReport, error) {
	if err := ctx.Err(); err != nil {
		return PermissionReport{}, err
	}

	return PermissionReport{
		AccessibilityGranted:   C.logi_accessibility_granted() != 0,
		InputMonitoringGranted: C.logi_input_monitoring_granted() != 0,
	}, nil
}
