package user

import (
	"context"

	"github.com/juicyluv/ReadyRead/pkg/logger"
)

// Service describes user service functionality.
type Service interface {
	Create(ctx context.Context, user *CreateUserDTO) (*User, error)
	GetByEmailAndPassword(ctx context.Context, email, password string) (*User, error)
	GetById(ctx context.Context, id int64) (*User, error)
	UpdatePartially(ctx context.Context, user *UpdateUserDTO) error
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

// Create will check whether provided email already taken.
// If it is, returns an error. Then it will hash user password
// and try to insert the user. Returns inserted UUID or an error
// on failure.
func (s *service) Create(ctx context.Context, input *CreateUserDTO) (*User, error) {
	return nil, nil
}

// GetByEmailAndPassword will find a user with provided email.
// If there's no such user with this email, returns No Rows error.
// If password doesn't match, returns Wrong Password error.
// Returns a user if everything is OK.
func (s *service) GetByEmailAndPassword(ctx context.Context, email, password string) (*User, error) {
	return nil, nil
}

// GetById will find a user with specified uuid in storage.
// Returns an error on failure of there's no user with this uuid.
func (s *service) GetById(ctx context.Context, id int64) (*User, error) {
	return nil, nil
}

// UpdatePartially will find the user with provided uuid.
// If there is no user with such id, returns No Rows error.
// Then passwords will be compared. If it don't match, returns
// Wrong Password error. Then updates the user. If something went wrong,
// returns an error and nil if everything is OK.
func (s *service) UpdatePartially(ctx context.Context, user *UpdateUserDTO) error {
	return nil
}

// Delete tries to delete the user with provided uuid.
// Returns an error on failure or nil if query has been executed.
func (s *service) Delete(ctx context.Context, id int64) error {
	return nil
}
