package config

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/tacokumo/tacoctl/internal/auth"
	"github.com/tacokumo/tacoctl/internal/config"
)

// NewInitCmd creates the config init command
func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize tacoctl configuration",
		Long: `Initialize tacoctl configuration interactively.

This command will guide you through setting up your Portal API access,
including authentication credentials and preferences.`,
		RunE: runInit,
	}

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Load existing configuration if it exists
	existingConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load existing config: %w", err)
	}

	// Check if configuration already exists
	if existingConfig.CurrentContext != "" {
		fmt.Printf("Configuration already exists (current context: %s)\n", existingConfig.CurrentContext)
		fmt.Println("This will overwrite your existing configuration.")
		fmt.Println()
	}

	// Run interactive setup
	newConfig, err := config.InteractiveSetup()
	if err != nil {
		return fmt.Errorf("failed to setup configuration: %w", err)
	}

	// Validate the configuration by testing the connection
	fmt.Println()
	fmt.Println("Validating configuration...")

	if err := validateConfig(ctx, newConfig); err != nil {
		fmt.Printf("Warning: Configuration validation failed: %v\n", err)
		fmt.Println("Configuration has been saved but may not work properly.")
	} else {
		fmt.Println("Configuration validation successful!")
	}

	// Save the configuration
	if err := config.SaveConfig(newConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("\nConfiguration saved successfully!\n")
	fmt.Printf("Current context: %s\n", newConfig.CurrentContext)
	fmt.Println("\nYou can now use tacoctl to interact with your Portal API.")
	fmt.Println("Try: tacoctl config view")

	return nil
}

func validateConfig(ctx context.Context, cfg *config.Config) error {
	// Get current context
	currentCtx, err := cfg.GetCurrentContext()
	if err != nil {
		return fmt.Errorf("failed to get current context: %w", err)
	}

	// Get portal info
	_, err = cfg.GetPortal(currentCtx.Context.Portal)
	if err != nil {
		return fmt.Errorf("failed to get portal info: %w", err)
	}

	// Get user info
	user, err := cfg.GetUser(currentCtx.Context.User)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	// Load token
	token, err := config.LoadToken(user.User.TokenFile)
	if err != nil {
		return fmt.Errorf("failed to load token: %w", err)
	}

	// Validate GitHub PAT if using PAT authentication
	if user.User.AuthType == "pat" {
		githubUser, err := auth.ValidateGitHubPAT(ctx, token)
		if err != nil {
			return fmt.Errorf("failed to validate GitHub PAT: %w", err)
		}

		if githubUser.Login != user.User.GitHubUser {
			return fmt.Errorf("GitHub username mismatch: expected %s, got %s",
				user.User.GitHubUser, githubUser.Login)
		}
	}

	// TODO: Add Portal API health check validation once we have a working Portal API
	// For now, we'll just validate the GitHub token

	return nil
}
