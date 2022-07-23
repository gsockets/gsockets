package server

func (srv *Server) routes() {
	srv.router.Get("/", srv.rootHandler)
	srv.router.Get("/app/{appKey}", srv.serveWs)
	srv.router.Post("/apps/{appId}/events", srv.trigger)
}
