package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gsockets/gsockets"
	appmanagers "github.com/gsockets/gsockets/app_managers"
	"github.com/gsockets/gsockets/config"
	"github.com/gsockets/gsockets/log"
)

func New(config config.Config, logger log.Logger) *Server {
	return &Server{closing: false, config: config, router: chi.NewRouter(), logger: logger.With("module", "server")}
}

// Server struct is the gsockets server.
type Server struct {
	// closing is set to true when the server process is shutting down.
	// No new connection is accepted when the server is in closing state.
	closing bool

	logger     log.Logger
	config     config.Config
	apps       gsockets.AppManager
	httpServer *http.Server
	router     chi.Router
}

func (srv *Server) Start() error {
	err := srv.initiate()
	if err != nil {
		return err
	}

	srv.logger.Info("msg", "http server started listening for requests", "port", srv.config.Server.Port)

	err = srv.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) Stop() {
	srv.closing = true

	shutdownCtx, _ := context.WithTimeout(context.Background(), time.Second*10)
	srv.logger.Info("msg", "shutdown sequence initiated", "graceful_timeout", 10)

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			srv.logger.Fatal("msg", "graceful shutdown timed out, exiting forcefully Goodbye!!")
		}
	}()

	err := srv.httpServer.Shutdown(shutdownCtx)

	if err != nil {
		srv.logger.Fatal("msg", "error shutting down the http server", "error", err.Error())
	}
}

func (srv *Server) initiate() error {
	apps, err := appmanagers.New(srv.config.AppManager)
	if err != nil {
		return err
	}

	srv.apps = apps
	srv.routes()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", srv.config.Server.Port),
		Handler: srv.router,
	}

	srv.httpServer = httpServer

	return nil
}
