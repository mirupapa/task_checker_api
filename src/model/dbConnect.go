package model

import (
	"database/sql"
	"log"
	"os"

	// postgresql driver
	_ "github.com/lib/pq"
	// dotenv
)

// DBConnect returns *sql.DB
func DBConnect() (db *sql.DB) {
	// env := os.Getenv("ENV")
	DBName := os.Getenv("DB_NAME")
	DBDriver := os.Getenv("DB_DRIVER")
	DBUser := os.Getenv("DB_USER")
	DBPass := os.Getenv("DB_PASS")
	DBPort := os.Getenv("DB_PORT")
	DBHost := os.Getenv("DB_HOST")
	db, dberr := sql.Open(DBDriver, "host="+DBHost+" port="+DBPort+" user="+DBUser+" password="+DBPass+" dbname="+DBName+" sslmode=disable")
	if dberr != nil {
		log.Fatal(dberr)
	}
	return db
}
