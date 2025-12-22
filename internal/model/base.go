// Package model contains data models and domain entities for the Yan application.
package model

import "time"

type BaseModel struct {
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

func (m *BaseModel) TouchUpdated() {
	m.UpdatedAt = time.Now()
}
