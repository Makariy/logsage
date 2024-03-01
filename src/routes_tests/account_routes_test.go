package routes_tests

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"main/db_connector"
	"main/forms"
	"main/models"
	"main/routes"
	"main/test_utils"
	"main/utils"
)

type AccountRoutesSuit struct {
	suite.Suite
	router *gin.Engine

	user        *models.User
	currency    *models.Currency
	authHeaders map[string]string
}

func (suite *AccountRoutesSuit) SetupTest() {
	test_utils.CreateTestDB()
	models.MigrateModels(db_connector.GetConnection())

	suite.user = CreateTestUser(userEmail, userPassword)
	suite.currency = CreateTestCurrency(currencyName)
	suite.authHeaders = GetAuthHeaders(suite.user)

	suite.router = gin.Default()
	routes.AddAccountRoutes(suite.router)
}

func (suite *AccountRoutesSuit) TearDownTest() {
	test_utils.DropTestDB()
}

func getAccountForm(suite *AccountRoutesSuit) []byte {
	form := &forms.AccountForm{
		Name:       accountName,
		Balance:    accountBalance,
		CurrencyID: suite.currency.ID,
	}
	stringForm, _ := json.Marshal(&form)
	return stringForm
}

func (suite *AccountRoutesSuit) TestHandleCreateAccount() {
	resp := PerformTestRequest(
		suite.router,
		"POST",
		"/account/create/",
		getAccountForm(suite),
		&suite.authHeaders,
	)
	AssertResponseSuccess(201, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.AccountResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	TestAccountsEqual(&accountResponse, response, &suite.Suite)
}

func (suite *AccountRoutesSuit) TestHandleGetAccount() {
	account := CreateTestAccount(accountName, accountBalance, suite.user.ID, suite.currency.ID)

	resp := PerformTestRequest(
		suite.router,
		"GET",
		fmt.Sprintf("/account/get/%d/", account.ID),
		nil,
		&suite.authHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.AccountResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	expected, err := utils.MarshalModelToForm[models.Account, forms.AccountResponse](account)
	if err != nil {
		suite.Error(err)
	}

	TestAccountsEqual(expected, response, &suite.Suite)
}

func (suite *AccountRoutesSuit) TestHandleGetAllAccounts() {
	first := CreateTestAccount(accountName, accountBalance, suite.user.ID, suite.currency.ID)
	second := CreateTestAccount("Another account", decimal.NewFromInt(20), suite.user.ID, suite.currency.ID)

	resp := PerformTestRequest(
		suite.router,
		"GET",
		"/account/all/",
		getAccountForm(suite),
		&suite.authHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.AccountsResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	suite.Equal(2, len(response.Accounts))

	firstForm, err := utils.MarshalModelToForm[models.Account, forms.AccountResponse](first)
	if err != nil {
		suite.Error(err)
	}
	secondForm, err := utils.MarshalModelToForm[models.Account, forms.AccountResponse](second)
	if err != nil {
		suite.Error(err)
	}

	TestAccountsEqual(firstForm, response.Accounts[0], &suite.Suite)
	TestAccountsEqual(secondForm, response.Accounts[1], &suite.Suite)
}

func (suite *AccountRoutesSuit) TestHandlePatchAccount() {
	var (
		newName     = "New account name"
		newBalance  = decimal.NewFromInt(20)
		newCurrency = suite.currency
	)

	account := CreateTestAccount(accountName, accountBalance, suite.user.ID, suite.currency.ID)
	patchedAccount := forms.AccountForm{
		Name:       newName,
		Balance:    newBalance,
		CurrencyID: newCurrency.ID,
	}
	stringPatch, _ := json.Marshal(&patchedAccount)
	resp := PerformTestRequest(
		suite.router,
		"PATCH",
		fmt.Sprintf("/account/patch/%d/", account.ID),
		stringPatch,
		&suite.authHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.AccountResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	expectedAccount := models.Account{
		ID:       1,
		Name:     newName,
		Balance:  newBalance,
		Currency: *newCurrency,
		User:     *suite.user,
	}
	expectedForm, err := utils.MarshalModelToForm[models.Account, forms.AccountResponse](&expectedAccount)
	if err != nil {
		suite.Error(err)
	}
	TestAccountsEqual(expectedForm, response, &suite.Suite)
}

func (suite *AccountRoutesSuit) TestHandleDeleteAccount() {
	account := CreateTestAccount(accountName, accountBalance, suite.user.ID, suite.currency.ID)

	resp := PerformTestRequest(
		suite.router,
		"DELETE",
		fmt.Sprintf("/account/delete/%d/", account.ID),
		getAccountForm(suite),
		&suite.authHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.AccountResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	expectedForm, err := utils.MarshalModelToForm[models.Account, forms.AccountResponse](account)
	if err != nil {
		suite.Error(err)
	}
	TestAccountsEqual(expectedForm, response, &suite.Suite)
}
