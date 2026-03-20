package telegram

import (
	"context"
)

func (s *service) DeleteSession(ctx context.Context, sessionID string) error {
	// TODO: закрыть соединение с Telegram
	// TODO: вызвать auth.logOut если нужно

	return s.sessionRepo.Delete(ctx, sessionID)
}
