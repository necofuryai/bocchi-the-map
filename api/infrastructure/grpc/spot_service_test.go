package grpc

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestSpotService_CreateSpot(t *testing.T) {
	// Create test database connection (using test database)
	// Note: For actual testing, you would use a test database or mock
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/bocchi_test?parseTime=true"
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Skip("Skipping test: test database not available")
	}
	defer db.Close()
	
	// Skip if database connection fails (test environment may not have DB)
	if err := db.Ping(); err != nil {
		t.Skip("Skipping test: cannot connect to test database")
	}

	service := NewSpotService(db)

	tests := []struct {
		name    string
		req     *CreateSpotRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &CreateSpotRequest{
				Name: "Test Cafe",
				Coordinates: &Coordinates{
					Latitude:  35.6762,
					Longitude: 139.6503,
				},
				Category:    "cafe",
				Address:     "Tokyo, Japan",
				CountryCode: "JP",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			req: &CreateSpotRequest{
				Coordinates: &Coordinates{
					Latitude:  35.6762,
					Longitude: 139.6503,
				},
				Category:    "cafe",
				CountryCode: "JP",
			},
			wantErr: true,
		},
		{
			name: "missing coordinates",
			req: &CreateSpotRequest{
				Name:        "Test Cafe",
				Category:    "cafe",
				CountryCode: "JP",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.CreateSpot(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSpot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil || resp.Spot == nil {
					t.Error("CreateSpot() returned nil response or spot")
					return
				}
				if resp.Spot.Name != tt.req.Name {
					t.Errorf("CreateSpot() spot name = %v, want %v", resp.Spot.Name, tt.req.Name)
				}
			}
		})
	}
}

func TestSpotService_GetSpot(t *testing.T) {
	// Create test database connection (using test database)
	// Note: For actual testing, you would use a test database or mock
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/bocchi_test?parseTime=true"
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Skip("Skipping test: test database not available")
	}
	defer db.Close()
	
	// Skip if database connection fails (test environment may not have DB)
	if err := db.Ping(); err != nil {
		t.Skip("Skipping test: cannot connect to test database")
	}

	service := NewSpotService(db)

	tests := []struct {
		name    string
		req     *GetSpotRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     &GetSpotRequest{ID: "test_id"},
			wantErr: false,
		},
		{
			name:    "empty id",
			req:     &GetSpotRequest{ID: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.GetSpot(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSpot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil || resp.Spot == nil {
					t.Error("GetSpot() returned nil response or spot")
					return
				}
				if resp.Spot.ID != tt.req.ID {
					t.Errorf("GetSpot() spot ID = %v, want %v", resp.Spot.ID, tt.req.ID)
				}
			}
		})
	}
}

func TestSpotService_ListSpots(t *testing.T) {
	// Create test database connection (using test database)
	// Note: For actual testing, you would use a test database or mock
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/bocchi_test?parseTime=true"
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Skip("Skipping test: test database not available")
	}
	defer db.Close()
	
	// Skip if database connection fails (test environment may not have DB)
	if err := db.Ping(); err != nil {
		t.Skip("Skipping test: cannot connect to test database")
	}

	service := NewSpotService(db)

	req := &ListSpotsRequest{
		Pagination: &PaginationRequest{
			Page:     1,
			PageSize: 20,
		},
	}

	resp, err := service.ListSpots(context.Background(), req)
	if err != nil {
		t.Errorf("ListSpots() error = %v", err)
		return
	}

	if resp == nil {
		t.Error("ListSpots() returned nil response")
		return
	}

	if len(resp.Spots) == 0 {
		t.Error("ListSpots() returned empty spots list")
	}

	if resp.Pagination == nil {
		t.Error("ListSpots() returned nil pagination")
	}
}

func TestSpotService_SearchSpots(t *testing.T) {
	// Create test database connection (using test database)
	// Note: For actual testing, you would use a test database or mock
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/bocchi_test?parseTime=true"
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Skip("Skipping test: test database not available")
	}
	defer db.Close()
	
	// Skip if database connection fails (test environment may not have DB)
	if err := db.Ping(); err != nil {
		t.Skip("Skipping test: cannot connect to test database")
	}

	service := NewSpotService(db)

	tests := []struct {
		name    string
		req     *SearchSpotsRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &SearchSpotsRequest{
				Query: "cafe",
				Pagination: &PaginationRequest{
					Page:     1,
					PageSize: 20,
				},
			},
			wantErr: false,
		},
		{
			name: "empty query",
			req: &SearchSpotsRequest{
				Query: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.SearchSpots(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchSpots() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("SearchSpots() returned nil response")
					return
				}
				if len(resp.Spots) == 0 {
					t.Error("SearchSpots() returned empty spots list")
				}
			}
		})
	}
}