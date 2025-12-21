// Package repo provides data access layer implementations for database operations.
package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ray-d-song/yan/internal/model"
)

type UserRepo interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, u *model.User) error
	Update(ctx context.Context, u *model.User) error
	DisableByID(ctx context.Context, id int64) error
	UpdatePassword(ctx context.Context, id int64, newHash string) error
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var u model.User
	err := r.db.GetContext(ctx, &u, `
		SELECT
			id, username, password_hash, email,
			status, is_admin, created_at, updated_at
		FROM users
		WHERE id = ?
		LIMIT 1
	`, id)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.db.GetContext(ctx, &u, `
		SELECT
			id, username, password_hash, email,
			status, is_admin, created_at, updated_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`, email)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepo) Create(ctx context.Context, u *model.User) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var hasAdmin bool
	err = tx.Get(&hasAdmin, `SELECT EXISTS (
		SELECT 1 FROM users WHERE is_admin = 1
		FOR UPDATE
	)`)

	var isAdmin int
	if !hasAdmin {
		isAdmin = 1
	} else {
		isAdmin = u.IsAdmin
	}

	res, err := tx.ExecContext(ctx, `
		INSERT INTO users (
			username,
			password_hash,
			email,
			status,
			is_admin
		) VALUES (?, ?, ?, ?, ?)
	`,
		u.Username,
		u.PasswordHash,
		u.Email,
		u.Status,
		isAdmin,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = id
	return nil
}

func (r *userRepo) Update(ctx context.Context, u *model.User) error {
	u.TouchUpdated()
	_, err := r.db.ExecContext(ctx, `
        UPDATE users
        SET
            username = ?,
            email = ?,
            status = ?,
            is_admin = ?,
            updated_at = datetime('now')
        WHERE id = ?
    `,
		u.Username,
		u.Email,
		u.Status,
		u.IsAdmin,
		u.ID,
	)
	return err
}

func (r *userRepo) DisableByID(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET status = 0 WHERE id = ?
	`, id)

	return err
}

func (r *userRepo) UpdatePassword(ctx context.Context, id int64, newHash string) error {
	_, err := r.db.ExecContext(ctx, `
	Update users SET password_hash = ? WHERE id = ?
	`, newHash, id)
	return err
}
