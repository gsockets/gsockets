package channels

import "github.com/gsockets/gsockets"

type publicChannel struct {
	name        string
	connections map[string]gsockets.Connection
}

func newPublicChannel(name string) gsockets.Channel {
	return &publicChannel{name: name, connections: make(map[string]gsockets.Connection)}
}

// Name returns the name of the channel. Channel names are unique in the server.
func (c *publicChannel) Name() string {
	return c.name
}

// Connections returns all the clients subscribed to this channel.
func (c *publicChannel) Connections() []gsockets.Connection {
	conns := make([]gsockets.Connection, len(c.connections))
	for _, conn := range c.connections {
		conns = append(conns, conn)
	}

	return conns
}

// Subscribe adds a new connection to the channel.
func (c *publicChannel) Subscribe(conn gsockets.Connection, payload any) {
	if c.IsSubscribed(conn) {
		return
	}

	c.connections[conn.Id()] = conn

	resp := struct {
		Event   string `json:"event"`
		Channel string `json:"channel"`
		Data    any    `json:"data"`
	}{
		Event:   "pusher_internal:subscription_succeeded",
		Channel: c.name,
		Data:    "{}",
	}

	conn.Send(resp)
}

// Unsubscribe removes a connection from the channel.
func (c *publicChannel) Unsubscribe(conn gsockets.Connection) {
	if !c.IsSubscribed(conn) {
		return
	}

	delete(c.connections, conn.Id())
}

// IsSubscribed determines if the given connection is already connected to the channel.
func (c *publicChannel) IsSubscribed(conn gsockets.Connection) bool {
	_, ok := c.connections[conn.Id()]
	return ok
}

// Broadcast will send the given data to all the subscribed connections to this channel.
func (c *publicChannel) Broadcast(data any) {
	for _, conn := range c.connections {
		conn.Send(data)
	}
}

// BroadcastExcept sends the data to all the channels except to the connection identified
// by the given connection id.
func (c *publicChannel) BroadcastExcept(data any, connToExclude string) {
	for _, conn := range c.connections {
		if conn.Id() == connToExclude {
			continue
		}

		conn.Send(data)
	}
}
