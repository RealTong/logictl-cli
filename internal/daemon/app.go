package daemon

import (
	"context"

	"github.com/realtong/logi-cli/internal/actions"
	appcore "github.com/realtong/logi-cli/internal/app"
	"github.com/realtong/logi-cli/internal/events"
	"github.com/realtong/logi-cli/internal/hidapi"
	"github.com/realtong/logi-cli/internal/ipc"
	platformmacos "github.com/realtong/logi-cli/internal/platform/macos"
)

type App struct {
	paths   appcore.Paths
	runtime *Runtime
}

func NewApp(paths appcore.Paths) *App {
	shortcutEmitter := actions.AppleScriptShortcutEmitter{}

	return &App{
		paths: paths,
		runtime: NewRuntimeWithDependencies(RuntimeDependencies{
			Source:      newMXMaster4EventSource(hidapi.NewClient(), events.NewHIDSource),
			AppResolver: platformmacos.NewEnvironment(),
			Executor: actions.Executor{
				ShortcutEmitter: shortcutEmitter,
				SystemRunner:    actions.MacOSSystemRunner{ShortcutEmitter: shortcutEmitter},
				ScriptRunner:    actions.CommandScriptRunner{},
			},
			ConfigPath: paths.ConfigFile,
		}),
	}
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
