package routes_tests

import (
	"github.com/shopspring/decimal"
	"main/auth"
	"main/models"
	"main/repository"
	"time"
)

func CreateTestUser(userEmail, userPassword string) *models.User {
	user, err := repository.CreateUser(userEmail, userPassword)
	if err != nil {
		panic("could not create test user")
	}
	return user
}

func CreateTestCurrency(currencyName string) *models.Currency {
	currency, err := repository.CreateCurrency(currencyName)
	if err != nil {
		panic("could not create test currency")
	}
	return currency
}

func CreateTestCategory(name string, categoryType string, userID uint) *models.Category {
	category, err := repository.CreateCategory(userID, name, categoryType)
	if err != nil {
		panic("could not create test category")
	}
	return category
}

func CreateTestAccount(name string, balance decimal.Decimal, userID, currencyID uint) *models.Account {
	account, err := repository.CreateAccount(
		userID,
		name,
		balance,
		currencyID,
	)
	if err != nil {
		panic("could not create account")
	}
	return account
}

func CreateTestTransaction(
	description string,
	amount decimal.Decimal,
	date time.Time,
	userID uint,
	categoryID uint,
	accountID uint,
) *models.Transaction {
	transaction, err := repository.CreateTransaction(description, amount, date, userID, categoryID, accountID)
	if err != nil {
		panic("could not create test transaction")
	}
	return transaction
}

func GetAuthHeaders(user *models.User) map[string]string {
	headers := make(map[string]string)

	token := auth.CreateAuthToken()
	err := auth.SetUserByToken(user, token)
	if err != nil {
		panic("Could not create auth token")
	}

	headers["Authorization"] = auth.RenderAuthorizationHeader(token)
	return headers
}
