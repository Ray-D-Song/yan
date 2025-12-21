package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ray-d-song/yan/internal/model"
)

type SessionRepo interface {
	GetByID(ctx context.Context, sessionID string) (*model.Session, error)
	GetByUserID(ctx context.Context, userID int64) ([]*model.Session, error)
	Create(ctx context.Context, s *model.Session) error
	Update(ctx context.Context, s *model.Session) error
	Delete(ctx context.Context, sessionID string) error
	DeleteByUserID(ctx context.Context, userID int64) error
	DeleteExpired(ctx context.Context) error
}

type sessionRepo struct {
	db *sqlx.DB
}

func NewSessionRepo(db *sqlx.DB) SessionRepo {
	return &sessionRepo{db: db}
}

func (r *sessionRepo) GetByID(ctx context.Context, sessionID string) (*model.Session, error) {
	var s model.Session
	err := r.db.GetContext(ctx, &s, `
		SELECT session_id, user_id, data, expires_at, created_at, updated_at
		FROM sessions
		WHERE session_id = ?
		LIMIT 1
	`, sessionID)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *sessionRepo) GetByUserID(ctx context.Context, userID int64) ([]*model.Session, error) {
	var sessions []*model.Session
	err := r.db.SelectContext(ctx, &sessions, `
		SELECT session_id, user_id, data, expires_at, created_at, updated_at
		FROM sessions
		WHERE user_id = ? AND expires_at > datetime('now')
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *sessionRepo) Create(ctx context.Context, s *model.Session) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO sessions (session_id, user_id, data, expires_at)
		VALUES (?, ?, ?, ?)
	`, s.SessionID, s.UserID, s.Data, s.ExpiresAt.Format(time.RFC3339))
	return err
}

func (r *sessionRepo) Update(ctx context.Context, s *model.Session) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE sessions
		SET data = ?, expires_at = ?, updated_at = datetime('now')
		WHERE session_id = ?
	`, s.Data, s.ExpiresAt.Format(time.RFC3339), s.SessionID)
	return err
}

func (r *sessionRepo) Delete(ctx context.Context, sessionID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM sessions WHERE session_id = ?
	`, sessionID)
	return err
}

func (r *sessionRepo) DeleteByUserID(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM sessions WHERE user_id = ?
	`, userID)
	return err
}

func (r *sessionRepo) DeleteExpired(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM sessions WHERE expires_at < datetime('now')
	`)
	return err
}
