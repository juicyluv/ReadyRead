package apperror

import (
	"errors"
	"net/http"

	"github.com/juicyluv/ReadyRead/pkg/apperror"
)

var (
	// ErrNotFound is used when server needs to response with 404 Not Found status code.
	ErrNotFound = apperror.NewAppError(
		http.StatusNotFound,
		"requested resource is not found",
		"maybe you have an error in your request or requested resource not found",
	)

	// ErrNoRows is used when no rows returned from storage.
	ErrNoRows = errors.New("no rows")

	// ErrValidationFailed is used when input validation failed.
	ErrValidationFailed = errors.New("input validation failed. please, provide valid values")

	// ErrEmailTaken is used when the user is being created and given email is already taken.
	ErrEmailTaken = errors.New("email already taken")

	// ErrWrongPassword is used when user entered wrong password.
	ErrWrongPassword = errors.New("wrong email or password")
)
