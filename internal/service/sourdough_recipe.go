package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	"dough-calculator/internal/domain"
	internalErrors "dough-calculator/internal/errors"
	"dough-calculator/internal/utils"
)

type sourdoughRecipeService struct {
	repository domain.SourdoughRecipeRepository
}

func (service *sourdoughRecipeService) Create(ctx context.Context, request domain.CreateSourdoughRecipeRequest) (domain.SourdoughRecipeDto, error) {
	recipe := service.toNewRecipe(request)

	createdEntity, err := service.repository.Create(ctx, recipe)

	if err != nil {
		log.Err(err).
			Str("name", recipe.Name).
			Msg("failed to create recipe")

		return domain.SourdoughRecipeDto{}, internalErrors.NewInternalServerErrorWrap(err, "failed to create recipe")
	}

	return createdEntity.ToDto(), nil
}

func (service *sourdoughRecipeService) FindById(ctx context.Context, id uuid.UUID) (domain.SourdoughRecipeDto, error) {
	recipe, err := service.repository.GetById(ctx, id)
	if err != nil {
		log.Err(err).
			Str("id", id.String()).
			Msg("failed to find recipe by id")

		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.SourdoughRecipeDto{},
				internalErrors.SourdoughRecipeNotFound(fmt.Sprintf("recipe with id %s not found", id.String()))
		}

		return domain.SourdoughRecipeDto{}, internalErrors.NewInternalServerErrorWrap(err, "failed to find recipe by id")
	}
	return recipe.ToDto(), nil
}

func (service *sourdoughRecipeService) Find(ctx context.Context, offset, limit int) ([]domain.SourdoughRecipeDto, error) {
	recipes, err := service.repository.Find(ctx, offset, limit)
	if err != nil {
		log.Err(err).
			Msg("failed to find recipes")

		return nil, internalErrors.NewInternalServerErrorWrap(err, "failed to find recipes")
	}

	return utils.Map(recipes, func(entity domain.SourdoughRecipeEntity) domain.SourdoughRecipeDto {
		return entity.ToDto()
	}), nil
}

func (service *sourdoughRecipeService) SearchByName(ctx context.Context, name string) ([]domain.SourdoughRecipeDto, error) {
	recipes, err := service.repository.SearchByName(ctx, name)
	if err != nil {
		log.Err(err).
			Str("name", name).
			Msg("failed to search recipes by name")

		return nil, internalErrors.NewInternalServerErrorWrap(err, "failed to search recipes by name")
	}

	return utils.Map(recipes, func(entity domain.SourdoughRecipeEntity) domain.SourdoughRecipeDto {
		return entity.ToDto()
	}), nil
}

func (service *sourdoughRecipeService) toNewRecipe(request domain.CreateSourdoughRecipeRequest) domain.SourdoughRecipeEntity {
	bakerAmountConverter := func(amount domain.BakerAmountDto) domain.BakerAmount {
		return amount.ToEntity()
	}

	nutritionFacts := make(map[string]domain.NutritionFacts, len(request.NutritionFacts))
	for key, value := range request.NutritionFacts {
		nutritionFacts[key] = value.ToEntity()
	}

	return domain.SourdoughRecipeEntity{
		RecipeEntity: domain.RecipeEntity{
			Id:                    uuid.New(),
			Name:                  request.Name,
			Description:           request.Description,
			Flour:                 utils.Map(request.Flour, func(amount domain.FlourAmountDto) domain.FlourAmount { return amount.ToEntity() }),
			Water:                 utils.Map(request.Water, bakerAmountConverter),
			AdditionalIngredients: utils.Map(request.AdditionalIngredients, bakerAmountConverter),
			Details:               service.calculateRecipeDetails(request).ToEntity(),
			NutritionFacts:        nutritionFacts,
			CreatedAt:             time.Now(),
			Yield:                 request.Yield.ToEntity(),
		},
		Levain: request.Levain.ToEntity(),
	}
}

func (service *sourdoughRecipeService) calculateRecipeDetails(request domain.CreateSourdoughRecipeRequest) domain.RecipeDetailsDto {
	flourAmount := service.calculateFlourAmount(request.Flour)
	waterAmount := service.calculateWaterAmount(flourAmount, request.Water)
	levainAmount := request.Levain.Amount
	additionalIngredientsAmount := service.calculateAdditionalIngredientsAmount(flourAmount, request.AdditionalIngredients)

	recipeDetails := domain.RecipeDetailsDto{
		Flour:                 flourAmount,
		Water:                 waterAmount,
		Levain:                levainAmount,
		AdditionalIngredients: additionalIngredientsAmount,
		TotalWeight:           service.calculateTotalWeight(flourAmount, waterAmount, levainAmount, additionalIngredientsAmount),
	}

	return recipeDetails

}

func (service *sourdoughRecipeService) calculateFlourAmount(flour []domain.FlourAmountDto) domain.BakerAmountDto {
	if len(flour) == 0 {
		return domain.BakerAmountDto{}
	}

	var totalFlourAmount float64
	for _, flourAmount := range flour {
		totalFlourAmount += flourAmount.Amount
	}

	return domain.BakerAmountDto{
		Amount:          totalFlourAmount,
		BakerPercentage: 100,
	}
}

func (service *sourdoughRecipeService) calculateWaterAmount(flour domain.BakerAmountDto, water []domain.BakerAmountDto) domain.BakerAmountDto {
	if len(water) == 0 {
		return domain.BakerAmountDto{}
	}

	var totalWaterAmount float64
	for _, waterAmount := range water {
		totalWaterAmount += waterAmount.Amount
	}

	return domain.BakerAmountDto{
		Amount:          totalWaterAmount,
		BakerPercentage: totalWaterAmount / flour.Amount * 100,
	}
}

func (service *sourdoughRecipeService) calculateAdditionalIngredientsAmount(flour domain.BakerAmountDto, ingredients []domain.BakerAmountDto) domain.BakerAmountDto {
	if len(ingredients) == 0 {
		return domain.BakerAmountDto{}
	}

	var totalIngredientsAmount float64
	for _, ingredient := range ingredients {
		totalIngredientsAmount += ingredient.Amount
	}

	return domain.BakerAmountDto{
		Amount:          totalIngredientsAmount,
		BakerPercentage: totalIngredientsAmount / flour.Amount * 100,
	}
}

func (service *sourdoughRecipeService) calculateTotalWeight(
	flourAmount domain.BakerAmountDto,
	waterAmount domain.BakerAmountDto,
	levainAmount domain.BakerAmountDto,
	additionalIngredientsAmount domain.BakerAmountDto,
) int {
	return int(flourAmount.Amount + waterAmount.Amount + levainAmount.Amount + additionalIngredientsAmount.Amount)
}

func NewSourdoughRecipeService(repository domain.SourdoughRecipeRepository) (domain.SourdoughRecipeService, error) {
	if repository == nil {
		return nil, errors.New("repository cannot be nil")
	}

	return &sourdoughRecipeService{
		repository: repository,
	}, nil
}
