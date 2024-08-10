package models

import "github.com/shopspring/decimal"

type Currency struct {
	ID     ModelID         `gorm:"column:id;primaryKey;unique;autoIncrement"`
	Name   string          `gorm:"column:name"`
	Symbol string          `gorm:"column:symbol"`
	Value  decimal.Decimal `gorm:"column:value"`
}

func (_ *Currency) TableName() string {
	return "currency"
}
