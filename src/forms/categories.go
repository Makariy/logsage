package forms

type CategoryForm struct {
	Name string `json:"name" form:"name"`
	Type string `json:"type" form:"type"`
}

type CategoryResponse struct {
	*SuccessResponse
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type CategoriesResponse struct {
	*SuccessResponse
	Categories []*CategoryResponse `json:"categories"`
}

func (CategoriesResponse) ListField() string {
	return "Categories"
}
