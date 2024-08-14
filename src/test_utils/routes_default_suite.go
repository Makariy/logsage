package test_utils

import (
	"github.com/gin-gonic/gin"
	"main/routes"
)

type RoutesDefaultSuite struct {
	Data RepositoryDefaultSuite

	Router      *gin.Engine
	AuthHeaders map[string]string
}

func (suite *RoutesDefaultSuite) SetupBase() {
	suite.Router = gin.Default()
	suite.Data.SetupDB()
}

func (suite *RoutesDefaultSuite) SetupAuth() {
	suite.AuthHeaders = GetAuthHeaders(suite.Data.User)
}

func (suite *RoutesDefaultSuite) SetupAllTestData() {
	suite.SetupBase()
	suite.Data.SetupAllTestData()
	suite.SetupRoutes()
	suite.SetupAuth()
}

func (suite *RoutesDefaultSuite) SetupRoutes() {
	routes.SetupAllRoutes(suite.Router)
}

func (suite *RoutesDefaultSuite) TearDownTest() {
	suite.Data.TearDownTest()
}
