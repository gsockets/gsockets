package gsockets

// Channel interface defines the methods required for a channel to implement.
type Channel interface {
	// Subscribe adds a new connection to the channel.
	Subscribe(appId string, conn Connection, payload MessageData) error

	// Unsubscribe removes a connection from the channel.
	Unsubscribe(appId, channel string, conn Connection) error

	// Broadcast will send the given data to all the subscribed connections to this channel.
	Broadcast(appId, channel string, data any)

	// BroadcastExcept sends the data to all the channels except to the connection identified
	// by the given connection id.
	BroadcastExcept(appId, channel string, data any, connToExclude string)
}
