package forms

import (
	"github.com/shopspring/decimal"
	"main/models"
)

type AccountStats struct {
	Account models.Account `json:"account"`

	TotalEarnedAmount decimal.Decimal `json:"totalEarnedAmount"`
	TotalSpentAmount  decimal.Decimal `json:"totalSpentAmount"`

	Transactions []*models.Transaction `json:"transactions"`
}

type CategoryStats struct {
	Category models.Category `json:"category"`

	TotalAmount decimal.Decimal `json:"totalAmount"`

	Transactions []*models.Transaction `json:"transactions"`
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

type CategoryStatsResponse struct {
	*SuccessResponse
	*DateRange
	Category CategoryResponse `json:"category"`
	Stats    CategoryStats    `json:"stats"`
}

type AccountStatsResponse struct {
	*SuccessResponse
	*DateRange
	Account AccountResponse `json:"account"`
	Stats   AccountStats    `json:"stats"`
}

type TotalCategoriesStatsResponse struct {
	*SuccessResponse
	*DateRange
	Stats TotalCategoriesStats `json:"stats"`
}

type TotalAccountsStatsResponse struct {
	*SuccessResponse
	*DateRange
	Stats TotalAccountsStats `json:"stats"`
}
