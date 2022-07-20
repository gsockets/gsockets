package appmanagers

import (
	"errors"

	"github.com/gsockets/gsockets"
	"github.com/gsockets/gsockets/config"
)

var (
	ErrInvalidAppManagerDriver = errors.New("invalid driver for app manager")
	ErrInvalidAppKey = errors.New("invalid appKey, no app exists with the given appKey")
	ErrInvalidAppId = errors.New("invalid app id, no app exists with the given id")
)

func New(appManagerConfig config.AppManager) (gsockets.AppManager, error) {
	switch appManagerConfig.Driver {
	case "array":
		return newConfigAppManager(appManagerConfig.Array), nil
	default:
		return nil, ErrInvalidAppManagerDriver
	}
}
