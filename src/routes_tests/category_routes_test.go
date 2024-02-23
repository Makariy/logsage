package routes_tests

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"main/db_connector"
	"main/forms"
	"main/models"
	"main/repository"
	"main/routes"
	"main/test_utils"
)

type CategoryRoutesSuit struct {
	suite.Suite
	router *gin.Engine

	user        *models.User
	authHeaders map[string]string
}

func (suite *CategoryRoutesSuit) SetupTest() {
	test_utils.CreateTestDB()
	models.MigrateModels(db_connector.GetConnection())

	suite.user = CreateTestUser(userEmail, userPassword)
	suite.authHeaders = GetAuthHeaders(suite.user)

	suite.router = gin.Default()
	routes.AddCategoryRoutes(suite.router)
}

func (suite *CategoryRoutesSuit) TearDownTest() {
	test_utils.DropTestDB()
}

func getCategoryForm() []byte {
	form := &forms.CategoryForm{
		Name: categoryName,
		Type: categoryType,
	}
	stringForm, _ := json.Marshal(&form)
	return stringForm
}

func (suite *CategoryRoutesSuit) TestHandleCreateCategory() {
	resp := PerformTestRequest(
		suite.router,
		"POST",
		"/category/create/",
		getCategoryForm(),
		&suite.authHeaders,
	)
	AssertResponseSuccess(201, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.CategoryResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	testCategoriesEqual(&categoryResponse, response, &suite.Suite)
}

func (suite *CategoryRoutesSuit) TestHandleGetCategory() {
	category, err := repository.CreateCategory(suite.user.ID, categoryName, categoryType)
	if err != nil {
		suite.Error(err)
	}

	resp := PerformTestRequest(
		suite.router,
		"GET",
		fmt.Sprintf("/category/get/%d/", category.ID),
		nil,
		&suite.authHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.CategoryResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	expected, err := MarshalModelToForm[models.Category, forms.CategoryResponse](category)
	if err != nil {
		suite.Error(err)
	}

	testCategoriesEqual(expected, response, &suite.Suite)
}

func (suite *CategoryRoutesSuit) TestHandleGetAllCategories() {
	first, err := repository.CreateCategory(suite.user.ID, categoryName, categoryType)
	if err != nil {
		suite.Error(err)
	}

	second, err := repository.CreateCategory(suite.user.ID, "Another Category", models.EARNING)
	if err != nil {
		suite.Error(err)
	}

	resp := PerformTestRequest(
		suite.router,
		"GET",
		"/category/all/",
		getCategoryForm(),
		&suite.authHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.CategoriesResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	suite.Equal(2, len(response.Categories))

	firstForm, err := MarshalModelToForm[models.Category, forms.CategoryResponse](first)
	if err != nil {
		suite.Error(err)
	}
	secondForm, err := MarshalModelToForm[models.Category, forms.CategoryResponse](second)
	if err != nil {
		suite.Error(err)
	}

	testCategoriesEqual(firstForm, response.Categories[0], &suite.Suite)
	testCategoriesEqual(secondForm, response.Categories[1], &suite.Suite)
}

func (suite *CategoryRoutesSuit) TestHandlePatchCategory() {
	var (
		newName = "New category name"
		newType = models.SPENDING
	)

	category, err := repository.CreateCategory(suite.user.ID, categoryName, categoryType)
	if err != nil {
		suite.Error(err)
	}

	patchedCategory := forms.CategoryForm{
		Name: newName,
		Type: newType,
	}
	stringPatch, _ := json.Marshal(&patchedCategory)
	resp := PerformTestRequest(
		suite.router,
		"PATCH",
		fmt.Sprintf("/category/patch/%d/", category.ID),
		stringPatch,
		&suite.authHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.CategoryResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	expectedCategory := models.Category{
		ID:   1,
		Name: newName,
		Type: newType,
		User: *suite.user,
	}
	expectedForm, err := MarshalModelToForm[models.Category, forms.CategoryResponse](&expectedCategory)
	if err != nil {
		suite.Error(err)
	}
	testCategoriesEqual(expectedForm, response, &suite.Suite)
}

func (suite *CategoryRoutesSuit) TestHandleDeleteCategory() {
	category, err := repository.CreateCategory(suite.user.ID, categoryName, categoryType)
	if err != nil {
		suite.Error(err)
	}

	resp := PerformTestRequest(
		suite.router,
		"DELETE",
		fmt.Sprintf("/category/delete/%d/", category.ID),
		getCategoryForm(),
		&suite.authHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.CategoryResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	expectedForm, err := MarshalModelToForm[models.Category, forms.CategoryResponse](category)
	if err != nil {
		suite.Error(err)
	}
	testCategoriesEqual(expectedForm, response, &suite.Suite)
}
