package repository

import (
	"fmt"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"main/db_connector"
	"main/forms"
	"main/models"
	"time"
)

func queryTransactions(db *gorm.DB) *gorm.DB {
	return db.Model(models.Transaction{})
}

func queryTransactionsByDate(db *gorm.DB, fromDate, toDate time.Time) *gorm.DB {
	return db.Model(models.Transaction{}).
		Where("? <= date", fromDate).
		Where("date <= ?", toDate)
}

func queryJoinCategory(db *gorm.DB) *gorm.DB {
	return db.Joins("LEFT JOIN category ON category_id = category.id")
}

func queryJoinCategoryByType(
	db *gorm.DB,
	categoryType models.CategoryType,
) *gorm.DB {
	return queryJoinCategory(
		queryTransactions(db),
	).Where("category.type = ?", categoryType)
}

func queryJoinAccountAndCurrency(db *gorm.DB) *gorm.DB {
	return db.Joins("LEFT JOIN account ON transaction.account_id = account.id").
		Joins("LEFT JOIN currency ON account.currency_id = currency.id")
}

func getTotalTransactionsAmount(
	db *gorm.DB,
) (decimal.Decimal, error) {
	var totalPrice decimal.Decimal

	tx := queryJoinAccountAndCurrency(
		queryTransactions(db),
	).
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
		queryTransactionsByDate(db, fromDate, toDate).
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
		TotalAmount:  ConvertToRelativeCurrency(totalCategoryAmount, outputCurrency),
		Transactions: transactions,
	}
	return &stats, nil
}

func getTotalTransactedAmountByType(
	db *gorm.DB,
	fromDate,
	toDate time.Time,
	categoryType models.CategoryType,
) (decimal.Decimal, error) {
	return getTotalTransactionsAmount(
		queryJoinCategoryByType(
			queryTransactionsByDate(db, fromDate, toDate),
			categoryType,
		),
	)
}

func getAccountTransactedAmountByType(
	db *gorm.DB,
	fromDate,
	toDate time.Time,
	accountID models.ModelID,
	categoryType models.CategoryType,
) (amount decimal.Decimal, err error) {
	return getTotalTransactionsAmount(
		queryJoinCategoryByType(
			queryTransactionsByDate(db, fromDate, toDate),
			categoryType,
		).
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
		TotalEarnedAmount: ConvertToRelativeCurrency(totalEarnedAmount, outputCurrency),
		TotalSpentAmount:  ConvertToRelativeCurrency(totalSpentAmount, outputCurrency),
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

type rawIntervalStat struct {
	TotalEarnedAmount decimal.Decimal
	TotalSpentAmount  decimal.Decimal
	Date              time.Time
}

func getSelectSumByCategoryType(categoryType models.CategoryType) string {
	return "COALESCE(" +
		fmt.Sprintf("SUM(case when category.type = '%s' ", categoryType) +
		"then transaction.amount::numeric * currency.value::numeric " +
		"else 0 end), " +
		"0)"
}

func renderQueryInterval(interval models.TimeIntervalStep) string {
	return fmt.Sprintf("'1 %s'::interval", string(interval))
}

func queryTimeIntervalSeries(
	db *gorm.DB,
	interval models.TimeIntervalStep,
	fromDate time.Time,
	toDate time.Time,
) *gorm.DB {
	return db.Table(""+
		fmt.Sprintf(
			"generate_series("+
				"?::date + '1 second'::interval, "+
				"?::date - '1 second'::interval, "+
				"%s) as date_start",
			renderQueryInterval(interval),
		), fromDate, toDate,
	)
}

func queryJoinTransactionsByUserInInterval(
	db *gorm.DB,
	interval models.TimeIntervalStep,
	userID models.ModelID,
) *gorm.DB {
	return db.Joins(
		"LEFT JOIN transaction ON date_start <= transaction.date "+
			"AND "+
			fmt.Sprintf(
				"transaction.date < date_start + %s",
				renderQueryInterval(interval),
			)+
			" AND transaction.user_id = ?", userID,
	)
}

func getSelectTotalAmounts() []string {
	return []string{
		getSelectSumByCategoryType(models.EARNING) +
			" as total_earned_amount",
		getSelectSumByCategoryType(models.SPENDING) +
			" as total_spent_amount",
	}
}

func getRawTimeIntervalStats(
	userID models.ModelID,
	fromDate,
	toDate time.Time,
	step models.TimeIntervalStep,
) ([]rawIntervalStat, error) {
	db := db_connector.GetConnection()

	var result []rawIntervalStat

	tx := queryJoinAccountAndCurrency(
		queryJoinCategory(
			queryJoinTransactionsByUserInInterval(
				queryTimeIntervalSeries(db, step, fromDate, toDate),
				step,
				userID,
			),
		),
	).
		Select(append(getSelectTotalAmounts(), "date_start as date")).
		Group("date_start").
		Order("date_start ASC").
		Find(&result)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return result, nil
}

func getToDateForStat(
	rawStats []rawIntervalStat,
	index int,
	step models.TimeIntervalStep,
) int64 {
	if index >= len(rawStats)-1 {
		return rawStats[index].Date.Unix() + models.ConvertTimeStepToTime(step)
	}

	return rawStats[index+1].Date.Unix()
}

func getRawStatDateRangeByIndex(
	rawStats []rawIntervalStat,
	index int,
	step models.TimeIntervalStep,
) *forms.DateRange {
	stat := rawStats[index]
	return &forms.DateRange{
		FromDate: stat.Date.Unix(),
		ToDate:   getToDateForStat(rawStats, index, step),
	}
}

func parseRawIntervalStat(
	rawStats []rawIntervalStat,
	index int,
	step models.TimeIntervalStep,
	outputCurrency *models.Currency,
) *forms.TimeIntervalStat {
	rawStat := rawStats[index]
	return &forms.TimeIntervalStat{
		TotalSpentAmount:  ConvertToRelativeCurrency(rawStat.TotalSpentAmount, outputCurrency),
		TotalEarnedAmount: ConvertToRelativeCurrency(rawStat.TotalEarnedAmount, outputCurrency),
		DateRange:         getRawStatDateRangeByIndex(rawStats, index, step),
	}
}

func parseRawIntervalStats(
	rawStats []rawIntervalStat,
	fromDate,
	toDate time.Time,
	step models.TimeIntervalStep,
	outputCurrency *models.Currency,
) *forms.TimeIntervalStats {
	result := forms.TimeIntervalStats{
		DateRange: &forms.DateRange{
			FromDate: fromDate.Unix(),
			ToDate:   toDate.Unix(),
		},
		TimeStep: models.ConvertTimeStepToTime(step),
	}

	for index := range rawStats {
		result.IntervalStats = append(
			result.IntervalStats,
			parseRawIntervalStat(
				rawStats,
				index,
				step,
				outputCurrency,
			),
		)
	}
	return &result
}

func GetTimeIntervalStats(
	userID models.ModelID,
	fromDate,
	toDate time.Time,
	step models.TimeIntervalStep,
	outputCurrency *models.Currency,
) (*forms.TimeIntervalStats, error) {
	rawStats, err := getRawTimeIntervalStats(
		userID,
		fromDate,
		toDate,
		step,
	)
	if err != nil {
		return nil, err
	}

	stats := parseRawIntervalStats(
		rawStats,
		fromDate,
		toDate,
		step,
		outputCurrency,
	)

	return stats, nil
}
