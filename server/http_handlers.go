package server

import (
	"net/http"
)

func (srv *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	RenderSuccess(w, "gsockets server", 200)
}
