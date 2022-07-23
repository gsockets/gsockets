package channelmanagers

import (
	"errors"

	"github.com/gsockets/gsockets"
	"github.com/gsockets/gsockets/config"
)

var (
	ErrInvalidChannelManagerDriver = errors.New("invalid channel manager driver")
)

func New(config config.ChannelManager) (gsockets.ChannelManager, error) {
	switch config.Driver {
	case "local":
		return newLocalChannelManager(), nil
	default:
		return nil, ErrInvalidChannelManagerDriver
	}
}
