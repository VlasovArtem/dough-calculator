//go:generate mockgen -destination=./mocks/flour.go -package=mocks -source=flour.go

package domain

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type FlourEntity struct {
	Id             uuid.UUID `bson:"_id"`
	FlourType      string
	Name           string
	Description    string
	NutritionFacts NutritionFacts
}

func (entity FlourEntity) ToDto() FlourDto {
	return FlourDto{
		Id:             entity.Id,
		FlourType:      entity.FlourType,
		Name:           entity.Name,
		Description:    entity.Description,
		NutritionFacts: entity.NutritionFacts.ToDto(),
	}
}

type FlourRepository interface {
	Create(ctx context.Context, flour FlourEntity) (FlourEntity, error)
	FindById(ctx context.Context, id uuid.UUID) (FlourEntity, error)
	Find(ctx context.Context, offset, limit int) ([]FlourEntity, error)
	SearchByName(ctx context.Context, name string) ([]FlourEntity, error)
}

type FlourDto struct {
	Id             uuid.UUID         `json:"id"`
	FlourType      string            `json:"flour_type"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	NutritionFacts NutritionFactsDto `json:"nutrition_facts"`
}

func (dto FlourDto) ToEntity() FlourEntity {
	return FlourEntity{
		Id:             dto.Id,
		FlourType:      dto.FlourType,
		Name:           dto.Name,
		Description:    dto.Description,
		NutritionFacts: dto.NutritionFacts.ToEntity(),
	}
}

type FlourService interface {
	Create(ctx context.Context, request CreateFlourRequest) (FlourDto, error)
	FindById(ctx context.Context, id uuid.UUID) (FlourDto, error)
	Find(ctx context.Context, offset, limit int) ([]FlourDto, error)
	SearchByName(ctx context.Context, name string) ([]FlourDto, error)
}

type CreateFlourRequest struct {
	FlourType      string            `json:"flour_type"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	NutritionFacts NutritionFactsDto `json:"nutrition_facts"`
}

type FlourHandler interface {
	Create() http.HandlerFunc
	FindById() http.HandlerFunc
	Find() http.HandlerFunc
	Search() http.HandlerFunc
}
