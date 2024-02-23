package repository

import (
	"main/models"
)

func CreateCategory(userID uint, name string, categoryType string) (*models.Category, error) {
	category := models.Category{
		Name:   name,
		Type:   categoryType,
		UserID: userID,
	}
	return CreateModel(&category)
}

func GetCategoryByID(id uint) (*models.Category, error) {
	return GetModelByID[models.Category](id)
}

func GetUserCategories(userID uint) ([]*models.Category, error) {
	return GetUserModels[models.Category](userID)
}

func PatchCategory(categoryID uint, name string, categoryType string, userID uint) (*models.Category, error) {
	category := models.Category{
		ID:     categoryID,
		Name:   name,
		Type:   categoryType,
		UserID: userID,
	}
	return PatchModel(&category)
}

func DeleteCategory(categoryID uint) (*models.Category, error) {
	return DeleteModel[models.Category](categoryID)
}
