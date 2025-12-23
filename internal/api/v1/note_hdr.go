// Package v1 provides HTTP handlers for API version 1 endpoints.
package v1

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ray-d-song/yan/internal/infra"
	"github.com/ray-d-song/yan/internal/model"
	"github.com/ray-d-song/yan/internal/service"
)

type NoteHandler struct {
	noteService service.NoteService
}

func NewNoteHandler(noteService service.NoteService) *NoteHandler {
	return &NoteHandler{
		noteService: noteService,
	}
}

// RegisterRoutes registers all note-related routes
// Note: Auth middleware should be applied before calling this
func (h *NoteHandler) RegisterRoutes(g *gin.RouterGroup) {
	g.POST("", h.CreateNote)
	g.GET("/:id", h.GetNote)
	g.GET("", h.ListNotes)
	g.PUT("/:id", h.UpdateNote)
	g.DELETE("/:id", h.DeleteNote)
	g.PUT("/:id/trash", h.TrashNote)
	g.PUT("/:id/restore", h.RestoreNote)
	g.PUT("/:id/favorite", h.ToggleFavorite)
	g.PUT("/:id/position", h.UpdatePosition)
}

// CreateNoteRequest represents the create note request payload
type CreateNoteRequest struct {
	ParentID   *int64  `json:"parent_id"`
	Title      string  `json:"title" binding:"required"`
	Content    string  `json:"content"`
	Icon       *string `json:"icon"`
	IsFavorite int     `json:"is_favorite"`
	Position   int     `json:"position"`
}

// UpdateNoteRequest represents the update note request payload
type UpdateNoteRequest struct {
	ParentID   *int64  `json:"parent_id"`
	Title      *string `json:"title"`
	Content    *string `json:"content"`
	Icon       *string `json:"icon"`
	IsFavorite *int    `json:"is_favorite"`
	Position   *int    `json:"position"`
	Status     *int    `json:"status"`
}

// UpdateNoteResponse represents the update note response payload
type UpdateNoteResponse struct {
	UpdatedAt string `json:"updatedAt"`
}

// UpdatePositionRequest represents the update position request payload
type UpdatePositionRequest struct {
	Position int `json:"position" binding:"required"`
}

// CreateNote creates a new note
// POST /api/v1/notes
func (h *NoteHandler) CreateNote(c *gin.Context) {
	userID, err := infra.UserIDFromCtx(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		return
	}

	var req CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	note := &model.Note{
		UserID:     userID,
		Title:      req.Title,
		Content:    req.Content,
		IsFavorite: req.IsFavorite,
		Position:   req.Position,
		Status:     model.NoteStatusNormal,
	}

	if req.ParentID != nil {
		note.ParentID = model.NullInt64{NullInt64: sql.NullInt64{Int64: *req.ParentID, Valid: true}}
	}

	if req.Icon != nil {
		note.Icon = model.NullString{NullString: sql.NullString{String: *req.Icon, Valid: true}}
	}

	if err := h.noteService.Create(c.Request.Context(), note); err != nil {
		if err == service.ErrInvalidParentNote {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		if err == service.ErrNoteUnauthorized {
			c.String(http.StatusForbidden, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, note)
}

// GetNote retrieves a note by ID
// GET /api/v1/notes/:id
func (h *NoteHandler) GetNote(c *gin.Context) {
	userID, err := infra.UserIDFromCtx(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid note id")
		return
	}

	note, err := h.noteService.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		if err == service.ErrNoteNotFound {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		if err == service.ErrNoteUnauthorized {
			c.String(http.StatusForbidden, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, note)
}

// ListNotes retrieves notes by parent_id or all user notes
// GET /api/v1/notes?parent_id=123&status=1
func (h *NoteHandler) ListNotes(c *gin.Context) {
	userID, err := infra.UserIDFromCtx(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	// Get query parameters
	parentIDStr := c.Query("parent_id")
	statusStr := c.DefaultQuery("status", "1")
	favoriteStr := c.Query("favorite")

	status, err := strconv.Atoi(statusStr)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid status")
		return
	}

	// Get favorites
	if favoriteStr == "true" || favoriteStr == "1" {
		notes, err := h.noteService.GetFavorites(c.Request.Context(), userID)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, notes)
		return
	}

	// Get notes by parent_id
	if parentIDStr != "" {
		var parentID sql.NullInt64
		if parentIDStr == "null" || parentIDStr == "0" {
			parentID = sql.NullInt64{Valid: false}
		} else {
			id, err := strconv.ParseInt(parentIDStr, 10, 64)
			if err != nil {
				c.String(http.StatusBadRequest, "invalid parent_id")
				return
			}
			parentID = sql.NullInt64{Int64: id, Valid: true}
		}

		notes, err := h.noteService.GetByParentID(c.Request.Context(), parentID, userID, status)
		if err != nil {
			if err == service.ErrInvalidParentNote {
				c.String(http.StatusBadRequest, err.Error())
				return
			}
			if err == service.ErrNoteUnauthorized {
				c.String(http.StatusForbidden, err.Error())
				return
			}
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, notes)
		return
	}

	// Get all notes for user
	notes, err := h.noteService.GetByUserID(c.Request.Context(), userID, status)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, notes)
}

// UpdateNote updates a note
// PUT /api/v1/notes/:id
func (h *NoteHandler) UpdateNote(c *gin.Context) {
	userID, err := infra.UserIDFromCtx(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid note id")
		return
	}

	var req UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// Get existing note to update only provided fields
	existingNote, err := h.noteService.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		if err == service.ErrNoteNotFound {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		if err == service.ErrNoteUnauthorized {
			c.String(http.StatusForbidden, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Update only provided fields
	note := &model.Note{
		ID:     id,
		UserID: existingNote.UserID,
		Status: existingNote.Status,
	}

	if req.Title != nil {
		note.Title = *req.Title
	} else {
		note.Title = existingNote.Title
	}

	if req.Content != nil {
		note.Content = *req.Content
	} else {
		note.Content = existingNote.Content
	}

	if req.IsFavorite != nil {
		note.IsFavorite = *req.IsFavorite
	} else {
		note.IsFavorite = existingNote.IsFavorite
	}

	if req.Position != nil {
		note.Position = *req.Position
	} else {
		note.Position = existingNote.Position
	}

	if req.Status != nil {
		note.Status = *req.Status
	}

	if req.ParentID != nil {
		note.ParentID = model.NullInt64{NullInt64: sql.NullInt64{Int64: *req.ParentID, Valid: true}}
	} else {
		note.ParentID = existingNote.ParentID
	}

	if req.Icon != nil {
		note.Icon = model.NullString{NullString: sql.NullString{String: *req.Icon, Valid: true}}
	} else {
		note.Icon = existingNote.Icon
	}

	if err := h.noteService.Update(c.Request.Context(), note, userID); err != nil {
		if err == service.ErrNoteNotFound {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		if err == service.ErrNoteUnauthorized {
			c.String(http.StatusForbidden, err.Error())
			return
		}
		if err == service.ErrInvalidParentNote {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, UpdateNoteResponse{
		UpdatedAt: note.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// DeleteNote permanently deletes a note
// DELETE /api/v1/notes/:id
func (h *NoteHandler) DeleteNote(c *gin.Context) {
	userID, err := infra.UserIDFromCtx(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid note id")
		return
	}

	if err := h.noteService.Delete(c.Request.Context(), id, userID); err != nil {
		if err == service.ErrNoteNotFound {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		if err == service.ErrNoteUnauthorized {
			c.String(http.StatusForbidden, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// TrashNote moves a note to trash (soft delete)
// PUT /api/v1/notes/:id/trash
func (h *NoteHandler) TrashNote(c *gin.Context) {
	userID, err := infra.UserIDFromCtx(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid note id")
		return
	}

	if err := h.noteService.Trash(c.Request.Context(), id, userID); err != nil {
		if err == service.ErrNoteNotFound {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		if err == service.ErrNoteUnauthorized {
			c.String(http.StatusForbidden, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// RestoreNote restores a note from trash
// PUT /api/v1/notes/:id/restore
func (h *NoteHandler) RestoreNote(c *gin.Context) {
	userID, err := infra.UserIDFromCtx(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid note id")
		return
	}

	if err := h.noteService.Restore(c.Request.Context(), id, userID); err != nil {
		if err == service.ErrNoteNotFound {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		if err == service.ErrNoteUnauthorized {
			c.String(http.StatusForbidden, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// ToggleFavorite toggles the favorite status of a note
// PUT /api/v1/notes/:id/favorite
func (h *NoteHandler) ToggleFavorite(c *gin.Context) {
	userID, err := infra.UserIDFromCtx(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid note id")
		return
	}

	if err := h.noteService.ToggleFavorite(c.Request.Context(), id, userID); err != nil {
		if err == service.ErrNoteNotFound {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		if err == service.ErrNoteUnauthorized {
			c.String(http.StatusForbidden, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// UpdatePosition updates the position of a note
// PUT /api/v1/notes/:id/position
func (h *NoteHandler) UpdatePosition(c *gin.Context) {
	userID, err := infra.UserIDFromCtx(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid note id")
		return
	}

	var req UpdatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := h.noteService.UpdatePosition(c.Request.Context(), id, req.Position, userID); err != nil {
		if err == service.ErrNoteNotFound {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		if err == service.ErrNoteUnauthorized {
			c.String(http.StatusForbidden, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
