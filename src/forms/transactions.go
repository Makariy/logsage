package forms

import (
	"github.com/shopspring/decimal"
	"time"
)

type TransactionForm struct {
	Description string          `json:"name" form:"name"`
	Amount      decimal.Decimal `json:"amount" form:"amount"`
	Date        time.Time       `json:"date" form:"date"`
	CategoryID  uint            `json:"categoryId" form:"categoryId"`
	AccountID   uint            `json:"accountID" form:"accountID"`
}

type TransactionResponse struct {
	*SuccessResponse
	ID          uint             `json:"id"`
	Description string           `json:"name"`
	Amount      decimal.Decimal  `json:"amount"`
	Date        time.Time        `json:"date"`
	Category    CategoryResponse `json:"category"`
	Account     AccountResponse  `json:"account"`
}

type TransactionsResponse struct {
	*SuccessResponse
	Transactions []*TransactionResponse `json:"transactions"`
}

func (TransactionsResponse) ListField() string {
	return "Transactions"
}
