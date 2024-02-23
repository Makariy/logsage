package repository_tests

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"main/db_connector"
	"main/models"
	"main/repository"
	"main/test_utils"
)

type CategoryRepositorySuit struct {
	suite.Suite
	router *gin.Engine
	user   *models.User
}

func (suite *CategoryRepositorySuit) SetupTest() {
	test_utils.CreateTestDB()
	models.MigrateModels(db_connector.GetConnection())

	suite.user = CreateTestUser(userEmail, userPassword)
}

func (suite *CategoryRepositorySuit) TearDownTest() {
	test_utils.DropTestDB()
}

func (suite *CategoryRepositorySuit) TestCreateCategory() {
	category, err := repository.CreateCategory(suite.user.ID, categoryName, categoryType)
	if err != nil {
		suite.Error(err)
	}

	expected := models.Category{
		ID:   1,
		Name: categoryName,
		Type: categoryType,
		User: *suite.user,
	}
	testCategoriesEqual(&expected, category, &suite.Suite)
}

func (suite *CategoryRepositorySuit) TestGetCategoryByID() {
	category, err := repository.CreateCategory(suite.user.ID, categoryName, categoryType)
	if err != nil {
		suite.Error(err)
	}

	result, err := repository.GetCategoryByID(category.ID)
	if err != nil {
		suite.Error(err)
	}

	testCategoriesEqual(category, result, &suite.Suite)
}

func (suite *CategoryRepositorySuit) TestGetAllCategories() {
	first, err := repository.CreateCategory(suite.user.ID, categoryName, categoryType)
	if err != nil {
		suite.Error(err)
	}
	second, err := repository.CreateCategory(suite.user.ID, "Other category", models.EARNING)
	if err != nil {
		suite.Error(err)
	}

	categories, err := repository.GetUserCategories(suite.user.ID)
	if err != nil {
		suite.Error(err)
	}
	suite.Equal(2, len(categories))

	for _, category := range categories {
		testUsersEqual(suite.user, &category.User, &suite.Suite)
	}

	isFirstFirst := categories[0].ID == first.ID
	if isFirstFirst {
		testCategoriesEqual(first, categories[0], &suite.Suite)
		testCategoriesEqual(second, categories[1], &suite.Suite)
	} else {
		testCategoriesEqual(first, categories[1], &suite.Suite)
		testCategoriesEqual(second, categories[0], &suite.Suite)
	}
}

func (suite *CategoryRepositorySuit) TestPatchCategory() {
	category, err := repository.CreateCategory(suite.user.ID, categoryName, categoryType)
	if err != nil {
		suite.Error(err)
	}

	var (
		newName = "New category"
		newType = "New type"
	)

	patched, err := repository.PatchCategory(category.ID, newName, newType, suite.user.ID)
	if err != nil {
		suite.Error(err)
	}

	expected := models.Category{
		ID:   category.ID,
		Name: newName,
		Type: newType,
		User: category.User,
	}

	testCategoriesEqual(&expected, patched, &suite.Suite)
}

func (suite *CategoryRepositorySuit) TestDeleteCategory() {
	category, err := repository.CreateCategory(suite.user.ID, categoryName, categoryType)
	if err != nil {
		suite.Error(err)
	}

	result, err := repository.DeleteCategory(category.ID)
	if err != nil {
		suite.Error(err)
	}
	testCategoriesEqual(category, result, &suite.Suite)

	categories, err := repository.GetUserCategories(category.User.ID)
	if err != nil {
		suite.Error(err)
	}
	suite.Equal(0, len(categories))
}
