package repository_tests

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"main/models"
	"main/repository"
	"main/test_utils"
)

type AccountRepositorySuit struct {
	suite.Suite
	base test_utils.RepositoryDefaultSuite
}

func (suite *AccountRepositorySuit) SetupTest() {
	suite.base.SetupDB()
	suite.base.CreateDefaultUser()
	suite.base.CreateDefaultCurrencies()
}

func (suite *AccountRepositorySuit) TearDownTest() {
	suite.base.TearDownTest()
}

func (suite *AccountRepositorySuit) TestCreateAccount() {
	account := test_utils.CreateTestAccount(
		accountName,
		accountBalance,
		suite.base.User.ID,
		suite.base.FirstCurrency.ID,
	)

	expected := models.Account{
		ID:       1,
		Name:     accountName,
		Balance:  accountBalance,
		Currency: *suite.base.FirstCurrency,
		User:     *suite.base.User,
	}

	TestAccountsEqual(&suite.Suite, &expected, account)
}

func (suite *AccountRepositorySuit) TestPatchAccount() {
	var (
		newAccountName = "New test account name"
		newBalance     = decimal.NewFromInt(2000)
	)
	account := test_utils.CreateTestAccount(
		accountName,
		accountBalance,
		suite.base.User.ID,
		suite.base.FirstCurrency.ID,
	)
	patched, err := repository.PatchAccount(
		account.ID,
		newAccountName,
		newBalance,
		account.CurrencyID,
		suite.base.User.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	expected := models.Account{
		ID:       account.ID,
		Name:     newAccountName,
		Balance:  newBalance,
		Currency: *suite.base.FirstCurrency,
		User:     *suite.base.User,
	}
	TestAccountsEqual(&suite.Suite, &expected, patched)
}

func (suite *AccountRepositorySuit) TestGetAccountByID() {
	account := test_utils.CreateTestAccount(accountName, accountBalance, suite.base.User.ID, suite.base.FirstCurrency.ID)

	foundAccount, err := repository.GetAccountByID(account.ID)
	if err != nil {
		suite.Error(err)
	}

	TestAccountsEqual(&suite.Suite, account, foundAccount)
}

func (suite *AccountRepositorySuit) TestGetUserAccounts() {
	secondBalance := decimal.NewFromInt(200)

	first := test_utils.CreateTestAccount(accountName, accountBalance, suite.base.User.ID, suite.base.FirstCurrency.ID)
	second, err := repository.CreateAccount(
		suite.base.User.ID,
		"Second account",
		secondBalance,
		suite.base.FirstCurrency.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	accounts, err := repository.GetUserAccounts(suite.base.User.ID)
	if err != nil {
		suite.Error(err)
	}

	suite.Equal(len(accounts), 2)

	for _, account := range accounts {
		TestUsersEqual(&suite.Suite, suite.base.User, &account.User)
	}

	isFirstFirst := accounts[0].ID == first.ID
	if isFirstFirst {
		TestAccountsEqual(&suite.Suite, first, accounts[0])
		TestAccountsEqual(&suite.Suite, second, accounts[1])
	} else {
		TestAccountsEqual(&suite.Suite, first, accounts[1])
		TestAccountsEqual(&suite.Suite, second, accounts[0])
	}
}

func (suite *AccountRepositorySuit) TestDeleteAccount() {
	account := test_utils.CreateTestAccount(accountName, accountBalance, suite.base.User.ID, suite.base.FirstCurrency.ID)

	result, err := repository.DeleteAccount(account.ID)
	if err != nil {
		suite.Error(err)
	}
	TestAccountsEqual(&suite.Suite, account, result)

	accounts, err := repository.GetUserAccounts(account.User.ID)
	if err != nil {
		suite.Error(err)
	}
	suite.Equal(0, len(accounts))
}
