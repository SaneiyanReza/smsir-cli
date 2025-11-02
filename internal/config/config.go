package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	// configDirName is the name of the configuration directory
	configDirName = ".smsir"
	// configFileName is the name of the configuration file
	configFileName = "config.json"
	// defaultConfigPerms are the default permissions for config directory
	defaultConfigPerms = 0755
	// defaultFilePerms are the default permissions for config file
	defaultFilePerms = 0644
	// defaultBaseURL is the default SMS.ir API base URL
	defaultBaseURL = "https://api.sms.ir/v1"
)

// Config holds the application configuration
type Config struct {
	APIKey     string `json:"api_key" mapstructure:"api_key"`
	LineNumber string `json:"line_number" mapstructure:"line_number"`
	BaseURL    string `json:"base_url" mapstructure:"base_url"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		BaseURL: defaultBaseURL,
	}
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()

	// Set config file path
	configFile, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Config file doesn't exist, create default one
		if err := createDefaultConfig(configFile); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
	}

	// Set viper configuration
	viper.SetConfigFile(configFile)
	viper.SetConfigType("json")
	viper.SetEnvPrefix("SMSIR")
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("base_url", defaultBaseURL)

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal config
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

// getConfigPath returns the path to the configuration file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, configDirName)
	configFile := filepath.Join(configDir, configFileName)

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, defaultConfigPerms); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configFile, nil
}

// SaveConfig saves the configuration to file
func (c *Config) SaveConfig() error {
	configFile, err := getConfigPath()
	if err != nil {
		return err
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configFile, data, defaultFilePerms); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("api key is required")
	}
	if c.LineNumber == "" {
		return fmt.Errorf("line number is required")
	}
	return nil
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig(configFile string) error {
	cfg := DefaultConfig()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, defaultFilePerms)
}
