package repository_tests

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestCategoriesSuite(t *testing.T) {
	suite.Run(t, new(CategoryRepositorySuit))
}

func TestAccountSuite(t *testing.T) {
	suite.Run(t, new(AccountRepositorySuit))
}

func TestTransactionSuite(t *testing.T) {
	suite.Run(t, new(TransactionRepositorySuit))
}
