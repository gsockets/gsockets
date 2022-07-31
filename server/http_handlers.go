package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gsockets/gsockets"
	appmanagers "github.com/gsockets/gsockets/app_managers"
)

type okResponse struct {
	Ok bool `json:"ok"`
}

func (srv *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Message  string `json:"message"`
		ServerId string `json:"server_id"`
	}{
		ServerId: srv.id,
		Message:  "Welcome to the Gsockets server",
	}

	RenderJSON(w, 200, "", resp)
}

// trigger handles the api endpoint that accepts events from the application backends and distributes
// them to the channels backend to be delivered to the subscirbed clients.
func (srv *Server) trigger(w http.ResponseWriter, r *http.Request) {
	appId := chi.URLParam(r, "appId")
	var body gsockets.PusherAPIMessage

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		RenderJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	go srv.broadcast(appId, body)

	RenderJSON(w, http.StatusOK, "", okResponse{Ok: true})
}

// triggerBatch works similar to the trigger endpoint, the only difference is instead of a single
// event, this endpoint accepts a batch of events.
func (srv *Server) triggerBatch(w http.ResponseWriter, r *http.Request) {
	appId := chi.URLParam(r, "appId")
	var body gsockets.PusherBatchApiMessage

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		RenderJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	for _, msg := range body.Batch {
		if msg.Channel != "" {
			msg.Channels = append(msg.Channels, msg.Channel)
		}

		go srv.broadcast(appId, msg)
	}

	RenderJSON(w, http.StatusOK, "", okResponse{Ok: true})
}

// allChannels returns all the active channels in the server along with how many connections are subscirbed
// to each of those channels.
func (srv *Server) allChannels(w http.ResponseWriter, r *http.Request) {
	channels := srv.channels.GetGlobalChannelsWithConnectionCount(chi.URLParam(r, "appId"))
	filter := r.URL.Query().Get("filter_by_prefix")

	listResponse := make(map[string]gsockets.ChannelResponse)

	for channel, count := range channels {
		if filter != "" && !strings.HasPrefix(channel, filter) {
			continue
		}

		listResponse[channel] = gsockets.ChannelResponse{SubscriptionCount: count, Occupied: count > 0}
	}

	RenderJSON(w, http.StatusOK, "", gsockets.ChannelListResponse{Channels: listResponse})
}

// channelDetails returns details about a single channel.
func (srv *Server) channelDetails(w http.ResponseWriter, r *http.Request) {
	appId := chi.URLParam(r, "appId")
	channelName := chi.URLParam(r, "channelName")

	count := srv.channels.GetChannelConnectionCount(appId, channelName)
	resp := gsockets.ChannelResponse{
		SubscriptionCount: count,
		Occupied:          count > 0,
	}

	RenderJSON(w, http.StatusOK, "", resp)
}

// channelMembers returns all the users subscribed to a persence channel.
func (srv *Server) channelMembers(w http.ResponseWriter, r *http.Request) {
	channelName := chi.URLParam(r, "channelName")
	if !strings.HasPrefix(channelName, "presence-") {
		RenderJSON(w, http.StatusBadRequest, "The channel must be presence channel", nil)
		return
	}

	channelMembers := srv.channels.GetChannelMembers(chi.URLParam(r, "appId"), channelName)
	members := make([]gsockets.ChannelMember, 0)

	for userId := range channelMembers {
		members = append(members, gsockets.ChannelMember{Id: userId})
	}

	resp := gsockets.ChannelMemberResponse{Users: members}
	RenderJSON(w, http.StatusOK, "", resp)
}

// terminateUserConnections will disconnect all the connection from a particular user.
func (srv *Server) terminateUserConnections(w http.ResponseWriter, r *http.Request) {
	srv.channels.TerminateUserConnections(chi.URLParam(r, "appId"), chi.URLParam(r, "userId"))
	res := okResponse{Ok: true}

	RenderJSON(w, http.StatusOK, "", res)
}

// broadcast distributes the messages to the channels backend. The message payload should be validated before
// calling broadcast, it doesn't do any validation or sanity checks, just pushes the message to channels.
func (srv *Server) broadcast(appId string, msg gsockets.PusherAPIMessage) {
	for _, channel := range msg.Channels {
		payload := gsockets.PusherSentMessage{
			Event:   msg.Name,
			Channel: channel,
			Data:    msg.Data,
		}

		if msg.SocketId == "" {
			srv.channels.BroadcastToChannel(appId, channel, payload)
		} else {
			srv.channels.BroadcastExcept(appId, channel, payload, msg.SocketId)
		}
	}
}

// serveWs handles the incoming websocket connections. After doing some validations, serveWs upgrades the connection
// to the websocket protocall and hands it off to the channels backend.
func (srv *Server) serveWs(w http.ResponseWriter, r *http.Request) {
	appKey := chi.URLParam(r, "appKey")
	if appKey == "" {
		srv.logger.Error("msg", "no appKey found on the url")
		RenderJSON(w, http.StatusForbidden, "appKey is required", nil)
		return
	}

	app, err := srv.apps.FindByKey(r.Context(), appKey)
	if err != nil {
		if errors.Is(err, appmanagers.ErrInvalidAppKey) {
			srv.logger.Error("msg", "could not fetch app details", "error", err.Error(), "app_key", appKey)
			RenderJSON(w, http.StatusForbidden, err.Error(), nil)
		} else {
			srv.logger.Error("msg", "error fetching app details", "error", err.Error())
			RenderJSON(w, http.StatusInternalServerError, err.Error(), nil)
		}

		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		srv.logger.Error("msg", "error upgrading to websocket connection", "error", err.Error())
		return
	}

	newConn := NewConnection(app, conn, srv.channels, srv.logger)
	srv.channels.AddConnection(app.ID, newConn)

	srv.logger.Info("msg", "received new connection", "connection", newConn.Id())

	resp := struct {
		Event string `json:"event"`
		Data  struct {
			SocketId        string `json:"socket_id"`
			ActivityTimeout int    `json:"activity_timeout"`
		} `json:"data"`
	}{
		Event: "pusher:connection_established",
		Data: struct {
			SocketId        string "json:\"socket_id\""
			ActivityTimeout int    "json:\"activity_timeout\""
		}{SocketId: newConn.Id(), ActivityTimeout: 120},
	}

	newConn.Send(resp)
}
