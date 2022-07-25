package gsockets

import (
	"encoding/json"
	"strings"
)

type MessageData struct {
	Channel     string `json:"channel,omitempty"`
	Auth        string `json:"auth,omitempty"`
	ChannelData string `json:"channel_data,omitempty"`
	UserData    string `json:"user_data,omitempty"`
}

type PusherMessage struct {
	Name    string          `jsone:"name"`
	Event   string          `json:"event"`
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data"`
}

func (p PusherMessage) IsClientEvent() bool {
	return strings.HasPrefix(p.Event, "client-")
}

type PusherAPIMessage struct {
	Name     string   `json:"name"`
	Event    string   `json:"event"`
	Channel  string   `json:"channel"`
	Channels []string `json:"channels"`
	Data     string   `json:"data"`
	SocketId string   `json:"socket_id"`
}

type PusherBatchApiMessage struct {
	Batch []PusherAPIMessage `json:"batch"`
}

type ChannelResponse struct {
	SubscriptionCount int  `json:"subscription_count,omitempty"`
	UserCount         int  `json:"user_count,omitempty"`
	Occupied          bool `json:"occupied"`
}

type ChannelListResponse struct {
	Channels map[string]ChannelResponse `json:"channels"`
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
