package gsockets

type MessageData struct {
	Channel     string `json:"channel"`
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
	Channel  string   `json:"channel"`
	Channels []string `json:"channels"`
	Data     string   `json:"data"`
	SocketId string   `json:"socket_id"`
}

type PusherSentMessage struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
}
