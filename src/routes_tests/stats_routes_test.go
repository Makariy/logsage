package routes_tests

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"main/forms"
	"main/models"
	"main/test_utils"
	"main/utils"
)

type StatsRoutesSuit struct {
	suite.Suite
	router test_utils.RoutesDefaultSuite
}

func (suite *StatsRoutesSuit) SetupTest() {
	suite.router.SetupAllTestData()
}

func (suite *StatsRoutesSuit) TearDownTest() {
	suite.router.TearDownTest()
}

func getExpectedDateRangeForTransactions(
	startTransaction *models.Transaction,
	stopTransaction *models.Transaction,
) *forms.DateRange {
	return &forms.DateRange{
		FromDate: startTransaction.Date.Unix(),
		ToDate:   stopTransaction.Date.Unix(),
	}
}

func (suite *StatsRoutesSuit) TestHandleGetCategoryStats() {
	startTransaction := suite.router.Data.Transaction1
	stopTransaction := suite.router.Data.Transaction4

	resp := PerformTestRequest(
		suite.router.Router,
		"GET",
		fmt.Sprintf(
			"/stats/category/%d/?currency=%s&%s",
			suite.router.Data.FirstCategory.ID,
			suite.router.Data.FirstCurrency.Symbol,
			renderDateRangeForTransactions(
				startTransaction,
				stopTransaction,
			),
		),
		nil,
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	categoryStats, err := UnmarshalResponse[forms.CategoryStatsResponse](resp)
	suite.True(err == nil)

	categoryForm, err := utils.MarshalModelToForm[models.Category, forms.CategoryResponse](
		suite.router.Data.FirstCategory,
	)

	if err != nil {
		panic("could not marshal category to form")
	}

	expected := forms.CategoryStatsResponse{
		Category: *categoryForm,
		DateRange: getExpectedDateRangeForTransactions(
			startTransaction,
			stopTransaction,
		),
		Stats: forms.CategoryStats{
			Category:    *suite.router.Data.FirstCategory,
			TotalAmount: decimal.NewFromInt(300),
			Transactions: []*models.Transaction{
				suite.router.Data.Transaction1,
				suite.router.Data.Transaction2,
			},
		},
	}

	TestCategoryStatsEqual(&suite.Suite, &expected, categoryStats)
}

func (suite *StatsRoutesSuit) TestHandleGetAccountStats() {
	startTransaction := suite.router.Data.Transaction1
	stopTransaction := suite.router.Data.Transaction4

	resp := PerformTestRequest(
		suite.router.Router,
		"GET",
		fmt.Sprintf(
			"/stats/account/%d/?%s",
			suite.router.Data.FirstAccount.ID,
			renderDateRangeForTransactions(
				startTransaction,
				stopTransaction,
			),
		),
		nil,
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	accountStats, err := UnmarshalResponse[forms.AccountStatsResponse](resp)
	suite.True(err == nil)

	accountForm, err := utils.MarshalModelToForm[models.Account, forms.AccountResponse](
		suite.router.Data.FirstAccount,
	)
	if err != nil {
		panic("could not marshall account to form")
	}

	expected := forms.AccountStatsResponse{
		Account: *accountForm,
		DateRange: getExpectedDateRangeForTransactions(
			startTransaction,
			stopTransaction,
		),
		Stats: forms.AccountStats{
			Account:           *suite.router.Data.FirstAccount,
			TotalSpentAmount:  decimal.NewFromInt(300),
			TotalEarnedAmount: decimal.Zero,
			Transactions: []*models.Transaction{
				suite.router.Data.Transaction1,
				suite.router.Data.Transaction2,
			},
		},
	}

	TestAccountStatsEqual(&suite.Suite, &expected, accountStats)
}

func (suite *StatsRoutesSuit) TestHandleGetTotalAccountsStats() {
	startTransaction := suite.router.Data.Transaction1
	stopTransaction := suite.router.Data.Transaction4

	resp := PerformTestRequest(
		suite.router.Router,
		"GET",
		fmt.Sprintf(
			"/stats/account/all/?%s",
			renderDateRangeForTransactions(
				startTransaction,
				stopTransaction,
			),
		),
		nil,
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	stats, err := UnmarshalResponse[forms.TotalAccountsStatsResponse](resp)
	suite.True(err == nil)

	expected := forms.TotalAccountsStatsResponse{
		DateRange: getExpectedDateRangeForTransactions(
			startTransaction,
			stopTransaction,
		),
		Stats: forms.TotalAccountsStats{
			TotalSpentAmount:  decimal.NewFromInt(300),
			TotalEarnedAmount: decimal.NewFromInt(7),
			AccountsStats: []*forms.AccountStats{
				{
					Account:           *suite.router.Data.FirstAccount,
					TotalEarnedAmount: decimal.Zero,
					TotalSpentAmount:  decimal.NewFromInt(300),
					Transactions: []*models.Transaction{
						suite.router.Data.Transaction1,
						suite.router.Data.Transaction2,
					},
				},
				{
					Account:           *suite.router.Data.SecondAccount,
					TotalEarnedAmount: decimal.NewFromInt(7),
					TotalSpentAmount:  decimal.NewFromInt(0),
					Transactions: []*models.Transaction{
						suite.router.Data.Transaction3,
						suite.router.Data.Transaction4,
					},
				},
			},
		},
	}

	TestTotalAccountsStatsEqual(&suite.Suite, &expected, stats)
}

func (suite *StatsRoutesSuit) TestHandleGetTotalCategoriesStats() {
	startTransaction := suite.router.Data.Transaction1
	stopTransaction := suite.router.Data.Transaction4

	resp := PerformTestRequest(
		suite.router.Router,
		"GET",
		fmt.Sprintf(
			"/stats/category/all/?currency=%s&%s",
			suite.router.Data.FirstCurrency.Symbol,
			renderDateRangeForTransactions(
				startTransaction,
				stopTransaction,
			),
		),
		nil,
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	stats, err := UnmarshalResponse[forms.TotalCategoriesStatsResponse](resp)
	suite.True(err == nil)

	expected := forms.TotalCategoriesStatsResponse{
		DateRange: getExpectedDateRangeForTransactions(
			startTransaction,
			stopTransaction,
		),
		Stats: forms.TotalCategoriesStats{
			TotalEarnedAmount: decimal.NewFromInt(7),
			TotalSpentAmount:  decimal.NewFromInt(300),
			CategoriesStats: []*forms.CategoryStats{
				{
					Category:    *suite.router.Data.FirstCategory,
					TotalAmount: decimal.NewFromInt(300),
					Transactions: []*models.Transaction{
						suite.router.Data.Transaction1,
						suite.router.Data.Transaction2,
					},
				},
				{
					Category:    *suite.router.Data.SecondCategory,
					TotalAmount: decimal.NewFromInt(7),
					Transactions: []*models.Transaction{
						suite.router.Data.Transaction3,
						suite.router.Data.Transaction4,
					},
				},
			},
		},
	}

	TestTotalCategoriesStatsEqual(&suite.Suite, &expected, stats)
}
