package model

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// postgresql driver
	_ "github.com/lib/pq"
	// dotenv
	"github.com/joho/godotenv"
)

// DBConnect returns *sql.DB
func DBConnect() (db *sql.DB) {
	err := godotenv.Load(fmt.Sprintf("../%s.env", os.Getenv("GO_ENV")))
	if err != nil {
<<<<<<< HEAD
		// .env読めなかった場合の処理
	}
	env := os.Getenv("ENV")
	DBName := os.Getenv("DB")
	DBDriver := os.Getenv("DB_DRIVAR")
=======
		print("error_env")
	}
	// env := os.Getenv("ENV")
	DBName := os.Getenv("DB")
	DBDriver := os.Getenv("DB_DRIVER")
>>>>>>> 225fbed8e2604733903d14e234c2e48bea9df3b4
	DBUser := os.Getenv("DB_USER")
	DBPass := os.Getenv("DB_PASS")
	DBPort := os.Getenv("DB_PORT")
	DBHost := os.Getenv("DB_HOST")
<<<<<<< HEAD
	db, dberr := sql.Open(DBDriver, "host=”+DBHost+” port="+DBPort+" user="+DBUser+" password="+DBPass+" dbname="+DBName+" sslmode=disable")
=======
	db, dberr := sql.Open(DBDriver, "host="+DBHost+" port="+DBPort+" user="+DBUser+" password="+DBPass+" dbname="+DBName+" sslmode=disable")
>>>>>>> 225fbed8e2604733903d14e234c2e48bea9df3b4
	if dberr != nil {
		log.Fatal(dberr)
	}
	return db
}
