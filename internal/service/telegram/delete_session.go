package telegram

import (
	"context"
	"fmt"
	"os"
)

func (s *service) DeleteSession(ctx context.Context, sessionID string) error {
	s.mu.Lock()
	cancel, ok := s.cancels[sessionID]
	delete(s.cancels, sessionID)
	delete(s.clients, sessionID)
	delete(s.dispatchers, sessionID)

	if subs, ok := s.subscribers[sessionID]; ok {
		for _, ch := range subs {
			close(ch)
		}
		delete(s.subscribers, sessionID)
	}
	s.mu.Unlock()

	if ok && cancel != nil {
		cancel() // просто отменяем контекст
	}

	// Удаляем файл
	os.Remove(fmt.Sprintf("./sessions/%s.json", sessionID))

	return s.sessionRepo.Delete(ctx, sessionID)
}
