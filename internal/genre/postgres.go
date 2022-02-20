package genre

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
	tableName = "genres"
)

// Check whether db implements genre storage interface.
var _ Storage = &db{}

// db implements genre storage interface.
type db struct {
	logger         logger.Logger
	conn           *pgx.Conn
	requestTimeout time.Duration
}

// NewStorage returns a new genre storage instance.
func NewStorage(storage *pgx.Conn, requestTimeout int) Storage {
	return &db{
		logger:         logger.GetLogger(),
		conn:           storage,
		requestTimeout: time.Duration(requestTimeout) * time.Second,
	}
}

// Create inserts a genre record in the database.
// Returns an error on failure or inserted genre with it's id on success.
func (d *db) Create(genre *Genre) (*Genre, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (genre)
	VALUES ($1)
	RETURNING id`, tableName)

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	err := d.conn.QueryRow(
		ctx,
		query,
		genre.Genre,
	).Scan(&genre.Id)

	if err != nil {
		err = fmt.Errorf("failed to execute create genre query: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	return genre, nil
}

func (d *db) FindById(id int16) (*Genre, error) {
	query := fmt.Sprintf(`
	SELECT id, genre
	FROM %s 
	WHERE id = $1`, tableName)

	var found Genre

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	err := d.conn.QueryRow(ctx, query, id).Scan(
		&found.Id,
		&found.Genre,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNoRows
		}
		err = fmt.Errorf("failed to execute find genre by id query: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	return &found, nil
}

func (d *db) Update(genre *UpdateGenreDTO) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET genre = $1
	WHERE id = $2`, tableName)

	args := []interface{}{
		genre.Genre,
		genre.Id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	result, err := d.conn.Exec(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrNoRows
		}
		err = fmt.Errorf("failed to execute update genre query: %v", err)
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
		return fmt.Errorf("failed to delete author: %v", err)
	}

	return nil
}
