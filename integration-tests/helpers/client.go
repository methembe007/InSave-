package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TestClient provides HTTP client functionality for integration tests
type TestClient struct {
	BaseURL    string
	HTTPClient *http.Client
	AuthToken  string
}

// NewTestClient creates a new test client
func NewTestClient(baseURL string) *TestClient {
	return &TestClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SetAuthToken sets the authentication token for subsequent requests
func (c *TestClient) SetAuthToken(token string) {
	c.AuthToken = token
}

// Get performs a GET request
func (c *TestClient) Get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}

	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.HTTPClient.Do(req)
}

// Post performs a POST request with JSON body
func (c *TestClient) Post(path string, body interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.HTTPClient.Do(req)
}

// Put performs a PUT request with JSON body
func (c *TestClient) Put(path string, body interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", c.BaseURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.HTTPClient.Do(req)
}

// Delete performs a DELETE request
func (c *TestClient) Delete(path string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}

	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.HTTPClient.Do(req)
}

// ParseResponse parses JSON response into target struct
func ParseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return json.Unmarshal(body, target)
}

// WaitForService waits for a service to become healthy
func WaitForService(baseURL string, maxRetries int) error {
	client := &http.Client{Timeout: 2 * time.Second}

	for i := 0; i < maxRetries; i++ {
		resp, err := client.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("service at %s did not become healthy after %d retries", baseURL, maxRetries)
}
