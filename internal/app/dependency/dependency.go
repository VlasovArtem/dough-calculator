package dependency

import (
	"context"

	"github.com/pkg/errors"

	"dough-calculator/internal/domain"
)

type dependencyManager struct {
	commonDependencyService               domain.CommonDependencyService
	sourdoughRecipeDependencyService      domain.SourdoughRecipeDependencyService
	sourdoughRecipeScaleDependencyService domain.SourdoughRecipeScaleDependencyService
	flourDependencyService                domain.FlourDependencyService
}

func (manager *dependencyManager) Initialize(ctx context.Context) error {
	err := manager.commonDependencyService.Initialize(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to initialize common dependency service")
	}

	ctx = context.WithValue(ctx, "configManager", manager.commonDependencyService.ConfigManager())
	ctx = context.WithValue(ctx, "mongoDBService", manager.commonDependencyService.MongoDBService())

	err = manager.sourdoughRecipeDependencyService.Initialize(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to initialize sourdough recipe dependency service")
	}

	ctx = context.WithValue(ctx, "sourdoughRecipeService", manager.sourdoughRecipeDependencyService.Service())

	err = manager.sourdoughRecipeScaleDependencyService.Initialize(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to initialize sourdough recipe scale dependency service")
	}

	err = manager.flourDependencyService.Initialize(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to initialize flour dependency service")
	}

	return nil
}

func (manager *dependencyManager) SourdoughRecipe() domain.SourdoughRecipeDependencyService {
	return manager.sourdoughRecipeDependencyService
}

func (manager *dependencyManager) SourdoughRecipeScale() domain.SourdoughRecipeScaleDependencyService {
	return manager.sourdoughRecipeScaleDependencyService
}

func (manager *dependencyManager) Common() domain.CommonDependencyService {
	return manager.commonDependencyService
}

func (manager *dependencyManager) Flour() domain.FlourDependencyService {
	return manager.flourDependencyService
}

func NewDependencyManager() domain.DependencyManager {
	return newDependencyManager(
		NewCommonDependencyService(),
		NewSourdoughRecipeDependencyService(),
		NewSourdoughRecipeScaleDependencyService(),
		NewFlourDependencyService(),
	)
}

func newDependencyManager(
	commonDependencyService domain.CommonDependencyService,
	sourdoughRecipeDependencyService domain.SourdoughRecipeDependencyService,
	sourdoughRecipeScaleDependencyService domain.SourdoughRecipeScaleDependencyService,
	flourDependencyService domain.FlourDependencyService,
) domain.DependencyManager {
	return &dependencyManager{
		commonDependencyService:               commonDependencyService,
		sourdoughRecipeDependencyService:      sourdoughRecipeDependencyService,
		sourdoughRecipeScaleDependencyService: sourdoughRecipeScaleDependencyService,
		flourDependencyService:                flourDependencyService,
	}
}

func getFromContext[T any](ctx context.Context, key string) (t T, err error) {
	if ctx == nil {
		return t, errors.New("context is nil")
	}
	value := ctx.Value(key)
	if isNil(value) {
		return t, errors.Errorf("%s is nil", key)
	}

	t, ok := value.(T)
	if !ok {
		return t, errors.Errorf("%s is not valid type", key)
	}

	return t, nil
}

func isNil[T any](t T) bool {
	switch casted := any(t).(type) {
	case interface{}:
		return casted == nil
	case nil:
		return true
	default:
		return false
	}
}
