package rest

import (
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"dough-calculator/internal/domain"
	internalErrors "dough-calculator/internal/errors"
)

const (
	flourIdNotFound = 20001
	flourIdNotValid = 20002
)

type SearchFlourInput struct {
	Name string `in:"query=name"`
}

type flourHandler struct {
	service domain.FlourService
}

func (handler *flourHandler) Create() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var request domain.CreateFlourRequest

		if err := render.DecodeJSON(req.Body, &request); err != nil {
			HandlerError(res, req, errors.Wrap(err, "error while decoding request body"))
			return
		}

		recipeDto, err := handler.service.Create(req.Context(), request)
		if err != nil {
			HandlerError(res, req, err)
			return
		}

		render.Status(req, http.StatusCreated)
		render.JSON(res, req, recipeDto)
	}
}

func (handler *flourHandler) FindById() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		flourId := handler.getIdParam(res, req)
		if flourId == nil {
			return
		}

		recipeDto, err := handler.service.FindById(req.Context(), *flourId)
		if err != nil {
			HandlerError(res, req, err)
			return
		}

		render.JSON(res, req, recipeDto)
	}
}

func (handler *flourHandler) getIdParam(res http.ResponseWriter, req *http.Request) *uuid.UUID {
	param := chi.URLParam(req, "id")
	if param == "" {
		HandlerError(res, req, internalErrors.NewBadRequestError(flourIdNotFound, "id is required", "id is required"))
		return nil
	}
	id, err := uuid.Parse(param)
	if err != nil {
		HandlerError(res, req, internalErrors.NewBadRequestError(flourIdNotValid, "id is not valid", "id is not valid"))
		return nil
	}
	return &id
}

func (handler *flourHandler) Find() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		page := req.Context().Value(httpin.Input).(*PageInput)

		flourDtos, err := handler.service.Find(req.Context(), page.Offset, page.Limit)
		if err != nil {
			HandlerError(res, req, err)
			return
		}

		render.JSON(res, req, flourDtos)
	}
}

func (handler *flourHandler) Search() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		search := req.Context().Value(httpin.Input).(*SearchFlourInput)

		flourDtos, err := handler.service.SearchByName(req.Context(), search.Name)
		if err != nil {
			HandlerError(res, req, err)
			return
		}

		render.JSON(res, req, flourDtos)
	}
}

func NewFlourHandler(service domain.FlourService) (domain.FlourHandler, error) {
	if service == nil {
		return nil, errors.New("service is nil")
	}

	return &flourHandler{service: service}, nil
}
