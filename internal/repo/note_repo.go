// Package repo provides data access layer implementations for database operations.
package repo

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/ray-d-song/yan/internal/model"
)

type NoteRepo interface {
	GetByID(ctx context.Context, id int64) (*model.Note, error)
	GetByUserID(ctx context.Context, userID int64, status int) ([]*model.Note, error)
	GetByParentID(ctx context.Context, parentID sql.NullInt64, userID int64, status int) ([]*model.Note, error)
	GetFavorites(ctx context.Context, userID int64) ([]*model.Note, error)
	Create(ctx context.Context, n *model.Note) error
	Update(ctx context.Context, n *model.Note) error
	Delete(ctx context.Context, id int64) error
	UpdateStatus(ctx context.Context, id int64, status int) error
	UpdateFavorite(ctx context.Context, id int64, isFavorite int) error
	UpdatePosition(ctx context.Context, id int64, position int) error
}

type noteRepo struct {
	db *sqlx.DB
}

func NewNoteRepo(db *sqlx.DB) NoteRepo {
	return &noteRepo{db: db}
}

func (r *noteRepo) GetByID(ctx context.Context, id int64) (*model.Note, error) {
	var n model.Note
	err := r.db.GetContext(ctx, &n, `
		SELECT
			id, parent_id, user_id, title, content,
			icon, is_favorite, position, status, created_at, updated_at
		FROM notes
		WHERE id = ?
		LIMIT 1
	`, id)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *noteRepo) GetByUserID(ctx context.Context, userID int64, status int) ([]*model.Note, error) {
	notes := make([]*model.Note, 0)
	err := r.db.SelectContext(ctx, &notes, `
		SELECT
			id, parent_id, user_id, title, content,
			icon, is_favorite, position, status, created_at, updated_at
		FROM notes
		WHERE user_id = ? AND status = ?
		ORDER BY position ASC, created_at DESC
	`, userID, status)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *noteRepo) GetByParentID(ctx context.Context, parentID sql.NullInt64, userID int64, status int) ([]*model.Note, error) {
	notes := make([]*model.Note, 0)
	var err error

	if parentID.Valid {
		err = r.db.SelectContext(ctx, &notes, `
			SELECT
				id, parent_id, user_id, title, content,
				icon, is_favorite, position, status, created_at, updated_at
			FROM notes
			WHERE parent_id = ? AND user_id = ? AND status = ?
			ORDER BY position ASC, created_at DESC
		`, parentID.Int64, userID, status)
	} else {
		err = r.db.SelectContext(ctx, &notes, `
			SELECT
				id, parent_id, user_id, title, content,
				icon, is_favorite, position, status, created_at, updated_at
			FROM notes
			WHERE parent_id IS NULL AND user_id = ? AND status = ?
			ORDER BY position ASC, created_at DESC
		`, userID, status)
	}

	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *noteRepo) GetFavorites(ctx context.Context, userID int64) ([]*model.Note, error) {
	notes := make([]*model.Note, 0)
	err := r.db.SelectContext(ctx, &notes, `
		SELECT
			id, parent_id, user_id, title, content,
			icon, is_favorite, position, status, created_at, updated_at
		FROM notes
		WHERE user_id = ? AND is_favorite = 1 AND status = 1
		ORDER BY position ASC, created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *noteRepo) Create(ctx context.Context, n *model.Note) error {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO notes (
			parent_id,
			user_id,
			title,
			content,
			icon,
			is_favorite,
			position,
			status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		n.ParentID,
		n.UserID,
		n.Title,
		n.Content,
		n.Icon,
		n.IsFavorite,
		n.Position,
		n.Status,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	n.ID = id
	return nil
}

func (r *noteRepo) Update(ctx context.Context, n *model.Note) error {
	n.TouchUpdated()
	_, err := r.db.ExecContext(ctx, `
		UPDATE notes
		SET
			parent_id = ?,
			title = ?,
			content = ?,
			icon = ?,
			is_favorite = ?,
			position = ?,
			status = ?,
			updated_at = datetime('now')
		WHERE id = ?
	`,
		n.ParentID,
		n.Title,
		n.Content,
		n.Icon,
		n.IsFavorite,
		n.Position,
		n.Status,
		n.ID,
	)
	return err
}

func (r *noteRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM notes WHERE id = ?
	`, id)

	return err
}

func (r *noteRepo) UpdateStatus(ctx context.Context, id int64, status int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notes
		SET status = ?, updated_at = datetime('now')
		WHERE id = ?
	`, status, id)

	return err
}

func (r *noteRepo) UpdateFavorite(ctx context.Context, id int64, isFavorite int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notes
		SET is_favorite = ?, updated_at = datetime('now')
		WHERE id = ?
	`, isFavorite, id)

	return err
}

func (r *noteRepo) UpdatePosition(ctx context.Context, id int64, position int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notes
		SET position = ?, updated_at = datetime('now')
		WHERE id = ?
	`, position, id)

	return err
}
