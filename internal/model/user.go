package model

import "time"

type User struct {
	BaseModel
	ID           int64     `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Email        string    `db:"email" json:"email"`
	Status       int       `db:"status" json:"status"`     // 1 normal, 0 disable
	IsAdmin      int       `db:"is_admin" json:"is_admin"` // 1 true, 0 false
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

const (
	// User status
	UserStatusDisabled = 0
	UserStatusNormal   = 1
)

const (
	// Admin flag
	UserAdminFalse = 0
	UserAdminTrue  = 1
)

func (User) TableName() string {
	return "users"
}

func (u User) IsActive() bool {
	return u.Status == 1
}

func (u User) IsAdministrator() bool {
	return u.IsAdmin == 1
}
