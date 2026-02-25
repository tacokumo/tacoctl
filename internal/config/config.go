package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the main configuration structure
type Config struct {
	APIVersion     string      `yaml:"apiVersion"`
	Kind           string      `yaml:"kind"`
	CurrentContext string      `yaml:"current-context"`
	Preferences    Preferences `yaml:"preferences"`
	Contexts       []Context   `yaml:"contexts"`
	Portals        []Portal    `yaml:"portals"`
	Users          []User      `yaml:"users"`
}

// Context represents a portal+user combination
type Context struct {
	Name    string        `yaml:"name"`
	Context ContextDetail `yaml:"context"`
}

// ContextDetail holds the portal and user references
type ContextDetail struct {
	Portal string `yaml:"portal"`
	User   string `yaml:"user"`
}

// Portal represents Portal API server information
type Portal struct {
	Name   string       `yaml:"name"`
	Portal PortalDetail `yaml:"portal"`
}

// PortalDetail holds the server configuration
type PortalDetail struct {
	Server                string `yaml:"server"`
	InsecureSkipTLSVerify bool   `yaml:"insecure-skip-tls-verify"`
}

// User represents authentication information
type User struct {
	Name string     `yaml:"name"`
	User UserDetail `yaml:"user"`
}

// UserDetail holds the authentication details
type UserDetail struct {
	TokenFile  string `yaml:"token-file"`
	GitHubUser string `yaml:"github-user"`
	AuthType   string `yaml:"auth-type"` // "pat" or "jwt"
}

// Preferences represents user preferences
type Preferences struct {
	OutputFormat string `yaml:"output-format"` // "json", "yaml", "table"
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".tacokumo"), nil
}

// GetConfigPath returns the configuration file path
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config"), nil
}

// GetTokensDir returns the tokens directory path
func GetTokensDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "tokens"), nil
}

// EnsureConfigDir creates the configuration directory if it doesn't exist
func EnsureConfigDir() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// Create config directory with 0700 permissions
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create tokens directory with 0700 permissions
	tokensDir, err := GetTokensDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(tokensDir, 0700); err != nil {
		return fmt.Errorf("failed to create tokens directory: %w", err)
	}

	return nil
}

// NewConfig creates a new empty configuration
func NewConfig() *Config {
	return &Config{
		APIVersion: "v1",
		Kind:       "Config",
		Preferences: Preferences{
			OutputFormat: "table",
		},
		Contexts: []Context{},
		Portals:  []Portal{},
		Users:    []User{},
	}
}

// GetCurrentContext returns the current context details
func (c *Config) GetCurrentContext() (*Context, error) {
	if c.CurrentContext == "" {
		return nil, fmt.Errorf("no current context set")
	}

	for _, ctx := range c.Contexts {
		if ctx.Name == c.CurrentContext {
			return &ctx, nil
		}
	}

	return nil, fmt.Errorf("current context %q not found", c.CurrentContext)
}

// GetPortal returns the portal details by name
func (c *Config) GetPortal(name string) (*Portal, error) {
	for _, portal := range c.Portals {
		if portal.Name == name {
			return &portal, nil
		}
	}

	return nil, fmt.Errorf("portal %q not found", name)
}

// GetUser returns the user details by name
func (c *Config) GetUser(name string) (*User, error) {
	for _, user := range c.Users {
		if user.Name == name {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user %q not found", name)
}

// AddContext adds a new context to the configuration
func (c *Config) AddContext(ctx Context) {
	// Remove existing context with the same name
	for i, existing := range c.Contexts {
		if existing.Name == ctx.Name {
			c.Contexts = append(c.Contexts[:i], c.Contexts[i+1:]...)
			break
		}
	}
	c.Contexts = append(c.Contexts, ctx)
}

// AddPortal adds a new portal to the configuration
func (c *Config) AddPortal(portal Portal) {
	// Remove existing portal with the same name
	for i, existing := range c.Portals {
		if existing.Name == portal.Name {
			c.Portals = append(c.Portals[:i], c.Portals[i+1:]...)
			break
		}
	}
	c.Portals = append(c.Portals, portal)
}

// AddUser adds a new user to the configuration
func (c *Config) AddUser(user User) {
	// Remove existing user with the same name
	for i, existing := range c.Users {
		if existing.Name == user.Name {
			c.Users = append(c.Users[:i], c.Users[i+1:]...)
			break
		}
	}
	c.Users = append(c.Users, user)
}
