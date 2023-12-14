package service

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"dough-calculator/internal/config"
	"dough-calculator/internal/domain"
)

type mongoDBService struct {
	*mongo.Client
	databases   sync.Map
	collections sync.Map
	config      config.Database
}

func (service *mongoDBService) GetDatabase(name string) (*mongo.Database, error) {
	value, ok := service.databases.Load(name)
	if ok {
		return value.(*mongo.Database), nil
	}

	database := service.Database(name)
	if database == nil {
		return nil, errors.New("database cannot be nil")
	}

	service.databases.Store(name, database)

	return database, nil
}

func (service *mongoDBService) GetCollection(databaseName, collectionName string) (*mongo.Collection, error) {
	database, err := service.GetDatabase(databaseName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get database")
	}

	value, ok := service.collections.Load(collectionName)
	if ok {
		return value.(*mongo.Collection), nil
	}

	collection := database.Collection(collectionName)
	if collection == nil {
		return nil, errors.New("collection cannot be nil")
	}

	service.collections.Store(collectionName, collection)

	return collection, nil
}

func (service *mongoDBService) GetClient() (*mongo.Client, error) {
	return service.Client, nil
}

func (service *mongoDBService) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), service.config.ConnectionTimeout)
	defer cancel()
	return service.Client.Disconnect(ctx)
}

func NewMongoDBService(config config.Database) (domain.MongoDBService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Uri))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create database")
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create database")
	}

	return newMongoDBService(config, client)
}

func newMongoDBService(config config.Database, client *mongo.Client) (domain.MongoDBService, error) {
	if client == nil {
		return nil, errors.New("client cannot be nil")
	}

	return &mongoDBService{config: config, Client: client}, nil
}
