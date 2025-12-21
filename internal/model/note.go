package model

import "database/sql"

type Note struct {
	BaseModel
	ID         int64          `db:"id" json:"id"`
	ParentID   sql.NullInt64  `db:"parent_id" json:"parent_id"`
	UserID     int64          `db:"user_id" json:"user_id"`
	Title      string         `db:"title" json:"title"`
	Content    string         `db:"content" json:"content"`
	Icon       sql.NullString `db:"icon" json:"icon"`
	IsFavorite int            `db:"is_favorite" json:"is_favorite"` // 1 favorite, 0 not favorite
	Position   int            `db:"position" json:"position"`
	Status     int            `db:"status" json:"status"` // 1 normal, 0 trashed
	CreatedAt  string         `db:"created_at" json:"created_at"`
}

const (
	// Note status
	NoteStatusTrashed = 0
	NoteStatusNormal  = 1
)

const (
	// Note favorite
	NoteFavoriteNo  = 0
	NoteFavoriteYes = 1
)

func (Note) TableName() string {
	return "notes"
}

func (n Note) IsNormal() bool {
	return n.Status == NoteStatusNormal
}

func (n Note) IsTrashed() bool {
	return n.Status == NoteStatusTrashed
}

func (n Note) IsFavorited() bool {
	return n.IsFavorite == NoteFavoriteYes
}

func (n Note) IsRoot() bool {
	return !n.ParentID.Valid
}
