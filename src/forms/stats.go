package forms

import "main/models"

type CategoryStatsResponse struct {
	*SuccessResponse
	*DateRange
	Category CategoryResponse     `json:"category"`
	Stats    models.CategoryStats `json:"stats"`
}

type AccountStatsResponse struct {
	*SuccessResponse
	*DateRange
	Account AccountResponse     `json:"account"`
	Stats   models.AccountStats `json:"stats"`
}

type TotalCategoriesStatsResponse struct {
	*SuccessResponse
	*DateRange
	Stats models.TotalCategoriesStats `json:"stats"`
}

type TotalAccountsStatsResponse struct {
	*SuccessResponse
	*DateRange
	Stats models.TotalAccountsStats `json:"stats"`
}
