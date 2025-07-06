package main

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		fmt.Println("ERROR: TEST_DATABASE_URL not set")
		return
	}
	
	fmt.Printf("Connecting to: %s\n", dsn)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("ERROR opening DB: %v\n", err)
		return
	}
	defer db.Close()
	
	if err := db.Ping(); err != nil {
		fmt.Printf("ERROR pinging DB: %v\n", err)
		return
	}
	
	fmt.Println("Database connection successful")
}