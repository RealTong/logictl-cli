package daemon

import (
	"context"

	appcore "github.com/realtong/logi-cli/internal/app"
	"github.com/realtong/logi-cli/internal/ipc"
)

type App struct {
	paths   appcore.Paths
	runtime *Runtime
}

func NewApp(paths appcore.Paths) *App {
	return &App{
		paths:   paths,
		runtime: NewRuntime(),
	}
}

func (a *App) Run(ctx context.Context) error {
	return newServer(a.paths.SocketFile, a.runtime).Run(ctx)
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
