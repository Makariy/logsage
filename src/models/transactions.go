package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	ID          uint            `gorm:"column:id;primaryKey;unique;autoIncrement"`
	Description string          `gorm:"column:description"`
	Amount      decimal.Decimal `gorm:"column:balance"`
	Date        time.Time       `gorm:"column:date"`
	UserID      uint            `gorm:"column:user_id"`
	User        User            `gorm:"foreignKey:UserID"`
	CategoryID  uint            `gorm:"column:category_id"`
	Category    Category        `gorm:"foreignKey:CategoryID"`
}

func (Transaction) TableName() string {
	return "transaction"
}

func (transaction Transaction) GetUser() *User {
	return &transaction.User
}
