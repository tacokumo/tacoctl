package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/tacokumo/tacoctl/internal/config"
)

// NewViewCmd creates the config view command
func NewViewCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "view",
		Short: "Display tacoctl configuration",
		Long: `Display the current tacoctl configuration.

This command shows your current configuration including contexts, portals,
and users. Sensitive information like tokens are masked for security.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runView(outputFormat)
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (json|yaml|table)")

	return cmd
}

func runView(outputFormat string) error {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check if configuration is empty
	if cfg.CurrentContext == "" && len(cfg.Contexts) == 0 {
		fmt.Println("No configuration found. Run 'tacoctl config init' to set up your configuration.")
		return nil
	}

	// Create a copy for display (mask sensitive information)
	displayConfig := maskSensitiveInfo(cfg)

	// Determine output format
	if outputFormat == "" {
		outputFormat = cfg.Preferences.OutputFormat
		if outputFormat == "" {
			outputFormat = "table"
		}
	}

	// Display configuration based on format
	switch outputFormat {
	case "json":
		return displayJSON(displayConfig)
	case "yaml":
		return displayYAML(displayConfig)
	case "table":
		return displayTable(displayConfig)
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
}

func maskSensitiveInfo(cfg *config.Config) *config.Config {
	// Create a deep copy
	masked := *cfg
	masked.Users = make([]config.User, len(cfg.Users))
	copy(masked.Users, cfg.Users)

	// Mask token files
	for i, user := range masked.Users {
		if user.User.TokenFile != "" {
			masked.Users[i].User.TokenFile = maskTokenPath(user.User.TokenFile)
		}
	}

	return &masked
}

func maskTokenPath(tokenPath string) string {
	// Show only the filename, not the full path
	parts := strings.Split(tokenPath, "/")
	if len(parts) > 0 {
		return "~/.tacokumo/tokens/" + parts[len(parts)-1]
	}
	return tokenPath
}

func displayJSON(cfg *config.Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func displayYAML(cfg *config.Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %w", err)
	}

	fmt.Print(string(data))
	return nil
}

func displayTable(cfg *config.Config) error {
	fmt.Printf("API Version: %s\n", cfg.APIVersion)
	fmt.Printf("Kind: %s\n", cfg.Kind)
	fmt.Printf("Current Context: %s\n", cfg.CurrentContext)
	fmt.Printf("Output Format: %s\n", cfg.Preferences.OutputFormat)
	fmt.Println()

	// Display contexts
	fmt.Println("CONTEXTS:")
	if len(cfg.Contexts) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, ctx := range cfg.Contexts {
			marker := " "
			if ctx.Name == cfg.CurrentContext {
				marker = "*"
			}
			fmt.Printf("%s %-15s portal: %-15s user: %s\n",
				marker, ctx.Name, ctx.Context.Portal, ctx.Context.User)
		}
	}
	fmt.Println()

	// Display portals
	fmt.Println("PORTALS:")
	if len(cfg.Portals) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, portal := range cfg.Portals {
			fmt.Printf("  %-15s server: %s\n", portal.Name, portal.Portal.Server)
			if portal.Portal.InsecureSkipTLSVerify {
				fmt.Printf("  %-15s insecure-skip-tls-verify: true\n", "")
			}
		}
	}
	fmt.Println()

	// Display users
	fmt.Println("USERS:")
	if len(cfg.Users) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, user := range cfg.Users {
			fmt.Printf("  %-15s github-user: %-15s auth-type: %-5s token-file: %s\n",
				user.Name, user.User.GitHubUser, user.User.AuthType, user.User.TokenFile)
		}
	}

	return nil
}
