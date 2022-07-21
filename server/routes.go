package server

func (srv *Server) routes() {
	srv.router.Get("/", srv.rootHandler)
	srv.router.Get("/app/{appKey}", srv.serveWs)
}
