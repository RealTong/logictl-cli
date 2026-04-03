package daemon

import (
	"context"
	"testing"
	"time"

	"github.com/realtong/logi-cli/internal/config"
	"github.com/realtong/logi-cli/internal/events"
	platformmacos "github.com/realtong/logi-cli/internal/platform/macos"
	"github.com/realtong/logi-cli/internal/rules"
)

type fakeEventSource struct {
	events []events.DeviceEvent
	err    error
}

func (s fakeEventSource) Stream(context.Context) (<-chan events.DeviceEvent, <-chan error) {
	eventsCh := make(chan events.DeviceEvent, len(s.events))
	for _, event := range s.events {
		eventsCh <- event
	}
	close(eventsCh)

	errCh := make(chan error, 1)
	if s.err != nil {
		errCh <- s.err
	}
	close(errCh)
	return eventsCh, errCh
}

type fakeMatcher struct {
	action      config.Action
	err         error
	lastContext rules.Context
	lastEvent   events.DeviceEvent
}

func (m *fakeMatcher) Match(ctx rules.Context, event events.DeviceEvent) (config.Action, error) {
	m.lastContext = ctx
	m.lastEvent = event
	return m.action, m.err
}

type fakeExecutor struct {
	actions []config.Action
}

func (e *fakeExecutor) Execute(_ context.Context, action config.Action) error {
	e.actions = append(e.actions, action)
	return nil
}

type fakeActiveAppResolver struct {
	bundleID string
}

func (r fakeActiveAppResolver) ActiveBundleID(context.Context) (string, error) {
	return r.bundleID, nil
}

type fakeScrollRewriter struct {
	records []scrollRewriteRecord
}

type scrollRewriteRecord struct {
	deviceID string
	gesture  string
	settings config.ScrollConfig
	at       time.Time
}

func (f *fakeScrollRewriter) Start(context.Context) error {
	return nil
}

func (f *fakeScrollRewriter) Record(deviceID, gesture string, settings config.ScrollConfig, at time.Time) {
	f.records = append(f.records, scrollRewriteRecord{
		deviceID: deviceID,
		gesture:  gesture,
		settings: settings,
		at:       at,
	})
}

func TestRuntimeApplyConfigWithoutRestart(t *testing.T) {
	runtime := NewRuntimeWithDependencies(RuntimeDependencies{})
	if err := runtime.ApplyConfig(sampleRuntimeConfig("global")); err != nil {
		t.Fatalf("ApplyConfig returned error: %v", err)
	}
	if err := runtime.ApplyConfig(sampleRuntimeConfig("chrome")); err != nil {
		t.Fatalf("ApplyConfig returned error: %v", err)
	}

	current := runtime.CurrentConfig()
	if current == nil {
		t.Fatal("CurrentConfig returned nil, want loaded config")
	}
	if got, want := current.Profiles[0].ID, "chrome"; got != want {
		t.Fatalf("CurrentConfig profile ID = %q, want %q", got, want)
	}
}

func TestRuntimeDispatchesMatchedAction(t *testing.T) {
	matcher := &fakeMatcher{
		action: config.Action{ID: "close_tab", Type: "shortcut", Keys: []string{"cmd", "w"}},
	}
	executor := &fakeExecutor{}
	runtime := NewRuntimeWithDependencies(RuntimeDependencies{
		Source:      fakeEventSource{events: []events.DeviceEvent{{DeviceID: "mx-master-4", Gesture: "hold(gesture_button)+move(down)"}}},
		Matcher:     matcher,
		Executor:    executor,
		AppResolver: fakeActiveAppResolver{bundleID: "com.google.Chrome"},
	})

	if err := runtime.Run(context.Background()); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if got, want := matcher.lastContext.AppBundleID, "com.google.Chrome"; got != want {
		t.Fatalf("matcher.lastContext.AppBundleID = %q, want %q", got, want)
	}
	if got := len(executor.actions); got != 1 {
		t.Fatalf("len(executor.actions) = %d, want 1", got)
	}
}

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

func TestRuntimeRecordsWheelGesturesForScrollRewrite(t *testing.T) {
	rewriter := &fakeScrollRewriter{}
	runtime := NewRuntimeWithDependencies(RuntimeDependencies{
		Source: fakeEventSource{events: []events.DeviceEvent{
			{DeviceID: "mx-master-4", Gesture: "wheel_up", At: time.Unix(1, 0)},
		}},
		Matcher:        &fakeMatcher{err: rules.ErrNoBinding},
		Executor:       &fakeExecutor{},
		AppResolver:    fakeActiveAppResolver{bundleID: "com.apple.finder"},
		ScrollRewriter: rewriter,
	})

	cfg := sampleRuntimeConfig("global")
	cfg.Devices[0].Scroll = config.ScrollConfig{
		Direction:    "standard",
		SmoothScroll: true,
	}
	if err := runtime.ApplyConfig(cfg); err != nil {
		t.Fatalf("ApplyConfig returned error: %v", err)
	}

	if err := runtime.Run(context.Background()); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if len(rewriter.records) != 1 {
		t.Fatalf("len(rewriter.records) = %d, want 1", len(rewriter.records))
	}
	if got, want := rewriter.records[0].settings, (config.ScrollConfig{Direction: "standard", SmoothScroll: true}); got != want {
		t.Fatalf("rewriter.records[0].settings = %#v, want %#v", got, want)
	}
	if got, want := rewriter.records[0].gesture, "wheel_up"; got != want {
		t.Fatalf("rewriter.records[0].gesture = %q, want %q", got, want)
	}
}

var _ platformmacos.ScrollRewriter = (*fakeScrollRewriter)(nil)

func sampleRuntimeConfig(profileID string) *config.Config {
	return &config.Config{
		Devices: []config.Device{
			{
				ID:             "mx-master-4",
				MatchVendorID:  1133,
				MatchProductID: 0xb042,
			},
		},
		Actions: []config.Action{
			{ID: "close_tab", Type: "shortcut", Keys: []string{"cmd", "w"}},
		},
		Profiles: []config.Profile{
			{
				ID:          profileID,
				AppBundleID: "com.google.Chrome",
				Bindings: []config.Binding{
					{
						Device:  "mx-master-4",
						Trigger: "hold(gesture_button)+move(down)",
						Action:  "close_tab",
					},
				},
			},
		},
	}
}
