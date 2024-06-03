package forms

import (
	"github.com/shopspring/decimal"
	"main/models"
)

type AccountForm struct {
	Name       string          `json:"name" form:"name"`
	CurrencyID models.ModelID  `json:"currencyId" form:"currencyId"`
	Balance    decimal.Decimal `json:"balance" form:"balance"`
}

type AccountResponse struct {
	*SuccessResponse
	ID       models.ModelID    `json:"id"`
	Name     string            `json:"name"`
	Currency *CurrencyResponse `json:"currency"`
	Balance  decimal.Decimal   `json:"balance"`
}

type AccountsResponse struct {
	*SuccessResponse
	Accounts []*AccountResponse `json:"accounts"`
}

func (AccountsResponse) ListField() string {
	return "Accounts"
}
