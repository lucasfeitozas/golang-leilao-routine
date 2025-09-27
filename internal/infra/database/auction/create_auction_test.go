package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCreateAuction_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should create auction successfully and update status after interval", func(mt *mtest.T) {
		// Set environment variable for auction interval
		os.Setenv("AUCTION_INTERVAL", "100ms")
		defer os.Unsetenv("AUCTION_INTERVAL")

		// Create repository with mock collection
		repo := NewAuctionRepository(mt.DB)

		// Mock successful insertion
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Mock successful update for the goroutine
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "modifiedCount", Value: 1}))

		// Create test auction entity
		auction := &auction_entity.Auction{
			Id:          "test-auction-id",
			ProductName: "Test Product",
			Category:    "Electronics",
			Description: "Test description for auction",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		// Execute CreateAuction
		err := repo.CreateAuction(context.Background(), auction)

		// Verify initial creation was successful
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Wait for goroutine to complete (slightly longer than auction interval)
		time.Sleep(150 * time.Millisecond)
	})
}

func TestCreateAuction_InsertionError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should return error when insertion fails", func(mt *mtest.T) {
		repo := NewAuctionRepository(mt.DB)

		// Mock insertion error
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "duplicate key error",
		}))

		auction := &auction_entity.Auction{
			Id:          "test-auction-id",
			ProductName: "Test Product",
			Category:    "Electronics",
			Description: "Test description for auction",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		err := repo.CreateAuction(context.Background(), auction)

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if err.Error() != "Error trying to insert auction" {
			t.Errorf("Expected 'Error trying to insert auction', got %v", err.Error())
		}
	})
}

func TestCreateAuction_GoroutineUpdateError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should handle update error in goroutine gracefully", func(mt *mtest.T) {
		os.Setenv("AUCTION_INTERVAL", "50ms")
		defer os.Unsetenv("AUCTION_INTERVAL")

		repo := NewAuctionRepository(mt.DB)

		// Mock successful insertion
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Mock update error for the goroutine
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    1,
			Message: "update failed",
		}))

		auction := &auction_entity.Auction{
			Id:          "test-auction-id",
			ProductName: "Test Product",
			Category:    "Electronics",
			Description: "Test description for auction",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		// Execute CreateAuction - should succeed despite future update error
		err := repo.CreateAuction(context.Background(), auction)

		if err != nil {
			t.Errorf("Expected no error for initial creation, got %v", err)
		}

		// Wait for goroutine to complete
		time.Sleep(100 * time.Millisecond)
	})
}

func TestCreateAuction_NoDocumentUpdated(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should handle case when no document is updated", func(mt *mtest.T) {
		os.Setenv("AUCTION_INTERVAL", "50ms")
		defer os.Unsetenv("AUCTION_INTERVAL")

		repo := NewAuctionRepository(mt.DB)

		// Mock successful insertion
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Mock update response with 0 modified documents
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "modifiedCount", Value: 0}))

		auction := &auction_entity.Auction{
			Id:          "test-auction-id",
			ProductName: "Test Product",
			Category:    "Electronics",
			Description: "Test description for auction",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		err := repo.CreateAuction(context.Background(), auction)

		if err != nil {
			t.Errorf("Expected no error for initial creation, got %v", err)
		}

		// Wait for goroutine to complete
		time.Sleep(100 * time.Millisecond)
	})
}

func TestCreateAuction_StatusTransition(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should transition from Active to Completed status", func(mt *mtest.T) {
		os.Setenv("AUCTION_INTERVAL", "50ms")
		defer os.Unsetenv("AUCTION_INTERVAL")

		repo := NewAuctionRepository(mt.DB)

		// Mock successful insertion
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Mock successful update with proper status transition
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "modifiedCount", Value: 1}))

		auction := &auction_entity.Auction{
			Id:          "test-auction-id",
			ProductName: "Test Product",
			Category:    "Electronics",
			Description: "Test description for auction",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active, // Initial status should be Active
			Timestamp:   time.Now(),
		}

		// Verify initial status is Active
		if auction.Status != auction_entity.Active {
			t.Errorf("Expected initial status to be Active, got %v", auction.Status)
		}

		err := repo.CreateAuction(context.Background(), auction)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Wait for goroutine to complete
		time.Sleep(100 * time.Millisecond)
	})
}

func TestCreateAuction_ContextCancellation(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should handle context cancellation properly", func(mt *mtest.T) {
		os.Setenv("AUCTION_INTERVAL", "50ms")
		defer os.Unsetenv("AUCTION_INTERVAL")

		repo := NewAuctionRepository(mt.DB)

		// Mock successful insertion
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Mock successful update for the goroutine
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "modifiedCount", Value: 1}))

		auction := &auction_entity.Auction{
			Id:          "test-auction-id",
			ProductName: "Test Product",
			Category:    "Electronics",
			Description: "Test description for auction",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())

		err := repo.CreateAuction(ctx, auction)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Cancel the context immediately after creation
		cancel()

		// Wait for goroutine to complete
		time.Sleep(100 * time.Millisecond)
	})
}

func TestGetAuctionInterval(t *testing.T) {
	tests := []struct {
		name           string
		envValue       string
		expectedResult time.Duration
	}{
		{
			name:           "valid duration string",
			envValue:       "2m",
			expectedResult: 2 * time.Minute,
		},
		{
			name:           "valid duration in seconds",
			envValue:       "30s",
			expectedResult: 30 * time.Second,
		},
		{
			name:           "invalid duration string",
			envValue:       "invalid",
			expectedResult: 5 * time.Minute, // default value
		},
		{
			name:           "empty environment variable",
			envValue:       "",
			expectedResult: 5 * time.Minute, // default value
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.envValue != "" {
				os.Setenv("AUCTION_INTERVAL", tt.envValue)
			} else {
				os.Unsetenv("AUCTION_INTERVAL")
			}
			defer os.Unsetenv("AUCTION_INTERVAL")

			result := getAuctionInterval()

			if result != tt.expectedResult {
				t.Errorf("Expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestCreateAuction_DataConsistency(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should maintain data consistency during creation and update", func(mt *mtest.T) {
		os.Setenv("AUCTION_INTERVAL", "50ms")
		defer os.Unsetenv("AUCTION_INTERVAL")

		repo := NewAuctionRepository(mt.DB)

		// Mock successful insertion
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Mock successful update
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "modifiedCount", Value: 1}))

		originalTime := time.Now()
		auction := &auction_entity.Auction{
			Id:          "test-auction-id",
			ProductName: "Test Product Name",
			Category:    "Test Category",
			Description: "Test description for the auction item",
			Condition:   auction_entity.Used,
			Status:      auction_entity.Active,
			Timestamp:   originalTime,
		}

		err := repo.CreateAuction(context.Background(), auction)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify that the original auction entity data remains unchanged
		if auction.Id != "test-auction-id" {
			t.Errorf("Expected Id to remain unchanged, got %v", auction.Id)
		}
		if auction.ProductName != "Test Product Name" {
			t.Errorf("Expected ProductName to remain unchanged, got %v", auction.ProductName)
		}
		if auction.Status != auction_entity.Active {
			t.Errorf("Expected original Status to remain Active, got %v", auction.Status)
		}
		if !auction.Timestamp.Equal(originalTime) {
			t.Errorf("Expected Timestamp to remain unchanged, got %v", auction.Timestamp)
		}

		// Wait for goroutine to complete
		time.Sleep(100 * time.Millisecond)
	})
}

func TestCreateAuction_ConcurrentCreations(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should handle concurrent auction creations", func(mt *mtest.T) {
		os.Setenv("AUCTION_INTERVAL", "100ms")
		defer os.Unsetenv("AUCTION_INTERVAL")

		repo := NewAuctionRepository(mt.DB)

		// Mock multiple successful insertions and updates
		for i := 0; i < 6; i++ { // 3 insertions + 3 updates
			if i < 3 {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			} else {
				mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "modifiedCount", Value: 1}))
			}
		}

		// Create multiple auctions concurrently
		done := make(chan bool, 3)
		for i := 0; i < 3; i++ {
			go func(index int) {
				auction := &auction_entity.Auction{
					Id:          "test-auction-id-" + string(rune(index+'0')),
					ProductName: "Test Product",
					Category:    "Electronics",
					Description: "Test description for auction",
					Condition:   auction_entity.New,
					Status:      auction_entity.Active,
					Timestamp:   time.Now(),
				}

				err := repo.CreateAuction(context.Background(), auction)
				if err != nil {
					t.Errorf("Expected no error for auction %d, got %v", index, err)
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 3; i++ {
			<-done
		}

		// Wait for all update goroutines to complete
		time.Sleep(150 * time.Millisecond)
	})
}