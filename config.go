package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

// Config represents the application configuration
type Config struct {
	TaskManager TaskManagerConfig `toml:"taskmanager"`
	Display     DisplayConfig     `toml:"display"`
}

// TaskManagerConfig holds the task manager specific settings
type TaskManagerConfig struct {
	Directory   string   `toml:"directory"`   // Single directory (deprecated, use Directories)
	Directories []string `toml:"directories"` // Multiple directories containing task markdown files
}

// DisplayConfig holds display customization settings
type DisplayConfig struct {
	StatusIndicators map[string]string `toml:"status_indicators"` // Custom status indicators
	DefaultStatus    string            `toml:"default_status"`    // Default status for tasks without one
}

// GetStatusIndicator returns the indicator for a given status
// Falls back to defaults if not configured
func (c *DisplayConfig) GetStatusIndicator(status string) string {
	// If custom indicator is defined, use it
	if indicator, ok := c.StatusIndicators[status]; ok {
		return indicator
	}

	// Fall back to defaults
	return getDefaultStatusIndicator(status)
}

// GetDefaultStatus returns the configured default status, or "todo" if not set
func (c *DisplayConfig) GetDefaultStatus() string {
	if c.DefaultStatus != "" {
		return c.DefaultStatus
	}
	return "todo"
}

// getDefaultStatusIndicator returns the default indicator for a status
func getDefaultStatusIndicator(status string) string {
	switch status {
	case "done", "completed":
		return "[✓]"
	case "in-progress", "doing":
		return "[~]"
	case "todo":
		return "[ ]"
	default:
		return "   "
	}
}

// GetDirectories returns all configured directories
// Handles both old single directory and new multiple directories config
func (c *TaskManagerConfig) GetDirectories() []string {
	// If Directories is set, use that
	if len(c.Directories) > 0 {
		return c.Directories
	}

	// Otherwise fall back to single Directory
	if c.Directory != "" {
		return []string{c.Directory}
	}

	// Default fallback
	return []string{"~/.tasks"}
}

// defaultConfig returns the default configuration
func defaultConfig() Config {
	return Config{
		TaskManager: TaskManagerConfig{
			Directories: []string{"~/.tasks"},
		},
		Display: DisplayConfig{
			StatusIndicators: map[string]string{
				"todo":        "[ ]",
				"in-progress": "[~]",
				"done":        "[✓]",
			},
			DefaultStatus: "todo",
		},
	}
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	var configDir string

	// Use ~/.config on Unix-like systems (macOS and Linux)
	// Use standard location on Windows
	if runtime.GOOS == "windows" {
		// On Windows, use the standard AppData location
		dir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("couldn't get config directory: %w", err)
		}
		configDir = dir
	} else {
		// On macOS and Linux, use ~/.config
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("couldn't get home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, ".config")
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
} // saveConfig writes the configuration to file
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
