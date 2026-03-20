package telegram

import (
	"context"
	"time"

	"github.com/KolManis/tt_pact/internal/model"
)

func (s *service) SubscribeMessages(ctx context.Context, sessionID string) (<-chan *model.Message, error) {
	// Проверяем, что сессия существует
	_, err := s.sessionRepo.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	ch := make(chan *model.Message)

	// TODO: реальная подписка на сообщения из Telegram через gotd
	// Пока тестовая заглушка
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				ch <- &model.Message{
					ID:        t.Unix(),
					From:      "@test_user",
					Text:      "Test message at " + t.String(),
					Timestamp: t.Unix(),
					SessionID: sessionID,
				}
			}
		}
	}()

	return ch, nil
}
