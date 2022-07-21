package server

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	appmanagers "github.com/gsockets/gsockets/app_managers"
)

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
		RenderJSON(w, http.StatusInternalServerError, err.Error(), nil)
	}

	newConn := NewConnection(app, conn, nil, srv.logger)
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
