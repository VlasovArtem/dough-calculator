package dependency

import (
	"context"

	"github.com/pkg/errors"

	"dough-calculator/internal/config"
	"dough-calculator/internal/controller/rest"
	"dough-calculator/internal/domain"
	"dough-calculator/internal/service"
)

type commonDependencyService struct {
	configManagerCreator func() domain.ConfigManager
	configManager        domain.ConfigManager

	actuatorHandlerCreator func() domain.ActuatorHandler
	actuatorHandler        domain.ActuatorHandler

	mongoDBServiceCreator func(config config.Database) (domain.MongoDBService, error)
	mongoDBService        domain.MongoDBService
}

func (dependencyService *commonDependencyService) Initialize(ctx context.Context) error {
	configManager := dependencyService.configManagerCreator()

	err := configManager.ParseConfig()
	if err != nil {
		return errors.Wrap(err, "failed to parse config")
	}

	mongoDBService, err := dependencyService.mongoDBServiceCreator(configManager.GetConfig().Database)
	if err != nil {
		return errors.Wrap(err, "failed to create mongodb service")
	}

	dependencyService.actuatorHandler = dependencyService.actuatorHandlerCreator()
	dependencyService.configManager = configManager
	dependencyService.mongoDBService = mongoDBService

	return nil
}

func (dependencyService *commonDependencyService) ConfigManager() domain.ConfigManager {
	return dependencyService.configManager
}

func (dependencyService *commonDependencyService) Actuator() domain.ActuatorHandler {
	return dependencyService.actuatorHandler
}

func (dependencyService *commonDependencyService) MongoDBService() domain.MongoDBService {
	return dependencyService.mongoDBService
}

func NewCommonDependencyService() domain.CommonDependencyService {
	return newCommonDependencyService(service.NewConfigManager, rest.NewActuatorHandler, service.NewMongoDBService)
}

func newCommonDependencyService(
	configManagerCreator func() domain.ConfigManager,
	actuatorHandlerCreator func() domain.ActuatorHandler,
	mongoDBServiceCreator func(config config.Database) (domain.MongoDBService, error),
) domain.CommonDependencyService {
	return &commonDependencyService{
		configManagerCreator:   configManagerCreator,
		actuatorHandlerCreator: actuatorHandlerCreator,
		mongoDBServiceCreator:  mongoDBServiceCreator,
	}
}
