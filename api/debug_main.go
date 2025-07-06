package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/necofuryai/bocchi-the-map/api/application/clients"
	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
)

func main() {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		fmt.Println("ERROR: TEST_DATABASE_URL not set")
		return
	}
	
	fmt.Printf("Step 1: Connecting to: %s\n", dsn)
	
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
	
	fmt.Println("Step 2: Database connection successful")
	
	// Test creating Queries
	fmt.Println("Step 3: Creating database.Queries")
	queries := database.New(db)
	if queries == nil {
		fmt.Println("ERROR: database.New returned nil")
		return
	}
	fmt.Println("Step 4: Queries created successfully")
	
	// Test creating UserClient
	fmt.Println("Step 5: Creating UserClient")
	userClient, err := clients.NewUserClient("internal", db)
	if err != nil {
		fmt.Printf("ERROR creating UserClient: %v\n", err)
		return
	}
	if userClient == nil {
		fmt.Println("ERROR: NewUserClient returned nil")
		return
	}
	
	fmt.Println("Step 6: UserClient created successfully")
	
	fmt.Println("All components created successfully!")
}