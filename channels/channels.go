package channels

import (
	"strings"

	"github.com/gsockets/gsockets"
)

func New(name string, cm gsockets.ChannelManager) gsockets.Channel {
	if strings.HasPrefix(name, "private-") {
		return newPrivateChannel(cm)
	}

	return newPublicChannel(cm)
}
