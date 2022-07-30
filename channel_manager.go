package gsockets

// ChannelManager interface defines the methods for a channel manager to implement. A channel
// manager is responsible of keeping track of the active channels and all the connected connections
// to the channels in the server.
type ChannelManager interface {
	// GetLocalConnections returns the local connections stored in this instance regardless
	// of the channel.
	GetLocalConnections(appId string) []Connection

	// GetLocalChannels returns all the channels for a specific app for the current instance.
	GetLocalChannels(appId string) []string

	// GetGlobalChannels returns all the channels accross all instances.
	GetGlobalChannels(appId string) []string

	// GetGlobalChannelsWithConnectionCount returns the list of all the channels with number of connections
	// subscribed to them.
	GetGlobalChannelsWithConnectionCount(appId string) map[string]int

	// GetChannelMembers returns all the subscribed user info for a presence channel.
	GetChannelMembers(appId, channelName string) map[string]PresenceMember

	// GetChannelConnectionCount returns the number of connections currently subscribed with the given channel.
	GetChannelConnectionCount(appId, channelName string) int

	// AddConnection adds a connection to this instance.
	AddConnection(appId string, conn Connection)

	// RemoveConnection removes a connection from this instance.
	RemoveConnection(appId string, conn Connection)

	// SetUser associates a connection with an user.
	SetUser(appId, userId, connId string)

	// RemoveUser removes the link between a connection and user.
	RemoveUser(appId, userId, connId string)

	// GetUserConnections returns all the connection associated with a particular user.
	GetUserConnections(appId, userId string) []Connection

	// SubscribeToChannel subscribe a connection to a specific channel.
	SubscribeToChannel(appId, channelName string, conn Connection, payload any)

	// UnsubscribeFromChannel will remove a connection from specific channel.
	UnsubscribeFromChannel(appId, channelName string, conn Connection)

	// UnsubscribeFromAllChannels will unsubscribe the connection accross all the channels it is
	// subscribed to.
	UnsubscribeFromAllChannels(appId, conn string)

	// IsInChannel returns a boolean indicating if a connection is subscribed to a channel.
	IsInChannel(appId string, channel string, conn Connection) bool

	// BroadcastToChannel sends the given data to all the connected clients to the channel.
	BroadcastToChannel(appId, channel string, data any)

	// BroadcastExcept sends the given data to all the channels except the given connection id.
	BroadcastExcept(appId, channel string, data any, connId string)
}
