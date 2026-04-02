package daemon

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/realtong/logi-cli/internal/config"
	"github.com/realtong/logi-cli/internal/events"
	"github.com/realtong/logi-cli/internal/ipc"
	"github.com/realtong/logi-cli/internal/rules"
)

type eventSource interface {
	Stream(context.Context) (<-chan events.DeviceEvent, <-chan error)
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
	Source       eventSource
	AppResolver  appResolver
	Matcher      ruleMatcher
	Executor     actionExecutor
	ConfigPath   string
	LoadConfig   configLoader
	BuildMatcher matcherBuilder
}

type Runtime struct {
	mu            sync.RWMutex
	status        ipc.Status
	source        eventSource
	appResolver   appResolver
	matcher       ruleMatcher
	executor      actionExecutor
	configPath    string
	loadConfig    configLoader
	buildMatcher  matcherBuilder
	currentConfig *config.Config
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
		source:       deps.Source,
		appResolver:  deps.AppResolver,
		matcher:      deps.Matcher,
		executor:     deps.Executor,
		configPath:   deps.ConfigPath,
		loadConfig:   loadConfig,
		buildMatcher: buildMatcher,
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
	return nil
}

func (r *Runtime) Initialize() error {
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
	for eventsCh != nil || errs != nil {
		select {
		case event, ok := <-eventsCh:
			if !ok {
				eventsCh = nil
				continue
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
