package repository

import (
	"github.com/shopspring/decimal"
	"main/models"
	"time"
)

func GetTransactionByID(id uint) (*models.Transaction, error) {
	return GetModelByID[models.Transaction](id)
}

func GetUserTransactions(userID uint) ([]*models.Transaction, error) {
	return GetUserModels[models.Transaction](userID)
}

func CreateTransaction(
	description string,
	amount decimal.Decimal,
	date time.Time,
	userID uint,
	categoryID uint,
) (*models.Transaction, error) {
	transaction := models.Transaction{
		Description: description,
		Amount:      amount,
		Date:        date,
		UserID:      userID,
		CategoryID:  categoryID,
	}
	return CreateModel(&transaction)
}

func PatchTransaction(
	transactionID uint,
	description string,
	amount decimal.Decimal,
	date time.Time,
	userID uint,
	categoryID uint,
) (*models.Transaction, error) {
	transaction := models.Transaction{
		ID:          transactionID,
		Description: description,
		Amount:      amount,
		Date:        date,
		UserID:      userID,
		CategoryID:  categoryID,
	}
	return PatchModel(&transaction)
}

func DeleteTransaction(id uint) (*models.Transaction, error) {
	return DeleteModel[models.Transaction](id)
}
