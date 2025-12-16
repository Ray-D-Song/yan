// Package model contains data models and domain entities for the Yan application.
package model

import "time"

type BaseModel struct {
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *BaseModel) TouchUpdated() {
	m.UpdatedAt = time.Now()
}
