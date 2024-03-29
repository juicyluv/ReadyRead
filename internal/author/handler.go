package author

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
	authorsURL = "/api/authors"
	authorURL  = "/api/authors/:id"
)

// Handler handles requests specified to author service.
type Handler struct {
	logger        logger.Logger
	authorService Service
}

// NewHandler returns a new author Handler instance.
func NewHandler(logger logger.Logger, authorService Service) handler.Handling {
	return &Handler{
		logger:        logger,
		authorService: authorService,
	}
}

// Register registers new routes for router.
func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, authorURL, h.GetAuthor)
	router.HandlerFunc(http.MethodPost, authorsURL, h.CreateAuthor)
	router.HandlerFunc(http.MethodPut, authorURL, h.UpdateAuthor)
	router.HandlerFunc(http.MethodPatch, authorURL, h.UpdateAuthorPartially)
	router.HandlerFunc(http.MethodDelete, authorURL, h.DeleteAuthor)
}

// GetAuthor godoc
// @Summary Show author information
// @Description Get author by id.
// @Tags authors
// @Accept json
// @Produce json
// @Param id path int64 true "Author id"
// @Success 200 {object} Author
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /authors/{id} [get]
func (h *Handler) GetAuthor(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GET AUTHOR")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	author, err := h.authorService.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) {
			response.NotFound(w)
			return
		}
		h.logger.Error(err)
		response.InternalError(w, err.Error(), "")
		return
	}

	response.JSON(w, http.StatusOK, author)
}

// CreateAuthor godoc
// @Summary Create author
// @Description Register a new author.
// @Tags authors
// @Accept json
// @Produce json
// @Param input body CreateAuthorDTO true "JSON input"
// @Success 201 {object} Author
// @Failure 400 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /authors [post]
func (h *Handler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("CREATE AUTHOR")

	var input CreateAuthorDTO
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}

	if err := input.Validate(); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrValidationFailed.Error())
		return
	}

	author, err := h.authorService.Create(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmailTaken) {
			response.BadRequest(w, err.Error(), "")
			return
		}
		response.InternalError(w, fmt.Sprintf("cannot create author: %v", err), "")
		return
	}

	response.JSON(w, http.StatusCreated, author)
}

// UpdateAuthor godoc
// @Summary Update author
// @Description Update author with specified id.
// @Tags authors
// @Accept json
// @Produce json
// @Param id path int64 true "Author id"
// @Param input body UpdateAuthorDTO true "JSON input"
// @Success 200
// @Failure 400 {object} apperror.AppError
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /authors/{id} [put]
func (h *Handler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("UPDATE AUTHOR")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	var input UpdateAuthorDTO
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), "please, fix your request body")
		return
	}

	if err := input.Validate(); err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	input.Id = id

	err = h.authorService.Update(r.Context(), &input)
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

// UpdateAuthorPartially godoc
// @Summary Update author
// @Description Partially update author with specified id.
// @Tags authors
// @Accept json
// @Produce json
// @Param id path int64 true "Author id"
// @Param input body UpdateAuthorPartiallyDTO true "JSON input"
// @Success 200
// @Failure 400 {object} apperror.AppError
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /authors/{id} [patch]
func (h *Handler) UpdateAuthorPartially(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("UPDATE AUTHOR PARTIALLY")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	var input UpdateAuthorPartiallyDTO
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), "please, fix your request body")
		return
	}

	if err := input.Validate(); err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	input.Id = id

	err = h.authorService.UpdatePartially(r.Context(), &input)
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

// DeleteAuthor godoc
// @Summary Delete author
// @Description Delete author with specified id.
// @Tags authors
// @Accept json
// @Produce json
// @Param id path int64 true "Author id"
// @Success 200
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /authors/{id} [delete]
func (h *Handler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("DELETE AUTHOR")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	err = h.authorService.Delete(r.Context(), id)
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
