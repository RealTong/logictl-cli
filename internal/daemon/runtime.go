package daemon

import (
	"context"
	"sync"

	"github.com/realtong/logi-cli/internal/ipc"
)

type Runtime struct {
	mu     sync.Mutex
	status ipc.Status
}

func NewRuntime() *Runtime {
	return &Runtime{
		status: ipc.Status{
			Running: true,
			Message: "running",
		},
	}
}

func (r *Runtime) Status() ipc.Status {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.status
}

func (r *Runtime) Reload(context.Context) (ipc.Status, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	status := r.status
	status.Message = "reload requested"
	return status, nil
}
