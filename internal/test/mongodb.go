package test

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"

	"dough-calculator/internal/config"
)

type MongoDBStarter interface {
	Start() error
	GetConfig() (config.Database, error)
	Stop() error
}

type mongoDBDockerStarter struct {
	dbContainer *mongodb.MongoDBContainer
	config      config.Database
}

func (starter *mongoDBDockerStarter) Start() error {
	ctx := context.Background()

	image := testcontainers.WithImage("mongo:latest")
	container, err := mongodb.RunContainer(ctx, image)
	if err != nil {
		return errors.Wrap(err, "failed to start mongodb container")
	}

	starter.dbContainer = container

	return nil
}

func (starter *mongoDBDockerStarter) GetConfig() (config.Database, error) {
	if starter.dbContainer == nil {
		return config.Database{}, errors.New("container is not started")
	}

	if starter.config == (config.Database{}) {
		var err error
		starter.config.Uri, err = starter.dbContainer.ConnectionString(context.Background())
		if err != nil {
			return config.Database{}, err
		}

		starter.config.ConnectionTimeout = time.Second * 5
	}

	return starter.config, nil
}

func (starter *mongoDBDockerStarter) Stop() error {
	err := starter.dbContainer.Terminate(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to stop mongodb container")
	}

	return nil
}

func NewMongoDBDockerStarter() MongoDBStarter {
	return &mongoDBDockerStarter{}
}
