//go:generate mockgen -destination=./mocks/mongo_db.go -package=mocks -source=mongo_db.go

package domain

import "go.mongodb.org/mongo-driver/mongo"

type MongoDBService interface {
	GetClient() (*mongo.Client, error)
	GetDatabase(name string) (*mongo.Database, error)
	GetCollection(databaseName, collectionName string) (*mongo.Collection, error)
	Disconnect() error
}
