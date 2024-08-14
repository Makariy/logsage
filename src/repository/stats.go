package repository

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"main/db_connector"
	"main/forms"
	"main/models"
	"time"
)

func getTransactionsByDateRangeBase(db *gorm.DB, fromDate, toDate time.Time) *gorm.DB {
	return db.Model(models.Transaction{}).
		Where("? <= date", fromDate).
		Where("date <= ?", toDate)
}

func getStatsQueryJoinCategory(db *gorm.DB, fromDate, toDate time.Time) *gorm.DB {
	return getTransactionsByDateRangeBase(db, fromDate, toDate).
		Joins("INNER JOIN category ON category_id = category.id")
}

func getStatsQueryJoinCategoryByType(
	db *gorm.DB,
	fromDate,
	toDate time.Time,
	categoryType string,
) *gorm.DB {
	return getStatsQueryJoinCategory(db, fromDate, toDate).
		Where("category.type = ?", categoryType)
}

func getTotalTransactionsAmount(
	db *gorm.DB,
) (decimal.Decimal, error) {
	var totalPrice decimal.Decimal

	tx := db.Joins("INNER JOIN account ON transaction.account_id = account.id").
		Joins("INNER JOIN currency ON account.currency_id = currency.id").
		Select("COALESCE(SUM(transaction.amount::numeric * currency.value::numeric), 0)").
		Scan(&totalPrice)

	if tx.Error != nil {
		return decimal.Zero, tx.Error
	}

	return totalPrice, nil
}

func GetCategoryStats(
	categoryID models.ModelID,
	fromDate,
	toDate time.Time,
	outputCurrency *models.Currency,
) (*forms.CategoryStats, error) {
	db := db_connector.GetConnection()

	category, err := GetCategoryByID(categoryID)
	if err != nil {
		return nil, err
	}

	totalCategoryAmount, err := getTotalTransactionsAmount(
		getTransactionsByDateRangeBase(db, fromDate, toDate).
			Where("category_id = ?", categoryID),
	)
	if err != nil {
		return nil, err
	}

	transactions, err := GetCategoryTransactionsByDate(categoryID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	stats := forms.CategoryStats{
		Category:     *category,
		TotalAmount:  totalCategoryAmount.Div(outputCurrency.Value),
		Transactions: transactions,
	}
	return &stats, nil
}

func getTotalTransactedAmountByType(
	db *gorm.DB,
	fromDate,
	toDate time.Time,
	categoryType string,
) (decimal.Decimal, error) {
	return getTotalTransactionsAmount(
		getStatsQueryJoinCategoryByType(db, fromDate, toDate, categoryType),
	)
}

func getAccountTransactedAmountByType(
	db *gorm.DB,
	fromDate,
	toDate time.Time,
	accountID models.ModelID,
	categoryType string,
) (amount decimal.Decimal, err error) {
	return getTotalTransactionsAmount(
		getStatsQueryJoinCategoryByType(db, fromDate, toDate, categoryType).
			Where("account_id = ?", accountID),
	)
}

func GetAccountStats(
	accountID models.ModelID,
	fromDate,
	toDate time.Time,
) (*forms.AccountStats, error) {
	db := db_connector.GetConnection()

	account, err := GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	accountEarnedAmount, err := getAccountTransactedAmountByType(db, fromDate, toDate, accountID, models.EARNING)
	if err != nil {
		return nil, err
	}

	accountSpentAmount, err := getAccountTransactedAmountByType(db, fromDate, toDate, accountID, models.SPENDING)
	if err != nil {
		return nil, err
	}

	transactions, err := GetAccountTransactionByDate(account.ID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	stats := forms.AccountStats{
		Account:           *account,
		TotalEarnedAmount: accountEarnedAmount,
		TotalSpentAmount:  accountSpentAmount,
		Transactions:      transactions,
	}

	return &stats, nil
}

func GetAccountsStats(
	userID models.ModelID,
	fromDate,
	toDate time.Time,
) ([]*forms.AccountStats, error) {
	accounts, err := GetUserAccounts(userID)
	if err != nil {
		return nil, err
	}

	accountsStats := make([]*forms.AccountStats, 0, len(accounts))
	for _, account := range accounts {
		accountStats, err := GetAccountStats(account.ID, fromDate, toDate)
		if err != nil {
			return nil, err
		}
		accountsStats = append(accountsStats, accountStats)
	}
	return accountsStats, nil
}

func GetCategoriesStats(
	userID models.ModelID,
	fromDate,
	toDate time.Time,
	outputCurrency *models.Currency,
) ([]*forms.CategoryStats, error) {
	categories, err := GetUserCategories(userID)
	if err != nil {
		return nil, err
	}

	categoriesStats := make([]*forms.CategoryStats, 0, len(categories))
	for _, category := range categories {
		categoryStats, err := GetCategoryStats(category.ID, fromDate, toDate, outputCurrency)
		if err != nil {
			return nil, err
		}
		categoriesStats = append(categoriesStats, categoryStats)
	}

	return categoriesStats, nil
}

func GetTotalCategoriesStats(
	userID models.ModelID,
	fromDate,
	toDate time.Time,
	outputCurrency *models.Currency,
) (*forms.TotalCategoriesStats, error) {
	db := db_connector.GetConnection()

	totalEarnedAmount, err := getTotalTransactedAmountByType(db, fromDate, toDate, models.EARNING)
	if err != nil {
		return nil, err
	}
	totalSpentAmount, err := getTotalTransactedAmountByType(db, fromDate, toDate, models.SPENDING)
	if err != nil {
		return nil, err
	}

	categoriesStats, err := GetCategoriesStats(userID, fromDate, toDate, outputCurrency)
	if err != nil {
		return nil, err
	}

	stats := forms.TotalCategoriesStats{
		TotalEarnedAmount: totalEarnedAmount.Div(outputCurrency.Value),
		TotalSpentAmount:  totalSpentAmount.Div(outputCurrency.Value),
		CategoriesStats:   categoriesStats,
	}
	return &stats, nil
}

func GetTotalAccountsStats(
	userID models.ModelID,
	fromDate,
	toDate time.Time,
) (*forms.TotalAccountsStats, error) {
	db := db_connector.GetConnection()

	totalEarnedAmount, err := getTotalTransactedAmountByType(db, fromDate, toDate, models.EARNING)
	if err != nil {
		return nil, err
	}
	totalSpentAmount, err := getTotalTransactedAmountByType(db, fromDate, toDate, models.SPENDING)
	if err != nil {
		return nil, err
	}

	accountsStats, err := GetAccountsStats(userID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	stats := forms.TotalAccountsStats{
		TotalEarnedAmount: totalEarnedAmount,
		TotalSpentAmount:  totalSpentAmount,
		AccountsStats:     accountsStats,
	}
	return &stats, nil
}
