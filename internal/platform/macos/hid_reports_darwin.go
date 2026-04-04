//go:build darwin

package macos

/*
#cgo LDFLAGS: -framework CoreFoundation -framework IOKit
#include <CoreFoundation/CoreFoundation.h>
#include <IOKit/IOReturn.h>
#include <IOKit/hid/IOHIDManager.h>
#include <IOKit/hid/IOHIDKeys.h>
#include <stdint.h>
#include <stdlib.h>

extern void logiHandleInputReport(uintptr_t handle, uint32_t reportID, uint8_t *report, CFIndex reportLength);

static void logi_bridge_input_report(void *context, IOReturn result, void *sender, IOHIDReportType type, uint32_t reportID, uint8_t *report, CFIndex reportLength) {
	if (result != kIOReturnSuccess || report == NULL || reportLength <= 0) {
		return;
	}
	logiHandleInputReport((uintptr_t)context, reportID, report, reportLength);
}

static void logi_register_input_report_callback(IOHIDManagerRef manager, void *context) {
	IOHIDManagerRegisterInputReportCallback(manager, logi_bridge_input_report, context);
}

static CFStringRef logi_create_cfstring(const char *value) {
	return CFStringCreateWithCString(kCFAllocatorDefault, value, kCFStringEncodingUTF8);
}

static void logi_dict_set_int(CFMutableDictionaryRef dict, const char *key, int value) {
	CFStringRef cfKey = logi_create_cfstring(key);
	if (!cfKey) {
		return;
	}
	CFNumberRef cfValue = CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &value);
	if (cfValue) {
		CFDictionarySetValue(dict, cfKey, cfValue);
		CFRelease(cfValue);
	}
	CFRelease(cfKey);
}

static void logi_dict_set_string(CFMutableDictionaryRef dict, const char *key, const char *value) {
	if (!value || value[0] == '\0') {
		return;
	}
	CFStringRef cfKey = logi_create_cfstring(key);
	CFStringRef cfValue = logi_create_cfstring(value);
	if (cfKey && cfValue) {
		CFDictionarySetValue(dict, cfKey, cfValue);
	}
	if (cfValue) {
		CFRelease(cfValue);
	}
	if (cfKey) {
		CFRelease(cfKey);
	}
}

static CFMutableDictionaryRef logi_create_match_dict(int vendor, int product, int usagePage, int usage, const char *serialNumber) {
	CFMutableDictionaryRef dict = CFDictionaryCreateMutable(
		kCFAllocatorDefault,
		0,
		&kCFTypeDictionaryKeyCallBacks,
		&kCFTypeDictionaryValueCallBacks
	);
	if (!dict) {
		return NULL;
	}

	logi_dict_set_int(dict, kIOHIDVendorIDKey, vendor);
	logi_dict_set_int(dict, kIOHIDProductIDKey, product);
	logi_dict_set_int(dict, kIOHIDPrimaryUsagePageKey, usagePage);
	logi_dict_set_int(dict, kIOHIDPrimaryUsageKey, usage);
	logi_dict_set_int(dict, kIOHIDDeviceUsagePageKey, usagePage);
	logi_dict_set_int(dict, kIOHIDDeviceUsageKey, usage);
	logi_dict_set_string(dict, kIOHIDSerialNumberKey, serialNumber);

	return dict;
}

static void logi_stop_run_loop(CFRunLoopRef runLoop) {
	if (runLoop != NULL) {
		CFRunLoopStop(runLoop);
	}
}
*/
import "C"

import (
	"context"
	"fmt"
	"runtime"
	"runtime/cgo"
	"sync"
	"time"
	"unsafe"

	"github.com/realtong/logictl-cli/internal/events"
)

type hidReportSourceFactory struct{}

func NewHIDReportSourceFactory() HIDReportSourceFactory {
	return hidReportSourceFactory{}
}

func (hidReportSourceFactory) Validate(match HIDReportMatch) error {
	count, err := matchedNativeDevices(match)
	if err != nil {
		return err
	}

	switch count {
	case 0:
		return fmt.Errorf("no native HID report devices matched %s", describeHIDMatch(match))
	case 1:
		return nil
	default:
		return fmt.Errorf("multiple native HID report devices matched %s", describeHIDMatch(match))
	}
}

func (hidReportSourceFactory) Open(match HIDReportMatch) events.Source {
	return hidReportSource{
		match: match,
		now:   time.Now,
	}
}

type hidReportSource struct {
	match HIDReportMatch
	now   func() time.Time
}

type hidReportCallbackSink struct {
	ctx     context.Context
	now     func() time.Time
	reports chan events.RawReport
	errs    chan error

	mu     sync.RWMutex
	closed bool
}

func (s hidReportSource) Stream(ctx context.Context) (<-chan events.RawReport, <-chan error) {
	reports := make(chan events.RawReport, 32)
	errs := make(chan error, 1)

	go s.run(ctx, reports, errs)

	return reports, errs
}

func (s hidReportSource) run(ctx context.Context, reports chan events.RawReport, errs chan error) {
	defer close(reports)
	defer close(errs)

	sink := &hidReportCallbackSink{
		ctx:     ctx,
		now:     s.now,
		reports: reports,
		errs:    errs,
	}
	handle := cgo.NewHandle(sink)
	defer handle.Delete()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	manager, matchDict, err := newNativeManager(s.match)
	if err != nil {
		sink.reportErr(err)
		return
	}
	defer C.CFRelease(C.CFTypeRef(matchDict))
	defer C.CFRelease(C.CFTypeRef(manager))

	C.IOHIDManagerSetDeviceMatching(manager, C.CFDictionaryRef(matchDict))
	C.logi_register_input_report_callback(manager, unsafe.Pointer(uintptr(handle)))

	runLoop := C.CFRunLoopGetCurrent()
	C.IOHIDManagerScheduleWithRunLoop(manager, runLoop, C.kCFRunLoopDefaultMode)
	defer C.IOHIDManagerUnscheduleFromRunLoop(manager, runLoop, C.kCFRunLoopDefaultMode)

	if rc := C.IOHIDManagerOpen(manager, C.kIOHIDOptionsTypeNone); rc != C.kIOReturnSuccess {
		sink.reportErr(fmt.Errorf("IOHIDManagerOpen failed for %s: 0x%x", describeHIDMatch(s.match), uint32(rc)))
		return
	}
	defer C.IOHIDManagerClose(manager, C.kIOHIDOptionsTypeNone)

	count, err := matchedDevicesFromManager(manager)
	if err != nil {
		sink.reportErr(err)
		return
	}
	switch count {
	case 0:
		sink.reportErr(fmt.Errorf("no native HID report devices matched %s", describeHIDMatch(s.match)))
		return
	case 1:
	default:
		sink.reportErr(fmt.Errorf("multiple native HID report devices matched %s", describeHIDMatch(s.match)))
		return
	}

	done := make(chan struct{})
	go func(loop C.CFRunLoopRef) {
		select {
		case <-ctx.Done():
			C.logi_stop_run_loop(loop)
		case <-done:
		}
	}(runLoop)

	C.CFRunLoopRun()
	close(done)
	sink.markClosed()
}

func (s *hidReportCallbackSink) push(reportID uint32, payload []byte) {
	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return
	}
	reports := s.reports
	ctx := s.ctx
	now := s.now
	s.mu.RUnlock()

	report := events.RawReport{
		At:    now(),
		Bytes: normalizeReportBytes(reportID, payload),
	}

	select {
	case reports <- report:
	case <-ctx.Done():
	}
}

func (s *hidReportCallbackSink) reportErr(err error) {
	if err == nil {
		return
	}

	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return
	}
	errs := s.errs
	s.mu.RUnlock()

	select {
	case errs <- err:
	default:
	}
}

func (s *hidReportCallbackSink) markClosed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
}

func matchedNativeDevices(match HIDReportMatch) (int, error) {
	manager, matchDict, err := newNativeManager(match)
	if err != nil {
		return 0, err
	}
	defer C.CFRelease(C.CFTypeRef(matchDict))
	defer C.CFRelease(C.CFTypeRef(manager))

	C.IOHIDManagerSetDeviceMatching(manager, C.CFDictionaryRef(matchDict))
	if rc := C.IOHIDManagerOpen(manager, C.kIOHIDOptionsTypeNone); rc != C.kIOReturnSuccess {
		return 0, fmt.Errorf("IOHIDManagerOpen failed for %s: 0x%x", describeHIDMatch(match), uint32(rc))
	}
	defer C.IOHIDManagerClose(manager, C.kIOHIDOptionsTypeNone)

	return matchedDevicesFromManager(manager)
}

func matchedDevicesFromManager(manager C.IOHIDManagerRef) (int, error) {
	set := C.IOHIDManagerCopyDevices(manager)
	if set == 0 {
		return 0, nil
	}
	defer C.CFRelease(C.CFTypeRef(set))

	return int(C.CFSetGetCount(set)), nil
}

func newNativeManager(match HIDReportMatch) (C.IOHIDManagerRef, C.CFMutableDictionaryRef, error) {
	manager := C.IOHIDManagerCreate(C.kCFAllocatorDefault, C.kIOHIDManagerOptionNone)
	if manager == 0 {
		return 0, 0, fmt.Errorf("IOHIDManagerCreate failed for %s", describeHIDMatch(match))
	}

	var serialPtr *C.char
	if match.SerialNumber != "" {
		serialPtr = C.CString(match.SerialNumber)
		defer C.free(unsafe.Pointer(serialPtr))
	}

	matchDict := C.logi_create_match_dict(
		C.int(match.VendorID),
		C.int(match.ProductID),
		C.int(match.UsagePage),
		C.int(match.Usage),
		serialPtr,
	)
	if matchDict == 0 {
		C.CFRelease(C.CFTypeRef(manager))
		return 0, 0, fmt.Errorf("failed to build native HID match for %s", describeHIDMatch(match))
	}

	return manager, matchDict, nil
}

func normalizeReportBytes(reportID uint32, payload []byte) []byte {
	report := append([]byte(nil), payload...)
	if reportID != 0 && (len(report) == 0 || report[0] != byte(reportID)) {
		report = append([]byte{byte(reportID)}, report...)
	}
	return report
}

func describeHIDMatch(match HIDReportMatch) string {
	product := match.Product
	if product == "" {
		product = "Logitech HID device"
	}
	return fmt.Sprintf(
		"%s (vendor=0x%04x product=0x%04x usagePage=0x%04x usage=0x%04x)",
		product,
		match.VendorID,
		match.ProductID,
		match.UsagePage,
		match.Usage,
	)
}
