package language

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
	languagesURL = "/api/languages"
	languageURL  = "/api/languages/:id"
)

// Handler handles requests specified to language service.
type Handler struct {
	logger          logger.Logger
	languageService Service
}

// NewHandler returns a new language Handler instance.
func NewHandler(logger logger.Logger, languageService Service) handler.Handling {
	return &Handler{
		logger:          logger,
		languageService: languageService,
	}
}

// Register registers new routes for router.
func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, languageURL, h.GetLanguage)
	router.HandlerFunc(http.MethodPost, languagesURL, h.CreateLanguage)
	router.HandlerFunc(http.MethodPut, languagesURL, h.UpdateLanguage)
	router.HandlerFunc(http.MethodDelete, languageURL, h.DeleteLanguage)
}

// GetLanguage godoc
// @Summary Show language information
// @Description Get language by id.
// @Tags languages
// @Accept json
// @Produce json
// @Param id path int64 true "Language id"
// @Success 200 {object} Language
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /languages/{id} [get]
func (h *Handler) GetLanguage(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GET LANGUAGE")

	id, err := handler.ReadIdParam16(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	author, err := h.languageService.GetById(r.Context(), id)
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

// CreateLanguage godoc
// @Summary Create language
// @Description Insert language in database.
// @Tags languages
// @Accept json
// @Produce json
// @Param input body CreateLanguageDTO true "JSON input"
// @Success 201 {object} Language
// @Failure 400 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /languages [post]
func (h *Handler) CreateLanguage(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("CREATE LANGUAGE")

	var input CreateLanguageDTO
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}

	if err := input.Validate(); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrValidationFailed.Error())
		return
	}

	author, err := h.languageService.Create(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmailTaken) {
			response.BadRequest(w, err.Error(), "")
			return
		}
		response.InternalError(w, fmt.Sprintf("cannot create language: %v", err), "")
		return
	}

	response.JSON(w, http.StatusCreated, author)
}

// UpdateLanguage godoc
// @Summary Update language
// @Description Update language with specified id.
// @Tags languages
// @Accept json
// @Produce json
// @Param id path int64 true "Language id"
// @Param input body UpdateLanguageDTO true "JSON input"
// @Success 200
// @Failure 400 {object} apperror.AppError
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /languages/{id} [put]
func (h *Handler) UpdateLanguage(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("UPDATE LANGUAGE")

	id, err := handler.ReadIdParam16(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	var input UpdateLanguageDTO
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), "please, fix your request body")
		return
	}

	if err := input.Validate(); err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	input.Id = id

	err = h.languageService.Update(r.Context(), &input)
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

// DeleteLanguage godoc
// @Summary Delete genre
// @Description Delete language with specified id.
// @Tags languages
// @Accept json
// @Produce json
// @Param id path int64 true "Language id"
// @Success 200
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /languages/{id} [delete]
func (h *Handler) DeleteLanguage(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("DELETE LANGUAGE")

	id, err := handler.ReadIdParam16(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	err = h.languageService.Delete(r.Context(), id)
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
