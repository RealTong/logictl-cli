//go:build darwin

package macos

/*
#cgo LDFLAGS: -framework ApplicationServices -framework CoreFoundation
#include <ApplicationServices/ApplicationServices.h>
#include <CoreFoundation/CoreFoundation.h>
#include <stdint.h>

extern int logiHandleScrollEvent(uintptr_t handle, int32_t axis1, int32_t axis2, int32_t pointAxis1, int32_t pointAxis2);

static const int64_t logi_scroll_event_marker = 0x6c6f6769;

static CGEventRef logi_bridge_scroll_event(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *context) {
	if (type != kCGEventScrollWheel) {
		return event;
	}

	if (CGEventGetIntegerValueField(event, kCGEventSourceUserData) == logi_scroll_event_marker) {
		return event;
	}

	int32_t axis1 = (int32_t)CGEventGetIntegerValueField(event, kCGScrollWheelEventDeltaAxis1);
	int32_t axis2 = (int32_t)CGEventGetIntegerValueField(event, kCGScrollWheelEventDeltaAxis2);
	int32_t pointAxis1 = (int32_t)CGEventGetIntegerValueField(event, kCGScrollWheelEventPointDeltaAxis1);
	int32_t pointAxis2 = (int32_t)CGEventGetIntegerValueField(event, kCGScrollWheelEventPointDeltaAxis2);

	if (logiHandleScrollEvent((uintptr_t)context, axis1, axis2, pointAxis1, pointAxis2) != 0) {
		return NULL;
	}
	return event;
}

static void logi_stop_scroll_run_loop(CFRunLoopRef runLoop) {
	if (runLoop != NULL) {
		CFRunLoopStop(runLoop);
	}
}

static CGEventMask logi_scroll_event_mask(void) {
	return CGEventMaskBit(kCGEventScrollWheel);
}

static CGEventRef logi_create_scroll_event(CGScrollEventUnit unit, int32_t axis1, int32_t axis2) {
	return CGEventCreateScrollWheelEvent(NULL, unit, 2, axis1, axis2);
}

static CFMachPortRef logi_create_scroll_event_tap(void *context) {
	return CGEventTapCreate(
		kCGSessionEventTap,
		kCGHeadInsertEventTap,
		kCGEventTapOptionDefault,
		logi_scroll_event_mask(),
		logi_bridge_scroll_event,
		context
	);
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

	"github.com/realtong/logictl-cli/internal/config"
)

type EventTapScrollRewriter struct {
	matcher *scrollMatcher
	emitter scrollEmitter
	now     func() time.Time

	errOnce  sync.Once
	mu       sync.Mutex
	firstErr error
	errCh    chan error
}

type scrollEmitter interface {
	Emit(scrollRewritePlan) error
}

type coreGraphicsScrollEmitter struct{}

func NewScrollRewriter() ScrollRewriter {
	return &EventTapScrollRewriter{
		matcher: newScrollMatcher(defaultScrollMatchWindow),
		emitter: coreGraphicsScrollEmitter{},
		now:     time.Now,
		errCh:   make(chan error, 1),
	}
}

func (r *EventTapScrollRewriter) Start(ctx context.Context) error {
	handle := cgo.NewHandle(r)
	defer handle.Delete()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	tap := C.logi_create_scroll_event_tap(unsafe.Pointer(uintptr(handle)))
	if tap == 0 {
		return fmt.Errorf("CGEventTapCreate failed for scroll rewrite; check Accessibility permission")
	}
	defer C.CFRelease(C.CFTypeRef(tap))

	source := C.CFMachPortCreateRunLoopSource(C.kCFAllocatorDefault, tap, 0)
	if source == 0 {
		return fmt.Errorf("CFMachPortCreateRunLoopSource failed for scroll rewrite")
	}
	defer C.CFRelease(C.CFTypeRef(source))

	runLoop := C.CFRunLoopGetCurrent()
	C.CFRunLoopAddSource(runLoop, source, C.kCFRunLoopCommonModes)
	C.CGEventTapEnable(tap, true)
	defer C.CFRunLoopRemoveSource(runLoop, source, C.kCFRunLoopCommonModes)

	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			C.logi_stop_scroll_run_loop(runLoop)
		case <-done:
		}
	}()

	go func() {
		select {
		case err := <-r.errCh:
			if err != nil {
				C.logi_stop_scroll_run_loop(runLoop)
			}
		case <-done:
		}
	}()

	C.CFRunLoopRun()
	close(done)

	r.mu.Lock()
	err := r.firstErr
	r.mu.Unlock()
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil && err != context.Canceled {
		return err
	}
	return nil
}

func (r *EventTapScrollRewriter) Record(deviceID, gesture string, settings config.ScrollConfig, at time.Time) {
	r.matcher.Record(deviceID, gesture, settings, at)
}

func (r *EventTapScrollRewriter) handleScroll(axis1, axis2, pointAxis1, pointAxis2 int) bool {
	plan, ok := r.matcher.Match(nativeScrollEvent{
		VerticalLine:    axis1,
		HorizontalLine:  axis2,
		VerticalPoint:   pointAxis1,
		HorizontalPoint: pointAxis2,
		At:              r.now(),
	})
	if !ok {
		return false
	}

	if err := r.emitter.Emit(plan); err != nil {
		r.reportErr(err)
		return false
	}
	return true
}

func (r *EventTapScrollRewriter) reportErr(err error) {
	if err == nil {
		return
	}
	r.errOnce.Do(func() {
		r.mu.Lock()
		r.firstErr = err
		r.mu.Unlock()
		r.errCh <- err
	})
}

func (coreGraphicsScrollEmitter) Emit(plan scrollRewritePlan) error {
	for _, emission := range plan.Emissions {
		if emission.Vertical == 0 && emission.Horizontal == 0 {
			continue
		}

		unit := C.CGScrollEventUnit(C.kCGScrollEventUnitLine)
		if emission.Unit == scrollUnitPixel {
			unit = C.CGScrollEventUnit(C.kCGScrollEventUnitPixel)
		}

		event := C.logi_create_scroll_event(
			unit,
			C.int32_t(emission.Vertical),
			C.int32_t(emission.Horizontal),
		)
		if event == 0 {
			return fmt.Errorf("CGEventCreateScrollWheelEvent failed")
		}

		C.CGEventSetIntegerValueField(event, C.kCGEventSourceUserData, C.logi_scroll_event_marker)
		C.CGEventPost(C.kCGSessionEventTap, event)
		C.CFRelease(C.CFTypeRef(event))
	}
	return nil
}
