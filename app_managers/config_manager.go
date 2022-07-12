package appmanagers

import (
	"context"

	"github.com/gsockets/gsockets"
)

type configAppManager struct {
	apps map[string]*gsockets.App
}

func newConfigAppManager(appsConfig []gsockets.App) gsockets.AppManager {
	apps := make(map[string]*gsockets.App)
	for _, app := range appsConfig {
		apps[app.ID] = &app
	}

	return &configAppManager{apps: apps}
}

// FindById returns an app instance by the app id.
func (config *configAppManager) FindById(ctx context.Context, id string) (*gsockets.App, error) {
	app, ok := config.apps[id]
	if !ok {
		return nil, nil
	}

	return app, nil
}

// FindByKey returns an app instance by app key.
func (config *configAppManager) FindByKey(ctx context.Context, key string) (*gsockets.App, error) {
	for _, app := range config.apps {
		if app.Key == key {
			return app, nil
		}
	}

	return nil, nil
}

// FindBySecret returns an app instance by app secret.
func (config *configAppManager) GetAppSecret(ctx context.Context, id string) (string, error) {
	app, _ := config.FindById(ctx, id)

	if app == nil {
		return "", nil
	}

	return app.Secret, nil
}
