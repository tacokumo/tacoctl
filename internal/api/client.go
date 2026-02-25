package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/tacokumo/tacoctl/internal/auth"
)

// Client represents a Portal API client
type Client struct {
	authClient *auth.AuthenticatedClient
}

// NewClient creates a new Portal API client
func NewClient(baseURL, authToken, authType string, opts auth.HTTPClientOptions) *Client {
	return &Client{
		authClient: auth.NewAuthenticatedClient(baseURL, authToken, authType, opts),
	}
}

// HealthCheck performs a health check against the Portal API
func (c *Client) HealthCheck(ctx context.Context) (*HealthResponse, error) {
	resp, err := c.authClient.Do(ctx, "GET", "/health", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make health check request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("health check failed with status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("health check failed: %s", errResp.Message)
	}

	var healthResp HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		return nil, fmt.Errorf("failed to decode health response: %w", err)
	}

	return &healthResp, nil
}

// GetUserProfile retrieves the current user's profile
func (c *Client) GetUserProfile(ctx context.Context) (*UserProfile, error) {
	resp, err := c.authClient.Do(ctx, "GET", "/api/v1/user/profile", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make user profile request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("get user profile failed with status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("get user profile failed: %s", errResp.Message)
	}

	var userProfile UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&userProfile); err != nil {
		return nil, fmt.Errorf("failed to decode user profile response: %w", err)
	}

	return &userProfile, nil
}

// ListApplications retrieves a list of applications
func (c *Client) ListApplications(ctx context.Context) ([]Application, error) {
	resp, err := c.authClient.Do(ctx, "GET", "/api/v1/applications", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make list applications request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("list applications failed with status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("list applications failed: %s", errResp.Message)
	}

	var applications []Application
	if err := json.NewDecoder(resp.Body).Decode(&applications); err != nil {
		return nil, fmt.Errorf("failed to decode applications response: %w", err)
	}

	return applications, nil
}
