package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gotd/td/tg"
)

func (s *service) SendMessage(ctx context.Context, sessionID, peer, text string) (int64, error) {
	s.mu.RLock()
	client, ok := s.clients[sessionID]
	s.mu.RUnlock()

	if !ok {
		return 0, errors.New("session not found")
	}

	// Для username
	if !strings.HasPrefix(peer, "@") {
		return 0, errors.New("peer must start with @")
	}

	username := strings.TrimPrefix(peer, "@")
	resolved, err := client.API().ContactsResolveUsername(ctx, &tg.ContactsResolveUsernameRequest{
		Username: username,
	})
	if err != nil {
		return 0, fmt.Errorf("resolve user: %w", err)
	}

	if len(resolved.Users) == 0 {
		return 0, errors.New("user not found")
	}

	user := resolved.Users[0].(*tg.User)

	result, err := client.API().MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
		Peer: &tg.InputPeerUser{
			UserID:     user.ID,
			AccessHash: user.AccessHash,
		},
		Message:  text,
		RandomID: time.Now().UnixNano(),
	})
	if err != nil {
		return 0, err
	}

	switch updates := result.(type) {
	case *tg.UpdateShortSentMessage:
		return int64(updates.ID), nil
	case *tg.Updates:
		for _, update := range updates.Updates {
			if newMsg, ok := update.(*tg.UpdateNewMessage); ok {
				if msg, ok := newMsg.Message.(*tg.Message); ok {
					return int64(msg.ID), nil
				}
			}
		}
	}

	return 0, errors.New("no message id")
}
