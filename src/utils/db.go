package utils

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB *sql.DB
)

func ConnectDB() (*sql.DB, error) {
	log.Println("Connecting to database connection...")
	connection := strings.Join([]string{
		GetEnv("DB_USER", ""),
		":",
		GetEnv("DB_PASS", ""),
		"@tcp(",
		GetEnv("DB_HOST", ""),
		":",
		GetEnv("DB_PORT", ""),
		")/",
		GetEnv("DB_NAME", ""),
	}, "")
	log.Println("Connecting to database...")
	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Println("Error connecting to database")
		log.Println("Error:", err)
		return nil, err
	}
	log.Println("Pinging database...")
	err = db.Ping()
	if err != nil {
		log.Println("Error pinging database")
		log.Println("Error:", err)
		return nil, err
	}
	log.Println("Connected to database")
	return db, nil
}
