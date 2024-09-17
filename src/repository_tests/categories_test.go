package repository_tests

import (
	"github.com/stretchr/testify/suite"
	"main/models"
	"main/repository"
	"main/test_utils"
)

type CategoryRepositorySuit struct {
	suite.Suite
	base test_utils.RepositoryDefaultSuite
}

func (suite *CategoryRepositorySuit) SetupTest() {
	suite.base.SetupDB()
	suite.base.CreateDefaultUser()
}

func (suite *CategoryRepositorySuit) TearDownTest() {
	suite.base.TearDownTest()
}

func (suite *CategoryRepositorySuit) TestCreateCategory() {
	category, err := repository.CreateCategory(
		suite.base.User.ID,
		categoryName,
		categoryType,
		suite.base.CategoryImage.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	expected := models.Category{
		ID:   1,
		Name: categoryName,
		Type: categoryType,
		User: *suite.base.User,
	}
	TestCategoriesEqual(&suite.Suite, &expected, category)
}

func (suite *CategoryRepositorySuit) TestGetCategoryByID() {
	category, err := repository.CreateCategory(
		suite.base.User.ID,
		categoryName,
		categoryType,
		suite.base.CategoryImage.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	result, err := repository.GetCategoryByID(category.ID)
	if err != nil {
		suite.Error(err)
	}

	TestCategoriesEqual(&suite.Suite, category, result)
}

func (suite *CategoryRepositorySuit) TestGetAllCategories() {
	first, err := repository.CreateCategory(
		suite.base.User.ID,
		categoryName,
		categoryType,
		suite.base.CategoryImage.ID,
	)
	if err != nil {
		suite.Error(err)
	}
	second, err := repository.CreateCategory(
		suite.base.User.ID,
		"Other category",
		models.EARNING,
		suite.base.CategoryImage.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	categories, err := repository.GetUserCategories(suite.base.User.ID)
	if err != nil {
		suite.Error(err)
	}
	suite.Equal(2, len(categories))

	for _, category := range categories {
		TestUsersEqual(&suite.Suite, suite.base.User, &category.User)
	}

	isFirstFirst := categories[0].ID == first.ID
	if isFirstFirst {
		TestCategoriesEqual(&suite.Suite, first, categories[0])
		TestCategoriesEqual(&suite.Suite, second, categories[1])
	} else {
		TestCategoriesEqual(&suite.Suite, first, categories[1])
		TestCategoriesEqual(&suite.Suite, second, categories[0])
	}
}

func (suite *CategoryRepositorySuit) TestPatchCategory() {
	category, err := repository.CreateCategory(
		suite.base.User.ID,
		categoryName,
		categoryType,
		suite.base.CategoryImage.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	var (
		newName                     = "New category"
		newType models.CategoryType = "New type"
	)

	patched, err := repository.PatchCategory(
		category.ID,
		newName,
		newType,
		suite.base.User.ID,
		suite.base.CategoryImage.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	expected := models.Category{
		ID:   category.ID,
		Name: newName,
		Type: newType,
		User: category.User,
	}

	TestCategoriesEqual(&suite.Suite, &expected, patched)
}

func (suite *CategoryRepositorySuit) TestDeleteCategory() {
	category, err := repository.CreateCategory(
		suite.base.User.ID,
		categoryName,
		categoryType,
		suite.base.CategoryImage.ID,
	)
	if err != nil {
		suite.Error(err)
	}

	result, err := repository.DeleteCategory(category.ID)
	if err != nil {
		suite.Error(err)
	}
	TestCategoriesEqual(&suite.Suite, category, result)

	categories, err := repository.GetUserCategories(category.User.ID)
	if err != nil {
		suite.Error(err)
	}
	suite.Equal(0, len(categories))
}
