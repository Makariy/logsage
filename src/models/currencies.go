package models

type Currency struct {
	ID     ModelID `gorm:"column:id;primaryKey;unique;autoIncrement"`
	Name   string  `gorm:"column:name"`
	Symbol string  `gorm:"column:symbol"`
}

func (_ *Currency) TableName() string {
	return "currency"
}
