package repository_tests

import (
	"github.com/stretchr/testify/suite"
	"main/models"
	"time"
)

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
	suite.True(expected.Balance.Equal(actual.Balance))
	suite.Equal(expected.Name, actual.Name)

	TestCurrencyEqual(&expected.Currency, &actual.Currency, suite)
	TestUsersEqual(&expected.User, &actual.User, suite)
}

func TestTransactionsEqual(expected, actual *models.Transaction, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Description, actual.Description)
	suite.True(expected.Date.Truncate(time.Second).Equal(actual.Date.Truncate(time.Second)))
	suite.True(expected.Amount.Equal(actual.Amount))

	TestUsersEqual(&expected.User, &actual.User, suite)
	TestCategoriesEqual(&expected.Category, &actual.Category, suite)
	TestAccountsEqual(&expected.Account, &actual.Account, suite)
}

func TestCategoriesStatsEqual(expected, actual *models.CategoryStats, suite *suite.Suite) {
	suite.True(expected.TotalAmount.Equal(actual.TotalAmount))
	suite.True(expected.TotalPercent.Equal(actual.TotalPercent))

	TestCategoriesEqual(&expected.Category, &actual.Category, suite)

	suite.Equal(len(expected.Transactions), len(actual.Transactions))
	for i := range expected.Transactions {
		TestTransactionsEqual(expected.Transactions[i], actual.Transactions[i], suite)
	}
}

func TestAccountStatsEqual(expected, actual *models.AccountStats, suite *suite.Suite) {
	suite.True(expected.TotalSpentAmount.Equal(actual.TotalSpentAmount))
	suite.True(expected.TotalSpentPercent.Equal(actual.TotalSpentPercent))
	suite.True(expected.TotalEarnedAmount.Equal(actual.TotalEarnedAmount))
	suite.True(expected.TotalEarnedPercent.Equal(actual.TotalEarnedPercent))

	TestAccountsEqual(&expected.Account, &actual.Account, suite)

	suite.Equal(len(expected.Transactions), len(actual.Transactions))
	for i := range expected.Transactions {
		TestTransactionsEqual(expected.Transactions[i], actual.Transactions[i], suite)
	}
}

func TestTotalStatsEqual(expected, actual *models.TotalStats, suite *suite.Suite) {
	suite.True(expected.TotalEarnedAmount.Equal(actual.TotalEarnedAmount))
	suite.True(expected.TotalSpentAmount.Equal(actual.TotalSpentAmount))

	suite.Equal(len(expected.AccountsStats), len(actual.AccountsStats))
	for i := range expected.AccountsStats {
		TestAccountStatsEqual(expected.AccountsStats[i], actual.AccountsStats[i], suite)
	}

	suite.Equal(len(expected.CategoriesStats), len(actual.CategoriesStats))
	for i := range expected.CategoriesStats {
		TestCategoriesStatsEqual(expected.CategoriesStats[i], actual.CategoriesStats[i], suite)
	}
}
