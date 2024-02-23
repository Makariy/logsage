package repository

import (
	"main/models"
)

func GetAllCurrencies() ([]*models.Currency, error) {
	return GetAllModels[models.Currency]()
}

func CreateCurrency(name string) (*models.Currency, error) {
	currency := models.Currency{
		Name: name,
	}
	return CreateModel(&currency)
}
