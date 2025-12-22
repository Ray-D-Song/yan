package model

import "time"

type Session struct {
	BaseModel
	SessionID string    `db:"session_id" json:"sessionId"`
	UserID    int64     `db:"user_id" json:"userId"`
	Data      string    `db:"data" json:"-"` // serialized session data
	ExpiresAt time.Time `db:"expires_at" json:"expiresAt"`
}

func (Session) TableName() string {
	return "sessions"
}

func (s Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
