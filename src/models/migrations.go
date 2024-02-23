package models

import (
	"fmt"
	"gorm.io/gorm"
)

func MigrateModels(db *gorm.DB) {
	modelsToMigrate := []interface{}{
		&User{},
		&Account{},
		&Currency{},
		&Category{},
		&Transaction{},
	}

	for _, model := range modelsToMigrate {
		err := db.AutoMigrate(model)
		if err != nil {
			panic(fmt.Sprintf("could not migrate model: %v", model))
		}
	}
}
