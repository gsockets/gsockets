package gsockets

// Connection interface defines the methods for interacting with a connection to
// the gsockets server. This interface is used to abstract away the underlying
// low level websocket connection to the server.
type Connection interface {
	// Id returns the unique connection id
	Id() string

	// Presence returns the current presence channel memberships for this connection.
	// The Connection implementaion should make sure then membership is added and removed
	// in a concurrent safe manner.
	Presence() map[string]PresenceMember

	// GetPresence returns the presence member info for a given presence channel, will return
	// false if no info found.
	GetPresence(channelName string) (PresenceMember, bool)

	// SetPresence sets a new presence channel subscription.
	SetPresence(channelName string, member PresenceMember)

	// RemovePresence removes subscription for a presence channel.
	RemovePresence(channelName string)

	// App returns the app to which this connection has been made
	App() *App

	// Send will send data back to the connected client
	Send(data any)

	// Close closes the current connection
	Close()
}
