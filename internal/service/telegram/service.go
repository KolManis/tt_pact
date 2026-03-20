package telegram

import (
	"github.com/KolManis/tt_pact/internal/repository"
	def "github.com/KolManis/tt_pact/internal/service"
)

// Проверяем, что service реализует интерфейс Service
var _ def.TelegramService = (*service)(nil)

type service struct {
	sessionRepo repository.SessionRepository // зависимость	 от репозитория
	// TODO другие зависимости  tgClient
}

func NewService(sessionRepo repository.SessionRepository) *service {
	return &service{
		sessionRepo: sessionRepo,
	}
}
