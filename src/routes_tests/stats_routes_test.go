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
	"main/utils"
	"time"
)

type StatsRoutesSuit struct {
	suite.Suite
	router *gin.Engine

	user          *models.User
	currency      *models.Currency
	category      *models.Category
	otherCategory *models.Category
	account       *models.Account
	authHeaders   map[string]string
}

func (suite *StatsRoutesSuit) SetupTest() {
	test_utils.CreateTestDB()
	models.MigrateModels(db_connector.GetConnection())

	suite.user = CreateTestUser(userEmail, userPassword)
	suite.currency = CreateTestCurrency(currencyName)
	suite.category = CreateTestCategory(categoryName, categoryType, suite.user.ID)
	suite.otherCategory = CreateTestCategory("Other category", models.EARNING, suite.user.ID)
	suite.account = CreateTestAccount(accountName, accountBalance, suite.user.ID, suite.currency.ID)
	suite.authHeaders = GetAuthHeaders(suite.user)

	createTestTransactions(suite)

	suite.router = gin.Default()
	routes.AddStatsRoutes(suite.router)
}

func (suite *StatsRoutesSuit) TearDownTest() {
	test_utils.DropTestDB()
}

var (
	transaction1, transaction2, transaction3, transaction4 *models.Transaction
)

var (
	year     = 2024
	month    = time.January
	fromDate = time.Date(year, month, 01, 00, 00, 00, 00, time.Local)
	toDate   = time.Date(year, month, 10, 00, 00, 00, 00, time.Local)
)

func createTestTransactions(suite *StatsRoutesSuit) {
	transaction1, _ = repository.CreateTransaction(
		"First transaction",
		decimal.NewFromInt(200),
		time.Date(year, month, 01, 00, 00, 00, 00, time.Local),
		suite.user.ID,
		suite.category.ID,
		suite.account.ID,
	)
	transaction2, _ = repository.CreateTransaction(
		"Second transaction",
		decimal.NewFromInt(200),
		time.Date(year, month, 02, 00, 00, 00, 00, time.Local),
		suite.user.ID,
		suite.category.ID,
		suite.account.ID,
	)
	transaction3, _ = repository.CreateTransaction(
		"Third transaction",
		decimal.NewFromInt(200),
		time.Date(year, month, 03, 00, 00, 00, 00, time.Local),
		suite.user.ID,
		suite.category.ID,
		suite.account.ID,
	)
	transaction4, _ = repository.CreateTransaction(
		"FORTH transaction",
		decimal.NewFromInt(200),
		time.Date(year, month, 04, 00, 00, 00, 00, time.Local),
		suite.user.ID,
		suite.otherCategory.ID,
		suite.account.ID,
	)
}

func getDateRange() []byte {
	form := forms.DateRange{
		FromDate: fromDate,
		ToDate:   toDate,
	}

	data, err := json.Marshal(form)
	if err != nil {
		panic("could not marshal DateRange to JSON")
	}
	return data
}

func (suite *StatsRoutesSuit) TestHandleGetCategoryStats() {
	resp := PerformTestRequest(
		suite.router,
		"GET",
		fmt.Sprintf("/stats/category/%d/", suite.category.ID),
		getDateRange(),
		&suite.authHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	categoryStats, err := UnmarshalResponse[forms.CategoryStatsResponse](resp)
	suite.True(err == nil)

	categoryForm, err := utils.MarshalModelToForm[models.Category, forms.CategoryResponse](suite.category)

	if err != nil {
		panic("could not marshal category to form")
	}

	expected := forms.CategoryStatsResponse{
		Category: *categoryForm,
		Stats: models.CategoryStats{
			Category:     *suite.category,
			TotalAmount:  decimal.NewFromInt(600),
			TotalPercent: decimal.NewFromFloat(0.75),
			Transactions: []*models.Transaction{
				transaction1,
				transaction2,
				transaction3,
			},
		},
	}

	TestCategoryStatsEqual(&expected, categoryStats, &suite.Suite)
}

func (suite *StatsRoutesSuit) TestHandleGetAccountStats() {
	resp := PerformTestRequest(
		suite.router,
		"GET",
		fmt.Sprintf("/stats/account/%d/", suite.account.ID),
		getDateRange(),
		&suite.authHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	accountStats, err := UnmarshalResponse[forms.AccountStatsResponse](resp)
	suite.True(err == nil)

	accountForm, err := utils.MarshalModelToForm[models.Account, forms.AccountResponse](suite.account)
	if err != nil {
		panic("could not marshall account to form")
	}

	expected := forms.AccountStatsResponse{
		Account: *accountForm,
		Stats: models.AccountStats{
			Account:            *suite.account,
			TotalEarnedAmount:  decimal.NewFromInt(200),
			TotalSpentAmount:   decimal.NewFromInt(600),
			TotalEarnedPercent: decimal.NewFromFloat(1),
			TotalSpentPercent:  decimal.NewFromFloat(1),
			Transactions: []*models.Transaction{
				transaction1,
				transaction2,
				transaction3,
				transaction4,
			},
		},
	}

	TestAccountStatsEqual(&expected, accountStats, &suite.Suite)
}

func (suite *StatsRoutesSuit) TestHandleGetTotalStats() {
	resp := PerformTestRequest(
		suite.router,
		"GET",
		"/stats/all/",
		getDateRange(),
		&suite.authHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	stats, err := UnmarshalResponse[forms.TotalStatsResponse](resp)
	suite.True(err == nil)

	expected := forms.TotalStatsResponse{
		Stats: models.TotalStats{
			TotalEarnedAmount: decimal.NewFromInt(200),
			TotalSpentAmount:  decimal.NewFromInt(600),
			AccountsStats: []*models.AccountStats{
				{
					Account:            *suite.account,
					TotalEarnedAmount:  decimal.NewFromInt(200),
					TotalSpentAmount:   decimal.NewFromInt(600),
					TotalEarnedPercent: decimal.NewFromFloat(1),
					TotalSpentPercent:  decimal.NewFromFloat(1),
					Transactions: []*models.Transaction{
						transaction1,
						transaction2,
						transaction3,
						transaction4,
					},
				},
			},
			CategoriesStats: []*models.CategoryStats{
				{
					Category:     *suite.category,
					TotalAmount:  decimal.NewFromInt(600),
					TotalPercent: decimal.NewFromFloat(0.75),
					Transactions: []*models.Transaction{
						transaction1,
						transaction2,
						transaction3,
					},
				},
				{
					Category:     *suite.otherCategory,
					TotalAmount:  decimal.NewFromInt(200),
					TotalPercent: decimal.NewFromFloat(0.25),
					Transactions: []*models.Transaction{
						transaction4,
					},
				},
			},
		},
	}

	TestTotalStatsEqual(&expected, stats, &suite.Suite)
}
