package v1

import (
	"github.com/KolManis/tt_pact/internal/service"
	telegramV1 "github.com/KolManis/tt_pact/pkg/proto/pact/telegram/v1"
)

type api struct {
	telegramV1.UnimplementedTelegramServiceServer

	telegramService service.TelegramService
}

func NewAPI(telegramService service.TelegramService) *api {
	return &api{
		telegramService: telegramService,
	}
}
