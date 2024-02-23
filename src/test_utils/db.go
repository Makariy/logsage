package test_utils

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"main/db_connector"
)

func getTestDBName() string {
	return fmt.Sprintf("test_%s", db_connector.DBName)
}

func databaseExists(dbName string) (bool, error) {
	db := db_connector.GetConnection()

	var count int64
	result := db.Raw("SELECT count(*) FROM pg_database WHERE datname = ?", dbName).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func disconnectDB(db *gorm.DB) {
	conn, err := db.DB()
	if err != nil {
		panic("could not get test database connection")
	}
	err = conn.Close()
	if err != nil {
		panic("could not close test database connection")
	}
}

func CreateTestDB() {
	db := db_connector.GetConnection()

	testDBName := getTestDBName()

	exists, err := databaseExists(testDBName)
	if err != nil {
		panic("could not check if the test database exists")
	}
	if exists {
		log.Println(fmt.Sprintf("Dropping test database to recreate"))
		DropTestDB()
	}

	tx := db.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName))
	if tx.Error != nil {
		panic("could not create test database")
	}
	db_connector.Connect(testDBName)
}

func DropTestDB() {
	disconnectDB(db_connector.GetConnection())

	db_connector.Connect(db_connector.DBName)
	db := db_connector.GetConnection()

	testDBName := getTestDBName()
	tx := db.Exec(fmt.Sprintf("DROP DATABASE %s", testDBName))
	if tx.Error != nil {
		panic("could not drop test database")
	}
}
