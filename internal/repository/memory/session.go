// тоже разбить потом
package memory

import (
	"context"
	"sync"

	"github.com/KolManis/tt_pact/internal/model"
	"github.com/KolManis/tt_pact/internal/repository"
)

var _ repository.SessionRepository = (*sessionRepo)(nil)

type sessionRepo struct {
	mu       sync.RWMutex
	sessions map[string]*model.Session
}

func NewSessionRepo() *sessionRepo {
	return &sessionRepo{
		sessions: make(map[string]*model.Session),
	}
}

func (r *sessionRepo) Create(ctx context.Context, session *model.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.ID] = session
	return nil
}

func (r *sessionRepo) Get(ctx context.Context, id string) (*model.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, ok := r.sessions[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return session, nil
}

func (r *sessionRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, id)
	return nil
}

func (r *sessionRepo) Update(ctx context.Context, session *model.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.ID] = session
	return nil
}
