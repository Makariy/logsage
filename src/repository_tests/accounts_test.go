package repository_tests

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"main/models"
	"main/repository"
)

type AccountRepositorySuit struct {
	suite.Suite
	base DefaultRepositorySuite
}

func (suite *AccountRepositorySuit) SetupTest() {
	suite.base.setupDB()
	suite.base.createDefaultUser()
	suite.base.createDefaultCurrencies()
}

func (suite *AccountRepositorySuit) TearDownTest() {
	suite.base.TearDownTest()
}

func (suite *AccountRepositorySuit) TestCreateAccount() {
	account := CreateTestAccount(
		accountName,
		accountBalance,
		suite.base.user.ID,
		suite.base.firstCurrency.ID,
	)

	expected := models.Account{
		ID:       1,
		Name:     accountName,
		Balance:  accountBalance,
		Currency: *suite.base.firstCurrency,
		User:     *suite.base.user,
	}

	TestAccountsEqual(&expected, account, &suite.Suite)
}

func (suite *AccountRepositorySuit) TestPatchAccount() {
	var (
		newAccountName = "New test account name"
		newBalance     = decimal.NewFromInt(2000)
	)
	account := CreateTestAccount(
		accountName,
		accountBalance,
		suite.base.user.ID,
		suite.base.firstCurrency.ID,
	)
	patched, err := repository.PatchAccount(
		account.ID,
		newAccountName,
		newBalance,
		account.CurrencyID,
		suite.base.user.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	expected := models.Account{
		ID:       account.ID,
		Name:     newAccountName,
		Balance:  newBalance,
		Currency: *suite.base.firstCurrency,
		User:     *suite.base.user,
	}
	TestAccountsEqual(&expected, patched, &suite.Suite)
}

func (suite *AccountRepositorySuit) TestGetAccountByID() {
	account := CreateTestAccount(accountName, accountBalance, suite.base.user.ID, suite.base.firstCurrency.ID)

	foundAccount, err := repository.GetAccountByID(account.ID)
	if err != nil {
		suite.Error(err)
	}

	TestAccountsEqual(account, foundAccount, &suite.Suite)
}

func (suite *AccountRepositorySuit) TestGetUserAccounts() {
	secondBalance := decimal.NewFromInt(200)

	first := CreateTestAccount(accountName, accountBalance, suite.base.user.ID, suite.base.firstCurrency.ID)
	second, err := repository.CreateAccount(
		suite.base.user.ID,
		"Second account",
		secondBalance,
		suite.base.firstCurrency.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	accounts, err := repository.GetUserAccounts(suite.base.user.ID)
	if err != nil {
		suite.Error(err)
	}

	suite.Equal(len(accounts), 2)

	for _, account := range accounts {
		TestUsersEqual(suite.base.user, &account.User, &suite.Suite)
	}

	isFirstFirst := accounts[0].ID == first.ID
	if isFirstFirst {
		TestAccountsEqual(first, accounts[0], &suite.Suite)
		TestAccountsEqual(second, accounts[1], &suite.Suite)
	} else {
		TestAccountsEqual(first, accounts[1], &suite.Suite)
		TestAccountsEqual(second, accounts[0], &suite.Suite)
	}
}

func (suite *AccountRepositorySuit) TestDeleteAccount() {
	account := CreateTestAccount(accountName, accountBalance, suite.base.user.ID, suite.base.firstCurrency.ID)

	result, err := repository.DeleteAccount(account.ID)
	if err != nil {
		suite.Error(err)
	}
	TestAccountsEqual(account, result, &suite.Suite)

	accounts, err := repository.GetUserAccounts(account.User.ID)
	if err != nil {
		suite.Error(err)
	}
	suite.Equal(0, len(accounts))
}
