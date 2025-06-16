package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config represents the configuration for the Slurm REST CLI.
type Config struct {
	Auth struct {
		Keycloak struct {
			ServerURL string `yaml:"server_url" mapstructure:"server_url"`
			Realm     string `yaml:"realm" mapstructure:"realm"`
			ClientID  string `yaml:"client_id" mapstructure:"client_id"`
			Username  string `yaml:"username" mapstructure:"username"`
		} `yaml:"keycloak" mapstructure:"keycloak"`
		JWT struct {
			Token  string `yaml:"token" mapstructure:"token"`
			Expiry string `yaml:"expiry" mapstructure:"expiry"`
		} `yaml:"jwt" mapstructure:"jwt"`
		Token    string `yaml:"token" mapstructure:"token"`
		Username string `yaml:"username" mapstructure:"username"`
	} `yaml:"auth" mapstructure:"auth"`
	API struct {
		URL     string `yaml:"url" mapstructure:"url"`
		Version string `yaml:"version" mapstructure:"version"`
	} `yaml:"api" mapstructure:"api"`
	Preferences struct {
		OutputFormat string `yaml:"output_format" mapstructure:"output_format"`
	} `yaml:"preferences" mapstructure:"preferences"`
}

// DefaultConfigPath returns the default path for the config file.
func DefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}
	return filepath.Join(home, ".config", "slurm-rest-cli.yaml"), nil
}

// LoadConfig loads the configuration from the specified file.
func LoadConfig(configPath string) (*Config, error) {
	var config Config

	// If no config path is provided, use the default
	if configPath == "" {
		var err error
		configPath, err = DefaultConfigPath()
		if err != nil {
			return nil, err
		}
	}

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &config, nil // Return empty config if file doesn't exist
	}

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse the config file
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to the specified file.
func SaveConfig(config *Config, configPath string) error {
	// If no config path is provided, use the default
	if configPath == "" {
		var err error
		configPath, err = DefaultConfigPath()
		if err != nil {
			return err
		}
	}

	// Create the directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	// Marshal the config to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	// Write the config file
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// GetAuthInfo returns the authentication information from the config.
func (c *Config) GetAuthInfo() (token, jwt, username string) {
	token = c.Auth.Token
	jwt = c.Auth.JWT.Token
	username = c.Auth.Username
	
	// If username is not set in the root, try to get it from Keycloak
	if username == "" {
		username = c.Auth.Keycloak.Username
	}
	
	return token, jwt, username
}

// GetAPIURL returns the API URL from the config.
func (c *Config) GetAPIURL() string {
	return c.API.URL
}

// GetOutputFormat returns the output format from the config.
func (c *Config) GetOutputFormat() string {
	if c.Preferences.OutputFormat == "" {
		return "table" // Default to table format
	}
	return c.Preferences.OutputFormat
}

// LoadFromViper loads the configuration from Viper.
func LoadFromViper() (*Config, error) {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config from viper: %w", err)
	}
	return &config, nil
}
