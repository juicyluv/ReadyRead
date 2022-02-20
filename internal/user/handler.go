package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/juicyluv/ReadyRead/internal/apperror"
	"github.com/juicyluv/ReadyRead/internal/handler"
	"github.com/juicyluv/ReadyRead/internal/response"
	"github.com/juicyluv/ReadyRead/pkg/logger"
	"github.com/julienschmidt/httprouter"
)

const (
	usersURL = "/api/users"
	userURL  = "/api/users/:id"
)

// Handler handles requests specified to user service.
type Handler struct {
	logger      logger.Logger
	userService Service
}

// NewHandler returns a new user Handler instance.
func NewHandler(logger logger.Logger, userService Service) handler.Handling {
	return &Handler{
		logger:      logger,
		userService: userService,
	}
}

// Register registers new routes for router.
func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, userURL, h.GetUser)
	router.HandlerFunc(http.MethodGet, usersURL, h.GetUserByEmailAndPassword)
	router.HandlerFunc(http.MethodPost, usersURL, h.CreateUser)
	router.HandlerFunc(http.MethodPut, userURL, h.UpdateUser)
	router.HandlerFunc(http.MethodPatch, userURL, h.UpdateUserPartially)
	router.HandlerFunc(http.MethodDelete, userURL, h.DeleteUser)
}

// GetUser godoc
// @Summary Show user information
// @Description Get user by id.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int64 true "User id"
// @Success 200 {object} User
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /users/{id} [get]
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GET USER")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	user, err := h.userService.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) {
			response.NotFound(w)
			return
		}
		h.logger.Error(err)
		response.InternalError(w, err.Error(), "")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

// CreateUser godoc
// @Summary Create user
// @Description Register a new user.
// @Tags users
// @Accept json
// @Produce json
// @Param input body CreateUserDTO true "JSON input"
// @Success 201 {object} openapi.CreateUserResponse
// @Failure 400 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("CREATE USER")

	var input CreateUserDTO
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}

	if err := input.Validate(); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrValidationFailed.Error())
		return
	}

	if input.Password != input.RepeatPassword {
		response.BadRequest(w, "passwords don't match", "")
		return
	}

	user, err := h.userService.Create(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmailTaken) {
			response.BadRequest(w, err.Error(), "")
			return
		}
		response.InternalError(w, fmt.Sprintf("cannot create user: %v", err), "")
		return
	}

	response.JSON(w, http.StatusCreated, user)
}

// GetUserByEmailAndPassword godoc
// @Summary Get user by email and password from query parameters
// @Description Get user by email and password.
// @Tags users
// @Accept json
// @Produce json
// @Param email query string true "user email"
// @Param password query string true "user raw password"
// @Success 200 {object} User
// @Failure 400 {object} apperror.AppError
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /users [get]
func (h *Handler) GetUserByEmailAndPassword(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GET USER BY EMAIL AND PASSWORD")

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	if email == "" || password == "" {
		response.BadRequest(w, "empty email or password", "email and password must be provided")
		return
	}

	user, err := h.userService.GetByEmailAndPassword(r.Context(), email, password)
	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error(), "")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update the user with provided current password.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int64 true "User id"
// @Param input body UpdateUserDTO true "JSON input"
// @Success 200
// @Failure 400 {object} apperror.AppError
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /users/{id} [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("UPDATE USER")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	var input UpdateUserDTO
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), "please, fix your request body")
		return
	}

	if err := input.Validate(); err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	input.Id = id

	err = h.userService.Update(r.Context(), &input)
	if err != nil {
		switch err {
		case apperror.ErrNoRows:
			response.NotFound(w)
		case apperror.ErrWrongPassword:
			response.BadRequest(w, err.Error(), "")
		default:
			response.InternalError(w, err.Error(), "")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateUserPartially godoc
// @Summary Update user
// @Description Partially update the user with provided current password.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int64 true "User id"
// @Param input body UpdateUserPartiallyDTO true "JSON input"
// @Success 200
// @Failure 400 {object} apperror.AppError
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /users/{id} [patch]
func (h *Handler) UpdateUserPartially(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("UPDATE USER PARTIALLY")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	var input UpdateUserPartiallyDTO
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), "please, fix your request body")
		return
	}

	if err := input.Validate(); err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	input.Id = id

	err = h.userService.UpdatePartially(r.Context(), &input)
	if err != nil {
		switch err {
		case apperror.ErrNoRows:
			response.NotFound(w)
		case apperror.ErrWrongPassword:
			response.BadRequest(w, err.Error(), "")
		default:
			response.InternalError(w, err.Error(), "")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete the user by id.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int64 true "User id"
// @Success 200
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("DELETE USER")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	err = h.userService.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) {
			response.NotFound(w)
			return
		}
		response.InternalError(w, err.Error(), "something went wrong on the server side")
		return
	}

	w.WriteHeader(http.StatusOK)
}
