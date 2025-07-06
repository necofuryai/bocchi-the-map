package main

import (
	"fmt"
	"os"

	"github.com/necofuryai/bocchi-the-map/api/tests/helpers"
)

func main() {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		fmt.Println("ERROR: TEST_DATABASE_URL not set")
		return
	}
	
	fmt.Printf("Step 1: TEST_DATABASE_URL is set: %s\n", dsn)
	
	// Test EnsureTestDatabase
	fmt.Println("Step 2: Running EnsureTestDatabase")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("EnsureTestDatabase panic: %v\n", r)
		}
	}()
	
	helpers.EnsureTestDatabase()
	fmt.Println("Step 3: EnsureTestDatabase completed successfully")
	
	// Test NewTestDatabase
	fmt.Println("Step 4: Creating NewTestDatabase")
	testDB, err := helpers.NewTestDatabase()
	if err != nil {
		fmt.Printf("ERROR creating NewTestDatabase: %v\n", err)
		return
	}
	if testDB == nil {
		fmt.Println("ERROR: NewTestDatabase returned nil")
		return
	}
	defer testDB.Close()
	
	fmt.Println("Step 5: NewTestDatabase created successfully")
	
	// Test NewCommonTestSuite
	fmt.Println("Step 6: Creating NewCommonTestSuite")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("NewCommonTestSuite panic: %v\n", r)
		}
	}()
	
	suite := helpers.NewCommonTestSuite()
	if suite == nil {
		fmt.Println("ERROR: NewCommonTestSuite returned nil")
		return
	}
	defer suite.Cleanup()
	
	fmt.Println("Step 7: NewCommonTestSuite created successfully")
	
	fmt.Println("All test suite components created successfully!")
}