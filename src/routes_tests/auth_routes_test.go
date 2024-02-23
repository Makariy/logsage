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

type UserRoutesSuit struct {
	suite.Suite
	router *gin.Engine
}

func (suite *UserRoutesSuit) SetupTest() {
	test_utils.CreateTestDB()
	models.MigrateModels(db_connector.GetConnection())

	suite.router = gin.Default()
	routes.AddAuthRoutes(suite.router)
}

func (suite *UserRoutesSuit) TearDownTest() {
	test_utils.DropTestDB()
}

func getLoginForm() []byte {
	form := &forms.UserForm{
		Email:    userEmail,
		Password: userPassword,
	}
	stringForm, _ := json.Marshal(&form)
	return stringForm
}

func (suite *UserRoutesSuit) TestHandleSignup() {
	resp := PerformTestRequest(suite.router, "POST", "/auth/signup/", getLoginForm(), nil)
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
	suite.Equal(userEmail, user.Email, "Signed up with different email")
}

func (suite *UserRoutesSuit) TestHandleSignupEmailExists() {
	resp := PerformTestRequest(suite.router, "POST", "/auth/signup/", getLoginForm(), nil)
	AssertResponseSuccess(201, resp, &suite.Suite)

	resp = PerformTestRequest(suite.router, "POST", "/auth/signup/", getLoginForm(), nil)
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
	_, err := repository.CreateUser(userEmail, userPassword)
	suite.Equal(err, nil, "Got an unexpected error creating user")

	resp := PerformTestRequest(suite.router, "POST", "/auth/login/", getLoginForm(), nil)
	AssertResponseSuccess(200, resp, &suite.Suite)
}

func (suite *UserRoutesSuit) TestHandleLogout() {
	_, err := repository.CreateUser(userEmail, userPassword)
	suite.Equal(err, nil, "Got an unexpected error creating user")

	resp := PerformTestRequest(suite.router, "POST", "/auth/login/", getLoginForm(), nil)
	AssertResponseSuccess(200, resp, &suite.Suite)

	headers := map[string]string{
		"Authorization": resp.Header().Get("Authorization"),
	}
	resp = PerformTestRequest(suite.router, "POST", "/auth/logout/", nil, &headers)
	AssertResponseSuccess(200, resp, &suite.Suite)
}
