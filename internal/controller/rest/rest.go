package rest

import (
	"net/http"

	"github.com/go-chi/render"

	internalErrors "dough-calculator/internal/errors"
)

type PageInput struct {
	Offset int `in:"query=offset;default=0"`
	Limit  int `in:"query=limit;default=25"`
}

func HandlerError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	if !internalErrors.IsServiceError(err) {
		err = internalErrors.NewInternalServerErrorWrap(err, "internal server error")
	}

	handlerServiceError(w, r, err.(*internalErrors.ServiceError))
}

func handlerServiceError(w http.ResponseWriter, r *http.Request, err *internalErrors.ServiceError) {
	w.WriteHeader(err.ResponseCode)
	render.JSON(w, r, err)
}
