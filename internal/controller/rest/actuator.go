package rest

import (
	"net/http"

	"github.com/go-chi/render"

	"dough-calculator/internal/domain"
)

type actuatorHandler struct{}

func (actuatorHandler *actuatorHandler) Health() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		render.JSON(writer, request, render.M{
			"status": "UP",
		})
	}
}

func NewActuatorHandler() domain.ActuatorHandler {
	return &actuatorHandler{}
}
