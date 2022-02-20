package genre

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
	genresURL = "/api/genres"
	genreURL  = "/api/genres/:id"
)

// Handler handles requests specified to genre service.
type Handler struct {
	logger       logger.Logger
	genreService Service
}

// NewHandler returns a new genre Handler instance.
func NewHandler(logger logger.Logger, genreService Service) handler.Handling {
	return &Handler{
		logger:       logger,
		genreService: genreService,
	}
}

// Register registers new routes for router.
func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, genreURL, h.GetGenre)
	router.HandlerFunc(http.MethodPost, genresURL, h.CreateGenre)
	router.HandlerFunc(http.MethodPut, genreURL, h.UpdateGenre)
	router.HandlerFunc(http.MethodDelete, genreURL, h.DeleteGenre)
}

// GetGenre godoc
// @Summary Show genre information
// @Description Get genre by id.
// @Tags genres
// @Accept json
// @Produce json
// @Param id path int64 true "Genre id"
// @Success 200 {object} Genre
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /genres/{id} [get]
func (h *Handler) GetGenre(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GET GENRE")

	id, err := handler.ReadIdParam16(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	author, err := h.genreService.GetById(r.Context(), id)
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

// CreateGenre godoc
// @Summary Create genre
// @Description Insert genre in database.
// @Tags genres
// @Accept json
// @Produce json
// @Param input body CreateGenreDTO true "JSON input"
// @Success 201 {object} Genre
// @Failure 400 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /genres [post]
func (h *Handler) CreateGenre(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("CREATE GENRE")

	var input CreateGenreDTO
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}

	if err := input.Validate(); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrValidationFailed.Error())
		return
	}

	author, err := h.genreService.Create(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmailTaken) {
			response.BadRequest(w, err.Error(), "")
			return
		}
		response.InternalError(w, fmt.Sprintf("cannot create genre: %v", err), "")
		return
	}

	response.JSON(w, http.StatusCreated, author)
}

// UpdateGenre godoc
// @Summary Update genre
// @Description Update genre with specified id.
// @Tags genres
// @Accept json
// @Produce json
// @Param id path int64 true "Genre id"
// @Param input body UpdateGenreDTO true "JSON input"
// @Success 200
// @Failure 400 {object} apperror.AppError
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /genres/{id} [put]
func (h *Handler) UpdateGenre(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("UPDATE GENRE")

	id, err := handler.ReadIdParam16(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	var input UpdateGenreDTO
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), "please, fix your request body")
		return
	}

	if err := input.Validate(); err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	input.Id = id

	err = h.genreService.Update(r.Context(), &input)
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

// DeleteGenre godoc
// @Summary Delete genre
// @Description Delete genre with specified id.
// @Tags genres
// @Accept json
// @Produce json
// @Param id path int64 true "Genre id"
// @Success 200
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /genres/{id} [delete]
func (h *Handler) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("DELETE GENRE")

	id, err := handler.ReadIdParam16(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	err = h.genreService.Delete(r.Context(), id)
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
