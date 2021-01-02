package model

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// postgresql driver
	_ "github.com/lib/pq"
	// dotenv
)

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

// DBConnect returns *sql.DB
func DBConnect() (db *sql.DB) {
	var (
		DBUser   = mustGetenv("DB_USER")
		DBPass   = mustGetenv("DB_PASS")
		DBHost   = mustGetenv("DB_HOST")
		DBName   = mustGetenv("DB_NAME")
		DBPort   = mustGetenv("DB_PORT")
		DBDriver = mustGetenv("DB_DRIVER")
	)
	ENV := os.Getenv("ENV")
	fmt.Println("env:" + ENV)
	if ENV == "development" {
		db, dberr := sql.Open(DBDriver, "host="+DBHost+" port="+DBPort+" user="+DBUser+" password="+DBPass+" dbname="+DBName+" sslmode=disable")
		if dberr != nil {
			log.Fatal(dberr)
		}
		return db
	}

	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	var dbURI string
	dbURI = fmt.Sprintf("user=%s password=%s database=%s host=%s/%s", DBUser, DBPass, DBName, socketDir, DBHost)
	fmt.Println("dbURI:" + dbURI)
	dbPool, err := sql.Open(DBDriver, dbURI)
	if err != nil {
		log.Fatal(err)
		fmt.Println("dbErr")
		return nil
	}

	configureConnectionPool(dbPool)

	return dbPool
}

// configureConnectionPool sets database connection pool properties.
// For more information, see https://golang.org/pkg/database/sql
func configureConnectionPool(dbPool *sql.DB) {
	dbPool.SetMaxIdleConns(5)
	dbPool.SetMaxOpenConns(7)
	dbPool.SetConnMaxLifetime(1800)
}
