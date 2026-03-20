package v1

import (
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/KolManis/tt_pact/internal/converter"
	telegramV1 "github.com/KolManis/tt_pact/pkg/proto/pact/telegram/v1"
)

func (a *api) SubscribeMessages(req *telegramV1.SubscribeMessagesRequest, stream telegramV1.TelegramService_SubscribeMessagesServer) error {
	msgCh, err := a.telegramService.SubscribeMessages(stream.Context(), req.SessionId)
	if err != nil {
		return status.Errorf(codes.NotFound, "failed to subscribe: %v", err)
	}

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case msg, ok := <-msgCh:
			if !ok {
				return nil
			}
			if err := stream.Send(converter.MessageToProtoUpdate(msg)); err != nil {
				log.Printf("failed to send message: %v", err)
				return err
			}
		}
	}
}
