package dependency

import (
	"context"

	"github.com/pkg/errors"

	"dough-calculator/internal/controller/rest"
	"dough-calculator/internal/domain"
	"dough-calculator/internal/repository"
	"dough-calculator/internal/service"
)

type sourdoughRecipeDependencyService struct {
	repositoryCreator func(mongoDBService domain.MongoDBService) (domain.SourdoughRecipeRepository, error)
	repository        domain.SourdoughRecipeRepository

	serviceCreator func(repository domain.SourdoughRecipeRepository) (domain.SourdoughRecipeService, error)
	service        domain.SourdoughRecipeService

	handlerCreator func(service domain.SourdoughRecipeService) (domain.SourdoughRecipeHandler, error)
	handler        domain.SourdoughRecipeHandler
}

func (dependencyService *sourdoughRecipeDependencyService) Initialize(ctx context.Context) error {
	mongoDBService, err := getFromContext[domain.MongoDBService](ctx, "mongoDBService")
	if err != nil {
		return errors.Wrap(err, "failed to get mongoDBService from context")
	}

	sourdoughRecipeRepository, err := dependencyService.repositoryCreator(mongoDBService)
	if err != nil {
		return errors.Wrap(err, "failed to create repository")
	}

	sourdoughRecipeService, err := dependencyService.serviceCreator(sourdoughRecipeRepository)
	if err != nil {
		return errors.Wrap(err, "failed to create service")
	}

	sourdoughRecipeHandler, err := dependencyService.handlerCreator(sourdoughRecipeService)
	if err != nil {
		return errors.Wrap(err, "failed to create handler")
	}

	dependencyService.repository = sourdoughRecipeRepository
	dependencyService.service = sourdoughRecipeService
	dependencyService.handler = sourdoughRecipeHandler

	return nil
}

func (dependencyService *sourdoughRecipeDependencyService) Repository() domain.SourdoughRecipeRepository {
	return dependencyService.repository
}

func (dependencyService *sourdoughRecipeDependencyService) Service() domain.SourdoughRecipeService {
	return dependencyService.service
}

func (dependencyService *sourdoughRecipeDependencyService) Router() domain.SourdoughRecipeHandler {
	return dependencyService.handler
}

func NewSourdoughRecipeDependencyService() domain.SourdoughRecipeDependencyService {
	return newSourdoughRecipeDependencyService(repository.NewSourdoughRecipeRepository, service.NewSourdoughRecipeService, rest.NewSourdoughRecipeHandler)
}

func newSourdoughRecipeDependencyService(
	repositoryCreator func(mongoDBService domain.MongoDBService) (domain.SourdoughRecipeRepository, error),
	serviceCreator func(repository domain.SourdoughRecipeRepository) (domain.SourdoughRecipeService, error),
	handlerCreator func(service domain.SourdoughRecipeService) (domain.SourdoughRecipeHandler, error),
) domain.SourdoughRecipeDependencyService {
	return &sourdoughRecipeDependencyService{
		repositoryCreator: repositoryCreator,
		serviceCreator:    serviceCreator,
		handlerCreator:    handlerCreator,
	}
}
