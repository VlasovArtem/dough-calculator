package test

import (
	"go.uber.org/mock/gomock"
)

type GoMockTestSuite struct {
	ApplicationTestSuite
	MockCtrl *gomock.Controller
}

func (suite *GoMockTestSuite) SetupTest() {
	suite.MockCtrl = gomock.NewController(suite.T())
}

func (suite *GoMockTestSuite) TearDownTest() {
	if suite.MockCtrl == nil {
		return
	}
	suite.MockCtrl.Finish()
}
