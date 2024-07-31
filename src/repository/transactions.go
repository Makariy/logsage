package repository

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"main/db_connector"
	"main/models"
	"time"
)

func preloadAccountToTransaction(transaction *models.Transaction) (*models.Transaction, error) {
	account, err := GetModelByID[models.Account](transaction.AccountID)
	if err != nil {
		return nil, err
	}
	transaction.Account = *account
	return transaction, nil
}

func GetTransactionByID(id models.ModelID) (*models.Transaction, error) {
	transaction, err := GetModelByID[models.Transaction](id)
	if err != nil {
		return nil, err
	}
	return preloadAccountToTransaction(transaction)
}

func GetUserTransactions(userID models.ModelID) ([]*models.Transaction, error) {
	db := db_connector.GetConnection()

	var transactions []*models.Transaction

	tx := db.Preload(clause.Associations).
		Preload("Account."+clause.Associations).
		Preload("Category."+clause.Associations).
		Find(&transactions, "user_id = ?", userID)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return transactions, nil
}

func CreateTransaction(
	description string,
	amount decimal.Decimal,
	date time.Time,
	userID models.ModelID,
	categoryID models.ModelID,
	accountID models.ModelID,
) (*models.Transaction, error) {
	transaction := models.Transaction{
		Description: description,
		Amount:      amount,
		Date:        date,
		UserID:      userID,
		CategoryID:  categoryID,
		AccountID:   accountID,
	}
	result, err := CreateModel(&transaction)
	if err != nil {
		return nil, err
	}
	return preloadAccountToTransaction(result)
}

func PatchTransaction(
	transactionID models.ModelID,
	description string,
	amount decimal.Decimal,
	date time.Time,
	userID models.ModelID,
	categoryID models.ModelID,
	accountID models.ModelID,
) (*models.Transaction, error) {
	transaction := models.Transaction{
		ID:          transactionID,
		Description: description,
		Amount:      amount,
		Date:        date,
		UserID:      userID,
		CategoryID:  categoryID,
		AccountID:   accountID,
	}
	result, err := PatchModel(&transaction)
	if err != nil {
		return nil, err
	}
	return preloadAccountToTransaction(result)
}

func DeleteTransaction(id models.ModelID) (*models.Transaction, error) {
	result, err := DeleteModel[models.Transaction](id)
	if err != nil {
		return nil, err
	}
	return preloadAccountToTransaction(result)
}

func getTransactionsBaseQuery(db *gorm.DB, fromDate, toDate time.Time) *gorm.DB {
	return db.Model(models.Transaction{}).
		Where("? <= date", fromDate).
		Where("date <= ?", toDate).
		Preload(clause.Associations).
		Preload("Account." + clause.Associations).
		Preload("Category." + clause.Associations)
}

func GetUserTransactionsByDate(
	userID models.ModelID,
	fromDate, toDate time.Time,
) ([]*models.Transaction, error) {
	db := db_connector.GetConnection()

	var transactions []*models.Transaction
	tx := getTransactionsBaseQuery(db, fromDate, toDate).
		Where("user_id = ?", userID).
		Order("date desc").
		Find(&transactions)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return transactions, nil
}

func GetCategoryTransactionsByDate(
	categoryID models.ModelID,
	fromDate, toDate time.Time,
) ([]*models.Transaction, error) {
	db := db_connector.GetConnection()

	var transactions []*models.Transaction
	tx := getTransactionsBaseQuery(db, fromDate, toDate).
		Where("category_id = ?", categoryID).
		Find(&transactions)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return transactions, nil
}

func GetAccountTransactionByDate(
	accountID models.ModelID,
	fromDate, toDate time.Time,
) ([]*models.Transaction, error) {
	db := db_connector.GetConnection()

	var transactions []*models.Transaction
	tx := getTransactionsBaseQuery(db, fromDate, toDate).
		Where("account_id = ?", accountID).
		Find(&transactions)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return transactions, nil
}

func CreateTransactionFromModel(transaction *models.Transaction) (*models.Transaction, error) {
	return CreateTransaction(
		transaction.Description,
		transaction.Amount,
		transaction.Date,
		transaction.UserID,
		transaction.CategoryID,
		transaction.AccountID,
	)
}

func PatchTransactionFromModel(transaction *models.Transaction) (*models.Transaction, error) {
	return PatchTransaction(
		transaction.ID,
		transaction.Description,
		transaction.Amount,
		transaction.Date,
		transaction.User.ID,
		transaction.Category.ID,
		transaction.Account.ID,
	)
}
