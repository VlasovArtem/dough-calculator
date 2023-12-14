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
	SourdoughRecipeDatabase   = "dough-calculator"
	SourdoughRecipeCollection = "sourdough-recipes"
)

type sourdoughRecipeRepository struct {
	mongoDBService domain.MongoDBService
}

func (repository *sourdoughRecipeRepository) Create(ctx context.Context, recipe domain.SourdoughRecipeEntity) (entity domain.SourdoughRecipeEntity, err error) {
	collection, err := repository.getCollection()
	if err != nil {
		return
	}

	result, err := collection.InsertOne(ctx, recipe)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to insert recipe")
		return domain.SourdoughRecipeEntity{}, errors.Wrap(err, "failed to insert sourdough recipe")
	}

	log.Debug().Msgf("Inserted a single document: %s", result.InsertedID)

	return recipe, nil
}

func (repository *sourdoughRecipeRepository) GetById(ctx context.Context, id uuid.UUID) (entity domain.SourdoughRecipeEntity, err error) {
	defer func() {
		if err != nil {
			log.Error().
				Err(err).
				Stringer("id", id).
				Msg("failed to get recipe by id")
		}
	}()

	collection, err := repository.getCollection()
	if err != nil {
		return
	}

	err = collection.
		FindOne(ctx, bson.D{{"_id", id}}).
		Decode(&entity)
	if err != nil {
		return domain.SourdoughRecipeEntity{}, errors.Wrap(err, "failed to find recipe")
	}

	return
}

func (repository *sourdoughRecipeRepository) Find(ctx context.Context, offset, limit int) (result []domain.SourdoughRecipeEntity, err error) {
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

func (repository *sourdoughRecipeRepository) SearchByName(ctx context.Context, name string) (recipes []domain.SourdoughRecipeEntity, err error) {
	defer func() {
		if err != nil {
			log.Error().
				Err(err).
				Str("name", name).
				Msg("failed to find recipe by name")
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
		return recipes, nil
	} else if err != nil {
		return recipes, errors.Wrap(err, "failed to find recipe")
	}

	if err = cur.All(ctx, &recipes); err != nil {
		return nil, errors.Wrap(err, "failed to decode recipes")
	}

	return
}

func (repository *sourdoughRecipeRepository) getCollection() (*mongo.Collection, error) {
	collection, err := repository.mongoDBService.GetCollection(SourdoughRecipeDatabase, SourdoughRecipeCollection)
	if err != nil {
		log.Error().
			Err(err).
			Str("database", SourdoughRecipeDatabase).
			Str("collection", SourdoughRecipeCollection).
			Msg("failed to get collection")
		return nil, errors.Wrap(err, "failed to get collection")
	}
	return collection, nil
}

func NewSourdoughRecipeRepository(service domain.MongoDBService) (domain.SourdoughRecipeRepository, error) {
	if service == nil {
		return nil, errors.New("service cannot be nil")
	}

	collection, err := service.GetCollection(SourdoughRecipeDatabase, SourdoughRecipeCollection)
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

	return &sourdoughRecipeRepository{mongoDBService: service}, nil
}
