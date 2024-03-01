package routes_tests

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"main/db_connector"
	"main/forms"
	"main/models"
	"main/repository"
	"main/routes"
	"main/test_utils"
	"main/utils"
)

type CurrencyRoutesSuit struct {
	suite.Suite
	router *gin.Engine
}

func (suite *CurrencyRoutesSuit) SetupTest() {
	test_utils.CreateTestDB()
	models.MigrateModels(db_connector.GetConnection())

	suite.router = gin.Default()
	routes.AddCurrencyRoutes(suite.router)
}

func (suite *CurrencyRoutesSuit) TearDownTest() {
	test_utils.DropTestDB()
}

func (suite *CurrencyRoutesSuit) TestHandleGetAllCurrencies() {
	currency, err := repository.CreateCurrency(currencyName)
	if err != nil {
		suite.Error(err)
	}

	resp := PerformTestRequest(suite.router, "GET", "/currency/all/", nil, nil)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.CurrenciesResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	suite.Equal(1, len(response.Currencies))

	expectedForm, err := utils.MarshalModelToForm[models.Currency, forms.CurrencyResponse](currency)
	if err != nil {
		suite.Error(err)
	}
	TestCurrenciesEqual(expectedForm, response.Currencies[0], &suite.Suite)
}
