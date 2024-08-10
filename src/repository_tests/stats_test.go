package repository_tests

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"main/models"
	"main/repository"
	"time"
)

type StatsRepositorySuit struct {
	suite.Suite
	base DefaultRepositorySuite
}

var (
	fromDate = transaction1Date.Add(time.Hour * -1)
	toDate   = transaction4Date.Add(time.Hour)
)

func (suite *StatsRepositorySuit) SetupTest() {
	suite.base.SetupTest()
}

func (suite *StatsRepositorySuit) TearDownTest() {
	suite.base.TearDownTest()
}

func (suite *StatsRepositorySuit) TestGetCategoryStats() {
	categoryStats, err := repository.GetCategoryStats(
		suite.base.firstCategory.ID,
		fromDate,
		toDate,
		suite.base.firstCurrency,
	)
	suite.True(err == nil)

	expected := models.CategoryStats{
		Category:    *suite.base.firstCategory,
		TotalAmount: transaction1Amount.Add(transaction2Amount),
		Transactions: []*models.Transaction{
			suite.base.transaction1,
			suite.base.transaction2,
		},
	}

	TestCategoriesStatsEqual(&expected, categoryStats, &suite.Suite)
}

func (suite *StatsRepositorySuit) TestGetAccountStats() {
	accountStats, err := repository.GetAccountStats(
		suite.base.firstAccount.ID,
		fromDate,
		toDate,
	)
	suite.True(err == nil)

	expected := models.AccountStats{
		Account:           *suite.base.firstAccount,
		TotalSpentAmount:  transaction1Amount.Add(transaction2Amount),
		TotalEarnedAmount: decimal.Zero,
		Transactions: []*models.Transaction{
			suite.base.transaction1,
			suite.base.transaction2,
		},
	}

	TestAccountStatsEqual(&expected, accountStats, &suite.Suite)
}

func (suite *StatsRepositorySuit) TestGetTotalAccountsStats() {
	totalStats, err := repository.GetTotalAccountsStats(suite.base.user.ID, fromDate, toDate)
	suite.True(err == nil)

	expectedSpentAmount := transaction1Amount.Add(transaction2Amount)
	expectedEarnedAmount := transaction3Amount.
		Add(transaction4Amount).
		Mul(secondCurrencyValue)

	firstAccountStats := models.AccountStats{
		Account:           *suite.base.firstAccount,
		TotalSpentAmount:  expectedSpentAmount,
		TotalEarnedAmount: decimal.Zero,
		Transactions: []*models.Transaction{
			suite.base.transaction1,
			suite.base.transaction2,
		},
	}

	secondAccountStats := models.AccountStats{
		Account:           *suite.base.secondAccount,
		TotalSpentAmount:  decimal.Zero,
		TotalEarnedAmount: expectedEarnedAmount,
		Transactions: []*models.Transaction{
			suite.base.transaction3,
			suite.base.transaction4,
		},
	}

	expected := models.TotalAccountsStats{
		TotalEarnedAmount: expectedEarnedAmount,
		TotalSpentAmount:  expectedSpentAmount,
		AccountsStats: []*models.AccountStats{
			&firstAccountStats,
			&secondAccountStats,
		},
	}

	TestTotalAccountsStatsEqual(&expected, totalStats, &suite.Suite)
}

func (suite *StatsRepositorySuit) TestGetTotalCategoriesStats() {
	totalStats, err := repository.GetTotalCategoriesStats(
		suite.base.user.ID,
		fromDate,
		toDate,
		suite.base.firstCurrency,
	)
	suite.True(err == nil)

	expectedSpentAmount := transaction1Amount.Add(transaction2Amount)
	expectedEarnedAmount := transaction3Amount.
		Add(transaction4Amount).
		Mul(secondCurrencyValue)

	firstCategoryStats := models.CategoryStats{
		Category:    *suite.base.firstCategory,
		TotalAmount: expectedSpentAmount,
		Transactions: []*models.Transaction{
			suite.base.transaction1,
			suite.base.transaction2,
		},
	}
	secondCategoryStats := models.CategoryStats{
		Category:    *suite.base.secondCategory,
		TotalAmount: expectedEarnedAmount,
		Transactions: []*models.Transaction{
			suite.base.transaction3,
			suite.base.transaction4,
		},
	}

	expected := models.TotalCategoriesStats{
		TotalEarnedAmount: expectedEarnedAmount,
		TotalSpentAmount:  expectedSpentAmount,
		CategoriesStats: []*models.CategoryStats{
			&firstCategoryStats,
			&secondCategoryStats,
		},
	}

	TestTotalCategoriesStatsEqual(&expected, totalStats, &suite.Suite)
}
