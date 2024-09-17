package test_utils

import (
	"github.com/shopspring/decimal"
	"main/auth"
	"main/models"
	"main/repository"
)

func CreateTestUser(userEmail, userPassword string) *models.User {
	user, err := repository.CreateUser(userEmail, userPassword)
	if err != nil {
		panic("could not create test user")
	}
	return user
}

func CreateTestCurrency(currencyName string, value decimal.Decimal) *models.Currency {
	currency, err := repository.CreateCurrency(currencyName, currencyName, value)
	if err != nil {
		panic("could not create test currency")
	}
	return currency
}

func CreateTestCategory(
	userID models.ModelID,
	categoryName string,
	categoryType models.CategoryType,
	categoryImageID models.ModelID,
) *models.Category {
	category, err := repository.CreateCategory(userID, categoryName, categoryType, categoryImageID)
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
