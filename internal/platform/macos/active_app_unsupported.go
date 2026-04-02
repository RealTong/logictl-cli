//go:build !darwin

package macos

import (
	"context"
	"errors"
)

func defaultActiveBundleID(context.Context) (string, error) {
	return "", errors.New("active app lookup is only supported on macOS")
}
