//go:build !darwin

package macos

import (
	"context"
	"errors"
)

func defaultPermissionReport(context.Context) (PermissionReport, error) {
	return PermissionReport{}, errors.New("permission inspection is only supported on macOS")
}
