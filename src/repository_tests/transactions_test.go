package repository_tests

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"main/models"
	"main/repository"
	"main/test_utils"
	"time"
)

type TransactionRepositorySuit struct {
	suite.Suite
	base test_utils.RepositoryDefaultSuite
}

func (suite *TransactionRepositorySuit) SetupTest() {
	suite.base.SetupDB()
	suite.base.CreateDefaultUser()
	suite.base.CreateDefaultCurrencies()
	suite.base.CreateDefaultCategories()
	suite.base.CreateDefaultAccounts()
}

func (suite *TransactionRepositorySuit) TearDownTest() {
	suite.base.TearDownTest()
}

func (suite *TransactionRepositorySuit) TestCreateTransaction() {
	transaction, err := repository.CreateTransaction(
		transactionDescription,
		transactionAmount,
		transactionDate,
		suite.base.User.ID,
		suite.base.FirstCategory.ID,
		suite.base.FirstAccount.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	expected := models.Transaction{
		ID:          1,
		Description: transactionDescription,
		Amount:      transactionAmount,
		Date:        transactionDate,
		User:        *suite.base.User,
		Category:    *suite.base.FirstCategory,
		Account:     *suite.base.FirstAccount,
	}

	TestTransactionsEqual(&suite.Suite, &expected, transaction)
}

func (suite *TransactionRepositorySuit) TestGetTransactionByID() {
	transaction, err := repository.CreateTransaction(
		transactionDescription,
		transactionAmount,
		transactionDate,
		suite.base.User.ID,
		suite.base.FirstCategory.ID,
		suite.base.FirstAccount.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	result, err := repository.GetTransactionByID(transaction.ID)
	if err != nil {
		suite.Error(err)
	}

	expected := models.Transaction{
		ID:          transaction.ID,
		Description: transaction.Description,
		Amount:      transaction.Amount,
		Date:        transaction.Date,
		User:        transaction.User,
		Category:    transaction.Category,
		Account:     transaction.Account,
	}

	TestTransactionsEqual(&suite.Suite, &expected, result)
}

func (suite *TransactionRepositorySuit) TestGetAllTransactions() {
	first, err := repository.CreateTransaction(
		transactionDescription,
		transactionAmount,
		transactionDate,
		suite.base.User.ID,
		suite.base.FirstCategory.ID,
		suite.base.FirstAccount.ID,
	)
	if err != nil {
		suite.Error(err)
	}
	second, err := repository.CreateTransaction(
		"Other transaction",
		decimal.NewFromInt(500),
		time.Now(),
		suite.base.User.ID,
		suite.base.FirstCategory.ID,
		suite.base.FirstAccount.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	transactions, err := repository.GetUserTransactions(suite.base.User.ID)
	if err != nil {
		suite.Error(err)
	}
	suite.Equal(2, len(transactions))

	for _, transaction := range transactions {
		TestUsersEqual(&suite.Suite, suite.base.User, &transaction.User)
	}

	isFirstFirst := transactions[0].ID == first.ID
	if isFirstFirst {
		TestTransactionsEqual(&suite.Suite, first, transactions[0])
		TestTransactionsEqual(&suite.Suite, second, transactions[1])
	} else {
		TestTransactionsEqual(&suite.Suite, first, transactions[1])
		TestTransactionsEqual(&suite.Suite, second, transactions[0])
	}
}

func (suite *TransactionRepositorySuit) TestPatchTransaction() {
	transaction, err := repository.CreateTransaction(
		transactionDescription,
		transactionAmount,
		transactionDate,
		suite.base.User.ID,
		suite.base.FirstCategory.ID,
		suite.base.FirstAccount.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	var (
		newDescription = "New transaction description"
		newAmount      = decimal.NewFromInt(600)
		newDate        = time.Now()
	)

	patched, err := repository.PatchTransaction(
		transaction.ID,
		newDescription,
		newAmount,
		newDate,
		suite.base.User.ID,
		suite.base.FirstCategory.ID,
		suite.base.FirstAccount.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	expected := models.Transaction{
		ID:          transaction.ID,
		Description: newDescription,
		Amount:      newAmount,
		Date:        newDate,
		User:        transaction.User,
		Category:    transaction.Category,
		Account:     transaction.Account,
	}

	TestTransactionsEqual(&suite.Suite, &expected, patched)
}

func (suite *TransactionRepositorySuit) TestDeleteTransaction() {
	transaction, err := repository.CreateTransaction(
		transactionDescription,
		transactionAmount,
		transactionDate,
		suite.base.User.ID,
		suite.base.FirstCategory.ID,
		suite.base.FirstAccount.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	result, err := repository.DeleteTransaction(transaction.ID)
	if err != nil {
		suite.Error(err)
	}

	TestTransactionsEqual(&suite.Suite, transaction, result)

	transactions, err := repository.GetUserTransactions(transaction.User.ID)
	if err != nil {
		suite.Error(err)
	}
	suite.Equal(0, len(transactions))
}
