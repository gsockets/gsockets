package gsockets

// Channel interface defines the methods required for a channel to implement.
type Channel interface {
	// Name returns the name of the channel. Channel names are unique in the server.
	Name() string

	// Connections returns all the clients subscribed to this channel.
	Connections() []Connection

	// Subscribe adds a new connection to the channel.
	Subscribe(conn Connection, payload any)

	// Unsubscribe removes a connection from the channel.
	Unsubscribe(conn Connection)

	// IsSubscribed determines if the given connection is already connected to the channel.
	IsSubscribed(conn Connection) bool

	// Broadcast will send the given data to all the subscribed connections to this channel.
	Broadcast(data any)

	// BroadcastExcept sends the data to all the channels except to the connection identified
	// by the given connection id.
	BroadcastExcept(data any, connToExclude string)
}