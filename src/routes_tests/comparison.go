package routes_tests

import (
	"github.com/stretchr/testify/suite"
	"main/forms"
	"main/repository_tests"
	"time"
)

func TestCurrenciesEqual(suite *suite.Suite, expected, actual *forms.CurrencyResponse) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Name, actual.Name)
	suite.Equal(expected.Symbol, actual.Symbol)
}

func TestCategoriesEqual(suite *suite.Suite, expected, actual *forms.CategoryResponse) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Name, actual.Name)
	suite.Equal(expected.Type, actual.Type)
}

func TestAccountsEqual(suite *suite.Suite, expected, actual *forms.AccountResponse) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Name, actual.Name)
	suite.True(expected.Balance.Equal(actual.Balance))

	TestCurrenciesEqual(suite, expected.Currency, actual.Currency)
}

func TestTransactionsEqual(suite *suite.Suite, expected, actual *forms.TransactionResponse) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Description, actual.Description)
	suite.True(expected.Amount.Equal(actual.Amount))
	suite.True(expected.Date.Truncate(time.Second).Equal(actual.Date.Truncate(time.Second)))

	TestCategoriesEqual(suite, &expected.Category, &actual.Category)
	TestAccountsEqual(suite, &expected.Account, &actual.Account)
}

func TestCategoryStatsEqual(suite *suite.Suite, expected, actual *forms.CategoryStatsResponse) {
	TestCategoriesEqual(suite, &expected.Category, &actual.Category)
	repository_tests.TestCategoriesStatsEqual(suite, &expected.Stats, &actual.Stats)
}

func TestAccountStatsEqual(suite *suite.Suite, expected, actual *forms.AccountStatsResponse) {
	TestAccountsEqual(suite, &expected.Account, &actual.Account)
	repository_tests.TestAccountStatsEqual(suite, &expected.Stats, &actual.Stats)
}

func TestTotalAccountsStatsEqual(suite *suite.Suite, expected, actual *forms.TotalAccountsStatsResponse) {
	repository_tests.TestTotalAccountsStatsEqual(suite, &expected.Stats, &actual.Stats)
}

func TestTotalCategoriesStatsEqual(suite *suite.Suite, expected, actual *forms.TotalCategoriesStatsResponse) {
	repository_tests.TestTotalCategoriesStatsEqual(suite, &expected.Stats, &actual.Stats)
}
