package channels

import "github.com/gsockets/gsockets"

type publicChannel struct {
	channelManager gsockets.ChannelManager
}

func newPublicChannel(cm gsockets.ChannelManager) gsockets.Channel {
	return &publicChannel{channelManager: cm}
}

// Subscribe adds a new connection to the channel.
func (c *publicChannel) Subscribe(appId string, conn gsockets.Connection, payload gsockets.MessageData) error {
	c.channelManager.SubscribeToChannel(appId, payload.Channel, conn, payload)
	return nil
}

// Unsubscribe removes a connection from the channel.
func (c *publicChannel) Unsubscribe(appId, channel string, conn gsockets.Connection) error {
	c.channelManager.UnsubscribeFromChannel(appId, channel, conn)
	return nil
}

// Broadcast will send the given data to all the subscribed connections to this channel.
func (c *publicChannel) Broadcast(appId, channel string, data any) {
	c.channelManager.BroadcastToChannel(appId, channel, data)
}

// BroadcastExcept sends the data to all the channels except to the connection identified
// by the given connection id.
func (c *publicChannel) BroadcastExcept(appId, channel string, data any, connToExclude string) {
	c.channelManager.BroadcastExcept(appId, channel, data, connToExclude)
}
