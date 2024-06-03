package forms

import (
	"github.com/shopspring/decimal"
	"main/models"
	"time"
)

type TransactionForm struct {
	Description string          `json:"description" form:"description"`
	Amount      decimal.Decimal `json:"amount" form:"amount"`
	Date        time.Time       `json:"date" form:"date"`
	CategoryID  models.ModelID  `json:"categoryId" form:"categoryId"`
	AccountID   models.ModelID  `json:"accountID" form:"accountID"`
}

type TransactionResponse struct {
	*SuccessResponse
	ID          models.ModelID   `json:"id"`
	Description string           `json:"description"`
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
