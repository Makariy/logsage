package repository_tests

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"main/db_connector"
	"main/models"
	"main/repository"
	"main/test_utils"
	"time"
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

var (
	year     = 2024
	month    = time.January
	fromDate = time.Date(year, month, 01, 00, 00, 00, 00, time.Local)
	toDate   = time.Date(year, month, 10, 00, 00, 00, 00, time.Local)
)

var (
	transaction1, transaction2, transaction3, transaction4 *models.Transaction
)

func createTestTransactions(suite *StatsRepositorySuit) {
	transaction1, _ = repository.CreateTransaction(
		"First transaction",
		decimal.NewFromInt(200),
		time.Date(year, month, 01, 00, 00, 00, 00, time.Local),
		suite.user.ID,
		suite.firstCategory.ID,
		suite.firstAccount.ID,
	)
	transaction2, _ = repository.CreateTransaction(
		"Second transaction",
		decimal.NewFromInt(200),
		time.Date(year, month, 02, 00, 00, 00, 00, time.Local),
		suite.user.ID,
		suite.firstCategory.ID,
		suite.firstAccount.ID,
	)
	transaction3, _ = repository.CreateTransaction(
		"Third transaction",
		decimal.NewFromInt(200),
		time.Date(year, month, 03, 00, 00, 00, 00, time.Local),
		suite.user.ID,
		suite.secondCategory.ID,
		suite.secondAccount.ID,
	)
	transaction4, _ = repository.CreateTransaction(
		"FORTH transaction",
		decimal.NewFromInt(200),
		time.Date(year, month, 04, 00, 00, 00, 00, time.Local),
		suite.user.ID,
		suite.secondCategory.ID,
		suite.secondAccount.ID,
	)
}

func (suite *StatsRepositorySuit) SetupTest() {
	test_utils.CreateTestDB()
	models.MigrateModels(db_connector.GetConnection())

	suite.user = CreateTestUser(userEmail, userPassword)
	suite.currency = CreateTestCurrency(currencyName)

	suite.firstCategory = CreateTestCategory(suite.user.ID, "First category", models.SPENDING)
	suite.secondCategory = CreateTestCategory(suite.user.ID, "Second category", models.EARNING)

	suite.firstAccount = CreateTestAccount("First account", accountBalance, suite.user.ID, suite.currency.ID)
	suite.secondAccount = CreateTestAccount("Second account", decimal.NewFromInt(10000), suite.user.ID, suite.currency.ID)

	createTestTransactions(suite)
}

func (suite *StatsRepositorySuit) TearDownTest() {
	test_utils.DropTestDB()
}

func (suite *StatsRepositorySuit) TestGetCategoryStats() {
	categoryStats, err := repository.GetCategoryStats(
		suite.firstCategory.ID,
		fromDate,
		toDate,
	)
	suite.True(err == nil)

	expected := models.CategoryStats{
		Category:     *suite.firstCategory,
		TotalAmount:  decimal.NewFromInt(400),
		TotalPercent: decimal.NewFromFloat(1. / 2),
		Transactions: []*models.Transaction{
			transaction1,
			transaction2,
		},
	}

	TestCategoriesStatsEqual(&expected, categoryStats, &suite.Suite)
}

func (suite *StatsRepositorySuit) TestGetAccountStats() {
	accountStats, err := repository.GetAccountStats(
		suite.firstAccount.ID,
		fromDate,
		toDate,
	)
	suite.True(err == nil)

	expected := models.AccountStats{
		Account:            *suite.firstAccount,
		TotalSpentAmount:   decimal.NewFromInt(400),
		TotalSpentPercent:  decimal.NewFromFloat(1),
		TotalEarnedAmount:  decimal.NewFromInt(0),
		TotalEarnedPercent: decimal.NewFromInt(0),
		Transactions: []*models.Transaction{
			transaction1,
			transaction2,
		},
	}

	TestAccountStatsEqual(&expected, accountStats, &suite.Suite)
}

func (suite *StatsRepositorySuit) TestGetTotalAccountsStats() {
	totalStats, err := repository.GetTotalAccountsStats(suite.user.ID, fromDate, toDate)
	suite.True(err == nil)

	firstAccountStats := models.AccountStats{
		Account:            *suite.firstAccount,
		TotalSpentAmount:   decimal.NewFromInt(400),
		TotalSpentPercent:  decimal.NewFromFloat(1),
		TotalEarnedAmount:  decimal.NewFromInt(0),
		TotalEarnedPercent: decimal.NewFromInt(0),
		Transactions: []*models.Transaction{
			transaction1,
			transaction2,
		},
	}

	secondAccountStats := models.AccountStats{
		Account:            *suite.secondAccount,
		TotalSpentAmount:   decimal.NewFromInt(0),
		TotalSpentPercent:  decimal.NewFromFloat(0),
		TotalEarnedAmount:  decimal.NewFromInt(400),
		TotalEarnedPercent: decimal.NewFromInt(1),
		Transactions: []*models.Transaction{
			transaction3,
			transaction4,
		},
	}

	expected := models.TotalAccountsStats{
		TotalEarnedAmount: decimal.NewFromInt(400),
		TotalSpentAmount:  decimal.NewFromInt(400),
		AccountsStats: []*models.AccountStats{
			&firstAccountStats,
			&secondAccountStats,
		},
	}

	TestTotalAccountsStatsEqual(&expected, totalStats, &suite.Suite)
}

func (suite *StatsRepositorySuit) TestGetTotalCategoriesStats() {
	totalStats, err := repository.GetTotalCategoriesStats(suite.user.ID, fromDate, toDate)
	suite.True(err == nil)

	firstCategoryStats := models.CategoryStats{
		Category:     *suite.firstCategory,
		TotalAmount:  decimal.NewFromInt(400),
		TotalPercent: decimal.NewFromFloat(1. / 2),
		Transactions: []*models.Transaction{
			transaction1,
			transaction2,
		},
	}
	secondCategoryStats := models.CategoryStats{
		Category:     *suite.secondCategory,
		TotalAmount:  decimal.NewFromInt(400),
		TotalPercent: decimal.NewFromFloat(1. / 2),
		Transactions: []*models.Transaction{
			transaction3,
			transaction4,
		},
	}

	expected := models.TotalCategoriesStats{
		TotalEarnedAmount: decimal.NewFromInt(400),
		TotalSpentAmount:  decimal.NewFromInt(400),
		CategoriesStats: []*models.CategoryStats{
			&firstCategoryStats,
			&secondCategoryStats,
		},
	}

	TestTotalCategoriesStatsEqual(&expected, totalStats, &suite.Suite)
}
