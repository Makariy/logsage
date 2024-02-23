package repository_tests

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"main/db_connector"
	"main/models"
	"main/repository"
	"main/test_utils"
)

type AccountRepositorySuit struct {
	suite.Suite
	router   *gin.Engine
	user     *models.User
	currency *models.Currency
}

func (suite *AccountRepositorySuit) SetupTest() {
	test_utils.CreateTestDB()
	models.MigrateModels(db_connector.GetConnection())

	suite.user = CreateTestUser(userEmail, userPassword)
	suite.currency = CreateTestCurrency(currencyName)
}

func (suite *AccountRepositorySuit) TearDownTest() {
	test_utils.DropTestDB()
}

func createTestAccount(suite *AccountRepositorySuit) *models.Account {
	account, err := repository.CreateAccount(
		suite.user.ID,
		accountName,
		accountBalance,
		suite.currency.ID,
	)
	if err != nil {
		suite.Error(err, "could not create account")
	}
	return account
}

func (suite *AccountRepositorySuit) TestCreateAccount() {
	account := createTestAccount(suite)

	expected := models.Account{
		ID:       1,
		Name:     accountName,
		Balance:  accountBalance,
		Currency: *suite.currency,
		User:     *suite.user,
	}

	testAccountsEqual(&expected, account, &suite.Suite)
}

func (suite *AccountRepositorySuit) TestPatchAccount() {
	var (
		newAccountName = "New test account name"
		newBalance     = decimal.New(2000, 10)
	)
	account := createTestAccount(suite)
	patched, err := repository.PatchAccount(account.ID, newAccountName, newBalance, account.CurrencyID, suite.user.ID)
	if err != nil {
		suite.Error(err)
	}

	expected := models.Account{
		ID:       account.ID,
		Name:     newAccountName,
		Balance:  newBalance,
		Currency: *suite.currency,
		User:     *suite.user,
	}
	testAccountsEqual(&expected, patched, &suite.Suite)
}

func (suite *AccountRepositorySuit) TestGetAccountByID() {
	account := createTestAccount(suite)

	foundAccount, err := repository.GetAccountByID(account.ID)
	if err != nil {
		suite.Error(err)
	}

	testAccountsEqual(account, foundAccount, &suite.Suite)
}

func (suite *AccountRepositorySuit) TestGetUserAccounts() {
	secondBalance := decimal.New(200, 10)

	first := createTestAccount(suite)
	second, err := repository.CreateAccount(
		suite.user.ID,
		"Second account",
		secondBalance,
		suite.currency.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	accounts, err := repository.GetUserAccounts(suite.user.ID)
	if err != nil {
		suite.Error(err)
	}

	suite.Equal(len(accounts), 2)

	for _, account := range accounts {
		testUsersEqual(suite.user, &account.User, &suite.Suite)
	}

	isFirstFirst := accounts[0].ID == first.ID
	if isFirstFirst {
		testAccountsEqual(first, accounts[0], &suite.Suite)
		testAccountsEqual(second, accounts[1], &suite.Suite)
	} else {
		testAccountsEqual(first, accounts[1], &suite.Suite)
		testAccountsEqual(second, accounts[0], &suite.Suite)
	}
}

func (suite *AccountRepositorySuit) TestDeleteAccount() {
	account := createTestAccount(suite)

	result, err := repository.DeleteAccount(account.ID)
	if err != nil {
		suite.Error(err)
	}
	testAccountsEqual(account, result, &suite.Suite)

	accounts, err := repository.GetUserAccounts(account.User.ID)
	if err != nil {
		suite.Error(err)
	}
	suite.Equal(0, len(accounts))
}
