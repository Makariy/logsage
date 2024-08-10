package repository

import (
	"github.com/shopspring/decimal"
	"main/db_connector"
	"main/models"
)

func GetAllCurrencies() ([]*models.Currency, error) {
	return GetAllModels[models.Currency]()
}

func GetCurrencyBySymbol(symbol string) (*models.Currency, error) {
	db := db_connector.GetConnection()

	var currency models.Currency
	tx := db.Model(models.Currency{}).
		Where("symbol = ?", symbol).
		Scan(&currency)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &currency, nil
}

func CreateCurrency(
	name,
	symbol string,
	value decimal.Decimal,
) (*models.Currency, error) {
	currency := models.Currency{
		Name:   name,
		Symbol: symbol,
		Value:  value,
	}
	return CreateModel(&currency)
}
