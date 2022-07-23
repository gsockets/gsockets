package channels

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/gsockets/gsockets"
)

type privateChannel struct {
	*publicChannel
}

func newPrivateChannel(cm gsockets.ChannelManager) gsockets.Channel {
	return &privateChannel{
		&publicChannel{
			channelManager: cm,
		},
	}
}

func (c *privateChannel) Subscribe(appId string, conn gsockets.Connection, payload gsockets.MessageData) error {
	err := c.verifySignature(conn, payload)
	if err != nil {
		return err
	}

	c.channelManager.SubscribeToChannel(appId, payload.Channel, conn, payload)
	return nil
}

func (c *privateChannel) verifySignature(conn gsockets.Connection, payload gsockets.MessageData) error {
	// The pusher auth signature is in the following format: "<pusher-key>:<signature>", we are interested in the
	// signature part. We'll verify this signature against the one we generated to verify it's not an unauthorized
	// request.
	sigSlice := strings.SplitAfter(payload.Auth, ":")
	sig, err := hex.DecodeString(strings.Join(sigSlice[1:], ""))

	if err != nil {
		return gsockets.PusherError{Code: gsockets.ERROR_CONNECTION_IS_UNAUTHORIZED, Message: "invalid signature string provided"}
	}

	hasher := hmac.New(sha256.New, []byte(conn.App().Secret))
	hasher.Write([]byte(c.getDataToSign(conn, payload)))

	if valid := hmac.Equal(sig, hasher.Sum(nil)); !valid {
		return gsockets.PusherError{Code: gsockets.ERROR_CONNECTION_IS_UNAUTHORIZED, Message: "signature does not match"}
	}

	return nil
}

func (c *privateChannel) getDataToSign(conn gsockets.Connection, payload gsockets.MessageData) string {
	// For private channels, the string to sign is in the following format: "<socket-id>:<channel-name>".
	var signatureString strings.Builder
	signatureString.WriteString(conn.Id())
	signatureString.WriteString(":")
	signatureString.WriteString(payload.Channel)

	if payload.ChannelData != "" {
		signatureString.WriteString(":")
		signatureString.WriteString(payload.ChannelData)
	}

	return signatureString.String()
}
