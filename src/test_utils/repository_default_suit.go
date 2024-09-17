package test_utils

import (
	"github.com/stretchr/testify/suite"
	"main/db_connector"
	"main/models"
	"main/repository"
	data "main/test_utils/test_data"
)

type RepositoryDefaultSuite struct {
	suite.Suite

	User *models.User

	CategoryImage *models.CategoryImage

	FirstCurrency  *models.Currency
	SecondCurrency *models.Currency

	FirstCategory  *models.Category
	SecondCategory *models.Category

	FirstAccount  *models.Account
	SecondAccount *models.Account

	Transaction1 *models.Transaction
	Transaction2 *models.Transaction
	Transaction3 *models.Transaction
	Transaction4 *models.Transaction
}

func (suite *RepositoryDefaultSuite) SetupDB() {
	CreateTestDB()
	models.MigrateModels(db_connector.GetConnection())
	suite.CategoryImage, _ = repository.CreateModel[models.CategoryImage](&models.CategoryImage{FileName: "test.svg"})
}

func (suite *RepositoryDefaultSuite) CreateDefaultUser() {
	suite.User = CreateTestUser(data.UserEmail, data.UserPassword)
}

func (suite *RepositoryDefaultSuite) CreateDefaultCurrencies() {
	suite.FirstCurrency = CreateTestCurrency(data.FirstCurrencySymbol, data.FirstCurrencyValue)
	suite.SecondCurrency = CreateTestCurrency(data.SecondCurrencySymbol, data.SecondCurrencyValue)
}

func (suite *RepositoryDefaultSuite) CreateDefaultCategories() {
	suite.FirstCategory = CreateTestCategory(
		suite.User.ID,
		data.FirstCategoryName,
		data.FirstCategoryType,
		suite.CategoryImage.ID,
	)
	suite.SecondCategory = CreateTestCategory(
		suite.User.ID,
		data.SecondCategoryName,
		data.SecondCategoryType,
		suite.CategoryImage.ID,
	)
}

func (suite *RepositoryDefaultSuite) CreateDefaultAccounts() {
	suite.FirstAccount = CreateTestAccount(
		data.FirstAccountName,
		data.FirstAccountBalance,
		suite.User.ID,
		suite.FirstCurrency.ID,
	)

	suite.SecondAccount = CreateTestAccount(
		data.SecondAccountName,
		data.SecondAccountBalance,
		suite.User.ID,
		suite.SecondCurrency.ID,
	)
}

func (suite *RepositoryDefaultSuite) CreateDefaultTransactions() {
	suite.Transaction1, _ = repository.CreateTransaction(
		"First transaction",
		data.Transaction1Amount,
		data.Transaction1Date,
		suite.User.ID,
		suite.FirstCategory.ID,
		suite.FirstAccount.ID,
	)
	suite.Transaction2, _ = repository.CreateTransaction(
		"Second transaction",
		data.Transaction2Amount,
		data.Transaction2Date,
		suite.User.ID,
		suite.FirstCategory.ID,
		suite.FirstAccount.ID,
	)
	suite.Transaction3, _ = repository.CreateTransaction(
		"Third transaction",
		data.Transaction3Amount,
		data.Transaction3Date,
		suite.User.ID,
		suite.SecondCategory.ID,
		suite.SecondAccount.ID,
	)
	suite.Transaction4, _ = repository.CreateTransaction(
		"Forth transaction",
		data.Transaction4Amount,
		data.Transaction4Date,
		suite.User.ID,
		suite.SecondCategory.ID,
		suite.SecondAccount.ID,
	)
}

func (suite *RepositoryDefaultSuite) SetupAllTestData() {
	suite.CreateDefaultCurrencies()
	suite.CreateDefaultUser()
	suite.CreateDefaultCategories()
	suite.CreateDefaultAccounts()
	suite.CreateDefaultTransactions()
}

func (suite *RepositoryDefaultSuite) TearDownTest() {
	DropTestDB()
}
