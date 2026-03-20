package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	telegramV1 "github.com/KolManis/tt_pact/pkg/proto/pact/telegram/v1"
)

func (a *api) DeleteSession(ctx context.Context, req *telegramV1.DeleteSessionRequest) (*telegramV1.DeleteSessionResponse, error) {
	err := a.telegramService.DeleteSession(ctx, req.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to delete session: %v", err)
	}

	return &telegramV1.DeleteSessionResponse{}, nil
}
