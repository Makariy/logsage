package routes_tests

import (
	"main/forms"
	data "main/test_utils/test_data"
)

var (
	currencyResponse = forms.CurrencyResponse{
		ID:     1,
		Name:   data.FirstCurrencySymbol,
		Symbol: data.FirstCurrencySymbol,
	}
	categoryResponse = forms.CategoryResponse{
		ID:   1,
		Name: data.FirstCategoryName,
		Type: data.FirstCategoryType,
	}
	accountResponse = forms.AccountResponse{
		ID:       1,
		Name:     data.FirstAccountName,
		Balance:  data.FirstAccountBalance,
		Currency: &currencyResponse,
	}
	transactionResponse = forms.TransactionResponse{
		ID:          1,
		Description: data.TransactionDescription,
		Amount:      data.TransactionAmount,
		Date:        data.TransactionDate,
		Category:    categoryResponse,
		Account:     accountResponse,
	}
)
