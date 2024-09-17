package forms

import "main/models"

type CategoryImageResponse struct {
	ID       models.ModelID `json:"id"`
	Filename string         `json:"filename"`
}

type CategoryImagesResponse struct {
	*SuccessResponse
	CategoryImages []*CategoryImageResponse `json:"images"`
}

func (CategoryImagesResponse) ListField() string {
	return "CategoryImages"
}

type CategoryForm struct {
	Name            string              `json:"name" form:"name"`
	Type            models.CategoryType `json:"type" form:"type"`
	CategoryImageID models.ModelID      `json:"category_image_id" form:"category_image_id"`
}

type CategoryResponse struct {
	*SuccessResponse
	ID            models.ModelID        `json:"id"`
	Name          string                `json:"name"`
	Type          models.CategoryType   `json:"type"`
	CategoryImage CategoryImageResponse `json:"image"`
}

type CategoriesResponse struct {
	*SuccessResponse
	Categories []*CategoryResponse `json:"categories"`
}

func (CategoriesResponse) ListField() string {
	return "Categories"
}
