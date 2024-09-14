package repository

import (
	"github.com/shopspring/decimal"
	"main/models"
)

func ConvertToRelativeCurrency(
	amount decimal.Decimal,
	currency *models.Currency,
) decimal.Decimal {
	return amount.Div(currency.Value)
}
