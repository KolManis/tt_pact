package repository

import (
	"context"
	"errors"

	"github.com/KolManis/tt_pact/internal/model"
)

var ErrNotFound = errors.New("not found")

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) error
	Get(ctx context.Context, id string) (*model.Session, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, session *model.Session) error
}
