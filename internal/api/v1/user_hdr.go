// Package v1 provides HTTP handlers for API version 1 endpoints.
package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ray-d-song/yan/internal/infra"
	"github.com/ray-d-song/yan/internal/mdw"
	"github.com/ray-d-song/yan/internal/service"
)

type UserHandler struct {
	userService service.UserService
	store       *infra.DBStore
}

func NewUserHandler(userService service.UserService, store *infra.DBStore) *UserHandler {
	return &UserHandler{
		userService: userService,
		store:       store,
	}
}

// RegisterRoutes registers all user-related routes
func (h *UserHandler) RegisterRoutes(g *gin.RouterGroup) {
	users := g.Group("/users")
	{
		users.POST("/register", h.Register)
		users.POST("/login", h.Login)
		users.POST("/logout", h.Logout)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateProfile)
		users.PUT("/:id/password", h.ChangePassword)
	}
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest represents the update profile request payload
type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// ChangePasswordRequest represents the change password request payload
type ChangePasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// Register handles user registration
// POST /api/v1/users/register
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req.Username, req.Password, req.Email)
	if err != nil {
		if err == service.ErrUsernameExists {
			c.String(http.StatusConflict, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login handles user login
// POST /api/v1/users/login
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials || err == service.ErrUserDisabled {
			c.String(http.StatusUnauthorized, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Create session
	session, err := h.store.Get(c.Request, mdw.SessionName)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to create session")
		return
	}

	session.Values[mdw.SessionKeyUserID] = user.ID
	if err := session.Save(c.Request, c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "failed to save session")
		return
	}

	c.JSON(http.StatusOK, user)
}

// Logout handles user logout
// POST /api/v1/users/logout
func (h *UserHandler) Logout(c *gin.Context) {
	session, err := h.store.Get(c.Request, mdw.SessionName)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to get session")
		return
	}

	// Set MaxAge to -1 to delete the session
	session.Options.MaxAge = -1
	if err := session.Save(c.Request, c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "failed to delete session")
		return
	}

	c.Status(http.StatusOK)
}

// GetUser retrieves a user by ID
// GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile updates user profile
// PUT /api/v1/users/:id
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid user id")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := h.userService.UpdateProfile(c.Request.Context(), user); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePassword changes user password
// PUT /api/v1/users/:id/password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid user id")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), id, req.NewPassword); err != nil {
		if err == service.ErrUserNotFound {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
