package server

import (
	"net/http"
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
