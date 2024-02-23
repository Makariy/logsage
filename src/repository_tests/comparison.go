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
}
