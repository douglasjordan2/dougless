package permissions

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSavePermissionToConfig(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".douglessrc")

	// Test 1: Create new config with a permission
	err := SavePermissionToConfig(configPath, PermissionRead, "/home/user/test")
	if err != nil {
		t.Fatalf("Failed to save permission: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Load and verify content
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	if len(config.Permissions.Read) != 1 {
		t.Fatalf("Expected 1 read permission, got %d", len(config.Permissions.Read))
	}

	if config.Permissions.Read[0] != "/home/user/test" {
		t.Errorf("Expected /home/user/test, got %s", config.Permissions.Read[0])
	}

	// Test 2: Add another permission to existing config
	err = SavePermissionToConfig(configPath, PermissionWrite, "/home/user/data")
	if err != nil {
		t.Fatalf("Failed to add second permission: %v", err)
	}

	// Reload and verify
	data, err = os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	if len(config.Permissions.Read) != 1 {
		t.Errorf("Read permissions changed unexpectedly")
	}

	if len(config.Permissions.Write) != 1 {
		t.Fatalf("Expected 1 write permission, got %d", len(config.Permissions.Write))
	}

	if config.Permissions.Write[0] != "/home/user/data" {
		t.Errorf("Expected /home/user/data, got %s", config.Permissions.Write[0])
	}

	// Test 3: Try to add duplicate permission (should be no-op)
	err = SavePermissionToConfig(configPath, PermissionRead, "/home/user/test")
	if err != nil {
		t.Fatalf("Failed to save duplicate permission: %v", err)
	}

	// Verify still only 1 read permission
	data, err = os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	if len(config.Permissions.Read) != 1 {
		t.Errorf("Expected 1 read permission after duplicate, got %d", len(config.Permissions.Read))
	}
}

func TestSavePermissionToConfigEmptyPath(t *testing.T) {
	// Save current working directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	// Create and change to temp directory
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	// Save with empty config path should use current directory
	err := SavePermissionToConfig("", PermissionNet, "localhost")
	if err != nil {
		t.Fatalf("Failed to save with empty path: %v", err)
	}

	// Verify file exists in current directory
	configPath := filepath.Join(tmpDir, ".douglessrc")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created in current directory")
	}

	// Verify content
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	if len(config.Permissions.Net) != 1 || config.Permissions.Net[0] != "localhost" {
		t.Errorf("Expected localhost net permission, got %v", config.Permissions.Net)
	}
}

func TestSavePermissionToConfigAllTypes(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".douglessrc")

	// Add one permission of each type
	testCases := []struct {
		perm     Permission
		resource string
	}{
		{PermissionRead, "/tmp/read"},
		{PermissionWrite, "/tmp/write"},
		{PermissionNet, "api.example.com"},
		{PermissionEnv, "API_KEY"},
		{PermissionRun, "curl"},
	}

	for _, tc := range testCases {
		err := SavePermissionToConfig(configPath, tc.perm, tc.resource)
		if err != nil {
			t.Fatalf("Failed to save %s permission: %v", tc.perm, err)
		}
	}

	// Load and verify all permissions
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	if len(config.Permissions.Read) != 1 || config.Permissions.Read[0] != "/tmp/read" {
		t.Errorf("Read permission mismatch: %v", config.Permissions.Read)
	}

	if len(config.Permissions.Write) != 1 || config.Permissions.Write[0] != "/tmp/write" {
		t.Errorf("Write permission mismatch: %v", config.Permissions.Write)
	}

	if len(config.Permissions.Net) != 1 || config.Permissions.Net[0] != "api.example.com" {
		t.Errorf("Net permission mismatch: %v", config.Permissions.Net)
	}

	if len(config.Permissions.Env) != 1 || config.Permissions.Env[0] != "API_KEY" {
		t.Errorf("Env permission mismatch: %v", config.Permissions.Env)
	}

	if len(config.Permissions.Run) != 1 || config.Permissions.Run[0] != "curl" {
		t.Errorf("Run permission mismatch: %v", config.Permissions.Run)
	}
}
