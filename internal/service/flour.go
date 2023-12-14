package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	"dough-calculator/internal/domain"
	internalErrors "dough-calculator/internal/errors"
)

type flourService struct {
	repository domain.FlourRepository
}

func (service *flourService) Create(ctx context.Context, request domain.CreateFlourRequest) (domain.FlourDto, error) {
	createdEntity, err := service.repository.Create(ctx, service.toEntity(request))

	if err != nil {
		log.Err(err).
			Str("name", request.Name).
			Msg("failed to create flour")

		return domain.FlourDto{}, internalErrors.NewInternalServerErrorWrap(err, "failed to create flour")
	}

	return createdEntity.ToDto(), nil
}

func (service *flourService) FindById(ctx context.Context, id uuid.UUID) (domain.FlourDto, error) {
	flourEntity, err := service.repository.FindById(ctx, id)
	if err != nil {
		log.Err(err).
			Str("id", id.String()).
			Msg("failed to find flour by id")

		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.FlourDto{},
				internalErrors.FlourByIdNotFound(id)
		}

		return domain.FlourDto{}, internalErrors.NewInternalServerErrorWrap(err, "failed to find flour by id")
	}
	return flourEntity.ToDto(), nil
}

func (service *flourService) Find(ctx context.Context, offset, limit int) ([]domain.FlourDto, error) {
	flourEntities, err := service.repository.Find(ctx, offset, limit)
	if err != nil {
		log.Err(err).
			Msg("failed to find flours")

		return nil, internalErrors.NewInternalServerErrorWrap(err, "failed to find flours")
	}

	flours := make([]domain.FlourDto, len(flourEntities))
	for i, flourEntity := range flourEntities {
		flours[i] = flourEntity.ToDto()
	}

	return flours, nil
}

func (service *flourService) SearchByName(ctx context.Context, name string) ([]domain.FlourDto, error) {
	flourEntities, err := service.repository.SearchByName(ctx, name)
	if err != nil {
		log.Err(err).
			Str("name", name).
			Msg("failed to search flours by name")

		return nil, internalErrors.NewInternalServerErrorWrap(err, "failed to search flours by name")
	}

	flours := make([]domain.FlourDto, len(flourEntities))
	for i, flourEntity := range flourEntities {
		flours[i] = flourEntity.ToDto()
	}

	return flours, nil
}

func (service *flourService) toEntity(request domain.CreateFlourRequest) domain.FlourEntity {
	return domain.FlourEntity{
		Id:             uuid.New(),
		FlourType:      request.FlourType,
		Name:           request.Name,
		Description:    request.Description,
		NutritionFacts: request.NutritionFacts.ToEntity(),
	}
}

func NewFlourService(repository domain.FlourRepository) (domain.FlourService, error) {
	if repository == nil {
		return nil, errors.New("repository is nil")
	}

	return &flourService{repository: repository}, nil
}
