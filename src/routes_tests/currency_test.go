package routes_tests

import (
	"github.com/stretchr/testify/suite"
	"main/forms"
	"main/models"
	"main/repository"
	"main/routes"
	"main/test_utils"
	data "main/test_utils/test_data"
	"main/utils"
)

type CurrencyRoutesSuit struct {
	suite.Suite
	router test_utils.RoutesDefaultSuite
}

func (suite *CurrencyRoutesSuit) SetupTest() {
	suite.router.SetupBase()
	routes.AddCurrencyRoutes(suite.router.Router)
}

func (suite *CurrencyRoutesSuit) TearDownTest() {
	test_utils.DropTestDB()
}

func (suite *CurrencyRoutesSuit) TestHandleGetAllCurrencies() {
	currency, err := repository.CreateCurrency(
		"Test currency",
		data.FirstCurrencySymbol,
		data.FirstCurrencyValue,
	)
	if err != nil {
		suite.Error(err)
	}

	resp := PerformTestRequest(
		suite.router.Router,
		"GET",
		"/currency/all/",
		nil,
		nil,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.CurrenciesResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	suite.Equal(1, len(response.Currencies))

	expectedForm, err := utils.MarshalModelToForm[models.Currency, forms.CurrencyResponse](
		currency,
	)
	if err != nil {
		suite.Error(err)
	}
	TestCurrenciesEqual(expectedForm, response.Currencies[0], &suite.Suite)
}
