package permissions

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Permissions PermissionSet `json:"permissions"`
}

type PermissionSet struct {
	Read  []string `json:"read"`
	Write []string `json:"write"`
	Net   []string `json:"net"`
	Env   []string `json:"env"`
	Run   []string `json:"run"`
}

func FindConfig(startDir string) (string, error) {
	configPath := filepath.Join(startDir, ".douglessrc")

	if _, err := os.Stat(configPath); err == nil {
		return configPath, nil
	}

	return "", fmt.Errorf("no .douglessrc found in %s", startDir)
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// SavePermissionToConfig adds a permission to .douglessrc, creating it if needed.
// If configPath is empty, creates .douglessrc in the current directory.
func SavePermissionToConfig(configPath string, perm Permission, resource string) error {
	// If no config path specified, use .douglessrc in current directory
	if configPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		configPath = filepath.Join(cwd, ".douglessrc")
	}

	// Load existing config or create new one
	var config Config
	if data, err := os.ReadFile(configPath); err == nil {
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse existing config: %w", err)
		}
	}

	// Add the permission to the appropriate array (avoid duplicates)
	var targetArray *[]string
	switch perm {
	case PermissionRead:
		targetArray = &config.Permissions.Read
	case PermissionWrite:
		targetArray = &config.Permissions.Write
	case PermissionNet:
		targetArray = &config.Permissions.Net
	case PermissionEnv:
		targetArray = &config.Permissions.Env
	case PermissionRun:
		targetArray = &config.Permissions.Run
	default:
		return fmt.Errorf("unknown permission type: %s", perm)
	}

	// Check if already exists
	for _, existing := range *targetArray {
		if existing == resource {
			return nil // Already exists, no need to save
		}
	}

	// Add the new permission
	*targetArray = append(*targetArray, resource)

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
