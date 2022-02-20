package language

import (
	"context"
	"errors"

	"github.com/juicyluv/ReadyRead/internal/apperror"
	"github.com/juicyluv/ReadyRead/pkg/logger"
)

// Service describes language service functionality.
type Service interface {
	Create(ctx context.Context, language *CreateLanguageDTO) (*Language, error)
	GetById(ctx context.Context, id int16) (*Language, error)
	Update(ctx context.Context, language *UpdateLanguageDTO) error
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

func (s *service) Create(ctx context.Context, input *CreateLanguageDTO) (*Language, error) {
	l := Language{
		Language: input.Language,
	}

	language, err := s.storage.Create(&l)
	if err != nil {
		return nil, err
	}

	return language, nil
}

func (s *service) GetById(ctx context.Context, id int16) (*Language, error) {
	language, err := s.storage.FindById(id)
	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) {
			return nil, err
		}
		s.logger.Warnf("cannot find language by id: %v", err)
		return nil, err
	}

	return language, nil
}

func (s *service) Update(ctx context.Context, genre *UpdateLanguageDTO) error {
	l, err := s.GetById(ctx, genre.Id)
	if err != nil {
		if !errors.Is(err, apperror.ErrNoRows) {
			s.logger.Errorf("failed to get language: %v", err)
		}
		return err
	}

	if l == nil {
		return apperror.ErrNoRows
	}

	err = s.storage.Update(genre)
	if err != nil {
		s.logger.Errorf("failed to update language: %v", err)
		return err
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id int16) error {
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrNoRows) {
			s.logger.Warnf("failed to delete language: %v", err)
		}
		return err
	}

	return nil
}
