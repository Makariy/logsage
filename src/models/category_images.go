package models

type CategoryImage struct {
	ID       ModelID `gorm:"column:id;primaryKey;unique;autoIncrement"`
	FileName string  `gorm:"column:filename"`
}

func (CategoryImage) TableName() string {
	return "category_image"
}
