package routes_tests

import (
	"github.com/shopspring/decimal"
	"main/forms"
	"main/models"
	"time"
)

var (
	userEmail    = "test@test.com"
	userPassword = "testpassword"

	accountName    = "KuCoin"
	accountBalance = decimal.NewFromInt(1000)

	currencyName   = "Dollar"
	currencySymbol = "USD"

	categoryName = "Test category"
	categoryType = models.SPENDING

	transactionDescription = "Test transaction"
	transactionAmount      = decimal.NewFromInt(200)
	transactionDate        = time.Now()
)

var (
	currencyResponse = forms.CurrencyResponse{
		ID:     1,
		Name:   currencyName,
		Symbol: currencySymbol,
	}
	categoryResponse = forms.CategoryResponse{
		ID:   1,
		Name: categoryName,
		Type: categoryType,
	}
	accountResponse = forms.AccountResponse{
		ID:       1,
		Name:     accountName,
		Balance:  accountBalance,
		Currency: &currencyResponse,
	}
	transactionResponse = forms.TransactionResponse{
		ID:          1,
		Description: transactionDescription,
		Amount:      transactionAmount,
		Date:        transactionDate,
		Category:    categoryResponse,
		Account:     accountResponse,
	}
)
