package db_connector

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
)

var (
	DBHost     = os.Getenv("DB_HOST")
	DBUser     = os.Getenv("DB_USER")
	DBName     = os.Getenv("DB_NAME")
	DBPassword = os.Getenv("DB_PASS")
	DBPort     = os.Getenv("DB_PORT")
)

var db *gorm.DB

func getConnConfig(dbName string) gorm.Dialector {
	port, err := strconv.Atoi(DBPort)
	if err != nil {
		errLog := fmt.Sprintf("Could not recognize port: %v", port)
		panic(errLog)
	}
	connConfig := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		DBHost, DBUser, DBPassword, dbName, port,
	)
	return postgres.Open(connConfig)
}

func Connect(dbName string) {
	newConnection, err := gorm.Open(getConnConfig(dbName), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		errLog := fmt.Sprintf(
			"Could not open connection to database: %s:%s %s : %v",
			DBHost, DBPort, DBName, err,
		)
		panic(errLog)
	}
	db = newConnection
}

func GetConnection() *gorm.DB {
	return db
}

func init() {
	Connect(DBName)
}
