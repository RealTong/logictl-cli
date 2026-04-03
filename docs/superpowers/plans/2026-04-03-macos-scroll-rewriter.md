# macOS Scroll Rewriter Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `logi-cli` daemon rewrite `MX Master 4` scroll events on macOS so `devices.scroll.direction` and `devices.scroll.smooth_scroll` actually change wheel behavior.

**Architecture:** Add a dedicated macOS scroll rewriter component that combines a global CoreGraphics scroll-wheel event tap with a short-lived matcher fed by decoded `MX Master 4` wheel HID events. The daemon will only swallow and re-emit scroll events when a recent HID wheel signal matches, so other devices and the trackpad stay on the native path.

**Tech Stack:** Go, CoreGraphics via cgo, existing macOS HID report source, Go `testing`.

---

### Task 1: Add the runtime plumbing for scroll rewrite configuration

**Files:**
- Modify: `internal/daemon/runtime.go`
- Modify: `internal/daemon/app.go`
- Modify: `internal/config/config.go`
- Modify: `internal/daemon/runtime_test.go`

- [ ] **Step 1: Write the failing runtime test**

```go
func TestRuntimeApplyConfigBuildsScrollSettings(t *testing.T) {
	runtime := NewRuntimeWithDependencies(RuntimeDependencies{})
	cfg := sampleRuntimeConfig("global")
	cfg.Devices[0].Scroll.Direction = "natural"
	cfg.Devices[0].Scroll.SmoothScroll = false

	if err := runtime.ApplyConfig(cfg); err != nil {
		t.Fatalf("ApplyConfig returned error: %v", err)
	}

	settings := runtime.ScrollSettings("mx-master-4")
	if settings.Direction != "natural" || settings.SmoothScroll {
		t.Fatalf("settings = %#v, want natural + smooth false", settings)
	}
}
```

- [ ] **Step 2: Run the focused runtime test to verify it fails**

Run: `go test ./internal/daemon -run TestRuntimeApplyConfigBuildsScrollSettings -v`  
Expected: FAIL with missing scroll settings accessors

- [ ] **Step 3: Implement scroll settings extraction**

```go
type DeviceScrollSettings struct {
	Direction    string
	SmoothScroll bool
}
```

- [ ] **Step 4: Re-run the runtime package**

Run: `go test ./internal/daemon -v`  
Expected: PASS

### Task 2: Add a macOS scroll rewrite engine with matching and re-emission

**Files:**
- Create: `internal/platform/macos/scroll_rewriter.go`
- Create: `internal/platform/macos/scroll_rewriter_darwin.go`
- Create: `internal/platform/macos/scroll_rewriter_unsupported.go`
- Create: `internal/platform/macos/scroll_rewriter_test.go`

- [ ] **Step 1: Write the failing matcher tests**

```go
func TestMatcherOnlyConsumesRecentMatchingWheelTicks(t *testing.T) {
	m := newScrollMatcher(50 * time.Millisecond)
	m.Record("mx-master-4", "wheel_up", DeviceScrollSettings{Direction: "standard", SmoothScroll: true}, time.Unix(1, 0))

	got, ok := m.Match(nativeScrollEvent{Axis1: 1, At: time.Unix(1, 10*1e6)})
	if !ok || got.Direction != "standard" {
		t.Fatalf("Match = %#v, %v, want matched standard settings", got, ok)
	}
}
```

- [ ] **Step 2: Run the matcher test to verify it fails**

Run: `go test ./internal/platform/macos -run TestMatcherOnlyConsumesRecentMatchingWheelTicks -v`  
Expected: FAIL with undefined matcher types

- [ ] **Step 3: Implement the matcher and emission plan**

```go
type nativeScrollEvent struct {
	Axis1 int
	Axis2 int
	At    time.Time
}
```

- [ ] **Step 4: Implement the Darwin event tap and scroll emitter**

```go
func NewScrollRewriter() ScrollRewriter
func (r *EventTapScrollRewriter) Start(ctx context.Context) error
func (r *EventTapScrollRewriter) Record(deviceID, gesture string, settings DeviceScrollSettings, at time.Time)
```

- [ ] **Step 5: Re-run platform tests**

Run: `go test ./internal/platform/macos -v`  
Expected: PASS

### Task 3: Feed HID wheel gestures into the rewriter and keep behavior scoped to MX Master 4

**Files:**
- Modify: `internal/daemon/runtime.go`
- Modify: `internal/events/types.go`
- Modify: `internal/daemon/runtime_test.go`

- [ ] **Step 1: Write the failing runtime integration test**

```go
func TestRuntimeRecordsWheelGesturesForScrollRewrite(t *testing.T) {
	rewriter := &fakeScrollRewriter{}
	runtime := NewRuntimeWithDependencies(RuntimeDependencies{
		Source: fakeEventSource{events: []events.DeviceEvent{{DeviceID: "mx-master-4", Gesture: "wheel_up"}}},
		ScrollRewriter: rewriter,
		AppResolver: fakeActiveAppResolver{bundleID: "com.apple.finder"},
	})
	_ = runtime.ApplyConfig(sampleRuntimeConfig("global"))

	if err := runtime.Run(context.Background()); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if len(rewriter.records) != 1 {
		t.Fatalf("len(records) = %d, want 1", len(rewriter.records))
	}
}
```

- [ ] **Step 2: Run the focused runtime test to verify it fails**

Run: `go test ./internal/daemon -run TestRuntimeRecordsWheelGesturesForScrollRewrite -v`  
Expected: FAIL with missing scroll rewriter integration

- [ ] **Step 3: Implement wheel-gesture recording without routing those gestures through rule matching**

```go
if runtime.scrollRewriter != nil && isScrollGesture(event.Gesture) {
	runtime.scrollRewriter.Record(event.DeviceID, event.Gesture, runtime.scrollSettingsFor(event.DeviceID), event.At)
	continue
}
```

- [ ] **Step 4: Re-run daemon tests**

Run: `go test ./internal/daemon -v`  
Expected: PASS

### Task 4: Document and verify the user-facing behavior

**Files:**
- Modify: `README.md`
- Modify: `docs/manual-smoke-test.md`
- Modify: `examples/config.toml`

- [ ] **Step 1: Update docs to describe real scroll rewrite behavior**

```md
- `devices.scroll.direction` now flips vertical and horizontal wheel output for MX Master 4 while the daemon is running.
- `devices.scroll.smooth_scroll = true` emits short pixel-scroll bursts instead of line ticks.
```

- [ ] **Step 2: Run full verification**

Run: `go test ./... && go run ./cmd/logi validate --config ./examples/config.toml`  
Expected: PASS
