package appmanagers

import (
	"testing"

	"github.com/gsockets/gsockets/config"
	"github.com/stretchr/testify/assert"
)

func TestNewReturnsErrForInvalidDriver(t *testing.T) {
	config := config.AppManager{Driver: "invalid"}
	appManager, err := New(config)

	assert.Nil(t, appManager, "no app manager instance should be returned")
	assert.NotNil(t, err, "an error must be returned if invalid driver is provide")
	assert.ErrorIs(t, err, ErrInvalidAppManagerDriver, "the error returned must be", ErrInvalidAppManagerDriver)
}

func TestNewReturnsValidConfigAppManager(t *testing.T) {
	config := config.AppManager{Driver: "array"}
	appManager, err := New(config)

	assert.Nil(t, err, "no error should be returned for valid config")
	assert.NotNil(t, appManager, "returned app manager instance should not be nil")

	_, ok := appManager.(*configAppManager)
	assert.True(t, ok, "got invalid app manager implementation")
}
