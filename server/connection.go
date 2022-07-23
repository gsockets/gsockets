package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gsockets/gsockets"
	"github.com/gsockets/gsockets/log"
)

const (
	// writeWait is the maximum time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// pongWait is the maximum time to read the next message from the peer.
	pongWait = 120 * time.Second

	// pingPeriod is the time on which the server sends the pings to the peers.
	pingPeriod = (pongWait * 6) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 100
)

var (
	newLine = []byte{'\n'}
	space   = []byte{' '}
	ping    = []byte("ping")
)

type connection struct {
	id string

	app *gsockets.App
	ws  *websocket.Conn

	channels gsockets.ChannelManager

	logger log.Logger

	sendCh chan []byte

	closeCh chan struct{}
}

func NewConnection(app *gsockets.App, conn *websocket.Conn, cm gsockets.ChannelManager, logger log.Logger) gsockets.Connection {
	connId := generateConnectionId()
	newConn := &connection{
		id:       connId,
		app:      app,
		ws:       conn,
		channels: cm,
		logger:   logger.With("connection", connId),
		closeCh:  make(chan struct{}),
		sendCh:   make(chan []byte),
	}

	go newConn.readPump()
	go newConn.writePump()

	return newConn
}

// Id returns the unique connection id
func (c *connection) Id() string {
	return c.id
}

// App returns the app to which this connection has been made
func (c *connection) App() *gsockets.App {
	return c.app
}

// Send will send data back to the connected client
func (c *connection) Send(data any) {
	msg, err := json.Marshal(data)
	if err != nil {
		c.logger.Error("msg", "error parsing message to json", "error", err.Error())
		return
	}

	c.sendCh <- msg
}

// Close closes the current connection
func (c *connection) Close() {
	c.channels.UnsubscribeFromAllChannels(c)
	c.closeCh <- struct{}{}
	close(c.sendCh)

	err := c.ws.Close()
	if err != nil {
		c.logger.Error("msg", "error closing websocket connection", "error", err.Error())
	}
}

func (c *connection) readPump() {
	defer func() {
		c.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	_ = c.ws.SetReadDeadline(time.Now().Add(pongWait))

	c.ws.SetPongHandler(func(appData string) error {
		c.logger.Info("msg", "setPongHandler run")
		_ = c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("msg", "unexpected close on websocket connection", "error", err.Error())
				return
			}

			c.logger.Error("msg", "error reading message from the websocket connection", "error", err.Error())
			return
		}

		message = bytes.TrimSpace(bytes.Replace(message, newLine, space, -1))
		c.logger.Info("msg", "received message from the websocket connection", "payload", string(message))
	}
}

func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case msg, ok := <-c.sendCh:
			if !ok {
				_ = c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				c.Close()

				return
			}

			w, err := c.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				c.logger.Error("msg", "error acquiring next writer", "error", err.Error())
				continue
			}

			c.logger.Info("msg", "sending message to the client", "payload", string(msg))

			_, err = w.Write(msg)
			if err != nil {
				c.logger.Error("msg", "error writing message to the connection", "error", err.Error())
			}

			if err = w.Close(); err != nil {
				c.logger.Error("msg", "error closing the writer", "error", err.Error())

				c.Close()
				return
			}
		case <-ticker.C:
			
			_ = c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, ping); err != nil {
				c.logger.Error("msg", "error writing ping message to connection", "error", err.Error())

				c.Close()
				return
			}
		case <-c.closeCh:
			return
		}
	}
}

func generateConnectionId() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d.%d", rand.Intn(1000000000), rand.Intn(99999999999999))
}