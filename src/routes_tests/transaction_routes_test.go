package routes_tests

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"main/forms"
	"main/models"
	"main/test_utils"
	data "main/test_utils/test_data"
	"main/utils"
	"time"
)

type TransactionRoutesSuit struct {
	suite.Suite
	router test_utils.RoutesDefaultSuite
}

func (suite *TransactionRoutesSuit) SetupTest() {
	suite.router.SetupBase()
	suite.router.Data.CreateDefaultUser()
	suite.router.SetupAuth()
	suite.router.SetupRoutes()
	suite.router.Data.CreateDefaultCurrencies()
	suite.router.Data.CreateDefaultAccounts()
	suite.router.Data.CreateDefaultCategories()
}

func (suite *TransactionRoutesSuit) TearDownTest() {
	suite.router.TearDownTest()
}

func getTransactionForm(suite *TransactionRoutesSuit) []byte {
	form := &forms.TransactionForm{
		Description: data.TransactionDescription,
		Amount:      data.TransactionAmount,
		Date:        data.TransactionDate,
		CategoryID:  suite.router.Data.FirstCategory.ID,
		AccountID:   suite.router.Data.FirstAccount.ID,
	}
	stringForm, _ := json.Marshal(&form)
	return stringForm
}

func (suite *TransactionRoutesSuit) TestHandleCreateTransaction() {
	resp := PerformTestRequest(
		suite.router.Router,
		"POST",
		"/transaction/create/",
		getTransactionForm(suite),
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(201, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.TransactionResponse](resp)
	if err != nil {
		suite.Error(err, "Пиздец")
	}

	TestTransactionsEqual(&suite.Suite, &transactionResponse, response)
}

func (suite *TransactionRoutesSuit) TestHandleGetTransaction() {
	suite.router.Data.CreateDefaultTransactions()
	transaction := suite.router.Data.Transaction1

	resp := PerformTestRequest(
		suite.router.Router,
		"GET",
		fmt.Sprintf("/transaction/get/%d/", transaction.ID),
		nil,
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.TransactionResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	expected, err := utils.MarshalModelToForm[models.Transaction, forms.TransactionResponse](transaction)
	if err != nil {
		suite.Error(err)
	}

	TestTransactionsEqual(&suite.Suite, expected, response)
}

func (suite *TransactionRoutesSuit) TestHandleGetAllTransactions() {
	suite.router.Data.CreateDefaultTransactions()
	firstTransaction := suite.router.Data.Transaction1
	secondTransaction := suite.router.Data.Transaction2

	resp := PerformTestRequest(
		suite.router.Router,
		"GET",
		fmt.Sprintf(
			"/transaction/all/?%s",
			renderDateRangeForTransactions(firstTransaction, secondTransaction),
		),
		getTransactionForm(suite),
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.TransactionsResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	suite.Equal(2, len(response.Transactions))

	firstForm, err := utils.MarshalModelToForm[models.Transaction, forms.TransactionResponse](
		firstTransaction,
	)
	if err != nil {
		suite.Error(err)
	}
	secondForm, err := utils.MarshalModelToForm[models.Transaction, forms.TransactionResponse](
		secondTransaction,
	)
	if err != nil {
		suite.Error(err)
	}

	TestTransactionsEqual(&suite.Suite, firstForm, response.Transactions[1])
	TestTransactionsEqual(&suite.Suite, secondForm, response.Transactions[0])
}

func (suite *TransactionRoutesSuit) TestHandlePatchTransaction() {
	var (
		newDescription = "New transaction description"
		newAmount      = decimal.NewFromInt(20)
		newDate        = time.Now()
	)

	suite.router.Data.CreateDefaultTransactions()
	transaction := suite.router.Data.Transaction1

	patchedTransaction := forms.TransactionForm{
		Description: newDescription,
		Amount:      newAmount,
		Date:        newDate,
	}
	stringPatch, _ := json.Marshal(&patchedTransaction)
	resp := PerformTestRequest(
		suite.router.Router,
		"PATCH",
		fmt.Sprintf("/transaction/patch/%d/", transaction.ID),
		stringPatch,
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.TransactionResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	expectedTransaction := models.Transaction{
		ID:          1,
		Description: newDescription,
		Amount:      newAmount,
		Date:        newDate,
		User:        *suite.router.Data.User,
		Category:    *suite.router.Data.FirstCategory,
		Account:     *suite.router.Data.FirstAccount,
	}
	expectedForm, err := utils.MarshalModelToForm[models.Transaction, forms.TransactionResponse](&expectedTransaction)
	if err != nil {
		suite.Error(err)
	}
	TestTransactionsEqual(&suite.Suite, expectedForm, response)
}

func (suite *TransactionRoutesSuit) TestHandleDeleteTransaction() {
	suite.router.Data.CreateDefaultTransactions()
	transaction := suite.router.Data.Transaction1

	resp := PerformTestRequest(
		suite.router.Router,
		"DELETE",
		fmt.Sprintf("/transaction/delete/%d/", transaction.ID),
		getTransactionForm(suite),
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.TransactionResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	expectedForm, err := utils.MarshalModelToForm[models.Transaction, forms.TransactionResponse](transaction)
	if err != nil {
		suite.Error(err)
	}
	TestTransactionsEqual(&suite.Suite, expectedForm, response)
}
