//go:build integration && docker

package integration_test

import (
	"os"
	"testing"

	"github.com/rs/zerolog/log"

	"dough-calculator/internal/test"
)

var dockerStarter test.MongoDBStarter

func TestMain(m *testing.M) {
	dockerStarter = test.NewMongoDBDockerStarter()
	err := dockerStarter.Start()

	if err != nil {
		log.Fatal().Err(err).Msg("failed to start mongodb")
	}

	code := m.Run()

	if err = dockerStarter.Stop(); err != nil {
		log.Fatal().Err(err).Msg("failed to stop mongodb")
	}

	os.Exit(code)
}
