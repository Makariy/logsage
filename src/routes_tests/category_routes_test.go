package routes_tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"main/forms"
	"main/models"
	"main/routes"
	"main/test_utils"
	data "main/test_utils/test_data"
	"main/utils"
)

type CategoryRoutesSuit struct {
	suite.Suite
	router test_utils.RoutesDefaultSuite
}

func (suite *CategoryRoutesSuit) SetupTest() {
	suite.router.SetupBase()
	suite.router.Data.CreateDefaultUser()
	suite.router.SetupAuth()
	routes.AddCategoryRoutes(suite.router.Router)
}

func (suite *CategoryRoutesSuit) TearDownTest() {
	suite.router.TearDownTest()
}

func getCategoryForm() []byte {
	form := &forms.CategoryForm{
		Name: data.FirstCategoryName,
		Type: data.FirstCategoryType,
	}
	stringForm, _ := json.Marshal(&form)
	return stringForm
}

func (suite *CategoryRoutesSuit) TestHandleCreateCategory() {
	resp := PerformTestRequest(
		suite.router.Router,
		"POST",
		"/category/create/",
		getCategoryForm(),
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(201, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.CategoryResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	TestCategoriesEqual(&suite.Suite, &categoryResponse, response)
}

func (suite *CategoryRoutesSuit) TestHandleGetCategory() {
	suite.router.Data.CreateDefaultCategories()
	category := suite.router.Data.FirstCategory

	resp := PerformTestRequest(
		suite.router.Router,
		"GET",
		fmt.Sprintf("/category/get/%d/", category.ID),
		nil,
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.CategoryResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	expected, err := utils.MarshalModelToForm[models.Category, forms.CategoryResponse](category)
	if err != nil {
		suite.Error(err)
	}

	TestCategoriesEqual(&suite.Suite, expected, response)
}

func (suite *CategoryRoutesSuit) TestHandleGetAllCategories() {
	suite.router.Data.CreateDefaultCategories()
	firstCategory := suite.router.Data.FirstCategory
	secondCategory := suite.router.Data.SecondCategory

	resp := PerformTestRequest(
		suite.router.Router,
		"GET",
		"/category/all/",
		getCategoryForm(),
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.CategoriesResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	suite.Equal(2, len(response.Categories))

	firstForm, err := utils.MarshalModelToForm[models.Category, forms.CategoryResponse](
		firstCategory,
	)
	if err != nil {
		suite.Error(err)
	}
	secondForm, err := utils.MarshalModelToForm[models.Category, forms.CategoryResponse](
		secondCategory,
	)
	if err != nil {
		suite.Error(err)
	}

	TestCategoriesEqual(&suite.Suite, firstForm, response.Categories[0])
	TestCategoriesEqual(&suite.Suite, secondForm, response.Categories[1])
}

func (suite *CategoryRoutesSuit) TestHandlePatchCategory() {
	var (
		newName = "New category name"
		newType = models.SPENDING
	)
	suite.router.Data.CreateDefaultCategories()
	category := suite.router.Data.FirstCategory

	patchedCategory := forms.CategoryForm{
		Name: newName,
		Type: newType,
	}
	stringPatch, _ := json.Marshal(&patchedCategory)
	resp := PerformTestRequest(
		suite.router.Router,
		"PATCH",
		fmt.Sprintf("/category/patch/%d/", category.ID),
		stringPatch,
		&suite.router.AuthHeaders,
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
		User: *suite.router.Data.User,
	}
	expectedForm, err := utils.MarshalModelToForm[models.Category, forms.CategoryResponse](&expectedCategory)
	if err != nil {
		suite.Error(err)
	}
	TestCategoriesEqual(&suite.Suite, expectedForm, response)
}

func (suite *CategoryRoutesSuit) TestHandleDeleteCategory() {
	suite.router.Data.CreateDefaultCategories()
	category := suite.router.Data.FirstCategory

	resp := PerformTestRequest(
		suite.router.Router,
		"DELETE",
		fmt.Sprintf("/category/delete/%d/", category.ID),
		getCategoryForm(),
		&suite.router.AuthHeaders,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	response, err := UnmarshalResponse[forms.CategoryResponse](resp)
	if err != nil {
		suite.Error(err)
	}

	expectedForm, err := utils.MarshalModelToForm[models.Category, forms.CategoryResponse](category)
	if err != nil {
		suite.Error(err)
	}
	TestCategoriesEqual(&suite.Suite, expectedForm, response)
}
