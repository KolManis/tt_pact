package telegram

import (
	"context"
	"time"

	"github.com/KolManis/tt_pact/internal/model"
	"github.com/google/uuid"
)

func (s *service) CreateSession(ctx context.Context) (*model.Session, error) {
	session := &model.Session{
		ID:      uuid.New().String(),
		QRCode:  "tg://qr/login?token=" + uuid.New().String(),
		Status:  model.SessionStatusReady,
		Created: time.Now().Unix(),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	// TODO: запустить процесс авторизации в Telegram
	// TODO: получить реальный QR код от gotd

	return session, nil
}
