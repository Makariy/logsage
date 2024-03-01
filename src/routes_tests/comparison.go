package routes_tests

import (
	"github.com/stretchr/testify/suite"
	"main/forms"
	"main/repository_tests"
	"time"
)

func TestCurrenciesEqual(expected, actual *forms.CurrencyResponse, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Name, actual.Name)
}

func TestCategoriesEqual(expected, actual *forms.CategoryResponse, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Name, actual.Name)
	suite.Equal(expected.Type, actual.Type)
}

func TestAccountsEqual(expected, actual *forms.AccountResponse, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Name, actual.Name)
	suite.True(expected.Balance.Equal(actual.Balance))

	TestCurrenciesEqual(expected.Currency, actual.Currency, suite)
}

func TestTransactionsEqual(expected, actual *forms.TransactionResponse, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Description, actual.Description)
	suite.True(expected.Amount.Equal(actual.Amount))
	suite.True(expected.Date.Truncate(time.Second).Equal(actual.Date.Truncate(time.Second)))

	TestCategoriesEqual(&expected.Category, &actual.Category, suite)
}

func TestCategoryStatsEqual(expected, actual *forms.CategoryStatsResponse, suite *suite.Suite) {
	TestCategoriesEqual(&expected.Category, &actual.Category, suite)
	repository_tests.TestCategoriesStatsEqual(&expected.Stats, &actual.Stats, suite)
}

func TestAccountStatsEqual(expected, actual *forms.AccountStatsResponse, suite *suite.Suite) {
	TestAccountsEqual(&expected.Account, &actual.Account, suite)
	repository_tests.TestAccountStatsEqual(&expected.Stats, &actual.Stats, suite)
}

func TestTotalStatsEqual(expected, actual *forms.TotalStatsResponse, suite *suite.Suite) {
	repository_tests.TestTotalStatsEqual(&expected.Stats, &actual.Stats, suite)
}
