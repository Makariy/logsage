package repository_tests

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"main/db_connector"
	"main/models"
	"main/test_utils"
)

type StatsRepositorySuit struct {
	suite.Suite
	user     *models.User
	currency *models.Currency

	firstCategory  *models.Category
	secondCategory *models.Category

	firstAccount  *models.Account
	secondAccount *models.Account
}

func (suite *StatsRepositorySuit) SetupTest() {
	test_utils.CreateTestDB()
	models.MigrateModels(db_connector.GetConnection())

	suite.user = CreateTestUser(userEmail, userPassword)
	suite.currency = CreateTestCurrency(currencyName)

	suite.firstCategory = CreateTestCategory(suite.user.ID, categoryName, categoryType)
	suite.secondCategory = CreateTestCategory(suite.user.ID, "Second category", models.EARNING)

	suite.firstAccount = CreateTestAccount(accountName, accountBalance, suite.user.ID, suite.currency.ID)
	suite.secondAccount = CreateTestAccount("Second account", decimal.New(10000, 10), suite.user.ID, suite.currency.ID)
}

func (suite *StatsRepositorySuit) TearDownTest() {
	test_utils.DropTestDB()
}
