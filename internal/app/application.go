//go:generate oapi-codegen -package=integration_test -o=integration_test/app_client_test.go ../../api/openapi.yaml

package app

import (
	"context"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"dough-calculator/internal/app/dependency"
	"dough-calculator/internal/config"
	"dough-calculator/internal/controller/rest"
	internalMiddleware "dough-calculator/internal/controller/rest/middleware"
	"dough-calculator/internal/domain"
)

type applicationInitializer struct {
	dependencyManager domain.DependencyManager
}

func (initializer *applicationInitializer) Initialize() (domain.Application, error) {
	err := initializer.dependencyManager.Initialize(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize dependency")
	}

	restConfig := initializer.getConfig().Application.Rest

	server := &http.Server{
		Handler:      initializer.initializeRouter(),
		IdleTimeout:  restConfig.IdleTimeoutDuration(),
		ReadTimeout:  restConfig.ReadTimeoutDuration(),
		WriteTimeout: restConfig.WriteTimeoutDuration(),
	}

	return &application{initializer.dependencyManager, server}, nil
}

func (initializer *applicationInitializer) initializeRouter() *chi.Mux {
	chiRouter := chi.NewRouter()

	chiRouter.Use(internalMiddleware.LoggerMiddleware(log.Logger))

	chiRouter.Get("/actuator/health", initializer.dependencyManager.Common().Actuator().Health())

	initializer.mountAPIRoutes(chiRouter)

	return chiRouter
}

func (initializer *applicationInitializer) mountAPIRoutes(router *chi.Mux) {
	restConfig := initializer.getConfig().Application.Rest

	router.Route(restConfig.ContextPath, func(contextPathRouter chi.Router) {
		contextPathRouter.Route("/recipe/sourdough", func(sourdoughRecipeRouter chi.Router) {
			initializer.mountSourdoughRecipeAPIRoutes(sourdoughRecipeRouter)
			initializer.mountSourdoughRecipeScaleAPIRoutes(sourdoughRecipeRouter)
		})
		contextPathRouter.Route("/flour", func(flourRouter chi.Router) {
			initializer.mountFlourAPIRoutes(flourRouter)
		})
	})

}

func (initializer *applicationInitializer) mountSourdoughRecipeAPIRoutes(router chi.Router) {
	sourdoughRecipeHandler := initializer.dependencyManager.SourdoughRecipe().Router()

	router.
		With(httpin.NewInput(rest.PageInput{})).
		Get("/", sourdoughRecipeHandler.Find())
	router.Post("/", sourdoughRecipeHandler.Create())
	router.Route("/{id}", func(idRouter chi.Router) {
		idRouter.Get("/", sourdoughRecipeHandler.FindById())
	})
	router.
		With(httpin.NewInput(rest.SearchRecipeInput{})).
		Get("/search", sourdoughRecipeHandler.Search())
}

func (initializer *applicationInitializer) mountSourdoughRecipeScaleAPIRoutes(router chi.Router) {
	router.Post("/{id}/scale", initializer.dependencyManager.SourdoughRecipeScale().Router().Scale())
}

func (initializer *applicationInitializer) getConfig() config.Config {
	return initializer.dependencyManager.Common().ConfigManager().GetConfig()
}

func (initializer *applicationInitializer) mountFlourAPIRoutes(router chi.Router) {
	flourHandler := initializer.dependencyManager.Flour().Router()

	router.
		With(httpin.NewInput(rest.PageInput{})).
		Get("/", flourHandler.Find())
	router.Post("/", flourHandler.Create())
	router.Route("/{id}", func(idRouter chi.Router) {
		idRouter.Get("/", flourHandler.FindById())
	})
	router.
		With(httpin.NewInput(rest.SearchFlourInput{})).
		Get("/search", flourHandler.Search())

}

func NewApplicationInitializer() domain.ApplicationInitializer {
	return &applicationInitializer{
		dependencyManager: dependency.NewDependencyManager(),
	}
}

type application struct {
	dependencyManager domain.DependencyManager
	server            *http.Server
}

func (app *application) Server() *http.Server {
	return app.server
}

func (app *application) Config() config.Config {
	return app.dependencyManager.Common().ConfigManager().GetConfig()
}
