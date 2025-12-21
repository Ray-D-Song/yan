package infra

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/ray-d-song/yan/internal/model"
)

const (
	// ContextKeyUser is the key for storing user in context
	ContextKeyUser = "user"
	// ContextKeyUserID is the key for storing user ID in context
	ContextKeyUserID = "user_id"
)

var (
	ErrUserNotInContext = errors.New("user not found in context")
)

// SetUserInContext stores user information in gin context
func SetUserInContext(c *gin.Context, user *model.User) {
	c.Set(ContextKeyUser, user)
	c.Set(ContextKeyUserID, user.ID)
}

// UserFromCtx retrieves user from gin context
func UserFromCtx(c *gin.Context) (*model.User, error) {
	val, exists := c.Get(ContextKeyUser)
	if !exists {
		return nil, ErrUserNotInContext
	}

	user, ok := val.(*model.User)
	if !ok {
		return nil, ErrUserNotInContext
	}

	return user, nil
}

// UserIDFromCtx retrieves user ID from gin context
func UserIDFromCtx(c *gin.Context) (int64, error) {
	val, exists := c.Get(ContextKeyUserID)
	if !exists {
		return 0, ErrUserNotInContext
	}

	userID, ok := val.(int64)
	if !ok {
		return 0, ErrUserNotInContext
	}

	return userID, nil
}

// MustUserFromCtx retrieves user from context or panics
// Use this only in handlers where auth middleware is guaranteed to run
func MustUserFromCtx(c *gin.Context) *model.User {
	user, err := UserFromCtx(c)
	if err != nil {
		panic(err)
	}
	return user
}

// MustUserIDFromCtx retrieves user ID from context or panics
// Use this only in handlers where auth middleware is guaranteed to run
func MustUserIDFromCtx(c *gin.Context) int64 {
	userID, err := UserIDFromCtx(c)
	if err != nil {
		panic(err)
	}
	return userID
}

// UserFromStdContext retrieves user from standard context
func UserFromStdContext(ctx context.Context) (*model.User, error) {
	val := ctx.Value(ContextKeyUser)
	if val == nil {
		return nil, ErrUserNotInContext
	}

	user, ok := val.(*model.User)
	if !ok {
		return nil, ErrUserNotInContext
	}

	return user, nil
}
