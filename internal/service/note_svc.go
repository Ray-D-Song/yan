package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ray-d-song/yan/internal/model"
	"github.com/ray-d-song/yan/internal/repo"
)

var (
	ErrNoteNotFound      = errors.New("note not found")
	ErrNoteUnauthorized  = errors.New("unauthorized to access this note")
	ErrInvalidParentNote = errors.New("invalid parent note")
)

type NoteService interface {
	GetByID(ctx context.Context, id int64, userID int64) (*model.Note, error)
	GetByUserID(ctx context.Context, userID int64, status int) ([]*model.Note, error)
	GetByParentID(ctx context.Context, parentID sql.NullInt64, userID int64, status int) ([]*model.Note, error)
	GetFavorites(ctx context.Context, userID int64) ([]*model.Note, error)
	Create(ctx context.Context, n *model.Note) error
	Update(ctx context.Context, n *model.Note, userID int64) error
	Trash(ctx context.Context, id int64, userID int64) error
	Restore(ctx context.Context, id int64, userID int64) error
	Delete(ctx context.Context, id int64, userID int64) error
	ToggleFavorite(ctx context.Context, id int64, userID int64) error
	UpdatePosition(ctx context.Context, id int64, position int, userID int64) error
}

type noteService struct {
	noteRepo repo.NoteRepo
}

func NewNoteService(noteRepo repo.NoteRepo) NoteService {
	return &noteService{
		noteRepo: noteRepo,
	}
}

func (s *noteService) GetByID(ctx context.Context, id int64, userID int64) (*model.Note, error) {
	note, err := s.noteRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}

	// Check if the user owns this note
	if note.UserID != userID {
		return nil, ErrNoteUnauthorized
	}

	return note, nil
}

func (s *noteService) GetByUserID(ctx context.Context, userID int64, status int) ([]*model.Note, error) {
	return s.noteRepo.GetByUserID(ctx, userID, status)
}

func (s *noteService) GetByParentID(ctx context.Context, parentID sql.NullInt64, userID int64, status int) ([]*model.Note, error) {
	// If parentID is valid, check if the parent note exists and belongs to the user
	if parentID.Valid {
		parentNote, err := s.noteRepo.GetByID(ctx, parentID.Int64)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrInvalidParentNote
			}
			return nil, err
		}

		if parentNote.UserID != userID {
			return nil, ErrNoteUnauthorized
		}
	}

	return s.noteRepo.GetByParentID(ctx, parentID, userID, status)
}

func (s *noteService) GetFavorites(ctx context.Context, userID int64) ([]*model.Note, error) {
	return s.noteRepo.GetFavorites(ctx, userID)
}

func (s *noteService) Create(ctx context.Context, n *model.Note) error {
	// If parent_id is provided, validate it
	if n.ParentID.Valid {
		parentNote, err := s.noteRepo.GetByID(ctx, n.ParentID.Int64)
		if err != nil {
			if err == sql.ErrNoRows {
				return ErrInvalidParentNote
			}
			return err
		}

		// Parent note must belong to the same user
		if parentNote.UserID != n.UserID {
			return ErrNoteUnauthorized
		}
	}

	return s.noteRepo.Create(ctx, n)
}

func (s *noteService) Update(ctx context.Context, n *model.Note, userID int64) error {
	// Check if note exists and belongs to user
	existingNote, err := s.GetByID(ctx, n.ID, userID)
	if err != nil {
		return err
	}

	// Keep the original user_id
	n.UserID = existingNote.UserID

	// If parent_id is being changed, validate it
	if n.ParentID.Valid {
		parentNote, err := s.noteRepo.GetByID(ctx, n.ParentID.Int64)
		if err != nil {
			if err == sql.ErrNoRows {
				return ErrInvalidParentNote
			}
			return err
		}

		// Parent note must belong to the same user
		if parentNote.UserID != userID {
			return ErrNoteUnauthorized
		}

		// Prevent circular reference (note cannot be its own parent)
		if n.ParentID.Int64 == n.ID {
			return ErrInvalidParentNote
		}
	}

	return s.noteRepo.Update(ctx, n)
}

func (s *noteService) Trash(ctx context.Context, id int64, userID int64) error {
	// Check if note exists and belongs to user
	_, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	return s.noteRepo.UpdateStatus(ctx, id, model.NoteStatusTrashed)
}

func (s *noteService) Restore(ctx context.Context, id int64, userID int64) error {
	// Check if note exists and belongs to user
	_, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	return s.noteRepo.UpdateStatus(ctx, id, model.NoteStatusNormal)
}

func (s *noteService) Delete(ctx context.Context, id int64, userID int64) error {
	// Check if note exists and belongs to user
	_, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	return s.noteRepo.Delete(ctx, id)
}

func (s *noteService) ToggleFavorite(ctx context.Context, id int64, userID int64) error {
	// Check if note exists and belongs to user
	note, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	// Toggle favorite status
	newFavoriteStatus := model.NoteFavoriteNo
	if note.IsFavorite == model.NoteFavoriteNo {
		newFavoriteStatus = model.NoteFavoriteYes
	}

	return s.noteRepo.UpdateFavorite(ctx, id, newFavoriteStatus)
}

func (s *noteService) UpdatePosition(ctx context.Context, id int64, position int, userID int64) error {
	// Check if note exists and belongs to user
	_, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	return s.noteRepo.UpdatePosition(ctx, id, position)
}
