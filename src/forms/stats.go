package forms

import "main/models"

type CategoryStatsResponse struct {
	*SuccessResponse
	Category CategoryResponse     `json:"category"`
	Stats    models.CategoryStats `json:"stats"`
}

type AccountStatsResponse struct {
	*SuccessResponse
	Account AccountResponse     `json:"account"`
	Stats   models.AccountStats `json:"stats"`
}

type TotalStatsResponse struct {
	*SuccessResponse
	Stats models.TotalStats `json:"stats"`
}
