package response

import (
	"net/http"

	"github.com/juicyluv/ReadyRead/internal/apperror"
)

// Error is a wrapper around JSON method.
// It responses with specified error and http code.
func Error(w http.ResponseWriter, code int, message, developerMessage string) {
	appError := apperror.NewAppError(code, message, developerMessage)
	JSON(w, code, appError)
}

// BadRequest is a wrapper around Error method.
// Responses with 400 Bad Request status code and specified error message.
func BadRequest(w http.ResponseWriter, message, developerMessage string) {
	Error(w, http.StatusBadRequest, message, developerMessage)
}

// Not Found is a wrapper around JSON method.
// Responses with 404 Not Found status code and specified error message.
func NotFound(w http.ResponseWriter) {
	JSON(w, http.StatusNotFound, apperror.ErrNotFound)
}

// InternalError is a wrapper around Error method.
// Responses with 500 Internal Server Error status code and specified error message.
func InternalError(w http.ResponseWriter, message, developerMessage string) {
	Error(w, http.StatusInternalServerError, message, developerMessage)
}
