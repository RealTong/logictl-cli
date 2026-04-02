package daemon

import (
	"context"
	"testing"

	"github.com/realtong/logi-cli/internal/config"
	"github.com/realtong/logi-cli/internal/events"
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
		Source:      fakeEventSource{events: []events.DeviceEvent{{DeviceID: "mx-master-4", Gesture: "hold(thumb_button)+move(down)"}}},
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
						Trigger: "hold(thumb_button)+move(down)",
						Action:  "close_tab",
					},
				},
			},
		},
	}
}
