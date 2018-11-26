package datastore

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

func GetDBConnectionString() string {
	dbUsername := os.Getenv("POSTGRES_DB_USER")
	dbPassword := os.Getenv("POSTGRES_DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOSTNAME")
	dbPort := os.Getenv("DB_PORT")

	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUsername, dbPassword, dbName, dbHost, dbPort)
}

func InitDB() error {
	var err error
	dbConnStr := GetDBConnectionString()

	db, err = sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Panic(err)
		return err
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
		return err
	}

	return nil
}
