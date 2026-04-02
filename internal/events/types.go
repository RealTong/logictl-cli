package events

import (
	"context"
	"time"
)

type RawReport struct {
	DeviceID string
	Bytes    []byte
	At       time.Time
}

type Source interface {
	Stream(ctx context.Context) (<-chan RawReport, <-chan error)
}
