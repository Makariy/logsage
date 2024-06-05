package repository_tests

import (
	"github.com/shopspring/decimal"
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

func CreateTestCurrency(currencyName string) *models.Currency {
	currency, err := repository.CreateCurrency(currencyName, currencyName)
	if err != nil {
		panic("could not create test currency")
	}
	return currency
}

func CreateTestCategory(userID models.ModelID, categoryName string, categoryType string) *models.Category {
	category, err := repository.CreateCategory(userID, categoryName, categoryType)
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
