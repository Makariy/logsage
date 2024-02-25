package repository

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
	"main/db_connector"
	"main/models"
	"time"
)

func loadAccountToTransaction(transaction *models.Transaction) (*models.Transaction, error) {
	account, err := GetModelByID[models.Account](transaction.AccountID)
	if err != nil {
		return nil, err
	}
	transaction.Account = *account
	return transaction, nil
}

func GetTransactionByID(id uint) (*models.Transaction, error) {
	transaction, err := GetModelByID[models.Transaction](id)
	if err != nil {
		return nil, err
	}
	return loadAccountToTransaction(transaction)
}

func GetUserTransactions(userID uint) ([]*models.Transaction, error) {
	db := db_connector.GetConnection()

	var transactions []*models.Transaction

	tx := db.Preload(clause.Associations).
		Preload("Account."+clause.Associations).
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
	userID uint,
	categoryID uint,
	accountID uint,
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
	return loadAccountToTransaction(result)
}

func PatchTransaction(
	transactionID uint,
	description string,
	amount decimal.Decimal,
	date time.Time,
	userID uint,
	categoryID uint,
	accountID uint,
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
	return loadAccountToTransaction(result)
}

func DeleteTransaction(id uint) (*models.Transaction, error) {
	result, err := DeleteModel[models.Transaction](id)
	if err != nil {
		return nil, err
	}
	return loadAccountToTransaction(result)
}

func GetTransactionsByDate(fromDate, toDate time.Time) ([]models.Transaction, error) {
	db := db_connector.GetConnection()

	var transactions []models.Transaction
	tx := db.Model(models.Transaction{}).
		Where("? <= date <= ?", fromDate, toDate).
		Select(&transactions)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return transactions, nil
}
