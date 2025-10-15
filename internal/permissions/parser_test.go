package permissions

import (
	"testing"
)

func TestParseFlags(t *testing.T) {
	t.Run("no flags", func(t *testing.T) {
		args := []string{"script.js"}
		manager, remaining, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if manager == nil {
			t.Fatal("manager should not be nil")
		}
		if len(remaining) != 1 || remaining[0] != "script.js" {
			t.Errorf("expected [script.js], got %v", remaining)
		}

		// Should be denied by default
		if manager.Check(PermissionRead, "/tmp") {
			t.Error("read should be denied by default")
		}
	})

	t.Run("allow-all flag", func(t *testing.T) {
		args := []string{"--allow-all", "script.js"}
		manager, remaining, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(remaining) != 1 || remaining[0] != "script.js" {
			t.Errorf("expected [script.js], got %v", remaining)
		}

		// Everything should be allowed
		if !manager.Check(PermissionRead, "/etc/passwd") {
			t.Error("read should be granted with --allow-all")
		}
		if !manager.Check(PermissionWrite, "/tmp/test.txt") {
			t.Error("write should be granted with --allow-all")
		}
		if !manager.Check(PermissionNet, "example.com") {
			t.Error("net should be granted with --allow-all")
		}
	})

	t.Run("allow-all short flag", func(t *testing.T) {
		args := []string{"-A", "script.js"}
		manager, _, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !manager.Check(PermissionRead, "/tmp") {
			t.Error("-A should grant all permissions")
		}
	})

	t.Run("allow-read all", func(t *testing.T) {
		args := []string{"--allow-read", "script.js"}
		manager, _, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Read should be allowed everywhere
		if !manager.Check(PermissionRead, "/etc/passwd") {
			t.Error("read should be allowed")
		}
		if !manager.Check(PermissionRead, "/tmp/test.txt") {
			t.Error("read should be allowed")
		}

		// Write should still be denied
		if manager.Check(PermissionWrite, "/tmp/test.txt") {
			t.Error("write should still be denied")
		}
	})

	t.Run("allow-read specific path", func(t *testing.T) {
		args := []string{"--allow-read=/tmp", "script.js"}
		manager, _, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should allow /tmp
		if !manager.Check(PermissionRead, "/tmp/test.txt") {
			t.Error("/tmp should be allowed")
		}

		// Should deny /etc
		if manager.Check(PermissionRead, "/etc/passwd") {
			t.Error("/etc should be denied")
		}
	})

	t.Run("allow-read multiple paths", func(t *testing.T) {
		args := []string{"--allow-read=/tmp,/home/user", "script.js"}
		manager, _, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !manager.Check(PermissionRead, "/tmp/file.txt") {
			t.Error("/tmp should be allowed")
		}
		if !manager.Check(PermissionRead, "/home/user/doc.txt") {
			t.Error("/home/user should be allowed")
		}
		if manager.Check(PermissionRead, "/etc/passwd") {
			t.Error("/etc should be denied")
		}
	})

	t.Run("allow-write specific path", func(t *testing.T) {
		args := []string{"--allow-write=/tmp", "script.js"}
		manager, _, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !manager.Check(PermissionWrite, "/tmp/output.txt") {
			t.Error("/tmp write should be allowed")
		}
		if manager.Check(PermissionWrite, "/etc/config") {
			t.Error("/etc write should be denied")
		}
	})

	t.Run("allow-net all", func(t *testing.T) {
		args := []string{"--allow-net", "script.js"}
		manager, _, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !manager.Check(PermissionNet, "example.com") {
			t.Error("net should be allowed everywhere")
		}
		if !manager.Check(PermissionNet, "api.github.com") {
			t.Error("net should be allowed everywhere")
		}
	})

	t.Run("allow-net specific host", func(t *testing.T) {
		args := []string{"--allow-net=example.com", "script.js"}
		manager, _, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !manager.Check(PermissionNet, "example.com") {
			t.Error("example.com should be allowed")
		}
		if manager.Check(PermissionNet, "evil.com") {
			t.Error("evil.com should be denied")
		}
	})

	t.Run("allow-net multiple hosts", func(t *testing.T) {
		args := []string{"--allow-net=api.github.com,example.com", "script.js"}
		manager, _, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !manager.Check(PermissionNet, "api.github.com") {
			t.Error("api.github.com should be allowed")
		}
		if !manager.Check(PermissionNet, "example.com") {
			t.Error("example.com should be allowed")
		}
		if manager.Check(PermissionNet, "evil.com") {
			t.Error("evil.com should be denied")
		}
	})

	t.Run("allow-env all", func(t *testing.T) {
		args := []string{"--allow-env", "script.js"}
		manager, _, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !manager.Check(PermissionEnv, "PATH") {
			t.Error("env access should be allowed")
		}
	})

	t.Run("allow-env specific vars", func(t *testing.T) {
		args := []string{"--allow-env=PATH,HOME", "script.js"}
		manager, _, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !manager.Check(PermissionEnv, "PATH") {
			t.Error("PATH should be allowed")
		}
		if !manager.Check(PermissionEnv, "HOME") {
			t.Error("HOME should be allowed")
		}
		if manager.Check(PermissionEnv, "API_KEY") {
			t.Error("API_KEY should be denied")
		}
	})

	t.Run("allow-run all", func(t *testing.T) {
		args := []string{"--allow-run", "script.js"}
		manager, _, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !manager.Check(PermissionRun, "ls") {
			t.Error("run should be allowed")
		}
	})

	t.Run("allow-run specific programs", func(t *testing.T) {
		args := []string{"--allow-run=git,node", "script.js"}
		manager, _, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !manager.Check(PermissionRun, "git") {
			t.Error("git should be allowed")
		}
		if !manager.Check(PermissionRun, "node") {
			t.Error("node should be allowed")
		}
		if manager.Check(PermissionRun, "rm") {
			t.Error("rm should be denied")
		}
	})

	t.Run("prompt mode", func(t *testing.T) {
		args := []string{"--prompt", "script.js"}
		manager, _, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		desc := PermissionDescriptor{
			Name:     PermissionRead,
			Resource: "/tmp/file.txt",
		}

		state := manager.Query(desc)
		if state != StatePrompt {
			t.Errorf("expected StatePrompt, got %v", state)
		}
	})

	t.Run("multiple flags", func(t *testing.T) {
		args := []string{
			"--allow-read=/tmp",
			"--allow-write=/tmp",
			"--allow-net=localhost",
			"script.js",
		}
		manager, remaining, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(remaining) != 1 || remaining[0] != "script.js" {
			t.Errorf("expected [script.js], got %v", remaining)
		}

		if !manager.Check(PermissionRead, "/tmp/test.txt") {
			t.Error("read /tmp should be allowed")
		}
		if !manager.Check(PermissionWrite, "/tmp/test.txt") {
			t.Error("write /tmp should be allowed")
		}
		if !manager.Check(PermissionNet, "localhost") {
			t.Error("net localhost should be allowed")
		}
		if manager.Check(PermissionNet, "example.com") {
			t.Error("net example.com should be denied")
		}
	})

	t.Run("flags before and after script", func(t *testing.T) {
		args := []string{"--allow-read=/tmp", "script.js", "arg1", "arg2"}
		manager, remaining, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify permission was granted
		if !manager.Check(PermissionRead, "/tmp/test.txt") {
			t.Error("read permission should be granted")
		}

		// Should preserve non-flag arguments
		expected := []string{"script.js", "arg1", "arg2"}
		if len(remaining) != len(expected) {
			t.Fatalf("expected %d args, got %d", len(expected), len(remaining))
		}
		for i, arg := range expected {
			if remaining[i] != arg {
				t.Errorf("arg[%d]: expected %q, got %q", i, arg, remaining[i])
			}
		}
	})
}

func TestParsePermissionValue(t *testing.T) {
	t.Run("no equals - allow all", func(t *testing.T) {
		values, err := parsePermissionValue("--allow-read", "--allow-read")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(values) != 0 {
			t.Errorf("expected empty slice for 'allow all', got %v", values)
		}
	})

	t.Run("single value", func(t *testing.T) {
		values, err := parsePermissionValue("--allow-read=/tmp", "--allow-read")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(values) != 1 || values[0] != "/tmp" {
			t.Errorf("expected [/tmp], got %v", values)
		}
	})

	t.Run("multiple values", func(t *testing.T) {
		values, err := parsePermissionValue("--allow-read=/tmp,/home", "--allow-read")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(values) != 2 {
			t.Fatalf("expected 2 values, got %d", len(values))
		}
		if values[0] != "/tmp" || values[1] != "/home" {
			t.Errorf("expected [/tmp /home], got %v", values)
		}
	})

	t.Run("values with spaces", func(t *testing.T) {
		values, err := parsePermissionValue("--allow-read=/tmp, /home , /var", "--allow-read")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(values) != 3 {
			t.Fatalf("expected 3 values, got %d", len(values))
		}
		// Spaces should be trimmed
		if values[0] != "/tmp" || values[1] != "/home" || values[2] != "/var" {
			t.Errorf("expected trimmed values, got %v", values)
		}
	})

	t.Run("empty value after equals", func(t *testing.T) {
		_, err := parsePermissionValue("--allow-read=", "--allow-read")
		if err == nil {
			t.Error("expected error for empty value after equals")
		}
	})
}

func TestComplexScenarios(t *testing.T) {
	t.Run("development workflow", func(t *testing.T) {
		// Typical dev usage with --allow-all
		args := []string{"--allow-all", "app.js"}
		manager, _, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Everything should work
		tests := []struct {
			perm     Permission
			resource string
		}{
			{PermissionRead, "/home/user/project/config.json"},
			{PermissionWrite, "/home/user/project/output.txt"},
			{PermissionNet, "api.example.com"},
			{PermissionEnv, "DATABASE_URL"},
		}

		for _, tt := range tests {
			if !manager.Check(tt.perm, tt.resource) {
				t.Errorf("%s %s should be allowed in dev mode", tt.perm, tt.resource)
			}
		}
	})

	t.Run("production workflow", func(t *testing.T) {
		// Strict production permissions
		args := []string{
			"--allow-read=/app/config,/app/data",
			"--allow-write=/app/logs",
			"--allow-net=api.internal.com",
			"--allow-env=NODE_ENV,PORT",
			"app.js",
		}
		manager, _, err := ParseFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Allowed operations
		if !manager.Check(PermissionRead, "/app/config/settings.json") {
			t.Error("should read config")
		}
		if !manager.Check(PermissionWrite, "/app/logs/app.log") {
			t.Error("should write logs")
		}
		if !manager.Check(PermissionNet, "api.internal.com") {
			t.Error("should access internal API")
		}

		// Denied operations
		if manager.Check(PermissionRead, "/etc/passwd") {
			t.Error("should not read /etc/passwd")
		}
		if manager.Check(PermissionWrite, "/app/config/settings.json") {
			t.Error("should not write to config directory")
		}
		if manager.Check(PermissionNet, "evil.com") {
			t.Error("should not access external hosts")
		}
	})
}
