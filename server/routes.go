package server

import (
	"github.com/go-chi/chi/v5"
)

func (srv *Server) routes() {
	authMiddleware := NewAuthMiddleware(srv.apps)

	srv.router.Get("/", srv.rootHandler)
	srv.router.Get("/app/{appKey}", srv.serveWs)

	// Authenticated routes
	srv.router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Handler)

		r.Post("/apps/{appId}/events", srv.trigger)
		r.Post("/apps/{appId}/batch_events", srv.triggerBatch)
		r.Get("/apps/{appId}/channels", srv.allChannels)
		r.Get("/apps/{appId}/channels/{channelName}", srv.channelDetails)
		r.Get("/apps/{appId}/channels/{channelName}/users", srv.channelMembers)
	})
}
