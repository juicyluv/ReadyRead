package author

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/juicyluv/ReadyRead/internal/apperror"
	"github.com/juicyluv/ReadyRead/pkg/logger"
)

const (
	tableName = "authors"
)

// Check whether db implements author storage interface.
var _ Storage = &db{}

// db implementes user storage interface.
type db struct {
	logger         logger.Logger
	conn           *pgx.Conn
	requestTimeout time.Duration
}

// NewAuthorStorage returns a new author storage instance.
func NewAuthorStorage(storage *pgx.Conn, requestTimeout int) Storage {
	return &db{
		logger:         logger.GetLogger(),
		conn:           storage,
		requestTimeout: time.Duration(requestTimeout) * time.Second,
	}
}

// Create inserts a author record in the database.
// Returns an error on failure or inserted author with it's id on success.
func (d *db) Create(author *Author) (*Author, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (name, surname)
	VALUES ($1, $2)
	RETURNING id`, tableName)

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	err := d.conn.QueryRow(
		ctx,
		query,
		author.Name,
		author.Surname,
	).Scan(&author.Id)

	if err != nil {
		err = fmt.Errorf("failed to execute create author query: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	return author, nil
}

func (d *db) FindById(id int64) (*Author, error) {
	query := fmt.Sprintf(`
	SELECT id, name, surname
	FROM %s 
	WHERE id = $1`, tableName)

	var found Author

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	err := d.conn.QueryRow(ctx, query, id).Scan(
		&found.Id,
		&found.Name,
		&found.Surname,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNoRows
		}
		err = fmt.Errorf("failed to execute find author by id query: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	return &found, nil
}

func (d *db) Update(author *UpdateAuthorDTO) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET name=$1, surname=$2
	WHERE id = $3`, tableName)

	args := []interface{}{
		author.Name,
		author.Surname,
		author.Id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	result, err := d.conn.Exec(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrNoRows
		}
		err = fmt.Errorf("failed to execute update author query: %v", err)
		d.logger.Error(err)
		return err
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrNoRows
	}

	return nil
}

func (d *db) UpdatePartially(author *UpdateAuthorPartiallyDTO) error {
	values := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if author.Name != nil {
		values = append(values, fmt.Sprintf("name=$%d", argId))
		args = append(args, *author.Name)
		argId++
	}

	if author.Surname != nil {
		values = append(values, fmt.Sprintf("surname=$%d", argId))
		args = append(args, *author.Surname)
		argId++
	}

	valuesQuery := strings.Join(values, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", tableName, valuesQuery, argId)
	args = append(args, author.Id)

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	_, err := d.conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update author partially: %v", err)
	}

	return nil
}

func (d *db) Delete(id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	_, err := d.conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete author: %v", err)
	}

	return nil
}
