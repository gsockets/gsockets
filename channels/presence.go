package channels

import (
	"encoding/json"

	"github.com/gsockets/gsockets"
)

type presenceChannel struct {
	*privateChannel
}

func newPresenceChannel(cm gsockets.ChannelManager) gsockets.Channel {
	return &presenceChannel{&privateChannel{&publicChannel{channelManager: cm}}}
}

func (pc *presenceChannel) Subscribe(appId string, conn gsockets.Connection, payload gsockets.MessageData) error {
	err := pc.verifySignature(conn, payload)
	if err != nil {
		return err
	}

	var presenceMember gsockets.PresenceMember

	err = json.Unmarshal([]byte(payload.ChannelData), &presenceMember)
	if err != nil {
		return err
	}

	if presenceMember.UserId == "" {
		return gsockets.PusherError{Code: gsockets.ERROR_CONNECTION_IS_UNAUTHORIZED, Message: "user_id must be present in presence channel"}
	}

	pc.channelManager.SubscribeToChannel(appId, payload.Channel, conn, payload)
	members := pc.channelManager.GetChannelMembers(appId, payload.Channel)

	// If this user has not previously joined this presence channel, we'll trigger the
	// member_added event and set the presence channel details for the connection.
	if _, ok := members[presenceMember.UserId]; !ok {
		resp := gsockets.PusherSentMessage{
			Event:   "pusher_internal:member_added",
			Channel: payload.Channel,
			Data:    presenceMember,
		}

		pc.channelManager.BroadcastExcept(appId, payload.Channel, resp, conn.Id())
		conn.SetPresence(payload.Channel, presenceMember)

		members[presenceMember.UserId] = presenceMember
	}

	userIds := make([]string, 0)
	userHash := make(map[string]map[string]any)

	for userId, presenceInfo := range members {
		userIds = append(userIds, userId)
		userHash[userId] = presenceInfo.UserInfo
	}

	// Presence channel subscription reply includes current state of the members currently
	// subscribed to the channel.
	resp := gsockets.PusherSentMessage{
		Event:   "pusher_internal:subscription_succeeded",
		Channel: payload.Channel,
		Data: gsockets.PusherPresencePayload{
			Presence: gsockets.PusherPresenceData{
				Ids:   userIds,
				Hash:  userHash,
				Count: len(members),
			},
		},
	}

	conn.Send(resp)
	return nil
}

func (pc *presenceChannel) Unsubscribe(appId, channel string, conn gsockets.Connection) error {
	pc.publicChannel.Unsubscribe(appId, channel, conn)

	member, _ := conn.GetPresence(channel)
	conn.RemovePresence(channel)

	members := pc.channelManager.GetChannelMembers(appId, channel)
	if _, ok := members[member.UserId]; !ok {
		resp := gsockets.PusherSentMessage{
			Event:   "pusher_internal:member_removed",
			Channel: channel,
			Data:    map[string]string{"user_id": member.UserId},
		}

		pc.channelManager.BroadcastExcept(appId, channel, resp, conn.Id())
	}

	return nil
}
