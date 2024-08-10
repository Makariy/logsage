package models

import "github.com/shopspring/decimal"

type CategoryStats struct {
	Category Category `json:"category"`

	TotalAmount decimal.Decimal `json:"totalAmount"`

	Transactions []*Transaction `json:"transactions"`
}

type AccountStats struct {
	Account Account `json:"account"`

	TotalEarnedAmount decimal.Decimal `json:"totalEarnedAmount"`
	TotalSpentAmount  decimal.Decimal `json:"totalSpentAmount"`

	Transactions []*Transaction `json:"transactions"`
}

type TotalAccountsStats struct {
	TotalEarnedAmount decimal.Decimal `json:"totalEarnedAmount"`
	TotalSpentAmount  decimal.Decimal `json:"totalSpentAmount"`

	AccountsStats []*AccountStats `json:"accountsStats"`
}

type TotalCategoriesStats struct {
	TotalEarnedAmount decimal.Decimal `json:"totalEarnedAmount"`
	TotalSpentAmount  decimal.Decimal `json:"totalSpentAmount"`

	CategoriesStats []*CategoryStats `json:"categoriesStats"`
}
