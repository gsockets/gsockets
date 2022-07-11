package gsockets

// App struct represents the main application instance. If you are familier with
// pusher apps, gsockets apps serve the exact same purpose. Each App gets an id,
// key and secret that can be used to authenticate with the gsockets server.
type App struct {
	// ID uniquely identifies a single app.
	ID string

	// Key is the publishable Key that the client libraries can use
	// to connect with this app instance.
	Key string

	// Secret is used to encrypt and decrypt communications from the server SDKs.
	Secret string

	// MaxConnections configures the maximum number of concurrent connections allowed
	// for this app.
	MaxConnections int

	// EnableClientMessages configures whether client side messaging is enabled for
	// this app.
	EnableClientMessages bool

	// MaxEventPayload configures the size of the maximum allowed payload size for events in
	// kilobytes. It applies to both http api and websockets. If the value is negative, there
	// is no payload size restriction.
	MaxEventPayload int
}