package integration_tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/insavein/integration-tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGoalProgressFlow tests the complete goal progress workflow
// Validates Requirements: 9.1 (Goal Creation), 10.1 (Progress Tracking), 10.4 (Milestone Completion)
func TestGoalProgressFlow(t *testing.T) {
	// Setup
	authClient := helpers.NewTestClient("http://localhost:18080")
	goalClient := helpers.NewTestClient("http://localhost:18005")

	// Wait for services
	require.NoError(t, helpers.WaitForService("http://localhost:18080", 30))
	require.NoError(t, helpers.WaitForService("http://localhost:18005", 30))

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

	goalClient.SetAuthToken(authResp.AccessToken)

	var goalID string
	var milestones []helpers.Milestone

	t.Run("Step1_CreateGoal", func(t *testing.T) {
		// Create a new financial goal
		targetDate := time.Now().AddDate(0, 6, 0).Format("2006-01-02")

		goalReq := helpers.CreateGoalRequest{
			Title:        "Vacation Fund",
			Description:  "Save for summer vacation",
			TargetAmount: 2000.00,
			Currency:     "USD",
			TargetDate:   targetDate,
			Milestones: []helpers.Milestone{
				{Title: "First Quarter", Amount: 500.00, Order: 1},
				{Title: "Halfway There", Amount: 1000.00, Order: 2},
				{Title: "Three Quarters", Amount: 1500.00, Order: 3},
				{Title: "Goal Complete", Amount: 2000.00, Order: 4},
			},
		}

		resp, err := goalClient.Post("/api/goals", goalReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Goal creation should return 201 Created")

		var goal helpers.Goal
		err = helpers.ParseResponse(resp, &goal)
		require.NoError(t, err)

		goalID = goal.ID
		assert.NotEmpty(t, goalID, "Goal ID should be assigned")
		assert.Equal(t, "Vacation Fund", goal.Title, "Title should match")
		assert.Equal(t, 2000.00, goal.TargetAmount, "Target amount should match")
		assert.Equal(t, 0.00, goal.CurrentAmount, "Initial current amount should be 0")
		assert.Equal(t, "active", goal.Status, "Initial status should be active")
		assert.Equal(t, 0.00, goal.ProgressPercent, "Initial progress should be 0%")
	})

	t.Run("Step2_GetGoalDetails", func(t *testing.T) {
		// Get goal details including milestones
		resp, err := goalClient.Get(fmt.Sprintf("/api/goals/%s", goalID))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Goal retrieval should return 200 OK")

		var goal helpers.Goal
		err = helpers.ParseResponse(resp, &goal)
		require.NoError(t, err)

		assert.Equal(t, goalID, goal.ID, "Goal ID should match")
		assert.NotEmpty(t, goal.Milestones, "Should have milestones")

		milestones = goal.Milestones
		assert.Len(t, milestones, 4, "Should have 4 milestones")

		// Verify milestones are ordered
		for i := 0; i < len(milestones)-1; i++ {
			assert.Less(t, milestones[i].Amount, milestones[i+1].Amount,
				"Milestones should be ordered by amount")
		}
	})

	t.Run("Step3_AddFirstContribution", func(t *testing.T) {
		// Add contribution to reach first milestone
		contributionReq := helpers.ContributionRequest{
			Amount: 600.00, // Exceeds first milestone of 500
		}

		resp, err := goalClient.Post(fmt.Sprintf("/api/goals/%s/contributions", goalID), contributionReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Contribution should be accepted")

		var goal helpers.Goal
		err = helpers.ParseResponse(resp, &goal)
		require.NoError(t, err)

		// Verify progress updated
		assert.Equal(t, 600.00, goal.CurrentAmount, "Current amount should be updated")
		assert.Equal(t, 30.00, goal.ProgressPercent, "Progress should be 30%")
		assert.Equal(t, "active", goal.Status, "Status should still be active")

		// Wait for milestone processing
		time.Sleep(1 * time.Second)

		// Get updated goal to check milestones
		resp2, err := goalClient.Get(fmt.Sprintf("/api/goals/%s", goalID))
		require.NoError(t, err)
		defer resp2.Body.Close()

		var updatedGoal helpers.Goal
		err = helpers.ParseResponse(resp2, &updatedGoal)
		require.NoError(t, err)

		// First milestone should be completed
		if len(updatedGoal.Milestones) > 0 {
			assert.True(t, updatedGoal.Milestones[0].IsCompleted,
				"First milestone should be completed")
			assert.NotNil(t, updatedGoal.Milestones[0].CompletedAt,
				"First milestone should have completion timestamp")
		}
	})

	t.Run("Step4_AddMultipleContributions", func(t *testing.T) {
		// Add more contributions to reach halfway point
		contributionReq := helpers.ContributionRequest{
			Amount: 500.00, // Total now 1100
		}

		resp, err := goalClient.Post(fmt.Sprintf("/api/goals/%s/contributions", goalID), contributionReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var goal helpers.Goal
		err = helpers.ParseResponse(resp, &goal)
		require.NoError(t, err)

		assert.Equal(t, 1100.00, goal.CurrentAmount, "Current amount should be 1100")
		assert.Equal(t, 55.00, goal.ProgressPercent, "Progress should be 55%")

		// Wait for milestone processing
		time.Sleep(1 * time.Second)

		// Verify second milestone completed
		resp2, err := goalClient.Get(fmt.Sprintf("/api/goals/%s", goalID))
		require.NoError(t, err)
		defer resp2.Body.Close()

		var updatedGoal helpers.Goal
		err = helpers.ParseResponse(resp2, &updatedGoal)
		require.NoError(t, err)

		if len(updatedGoal.Milestones) >= 2 {
			assert.True(t, updatedGoal.Milestones[1].IsCompleted,
				"Second milestone should be completed")
		}
	})

	t.Run("Step5_CompleteGoal", func(t *testing.T) {
		// Add final contribution to complete goal
		contributionReq := helpers.ContributionRequest{
			Amount: 900.00, // Total now 2000 - goal complete
		}

		resp, err := goalClient.Post(fmt.Sprintf("/api/goals/%s/contributions", goalID), contributionReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var goal helpers.Goal
		err = helpers.ParseResponse(resp, &goal)
		require.NoError(t, err)

		// Verify goal completion
		assert.Equal(t, 2000.00, goal.CurrentAmount, "Current amount should equal target")
		assert.Equal(t, 100.00, goal.ProgressPercent, "Progress should be 100%")
		assert.Equal(t, "completed", goal.Status, "Status should be completed")

		// Wait for milestone processing
		time.Sleep(1 * time.Second)

		// Verify all milestones completed
		resp2, err := goalClient.Get(fmt.Sprintf("/api/goals/%s", goalID))
		require.NoError(t, err)
		defer resp2.Body.Close()

		var completedGoal helpers.Goal
		err = helpers.ParseResponse(resp2, &completedGoal)
		require.NoError(t, err)

		// All milestones should be completed
		for i, milestone := range completedGoal.Milestones {
			assert.True(t, milestone.IsCompleted,
				fmt.Sprintf("Milestone %d should be completed", i+1))
			assert.NotNil(t, milestone.CompletedAt,
				fmt.Sprintf("Milestone %d should have completion timestamp", i+1))
		}
	})

	t.Run("Step6_GetActiveGoals", func(t *testing.T) {
		// Get active goals (completed goal should not appear)
		resp, err := goalClient.Get("/api/goals")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var goals []helpers.Goal
		err = helpers.ParseResponse(resp, &goals)
		require.NoError(t, err)

		// Our completed goal should not be in active goals
		for _, goal := range goals {
			if goal.ID == goalID {
				assert.NotEqual(t, "completed", goal.Status,
					"Completed goal should not appear in active goals list")
			}
		}
	})
}
