package daemon

import (
	"context"

	"github.com/realtong/logictl-cli/internal/actions"
	appcore "github.com/realtong/logictl-cli/internal/app"
	"github.com/realtong/logictl-cli/internal/config"
	"github.com/realtong/logictl-cli/internal/events"
	"github.com/realtong/logictl-cli/internal/hidapi"
	"github.com/realtong/logictl-cli/internal/ipc"
	platformmacos "github.com/realtong/logictl-cli/internal/platform/macos"
	"github.com/realtong/logictl-cli/internal/rules"
)

type App struct {
	paths   appcore.Paths
	runtime *Runtime
}

type nativeReportSourceFactoryAdapter struct {
	factory platformmacos.HIDReportSourceFactory
}

func (a nativeReportSourceFactoryAdapter) Validate(spec nativeMatchSpec) error {
	return a.factory.Validate(platformmacos.HIDReportMatch{
		VendorID:     spec.VendorID,
		ProductID:    spec.ProductID,
		UsagePage:    spec.UsagePage,
		Usage:        spec.Usage,
		SerialNumber: spec.SerialNumber,
		Product:      spec.Product,
	})
}

func (a nativeReportSourceFactoryAdapter) Open(spec nativeMatchSpec) events.Source {
	return a.factory.Open(platformmacos.HIDReportMatch{
		VendorID:     spec.VendorID,
		ProductID:    spec.ProductID,
		UsagePage:    spec.UsagePage,
		Usage:        spec.Usage,
		SerialNumber: spec.SerialNumber,
		Product:      spec.Product,
	})
}

func NewApp(paths appcore.Paths) *App {
	shortcutEmitter := actions.AppleScriptShortcutEmitter{}

	return &App{
		paths: paths,
		runtime: NewRuntimeWithDependencies(RuntimeDependencies{
			Source: newMXMaster4EventSource(
				hidapi.NewClient(),
				events.NewHIDSource,
				nativeReportSourceFactoryAdapter{factory: platformmacos.NewHIDReportSourceFactory()},
			),
			AppResolver: platformmacos.NewEnvironment(),
			Executor: actions.Executor{
				ShortcutEmitter: shortcutEmitter,
				SystemRunner:    actions.MacOSSystemRunner{ShortcutEmitter: shortcutEmitter},
				ScriptRunner:    actions.CommandScriptRunner{},
			},
			ScrollRewriter: platformmacos.NewScrollRewriter(),
			ConfigPath:     paths.ConfigFile,
		}),
	}
}

func NewAppWithRuntime(paths appcore.Paths, runtime *Runtime) *App {
	return &App{
		paths:   paths,
		runtime: runtime,
	}
}

func NewFromConfig(cfg *config.Config) (*Runtime, error) {
	runtime := NewRuntimeWithDependencies(RuntimeDependencies{
		BuildMatcher: func(cfg *config.Config) (ruleMatcher, error) {
			return rules.NewEngine(cfg), nil
		},
	})
	if err := runtime.ApplyConfig(cfg); err != nil {
		return nil, err
	}
	return runtime, nil
}

func (a *App) Run(ctx context.Context) error {
	if err := a.runtime.Initialize(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 2)
	go func() {
		errCh <- a.runtime.Run(ctx)
	}()
	go func() {
		errCh <- newServer(a.paths.SocketFile, a.runtime).Run(ctx)
	}()

	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil {
			cancel()
			return err
		}
	}

	return nil
}

func (a *App) Preflight() error {
	return a.runtime.Initialize()
}

func (a *App) Status() (ipc.Status, error) {
	status, err := ipc.QueryStatus(a.paths.SocketFile)
	if err != nil {
		if ipc.IsUnavailable(err) {
			return ipc.Status{Message: "stopped"}, nil
		}
		return ipc.Status{}, err
	}
	return status, nil
}

func (a *App) Reload() (ipc.Status, error) {
	return ipc.RequestReload(a.paths.SocketFile)
}

func (a *App) SocketPath() string {
	return a.paths.SocketFile
}
