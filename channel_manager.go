package gsockets

// ChannelManager interface defines the methods for a channel manager to implement. A channel
// manager is responsible of keeping track of the active channels and all the connected connections
// to the channels in the server.
type ChannelManager interface {
	// Find finds a channel by the given app and channel name.
	Find(appId, channelName string) (Channel, error)

	// FindOrCreate returns the channel with the given name and app if exists, otherwise
	// it will create a new app and return it.
	FindOrCreate(appId, channelName string) (Channel, error)

	// GetLocalConnections returns the local connections stored in this instance regardless
	// of the channel.
	GetLocalConnections() ([]Connection, error)

	// GetLocalChannels returns all the channels for a specific app for the current instance.
	GetLocalChannels(appId string) ([]Channel, error)

	// SubscribeToChannel subscribe a connection to a specific channel.
	SubscribeToChannel(conn Connection, appId, channelName string, payload any) error

	// UnsubscribeFromChannel will remove a connection from specific channel.
	UnsubscribeFromChannel(conn Connection, appId, channelName string, payload any) error

	// UnsubscribeFromAllChannels will unsubscribe the connection accross all the channels it is
	// subscribed to.
	UnsubscribeFromAllChannels(conn Connection)
}
