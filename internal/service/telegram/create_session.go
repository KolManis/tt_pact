package telegram

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/KolManis/tt_pact/internal/model"
	"github.com/google/uuid"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth/qrlogin"
	"github.com/gotd/td/tg"
)

func (s *service) CreateSession(ctx context.Context) (*model.Session, error) {
	sessionID := uuid.New().String()

	// Проверяем, есть ли уже такая сессия в памяти
	s.mu.RLock()
	existingClient, exists := s.clients[sessionID]
	s.mu.RUnlock()

	if exists && existingClient != nil {
		// Сессия уже активна
		log.Printf("Session %s already active", sessionID)
		session, _ := s.sessionRepo.Get(ctx, sessionID)
		return session, nil
	}

	disp := tg.NewUpdateDispatcher()

	client := telegram.NewClient(s.appID, s.appHash, telegram.Options{
		UpdateHandler: &disp,
		SessionStorage: &telegram.FileSessionStorage{
			Path: fmt.Sprintf("./sessions/%s.json", sessionID),
		},
	})

	qrChan := make(chan string, 1)

	session := &model.Session{
		ID:      sessionID,
		QRCode:  "",
		Status:  model.SessionStatusPending,
		Created: time.Now().Unix(),
	}

	s.mu.Lock()
	s.clients[sessionID] = client
	s.dispatchers[sessionID] = &disp
	s.mu.Unlock()

	// Запускаем клиента в фоне
	go func() {
		clientCtx, cancel := context.WithCancel(context.Background())
		s.mu.Lock()
		s.cancels[sessionID] = cancel
		s.mu.Unlock()

		defer func() {
			s.mu.Lock()
			delete(s.cancels, sessionID)
			s.mu.Unlock()
		}()

		err := client.Run(clientCtx, func(ctx context.Context) error {
			status, err := client.Auth().Status(ctx)
			if err != nil {
				return err
			}

			if status.Authorized {
				log.Printf("Session %s already authorized", sessionID)
				s.updateSessionStatus(sessionID, model.SessionStatusReady)
				// Держим клиента живым
				<-ctx.Done()
				return nil
			}

			loginToken := qrlogin.OnLoginToken(&disp)

			_, err = client.QR().Auth(ctx, loginToken, func(ctx context.Context, token qrlogin.Token) error {
				qrChan <- token.URL()
				return nil
			})
			if err != nil {
				return err
			}

			log.Printf("Session %s ready", sessionID)
			s.updateSessionStatus(sessionID, model.SessionStatusReady)

			// Держим клиента живым для обработки сообщений
			<-ctx.Done()
			return nil
		})

		if err != nil && err != context.Canceled {
			log.Printf("Session %s error: %v", sessionID, err)
			s.updateSessionStatus(sessionID, model.SessionStatusClosed)
		}
	}()

	// Ждем QR код
	var qrCode string
	select {
	case qrCode = <-qrChan:
		log.Printf("✅ QR Code for session %s", sessionID)
	case <-time.After(10 * time.Second):
		// Если QR не пришел, возможно сессия уже авторизована
		time.Sleep(2 * time.Second)
		s.mu.RLock()
		client := s.clients[sessionID]
		s.mu.RUnlock()
		if client != nil {
			log.Printf("Session %s using existing auth", sessionID)
			s.updateSessionStatus(sessionID, model.SessionStatusReady)
		}
		qrCode = ""
	}

	session.QRCode = qrCode

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession - получить существующую сессию по ID
func (s *service) GetSession(sessionID string) (*model.Session, error) {
	s.mu.RLock()
	client, ok := s.clients[sessionID]
	s.mu.RUnlock()

	if !ok || client == nil {
		return nil, fmt.Errorf("session %s not active", sessionID)
	}

	return s.sessionRepo.Get(context.Background(), sessionID)
}

func (s *service) updateSessionStatus(sessionID string, status model.SessionStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, err := s.sessionRepo.Get(context.Background(), sessionID)
	if err != nil {
		return
	}
	session.Status = status
	s.sessionRepo.Update(context.Background(), session)
}

func (s *service) broadcastMessage(sessionID string, msg *model.Message) {
	s.mu.RLock()
	subs := s.subscribers[sessionID]
	s.mu.RUnlock()

	for _, ch := range subs {
		select {
		case ch <- msg:
		default:
		}
	}
}
