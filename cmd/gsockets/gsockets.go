package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gsockets/gsockets/config"
	"github.com/gsockets/gsockets/log"
	"github.com/gsockets/gsockets/server"
)

func main() {
	logger := log.New().With("module", "main")
	config, err := config.Load("./../../")
	if err != nil {
		logger.Fatal(err)
		return
	}

	server := server.New(config, logger)

	serverCtx, serverCancel := context.WithCancel(context.Background())

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-exit
		logger.Info("msg", "received inturrupt signal, stopping the server.")

		go func() {
			<-exit
			logger.Fatal("msg", "forcefully exiting, goodbye!!")
		}()

		server.Stop()
		serverCancel()
	}()

	err = server.Start()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}

	<-serverCtx.Done()
	logger.Info("msg", "shutdown complete, goodbye!!")
}
