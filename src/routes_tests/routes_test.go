package routes_tests

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestUserRoutesSuite(t *testing.T) {
	suite.Run(t, new(UserRoutesSuit))
}

func TestAccountRoutesSuite(t *testing.T) {
	suite.Run(t, new(AccountRoutesSuit))
}

func TestCategoryRoutesSuite(t *testing.T) {
	suite.Run(t, new(CategoryRoutesSuit))
}
func TestCurrencyRoutesSuite(t *testing.T) {
	suite.Run(t, new(CurrencyRoutesSuit))
}

func TestTransactionRoutesSuite(t *testing.T) {
	suite.Run(t, new(TransactionRoutesSuit))
}

func TestStatsRoutesSuite(t *testing.T) {
	suite.Run(t, new(StatsRoutesSuit))
}
