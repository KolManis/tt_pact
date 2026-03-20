package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/KolManis/tt_pact/internal/converter"
	telegramV1 "github.com/KolManis/tt_pact/pkg/proto/pact/telegram/v1"
)

func (a *api) SendMessage(ctx context.Context, req *telegramV1.SendMessageRequest) (*telegramV1.SendMessageResponse, error) {
	sessionID, peer, text := converter.SendMessageRequestToModel(req)

	messageID, err := a.telegramService.SendMessage(ctx, sessionID, peer, text)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send message: %v", err)
	}

	return &telegramV1.SendMessageResponse{
		MessageId: messageID,
	}, nil
}
