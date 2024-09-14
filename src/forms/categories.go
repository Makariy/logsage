package forms

import "main/models"

type CategoryForm struct {
	Name string              `json:"name" form:"name"`
	Type models.CategoryType `json:"type" form:"type"`
}

type CategoryResponse struct {
	*SuccessResponse
	ID   models.ModelID      `json:"id"`
	Name string              `json:"name"`
	Type models.CategoryType `json:"type"`
}

type CategoriesResponse struct {
	*SuccessResponse
	Categories []*CategoryResponse `json:"categories"`
}

func (CategoriesResponse) ListField() string {
	return "Categories"
}
