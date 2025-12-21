package model

import "time"

type Session struct {
	SessionID string    `db:"session_id" json:"session_id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Data      string    `db:"data" json:"-"` // serialized session data
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt string    `db:"created_at" json:"created_at"`
	UpdatedAt string    `db:"updated_at" json:"updated_at"`
}

func (Session) TableName() string {
	return "sessions"
}

func (s Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
