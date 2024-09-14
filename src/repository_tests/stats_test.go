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

	TestCategoriesStatsEqual(&suite.Suite, &expected, categoryStats)
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

	TestAccountStatsEqual(&suite.Suite, &expected, accountStats)
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

	TestTotalAccountsStatsEqual(&suite.Suite, &expected, totalStats)
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

	TestTotalCategoriesStatsEqual(&suite.Suite, &expected, totalStats)
}

func (suite *StatsRepositorySuit) TestGetTimeIntervalStats() {
	stats, err := repository.GetTimeIntervalStats(
		suite.base.User.ID,
		fromDate,
		toDate,
		"day",
		suite.base.FirstCurrency,
	)
	suite.True(err == nil)

	dayDuration := int64(24 * time.Hour / time.Second)
	expected := &forms.TimeIntervalStats{
		TimeStep: int64(24 * time.Hour / time.Second),
		DateRange: &forms.DateRange{
			FromDate: fromDate.Unix(),
			ToDate:   toDate.Unix(),
		},
		IntervalStats: []*forms.TimeIntervalStat{
			{
				TotalSpentAmount:  decimal.NewFromInt(100),
				TotalEarnedAmount: decimal.Zero,
				DateRange: &forms.DateRange{
					FromDate: fromDate.Unix(),
					ToDate:   fromDate.Unix() + dayDuration,
				},
			},
			{
				TotalSpentAmount:  decimal.NewFromInt(200),
				TotalEarnedAmount: decimal.Zero,
				DateRange: &forms.DateRange{
					FromDate: fromDate.Unix() + dayDuration,
					ToDate:   fromDate.Unix() + dayDuration*2,
				},
			},
			{
				TotalSpentAmount:  decimal.Zero,
				TotalEarnedAmount: decimal.NewFromInt(3),
				DateRange: &forms.DateRange{
					FromDate: fromDate.Unix() + dayDuration*2,
					ToDate:   fromDate.Unix() + dayDuration*3,
				},
			},
			{
				TotalSpentAmount:  decimal.Zero,
				TotalEarnedAmount: decimal.NewFromInt(4),
				DateRange: &forms.DateRange{
					FromDate: fromDate.Unix() + dayDuration*3,
					ToDate:   fromDate.Unix() + dayDuration*4,
				},
			},
		},
	}

	TestTimeIntervalStatsEqual(&suite.Suite, expected, stats)
}
