package repository

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"main/db_connector"
	"main/models"
	"time"
)

func getStatsBaseQuery(db *gorm.DB, fromDate, toDate time.Time) *gorm.DB {
	return db.Model(models.Transaction{}).
		Where("? <= date", fromDate).
		Where("date <= ?", toDate)
}

func getStatsQueryJoinCategory(db *gorm.DB, fromDate, toDate time.Time) *gorm.DB {
	return getStatsBaseQuery(db, fromDate, toDate).
		Joins("inner join category on category_id = category.id")
}

func GetCategoryStats(categoryID uint, fromDate, toDate time.Time) (*models.CategoryStats, error) {
	db := db_connector.GetConnection()
	var (
		totalCategoryAmount, totalAmount decimal.Decimal
	)

	category, err := GetCategoryByID(categoryID)
	if err != nil {
		return nil, err
	}

	tx := getStatsBaseQuery(db, fromDate, toDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalAmount)
	if tx.Error != nil {
		return nil, tx.Error
	}
	tx = getStatsBaseQuery(db, fromDate, toDate).
		Where("category_id = ?", categoryID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalCategoryAmount)
	if tx.Error != nil {
		return nil, tx.Error
	}

	transactions, err := GetCategoryTransactionsByDate(categoryID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	totalPercent := decimal.NewFromInt(0)
	if !totalAmount.Equal(decimal.NewFromInt(0)) {
		totalPercent = totalCategoryAmount.Div(totalAmount)
	}
	stats := models.CategoryStats{
		Category:     *category,
		TotalAmount:  totalCategoryAmount,
		TotalPercent: totalPercent,
		Transactions: transactions,
	}
	return &stats, nil
}

func getStatsQueryJoinCategoryByType(db *gorm.DB, fromDate, toDate time.Time, categoryType string) *gorm.DB {
	return getStatsQueryJoinCategory(db, fromDate, toDate).
		Where("category.type = ?", categoryType)
}

func getTotalTransactedAmount(
	db *gorm.DB,
	fromDate,
	toDate time.Time,
	categoryType string,
) (amount decimal.Decimal, err error) {
	err = getStatsQueryJoinCategoryByType(db, fromDate, toDate, categoryType).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&amount).Error
	return
}

func getAccountTransactedAmount(
	db *gorm.DB,
	fromDate,
	toDate time.Time,
	accountID uint,
	categoryType string,
) (amount decimal.Decimal, err error) {
	err = getStatsQueryJoinCategoryByType(db, fromDate, toDate, categoryType).
		Where("account_id = ?", accountID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&amount).Error
	return
}

func GetAccountStats(accountID uint, fromDate, toDate time.Time) (*models.AccountStats, error) {
	db := db_connector.GetConnection()

	account, err := GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	totalEarnedAmount, err := getTotalTransactedAmount(db, fromDate, toDate, models.EARNING)
	totalSpentAmount, err := getTotalTransactedAmount(db, fromDate, toDate, models.SPENDING)

	accountEarnedAmount, err := getAccountTransactedAmount(db, fromDate, toDate, accountID, models.EARNING)
	if err != nil {
		return nil, err
	}

	accountSpentAmount, err := getAccountTransactedAmount(db, fromDate, toDate, accountID, models.SPENDING)
	if err != nil {
		return nil, err
	}

	transactions, err := GetAccountTransactionByDate(account.ID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	stats := models.AccountStats{
		Account:            *account,
		TotalEarnedAmount:  accountEarnedAmount,
		TotalEarnedPercent: accountEarnedAmount.Div(totalEarnedAmount),
		TotalSpentAmount:   accountSpentAmount,
		TotalSpentPercent:  accountSpentAmount.Div(totalSpentAmount),
		Transactions:       transactions,
	}

	return &stats, nil
}

func GetAccountsStats(userID uint, fromDate, toDate time.Time) ([]*models.AccountStats, error) {
	accounts, err := GetUserAccounts(userID)
	if err != nil {
		return nil, err
	}

	accountsStats := make([]*models.AccountStats, 0, len(accounts))
	for _, account := range accounts {
		accountStats, err := GetAccountStats(account.ID, fromDate, toDate)
		if err != nil {
			return nil, err
		}
		accountsStats = append(accountsStats, accountStats)
	}
	return accountsStats, nil
}

func GetCategoriesStats(userID uint, fromDate, toDate time.Time) ([]*models.CategoryStats, error) {
	categories, err := GetUserCategories(userID)
	if err != nil {
		return nil, err
	}

	categoriesStats := make([]*models.CategoryStats, 0, len(categories))
	for _, category := range categories {
		categoryStats, err := GetCategoryStats(category.ID, fromDate, toDate)
		if err != nil {
			return nil, err
		}
		categoriesStats = append(categoriesStats, categoryStats)
	}

	return categoriesStats, nil
}

func GetTotalStats(userID uint, fromDate, toDate time.Time) (*models.TotalStats, error) {
	db := db_connector.GetConnection()

	totalEarnedAmount, err := getTotalTransactedAmount(db, fromDate, toDate, models.EARNING)
	if err != nil {
		return nil, err
	}
	totalSpentAmount, err := getTotalTransactedAmount(db, fromDate, toDate, models.SPENDING)
	if err != nil {
		return nil, err
	}

	accountsStats, err := GetAccountsStats(userID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	categoriesStats, err := GetCategoriesStats(userID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	stats := models.TotalStats{
		TotalEarnedAmount: totalEarnedAmount,
		TotalSpentAmount:  totalSpentAmount,
		AccountsStats:     accountsStats,
		CategoriesStats:   categoriesStats,
	}
	return &stats, nil
}
