package repository_tests

import (
	"github.com/shopspring/decimal"
	"main/models"
	"time"
)

var (
	userEmail    = "test@test.com"
	userPassword = "testpassword"

	firstCurrencySymbol = "USDT"
	firstCurrencyValue  = decimal.NewFromInt(1)

	secondCurrencySymbol = "RUB"
	secondCurrencyValue  = decimal.NewFromFloat(0.01)

	accountName    = "KuCoin"
	accountBalance = decimal.NewFromInt(1000)

	categoryName = "Products"
	categoryType = models.SPENDING

	transactionDescription = "Test transaction description"
	transactionDate        = time.Now()
	transactionAmount      = decimal.NewFromInt(100)

	transaction1Amount = decimal.NewFromInt(100)
	transaction2Amount = decimal.NewFromInt(200)
	transaction3Amount = decimal.NewFromInt(300)
	transaction4Amount = decimal.NewFromInt(400)

	transactionDateYear  = 2024
	transactionDateMonth = time.January
	transaction1Date     = time.Date(transactionDateYear, transactionDateMonth, 1, 0, 0, 0, 0, time.UTC)
	transaction2Date     = transaction1Date.AddDate(0, 0, 1)
	transaction3Date     = transaction2Date.AddDate(0, 0, 1)
	transaction4Date     = transaction3Date.AddDate(0, 0, 1)
)
