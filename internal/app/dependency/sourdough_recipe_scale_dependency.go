package dependency

import (
	"context"

	"github.com/pkg/errors"

	"dough-calculator/internal/controller/rest"
	"dough-calculator/internal/domain"
	"dough-calculator/internal/service"
)

type sourdoughRecipeScaleDependencyService struct {
	serviceCreator func(repository domain.SourdoughRecipeService) (domain.SourdoughRecipeScaleService, error)
	service        domain.SourdoughRecipeScaleService

	handlerCreator func(service domain.SourdoughRecipeScaleService) (domain.SourdoughRecipeScaleHandler, error)
	handler        domain.SourdoughRecipeScaleHandler
}

func (dependencyService *sourdoughRecipeScaleDependencyService) Initialize(ctx context.Context) error {
	sourdoughRecipeService, err := getFromContext[domain.SourdoughRecipeService](ctx, "sourdoughRecipeService")
	if err != nil {
		return errors.Wrap(err, "failed to get sourdoughRecipeService from context")
	}

	sourdoughRecipeScaleService, err := dependencyService.serviceCreator(sourdoughRecipeService)
	if err != nil {
		return errors.Wrap(err, "failed to create service")
	}

	sourdoughRecipeScaleHandler, err := dependencyService.handlerCreator(sourdoughRecipeScaleService)
	if err != nil {
		return errors.Wrap(err, "failed to create handler")
	}

	dependencyService.service = sourdoughRecipeScaleService
	dependencyService.handler = sourdoughRecipeScaleHandler

	return nil
}

func (dependencyService *sourdoughRecipeScaleDependencyService) Service() domain.SourdoughRecipeScaleService {
	return dependencyService.service
}

func (dependencyService *sourdoughRecipeScaleDependencyService) Router() domain.SourdoughRecipeScaleHandler {
	return dependencyService.handler
}

func NewSourdoughRecipeScaleDependencyService() domain.SourdoughRecipeScaleDependencyService {
	return newSourdoughRecipeScaleDependencyService(service.NewSourdoughRecipeScaleService, rest.NewSourdoughRecipeScaleHandler)
}

func newSourdoughRecipeScaleDependencyService(
	serviceCreator func(repository domain.SourdoughRecipeService) (domain.SourdoughRecipeScaleService, error),
	handlerCreator func(service domain.SourdoughRecipeScaleService) (domain.SourdoughRecipeScaleHandler, error),
) domain.SourdoughRecipeScaleDependencyService {
	return &sourdoughRecipeScaleDependencyService{
		serviceCreator: serviceCreator,
		handlerCreator: handlerCreator,
	}
}
