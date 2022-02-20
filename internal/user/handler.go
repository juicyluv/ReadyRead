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
	userURL  = "/api/users/:uuid"
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
// @Description Get user by uuid.
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User id"
// @Success 200 {object} User
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /users/{uuid} [get]
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GET USER")
}

// CreateUser godoc
// @Summary Create user
// @Description Register a new user.
// @Tags users
// @Accept json
// @Produce json
// @Param input body user.CreateUserDTO true "JSON input"
// @Success 201 {object} internal.CreateUserResponse
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
// @Description Create a new user
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
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("UPDATE USER")
}

// UpdateUserPartially godoc
// @Summary Update user
// @Description Partially update the user with provided current password.
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User id"
// @Param input body user.UpdateUserDTO true "JSON input"
// @Success 200
// @Failure 400 {object} apperror.AppError
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /users/{uuid} [patch]
func (h *Handler) UpdateUserPartially(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("UPDATE USER PARTIALLY")
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete the user by uuid.
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User id"
// @Success 200
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /users/{uuid} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("DELETE USER")
}
