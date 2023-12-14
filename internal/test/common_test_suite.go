package test

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type ApplicationTestSuite struct {
	suite.Suite
}

func (suite *ApplicationTestSuite) SetupSuite() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}
