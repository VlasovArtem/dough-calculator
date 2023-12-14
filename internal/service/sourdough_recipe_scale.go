package service

import (
	"context"
	"math"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"dough-calculator/internal/domain"
)

type scaledKey struct {
	id               uuid.UUID
	finalDoughWeight int
}

type sourdoughRecipeScaleService struct {
	sourdoughRecipeService domain.SourdoughRecipeService
	scaledRecipes          sync.Map
}

func (service *sourdoughRecipeScaleService) Scale(ctx context.Context, id uuid.UUID, request domain.SourdoughRecipeScaleRequestDto) (domain.SourdoughRecipeDto, error) {
	key := scaledKey{
		id:               id,
		finalDoughWeight: request.FinalDoughWeight,
	}

	if scaledRecipe, ok := service.scaledRecipes.Load(key); ok {
		return scaledRecipe.(domain.SourdoughRecipeDto), nil
	}

	recipeDto, err := service.sourdoughRecipeService.FindById(ctx, id)
	if err != nil {
		return domain.SourdoughRecipeDto{}, err
	}

	scaledRecipe := service.scale(recipeDto, request)

	service.scaledRecipes.Store(key, scaledRecipe)

	return scaledRecipe, nil
}

func (service *sourdoughRecipeScaleService) scale(dto domain.SourdoughRecipeDto, request domain.SourdoughRecipeScaleRequestDto) domain.SourdoughRecipeDto {
	scaledDto := dto

	scaledDto.Flour = service.scaleFlourAmounts(dto.Details.TotalWeight, dto.Flour, request.FinalDoughWeight)
	scaledDto.Water = service.scaleBakerAmounts(dto.Details.TotalWeight, dto.Water, request.FinalDoughWeight)
	scaledDto.AdditionalIngredients = service.scaleBakerAmounts(dto.Details.TotalWeight, dto.AdditionalIngredients, request.FinalDoughWeight)
	scaledDto.Levain = service.scaleLevain(dto.Details.TotalWeight, dto.Levain, request.FinalDoughWeight)
	scaledDto.Details = service.scaleRecipeDetails(dto.Details.TotalWeight, dto.Details, request.FinalDoughWeight)
	scaledDto.Yield = domain.RecipeYieldDto{}

	return scaledDto
}

func (service *sourdoughRecipeScaleService) scaleLevain(totalWeight int, levain domain.SourdoughLevainAgentDto, newTotalWeight int) domain.SourdoughLevainAgentDto {
	scaledLevain := levain

	scaledLevain.Starter = service.scaleBakerAmount(totalWeight, levain.Starter, newTotalWeight)
	scaledLevain.Flour = service.scaleFlourAmounts(totalWeight, levain.Flour, newTotalWeight)
	scaledLevain.Water = service.scaleBakerAmount(totalWeight, levain.Water, newTotalWeight)
	scaledLevain.Amount = service.scaleBakerAmount(totalWeight, levain.Amount, newTotalWeight)

	return scaledLevain
}

func (service *sourdoughRecipeScaleService) scaleRecipeDetails(totalWeight int, details domain.RecipeDetailsDto, newTotalWeight int) domain.RecipeDetailsDto {
	scaledDetails := details

	scaledDetails.Flour = service.scaleBakerAmount(totalWeight, details.Flour, newTotalWeight)
	scaledDetails.Water = service.scaleBakerAmount(totalWeight, details.Water, newTotalWeight)
	scaledDetails.Levain = service.scaleBakerAmount(totalWeight, details.Levain, newTotalWeight)
	scaledDetails.AdditionalIngredients = service.scaleBakerAmount(totalWeight, details.AdditionalIngredients, newTotalWeight)
	scaledDetails.TotalWeight = newTotalWeight

	return scaledDetails
}

func (service *sourdoughRecipeScaleService) scaleFlourAmounts(totalWeight int, flours []domain.FlourAmountDto, newTotalWeight int) []domain.FlourAmountDto {
	scaledFlour := make([]domain.FlourAmountDto, len(flours))

	for i, flour := range flours {
		scaledFlour[i] = service.scaleFlourAmount(totalWeight, flour, newTotalWeight)
	}

	return scaledFlour
}

func (service *sourdoughRecipeScaleService) scaleFlourAmount(totalWeight int, flour domain.FlourAmountDto, newTotalWeight int) domain.FlourAmountDto {
	return domain.FlourAmountDto{
		FlourDto: flour.FlourDto,
		Amount:   service.scaleAmount(totalWeight, flour.Amount, newTotalWeight),
	}
}

func (service *sourdoughRecipeScaleService) scaleBakerAmounts(totalWeight int, items []domain.BakerAmountDto, newTotalWeight int) []domain.BakerAmountDto {
	scaledItems := make([]domain.BakerAmountDto, len(items))

	for i, item := range items {
		scaledItems[i] = service.scaleBakerAmount(totalWeight, item, newTotalWeight)
	}

	return scaledItems
}

func (service *sourdoughRecipeScaleService) scaleBakerAmount(totalWeight int, item domain.BakerAmountDto, newTotalWeight int) domain.BakerAmountDto {
	return domain.BakerAmountDto{
		Amount:          service.scaleAmount(totalWeight, item.Amount, newTotalWeight),
		BakerPercentage: item.BakerPercentage,
		Name:            item.Name,
	}
}

func (service *sourdoughRecipeScaleService) scaleAmount(originalTotalWeight int, originalAmount float64, newTotalWeight int) float64 {
	return math.Round(float64(newTotalWeight) * originalAmount / float64(originalTotalWeight))
}

func NewSourdoughRecipeScaleService(sourdoughRecipeService domain.SourdoughRecipeService) (domain.SourdoughRecipeScaleService, error) {
	if sourdoughRecipeService == nil {
		return nil, errors.New("sourdoughRecipeService cannot be nil")
	}

	return &sourdoughRecipeScaleService{
		sourdoughRecipeService: sourdoughRecipeService,
	}, nil
}
