package channelmanagers

import (
	"sync"

	"github.com/gsockets/gsockets"
)

type localChannelManager struct {
	// namespaces is a map that stores channels and connections for each app in the
	// server. Here the app id is the key, so all the connections and channels are
	// separate for individual apps.
	namespaces map[string]*gsockets.Namespace

	namespaceLock sync.Mutex
}

func newLocalChannelManager() gsockets.ChannelManager {
	return &localChannelManager{namespaces: make(map[string]*gsockets.Namespace)}
}

// getNamespace returns the namespace associated with the given appId, if no namespace exists
// for the given app, a new namespace is created.
func (l *localChannelManager) getNamespace(appId string) *gsockets.Namespace {
	l.namespaceLock.Lock()
	defer l.namespaceLock.Unlock()

	namespace, ok := l.namespaces[appId]
	if !ok {
		namespace = gsockets.NewNamespace()
		l.namespaces[appId] = namespace
	}

	return namespace
}

func (l *localChannelManager) AddConnection(appId string, conn gsockets.Connection) {
	l.getNamespace(appId).AddConnection(conn)
}

func (l *localChannelManager) RemoveConnection(appId string, conn gsockets.Connection) {
	l.getNamespace(appId).RemoveConnection(conn.Id())
}

func (l *localChannelManager) GetLocalConnections(appId string) []gsockets.Connection {
	conns := make([]gsockets.Connection, 0)
	for _, conn := range l.getNamespace(appId).GetConnections() {
		conns = append(conns, conn)
	}

	return conns
}

func (l *localChannelManager) GetLocalChannels(appId string) []string {
	return l.getNamespace(appId).GetChannels()
}

func (l *localChannelManager) GetGlobalChannels(appId string) []string {
	return l.GetLocalChannels(appId)
}

func (l *localChannelManager) GetGlobalChannelsWithConnectionCount(appId string) map[string]int {
	channels := l.GetGlobalChannels(appId)

	ret := make(map[string]int)
	for _, channel := range channels {
		conns := l.getNamespace(appId).GetChannelConnections(channel)
		ret[channel] = len(conns)
	}

	return ret
}

func (l *localChannelManager) GetChannelConnectionCount(appId string, channelName string) int {
	conns := l.getNamespace(appId).GetChannelConnections(channelName)
	return len(conns)
}

func (l *localChannelManager) GetChannelMembers(appId, channelName string) map[string]gsockets.PresenceMember {
	return l.getNamespace(appId).GetChannelMembers(channelName)
}

func (l *localChannelManager) SetUser(appId, userId, connId string) {
	l.getNamespace(appId).AddUser(userId, connId)
}

func (l *localChannelManager) RemoveUser(appId, userId, connId string) {
	l.getNamespace(appId).RemoveUser(userId, connId)
}

func (l *localChannelManager) GetUserConnections(appId, userId string) []gsockets.Connection {
	return l.getNamespace(appId).GetUserConnections(userId)
}

func (l *localChannelManager) TerminateUserConnections(appId, userId string) {
	conns := l.getNamespace(appId).GetChannelConnections(userId)
	for _, conn := range conns {
		pusherErr := gsockets.NewPusherError("pusher:error", "disconnected by the server", "", gsockets.ERROR_CONNECTION_IS_UNAUTHORIZED)

		go func(conn gsockets.Connection) {
			conn.Send(pusherErr)
			conn.Close()
		}(conn)
	}
}

func (l *localChannelManager) SubscribeToChannel(appId string, channelName string, conn gsockets.Connection, payload any) {
	l.getNamespace(appId).AddConnectionToChannel(channelName, conn)
}

func (l *localChannelManager) UnsubscribeFromChannel(appId string, channelName string, conn gsockets.Connection) {
	l.getNamespace(appId).RemoveConnectionFromChannel(conn.Id(), channelName)
}

func (l *localChannelManager) UnsubscribeFromAllChannels(appId string, conn string) {
	l.getNamespace(appId).RemoveConnectionFromChannel(conn, l.GetLocalChannels(appId)...)
}

func (l *localChannelManager) IsInChannel(appId string, channel string, conn gsockets.Connection) bool {
	return l.getNamespace(appId).IsInChannel(conn.Id(), channel)
}

func (l *localChannelManager) BroadcastToChannel(appId string, channel string, data any) {
	conns := l.getNamespace(appId).GetChannelConnections(channel)

	for _, conn := range conns {
		conn.Send(data)
	}
}

func (l *localChannelManager) BroadcastExcept(appId string, channel string, data any, connId string) {
	conns := l.getNamespace(appId).GetChannelConnections(channel)

	for _, conn := range conns {
		if conn.Id() == connId {
			continue
		}

		conn.Send(data)
	}
}
