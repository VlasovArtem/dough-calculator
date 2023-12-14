//go:generate mockgen -source=dependency.go -destination=mocks/dependency.go -package mocks

package domain

import "context"

type DependencyInitializer interface {
	Initialize(ctx context.Context) error
}

type DependencyManager interface {
	DependencyInitializer
	Common() CommonDependencyService
	SourdoughRecipe() SourdoughRecipeDependencyService
	SourdoughRecipeScale() SourdoughRecipeScaleDependencyService
	Flour() FlourDependencyService
}

type SourdoughRecipeDependencyService interface {
	DependencyInitializer
	Repository() SourdoughRecipeRepository
	Service() SourdoughRecipeService
	Router() SourdoughRecipeHandler
}

type SourdoughRecipeScaleDependencyService interface {
	DependencyInitializer
	Service() SourdoughRecipeScaleService
	Router() SourdoughRecipeScaleHandler
}

type CommonDependencyService interface {
	DependencyInitializer
	Actuator() ActuatorHandler
	MongoDBService() MongoDBService
	ConfigManager() ConfigManager
}

type FlourDependencyService interface {
	DependencyInitializer
	Repository() FlourRepository
	Service() FlourService
	Router() FlourHandler
}
