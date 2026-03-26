package telegram

import (
	"context"
	"errors"

	"github.com/KolManis/tt_pact/internal/model"
)

func (s *service) SubscribeMessages(ctx context.Context, sessionID string) (<-chan *model.Message, error) {
	session, err := s.sessionRepo.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if session.Status != model.SessionStatusReady {
		return nil, errors.New("session not ready")
	}

	ch := make(chan *model.Message, 100)

	s.mu.Lock()
	s.subscribers[sessionID] = append(s.subscribers[sessionID], ch)
	s.mu.Unlock()

	go func() {
		<-ctx.Done()
		s.mu.Lock()
		defer s.mu.Unlock()

		subs := s.subscribers[sessionID]
		for i, sub := range subs {
			if sub == ch {
				s.subscribers[sessionID] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
		if len(s.subscribers[sessionID]) == 0 {
			delete(s.subscribers, sessionID)
		}
		close(ch)
	}()

	return ch, nil
}
