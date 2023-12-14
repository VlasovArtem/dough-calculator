package test

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"dough-calculator/internal/domain"
)

type MongoServiceStub struct {
	client *mongo.Client
}

func (serviceStub *MongoServiceStub) GetClient() (*mongo.Client, error) {
	return serviceStub.client, nil
}

func (serviceStub *MongoServiceStub) GetDatabase(name string) (*mongo.Database, error) {
	database := serviceStub.client.Database(name)

	return database, nil
}

func (serviceStub *MongoServiceStub) GetCollection(databaseName, collectionName string) (*mongo.Collection, error) {
	database, err := serviceStub.GetDatabase(databaseName)
	if err != nil {
		return nil, err
	}

	collection := database.Collection(collectionName)

	return collection, nil
}

func (serviceStub *MongoServiceStub) MustGetCollection(databaseName, collectionName string) *mongo.Collection {
	return Must(func() (*mongo.Collection, error) {
		return serviceStub.GetCollection(databaseName, collectionName)
	})
}

func (serviceStub *MongoServiceStub) Disconnect() error {
	return serviceStub.client.Disconnect(context.Background())
}

type MongoDBServiceDockerIntegrationTestSuite struct {
	MongoDBDockerIntegrationTestSuite

	Stub domain.MongoDBService
}

func (suite *MongoDBServiceDockerIntegrationTestSuite) SetupSuite() {
	suite.ApplicationTestSuite.SetupSuite()

	config := suite.GetConfig()

	timeout, cancelFunc := context.WithTimeout(context.Background(), config.ConnectionTimeout)
	defer cancelFunc()

	client, err := mongo.Connect(timeout, options.Client().ApplyURI(config.Uri))
	suite.Require().NoError(err)

	err = client.Ping(timeout, nil)
	suite.Require().NoError(err)

	suite.Stub = &MongoServiceStub{
		client: client,
	}
}

func (suite *MongoDBServiceDockerIntegrationTestSuite) Drop(database, collection string) error {
	getCollection, err := suite.Stub.GetCollection(database, collection)
	if err != nil {
		return err
	}
	return getCollection.Drop(context.Background())
}

func (suite *MongoDBServiceDockerIntegrationTestSuite) MStub() *MongoServiceStub {
	return suite.Stub.(*MongoServiceStub)
}

func NewMongoDBServiceDockerIntegrationTestSuite(starter MongoDBStarter) MongoDBServiceDockerIntegrationTestSuite {
	return MongoDBServiceDockerIntegrationTestSuite{
		MongoDBDockerIntegrationTestSuite: NewMongoDBDockerIntegrationTestSuite(starter),
	}
}
