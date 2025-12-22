// Package model contains data models and domain entities for the Yan application.
package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type BaseModel struct {
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

func (m *BaseModel) TouchUpdated() {
	m.UpdatedAt = time.Now()
}

type NullString struct {
	sql.NullString
}

func (v *NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

func (v *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.String = *s
	} else {
		v.Valid = false
	}
	return nil
}

type NullInt64 struct {
	sql.NullInt64
}

func (v *NullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	}
	return json.Marshal(nil)
}

func (v *NullInt64) UnmarshalJSON(data []byte) error {
	var i *int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		v.Valid = true
		v.Int64 = *i
	} else {
		v.Valid = false
	}
	return nil
}

type NullInt32 struct {
	sql.NullInt32
}

func (v *NullInt32) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int32)
	}
	return json.Marshal(nil)
}

func (v *NullInt32) UnmarshalJSON(data []byte) error {
	var i *int32
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		v.Valid = true
		v.Int32 = *i
	} else {
		v.Valid = false
	}
	return nil
}

type NullFloat64 struct {
	sql.NullFloat64
}

func (v *NullFloat64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Float64)
	}
	return json.Marshal(nil)
}

func (v *NullFloat64) UnmarshalJSON(data []byte) error {
	var f *float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	if f != nil {
		v.Valid = true
		v.Float64 = *f
	} else {
		v.Valid = false
	}
	return nil
}

type NullBool struct {
	sql.NullBool
}

func (v *NullBool) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Bool)
	}
	return json.Marshal(nil)
}

func (v *NullBool) UnmarshalJSON(data []byte) error {
	var b *bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if b != nil {
		v.Valid = true
		v.Bool = *b
	} else {
		v.Valid = false
	}
	return nil
}

type NullTime struct {
	sql.NullTime
}

func (v *NullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	}
	return json.Marshal(nil)
}

func (v *NullTime) UnmarshalJSON(data []byte) error {
	var t *time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	if t != nil {
		v.Valid = true
		v.Time = *t
	} else {
		v.Valid = false
	}
	return nil
}
