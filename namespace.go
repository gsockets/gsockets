package gsockets

import (
	"errors"
	"sync"
)

// ErrConnectionNotFound is returned when we try to get a connection that does not exist in the instance.
var ErrConnectionNotFound = errors.New("namespace: connection not found")

type Namespace struct {
	// channels store the channels with all the connections that are currently subscribed to it.
	// It only stores the connection's id, the connections themselves are stored in the conns map.
	channels map[string]map[string]bool

	// conns stores all the connections to this instance. It's a map with connection's
	// id as key and the connection itself as the value.
	conns map[string]Connection

	channelLock sync.Mutex
	connLock    sync.Mutex
}

func NewNamespace() *Namespace {
	return &Namespace{
		channels: make(map[string]map[string]bool),
		conns:    make(map[string]Connection),
	}
}

// GetChannels returns all the channel names currently maintained in this instance.
func (n *Namespace) GetChannels() []string {
	n.channelLock.Lock()
	defer n.channelLock.Unlock()

	channels := make([]string, len(n.channels))

	i := 0
	for channel := range n.channels {
		channels[i] = channel
		i++
	}

	return channels
}

// GetConnections returns all the connections maintained in this instance.
func (n *Namespace) GetConnections() map[string]Connection {
	n.connLock.Lock()
	defer n.connLock.Unlock()

	return n.conns
}

// AddConnection adds a connections to this instance. Generally should be called when
// the websocket connection is first established.
func (n *Namespace) AddConnection(conn Connection) {
	n.connLock.Lock()
	defer n.connLock.Unlock()

	if _, ok := n.conns[conn.Id()]; ok {
		return
	}

	n.conns[conn.Id()] = conn
}

// RemoveConnection will remove a connection from this instance. Removing a connection will
// cause it to be removed from all the channels. Should be called when the websocket connection
// is closed.
func (n *Namespace) RemoveConnection(connId string) {
	n.RemoveConnectionFromChannel(connId, n.GetChannels()...)

	n.connLock.Lock()
	defer n.connLock.Unlock()

	delete(n.conns, connId)
}

// AddConnectionToChannel will add a connection to a channel. If the channel is not present in this
// instance, the channel will be created and then the connection will be added to it. This only adds
// the connection id to the channel connection map, the actual connection should already be present
// on the conns map.
func (n *Namespace) AddConnectionToChannel(channelName string, conn Connection) {
	n.channelLock.Lock()
	defer n.channelLock.Unlock()

	if _, ok := n.channels[channelName]; !ok {
		n.channels[channelName] = make(map[string]bool)
	}

	channelConnections := n.channels[channelName]
	if _, ok := channelConnections[conn.Id()]; ok {
		return
	}

	channelConnections[conn.Id()] = true
	n.channels[channelName] = channelConnections
}

// RemoveConnectionFromChannel will remove a connection from a channel. If after removal the channel
// does not have any more connection, it will remove the channel from the instance too.
func (n *Namespace) RemoveConnectionFromChannel(connId string, channels ...string) {
	n.channelLock.Lock()
	defer n.channelLock.Unlock()

	remove := func(channelName string) {
		channelConnections, ok := n.channels[channelName]
		if !ok {
			return
		}

		delete(channelConnections, connId)

		if len(channelConnections) == 0 {
			delete(n.channels, channelName)
		} else {
			n.channels[channelName] = channelConnections
		}
	}

	for _, channelName := range channels {
		remove(channelName)
	}
}

// IsInChannel returns boolean indicating whether a given connection is subscribed to a channel.
func (n *Namespace) IsInChannel(connId string, channelName string) bool {
	n.channelLock.Lock()
	defer n.channelLock.Unlock()

	channelConnections, ok := n.channels[channelName]
	if !ok {
		return false
	}

	_, ok = channelConnections[connId]
	return ok
}

// GetConnection returns a connection instance. Will return an error if no connection found with
// the given id.
func (n *Namespace) GetConnection(connId string) (Connection, error) {
	n.connLock.Lock()
	defer n.connLock.Unlock()

	conn, ok := n.conns[connId]
	if !ok {
		return nil, ErrConnectionNotFound
	}

	return conn, nil
}

// GetChannelConnections returns all the connections attached to a channel.
func (n *Namespace) GetChannelConnections(channelName string) []Connection {
	n.channelLock.Lock()
	defer n.channelLock.Unlock()

	channelConns, ok := n.channels[channelName]
	if !ok {
		return []Connection{}
	}

	conns := make([]Connection, len(channelConns))

	i := 0
	for conn := range channelConns {
		conn, err := n.GetConnection(conn)
		if err != nil {
			continue
		}

		conns[i] = conn
		i++
	}

	return conns
}
