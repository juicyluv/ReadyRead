package user

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
	tableName = "users"
)

// Check whether db implements user storage interface.
var _ Storage = &db{}

// db implementes user storage interface.
type db struct {
	logger         logger.Logger
	conn           *pgx.Conn
	requestTimeout time.Duration
}

// NewUserStorage returns a new user storage instance.
func NewUserStorage(storage *pgx.Conn, requestTimeout int) Storage {
	return &db{
		logger:         logger.GetLogger(),
		conn:           storage,
		requestTimeout: time.Duration(requestTimeout) * time.Second,
	}
}

// Create inserts a user record in the database.
// Returns an error on failure or inserted user with it's id on success.
func (d *db) Create(user *User) (*User, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id, TO_CHAR(registered_at, 'DD-MM-YYYY')`, tableName)

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	err := d.conn.QueryRow(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&user.Id, &user.RegisteredAt)
	if err != nil {
		err = fmt.Errorf("failed to execute create user query: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	return user, nil
}

func (d *db) FindByEmail(email string) (*User, error) {
	query := fmt.Sprintf(`
	SELECT id, username, email, password, verified, address, phone_number, TO_CHAR(registered_at, 'DD-MM-YYYY')
	FROM %s 
	WHERE email = $1`, tableName)

	var found User

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	err := d.conn.QueryRow(ctx, query, email).Scan(
		&found.Id,
		&found.Username,
		&found.Email,
		&found.Password,
		&found.Verified,
		&found.Address,
		&found.PhoneNumber,
		&found.RegisteredAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNoRows
		}
		err = fmt.Errorf("failed to execute find user by email query: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	return &found, nil
}

func (d *db) FindByUsername(username string) (*User, error) {
	query := fmt.Sprintf(`
	SELECT id, username, email, password, verified, address, phone_number, TO_CHAR(registered_at, 'DD-MM-YYYY')
	FROM %s 
	WHERE username = $1`, tableName)

	var found User

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	err := d.conn.QueryRow(ctx, query, username).Scan(
		&found.Id,
		&found.Username,
		&found.Email,
		&found.Password,
		&found.Verified,
		&found.Address,
		&found.PhoneNumber,
		&found.RegisteredAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNoRows
		}
		err = fmt.Errorf("failed to execute find user by username query: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	return &found, nil
}

func (d *db) FindById(id int64) (*User, error) {
	query := fmt.Sprintf(`
	SELECT id, username, email, password, verified, address, phone_number, TO_CHAR(registered_at, 'DD-MM-YYYY')
	FROM %s 
	WHERE id = $1`, tableName)

	var found User

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	err := d.conn.QueryRow(ctx, query, id).Scan(
		&found.Id,
		&found.Username,
		&found.Email,
		&found.Password,
		&found.Verified,
		&found.Address,
		&found.PhoneNumber,
		&found.RegisteredAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNoRows
		}
		err = fmt.Errorf("failed to execute find user by id query: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	return &found, nil
}

func (d *db) Update(user *UpdateUserDTO) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET username=$1, email=$2, password=$3, address=$4, phone_number=$5
	WHERE id = $6`, tableName)

	args := []interface{}{
		user.Username,
		user.Email,
		user.Password,
		user.Address,
		user.PhoneNumber,
		user.Id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	result, err := d.conn.Exec(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrNoRows
		}
		err = fmt.Errorf("failed to execute update user query: %v", err)
		d.logger.Error(err)
		return err
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrNoRows
	}

	return nil
}

func (d *db) UpdatePartially(user *UpdateUserPartiallyDTO) error {
	values := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if user.Username != nil {
		values = append(values, fmt.Sprintf("username=$%d", argId))
		args = append(args, *user.Username)
		argId++
	}

	if user.Email != nil {
		values = append(values, fmt.Sprintf("email=$%d", argId))
		args = append(args, *user.Email)
		argId++
	}

	if user.Address != nil {
		values = append(values, fmt.Sprintf("address=$%d", argId))
		args = append(args, *user.Address)
		argId++
	}

	if user.PhoneNumber != nil {
		values = append(values, fmt.Sprintf("phone_number=$%d", argId))
		args = append(args, *user.PhoneNumber)
		argId++
	}

	if user.NewPassword != nil {
		values = append(values, fmt.Sprintf("password=$%d", argId))
		args = append(args, *user.NewPassword)
		argId++
	}

	valuesQuery := strings.Join(values, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", tableName, valuesQuery, argId)
	args = append(args, user.Id)

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	_, err := d.conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user partially: %v", err)
	}

	return nil
}

func (d *db) Delete(id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	_, err := d.conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}
