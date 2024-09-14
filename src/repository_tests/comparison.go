package repository_tests

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"main/forms"
	"main/models"
	"time"
)

func compareDecimal(suite *suite.Suite, expected, actual decimal.Decimal, message ...string) {
	errMsg := "Expected = " + expected.String() + " != " + actual.String() + " = Actual"

	for _, msg := range message {
		errMsg += "\n" + msg
	}

	suite.True(
		expected.Equal(actual),
		errMsg,
	)
}

func TestDateRangeEqual(suite *suite.Suite, expected, actual *forms.DateRange) {
	suite.Equal(expected.FromDate, actual.FromDate)
	suite.Equal(expected.ToDate, actual.ToDate)
}

func TestUsersEqual(suite *suite.Suite, expected, actual *models.User) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Email, actual.Email)
}

func TestCategoriesEqual(suite *suite.Suite, expected, actual *models.Category) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Name, actual.Name)
	suite.Equal(expected.Type, actual.Type)
}

func TestCurrencyEqual(suite *suite.Suite, expected, actual *models.Currency) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Name, actual.Name)
}

func TestAccountsEqual(suite *suite.Suite, expected, actual *models.Account) {
	suite.Equal(expected.Name, actual.Name)
	compareDecimal(suite, expected.Balance, actual.Balance)
	suite.Equal(expected.Name, actual.Name)

	TestCurrencyEqual(suite, &expected.Currency, &actual.Currency)
	TestUsersEqual(suite, &expected.User, &actual.User)
}

func TestTransactionsEqual(suite *suite.Suite, expected, actual *models.Transaction) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Description, actual.Description)
	suite.True(expected.Date.Truncate(time.Second).Equal(actual.Date.Truncate(time.Second)))
	compareDecimal(suite, expected.Amount, actual.Amount)

	TestUsersEqual(suite, &expected.User, &actual.User)
	TestCategoriesEqual(suite, &expected.Category, &actual.Category)
	TestAccountsEqual(suite, &expected.Account, &actual.Account)
}

func TestCategoriesStatsEqual(suite *suite.Suite, expected, actual *forms.CategoryStats) {
	compareDecimal(suite, expected.TotalAmount, actual.TotalAmount)

	TestCategoriesEqual(suite, &expected.Category, &actual.Category)

	suite.Equal(len(expected.Transactions), len(actual.Transactions))
	for i := range expected.Transactions {
		TestTransactionsEqual(suite, expected.Transactions[i], actual.Transactions[i])
	}
}

func TestAccountStatsEqual(suite *suite.Suite, expected, actual *forms.AccountStats) {
	compareDecimal(suite, expected.TotalSpentAmount, actual.TotalSpentAmount)
	compareDecimal(suite, expected.TotalEarnedAmount, actual.TotalEarnedAmount)

	TestAccountsEqual(suite, &expected.Account, &actual.Account)

	suite.Equal(len(expected.Transactions), len(actual.Transactions))
	for i := range expected.Transactions {
		TestTransactionsEqual(suite, expected.Transactions[i], actual.Transactions[i])
	}
}

func TestTotalCategoriesStatsEqual(suite *suite.Suite, expected, actual *forms.TotalCategoriesStats) {
	compareDecimal(suite, expected.TotalEarnedAmount, actual.TotalEarnedAmount)
	compareDecimal(suite, expected.TotalSpentAmount, actual.TotalSpentAmount)

	suite.Equal(len(expected.CategoriesStats), len(actual.CategoriesStats))
	for i := range expected.CategoriesStats {
		TestCategoriesStatsEqual(suite, expected.CategoriesStats[i], actual.CategoriesStats[i])
	}
}

func TestTotalAccountsStatsEqual(suite *suite.Suite, expected, actual *forms.TotalAccountsStats) {
	compareDecimal(
		suite,
		expected.TotalEarnedAmount,
		actual.TotalEarnedAmount,
		"Total EARNED amount do not match",
	)
	compareDecimal(
		suite,
		expected.TotalSpentAmount,
		actual.TotalSpentAmount,
		"Total SPENT amount do not match",
	)

	suite.Equal(len(expected.AccountsStats), len(actual.AccountsStats))
	for i := range expected.AccountsStats {
		TestAccountStatsEqual(suite, expected.AccountsStats[i], actual.AccountsStats[i])
	}
}

func TestTimeIntervalStatEqual(suite *suite.Suite, expected, actual *forms.TimeIntervalStat) {
	compareDecimal(suite, expected.TotalEarnedAmount, actual.TotalEarnedAmount)
	compareDecimal(suite, expected.TotalSpentAmount, actual.TotalSpentAmount)
	TestDateRangeEqual(suite, expected.DateRange, actual.DateRange)
}

func TestTimeIntervalStatsEqual(suite *suite.Suite, expected, actual *forms.TimeIntervalStats) {
	suite.Equal(expected.TimeStep, actual.TimeStep, "interval stats time step do not match")
	TestDateRangeEqual(suite, expected.DateRange, actual.DateRange)

	suite.Equal(len(expected.IntervalStats), len(actual.IntervalStats))

	for index := range expected.IntervalStats {
		TestTimeIntervalStatEqual(suite, expected.IntervalStats[index], actual.IntervalStats[index])
	}
}
