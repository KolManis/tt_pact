package service

import (
	"context"

	"github.com/KolManis/tt_pact/internal/model"
)

// TelegramService - это то, что API ожидает от сервисного слоя
type TelegramService interface {
	CreateSession(ctx context.Context) (*model.Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
	SendMessage(ctx context.Context, sessionID, peer, text string) (int64, error)
	SubscribeMessages(ctx context.Context, sessionID string) (<-chan *model.Message, error)
	// UnsubscribeMessages ... пока под вопросом
}
