package repository_tests

import (
	"github.com/stretchr/testify/suite"
	"main/models"
	"main/repository"
)

type CategoryRepositorySuit struct {
	suite.Suite
	base DefaultRepositorySuite
}

func (suite *CategoryRepositorySuit) SetupTest() {
	suite.base.setupDB()
	suite.base.createDefaultUser()
}

func (suite *CategoryRepositorySuit) TearDownTest() {
	suite.base.TearDownTest()
}

func (suite *CategoryRepositorySuit) TestCreateCategory() {
	category, err := repository.CreateCategory(suite.base.user.ID, categoryName, categoryType)
	if err != nil {
		suite.Error(err)
	}

	expected := models.Category{
		ID:   1,
		Name: categoryName,
		Type: categoryType,
		User: *suite.base.user,
	}
	TestCategoriesEqual(&expected, category, &suite.Suite)
}

func (suite *CategoryRepositorySuit) TestGetCategoryByID() {
	category, err := repository.CreateCategory(suite.base.user.ID, categoryName, categoryType)
	if err != nil {
		suite.Error(err)
	}

	result, err := repository.GetCategoryByID(category.ID)
	if err != nil {
		suite.Error(err)
	}

	TestCategoriesEqual(category, result, &suite.Suite)
}

func (suite *CategoryRepositorySuit) TestGetAllCategories() {
	first, err := repository.CreateCategory(suite.base.user.ID, categoryName, categoryType)
	if err != nil {
		suite.Error(err)
	}
	second, err := repository.CreateCategory(suite.base.user.ID, "Other category", models.EARNING)
	if err != nil {
		suite.Error(err)
	}

	categories, err := repository.GetUserCategories(suite.base.user.ID)
	if err != nil {
		suite.Error(err)
	}
	suite.Equal(2, len(categories))

	for _, category := range categories {
		TestUsersEqual(suite.base.user, &category.User, &suite.Suite)
	}

	isFirstFirst := categories[0].ID == first.ID
	if isFirstFirst {
		TestCategoriesEqual(first, categories[0], &suite.Suite)
		TestCategoriesEqual(second, categories[1], &suite.Suite)
	} else {
		TestCategoriesEqual(first, categories[1], &suite.Suite)
		TestCategoriesEqual(second, categories[0], &suite.Suite)
	}
}

func (suite *CategoryRepositorySuit) TestPatchCategory() {
	category, err := repository.CreateCategory(suite.base.user.ID, categoryName, categoryType)
	if err != nil {
		suite.Error(err)
	}

	var (
		newName = "New category"
		newType = "New type"
	)

	patched, err := repository.PatchCategory(category.ID, newName, newType, suite.base.user.ID)
	if err != nil {
		suite.Error(err)
	}

	expected := models.Category{
		ID:   category.ID,
		Name: newName,
		Type: newType,
		User: category.User,
	}

	TestCategoriesEqual(&expected, patched, &suite.Suite)
}

func (suite *CategoryRepositorySuit) TestDeleteCategory() {
	category, err := repository.CreateCategory(suite.base.user.ID, categoryName, categoryType)
	if err != nil {
		suite.Error(err)
	}

	result, err := repository.DeleteCategory(category.ID)
	if err != nil {
		suite.Error(err)
	}
	TestCategoriesEqual(category, result, &suite.Suite)

	categories, err := repository.GetUserCategories(category.User.ID)
	if err != nil {
		suite.Error(err)
	}
	suite.Equal(0, len(categories))
}
