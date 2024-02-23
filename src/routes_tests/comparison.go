package routes_tests

import (
	"github.com/stretchr/testify/suite"
	"main/forms"
	"time"
)

func testCurrenciesEqual(expected, actual *forms.CurrencyResponse, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Name, actual.Name)
}

func testCategoriesEqual(expected, actual *forms.CategoryResponse, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Name, actual.Name)
	suite.Equal(expected.Type, actual.Type)
}

func testAccountsEqual(expected, actual *forms.AccountResponse, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Name, actual.Name)
	suite.True(expected.Balance.Equal(actual.Balance))

	testCurrenciesEqual(expected.Currency, actual.Currency, suite)
}

func testTransactionsEqual(expected, actual *forms.TransactionResponse, suite *suite.Suite) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.Description, actual.Description)
	suite.True(expected.Amount.Equal(actual.Amount))
	suite.True(expected.Date.Truncate(time.Second).Equal(actual.Date.Truncate(time.Second)))

	testCategoriesEqual(&expected.Category, &actual.Category, suite)
}
