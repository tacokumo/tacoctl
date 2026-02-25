package auth

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

// GitHubUser represents a GitHub user
type GitHubUser struct {
	Login string `json:"login"`
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// ValidateGitHubPAT validates a GitHub Personal Access Token
func ValidateGitHubPAT(ctx context.Context, token string) (*GitHubUser, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}

// HTTPClientOptions represents options for creating an HTTP client
type HTTPClientOptions struct {
	InsecureSkipTLSVerify bool
	Timeout               time.Duration
}

// NewHTTPClient creates a new HTTP client with the specified options
func NewHTTPClient(opts HTTPClientOptions) *http.Client {
	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}

	transport := &http.Transport{}
	if opts.InsecureSkipTLSVerify {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return &http.Client{
		Timeout:   opts.Timeout,
		Transport: transport,
	}
}

// AuthenticatedClient represents an HTTP client with authentication
type AuthenticatedClient struct {
	client    *http.Client
	baseURL   string
	authToken string
	authType  string // "pat" or "jwt"
}

// NewAuthenticatedClient creates a new authenticated HTTP client
func NewAuthenticatedClient(baseURL, authToken, authType string, opts HTTPClientOptions) *AuthenticatedClient {
	return &AuthenticatedClient{
		client:    NewHTTPClient(opts),
		baseURL:   baseURL,
		authToken: authToken,
		authType:  authType,
	}
}

// Do performs an HTTP request with authentication
func (c *AuthenticatedClient) Do(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	u, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	var req *http.Request
	if body != nil {
		// For simplicity, we'll implement JSON body support when needed
		return nil, fmt.Errorf("request body not yet implemented")
	}

	req, err = http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication header
	switch c.authType {
	case "pat":
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	case "jwt":
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	default:
		return nil, fmt.Errorf("unsupported auth type: %s", c.authType)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return c.client.Do(req)
}
