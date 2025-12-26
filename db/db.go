package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDB() *sql.DB {
	dsn := "root@tcp(localhost:3306)/rmsk?parseTime=true"
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Ping(); err != nil {
		log.Fatal(err)
	}

	return database
}
