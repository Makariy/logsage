package routes_tests

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"main/forms"
	"main/models"
	"main/routes"
	"main/test_utils"
	data "main/test_utils/test_data"
	"main/utils"
)

type AccountRoutesSuit struct {
	suite.Suite
	router test_utils.RoutesDefaultSuite
}

func (suite *AccountRoutesSuit) SetupTest() {
	suite.router.SetupBase()
	suite.router.Data.CreateDefaultUser()
	suite.router.SetupAuth()
	suite.router.Data.CreateDefaultCurrencies()
	routes.AddAccountRoutes(suite.router.Router)
}

func (suite *AccountRoutesSuit) TearDownTest() {
	suite.router.TearDownTest()
}

func getAccountForm(suite *AccountRoutesSuit) []byte {
	form := &forms.AccountForm{
		Name:       data.FirstAccountName,
		Balance:    data.FirstAccountBalance,
		CurrencyID: suite.router.Data.FirstCurrency.ID,
	}
	stringForm, _ := json.Marshal(&form)
	return stringForm
}

func (suite *AccountRoutesSuit) TestHandleCreateAccount() {
	resp := PerformTestRequest(
		suite.router.Router,
		"POST",
		"/account/create/",
		getAccountForm(suite),
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(201, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.AccountResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	TestAccountsEqual(&suite.Suite, &accountResponse, response)
}

func (suite *AccountRoutesSuit) TestHandleGetAccount() {
	suite.router.Data.CreateDefaultAccounts()

	account := suite.router.Data.FirstAccount

	resp := PerformTestRequest(
		suite.router.Router,
		"GET",
		fmt.Sprintf("/account/get/%d/", account.ID),
		nil,
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.AccountResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	expected, err := utils.MarshalModelToForm[models.Account, forms.AccountResponse](
		account,
	)
	if err != nil {
		suite.Error(err)
	}

	TestAccountsEqual(&suite.Suite, expected, response)
}

func (suite *AccountRoutesSuit) TestHandleGetAllAccounts() {
	suite.router.Data.CreateDefaultAccounts()
	firstAccount := suite.router.Data.FirstAccount
	secondAccount := suite.router.Data.SecondAccount

	resp := PerformTestRequest(
		suite.router.Router,
		"GET",
		"/account/all/",
		getAccountForm(suite),
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.AccountsResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	suite.Equal(2, len(response.Accounts))

	firstForm, err := utils.MarshalModelToForm[models.Account, forms.AccountResponse](
		firstAccount,
	)
	if err != nil {
		suite.Error(err)
	}
	secondForm, err := utils.MarshalModelToForm[models.Account, forms.AccountResponse](
		secondAccount,
	)
	if err != nil {
		suite.Error(err)
	}

	TestAccountsEqual(&suite.Suite, firstForm, response.Accounts[0])
	TestAccountsEqual(&suite.Suite, secondForm, response.Accounts[1])
}

func (suite *AccountRoutesSuit) TestHandlePatchAccount() {
	var (
		newName     = "New account name"
		newBalance  = decimal.NewFromInt(20)
		newCurrency = suite.router.Data.FirstCurrency
	)

	suite.router.Data.CreateDefaultAccounts()
	account := suite.router.Data.FirstAccount

	patchedAccount := forms.AccountForm{
		Name:       newName,
		Balance:    newBalance,
		CurrencyID: newCurrency.ID,
	}
	stringPatch, _ := json.Marshal(&patchedAccount)
	resp := PerformTestRequest(
		suite.router.Router,
		"PATCH",
		fmt.Sprintf("/account/patch/%d/", account.ID),
		stringPatch,
		&suite.router.AuthHeaders,
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
		User:     *suite.router.Data.User,
	}
	expectedForm, err := utils.MarshalModelToForm[models.Account, forms.AccountResponse](&expectedAccount)
	if err != nil {
		suite.Error(err)
	}
	TestAccountsEqual(&suite.Suite, expectedForm, response)
}

func (suite *AccountRoutesSuit) TestHandleDeleteAccount() {
	suite.router.Data.CreateDefaultAccounts()
	account := suite.router.Data.FirstAccount

	resp := PerformTestRequest(
		suite.router.Router,
		"DELETE",
		fmt.Sprintf("/account/delete/%d/", account.ID),
		getAccountForm(suite),
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.AccountResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	expectedForm, err := utils.MarshalModelToForm[models.Account, forms.AccountResponse](
		account,
	)
	if err != nil {
		suite.Error(err)
	}
	TestAccountsEqual(&suite.Suite, expectedForm, response)
}
