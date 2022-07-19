package gsockets

// Connection interface defines the methods for interacting with a connection to 
// the gsockets server. This interface is used to abstract away the underlying
// low level websocket connection to the server.
type Connection interface {
	// Id returns the unique connection id
	Id() string

	// App returns the app to which this connection has been made
	App() *App

	// Send will send data back to the connected client
	Send(data any)

	// Close closes the current connection
	Close()
}