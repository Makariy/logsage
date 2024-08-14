package routes_tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"main/db_connector"
	"main/forms"
	"main/models"
	"main/repository"
	"main/routes"
	"main/test_utils"
	data "main/test_utils/test_data"
)

type UserRoutesSuit struct {
	suite.Suite
	router test_utils.RoutesDefaultSuite
}

func (suite *UserRoutesSuit) SetupTest() {
	suite.router.SetupBase()
	routes.AddAuthRoutes(suite.router.Router)
}

func (suite *UserRoutesSuit) TearDownTest() {
	suite.router.TearDownTest()
}

func getLoginForm() []byte {
	form := &forms.UserForm{
		Email:    data.UserEmail,
		Password: data.UserPassword,
	}
	stringForm, _ := json.Marshal(&form)
	return stringForm
}

func (suite *UserRoutesSuit) TestHandleSignup() {
	resp := PerformTestRequest(
		suite.router.Router,
		"POST",
		"/auth/signup/",
		getLoginForm(),
		nil,
	)
	AssertResponseSuccess(201, resp, &suite.Suite)
	suite.NotNil(resp.Header().Get("Authorization"), "No authorization token provided in response")

	db := db_connector.GetConnection()

	var users []models.User
	tx := db.Find(&users)
	if tx.Error != nil {
		panic(fmt.Sprintf("got an unexpected error reading users: %v", tx.Error))
	}
	suite.Equal(1, len(users), "Created an unexpected amount of users")

	user := users[0]
	suite.Equal(
		data.UserEmail,
		user.Email,
		"Signed up with different email",
	)
}

func (suite *UserRoutesSuit) TestHandleSignupEmailExists() {
	resp := PerformTestRequest(
		suite.router.Router,
		"POST",
		"/auth/signup/",
		getLoginForm(),
		nil,
	)
	AssertResponseSuccess(201, resp, &suite.Suite)

	resp = PerformTestRequest(
		suite.router.Router,
		"POST",
		"/auth/signup/",
		getLoginForm(),
		nil,
	)
	suite.Equal(400, resp.Code)
	suite.Equal("", resp.Header().Get("Authorization"))

	db := db_connector.GetConnection()
	var users []models.User
	tx := db.Find(&users)
	if tx.Error != nil {
		suite.Error(tx.Error)
	}
	suite.Equal(1, len(users))
}

func (suite *UserRoutesSuit) TestHandleLogin() {
	_, err := repository.CreateUser(
		data.UserEmail,
		data.UserPassword,
	)
	suite.Equal(err, nil, "Got an unexpected error creating user")

	resp := PerformTestRequest(
		suite.router.Router,
		"POST",
		"/auth/login/",
		getLoginForm(),
		nil,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)
}

func (suite *UserRoutesSuit) TestHandleLogout() {
	_, err := repository.CreateUser(
		data.UserEmail,
		data.UserPassword,
	)
	suite.Equal(err, nil, "Got an unexpected error creating user")

	resp := PerformTestRequest(
		suite.router.Router,
		"POST",
		"/auth/login/",
		getLoginForm(),
		nil,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)

	headers := map[string]string{
		"Authorization": resp.Header().Get("Authorization"),
	}
	resp = PerformTestRequest(
		suite.router.Router,
		"POST",
		"/auth/logout/",
		nil,
		&headers,
	)
	AssertResponseSuccess(200, resp, &suite.Suite)
}
