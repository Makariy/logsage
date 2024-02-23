package forms

import (
	"github.com/shopspring/decimal"
)

type AccountForm struct {
	Name       string          `json:"name" form:"name"`
	CurrencyID uint            `json:"currencyId" form:"currencyId"`
	Balance    decimal.Decimal `json:"balance" form:"balance"`
}

type AccountResponse struct {
	*SuccessResponse
	ID       uint              `json:"id"`
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
