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

// TestUserRegistrationFlow tests the complete user registration workflow
// Validates Requirements: 1.1 (User Registration), 3.2 (Profile Management), 3.3 (Preferences)
func TestUserRegistrationFlow(t *testing.T) {
	// Setup
	authClient := helpers.NewTestClient("http://localhost:18080")
	userClient := helpers.NewTestClient("http://localhost:18081")

	// Wait for services to be ready
	require.NoError(t, helpers.WaitForService("http://localhost:18080", 30))
	require.NoError(t, helpers.WaitForService("http://localhost:18081", 30))

	// Generate unique email for this test run
	testEmail := fmt.Sprintf("newuser_%d@example.com", time.Now().Unix())

	t.Run("Step1_RegisterNewUser", func(t *testing.T) {
		// Prepare registration request
		registerReq := helpers.RegisterRequest{
			Email:       testEmail,
			Password:    "SecurePassword123!",
			FirstName:   "Integration",
			LastName:    "Test",
			DateOfBirth: "1995-06-15",
		}

		// Send registration request
		resp, err := authClient.Post("/api/auth/register", registerReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify successful registration
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Registration should return 201 Created")

		// Parse response
		var authResp helpers.AuthResponse
		err = helpers.ParseResponse(resp, &authResp)
		require.NoError(t, err)

		// Validate response structure
		assert.NotEmpty(t, authResp.AccessToken, "Access token should be present")
		assert.NotEmpty(t, authResp.RefreshToken, "Refresh token should be present")
		assert.Greater(t, authResp.ExpiresIn, int64(0), "ExpiresIn should be positive")
		assert.Equal(t, testEmail, authResp.User.Email, "Email should match")
		assert.Equal(t, "Integration", authResp.User.FirstName, "First name should match")
		assert.Equal(t, "Test", authResp.User.LastName, "Last name should match")
		assert.NotEmpty(t, authResp.User.ID, "User ID should be assigned")

		// Store token for subsequent requests
		authClient.SetAuthToken(authResp.AccessToken)
		userClient.SetAuthToken(authResp.AccessToken)

		// Store user ID for cleanup
		t.Cleanup(func() {
			// Cleanup is handled by database reset between test runs
		})
	})

	t.Run("Step2_GetUserProfile", func(t *testing.T) {
		// Get user profile
		resp, err := userClient.Get("/api/users/profile")
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify successful retrieval
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Profile retrieval should return 200 OK")

		// Parse response
		var profile helpers.UserProfile
		err = helpers.ParseResponse(resp, &profile)
		require.NoError(t, err)

		// Validate profile data
		assert.Equal(t, testEmail, profile.Email, "Email should match")
		assert.Equal(t, "Integration", profile.FirstName, "First name should match")
		assert.Equal(t, "Test", profile.LastName, "Last name should match")
		assert.Equal(t, "1995-06-15", profile.DateOfBirth, "Date of birth should match")
		assert.NotNil(t, profile.Preferences, "Preferences should be initialized")
	})

	t.Run("Step3_UpdateUserProfile", func(t *testing.T) {
		// Prepare update request
		updateReq := helpers.UpdateProfileRequest{
			FirstName:       "Updated",
			LastName:        "Name",
			ProfileImageURL: "https://example.com/avatar.jpg",
		}

		// Send update request
		resp, err := userClient.Put("/api/users/profile", updateReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify successful update
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Profile update should return 200 OK")

		// Parse response
		var profile helpers.UserProfile
		err = helpers.ParseResponse(resp, &profile)
		require.NoError(t, err)

		// Validate updated data
		assert.Equal(t, "Updated", profile.FirstName, "First name should be updated")
		assert.Equal(t, "Name", profile.LastName, "Last name should be updated")
		assert.Equal(t, "https://example.com/avatar.jpg", profile.ProfileImageURL, "Profile image URL should be updated")
	})

	t.Run("Step4_UpdateUserPreferences", func(t *testing.T) {
		// Prepare preferences update
		prefsReq := helpers.UserPreferences{
			Currency:             "USD",
			NotificationsEnabled: true,
			EmailNotifications:   true,
			PushNotifications:    false,
			SavingsReminders:     true,
			ReminderTime:         "09:00",
			Theme:                "dark",
		}

		// Send preferences update
		resp, err := userClient.Put("/api/users/preferences", prefsReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify successful update
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Preferences update should return 200 OK")

		// Get preferences to verify
		resp2, err := userClient.Get("/api/users/preferences")
		require.NoError(t, err)
		defer resp2.Body.Close()

		var prefs helpers.UserPreferences
		err = helpers.ParseResponse(resp2, &prefs)
		require.NoError(t, err)

		// Validate preferences
		assert.Equal(t, "USD", prefs.Currency, "Currency should be USD")
		assert.True(t, prefs.NotificationsEnabled, "Notifications should be enabled")
		assert.True(t, prefs.EmailNotifications, "Email notifications should be enabled")
		assert.False(t, prefs.PushNotifications, "Push notifications should be disabled")
		assert.True(t, prefs.SavingsReminders, "Savings reminders should be enabled")
		assert.Equal(t, "09:00", prefs.ReminderTime, "Reminder time should be 09:00")
		assert.Equal(t, "dark", prefs.Theme, "Theme should be dark")
	})

	t.Run("Step5_LoginWithNewCredentials", func(t *testing.T) {
		// Create new client without token
		newAuthClient := helpers.NewTestClient("http://localhost:18080")

		// Prepare login request
		loginReq := helpers.LoginRequest{
			Email:    testEmail,
			Password: "SecurePassword123!",
		}

		// Send login request
		resp, err := newAuthClient.Post("/api/auth/login", loginReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify successful login
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Login should return 200 OK")

		// Parse response
		var authResp helpers.AuthResponse
		err = helpers.ParseResponse(resp, &authResp)
		require.NoError(t, err)

		// Validate tokens are issued
		assert.NotEmpty(t, authResp.AccessToken, "Access token should be present")
		assert.NotEmpty(t, authResp.RefreshToken, "Refresh token should be present")
		assert.Equal(t, testEmail, authResp.User.Email, "Email should match")
	})
}

// TestRegistrationValidation tests registration input validation
// Validates Requirement: 1.1 (Registration validation)
func TestRegistrationValidation(t *testing.T) {
	authClient := helpers.NewTestClient("http://localhost:18080")
	require.NoError(t, helpers.WaitForService("http://localhost:18080", 30))

	t.Run("RejectShortPassword", func(t *testing.T) {
		registerReq := helpers.RegisterRequest{
			Email:       fmt.Sprintf("short_%d@example.com", time.Now().Unix()),
			Password:    "short",
			FirstName:   "Test",
			LastName:    "User",
			DateOfBirth: "1995-01-01",
		}

		resp, err := authClient.Post("/api/auth/register", registerReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should reject short password
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should reject password shorter than 8 characters")
	})

	t.Run("RejectDuplicateEmail", func(t *testing.T) {
		// Use existing test user email
		registerReq := helpers.RegisterRequest{
			Email:       "test1@example.com",
			Password:    "ValidPassword123!",
			FirstName:   "Test",
			LastName:    "User",
			DateOfBirth: "1995-01-01",
		}

		resp, err := authClient.Post("/api/auth/register", registerReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should reject duplicate email
		assert.True(t, resp.StatusCode == http.StatusConflict || resp.StatusCode == http.StatusBadRequest,
			"Should reject duplicate email with 409 Conflict or 400 Bad Request")
	})

	t.Run("RejectInvalidEmail", func(t *testing.T) {
		registerReq := helpers.RegisterRequest{
			Email:       "not-an-email",
			Password:    "ValidPassword123!",
			FirstName:   "Test",
			LastName:    "User",
			DateOfBirth: "1995-01-01",
		}

		resp, err := authClient.Post("/api/auth/register", registerReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should reject invalid email format
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should reject invalid email format")
	})
}
