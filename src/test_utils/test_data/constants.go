package test_data

import (
	"github.com/shopspring/decimal"
	"main/models"
	"time"
)

var (
	UserEmail    = "test@test.com"
	UserPassword = "testpassword"

	FirstCurrencySymbol  = "USDT"
	FirstCurrencyValue   = decimal.NewFromInt(1)
	SecondCurrencySymbol = "RUB"
	SecondCurrencyValue  = decimal.NewFromFloat(0.01)

	FirstAccountName     = "KuCoin"
	FirstAccountBalance  = decimal.NewFromInt(1000)
	SecondAccountName    = "Revolut"
	SecondAccountBalance = decimal.NewFromInt(100000)

	FirstCategoryName  = "Groceries"
	FirstCategoryType  = models.SPENDING
	SecondCategoryName = "Work"
	SecondCategoryType = models.EARNING

	TransactionDescription = "Test transaction description"
	TransactionDate        = time.Now()
	TransactionAmount      = decimal.NewFromInt(100)

	Transaction1Amount = decimal.NewFromInt(100)
	Transaction2Amount = decimal.NewFromInt(200)
	Transaction3Amount = decimal.NewFromInt(300)
	Transaction4Amount = decimal.NewFromInt(400)

	TransactionDateYear  = 2024
	TransactionDateMonth = time.January
	Transaction1Date     = time.Date(TransactionDateYear, TransactionDateMonth, 1, 0, 0, 0, 0, time.UTC)
	Transaction2Date     = Transaction1Date.AddDate(0, 0, 1)
	Transaction3Date     = Transaction2Date.AddDate(0, 0, 1)
	Transaction4Date     = Transaction3Date.AddDate(0, 0, 1)
)
