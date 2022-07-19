package server

import (
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
)

var (
	newLine = []byte{'\n'}
	space   = []byte{' '}
)

type connection struct {
	id string

	app *gsockets.App
	ws  *websocket.Conn

	logger log.Logger

	sendCh chan []byte

	closeCh chan struct{}
}

func NewConnection(app *gsockets.App, conn *websocket.Conn, logger log.Logger) gsockets.Connection {
	connId := generateConnectionId()
	newConn := &connection{
		id:      connId,
		app:     app,
		ws:      conn,
		logger:  logger.With("connection", connId),
		closeCh: make(chan struct{}),
		sendCh:  make(chan []byte),
	}

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
	panic("not implemented") // TODO: Implement
}

// Close closes the current connection
func (c *connection) Close() {
	panic("not implemented") // TODO: Implement
}

func generateConnectionId() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d.%d", rand.Intn(1000000000), rand.Intn(99999999999999))
}
