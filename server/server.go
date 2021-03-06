package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/gsockets/gsockets"
	appmanagers "github.com/gsockets/gsockets/app_managers"
	channelmanagers "github.com/gsockets/gsockets/channel_managers"
	"github.com/gsockets/gsockets/config"
	"github.com/gsockets/gsockets/log"
	"github.com/oklog/ulid/v2"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func New(config config.Config, logger log.Logger) *Server {
	serverId := ulid.Make().String()
	return &Server{
		id:      serverId,
		closing: false,
		config:  config,
		router:  chi.NewRouter(),
		logger:  logger.With("module", "server", "server_id", serverId),
	}
}

// Server struct is the gsockets server.
type Server struct {
	// id is the unique id for this server instance.
	id string

	// closing is set to true when the server process is shutting down.
	// No new connection is accepted when the server is in closing state.
	closing bool

	apps     gsockets.AppManager
	channels gsockets.ChannelManager

	logger     log.Logger
	config     config.Config
	httpServer *http.Server
	router     chi.Router
}

func (srv *Server) Id() string {
	return srv.id
}

func (srv *Server) Start() error {
	err := srv.initiate()
	if err != nil {
		return err
	}

	srv.logger.Info("msg", "http server started listening for requests", "port", srv.config.Server.Port, "server_id", srv.id)

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

	cm, err := channelmanagers.New(srv.config.ChannelManager)
	if err != nil {
		return err
	}

	srv.apps = apps
	srv.channels = cm

	srv.routes()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", srv.config.Server.Port),
		Handler: srv.router,
	}

	srv.httpServer = httpServer

	return nil
}
