package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"dough-calculator/internal/domain"
)

const (
	FlourDatabase   = "dough-calculator"
	FlourCollection = "flour"
)

type flourRepository struct {
	mongoDBService domain.MongoDBService
}

func (repository *flourRepository) Create(ctx context.Context, flour domain.FlourEntity) (entity domain.FlourEntity, err error) {
	collection, err := repository.getCollection()
	if err != nil {
		return
	}

	result, err := collection.InsertOne(ctx, flour)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to insert flour")
		return entity, errors.Wrap(err, "failed to insert flour")
	}

	log.Debug().Msgf("Inserted a single document: %s", result.InsertedID)

	return flour, nil
}

func (repository *flourRepository) FindById(ctx context.Context, id uuid.UUID) (entity domain.FlourEntity, err error) {
	collection, err := repository.getCollection()
	if err != nil {
		return
	}

	err = collection.
		FindOne(ctx, bson.D{{"_id", id}}).
		Decode(&entity)
	if err != nil {
		log.Error().
			Err(err).
			Stringer("id", id).
			Msg("failed to get flour by id")
		return entity, errors.Wrap(err, "failed to get flour by id")
	}

	return entity, nil
}

func (repository *flourRepository) Find(ctx context.Context, offset, limit int) (result []domain.FlourEntity, err error) {
	defer func() {
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to find all recipes")
		}
	}()

	collection, err := repository.getCollection()
	if err != nil {
		return
	}

	cursor, err := collection.Find(ctx, bson.D{}, options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{"created_at", -1}}))
	if err != nil {
		return nil, errors.Wrap(err, "failed to find recipes")
	}

	if err = cursor.All(ctx, &result); err != nil {
		return nil, errors.Wrap(err, "failed to decode recipes")
	}

	return
}

func (repository *flourRepository) SearchByName(ctx context.Context, name string) (result []domain.FlourEntity, err error) {
	defer func() {
		if err != nil {
			log.Error().
				Err(err).
				Str("name", name).
				Msg("failed to find flour by name")
		}
	}()

	collection, err := repository.getCollection()
	if err != nil {
		return
	}

	cur, err := collection.Find(ctx, bson.D{{
		"name", bson.D{{
			"$regex", primitive.Regex{Pattern: name, Options: "i"},
		}},
	}})

	if errors.Is(err, mongo.ErrNoDocuments) {
		return result, nil
	} else if err != nil {
		return result, errors.Wrap(err, "failed to find flour")
	}

	if err = cur.All(ctx, &result); err != nil {
		return nil, errors.Wrap(err, "failed to decode flour")
	}

	return
}

func (repository *flourRepository) getCollection() (*mongo.Collection, error) {
	collection, err := repository.mongoDBService.GetCollection(FlourDatabase, FlourCollection)
	if err != nil {
		log.Error().
			Err(err).
			Str("database", FlourDatabase).
			Str("collection", FlourCollection).
			Msg("failed to get collection")
		return nil, errors.Wrap(err, "failed to get collection")
	}
	return collection, nil
}

func NewFlourRepository(mongoDBService domain.MongoDBService) (domain.FlourRepository, error) {
	if mongoDBService == nil {
		return nil, errors.New("service cannot be nil")
	}

	collection, err := mongoDBService.GetCollection(FlourDatabase, FlourCollection)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get collection")
	}

	_, err = collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{"name", 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create index")
	}

	return &flourRepository{
		mongoDBService: mongoDBService,
	}, nil
}
