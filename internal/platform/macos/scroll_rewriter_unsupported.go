//go:build !darwin

package macos

import (
	"context"
	"time"

	"github.com/realtong/logi-cli/internal/config"
)

type unsupportedScrollRewriter struct{}

func NewScrollRewriter() ScrollRewriter {
	return unsupportedScrollRewriter{}
}

func (unsupportedScrollRewriter) Start(context.Context) error {
	return nil
}

func (unsupportedScrollRewriter) Record(string, string, config.ScrollConfig, time.Time) {}
