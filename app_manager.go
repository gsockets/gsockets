package gsockets

import "context"

// AppManager interface defines the methods to work with apps.
type AppManager interface {
	// FindById returns an app instance by the app id.
	FindById(ctx context.Context, id string) (*App, error)

	// FindByKey returns an app instance by app key.
	FindByKey(ctx context.Context, key string) (*App, error)

	// GetAppSecret returns an app secret for the given app id.
	GetAppSecret(ctx context.Context, id string) (string, error)
}
