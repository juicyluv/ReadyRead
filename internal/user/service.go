package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/juicyluv/ReadyRead/internal/apperror"
	"github.com/juicyluv/ReadyRead/pkg/logger"
)

// Service describes user service functionality.
type Service interface {
	Create(ctx context.Context, user *CreateUserDTO) (*User, error)
	GetByEmailAndPassword(ctx context.Context, email, password string) (*User, error)
	GetById(ctx context.Context, id int64) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *UpdateUserDTO) error
	UpdatePartially(ctx context.Context, user *UpdateUserPartiallyDTO) error
	Delete(ctx context.Context, id int64) error
}

type service struct {
	logger  logger.Logger
	storage Storage
}

// NewService returns a new instance that implements Service interface.
func NewService(storage Storage, logger logger.Logger) Service {
	return &service{
		logger:  logger,
		storage: storage,
	}
}

// CreateUser inserts a new user record in storage. Returns inserted user on success
// or an error on failure.
func (s *service) Create(ctx context.Context, input *CreateUserDTO) (*User, error) {
	found, err := s.storage.FindByEmail(input.Email)
	if err != nil {
		if !errors.Is(err, apperror.ErrNoRows) {
			return nil, err
		}
	}

	if found != nil {
		return nil, apperror.ErrEmailTaken
	}

	u := User{
		Email:    input.Email,
		Username: input.Username,
		Password: input.Password,
	}

	err = u.HashPassword()
	if err != nil {
		return nil, fmt.Errorf("cannot hash password")
	}

	user, err := s.storage.Create(&u)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmailAndPassword finds a user record in storage by email and validates specified password.
// Returns ErrWrongPassword if passwords don't match. Returns an error on failure.
func (s *service) GetByEmailAndPassword(ctx context.Context, email, password string) (*User, error) {
	user, err := s.storage.FindByEmail(email)
	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) {
			return nil, err
		}
		s.logger.Warnf("cannot find user by email: %v", err)
		return nil, err
	}

	if !user.ComparePassword(password) {
		return nil, apperror.ErrWrongPassword
	}

	return user, nil
}

// GetById finds a user record in storage by specified id.
// Returns ErrNoRows user with this id doesn't exist.
// Returns an error on failure.
func (s *service) GetById(ctx context.Context, id int64) (*User, error) {
	user, err := s.storage.FindById(id)
	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) {
			return nil, err
		}
		s.logger.Warnf("cannot find user by id: %v", err)
		return nil, err
	}

	return user, nil
}

// GetByUsername finds a user record in storage by specified username.
// Returns ErrNoRows user with this username doesn't exist.
// Returns an error on failure.
func (s *service) GetByUsername(ctx context.Context, username string) (*User, error) {
	user, err := s.storage.FindByUsername(username)
	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) {
			return nil, err
		}
		s.logger.Warnf("cannot find user by username: %v", err)
		return nil, err
	}

	return user, nil
}

// Update updates a user record in storage by specified id.
// Returns ErrNoRows user with this id doesn't exist,
// ErrWrongPassword if passwords don't match or an error on failure.
func (s *service) Update(ctx context.Context, user *UpdateUserDTO) error {
	u, err := s.GetById(ctx, user.Id)
	if err != nil {
		if !errors.Is(err, apperror.ErrNoRows) {
			s.logger.Errorf("failed to get user: %v", err)
		}
		return err
	}

	if !u.ComparePassword(user.Password) {
		return apperror.ErrWrongPassword
	}

	err = s.storage.Update(user)
	if err != nil {
		s.logger.Errorf("failed to update user: %v", err)
		return err
	}

	return nil
}

// Update partially updates a user record in storage by specified id.
// Returns ErrNoRows user with this id doesn't exist,
// ErrWrongPassword if passwords don't match or an error on failure.
func (s *service) UpdatePartially(ctx context.Context, user *UpdateUserPartiallyDTO) error {
	u, err := s.GetById(ctx, user.Id)
	if err != nil {
		if !errors.Is(err, apperror.ErrNoRows) {
			s.logger.Errorf("failed to get user: %v", err)
		}
		return err
	}

	if !u.ComparePassword(*user.OldPassword) {
		return apperror.ErrWrongPassword
	}

	if user.NewPassword != nil {
		u.Password = *user.NewPassword
		err = user.HashPassword()
		if err != nil {
			s.logger.Errorf("failed ot hash password: %v", err)
			return err
		}
	}

	err = s.storage.UpdatePartially(user)
	if err != nil {
		s.logger.Errorf("failed to partially update user: %v", err)
		return err
	}

	return nil
}

// Delete deletes a user record in storage by specified id.
// Returns ErrNoRows user with this id doesn't exist, or an error on failure.
func (s *service) Delete(ctx context.Context, id int64) error {
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrNoRows) {
			s.logger.Warnf("failed to delete user: %v", err)
		}
		return err
	}

	return nil
}
