//go:build darwin

package macos

/*
#cgo LDFLAGS: -framework ApplicationServices -framework CoreFoundation
#include <ApplicationServices/ApplicationServices.h>
#include <stdint.h>
*/
import "C"

import "runtime/cgo"

//export logiHandleScrollEvent
func logiHandleScrollEvent(handle C.uintptr_t, axis1 C.int32_t, axis2 C.int32_t, pointAxis1 C.int32_t, pointAxis2 C.int32_t) C.int {
	rewriter, ok := cgo.Handle(handle).Value().(*EventTapScrollRewriter)
	if !ok {
		return 0
	}
	if rewriter.handleScroll(int(axis1), int(axis2), int(pointAxis1), int(pointAxis2)) {
		return 1
	}
	return 0
}
