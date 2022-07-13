package server

import (
	"net/http"
)

func (srv *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	RenderJSON(w, 200, "gsockets server", nil)
}
