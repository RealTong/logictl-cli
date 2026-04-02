//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit -framework Foundation
#include <stdlib.h>
#import <AppKit/AppKit.h>

static char *logi_frontmost_bundle_id(void) {
	@autoreleasepool {
		NSRunningApplication *app = [[NSWorkspace sharedWorkspace] frontmostApplication];
		if (app == nil || app.bundleIdentifier == nil) {
			return NULL;
		}
		return strdup([[app bundleIdentifier] UTF8String]);
	}
}
*/
import "C"

import (
	"context"
	"errors"
	"unsafe"
)

func defaultActiveBundleID(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	bundleID := C.logi_frontmost_bundle_id()
	if bundleID == nil {
		return "", errors.New("frontmost application is unavailable")
	}
	defer C.free(unsafe.Pointer(bundleID))

	return C.GoString(bundleID), nil
}
