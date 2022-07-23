package gsockets

const (
	// 4000-4009 error codes indicate  connections being closed by the server,
	// and attempting to reconnect using the same parameters will not succeed.
	ERROR_SSL_ONLY                          = 4000
	ERROR_APPLICATION_DOES_NOT_EXIST        = 4001
	ERROR_APPLICATION_DISABLED              = 4003
	ERROR_APPLICATION_OVER_CONNECTION_QUOTA = 4004
	ERROR_PATH_NOT_FOUND                    = 4005
	ERROR_INVALID_VERSION_STRING_FORMAT     = 4006
	ERROR_UNSUPPORTED_PROTOCOL_VERSION      = 4007
	ERROR_NO_PROTOCOL_VERSION_SUPPLIED      = 4008
	ERROR_CONNECTION_IS_UNAUTHORIZED        = 4009

	// 4100-4199 error coedes indicate errors resulting in connection being closed
	// by the server and the client may reconnect after some delay.
	ERROR_OVER_CAPACITY = 4100

	// 4200-4299 indicates an error resulting in the connection being closed by the
	// server and the client may reconnect immediately.
	ERROR_GENERIC_RECONNECT_IMMEDIATELY = 4200
	ERROR_PONG_NOT_RECEIVED             = 4201
	ERROR_CLOSED_AFTER_INACTIVITY       = 4202

	// 4300-4399 any kind of other errors.
	ERROR_CLIENT_EVENT_RATE_LIMIT = 4301
)

type PusherError struct {
	Code    int
	Message string
}

func (pe PusherError) Error() string {
	return pe.Message
}
