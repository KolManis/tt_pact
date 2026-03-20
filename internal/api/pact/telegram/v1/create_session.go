package v1

import (
	"context"

	"github.com/KolManis/tt_pact/internal/converter"
	telegramV1 "github.com/KolManis/tt_pact/pkg/proto/pact/telegram/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) CreateSession(ctx context.Context, req *telegramV1.CreateSessionRequest) (*telegramV1.CreateSessionResponse, error) {
	session, err := a.telegramService.CreateSession(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}

	return converter.SessionToProtoResponse(session), nil
}
