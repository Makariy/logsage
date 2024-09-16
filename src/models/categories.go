package models

type CategoryType string

const (
	SPENDING CategoryType = "SPENDING"
	EARNING  CategoryType = "EARNING"
)

type Category struct {
	ID              ModelID       `gorm:"column:id;primaryKey;unique;autoIncrement"`
	Name            string        `gorm:"column:name;unique"`
	Type            CategoryType  `gorm:"column:type"`
	UserID          ModelID       `gorm:"column:user_id"`
	User            User          `gorm:"foreignKey:UserID"`
	CategoryImageID ModelID       `gorm:"column:category_image_id"`
	CategoryImage   CategoryImage `gorm:"foreignKey:CategoryImageID"`
}

func (Category) TableName() string {
	return "category"
}

func (category Category) GetUser() *User {
	return &category.User
}
func (category *Category) SetUser(user *User) {
	category.UserID = user.ID
}
