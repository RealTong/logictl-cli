package daemon

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/realtong/logictl-cli/internal/config"
	"github.com/realtong/logictl-cli/internal/events"
	"github.com/realtong/logictl-cli/internal/ipc"
	platformmacos "github.com/realtong/logictl-cli/internal/platform/macos"
	"github.com/realtong/logictl-cli/internal/rules"
)

type eventSource interface {
	Stream(context.Context) (<-chan events.DeviceEvent, <-chan error)
}

type validatingEventSource interface {
	Validate() error
}

type appResolver interface {
	ActiveBundleID(context.Context) (string, error)
}

type ruleMatcher interface {
	Match(rules.Context, events.DeviceEvent) (config.Action, error)
}

type actionExecutor interface {
	Execute(context.Context, config.Action) error
}

type matcherBuilder func(*config.Config) (ruleMatcher, error)
type configLoader func(string) (*config.Config, error)

type RuntimeDependencies struct {
	Source         eventSource
	AppResolver    appResolver
	Matcher        ruleMatcher
	Executor       actionExecutor
	ScrollRewriter platformmacos.ScrollRewriter
	ConfigPath     string
	LoadConfig     configLoader
	BuildMatcher   matcherBuilder
}

type Runtime struct {
	mu             sync.RWMutex
	status         ipc.Status
	source         eventSource
	appResolver    appResolver
	matcher        ruleMatcher
	executor       actionExecutor
	scrollRewriter platformmacos.ScrollRewriter
	configPath     string
	loadConfig     configLoader
	buildMatcher   matcherBuilder
	currentConfig  *config.Config
	scrollSettings map[string]config.ScrollConfig
}

func NewRuntime() *Runtime {
	return NewRuntimeWithDependencies(RuntimeDependencies{})
}

func NewRuntimeWithDependencies(deps RuntimeDependencies) *Runtime {
	loadConfig := deps.LoadConfig
	if loadConfig == nil {
		loadConfig = config.LoadFile
	}

	buildMatcher := deps.BuildMatcher
	if buildMatcher == nil {
		buildMatcher = func(cfg *config.Config) (ruleMatcher, error) {
			return rules.NewEngine(cfg), nil
		}
	}

	return &Runtime{
		source:         deps.Source,
		appResolver:    deps.AppResolver,
		matcher:        deps.Matcher,
		executor:       deps.Executor,
		scrollRewriter: deps.ScrollRewriter,
		configPath:     deps.ConfigPath,
		loadConfig:     loadConfig,
		buildMatcher:   buildMatcher,
		scrollSettings: map[string]config.ScrollConfig{},
		status: ipc.Status{
			Running: true,
			Message: "running",
		},
	}
}

func (r *Runtime) Status() ipc.Status {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.status
}

func (r *Runtime) CurrentConfig() *config.Config {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.currentConfig
}

func (r *Runtime) ScrollSettings(deviceID string) config.ScrollConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.scrollSettings[deviceID]
}

func (r *Runtime) ApplyConfig(cfg *config.Config) error {
	if err := config.Validate(cfg); err != nil {
		return err
	}

	matcher, err := r.buildMatcher(cfg)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.currentConfig = cfg
	r.matcher = matcher
	r.scrollSettings = buildScrollSettings(cfg)
	return nil
}

func (r *Runtime) Initialize() error {
	if source, ok := r.source.(validatingEventSource); ok {
		if err := source.Validate(); err != nil {
			return err
		}
	}

	if r.configPath == "" {
		return nil
	}

	cfg, err := r.loadConfig(r.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return r.ApplyConfig(cfg)
}

func (r *Runtime) Run(ctx context.Context) error {
	if r.source == nil || r.matcher == nil || r.executor == nil || r.appResolver == nil {
		<-ctx.Done()
		return nil
	}

	eventsCh, errs := r.source.Stream(ctx)
	rewriterErrs := startScrollRewriter(ctx, r.scrollRewriter)
	for eventsCh != nil || errs != nil {
		select {
		case event, ok := <-eventsCh:
			if !ok {
				eventsCh = nil
				continue
			}

			if r.scrollRewriter != nil && isScrollGesture(event.Gesture) {
				settings := r.scrollSettingsFor(event.DeviceID)
				if settings != (config.ScrollConfig{}) && shouldRewriteScroll(settings) {
					r.scrollRewriter.Record(event.DeviceID, event.Gesture, settings, event.At)
				}
			}

			appBundleID, err := r.appResolver.ActiveBundleID(ctx)
			if err != nil {
				return err
			}

			matcher := r.currentMatcher()
			if matcher == nil {
				continue
			}

			action, err := matcher.Match(rules.Context{AppBundleID: appBundleID}, event)
			if err != nil {
				if errors.Is(err, rules.ErrNoBinding) || errors.Is(err, rules.ErrAmbiguousBinding) {
					continue
				}
				return err
			}

			if err := r.executor.Execute(ctx, action); err != nil {
				return err
			}
		case err, ok := <-errs:
			if !ok {
				errs = nil
				continue
			}
			if err != nil {
				return err
			}
		case err, ok := <-rewriterErrs:
			if !ok {
				rewriterErrs = nil
				continue
			}
			if err != nil && !errors.Is(err, context.Canceled) {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}

	return nil
}

func (r *Runtime) Reload(context.Context) (ipc.Status, error) {
	if r.configPath != "" {
		cfg, err := r.loadConfig(r.configPath)
		if err != nil {
			if !os.IsNotExist(err) {
				return ipc.Status{}, err
			}
		} else {
			if err := r.ApplyConfig(cfg); err != nil {
				return ipc.Status{}, err
			}
		}
	}

	status := r.Status()
	status.Message = "reload requested"
	return status, nil
}

func (r *Runtime) currentMatcher() ruleMatcher {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.matcher
}

func (r *Runtime) scrollSettingsFor(deviceID string) config.ScrollConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.scrollSettings[deviceID]
}

func buildScrollSettings(cfg *config.Config) map[string]config.ScrollConfig {
	settings := make(map[string]config.ScrollConfig)
	if cfg == nil {
		return settings
	}
	for _, device := range cfg.Devices {
		settings[device.ID] = device.Scroll
	}
	return settings
}

func isScrollGesture(gesture string) bool {
	switch gesture {
	case "wheel_up", "wheel_down", "thumb_wheel_left", "thumb_wheel_right":
		return true
	default:
		return false
	}
}

func shouldRewriteScroll(settings config.ScrollConfig) bool {
	return settings.SmoothScroll || settings.Direction == "standard"
}

func startScrollRewriter(ctx context.Context, rewriter platformmacos.ScrollRewriter) <-chan error {
	if rewriter == nil {
		return nil
	}

	errs := make(chan error, 1)
	go func() {
		defer close(errs)
		errs <- rewriter.Start(ctx)
	}()
	return errs
}
