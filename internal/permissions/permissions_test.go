package permissions

import (
	"path/filepath"
	"testing"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()
	
	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}
	
	// All permissions should be nil (denied) by default
	if manager.allowRead != nil {
		t.Error("allowRead should be nil by default")
	}
	if manager.allowWrite != nil {
		t.Error("allowWrite should be nil by default")
	}
	if manager.allowNet != nil {
		t.Error("allowNet should be nil by default")
	}
	// Note: promptMode is context-aware (true in interactive, false in non-interactive)
	// so we don't assert a specific default value here
}

func TestGrantAll(t *testing.T) {
	manager := NewManager()
	manager.GrantAll()
	
	// All permissions should be granted to everything
	if manager.allowRead == nil || len(*manager.allowRead) != 0 {
		t.Error("allowRead should be empty slice (allow all)")
	}
	if manager.allowWrite == nil || len(*manager.allowWrite) != 0 {
		t.Error("allowWrite should be empty slice (allow all)")
	}
	if manager.allowNet == nil || len(*manager.allowNet) != 0 {
		t.Error("allowNet should be empty slice (allow all)")
	}
}

func TestGrantSpecificPermissions(t *testing.T) {
	manager := NewManager()
	
	t.Run("grant specific read paths", func(t *testing.T) {
		paths := []string{"/tmp", "/home/user"}
		manager.GrantRead(paths)
		
		if manager.allowRead == nil {
			t.Fatal("allowRead should not be nil after grant")
		}
		if len(*manager.allowRead) != 2 {
			t.Errorf("expected 2 paths, got %d", len(*manager.allowRead))
		}
	})
	
	t.Run("grant all read access", func(t *testing.T) {
		manager.GrantRead([]string{})
		
		if manager.allowRead == nil {
			t.Fatal("allowRead should not be nil")
		}
		if len(*manager.allowRead) != 0 {
			t.Error("empty slice should mean allow all")
		}
	})
}

func TestCheckPermissionDenied(t *testing.T) {
	manager := NewManager()
	
	// No permissions granted - everything should be denied
	tests := []struct {
		name     string
		perm     Permission
		resource string
	}{
		{"read file", PermissionRead, "/etc/passwd"},
		{"write file", PermissionWrite, "/tmp/test.txt"},
		{"network", PermissionNet, "example.com"},
		{"env var", PermissionEnv, "PATH"},
		{"run program", PermissionRun, "ls"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if manager.Check(tt.perm, tt.resource) {
				t.Errorf("permission should be denied for %s", tt.resource)
			}
		})
	}
}

func TestCheckPermissionGrantedAll(t *testing.T) {
	manager := NewManager()
	manager.GrantAll()
	
	// Everything should be allowed
	tests := []struct {
		name     string
		perm     Permission
		resource string
	}{
		{"read file", PermissionRead, "/etc/passwd"},
		{"write file", PermissionWrite, "/tmp/test.txt"},
		{"network", PermissionNet, "example.com"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !manager.Check(tt.perm, tt.resource) {
				t.Errorf("permission should be granted for %s", tt.resource)
			}
		})
	}
}

func TestCheckPermissionGrantedSpecific(t *testing.T) {
	manager := NewManager()
	
	t.Run("read specific paths", func(t *testing.T) {
		manager.GrantRead([]string{"/tmp"})
		
		// Should allow /tmp and subdirectories
		if !manager.Check(PermissionRead, "/tmp/test.txt") {
			t.Error("/tmp/test.txt should be allowed")
		}
		if !manager.Check(PermissionRead, "/tmp/subdir/file.txt") {
			t.Error("/tmp/subdir/file.txt should be allowed")
		}
		
		// Should deny other paths
		if manager.Check(PermissionRead, "/etc/passwd") {
			t.Error("/etc/passwd should be denied")
		}
		if manager.Check(PermissionRead, "/home/user/file.txt") {
			t.Error("/home/user/file.txt should be denied")
		}
	})
	
	t.Run("multiple read paths", func(t *testing.T) {
		manager.GrantRead([]string{"/tmp", "/home/user"})
		
		if !manager.Check(PermissionRead, "/tmp/test.txt") {
			t.Error("/tmp/test.txt should be allowed")
		}
		if !manager.Check(PermissionRead, "/home/user/doc.txt") {
			t.Error("/home/user/doc.txt should be allowed")
		}
		if manager.Check(PermissionRead, "/etc/passwd") {
			t.Error("/etc/passwd should be denied")
		}
	})
	
	t.Run("write specific paths", func(t *testing.T) {
		manager.GrantWrite([]string{"/tmp"})
		
		if !manager.Check(PermissionWrite, "/tmp/output.txt") {
			t.Error("/tmp/output.txt should be allowed")
		}
		if manager.Check(PermissionWrite, "/etc/config") {
			t.Error("/etc/config should be denied")
		}
	})
}

func TestMatchPath(t *testing.T) {
	tests := []struct {
		name          string
		allowedPath   string
		requestedPath string
		shouldMatch   bool
	}{
		{"exact match", "/tmp", "/tmp", true},
		{"subdirectory", "/tmp", "/tmp/test.txt", true},
		{"nested subdirectory", "/tmp", "/tmp/sub/dir/file.txt", true},
		{"different path", "/tmp", "/etc/passwd", false},
		{"parent directory escape attempt", "/tmp", "/tmp/../etc/passwd", false},
		{"relative path allowed", "/tmp", "tmp/test.txt", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchPath(tt.allowedPath, tt.requestedPath)
			if result != tt.shouldMatch {
				t.Errorf("matchPath(%q, %q) = %v, want %v", 
					tt.allowedPath, tt.requestedPath, result, tt.shouldMatch)
			}
		})
	}
}

func TestMatchPathTraversalPrevention(t *testing.T) {
	// Test that path traversal attacks are blocked
	manager := NewManager()
	manager.GrantRead([]string{"/tmp"})
	
	// These should all be blocked
	traversalAttempts := []string{
		"/tmp/../etc/passwd",
		"/tmp/../../root/.ssh/id_rsa",
		"/tmp/./../etc/shadow",
	}
	
	for _, attempt := range traversalAttempts {
		if manager.Check(PermissionRead, attempt) {
			t.Errorf("traversal attempt should be blocked: %s", attempt)
		}
	}
}

func TestMatchHost(t *testing.T) {
	tests := []struct {
		name          string
		allowedHost   string
		requestedHost string
		shouldMatch   bool
	}{
		{"exact match", "example.com", "example.com", true},
		{"different domain", "example.com", "evil.com", false},
		{"subdomain with wildcard", "*.example.com", "api.example.com", true},
		{"subdomain without wildcard", "example.com", "api.example.com", false},
		{"localhost variations", "localhost", "127.0.0.1", true},
		{"localhost to ipv6", "localhost", "::1", true},
		{"port 80 normalization", "example.com:80", "example.com", true},
		{"port 443 normalization", "example.com:443", "example.com", true},
		{"custom port match", "example.com:8080", "example.com:8080", true},
		{"custom port mismatch", "example.com:8080", "example.com:3000", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchHost(tt.allowedHost, tt.requestedHost)
			if result != tt.shouldMatch {
				t.Errorf("matchHost(%q, %q) = %v, want %v",
					tt.allowedHost, tt.requestedHost, result, tt.shouldMatch)
			}
		})
	}
}

func TestMatchExact(t *testing.T) {
	tests := []struct {
		allowed   string
		requested string
		expected  bool
	}{
		{"PATH", "PATH", true},
		{"PATH", "HOME", false},
		{"API_KEY", "API_KEY", true},
		{"", "", true},
	}
	
	for _, tt := range tests {
		result := matchExact(tt.allowed, tt.requested)
		if result != tt.expected {
			t.Errorf("matchExact(%q, %q) = %v, want %v",
				tt.allowed, tt.requested, result, tt.expected)
		}
	}
}

func TestQuery(t *testing.T) {
	t.Run("denied by default", func(t *testing.T) {
		manager := NewManager()
		// Explicitly disable prompt mode for consistent test behavior
		manager.SetPromptMode(false)
		desc := PermissionDescriptor{
			Name:     PermissionRead,
			Resource: "/etc/passwd",
		}
		
		state := manager.Query(desc)
		if state != StateDenied {
			t.Errorf("expected StateDenied, got %v", state)
		}
	})
	
	t.Run("granted when allowed", func(t *testing.T) {
		manager := NewManager()
		manager.GrantRead([]string{"/tmp"})
		desc := PermissionDescriptor{
			Name:     PermissionRead,
			Resource: "/tmp/file.txt",
		}
		
		state := manager.Query(desc)
		if state != StateGranted {
			t.Errorf("expected StateGranted, got %v", state)
		}
	})
	
	t.Run("prompt mode", func(t *testing.T) {
		manager := NewManager()
		manager.SetPromptMode(true)
		desc := PermissionDescriptor{
			Name:     PermissionRead,
			Resource: "/etc/passwd",
		}
		
		state := manager.Query(desc)
		if state != StatePrompt {
			t.Errorf("expected StatePrompt, got %v", state)
		}
	})
}

func TestPermissionDescriptorString(t *testing.T) {
	tests := []struct {
		desc     PermissionDescriptor
		expected string
	}{
		{
			PermissionDescriptor{Name: PermissionRead, Resource: "/tmp/file.txt"},
			"read access to '/tmp/file.txt'",
		},
		{
			PermissionDescriptor{Name: PermissionWrite, Resource: ""},
			"write access",
		},
		{
			PermissionDescriptor{Name: PermissionNet, Resource: "example.com"},
			"net access to 'example.com'",
		},
	}
	
	for _, tt := range tests {
		result := tt.desc.String()
		if result != tt.expected {
			t.Errorf("String() = %q, want %q", result, tt.expected)
		}
	}
}

func TestErrorMessage(t *testing.T) {
	manager := NewManager()
	
	t.Run("read permission error", func(t *testing.T) {
		msg := manager.ErrorMessage(PermissionRead, "/etc/passwd")
		
		// Check that error message contains key information
		if msg == "" {
			t.Error("error message should not be empty")
		}
		// Should mention the permission type
		if len(msg) < 50 {
			t.Error("error message seems too short")
		}
	})
	
	t.Run("net permission error", func(t *testing.T) {
		msg := manager.ErrorMessage(PermissionNet, "example.com")
		
		if msg == "" {
			t.Error("error message should not be empty")
		}
	})
}

func TestGlobalManager(t *testing.T) {
	// Save and restore original global manager
	originalManager := globalManager
	defer func() {
		globalManager = originalManager
	}()
	
	t.Run("GetManager creates default", func(t *testing.T) {
		globalManager = nil
		manager := GetManager()
		
		if manager == nil {
			t.Fatal("GetManager should create a default manager")
		}
		if manager.allowRead != nil {
			t.Error("default manager should have nil allowRead")
		}
	})
	
	t.Run("SetGlobalManager and GetManager", func(t *testing.T) {
		custom := NewManager()
		custom.GrantAll()
		SetGlobalManager(custom)
		
		retrieved := GetManager()
		if retrieved != custom {
			t.Error("GetManager should return the set manager")
		}
		if retrieved.allowRead == nil {
			t.Error("should have the GrantAll settings")
		}
	})
}

func TestRealWorldPaths(t *testing.T) {
	manager := NewManager()
	
	// Get actual temp dir
	tmpDir := filepath.Clean("/tmp")
	manager.GrantRead([]string{tmpDir})
	
	// Test with real path
	testPath := filepath.Join(tmpDir, "test.txt")
	if !manager.Check(PermissionRead, testPath) {
		t.Errorf("should allow read to %s", testPath)
	}
}

func TestNetPermissionCheckSemantics(t *testing.T) {
	manager := NewManager()

	t.Run("localhost without port allows any loopback port", func(t *testing.T) {
		manager.GrantNet([]string{"localhost"})
		cases := []string{"localhost:3000", "127.0.0.1:9229", "[::1]:8080", "[::1]"}
		for _, host := range cases {
			if !manager.Check(PermissionNet, host) {
				t.Errorf("expected net access allowed for %s", host)
			}
		}
	})

	t.Run("localhost with port is strict", func(t *testing.T) {
		manager.GrantNet([]string{"localhost:3000"})
		if !manager.Check(PermissionNet, "localhost:3000") {
			t.Error("expected localhost:3000 to be allowed")
		}
		if manager.Check(PermissionNet, "localhost:9229") {
			t.Error("expected localhost:9229 to be denied")
		}
		if manager.Check(PermissionNet, "127.0.0.1:9229") {
			t.Error("expected 127.0.0.1:9229 to be denied")
		}
	})

	t.Run("wildcard domain matches subdomains and apex", func(t *testing.T) {
		manager.GrantNet([]string{"*.example.com"})
		if !manager.Check(PermissionNet, "api.example.com") {
			t.Error("expected api.example.com to be allowed")
		}
		if !manager.Check(PermissionNet, "example.com") {
			t.Error("expected example.com to be allowed by wildcard")
		}
		if manager.Check(PermissionNet, "evil.com") {
			t.Error("expected evil.com to be denied")
		}
		if manager.Check(PermissionNet, "app.evil.com") {
			t.Error("expected app.evil.com to be denied")
		}
	})

	t.Run("exact host without port does not allow custom port", func(t *testing.T) {
		manager.GrantNet([]string{"example.com"})
		if !manager.Check(PermissionNet, "example.com") {
			t.Error("expected example.com to be allowed")
		}
		if manager.Check(PermissionNet, "example.com:3000") {
			t.Error("expected example.com:3000 to be denied")
		}
		// default ports are normalized
		if !manager.Check(PermissionNet, "example.com:80") {
			t.Error("expected example.com:80 to be allowed (default port)")
		}
		if !manager.Check(PermissionNet, "example.com:443") {
			t.Error("expected example.com:443 to be allowed (default port)")
		}
	})
}

func TestGrantDefensiveCopy(t *testing.T) {
	manager := NewManager()

	// Net slice copy
	hosts := []string{"api.github.com"}
	manager.GrantNet(hosts)
	hosts[0] = "evil.com" // mutate original slice
	if !manager.Check(PermissionNet, "api.github.com") {
		t.Error("expected api.github.com still allowed after external slice mutation")
	}
	if manager.Check(PermissionNet, "evil.com") {
		t.Error("expected evil.com not allowed due to defensive copy")
	}

	// Env slice copy
	vars := []string{"PATH"}
	manager.GrantEnv(vars)
	vars[0] = "HOME"
	if !manager.Check(PermissionEnv, "PATH") {
		t.Error("expected PATH still allowed after external slice mutation (env)")
	}
	if manager.Check(PermissionEnv, "HOME") {
		t.Error("expected HOME not allowed due to defensive copy (env)")
	}

	// Run slice copy
	programs := []string{"git"}
	manager.GrantRun(programs)
	programs[0] = "rm"
	if !manager.Check(PermissionRun, "git") {
		t.Error("expected git still allowed after external slice mutation (run)")
	}
	if manager.Check(PermissionRun, "rm") {
		t.Error("expected rm not allowed due to defensive copy (run)")
	}
}
