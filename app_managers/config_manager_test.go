package appmanagers

import (
	"context"
	"testing"

	"github.com/gsockets/gsockets"
	"github.com/stretchr/testify/assert"
)

func getConfig() []gsockets.App {
	return []gsockets.App{
		{
			ID:                   "1234",
			Key:                  "app-key",
			Secret:               "secret",
			EnableClientMessages: true,
			MaxConnections:       -1,
			MaxEventPayload:      1024,
		},
	}
}

func TestNewConfigAppManagerCreation(t *testing.T) {
	config := getConfig()
	appManager := newConfigAppManager(config)

	assert.NotNil(t, appManager, "must return a valid configAppManager instance")

	_, ok := appManager.(*configAppManager)
	assert.True(t, ok, "returned instance must be of configAppManager type")
}

func TestFindByIdReturnsIfExists(t *testing.T) {
	confg := getConfig()
	manager := newConfigAppManager(confg)

	app, err := manager.FindById(context.Background(), "1234")

	assert.Nil(t, err, "no error should be returned if value exists")
	assert.NotNil(t, app, "returned app should not be nil")

	// as FindById returns a pointer to app, we verify if the values match
	assert.Equal(t, confg[0], *app, "the retuned app is not matching")
}

func TestFindByIdReturnsIfDoesNotExist(t *testing.T) {
	confg := getConfig()
	manager := newConfigAppManager(confg)

	app, err := manager.FindById(context.Background(), "invalid")

	assert.Equal(t, ErrInvalidAppId, err, "ErrInvalidAppId error should be returned if value does not exists")
	assert.Nil(t, app, "returned app should be nil if not found")
}

func TestFindByKeyReturnsIfExists(t *testing.T) {
	confg := getConfig()
	manager := newConfigAppManager(confg)

	app, err := manager.FindByKey(context.Background(), "app-key")

	assert.Nil(t, err, "no error should be returned if value exists")
	assert.NotNil(t, app, "returned app should not be nil")
	assert.Equal(t, confg[0], *app, "the retuned app is not matching")
}

func TestFindByKeyReturnsIfDoesNotExists(t *testing.T) {
	confg := getConfig()
	manager := newConfigAppManager(confg)

	app, err := manager.FindByKey(context.Background(), "invalid-app-key")

	assert.Equal(t, ErrInvalidAppKey, err, "ErrInvalidAppKey error should be returned if value does not exists")
	assert.Nil(t, app, "returned app should be nil if not found")
}

func TestGetAppSecretForValidApp(t *testing.T) {
	confg := getConfig()
	manager := newConfigAppManager(confg)

	secret, err := manager.GetAppSecret(context.Background(), "1234")

	assert.Nil(t, err, "no error should be returned for valid app id")
	assert.Equal(t, "secret", secret, "secret value should match")
}

func TestGetAppSecretForInvalidApp(t *testing.T) {
	confg := getConfig()
	manager := newConfigAppManager(confg)

	secret, err := manager.GetAppSecret(context.Background(), "invalid")

	assert.Nil(t, err, "no error should be returned even if value does not exists")
	assert.Equal(t, "", secret, "secret should return blank string for invalid app")
}
