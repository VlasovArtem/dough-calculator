//go:generate mockgen -destination=./mocks/recipe.go -package=mocks -source=recipe.go

package domain

import (
	"time"

	"github.com/google/uuid"

	"dough-calculator/internal/utils"
)

type RecipeDetails struct {
	Flour                 BakerAmount
	Water                 BakerAmount
	Levain                BakerAmount
	AdditionalIngredients BakerAmount `bson:"additional_ingredients"`
	TotalWeight           int         `bson:"total_weight"`
}

func (details RecipeDetails) ToDto() RecipeDetailsDto {
	return RecipeDetailsDto{
		Flour:                 details.Flour.ToDto(),
		Water:                 details.Water.ToDto(),
		Levain:                details.Levain.ToDto(),
		AdditionalIngredients: details.AdditionalIngredients.ToDto(),
		TotalWeight:           details.TotalWeight,
	}
}

type RecipeEntity struct {
	Id                    uuid.UUID `bson:"_id"`
	Name                  string
	Description           string
	Flour                 []FlourAmount
	Water                 []BakerAmount
	AdditionalIngredients []BakerAmount             `bson:"additional_ingredients"`
	Details               RecipeDetails             `bson:"recipe_details"`
	NutritionFacts        map[string]NutritionFacts `bson:"nutrition_facts"`
	CreatedAt             time.Time                 `bson:"created_at"`
	UpdatedAt             *time.Time                `bson:"updated_at,omitempty"`
	Yield                 RecipeYield
}

func (entity RecipeEntity) ToDto() RecipeDto {
	toBakerAmountDto := func(amount BakerAmount) BakerAmountDto { return amount.ToDto() }

	nutritionFactsDtoMap := make(map[string]NutritionFactsDto, len(entity.NutritionFacts))
	for key, value := range entity.NutritionFacts {
		nutritionFactsDtoMap[key] = value.ToDto()
	}

	return RecipeDto{
		Id:                    entity.Id,
		Name:                  entity.Name,
		Description:           entity.Description,
		Flour:                 utils.Map(entity.Flour, func(amount FlourAmount) FlourAmountDto { return amount.ToDto() }),
		Water:                 utils.Map(entity.Water, toBakerAmountDto),
		AdditionalIngredients: utils.Map(entity.AdditionalIngredients, toBakerAmountDto),
		Details:               entity.Details.ToDto(),
		NutritionFacts:        nutritionFactsDtoMap,
		CreatedAt:             entity.CreatedAt,
		UpdatedAt:             entity.UpdatedAt,
		Yield:                 entity.Yield.ToDto(),
	}
}

type RecipeYield struct {
	Unit   string
	Amount int
}

func (yield RecipeYield) ToDto() RecipeYieldDto {
	return RecipeYieldDto{
		Unit:   yield.Unit,
		Amount: yield.Amount,
	}
}

type FlourAmount struct {
	FlourEntity
	Amount float64
}

func (flourAmount FlourAmount) ToDto() FlourAmountDto {
	return FlourAmountDto{
		FlourDto: flourAmount.FlourEntity.ToDto(),
		Amount:   flourAmount.Amount,
	}
}

type BakerAmount struct {
	Amount          float64
	BakerPercentage float64
	Name            string
}

func (bakerAmount BakerAmount) ToDto() BakerAmountDto {
	return BakerAmountDto{
		Amount:          bakerAmount.Amount,
		BakerPercentage: bakerAmount.BakerPercentage,
		Name:            bakerAmount.Name,
	}
}

type RecipeDetailsDto struct {
	Flour                 BakerAmountDto `json:"flour"`
	Water                 BakerAmountDto `json:"water"`
	Levain                BakerAmountDto `json:"levain"`
	AdditionalIngredients BakerAmountDto `json:"additional_ingredients"`
	TotalWeight           int            `json:"total_weight"`
}

func (dto RecipeDetailsDto) ToEntity() RecipeDetails {
	return RecipeDetails{
		Flour:                 dto.Flour.ToEntity(),
		Water:                 dto.Water.ToEntity(),
		Levain:                dto.Levain.ToEntity(),
		AdditionalIngredients: dto.AdditionalIngredients.ToEntity(),
		TotalWeight:           dto.TotalWeight,
	}
}

type RecipeYieldDto struct {
	Unit   string `json:"unit"`
	Amount int    `json:"amount"`
}

func (yield RecipeYieldDto) ToEntity() RecipeYield {
	return RecipeYield{
		Unit:   yield.Unit,
		Amount: yield.Amount,
	}
}

type RecipeDto struct {
	Id                    uuid.UUID                    `json:"id"`
	Name                  string                       `json:"name"`
	Description           string                       `json:"description"`
	Flour                 []FlourAmountDto             `json:"flour"`
	Water                 []BakerAmountDto             `json:"water"`
	AdditionalIngredients []BakerAmountDto             `json:"additional_ingredients"`
	Details               RecipeDetailsDto             `json:"recipe_details"`
	NutritionFacts        map[string]NutritionFactsDto `json:"nutrition_facts"`
	CreatedAt             time.Time                    `json:"created_at"`
	UpdatedAt             *time.Time                   `json:"updated_at,omitempty"`
	Yield                 RecipeYieldDto               `json:"yield"`
}

func (dto RecipeDto) ToEntity() RecipeEntity {
	toFlourAmount := func(amount FlourAmountDto) FlourAmount { return amount.ToEntity() }
	toBakerAmount := func(amount BakerAmountDto) BakerAmount { return amount.ToEntity() }

	nutritionFactsMap := make(map[string]NutritionFacts, len(dto.NutritionFacts))
	for key, value := range dto.NutritionFacts {
		nutritionFactsMap[key] = value.ToEntity()
	}

	return RecipeEntity{
		Id:                    dto.Id,
		Name:                  dto.Name,
		Description:           dto.Description,
		Flour:                 utils.Map(dto.Flour, toFlourAmount),
		Water:                 utils.Map(dto.Water, toBakerAmount),
		AdditionalIngredients: utils.Map(dto.AdditionalIngredients, toBakerAmount),
		Details:               dto.Details.ToEntity(),
		NutritionFacts:        nutritionFactsMap,
		CreatedAt:             dto.CreatedAt,
		UpdatedAt:             dto.UpdatedAt,
		Yield:                 dto.Yield.ToEntity(),
	}
}

type FlourAmountDto struct {
	FlourDto
	Amount float64 `json:"amount"`
}

func (dto FlourAmountDto) ToEntity() FlourAmount {
	return FlourAmount{
		FlourEntity: dto.FlourDto.ToEntity(),
		Amount:      dto.Amount,
	}
}

type BakerAmountDto struct {
	Amount          float64 `json:"amount"`
	BakerPercentage float64 `json:"baker_percentage,omitempty"`
	Name            string  `json:"name,omitempty"`
}

func (dto BakerAmountDto) ToEntity() BakerAmount {
	return BakerAmount{
		Amount:          dto.Amount,
		BakerPercentage: dto.BakerPercentage,
		Name:            dto.Name,
	}
}
