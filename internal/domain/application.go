//go:generate mockgen -source=application.go -destination=mocks/application.go -package=mocks

package domain

import (
	"net/http"

	"dough-calculator/internal/config"
)

type ApplicationInitializer interface {
	Initialize() (Application, error)
}

type Application interface {
	Server() *http.Server
	Config() config.Config
}
