package integration_tests

import (
	"net/http"
	"testing"
	"time"

	"github.com/insavein/integration-tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSavingsFlow tests the complete savings workflow
// Validates Requirements: 4.1 (Savings Transactions), 4.3 (Transaction Recording), 5.1 (Streak Calculation)
func TestSavingsFlow(t *testing.T) {
	// Setup
	authClient := helpers.NewTestClient("http://localhost:18080")
	savingsClient := helpers.NewTestClient("http://localhost:18082")
	notificationClient := helpers.NewTestClient("http://localhost:18086")

	// Wait for services
	require.NoError(t, helpers.WaitForService("http://localhost:18080", 30))
	require.NoError(t, helpers.WaitForService("http://localhost:18082", 30))
	require.NoError(t, helpers.WaitForService("http://localhost:18086", 30))

	// Login as test user
	loginReq := helpers.LoginRequest{
		Email:    "test1@example.com",
		Password: "TestPassword123!",
	}

	resp, err := authClient.Post("/api/auth/login", loginReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	var authResp helpers.AuthResponse
	err = helpers.ParseResponse(resp, &authResp)
	require.NoError(t, err)

	// Set auth token for all clients
	token := authResp.AccessToken
	savingsClient.SetAuthToken(token)
	notificationClient.SetAuthToken(token)

	t.Run("Step1_CreateSavingsTransaction", func(t *testing.T) {
		// Create savings transaction
		savingsReq := helpers.CreateSavingsRequest{
			Amount:      75.50,
			Currency:    "USD",
			Description: "Weekly savings deposit",
			Category:    "general",
		}

		resp, err := savingsClient.Post("/api/savings/transactions", savingsReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify successful creation
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Transaction creation should return 201 Created")

		// Parse response
		var transaction helpers.SavingsTransaction
		err = helpers.ParseResponse(resp, &transaction)
		require.NoError(t, err)

		// Validate transaction data
		assert.NotEmpty(t, transaction.ID, "Transaction ID should be assigned")
		assert.Equal(t, 75.50, transaction.Amount, "Amount should match")
		assert.Equal(t, "USD", transaction.Currency, "Currency should match")
		assert.Equal(t, "Weekly savings deposit", transaction.Description, "Description should match")
		assert.Equal(t, "general", transaction.Category, "Category should match")
	})

	t.Run("Step2_GetSavingsSummary", func(t *testing.T) {
		// Get savings summary
		resp, err := savingsClient.Get("/api/savings/summary")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Summary retrieval should return 200 OK")

		var summary helpers.SavingsSummary
		err = helpers.ParseResponse(resp, &summary)
		require.NoError(t, err)

		// Validate summary data
		assert.Greater(t, summary.TotalSaved, 0.0, "Total saved should be positive")
		assert.GreaterOrEqual(t, summary.CurrentStreak, 0, "Current streak should be non-negative")
		assert.GreaterOrEqual(t, summary.LongestStreak, summary.CurrentStreak, "Longest streak should be >= current streak")
	})

	t.Run("Step3_GetSavingsHistory", func(t *testing.T) {
		// Get savings history
		resp, err := savingsClient.Get("/api/savings/history")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "History retrieval should return 200 OK")

		var transactions []helpers.SavingsTransaction
		err = helpers.ParseResponse(resp, &transactions)
		require.NoError(t, err)

		// Validate history
		assert.NotEmpty(t, transactions, "Should have at least one transaction")

		// Verify transactions are ordered by date descending
		if len(transactions) > 1 {
			for i := 0; i < len(transactions)-1; i++ {
				assert.True(t, transactions[i].CreatedAt.After(transactions[i+1].CreatedAt) ||
					transactions[i].CreatedAt.Equal(transactions[i+1].CreatedAt),
					"Transactions should be ordered by date descending")
			}
		}
	})

	t.Run("Step4_VerifyStreakCalculation", func(t *testing.T) {
		// Create another transaction to build streak
		savingsReq := helpers.CreateSavingsRequest{
			Amount:      50.00,
			Currency:    "USD",
			Description: "Building my streak",
			Category:    "general",
		}

		resp, err := savingsClient.Post("/api/savings/transactions", savingsReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Wait a moment for streak calculation
		time.Sleep(1 * time.Second)

		// Get updated streak
		resp2, err := savingsClient.Get("/api/savings/streak")
		require.NoError(t, err)
		defer resp2.Body.Close()

		var streak helpers.SavingsStreak
		err = helpers.ParseResponse(resp2, &streak)
		require.NoError(t, err)

		// Validate streak
		assert.GreaterOrEqual(t, streak.CurrentStreak, 1, "Should have at least 1 day streak")
		assert.GreaterOrEqual(t, streak.LongestStreak, streak.CurrentStreak, "Longest streak should be >= current")
	})
}

// TestSavingsValidation tests savings transaction validation
// Validates Requirement: 4.1 (Transaction validation)
func TestSavingsValidation(t *testing.T) {
	authClient := helpers.NewTestClient("http://localhost:18080")
	savingsClient := helpers.NewTestClient("http://localhost:18082")

	require.NoError(t, helpers.WaitForService("http://localhost:18080", 30))
	require.NoError(t, helpers.WaitForService("http://localhost:18082", 30))

	// Login
	loginReq := helpers.LoginRequest{
		Email:    "test1@example.com",
		Password: "TestPassword123!",
	}

	resp, err := authClient.Post("/api/auth/login", loginReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	var authResp helpers.AuthResponse
	err = helpers.ParseResponse(resp, &authResp)
	require.NoError(t, err)
	savingsClient.SetAuthToken(authResp.AccessToken)

	t.Run("RejectNegativeAmount", func(t *testing.T) {
		savingsReq := helpers.CreateSavingsRequest{
			Amount:      -10.00,
			Currency:    "USD",
			Description: "Invalid negative amount",
			Category:    "general",
		}

		resp, err := savingsClient.Post("/api/savings/transactions", savingsReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should reject negative amount")
	})

	t.Run("RejectZeroAmount", func(t *testing.T) {
		savingsReq := helpers.CreateSavingsRequest{
			Amount:      0.00,
			Currency:    "USD",
			Description: "Invalid zero amount",
			Category:    "general",
		}

		resp, err := savingsClient.Post("/api/savings/transactions", savingsReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should reject zero amount")
	})
}
