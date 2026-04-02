//go:build darwin

package macos

/*
#cgo LDFLAGS: -framework CoreFoundation -framework IOKit
#include <CoreFoundation/CoreFoundation.h>
#include <stdint.h>
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

//export logiHandleInputReport
func logiHandleInputReport(handle C.uintptr_t, reportID C.uint32_t, report *C.uint8_t, reportLength C.CFIndex) {
	sink, ok := cgo.Handle(handle).Value().(*hidReportCallbackSink)
	if !ok || report == nil || reportLength <= 0 {
		return
	}

	sink.push(uint32(reportID), C.GoBytes(unsafe.Pointer(report), C.int(reportLength)))
}
