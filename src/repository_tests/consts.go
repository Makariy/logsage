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
	accountBalance = decimal.New(1000, 10)

	categoryName = "Products"
	categoryType = models.SPENDING

	transactionDescription = "Test transaction description"
	transactionDate        = time.Now()
	transactionAmount      = decimal.New(100, 10)
)
