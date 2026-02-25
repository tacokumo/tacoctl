package config

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/tacokumo/tacoctl/internal/prompt"
)

// InteractiveSetup performs interactive setup for a new configuration
func InteractiveSetup() (*Config, error) {
	config := NewConfig()

	fmt.Println("Welcome to tacoctl configuration setup!")
	fmt.Println("This will help you configure your Portal API access.")
	fmt.Println()

	// Get context name
	contextName, err := prompt.Input("Enter context name", "default")
	if err != nil {
		return nil, fmt.Errorf("failed to get context name: %w", err)
	}

	// Get Portal URL
	portalURL, err := getPortalURL()
	if err != nil {
		return nil, err
	}

	// Get authentication method
	authType, err := getAuthType()
	if err != nil {
		return nil, err
	}

	// Get authentication credentials
	var githubUser, token string
	if authType == "pat" {
		githubUser, token, err = getGitHubPATCredentials()
		if err != nil {
			return nil, err
		}
	} else {
		// JWT authentication would be implemented here in the future
		return nil, fmt.Errorf("JWT authentication not yet implemented")
	}

	// Get output format preference
	outputFormat, err := getOutputFormat()
	if err != nil {
		return nil, err
	}

	// Get TLS verification preference
	insecureSkipTLS, err := getTLSPreference()
	if err != nil {
		return nil, err
	}

	// Create token file path
	tokensDir, err := GetTokensDir()
	if err != nil {
		return nil, err
	}
	tokenFile := filepath.Join(tokensDir, contextName)

	// Save token to file
	if err := SaveToken(tokenFile, token); err != nil {
		return nil, fmt.Errorf("failed to save token: %w", err)
	}

	// Build configuration
	portalName := contextName + "-portal"
	userName := contextName + "-user"

	config.CurrentContext = contextName
	config.Preferences.OutputFormat = outputFormat

	config.AddContext(Context{
		Name: contextName,
		Context: ContextDetail{
			Portal: portalName,
			User:   userName,
		},
	})

	config.AddPortal(Portal{
		Name: portalName,
		Portal: PortalDetail{
			Server:                portalURL,
			InsecureSkipTLSVerify: insecureSkipTLS,
		},
	})

	config.AddUser(User{
		Name: userName,
		User: UserDetail{
			TokenFile:  tokenFile,
			GitHubUser: githubUser,
			AuthType:   authType,
		},
	})

	return config, nil
}

func getPortalURL() (string, error) {
	for {
		portalURL, err := prompt.Input("Enter Portal API server URL", "https://portal.tacokumo.com")
		if err != nil {
			return "", err
		}

		// Validate URL format
		if _, err := url.Parse(portalURL); err != nil {
			fmt.Printf("Invalid URL format: %v\n", err)
			continue
		}

		// Ensure https:// or http:// prefix
		if !strings.HasPrefix(portalURL, "http://") && !strings.HasPrefix(portalURL, "https://") {
			portalURL = "https://" + portalURL
		}

		return portalURL, nil
	}
}

func getAuthType() (string, error) {
	options := []string{
		"GitHub Personal Access Token (recommended)",
		"JWT Token (via OAuth)",
	}

	index, err := prompt.Select("Select authentication method:", options, 0)
	if err != nil {
		return "", err
	}

	if index == 0 {
		return "pat", nil
	}
	return "jwt", nil
}

func getGitHubPATCredentials() (string, string, error) {
	fmt.Println()
	fmt.Println("GitHub Personal Access Token (PAT) Authentication")
	fmt.Println("Please create a PAT at: https://github.com/settings/tokens")
	fmt.Println("Required scopes: read:user, user:email")
	fmt.Println()

	githubUser, err := prompt.Input("Enter your GitHub username", "")
	if err != nil {
		return "", "", err
	}

	if githubUser == "" {
		return "", "", fmt.Errorf("GitHub username is required")
	}

	token, err := prompt.Password("Enter your GitHub Personal Access Token")
	if err != nil {
		return "", "", err
	}

	if token == "" {
		return "", "", fmt.Errorf("GitHub Personal Access Token is required")
	}

	return githubUser, token, nil
}

func getOutputFormat() (string, error) {
	options := []string{"table", "json", "yaml"}
	index, err := prompt.Select("Select default output format:", options, 0)
	if err != nil {
		return "", err
	}
	return options[index], nil
}

func getTLSPreference() (bool, error) {
	return prompt.Confirm("Skip TLS certificate verification? (only for development)", false)
}
