package genre

import (
	"context"
	"errors"

	"github.com/juicyluv/ReadyRead/internal/apperror"
	"github.com/juicyluv/ReadyRead/pkg/logger"
)

// Service describes genre service functionality.
type Service interface {
	Create(ctx context.Context, genre *CreateGenreDTO) (*Genre, error)
	GetById(ctx context.Context, id int16) (*Genre, error)
	Update(ctx context.Context, genre *UpdateGenreDTO) error
	Delete(ctx context.Context, id int16) error
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

func (s *service) Create(ctx context.Context, input *CreateGenreDTO) (*Genre, error) {
	g := Genre{
		Genre: input.Genre,
	}

	genre, err := s.storage.Create(&g)
	if err != nil {
		return nil, err
	}

	return genre, nil
}

func (s *service) GetById(ctx context.Context, id int16) (*Genre, error) {
	genre, err := s.storage.FindById(id)
	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) {
			return nil, err
		}
		s.logger.Warnf("cannot find genre by id: %v", err)
		return nil, err
	}

	return genre, nil
}

func (s *service) Update(ctx context.Context, genre *UpdateGenreDTO) error {
	a, err := s.GetById(ctx, genre.Id)
	if err != nil {
		if !errors.Is(err, apperror.ErrNoRows) {
			s.logger.Errorf("failed to get genre: %v", err)
		}
		return err
	}

	if a == nil {
		return apperror.ErrNoRows
	}

	err = s.storage.Update(genre)
	if err != nil {
		s.logger.Errorf("failed to update genre: %v", err)
		return err
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id int16) error {
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrNoRows) {
			s.logger.Warnf("failed to delete genre: %v", err)
		}
		return err
	}

	return nil
}
