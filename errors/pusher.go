package errors

const (
	// 4000-4009 error codes indicate  connections being closed by the server,
	// and attempting to reconnect using the same parameters will not succeed.
	PUSHER_SSL_ONLY=4000
	PUSHER_APPLICATION_DOES_NOT_EXIST=4001
	PUSHER_APPLICATION_DISABLED=4003
	PUSHER_APPLICATION_OVER_CONNECTION_QUOTA=4004
	PUSHER_PATH_NOT_FOUND=4005
	PUSHER_INVALID_VERSION_STRING_FORMAT=4006
	PUSHER_UNSUPPORTED_PROTOCOL_VERSION=4007
	PUSHER_NO_PROTOCOL_VERSION_SUPPLIED=4008
	PUSHER_CONNECTION_IS_UNAUTHORIZED=4009

	// 4100-4199 error coedes indicate errors resulting in connection being closed
	// by the server and the client may reconnect after some delay.
	PUSHER_OVER_CAPACITY=4100

	// 4200-4299 indicates an error resulting in the connection being closed by the
	// server and the client may reconnect immediately.
	PUSHER_GENERIC_RECONNECT_IMMEDIATELY=4200
	PUSHER_PONG_NOT_RECEIVED=4201
	PUSHER_CLOSED_AFTER_INACTIVITY=4202

	// 4300-4399 any kind of other errors.
	PUSHER_CLIENT_EVENT_RATE_LIMIT=4301
)