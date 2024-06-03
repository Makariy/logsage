package models

import "github.com/shopspring/decimal"

type Account struct {
	ID         ModelID         `gorm:"column:id;primaryKey;unique;autoIncrement"`
	Name       string          `gorm:"column:name"`
	Balance    decimal.Decimal `gorm:"column:balance;type:numeric"`
	CurrencyID ModelID         `gorm:"column:currency_id"`
	Currency   Currency        `gorm:"foreignKey:CurrencyID"`
	UserID     ModelID         `gorm:"column:user_id"`
	User       User            `gorm:"foreignKey:UserID"`
}

func (Account) TableName() string {
	return "account"
}

func (account Account) GetUser() *User {
	return &account.User
}

func (account *Account) SetUser(user *User) {
	account.UserID = user.ID
}
