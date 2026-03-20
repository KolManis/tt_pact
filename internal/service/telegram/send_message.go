package telegram

import (
	"context"
	"errors"
	"time"

	"github.com/KolManis/tt_pact/internal/model"
)

func (s *service) SendMessage(ctx context.Context, sessionID, peer, text string) (int64, error) {
	// Проверяем, что сессия существует
	session, err := s.sessionRepo.Get(ctx, sessionID)
	if err != nil {
		return 0, err
	}

	// Проверяем, что сессия готова
	if session.Status != model.SessionStatusReady {
		return 0, errors.New("session not ready")
	}

	// TODO: реальная отправка через gotd
	// msg, err := s.tgClient.SendMessage(ctx, sessionID, peer, text)

	// Пока возвращаем заглушку
	messageID := time.Now().Unix()

	return messageID, nil
}
