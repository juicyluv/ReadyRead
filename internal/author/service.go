package author

import (
	"context"
	"errors"

	"github.com/juicyluv/ReadyRead/internal/apperror"
	"github.com/juicyluv/ReadyRead/pkg/logger"
)

// Service describes author service functionality.
type Service interface {
	Create(ctx context.Context, author *CreateAuthorDTO) (*Author, error)
	GetById(ctx context.Context, id int64) (*Author, error)
	Update(ctx context.Context, author *UpdateAuthorDTO) error
	UpdatePartially(ctx context.Context, author *UpdateAuthorPartiallyDTO) error
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

func (s *service) Create(ctx context.Context, input *CreateAuthorDTO) (*Author, error) {
	a := Author{
		Name:    input.Name,
		Surname: input.Surname,
	}

	author, err := s.storage.Create(&a)
	if err != nil {
		return nil, err
	}

	return author, nil
}

func (s *service) GetById(ctx context.Context, id int64) (*Author, error) {
	author, err := s.storage.FindById(id)
	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) {
			return nil, err
		}
		s.logger.Warnf("cannot find author by id: %v", err)
		return nil, err
	}

	return author, nil
}

func (s *service) Update(ctx context.Context, author *UpdateAuthorDTO) error {
	a, err := s.GetById(ctx, author.Id)
	if err != nil {
		if !errors.Is(err, apperror.ErrNoRows) {
			s.logger.Errorf("failed to get author: %v", err)
		}
		return err
	}

	if a == nil {
		return apperror.ErrNoRows
	}

	err = s.storage.Update(author)
	if err != nil {
		s.logger.Errorf("failed to update user: %v", err)
		return err
	}

	return nil
}

func (s *service) UpdatePartially(ctx context.Context, author *UpdateAuthorPartiallyDTO) error {
	a, err := s.GetById(ctx, author.Id)
	if err != nil {
		if !errors.Is(err, apperror.ErrNoRows) {
			s.logger.Errorf("failed to get user: %v", err)
		}
		return err
	}

	if a == nil {
		return apperror.ErrNoRows
	}

	err = s.storage.UpdatePartially(author)
	if err != nil {
		s.logger.Errorf("failed to partially update author: %v", err)
		return err
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrNoRows) {
			s.logger.Warnf("failed to delete author: %v", err)
		}
		return err
	}

	return nil
}
