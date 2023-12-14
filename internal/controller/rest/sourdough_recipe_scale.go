package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"dough-calculator/internal/domain"
	internalErrors "dough-calculator/internal/errors"
)

type sourdoughRecipeScaleHandler struct {
	service domain.SourdoughRecipeScaleService
}

func (handler *sourdoughRecipeScaleHandler) Scale() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		recipeId := handler.getIdParam(res, req)
		if recipeId == nil {
			return
		}

		var request domain.SourdoughRecipeScaleRequestDto

		if err := render.DecodeJSON(req.Body, &request); err != nil {
			HandlerError(res, req, errors.Wrap(err, "error while decoding request body"))
			return
		}

		recipeDto, err := handler.service.Scale(req.Context(), *recipeId, request)
		if err != nil {
			HandlerError(res, req, err)
			return
		}

		render.JSON(res, req, recipeDto)
	}
}

func (handler *sourdoughRecipeScaleHandler) getIdParam(res http.ResponseWriter, req *http.Request) *uuid.UUID {
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

func NewSourdoughRecipeScaleHandler(service domain.SourdoughRecipeScaleService) (domain.SourdoughRecipeScaleHandler, error) {
	if service == nil {
		return nil, errors.New("service is nil")
	}

	return &sourdoughRecipeScaleHandler{
		service: service,
	}, nil
}
