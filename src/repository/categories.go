package repository

import (
	"main/models"
)

func CreateCategory(userID models.ModelID, name string, categoryType string) (*models.Category, error) {
	category := models.Category{
		Name:   name,
		Type:   categoryType,
		UserID: userID,
	}
	return CreateModel(&category)
}

func GetCategoryByID(id models.ModelID) (*models.Category, error) {
	return GetModelByID[models.Category](id)
}

func GetUserCategories(userID models.ModelID) ([]*models.Category, error) {
	return GetUserModels[models.Category](userID)
}

func PatchCategory(categoryID models.ModelID, name string, categoryType string, userID models.ModelID) (*models.Category, error) {
	category := models.Category{
		ID:     categoryID,
		Name:   name,
		Type:   categoryType,
		UserID: userID,
	}
	return PatchModel(&category)
}

func DeleteCategory(categoryID models.ModelID) (*models.Category, error) {
	return DeleteModel[models.Category](categoryID)
}
