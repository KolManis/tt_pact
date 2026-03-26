package telegram

import (
	"context"
	"fmt"
	"log"
	"os"
)

func (s *service) DeleteSession(ctx context.Context, sessionID string) error {
	s.mu.Lock()

	cancel, hasCancel := s.cancels[sessionID]

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

	if hasCancel && cancel != nil {
		cancel()
		log.Printf("Session %s: client stopped", sessionID)
	}

	sessionFile := fmt.Sprintf("./sessions/%s.json", sessionID)
	if err := os.Remove(sessionFile); err != nil && !os.IsNotExist(err) {
		log.Printf("Session %s: failed to remove session file: %v", sessionID, err)
	} else {
		log.Printf("Session %s: session file removed", sessionID)
	}

	if err := s.sessionRepo.Delete(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to delete session from repo: %w", err)
	}

	log.Printf(" Session %s deleted successfully", sessionID)
	return nil
}
