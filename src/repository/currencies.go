package repository

import (
	"main/models"
)

func GetAllCurrencies() ([]*models.Currency, error) {
	return GetAllModels[models.Currency]()
}

func CreateCurrency(name, symbol string) (*models.Currency, error) {
	currency := models.Currency{
		Name:   name,
		Symbol: symbol,
	}
	return CreateModel(&currency)
}
