package macos

import (
	"context"
	"errors"
)

type PermissionReport struct {
	AccessibilityGranted   bool
	InputMonitoringGranted bool
}

func (e Environment) Permissions(ctx context.Context) (PermissionReport, error) {
	if e.lookupPermissions == nil {
		return PermissionReport{}, errors.New("permission lookup is unavailable")
	}
	return e.lookupPermissions(ctx)
}
