package language

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/juicyluv/ReadyRead/internal/apperror"
	"github.com/juicyluv/ReadyRead/pkg/logger"
)

const (
	tableName = "languages"
)

// Check whether db implements language storage interface.
var _ Storage = &db{}

// db implements language storage interface.
type db struct {
	logger         logger.Logger
	conn           *pgx.Conn
	requestTimeout time.Duration
}

// NewStorage returns a new language storage instance.
func NewStorage(storage *pgx.Conn, requestTimeout int) Storage {
	return &db{
		logger:         logger.GetLogger(),
		conn:           storage,
		requestTimeout: time.Duration(requestTimeout) * time.Second,
	}
}

// Create inserts a language record in the database.
// Returns an error on failure or inserted genre with it's id on success.
func (d *db) Create(language *Language) (*Language, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (language)
	VALUES ($1)
	RETURNING id`, tableName)

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	err := d.conn.QueryRow(
		ctx,
		query,
		language.Language,
	).Scan(&language.Id)

	if err != nil {
		err = fmt.Errorf("failed to execute create language query: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	return language, nil
}

func (d *db) FindById(id int16) (*Language, error) {
	query := fmt.Sprintf(`
	SELECT id, language
	FROM %s 
	WHERE id = $1`, tableName)

	var found Language

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	err := d.conn.QueryRow(ctx, query, id).Scan(
		&found.Id,
		&found.Language,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNoRows
		}
		err = fmt.Errorf("failed to execute find language by id query: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	return &found, nil
}

func (d *db) Update(language *UpdateLanguageDTO) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET language = $1
	WHERE id = $2`, tableName)

	args := []interface{}{
		language.Language,
		language.Id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	result, err := d.conn.Exec(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrNoRows
		}
		err = fmt.Errorf("failed to execute update language query: %v", err)
		d.logger.Error(err)
		return err
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrNoRows
	}

	return nil
}

func (d *db) Delete(id int16) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	_, err := d.conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete language: %v", err)
	}

	return nil
}
