//go:generate mockgen -source=sourdough_recipe.go -destination=mocks/sourdough_recipe.go -package mocks

package domain

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"dough-calculator/internal/utils"
)

type SourdoughRecipeEntity struct {
	RecipeEntity `bson:",inline"`
	Levain       SourdoughLevainAgent
}

func (entity SourdoughRecipeEntity) ToDto() SourdoughRecipeDto {
	return SourdoughRecipeDto{
		RecipeDto: entity.RecipeEntity.ToDto(),
		Levain:    entity.Levain.ToDto(),
	}
}

type SourdoughLevainAgent struct {
	Amount  BakerAmount
	Starter BakerAmount
	Flour   []FlourAmount
	Water   BakerAmount
}

func (agent SourdoughLevainAgent) ToDto() SourdoughLevainAgentDto {
	return SourdoughLevainAgentDto{
		Amount:  agent.Amount.ToDto(),
		Starter: agent.Starter.ToDto(),
		Flour:   utils.Map(agent.Flour, func(f FlourAmount) FlourAmountDto { return f.ToDto() }),
		Water:   agent.Water.ToDto(),
	}
}

type SourdoughRecipeRepository interface {
	Create(ctx context.Context, recipe SourdoughRecipeEntity) (SourdoughRecipeEntity, error)
	GetById(ctx context.Context, id uuid.UUID) (SourdoughRecipeEntity, error)
	Find(ctx context.Context, offset, limit int) ([]SourdoughRecipeEntity, error)
	SearchByName(ctx context.Context, name string) ([]SourdoughRecipeEntity, error)
}

type SourdoughLevainAgentDto struct {
	Amount  BakerAmountDto   `json:"amount"`
	Starter BakerAmountDto   `json:"starter"`
	Flour   []FlourAmountDto `json:"flour"`
	Water   BakerAmountDto   `json:"water"`
}

func (dto SourdoughLevainAgentDto) ToEntity() SourdoughLevainAgent {
	return SourdoughLevainAgent{
		Amount:  dto.Amount.ToEntity(),
		Starter: dto.Starter.ToEntity(),
		Flour:   utils.Map(dto.Flour, func(f FlourAmountDto) FlourAmount { return f.ToEntity() }),
		Water:   dto.Water.ToEntity(),
	}
}

type SourdoughRecipeDto struct {
	RecipeDto
	Levain SourdoughLevainAgentDto `json:"levain"`
}

func (dto SourdoughRecipeDto) ToEntity() SourdoughRecipeEntity {
	return SourdoughRecipeEntity{
		RecipeEntity: dto.RecipeDto.ToEntity(),
		Levain:       dto.Levain.ToEntity(),
	}
}

type SourdoughRecipeService interface {
	Create(ctx context.Context, request CreateSourdoughRecipeRequest) (SourdoughRecipeDto, error)
	FindById(ctx context.Context, id uuid.UUID) (SourdoughRecipeDto, error)
	Find(ctx context.Context, offset, limit int) ([]SourdoughRecipeDto, error)
	SearchByName(ctx context.Context, name string) ([]SourdoughRecipeDto, error)
}

type SourdoughRecipeScaleService interface {
	Scale(ctx context.Context, id uuid.UUID, request SourdoughRecipeScaleRequestDto) (SourdoughRecipeDto, error)
}

type CreateSourdoughRecipeRequest struct {
	Name                  string                       `json:"name"`
	Description           string                       `json:"description"`
	Flour                 []FlourAmountDto             `json:"flour"`
	Water                 []BakerAmountDto             `json:"water"`
	Levain                SourdoughLevainAgentDto      `json:"levain"`
	AdditionalIngredients []BakerAmountDto             `json:"additional_ingredients"`
	NutritionFacts        map[string]NutritionFactsDto `json:"nutrition_facts"`
	Yield                 RecipeYieldDto               `json:"yield"`
}

type SourdoughRecipeScaleRequestDto struct {
	FinalDoughWeight int `json:"final_dough_weight"`
}

type SourdoughRecipeHandler interface {
	Create() http.HandlerFunc
	FindById() http.HandlerFunc
	Find() http.HandlerFunc
	Search() http.HandlerFunc
}

type SourdoughRecipeScaleHandler interface {
	Scale() http.HandlerFunc
}
