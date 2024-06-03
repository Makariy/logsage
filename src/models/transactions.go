package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	ID          ModelID         `gorm:"column:id;primaryKey;unique;autoIncrement"`
	Description string          `gorm:"column:description"`
	Amount      decimal.Decimal `gorm:"column:amount;type:numeric"`
	Date        time.Time       `gorm:"column:date"`
	UserID      ModelID         `gorm:"column:user_id"`
	User        User            `gorm:"foreignKey:UserID"`
	CategoryID  ModelID         `gorm:"column:category_id"`
	Category    Category        `gorm:"foreignKey:CategoryID"`
	AccountID   ModelID         `gorm:"column:account_id"`
	Account     Account         `gorm:"foreignKey:AccountID"`
}

func (Transaction) TableName() string {
	return "transaction"
}

func (transaction Transaction) GetUser() *User {
	return &transaction.User
}
func (transaction Transaction) SetUser(user *User) {
	transaction.UserID = user.ID
}
