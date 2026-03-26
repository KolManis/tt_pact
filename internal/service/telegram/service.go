package telegram

import (
	"context"
	"sync"

	"github.com/KolManis/tt_pact/internal/model"
	"github.com/KolManis/tt_pact/internal/repository"
	def "github.com/KolManis/tt_pact/internal/service"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

// Проверяем, что service реализует интерфейс Service
var _ def.TelegramService = (*service)(nil)

type service struct {
	sessionRepo repository.SessionRepository
	clients     map[string]*telegram.Client
	dispatchers map[string]*tg.UpdateDispatcher
	subscribers map[string][]chan *model.Message
	cancels     map[string]context.CancelFunc
	mu          sync.RWMutex
	appID       int
	appHash     string
}

func NewService(sessionRepo repository.SessionRepository, appID int, appHash string) *service {
	return &service{
		sessionRepo: sessionRepo,
		clients:     make(map[string]*telegram.Client),
		dispatchers: make(map[string]*tg.UpdateDispatcher),
		subscribers: make(map[string][]chan *model.Message),
		cancels:     make(map[string]context.CancelFunc),
		appID:       appID,
		appHash:     appHash,
	}
}

// Shutdown останавливает все сессии
func (s *service) Shutdown(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Отменяем все контексты
	for id, cancel := range s.cancels {
		cancel()
		delete(s.cancels, id)
	}

	// Закрываем все каналы подписчиков
	for _, subs := range s.subscribers {
		for _, ch := range subs {
			close(ch)
		}
	}
	s.subscribers = make(map[string][]chan *model.Message)

	// Очищаем клиенты
	s.clients = make(map[string]*telegram.Client)
	s.dispatchers = make(map[string]*tg.UpdateDispatcher)
}
