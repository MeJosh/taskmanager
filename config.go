package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents the application configuration
type Config struct {
	TaskManager TaskManagerConfig `toml:"taskmanager"`
}

// TaskManagerConfig holds the task manager specific settings
type TaskManagerConfig struct {
	Directory string `toml:"directory"` // Directory containing task markdown files
}

// defaultConfig returns the default configuration
func defaultConfig() Config {
	return Config{
		TaskManager: TaskManagerConfig{
			Directory: "~/.tasks",
		},
	}
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	// Get user's config directory (~/.config)
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("couldn't get config directory: %w", err)
	}

	// Path to our app's config directory
	appConfigDir := filepath.Join(configDir, "taskmanager")
	configFile := filepath.Join(appConfigDir, "config.toml")

	return configFile, nil
}

// loadConfig loads the configuration from file, or creates a default one if it doesn't exist
func loadConfig() (Config, error) {
	configFile, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Config doesn't exist, create it with defaults
		cfg := defaultConfig()
		if err := saveConfig(cfg); err != nil {
			return Config{}, fmt.Errorf("failed to create default config: %w", err)
		}
		return cfg, nil
	}

	// Read the config file
	var cfg Config
	if _, err := toml.DecodeFile(configFile, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// saveConfig writes the configuration to file
func saveConfig(cfg Config) error {
	configFile, err := getConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create the config file
	f, err := os.Create(configFile)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer f.Close()

	// Write the config as TOML
	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
