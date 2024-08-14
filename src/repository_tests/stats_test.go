package repository_tests

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"main/forms"
	"main/models"
	"main/repository"
	"main/test_utils"
	"time"
)

type StatsRepositorySuit struct {
	suite.Suite
	base test_utils.RepositoryDefaultSuite
}

var (
	fromDate = transaction1Date.Add(time.Hour * -1)
	toDate   = transaction4Date.Add(time.Hour)
)

func (suite *StatsRepositorySuit) SetupTest() {
	suite.base.SetupDB()
	suite.base.SetupAllTestData()
}

func (suite *StatsRepositorySuit) TearDownTest() {
	suite.base.TearDownTest()
}

func (suite *StatsRepositorySuit) TestGetCategoryStats() {
	categoryStats, err := repository.GetCategoryStats(
		suite.base.FirstCategory.ID,
		fromDate,
		toDate,
		suite.base.FirstCurrency,
	)
	suite.True(err == nil)

	expected := forms.CategoryStats{
		Category:    *suite.base.FirstCategory,
		TotalAmount: transaction1Amount.Add(transaction2Amount),
		Transactions: []*models.Transaction{
			suite.base.Transaction1,
			suite.base.Transaction2,
		},
	}

	TestCategoriesStatsEqual(&expected, categoryStats, &suite.Suite)
}

func (suite *StatsRepositorySuit) TestGetAccountStats() {
	accountStats, err := repository.GetAccountStats(
		suite.base.FirstAccount.ID,
		fromDate,
		toDate,
	)
	suite.True(err == nil)

	expected := forms.AccountStats{
		Account:           *suite.base.FirstAccount,
		TotalSpentAmount:  transaction1Amount.Add(transaction2Amount),
		TotalEarnedAmount: decimal.Zero,
		Transactions: []*models.Transaction{
			suite.base.Transaction1,
			suite.base.Transaction2,
		},
	}

	TestAccountStatsEqual(&expected, accountStats, &suite.Suite)
}

func (suite *StatsRepositorySuit) TestGetTotalAccountsStats() {
	totalStats, err := repository.GetTotalAccountsStats(suite.base.User.ID, fromDate, toDate)
	suite.True(err == nil)

	expectedSpentAmount := transaction1Amount.Add(transaction2Amount)
	expectedEarnedAmount := transaction3Amount.
		Add(transaction4Amount).
		Mul(secondCurrencyValue)

	firstAccountStats := forms.AccountStats{
		Account:           *suite.base.FirstAccount,
		TotalSpentAmount:  expectedSpentAmount,
		TotalEarnedAmount: decimal.Zero,
		Transactions: []*models.Transaction{
			suite.base.Transaction1,
			suite.base.Transaction2,
		},
	}

	secondAccountStats := forms.AccountStats{
		Account:           *suite.base.SecondAccount,
		TotalSpentAmount:  decimal.Zero,
		TotalEarnedAmount: expectedEarnedAmount,
		Transactions: []*models.Transaction{
			suite.base.Transaction3,
			suite.base.Transaction4,
		},
	}

	expected := forms.TotalAccountsStats{
		TotalEarnedAmount: expectedEarnedAmount,
		TotalSpentAmount:  expectedSpentAmount,
		AccountsStats: []*forms.AccountStats{
			&firstAccountStats,
			&secondAccountStats,
		},
	}

	TestTotalAccountsStatsEqual(&expected, totalStats, &suite.Suite)
}

func (suite *StatsRepositorySuit) TestGetTotalCategoriesStats() {
	totalStats, err := repository.GetTotalCategoriesStats(
		suite.base.User.ID,
		fromDate,
		toDate,
		suite.base.FirstCurrency,
	)
	suite.True(err == nil)

	expectedSpentAmount := transaction1Amount.Add(transaction2Amount)
	expectedEarnedAmount := transaction3Amount.
		Add(transaction4Amount).
		Mul(secondCurrencyValue)

	firstCategoryStats := forms.CategoryStats{
		Category:    *suite.base.FirstCategory,
		TotalAmount: expectedSpentAmount,
		Transactions: []*models.Transaction{
			suite.base.Transaction1,
			suite.base.Transaction2,
		},
	}
	secondCategoryStats := forms.CategoryStats{
		Category:    *suite.base.SecondCategory,
		TotalAmount: expectedEarnedAmount,
		Transactions: []*models.Transaction{
			suite.base.Transaction3,
			suite.base.Transaction4,
		},
	}

	expected := forms.TotalCategoriesStats{
		TotalEarnedAmount: expectedEarnedAmount,
		TotalSpentAmount:  expectedSpentAmount,
		CategoriesStats: []*forms.CategoryStats{
			&firstCategoryStats,
			&secondCategoryStats,
		},
	}

	TestTotalCategoriesStatsEqual(&expected, totalStats, &suite.Suite)
}
