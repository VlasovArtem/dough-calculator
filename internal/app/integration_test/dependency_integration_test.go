//go:build integration && docker

package integration_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"

	"dough-calculator/internal/app/dependency"
	"dough-calculator/internal/config"
	"dough-calculator/internal/domain"
	"dough-calculator/internal/test"
)

func TestDependencyTestSuite(t *testing.T) {
	suite.Run(t, &DependencyTestSuite{
		MongoDBDockerIntegrationTestSuite: test.NewMongoDBDockerIntegrationTestSuite(dockerStarter),
	})
}

type DependencyTestSuite struct {
	test.MongoDBDockerIntegrationTestSuite

	configManager domain.ConfigManager

	target domain.DependencyManager
}

func (suite *DependencyTestSuite) SetupSuite() {
	suite.MongoDBDockerIntegrationTestSuite.SetupSuite()

	configPath := suite.createConfigFile()

	err := os.Setenv("CONFIG_PATH", configPath)
	suite.Require().NoError(err)

	dependencyManager := dependency.NewDependencyManager()
	suite.Require().NoError(err)

	suite.target = dependencyManager
}

func (suite *DependencyTestSuite) TestInitialize() {
	err := suite.target.Initialize(context.Background())

	suite.NoError(err)

	suite.NotNil(suite.target.SourdoughRecipe().Router())
	suite.NotNil(suite.target.SourdoughRecipeScale().Router())
}

func (suite *DependencyTestSuite) createConfigFile() string {
	tempDir := suite.T().TempDir()

	configFilePath := tempDir + "/config.yaml"

	databaseConfig := suite.GetConfig()

	cfg := config.Config{
		Application: config.Application{
			Name: "loan",
			Rest: config.Rest{
				Server:               ":0",
				ContextPath:          "/v1",
				ReadTimeout:          10,
				WriteTimeout:         10,
				IdleTimeout:          10,
				GraceShutdownTimeout: 10,
			},
		},
		Database: databaseConfig,
	}

	cfgBytes, err := yaml.Marshal(cfg)
	suite.Require().NoError(err)
	err = os.WriteFile(configFilePath, cfgBytes, 0644)
	suite.Require().NoError(err)

	return configFilePath
}
