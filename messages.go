package gsockets

type MessageData struct {
	Channel     string `json:"channel"`
	Auth        string `json:"auth"`
	ChannelData string `json:"channel_data"`
	UserData    string `json:"user_data"`
}

type PusherMessage struct {
	Name    string      `jsone:"name"`
	Event   string      `json:"event"`
	Channel string      `json:"channel"`
	Data    MessageData `json:"data"`
}

type PusherAPIMessage struct {
	Name     string   `json:"name"`
	Event    string   `json:"event"`
	Channel  string   `json:'channel"`
	Channels []string `json:"channels"`
	Data     string   `json:"data"`
	SocketId string   `json:"socket_id"`
}

type PusherBatchApiMessage struct {
	Batch []PusherAPIMessage `json:"batch"`
}

type PusherSentMessage struct {
	Event   string `json:"event"`
	Channel string `json:"channel,omitempty"`
	Data    any    `json:"data"`
}

func NewPusherError(errorEvent, message, channel string, code int) PusherSentMessage {
	data := struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}{
		Message: message,
		Code:    code,
	}

	return PusherSentMessage{Event: errorEvent, Channel: channel, Data: data}
}
