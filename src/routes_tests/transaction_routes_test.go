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
	"main/repository"
	"main/routes"
	"main/test_utils"
	"time"
)

type TransactionRoutesSuit struct {
	suite.Suite
	router *gin.Engine

	user        *models.User
	currency    *models.Currency
	category    *models.Category
	account     *models.Account
	authHeaders map[string]string
}

func (suite *TransactionRoutesSuit) SetupTest() {
	test_utils.CreateTestDB()
	models.MigrateModels(db_connector.GetConnection())

	suite.user = CreateTestUser(userEmail, userPassword)
	suite.currency = CreateTestCurrency(currencyName)
	suite.category = CreateTestCategory(categoryName, categoryType, suite.user.ID)
	suite.account = CreateTestAccount(accountName, accountBalance, suite.user.ID, suite.currency.ID)
	suite.authHeaders = GetAuthHeaders(suite.user)

	suite.router = gin.Default()
	routes.AddTransactionRoutes(suite.router)
}

func (suite *TransactionRoutesSuit) TearDownTest() {
	test_utils.DropTestDB()
}

func getTransactionForm(suite *TransactionRoutesSuit) []byte {
	form := &forms.TransactionForm{
		Description: transactionDescription,
		Amount:      transactionAmount,
		Date:        transactionDate,
		CategoryID:  suite.category.ID,
		AccountID:   suite.account.ID,
	}
	stringForm, _ := json.Marshal(&form)
	return stringForm
}

func (suite *TransactionRoutesSuit) TestHandleCreateTransaction() {
	resp := PerformTestRequest(
		suite.router,
		"POST",
		"/transaction/create/",
		getTransactionForm(suite),
		&suite.authHeaders,
	)
	AssertResponseSuccess(201, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.TransactionResponse](resp)
	if err != nil {
		suite.Error(err, "Пиздец")
	}

	testTransactionsEqual(&transactionResponse, response, &suite.Suite)
}

func (suite *TransactionRoutesSuit) TestHandleGetTransaction() {
	transaction, err := repository.CreateTransaction(
		transactionDescription,
		transactionAmount,
		transactionDate,
		suite.user.ID,
		suite.category.ID,
		suite.account.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	resp := PerformTestRequest(
		suite.router,
		"GET",
		fmt.Sprintf("/transaction/get/%d/", transaction.ID),
		nil,
		&suite.authHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.TransactionResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	expected, err := MarshalModelToForm[models.Transaction, forms.TransactionResponse](transaction)
	if err != nil {
		suite.Error(err)
	}

	testTransactionsEqual(expected, response, &suite.Suite)
}

func (suite *TransactionRoutesSuit) TestHandleGetAllTransactions() {
	first, err := repository.CreateTransaction(
		transactionDescription,
		transactionAmount,
		transactionDate,
		suite.user.ID,
		suite.category.ID,
		suite.account.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	second, err := repository.CreateTransaction(
		"New description",
		decimal.NewFromInt(21),
		time.Now(),
		suite.user.ID,
		suite.category.ID,
		suite.account.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	resp := PerformTestRequest(
		suite.router,
		"GET",
		"/transaction/all/",
		getTransactionForm(suite),
		&suite.authHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.TransactionsResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	suite.Equal(2, len(response.Transactions))

	firstForm, err := MarshalModelToForm[models.Transaction, forms.TransactionResponse](first)
	if err != nil {
		suite.Error(err)
	}
	secondForm, err := MarshalModelToForm[models.Transaction, forms.TransactionResponse](second)
	if err != nil {
		suite.Error(err)
	}

	testTransactionsEqual(firstForm, response.Transactions[0], &suite.Suite)
	testTransactionsEqual(secondForm, response.Transactions[1], &suite.Suite)
}

func (suite *TransactionRoutesSuit) TestHandlePatchTransaction() {
	var (
		newDescription = "New transaction description"
		newAmount      = decimal.NewFromInt(20)
		newDate        = time.Now()
	)

	transaction, err := repository.CreateTransaction(
		transactionDescription,
		transactionAmount,
		transactionDate,
		suite.user.ID,
		suite.category.ID,
		suite.account.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	patchedTransaction := forms.TransactionForm{
		Description: newDescription,
		Amount:      newAmount,
		Date:        newDate,
	}
	stringPatch, _ := json.Marshal(&patchedTransaction)
	resp := PerformTestRequest(
		suite.router,
		"PATCH",
		fmt.Sprintf("/transaction/patch/%d/", transaction.ID),
		stringPatch,
		&suite.authHeaders,
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
		User:        *suite.user,
		Category:    *suite.category,
	}
	expectedForm, err := MarshalModelToForm[models.Transaction, forms.TransactionResponse](&expectedTransaction)
	if err != nil {
		suite.Error(err)
	}
	testTransactionsEqual(expectedForm, response, &suite.Suite)
}

func (suite *TransactionRoutesSuit) TestHandleDeleteTransaction() {
	transaction, err := repository.CreateTransaction(
		transactionDescription,
		transactionAmount,
		transactionDate,
		suite.user.ID,
		suite.category.ID,
		suite.account.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	resp := PerformTestRequest(
		suite.router,
		"DELETE",
		fmt.Sprintf("/transaction/delete/%d/", transaction.ID),
		getTransactionForm(suite),
		&suite.authHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.TransactionResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	expectedForm, err := MarshalModelToForm[models.Transaction, forms.TransactionResponse](transaction)
	if err != nil {
		suite.Error(err)
	}
	testTransactionsEqual(expectedForm, response, &suite.Suite)
}
