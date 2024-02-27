package repository_tests

import (
	"github.com/shopspring/decimal"
	"main/models"
	"time"
)

var (
	userEmail    = "test@test.com"
	userPassword = "testpassword"

	currencyName = "USDT"

	accountName    = "KuCoin"
	accountBalance = decimal.NewFromInt(1000)

	categoryName = "Products"
	categoryType = models.SPENDING

	transactionDescription = "Test transaction description"
	transactionDate        = time.Now()
	transactionAmount      = decimal.NewFromInt(100)
)
