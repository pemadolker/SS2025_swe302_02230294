// repository/user_repository_test.go
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	redisModule "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ============================================================================
// EXERCISE 1: Basic TestContainers Setup
// ============================================================================

// Global test database connection (Exercise 1)
var testDB *sql.DB

// Global Redis client for Exercise 5
var testRedis *redis.Client

// TestMain sets up the test environment (Exercise 1 & Exercise 5)
// This runs ONCE before all tests in this package
func TestMain(m *testing.M) {
	ctx := context.Background()

	// ========== EXERCISE 1: Setup PostgreSQL Container ==========
	// Create a PostgreSQL container
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.WithInitScripts("../migrations/init.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start PostgreSQL container: %v\n", err)
		os.Exit(1)
	}

	// Ensure PostgreSQL container is terminated at the end
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to terminate PostgreSQL container: %v\n", err)
		}
	}()

	// Get PostgreSQL connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get connection string: %v\n", err)
		os.Exit(1)
	}

	// Connect to the PostgreSQL database
	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Verify PostgreSQL connection
	if err = testDB.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to ping database: %v\n", err)
		os.Exit(1)
	}

	// ========== EXERCISE 5: Setup Redis Container ==========
	// Create a Redis container
	redisContainer, err := redisModule.RunContainer(ctx,
		testcontainers.WithImage("redis:7-alpine"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start Redis container: %v\n", err)
		os.Exit(1)
	}

	// Ensure Redis container is terminated at the end
	defer func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to terminate Redis container: %v\n", err)
		}
	}()

	// Get Redis connection string
	redisEndpoint, err := redisContainer.Endpoint(ctx, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get Redis endpoint: %v\n", err)
		os.Exit(1)
	}

	// Connect to Redis
	testRedis = redis.NewClient(&redis.Options{
		Addr: redisEndpoint,
	})

	// Verify Redis connection
	if err = testRedis.Ping(ctx).Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to ping Redis: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	testDB.Close()
	testRedis.Close()
	os.Exit(code)
}

// ============================================================================
// EXERCISE 1: Basic Integration Tests
// ============================================================================

// TestGetByID tests retrieving a user by ID (Exercise 1)
func TestGetByID(t *testing.T) {
	repo := NewUserRepository(testDB)

	// Test case 1: User exists (from init.sql)
	t.Run("User Exists", func(t *testing.T) {
		user, err := repo.GetByID(1)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if user.Email != "alice@example.com" {
			t.Errorf("Expected email 'alice@example.com', got: %s", user.Email)
		}

		if user.Name != "Alice Smith" {
			t.Errorf("Expected name 'Alice Smith', got: %s", user.Name)
		}

		if user.ID != 1 {
			t.Errorf("Expected ID 1, got: %d", user.ID)
		}
	})

	// Test case 2: User does not exist
	t.Run("User Not Found", func(t *testing.T) {
		_, err := repo.GetByID(9999)
		if err == nil {
			t.Fatal("Expected error for non-existent user, got nil")
		}
	})
}

// TestGetByEmail tests retrieving a user by email (Exercise 1)
func TestGetByEmail(t *testing.T) {
	repo := NewUserRepository(testDB)

	t.Run("User Exists", func(t *testing.T) {
		user, err := repo.GetByEmail("bob@example.com")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if user.Name != "Bob Johnson" {
			t.Errorf("Expected name 'Bob Johnson', got: %s", user.Name)
		}

		if user.Email != "bob@example.com" {
			t.Errorf("Expected email 'bob@example.com', got: %s", user.Email)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		_, err := repo.GetByEmail("nonexistent@example.com")
		if err == nil {
			t.Fatal("Expected error for non-existent email, got nil")
		}
	})
}

// ============================================================================
// EXERCISE 2: Complete CRUD Testing
// ============================================================================

// TestCreate tests user creation (Exercise 2)
func TestCreate(t *testing.T) {
	repo := NewUserRepository(testDB)

	t.Run("Create New User - Successful Creation", func(t *testing.T) {
		user, err := repo.Create("charlie@example.com", "Charlie Brown")
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Verify auto-generated ID
		if user.ID == 0 {
			t.Error("Expected non-zero ID for created user")
		}

		// Verify email
		if user.Email != "charlie@example.com" {
			t.Errorf("Expected email 'charlie@example.com', got: %s", user.Email)
		}

		// Verify name
		if user.Name != "Charlie Brown" {
			t.Errorf("Expected name 'Charlie Brown', got: %s", user.Name)
		}

		// Check created_at timestamp
		if user.CreatedAt.IsZero() {
			t.Error("Expected non-zero created_at timestamp")
		}

		// Cleanup: delete the created user
		defer repo.Delete(user.ID)
	})

	t.Run("Create Duplicate Email - Test Constraint", func(t *testing.T) {
		// Try to create user with existing email (from init.sql)
		_, err := repo.Create("alice@example.com", "Another Alice")
		if err == nil {
			t.Fatal("Expected error when creating user with duplicate email")
		}
	})
}

// TestUpdate tests user updates (Exercise 2)
func TestUpdate(t *testing.T) {
	repo := NewUserRepository(testDB)

	t.Run("Update Existing User - Successful Update", func(t *testing.T) {
		// First, create a user to update
		user, err := repo.Create("david@example.com", "David Davis")
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
		defer repo.Delete(user.ID)

		// Update the user
		err = repo.Update(user.ID, "david.updated@example.com", "David Updated")
		if err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		// Verify changes persist
		updatedUser, err := repo.GetByID(user.ID)
		if err != nil {
			t.Fatalf("Failed to retrieve updated user: %v", err)
		}

		if updatedUser.Email != "david.updated@example.com" {
			t.Errorf("Expected email 'david.updated@example.com', got: %s", updatedUser.Email)
		}

		if updatedUser.Name != "David Updated" {
			t.Errorf("Expected name 'David Updated', got: %s", updatedUser.Name)
		}
	})

	t.Run("Update Non-Existent User", func(t *testing.T) {
		err := repo.Update(9999, "nobody@example.com", "Nobody")
		if err == nil {
			t.Fatal("Expected error when updating non-existent user")
		}
	})
}

// TestDelete tests user deletion (Exercise 2)
func TestDelete(t *testing.T) {
	repo := NewUserRepository(testDB)

	t.Run("Delete Existing User - Successful Deletion", func(t *testing.T) {
		// Create a user to delete
		user, err := repo.Create("temp@example.com", "Temporary User")
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// Delete the user
		err = repo.Delete(user.ID)
		if err != nil {
			t.Fatalf("Failed to delete user: %v", err)
		}

		// Verify user is gone after deletion
		_, err = repo.GetByID(user.ID)
		if err == nil {
			t.Fatal("Expected error when retrieving deleted user")
		}
	})

	t.Run("Delete Non-Existent User", func(t *testing.T) {
		err := repo.Delete(9999)
		if err == nil {
			t.Fatal("Expected error when deleting non-existent user")
		}
	})
}

// TestList tests listing all users (Exercise 2)
func TestList(t *testing.T) {
	repo := NewUserRepository(testDB)

	users, err := repo.List()
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}

	// Should have at least 2 users from init.sql
	if len(users) < 2 {
		t.Errorf("Expected at least 2 users, got: %d", len(users))
	}

	// Verify order (ordered by id)
	for i := 1; i < len(users); i++ {
		if users[i].ID < users[i-1].ID {
			t.Error("Users are not ordered by ID")
			break
		}
	}

	// Check count
	count, err := repo.CountUsers()
	if err != nil {
		t.Fatalf("Failed to count users: %v", err)
	}

	if len(users) != count {
		t.Errorf("List count (%d) doesn't match CountUsers (%d)", len(users), count)
	}
}

// ============================================================================
// EXERCISE 3: Advanced Queries Testing
// ============================================================================

// TestFindByNamePattern tests pattern matching queries (Exercise 3)
func TestFindByNamePattern(t *testing.T) {
	repo := NewUserRepository(testDB)

	// Create test data with various patterns
	testUsers := []struct {
		email string
		name  string
	}{
		{"john.smith@example.com", "John Smith"},
		{"jane.smith@example.com", "Jane Smith"},
		{"mike.jones@example.com", "Mike Jones"},
	}

	var createdIDs []int
	for _, tu := range testUsers {
		user, err := repo.Create(tu.email, tu.name)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
		createdIDs = append(createdIDs, user.ID)
	}

	// Cleanup after test
	defer func() {
		for _, id := range createdIDs {
			repo.Delete(id)
		}
	}()

	t.Run("Find Users with Smith Pattern", func(t *testing.T) {
		users, err := repo.FindByNamePattern("%Smith%")
		if err != nil {
			t.Fatalf("Failed to find users: %v", err)
		}

		// Should find at least 2 users (John Smith and Jane Smith)
		smithCount := 0
		for _, user := range users {
			if user.Name == "John Smith" || user.Name == "Jane Smith" {
				smithCount++
			}
		}

		if smithCount < 2 {
			t.Errorf("Expected at least 2 users with 'Smith', got: %d", smithCount)
		}
	})

	t.Run("Find Users with No Match", func(t *testing.T) {
		users, err := repo.FindByNamePattern("%NonExistent%")
		if err != nil {
			t.Fatalf("Failed to find users: %v", err)
		}

		if len(users) != 0 {
			t.Errorf("Expected 0 users with 'NonExistent', got: %d", len(users))
		}
	})
}

// TestCountUsers tests counting users (Exercise 3)
func TestCountUsers(t *testing.T) {
	repo := NewUserRepository(testDB)

	// Count users before
	countBefore, err := repo.CountUsers()
	if err != nil {
		t.Fatalf("Failed to count users: %v", err)
	}

	// Create a new user
	user, err := repo.Create("count.test@example.com", "Count Test")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	defer repo.Delete(user.ID)

	// Count users after
	countAfter, err := repo.CountUsers()
	if err != nil {
		t.Fatalf("Failed to count users: %v", err)
	}

	// Verify count increased by 1
	if countAfter != countBefore+1 {
		t.Errorf("Expected count to increase by 1, before: %d, after: %d", countBefore, countAfter)
	}
}

// TestGetRecentUsers tests date filtering (Exercise 3)
func TestGetRecentUsers(t *testing.T) {
	repo := NewUserRepository(testDB)

	// Create a recent user
	user, err := repo.Create("recent@example.com", "Recent User")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	defer repo.Delete(user.ID)

	t.Run("Get Users Created in Last 7 Days", func(t *testing.T) {
		users, err := repo.GetRecentUsers(7)
		if err != nil {
			t.Fatalf("Failed to get recent users: %v", err)
		}

		// Should find at least the user we just created
		found := false
		for _, u := range users {
			if u.ID == user.ID {
				found = true
				break
			}
		}

		if !found {
			t.Error("Expected to find recently created user")
		}

		// Verify date filtering is accurate
		sevenDaysAgo := time.Now().AddDate(0, 0, -7)
		for _, u := range users {
			if u.CreatedAt.Before(sevenDaysAgo) {
				t.Errorf("User %d created at %v is older than 7 days", u.ID, u.CreatedAt)
			}
		}
	})

	t.Run("Get Users Created in Last 0 Days", func(t *testing.T) {
		users, err := repo.GetRecentUsers(0)
		if err != nil {
			t.Fatalf("Failed to get recent users: %v", err)
		}

		// Should find users created today
		today := time.Now().Truncate(24 * time.Hour)
		for _, u := range users {
			if u.CreatedAt.Before(today) {
				t.Errorf("User %d created at %v is not from today", u.ID, u.CreatedAt)
			}
		}
	})
}

// ============================================================================
// EXERCISE 4: Transaction Testing
// ============================================================================

// TestTransactionRollback tests that failed transactions roll back properly (Exercise 4)
func TestTransactionRollback(t *testing.T) {
	repo := NewUserRepository(testDB)

	// Count users before
	countBefore, err := repo.CountUsers()
	if err != nil {
		t.Fatal(err)
	}

	// Start a transaction that will fail
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}

	// Create user in transaction
	_, err = tx.Exec("INSERT INTO users (email, name) VALUES ($1, $2)",
		"tx@example.com", "TX User")
	if err != nil {
		t.Fatal(err)
	}

	// Rollback transaction
	tx.Rollback()

	// Verify count is unchanged (data consistency after rollback)
	countAfter, err := repo.CountUsers()
	if err != nil {
		t.Fatal(err)
	}

	if countAfter != countBefore {
		t.Errorf("Transaction was not rolled back properly. Before: %d, After: %d", countBefore, countAfter)
	}
}

// TestTransactionCommit tests that successful transactions commit properly (Exercise 4)
func TestTransactionCommit(t *testing.T) {
	repo := NewUserRepository(testDB)

	// Count users before
	countBefore, err := repo.CountUsers()
	if err != nil {
		t.Fatal(err)
	}

	// Start a transaction that will succeed
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}

	// Create user in transaction
	var userID int
	err = tx.QueryRow("INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id",
		"txcommit@example.com", "TX Commit User").Scan(&userID)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	// Cleanup
	defer repo.Delete(userID)

	// Verify count increased (successful transaction commits)
	countAfter, err := repo.CountUsers()
	if err != nil {
		t.Fatal(err)
	}

	if countAfter != countBefore+1 {
		t.Errorf("Transaction was not committed properly. Before: %d, After: %d", countBefore, countAfter)
	}
}

// TestTransferUserData tests the transaction-aware transfer method (Exercise 4)
func TestTransferUserData(t *testing.T) {
	repo := NewUserRepository(testDB)

	t.Run("Successful Transfer", func(t *testing.T) {
		// Create two users for transfer
		fromUser, err := repo.Create("from@example.com", "From User")
		if err != nil {
			t.Fatalf("Failed to create from user: %v", err)
		}

		toUser, err := repo.Create("to@example.com", "To User")
		if err != nil {
			t.Fatalf("Failed to create to user: %v", err)
		}
		defer repo.Delete(toUser.ID)

		// Transfer data
		err = repo.TransferUserData(fromUser.ID, toUser.ID)
		if err != nil {
			t.Fatalf("Failed to transfer user data: %v", err)
		}

		// Verify from user is deleted
		_, err = repo.GetByID(fromUser.ID)
		if err == nil {
			t.Error("Expected from user to be deleted")
		}

		// Verify to user has updated name
		updatedToUser, err := repo.GetByID(toUser.ID)
		if err != nil {
			t.Fatalf("Failed to get updated to user: %v", err)
		}

		if updatedToUser.Name != fromUser.Name {
			t.Errorf("Expected name '%s', got '%s'", fromUser.Name, updatedToUser.Name)
		}
	})

	t.Run("Failed Transfer - Non-existent Target", func(t *testing.T) {
		// Create source user
		fromUser, err := repo.Create("fail@example.com", "Fail User")
		if err != nil {
			t.Fatalf("Failed to create from user: %v", err)
		}
		defer repo.Delete(fromUser.ID)

		// Try to transfer to non-existent user
		err = repo.TransferUserData(fromUser.ID, 9999)
		if err == nil {
			t.Error("Expected error when transferring to non-existent user")
		}

		// Verify source user still exists (rollback occurred)
		_, err = repo.GetByID(fromUser.ID)
		if err != nil {
			t.Error("Expected source user to still exist after failed transfer")
		}
	})
}

// ============================================================================
// EXERCISE 5: Multi-Container Testing (PostgreSQL + Redis)
// ============================================================================

// TestCachedGetByID tests the cached repository (Exercise 5)
func TestCachedGetByID(t *testing.T) {
	cachedRepo := NewCachedUserRepository(testDB, testRedis)
	ctx := context.Background()

	// Create a test user
	repo := NewUserRepository(testDB)
	user, err := repo.Create("cached@example.com", "Cached User")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	defer repo.Delete(user.ID)

	t.Run("Cache Miss - Queries Database", func(t *testing.T) {
		// Clear cache first
		cachedRepo.InvalidateCache(user.ID)

		// First call should be a cache miss
		fetchedUser, err := cachedRepo.GetByIDCached(user.ID)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}

		if fetchedUser.Email != user.Email {
			t.Errorf("Expected email '%s', got '%s'", user.Email, fetchedUser.Email)
		}

		// Verify data is now in cache
		cacheKey := fmt.Sprintf("user:%d", user.ID)
		exists, err := testRedis.Exists(ctx, cacheKey).Result()
		if err != nil {
			t.Fatalf("Failed to check cache: %v", err)
		}

		if exists != 1 {
			t.Error("Expected user to be cached after database query")
		}
	})

	t.Run("Cache Hit - Returns Cached Data", func(t *testing.T) {
		// This call should hit the cache (from previous test)
		fetchedUser, err := cachedRepo.GetByIDCached(user.ID)
		if err != nil {
			t.Fatalf("Failed to get cached user: %v", err)
		}

		if fetchedUser.Email != user.Email {
			t.Errorf("Expected email '%s', got '%s'", user.Email, fetchedUser.Email)
		}

		// Verify cache is still there
		cacheKey := fmt.Sprintf("user:%d", user.ID)
		exists, err := testRedis.Exists(ctx, cacheKey).Result()
		if err != nil {
			t.Fatalf("Failed to check cache: %v", err)
		}

		if exists != 1 {
			t.Error("Expected user to still be in cache")
		}
	})

	t.Run("Cache Invalidation", func(t *testing.T) {
		// Invalidate cache
		err := cachedRepo.InvalidateCache(user.ID)
		if err != nil {
			t.Fatalf("Failed to invalidate cache: %v", err)
		}

		// Verify cache is empty
		cacheKey := fmt.Sprintf("user:%d", user.ID)
		exists, err := testRedis.Exists(ctx, cacheKey).Result()
		if err != nil {
			t.Fatalf("Failed to check cache: %v", err)
		}

		if exists != 0 {
			t.Error("Expected cache to be empty after invalidation")
		}
	})
}

// TestMultiContainerCommunication verifies both containers work together (Exercise 5)
func TestMultiContainerCommunication(t *testing.T) {
	cachedRepo := NewCachedUserRepository(testDB, testRedis)
	repo := NewUserRepository(testDB)
	ctx := context.Background()

	// Create test user
	user, err := repo.Create("multi@example.com", "Multi Container User")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	defer repo.Delete(user.ID)

	// Clear cache
	cachedRepo.InvalidateCache(user.ID)

	// First fetch - should query database and populate cache
	_, err = cachedRepo.GetByIDCached(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Update user in database
	err = repo.Update(user.ID, "multi.updated@example.com", "Multi Updated")
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Fetch from cache - should return OLD data (cache not invalidated yet)
	cachedUser, err := cachedRepo.GetByIDCached(user.ID)
	if err != nil {
		t.Fatalf("Failed to get cached user: %v", err)
	}

	if cachedUser.Email == "multi.updated@example.com" {
		t.Error("Expected cached data, but got fresh data from database")
	}

	// Invalidate cache
	cachedRepo.InvalidateCache(user.ID)

	// Fetch again - should query database and get NEW data
	freshUser, err := cachedRepo.GetByIDCached(user.ID)
	if err != nil {
		t.Fatalf("Failed to get fresh user: %v", err)
	}

	if freshUser.Email != "multi.updated@example.com" {
		t.Errorf("Expected email 'multi.updated@example.com', got '%s'", freshUser.Email)
	}

	// Verify Redis is working
	pingResult := testRedis.Ping(ctx).Val()
	if pingResult != "PONG" {
		t.Errorf("Redis ping failed: %s", pingResult)
	}

	// Verify PostgreSQL is working
	if err := testDB.Ping(); err != nil {
		t.Errorf("PostgreSQL ping failed: %v", err)
	}

	t.Log("✅ Both PostgreSQL and Redis containers are working correctly together!")
}