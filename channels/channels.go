package channels

import (
	"strings"

	"github.com/gsockets/gsockets"
)

func New(name string, cm gsockets.ChannelManager) gsockets.Channel {
	if strings.HasPrefix(name, "private-") {
		return newPrivateChannel(cm)
	} else if strings.HasPrefix(name, "presence-") {
		return newPresenceChannel(cm)
	}

	return newPublicChannel(cm)
}
