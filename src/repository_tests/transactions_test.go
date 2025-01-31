package repository_tests

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"main/db_connector"
	"main/models"
	"main/repository"
	"main/test_utils"
	"time"
)

type TransactionRepositorySuit struct {
	suite.Suite
	router   *gin.Engine
	user     *models.User
	currency *models.Currency
	category *models.Category
	account  *models.Account
}

func (suite *TransactionRepositorySuit) SetupTest() {
	test_utils.CreateTestDB()
	models.MigrateModels(db_connector.GetConnection())

	suite.user = CreateTestUser(userEmail, userPassword)
	suite.currency = CreateTestCurrency(currencyName)
	suite.category = CreateTestCategory(suite.user.ID, categoryName, categoryType)
	suite.account = CreateTestAccount(accountName, accountBalance, suite.user.ID, suite.currency.ID)
}

func (suite *TransactionRepositorySuit) TearDownTest() {
	test_utils.DropTestDB()
}

func (suite *TransactionRepositorySuit) TestCreateTransaction() {
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

	expected := models.Transaction{
		ID:          1,
		Description: transactionDescription,
		Amount:      transactionAmount,
		Date:        transactionDate,
		User:        *suite.user,
		Category:    *suite.category,
		Account:     *suite.account,
	}

	TestTransactionsEqual(&expected, transaction, &suite.Suite)
}

func (suite *TransactionRepositorySuit) TestGetTransactionByID() {
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

	result, err := repository.GetTransactionByID(transaction.ID)
	if err != nil {
		suite.Error(err)
	}

	expected := models.Transaction{
		ID:          transaction.ID,
		Description: transaction.Description,
		Amount:      transaction.Amount,
		Date:        transaction.Date,
		User:        transaction.User,
		Category:    transaction.Category,
		Account:     transaction.Account,
	}

	TestTransactionsEqual(&expected, result, &suite.Suite)
}

func (suite *TransactionRepositorySuit) TestGetAllTransactions() {
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
		"Other transaction",
		decimal.NewFromInt(500),
		time.Now(),
		suite.user.ID,
		suite.category.ID,
		suite.account.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	transactions, err := repository.GetUserTransactions(suite.user.ID)
	if err != nil {
		suite.Error(err)
	}
	suite.Equal(2, len(transactions))

	for _, transaction := range transactions {
		TestUsersEqual(suite.user, &transaction.User, &suite.Suite)
	}

	isFirstFirst := transactions[0].ID == first.ID
	if isFirstFirst {
		TestTransactionsEqual(first, transactions[0], &suite.Suite)
		TestTransactionsEqual(second, transactions[1], &suite.Suite)
	} else {
		TestTransactionsEqual(first, transactions[1], &suite.Suite)
		TestTransactionsEqual(second, transactions[0], &suite.Suite)
	}
}

func (suite *TransactionRepositorySuit) TestPatchTransaction() {
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

	var (
		newDescription = "New transaction description"
		newAmount      = decimal.NewFromInt(600)
		newDate        = time.Now()
	)

	patched, err := repository.PatchTransaction(
		transaction.ID,
		newDescription,
		newAmount,
		newDate,
		suite.user.ID,
		suite.category.ID,
		suite.account.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	expected := models.Transaction{
		ID:          transaction.ID,
		Description: newDescription,
		Amount:      newAmount,
		Date:        newDate,
		User:        transaction.User,
		Category:    transaction.Category,
		Account:     transaction.Account,
	}

	TestTransactionsEqual(&expected, patched, &suite.Suite)
}

func (suite *TransactionRepositorySuit) TestDeleteTransaction() {
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

	result, err := repository.DeleteTransaction(transaction.ID)
	if err != nil {
		suite.Error(err)
	}

	TestTransactionsEqual(transaction, result, &suite.Suite)

	transactions, err := repository.GetUserTransactions(transaction.User.ID)
	if err != nil {
		suite.Error(err)
	}
	suite.Equal(0, len(transactions))
}
