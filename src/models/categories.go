package models

const (
	SPENDING = "SPENDING"
	EARNING  = "EARNING"
)

type Category struct {
	ID     uint   `gorm:"column:id;primaryKey;unique;autoIncrement"`
	Name   string `gorm:"column:name;unique"`
	Type   string `gorm:"column:type"`
	UserID uint   `gorm:"column:user_id"`
	User   User   `gorm:"foreignKey:UserID"`
}

func (Category) TableName() string {
	return "category"
}

func (category Category) GetUser() *User {
	return &category.User
}
