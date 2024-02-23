package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"main/db_connector"
)

func CreateModel[T any](model *T) (*T, error) {
	db := db_connector.GetConnection()

	err := db.Transaction(func(tx *gorm.DB) error {
		err := db.Create(model).Error
		if err != nil {
			return err
		}
		err = db.Preload(clause.Associations).Find(model).Error
		return err
	})

	if err != nil {
		return nil, err
	}
	return model, nil
}

func PatchModel[T any](model *T) (*T, error) {
	db := db_connector.GetConnection()

	err := db.Transaction(func(tx *gorm.DB) error {
		err := db.Save(model).Error
		if err != nil {
			return err
		}
		err = db.Preload(clause.Associations).Find(model).Error
		return err
	})

	if err != nil {
		return nil, err
	}

	return model, nil
}

func GetModelByID[T any](id uint) (*T, error) {
	db := db_connector.GetConnection()

	var model T
	tx := db.Preload(clause.Associations).Find(&model, "id = ?", id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &model, nil
}

func GetAllModels[T any]() ([]*T, error) {
	db := db_connector.GetConnection()

	var items []*T
	tx := db.Find(&items)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return items, nil
}

func GetUserModels[T any](userId uint) ([]*T, error) {
	db := db_connector.GetConnection()

	var items []*T
	tx := db.Preload(clause.Associations).Find(&items, "user_id = ?", userId)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return items, nil
}

func DeleteModel[T any](id uint) (*T, error) {
	db := db_connector.GetConnection()

	var model T

	err := db.Transaction(func(tx *gorm.DB) error {
		err := db.Preload(clause.Associations).Find(&model, "id = ?", id).Error
		if err != nil {
			return err
		}
		err = db.Delete(&model).Error
		return err
	})

	if err != nil {
		return nil, err
	}

	return &model, nil
}
