// Package service contains business logic and service layer implementations.
package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ray-d-song/yan/internal/repo"
	"golang.org/x/crypto/bcrypt"

	"github.com/ray-d-song/yan/internal/model"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserDisabled       = errors.New("user is disabled")
)

//
// UserService interface
//

type UserService interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Register(ctx context.Context, username, password, email string) (*model.User, error)
	Login(ctx context.Context, email, password string) (*model.User, error)
	UpdateProfile(ctx context.Context, u *model.User) error
	ChangePassword(ctx context.Context, userID int64, newPassword string) error
}

//
// implementation
//

type userService struct {
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

//
// read
//

func (s *userService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	u, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrUserNotFound
	}
	return u, nil
}

//
// write / business logic
//

func (s *userService) Register(
	ctx context.Context,
	username, password, email string,
) (*model.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	if username == "" || password == "" || email == "" {
		return nil, errors.New("username/password/email required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		Username:     username,
		PasswordHash: string(hash),
		Email:        email,
		Status:       model.UserStatusNormal,
		IsAdmin:      model.UserAdminFalse,
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (*model.User, error) {
	u, err := s.userRepo.GetByEmail(ctx, strings.TrimSpace(email))
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrInvalidCredentials
	}

	if u.Status != model.UserStatusNormal {
		return nil, ErrUserDisabled
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(u.PasswordHash),
		[]byte(password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	return u, nil
}

func (s *userService) UpdateProfile(ctx context.Context, u *model.User) error {
	if u == nil || u.ID == 0 {
		return errors.New("invalid user")
	}

	existing, err := s.userRepo.GetByID(ctx, u.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrUserNotFound
	}

	return s.userRepo.Update(ctx, u)
}

func (s *userService) ChangePassword(
	ctx context.Context,
	userID int64,
	newPassword string,
) error {
	if newPassword == "" {
		return errors.New("password required")
	}

	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if u == nil {
		return ErrUserNotFound
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(newPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, userID, string(hash))
}
