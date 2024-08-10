package repository_tests

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"main/db_connector"
	"main/models"
	"main/repository"
	"main/test_utils"
)

type DefaultRepositorySuite struct {
	suite.Suite

	user *models.User

	firstCurrency  *models.Currency
	secondCurrency *models.Currency

	firstCategory  *models.Category
	secondCategory *models.Category

	firstAccount  *models.Account
	secondAccount *models.Account

	transaction1 *models.Transaction
	transaction2 *models.Transaction
	transaction3 *models.Transaction
	transaction4 *models.Transaction
}

func (suite *DefaultRepositorySuite) setupDB() {
	test_utils.CreateTestDB()
	models.MigrateModels(db_connector.GetConnection())
}

func (suite *DefaultRepositorySuite) createDefaultUser() {
	suite.user = CreateTestUser(userEmail, userPassword)
}

func (suite *DefaultRepositorySuite) createDefaultCurrencies() {
	suite.firstCurrency = CreateTestCurrency(firstCurrencySymbol, firstCurrencyValue)
	suite.secondCurrency = CreateTestCurrency(secondCurrencySymbol, secondCurrencyValue)
}

func (suite *DefaultRepositorySuite) createDefaultCategories() {
	suite.firstCategory = CreateTestCategory(
		suite.user.ID,
		"First category",
		models.SPENDING,
	)
	suite.secondCategory = CreateTestCategory(
		suite.user.ID,
		"Second category",
		models.EARNING,
	)
}

func (suite *DefaultRepositorySuite) createDefaultAccounts() {
	suite.firstAccount = CreateTestAccount(
		"First account",
		accountBalance,
		suite.user.ID,
		suite.firstCurrency.ID,
	)

	suite.secondAccount = CreateTestAccount(
		"Second account",
		decimal.NewFromInt(10000),
		suite.user.ID,
		suite.secondCurrency.ID,
	)
}

func (suite *DefaultRepositorySuite) createDefaultTransactions() {
	suite.transaction1, _ = repository.CreateTransaction(
		"First transaction",
		transaction1Amount,
		transaction1Date,
		suite.user.ID,
		suite.firstCategory.ID,
		suite.firstAccount.ID,
	)
	suite.transaction2, _ = repository.CreateTransaction(
		"Second transaction",
		transaction2Amount,
		transaction2Date,
		suite.user.ID,
		suite.firstCategory.ID,
		suite.firstAccount.ID,
	)
	suite.transaction3, _ = repository.CreateTransaction(
		"Third transaction",
		transaction3Amount,
		transaction3Date,
		suite.user.ID,
		suite.secondCategory.ID,
		suite.secondAccount.ID,
	)
	suite.transaction4, _ = repository.CreateTransaction(
		"Forth transaction",
		transaction4Amount,
		transaction4Date,
		suite.user.ID,
		suite.secondCategory.ID,
		suite.secondAccount.ID,
	)
}

func (suite *DefaultRepositorySuite) SetupTest() {
	suite.setupDB()

	suite.createDefaultCurrencies()
	suite.createDefaultUser()
	suite.createDefaultCategories()
	suite.createDefaultAccounts()
	suite.createDefaultTransactions()
}

func (suite *DefaultRepositorySuite) TearDownTest() {
	test_utils.DropTestDB()
}
