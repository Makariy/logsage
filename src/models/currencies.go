package models

type Currency struct {
	ID   uint   `gorm:"column:id;primaryKey;unique;autoIncrement"`
	Name string `gorm:"column:name"`
}

func (_ *Currency) TableName() string {
	return "currency"
}
