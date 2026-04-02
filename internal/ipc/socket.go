package ipc

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"path/filepath"
)

const (
	CommandStatus = "status"
	CommandReload = "reload"
)

type Request struct {
	Command string `json:"command"`
}

type Status struct {
	Running bool   `json:"running"`
	Message string `json:"message"`
}

type response struct {
	Status Status `json:"status"`
	Error  string `json:"error,omitempty"`
}

func Listen(socketPath string) (net.Listener, error) {
	if err := os.MkdirAll(filepath.Dir(socketPath), 0o755); err != nil {
		return nil, err
	}
	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return net.Listen("unix", socketPath)
}

func ServeConn(conn net.Conn, handler func(Request) (Status, error)) error {
	defer conn.Close()

	var req Request
	if err := json.NewDecoder(conn).Decode(&req); err != nil {
		return err
	}

	status, err := handler(req)
	resp := response{Status: status}
	if err != nil {
		resp.Error = err.Error()
	}

	return json.NewEncoder(conn).Encode(resp)
}

func QueryStatus(socketPath string) (Status, error) {
	return send(socketPath, Request{Command: CommandStatus})
}

func RequestReload(socketPath string) (Status, error) {
	return send(socketPath, Request{Command: CommandReload})
}

func send(socketPath string, req Request) (Status, error) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return Status{}, err
	}
	defer conn.Close()

	if err := json.NewEncoder(conn).Encode(req); err != nil {
		return Status{}, err
	}

	var resp response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return Status{}, err
	}
	if resp.Error != "" {
		return resp.Status, errors.New(resp.Error)
	}

	return resp.Status, nil
}
