package server

import (
	"github.com/go-chi/chi/v5"
)

func (srv *Server) routes() {
	authMiddleware := NewAuthMiddleware(srv.apps)

	srv.router.Get("/", srv.rootHandler)
	srv.router.Get("/app/{appKey}", srv.serveWs)

	srv.router.Get("/apps/{appId}/channels", srv.allChannels)
	srv.router.Get("/apps/{appId}/channels/:channelName", srv.channelDetails)

	// Authenticated routes
	srv.router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Handler)

		r.Post("/apps/{appId}/events", srv.trigger)
		r.Post("/apps/{appId}/batch_events", srv.triggerBatch)
	})
}
