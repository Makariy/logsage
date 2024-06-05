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

func CreateTestCurrency(name, symbol string) *models.Currency {
	currency, err := repository.CreateCurrency(name, symbol)
	if err != nil {
		panic("could not create test currency")
	}
	return currency
}

func CreateTestCategory(name string, categoryType string, userID models.ModelID) *models.Category {
	category, err := repository.CreateCategory(userID, name, categoryType)
	if err != nil {
		panic("could not create test category")
	}
	return category
}

func CreateTestAccount(name string, balance decimal.Decimal, userID, currencyID models.ModelID) *models.Account {
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
	userID models.ModelID,
	categoryID models.ModelID,
	accountID models.ModelID,
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
