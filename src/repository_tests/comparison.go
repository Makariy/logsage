package repository_tests

import (
	"github.com/stretchr/testify/suite"
	"main/models"
	"time"
)

func testUsersEqual(expected, actual *models.User, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Email, actual.Email)
}

func testCategoriesEqual(expected, actual *models.Category, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Name, actual.Name)
	suite.Equal(expected.Type, actual.Type)
}

func testCurrencyEqual(expected, actual *models.Currency, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Name, actual.Name)
}

func testAccountsEqual(expected, actual *models.Account, suite *suite.Suite) {
	suite.Equal(expected.Name, actual.Name)
	suite.True(expected.Balance.Equal(actual.Balance))
	suite.Equal(expected.Name, actual.Name)

	testCurrencyEqual(&expected.Currency, &actual.Currency, suite)
	testUsersEqual(&expected.User, &actual.User, suite)
}

func testTransactionsEqual(expected, actual *models.Transaction, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Description, actual.Description)
	suite.True(expected.Date.Truncate(time.Second).Equal(actual.Date.Truncate(time.Second)))
	suite.True(expected.Amount.Equal(actual.Amount))

	testUsersEqual(&expected.User, &actual.User, suite)
	testCategoriesEqual(&expected.Category, &actual.Category, suite)
	testAccountsEqual(&expected.Account, &actual.Account, suite)
}

func testCategoriesStatsEqual(expected, actual *models.CategoryStats, suite *suite.Suite) {
	suite.True(expected.TotalAmount.Equal(actual.TotalAmount))
	suite.True(expected.TotalPercent.Equal(actual.TotalPercent))

	testCategoriesEqual(&expected.Category, &actual.Category, suite)

	suite.Equal(len(expected.Transactions), len(actual.Transactions))
	for i := range expected.Transactions {
		testTransactionsEqual(expected.Transactions[i], actual.Transactions[i], suite)
	}
}

func testAccountStatsEqual(expected, actual *models.AccountStats, suite *suite.Suite) {
	suite.True(expected.TotalSpentAmount.Equal(actual.TotalSpentAmount))
	suite.True(expected.TotalSpentPercent.Equal(actual.TotalSpentPercent))
	suite.True(expected.TotalEarnedAmount.Equal(actual.TotalEarnedAmount))
	suite.True(expected.TotalEarnedPercent.Equal(actual.TotalEarnedPercent))

	testAccountsEqual(&expected.Account, &actual.Account, suite)

	suite.Equal(len(expected.Transactions), len(actual.Transactions))
	for i := range expected.Transactions {
		testTransactionsEqual(expected.Transactions[i], actual.Transactions[i], suite)
	}
}

func testTotalStatsEqual(expected, actual *models.TotalStats, suite *suite.Suite) {
	suite.True(expected.TotalEarnedAmount.Equal(actual.TotalEarnedAmount))
	suite.True(expected.TotalSpentAmount.Equal(actual.TotalSpentAmount))

	suite.Equal(len(expected.AccountsStats), len(actual.AccountsStats))
	for i := range expected.AccountsStats {
		testAccountStatsEqual(expected.AccountsStats[i], actual.AccountsStats[i], suite)
	}

	suite.Equal(len(expected.CategoriesStats), len(actual.CategoriesStats))
	for i := range expected.CategoriesStats {
		testCategoriesStatsEqual(expected.CategoriesStats[i], actual.CategoriesStats[i], suite)
	}
}
