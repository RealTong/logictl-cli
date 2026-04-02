package macos

import (
	"context"
	"errors"
)

type Environment struct {
	lookupActiveBundleID func(context.Context) (string, error)
	lookupPermissions    func(context.Context) (PermissionReport, error)
}

func NewEnvironment() Environment {
	return Environment{
		lookupActiveBundleID: defaultActiveBundleID,
		lookupPermissions:    defaultPermissionReport,
	}
}

func (e Environment) ActiveBundleID(ctx context.Context) (string, error) {
	if e.lookupActiveBundleID == nil {
		return "", errors.New("active app lookup is unavailable")
	}
	return e.lookupActiveBundleID(ctx)
}
