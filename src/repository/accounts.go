package repository

import (
	"github.com/shopspring/decimal"
	"main/models"
)

func CreateAccount(userId models.ModelID, name string, balance decimal.Decimal, currencyId models.ModelID) (*models.Account, error) {
	account := models.Account{
		UserID:     userId,
		Name:       name,
		CurrencyID: currencyId,
		Balance:    balance,
	}

	return CreateModel(&account)
}

func PatchAccount(id models.ModelID, name string, balance decimal.Decimal, currencyId models.ModelID, userID models.ModelID) (*models.Account, error) {
	account := models.Account{
		ID:         id,
		Name:       name,
		Balance:    balance,
		CurrencyID: currencyId,
		UserID:     userID,
	}
	return PatchModel(&account)
}

func GetAccountByID(id models.ModelID) (*models.Account, error) {
	return GetModelByID[models.Account](id)
}

func GetUserAccounts(userId models.ModelID) ([]*models.Account, error) {
	return GetUserModels[models.Account](userId)
}

func DeleteAccount(id models.ModelID) (*models.Account, error) {
	return DeleteModel[models.Account](id)
}
