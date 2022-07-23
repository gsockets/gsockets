package server

import (
	"encoding/json"
	"errors"
	"net/http"

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

func (srv *Server) trigger(w http.ResponseWriter, r *http.Request) {
	appId := chi.URLParam(r, "appId")
	var body gsockets.PusherAPIMessage

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		RenderJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	app, err := srv.apps.FindById(r.Context(), appId)
	if err != nil {
		if errors.Is(err, appmanagers.ErrInvalidAppId) {
			RenderJSON(w, http.StatusBadRequest, err.Error(), nil)
			return
		}

		RenderJSON(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	broadcast := func(channel string) {
		payload := gsockets.PusherSentMessage{
			Event:   body.Name,
			Channel: channel,
			Data:    body.Data,
		}

		if body.SocketId == "" {
			srv.channels.BroadcastToChannel(app.ID, channel, payload)
		} else {
			srv.channels.BroadcastExcept(app.ID, channel, payload, body.SocketId)
		}
	}

	for _, channel := range body.Channels {
		go broadcast(channel)
	}

	RenderJSON(w, http.StatusOK, "", okResponse{Ok: true})
}

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
