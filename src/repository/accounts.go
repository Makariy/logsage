package repository

import (
	"github.com/shopspring/decimal"
	"main/models"
)

func CreateAccount(userId uint, name string, balance decimal.Decimal, currencyId uint) (*models.Account, error) {
	account := models.Account{
		UserID:     userId,
		Name:       name,
		CurrencyID: currencyId,
		Balance:    balance,
	}

	return CreateModel(&account)
}

func PatchAccount(id uint, name string, balance decimal.Decimal, currencyId uint, userID uint) (*models.Account, error) {
	account := models.Account{
		ID:         id,
		Name:       name,
		Balance:    balance,
		CurrencyID: currencyId,
		UserID:     userID,
	}
	return PatchModel(&account)
}

func GetAccountByID(id uint) (*models.Account, error) {
	return GetModelByID[models.Account](id)
}

func GetUserAccounts(userId uint) ([]*models.Account, error) {
	return GetUserModels[models.Account](userId)
}

func DeleteAccount(id uint) (*models.Account, error) {
	return DeleteModel[models.Account](id)
}
