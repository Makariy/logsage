package repository_tests

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
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

func TestUsersEqual(expected, actual *models.User, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Email, actual.Email)
}

func TestCategoriesEqual(expected, actual *models.Category, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Name, actual.Name)
	suite.Equal(expected.Type, actual.Type)
}

func TestCurrencyEqual(expected, actual *models.Currency, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Name, actual.Name)
}

func TestAccountsEqual(expected, actual *models.Account, suite *suite.Suite) {
	suite.Equal(expected.Name, actual.Name)
	compareDecimal(suite, expected.Balance, actual.Balance)
	suite.Equal(expected.Name, actual.Name)

	TestCurrencyEqual(&expected.Currency, &actual.Currency, suite)
	TestUsersEqual(&expected.User, &actual.User, suite)
}

func TestTransactionsEqual(expected, actual *models.Transaction, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Description, actual.Description)
	suite.True(expected.Date.Truncate(time.Second).Equal(actual.Date.Truncate(time.Second)))
	compareDecimal(suite, expected.Amount, actual.Amount)

	TestUsersEqual(&expected.User, &actual.User, suite)
	TestCategoriesEqual(&expected.Category, &actual.Category, suite)
	TestAccountsEqual(&expected.Account, &actual.Account, suite)
}

func TestCategoriesStatsEqual(expected, actual *models.CategoryStats, suite *suite.Suite) {
	compareDecimal(suite, expected.TotalAmount, actual.TotalAmount)

	TestCategoriesEqual(&expected.Category, &actual.Category, suite)

	suite.Equal(len(expected.Transactions), len(actual.Transactions))
	for i := range expected.Transactions {
		TestTransactionsEqual(expected.Transactions[i], actual.Transactions[i], suite)
	}
}

func TestAccountStatsEqual(expected, actual *models.AccountStats, suite *suite.Suite) {
	compareDecimal(suite, expected.TotalSpentAmount, actual.TotalSpentAmount)
	compareDecimal(suite, expected.TotalEarnedAmount, actual.TotalEarnedAmount)

	TestAccountsEqual(&expected.Account, &actual.Account, suite)

	suite.Equal(len(expected.Transactions), len(actual.Transactions))
	for i := range expected.Transactions {
		TestTransactionsEqual(expected.Transactions[i], actual.Transactions[i], suite)
	}
}

func TestTotalCategoriesStatsEqual(expected, actual *models.TotalCategoriesStats, suite *suite.Suite) {
	compareDecimal(suite, expected.TotalEarnedAmount, actual.TotalEarnedAmount)
	compareDecimal(suite, expected.TotalSpentAmount, actual.TotalSpentAmount)

	suite.Equal(len(expected.CategoriesStats), len(actual.CategoriesStats))
	for i := range expected.CategoriesStats {
		TestCategoriesStatsEqual(expected.CategoriesStats[i], actual.CategoriesStats[i], suite)
	}
}

func TestTotalAccountsStatsEqual(expected, actual *models.TotalAccountsStats, suite *suite.Suite) {
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
		TestAccountStatsEqual(expected.AccountsStats[i], actual.AccountsStats[i], suite)
	}
}
