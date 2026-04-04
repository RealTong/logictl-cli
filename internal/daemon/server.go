package daemon

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/realtong/logictl-cli/internal/ipc"
)

type controller interface {
	Status() ipc.Status
	Reload(context.Context) (ipc.Status, error)
}

type Server struct {
	socketPath string
	handle     func(context.Context, ipc.Request) (ipc.Status, error)
}

func NewServer(socketPath string, status ipc.Status) *Server {
	return newServer(socketPath, staticController{status: status})
}

func newServer(socketPath string, ctl controller) *Server {
	return &Server{
		socketPath: socketPath,
		handle: func(ctx context.Context, req ipc.Request) (ipc.Status, error) {
			switch req.Command {
			case ipc.CommandStatus:
				return ctl.Status(), nil
			case ipc.CommandReload:
				return ctl.Reload(ctx)
			default:
				return ipc.Status{}, fmt.Errorf("unsupported command %q", req.Command)
			}
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	listener, err := ipc.Listen(s.socketPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = listener.Close()
		_ = os.Remove(s.socketPath)
	}()

	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil || errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}

		go func() {
			_ = ipc.ServeConn(conn, func(req ipc.Request) (ipc.Status, error) {
				return s.handle(ctx, req)
			})
		}()
	}
}

type staticController struct {
	status ipc.Status
}

func (c staticController) Status() ipc.Status {
	return c.status
}

func (c staticController) Reload(context.Context) (ipc.Status, error) {
	status := c.status
	if status.Message == "" {
		status.Message = "reload requested"
	}
	return status, nil
}
