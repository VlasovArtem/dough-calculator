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
	recipeIdNotFound = 10001
	recipeIdNotValid = 10002
)

type SearchRecipeInput struct {
	Name string `in:"query=name"`
}

type sourdoughRecipeHandler struct {
	service domain.SourdoughRecipeService
}

func (handler *sourdoughRecipeHandler) Create() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var request domain.CreateSourdoughRecipeRequest

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

func (handler *sourdoughRecipeHandler) FindById() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		recipeId := handler.getIdParam(res, req)
		if recipeId == nil {
			return
		}

		recipeDto, err := handler.service.FindById(req.Context(), *recipeId)
		if err != nil {
			HandlerError(res, req, err)
			return
		}

		render.JSON(res, req, recipeDto)
	}
}

func (handler *sourdoughRecipeHandler) getIdParam(res http.ResponseWriter, req *http.Request) *uuid.UUID {
	param := chi.URLParam(req, "id")
	if param == "" {
		HandlerError(res, req, internalErrors.NewBadRequestError(recipeIdNotFound, "id is required", "id is required"))
		return nil
	}
	id, err := uuid.Parse(param)
	if err != nil {
		HandlerError(res, req, internalErrors.NewBadRequestError(recipeIdNotValid, "id is not valid", "id is not valid"))
		return nil
	}
	return &id
}

func (handler *sourdoughRecipeHandler) Find() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		page := req.Context().Value(httpin.Input).(*PageInput)

		recipes, err := handler.service.Find(req.Context(), page.Offset, page.Limit)
		if err != nil {
			HandlerError(res, req, err)
			return
		}

		render.JSON(res, req, recipes)
	}
}

func (handler *sourdoughRecipeHandler) Search() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		name := req.Context().Value(httpin.Input).(*SearchRecipeInput)

		recipes, err := handler.service.SearchByName(req.Context(), name.Name)
		if err != nil {
			HandlerError(res, req, err)
			return
		}

		render.JSON(res, req, recipes)
	}
}

func NewSourdoughRecipeHandler(sourdoughRecipeService domain.SourdoughRecipeService) (domain.SourdoughRecipeHandler, error) {
	if sourdoughRecipeService == nil {
		return nil, errors.New("service cannot be nil")
	}

	return &sourdoughRecipeHandler{service: sourdoughRecipeService}, nil
}
