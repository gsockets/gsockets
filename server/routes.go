package server

func (srv *Server) routes() {
	srv.router.Get("/", srv.rootHandler)
}
