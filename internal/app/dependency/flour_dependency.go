package dependency

import (
	"context"

	"github.com/pkg/errors"

	"dough-calculator/internal/controller/rest"
	"dough-calculator/internal/domain"
	"dough-calculator/internal/repository"
	"dough-calculator/internal/service"
)

type flourDependencyService struct {
	repositoryCreator func(mongoDBService domain.MongoDBService) (domain.FlourRepository, error)
	repository        domain.FlourRepository

	serviceCreator func(repository domain.FlourRepository) (domain.FlourService, error)
	service        domain.FlourService

	handlerCreator func(service domain.FlourService) (domain.FlourHandler, error)
	handler        domain.FlourHandler
}

func (dependencyService *flourDependencyService) Initialize(ctx context.Context) error {
	mongoDBService, err := getFromContext[domain.MongoDBService](ctx, "mongoDBService")
	if err != nil {
		return errors.Wrap(err, "failed to get mongoDBService from context")
	}

	flourRepository, err := dependencyService.repositoryCreator(mongoDBService)
	if err != nil {
		return errors.Wrap(err, "failed to create repository")
	}

	flourService, err := dependencyService.serviceCreator(flourRepository)
	if err != nil {
		return errors.Wrap(err, "failed to create service")
	}

	flourHandler, err := dependencyService.handlerCreator(flourService)
	if err != nil {
		return errors.Wrap(err, "failed to create handler")
	}

	dependencyService.repository = flourRepository
	dependencyService.service = flourService
	dependencyService.handler = flourHandler

	return nil
}

func (dependencyService *flourDependencyService) Repository() domain.FlourRepository {
	return dependencyService.repository
}

func (dependencyService *flourDependencyService) Service() domain.FlourService {
	return dependencyService.service
}

func (dependencyService *flourDependencyService) Router() domain.FlourHandler {
	return dependencyService.handler
}

func NewFlourDependencyService() domain.FlourDependencyService {
	return newFlourDependencyService(repository.NewFlourRepository, service.NewFlourService, rest.NewFlourHandler)
}

func newFlourDependencyService(
	repositoryCreator func(mongoDBService domain.MongoDBService) (domain.FlourRepository, error),
	serviceCreator func(repository domain.FlourRepository) (domain.FlourService, error),
	handlerCreator func(service domain.FlourService) (domain.FlourHandler, error),
) domain.FlourDependencyService {
	return &flourDependencyService{
		repositoryCreator: repositoryCreator,
		serviceCreator:    serviceCreator,
		handlerCreator:    handlerCreator,
	}
}
