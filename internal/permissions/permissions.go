// Package permissions provides fine-grained control over file system access,
// network operations, environment variables, and program execution.
//
// Permissions can be granted via CLI flags (--allow-read, --allow-net, etc.) or
// through interactive prompts when the runtime detects a terminal.
//
// Design principles:
//   - Secure by default: All operations require explicit permission
//   - Granular control: Permissions can be scoped to specific paths, hosts, or resources
//   - Interactive prompts: Users can grant permissions at runtime in terminal mode
//   - Clear error messages: Denied operations provide helpful suggestions
package permissions

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Permission represents a category of system access.
type Permission string

// Available permission categories.
const (
	PermissionRead  Permission = "read"  // File system read access
	PermissionWrite Permission = "write" // File system write access
	PermissionNet   Permission = "net"   // Network access (HTTP, WebSocket)
	PermissionEnv   Permission = "env"   // Environment variable access
	PermissionRun   Permission = "run"   // Program execution access
)

// PermissionState represents the granted/denied status of a permission.
type PermissionState string

// Permission states.
const (
	StateGranted PermissionState = "granted" // Permission is explicitly granted
	StateDenied  PermissionState = "denied"  // Permission is explicitly denied
	StatePrompt  PermissionState = "prompt"  // Permission requires interactive prompt
)

// PermissionDescriptor describes a specific permission request.
// It includes the permission type and the resource being accessed.
type PermissionDescriptor struct {
	Name     Permission // The permission category (read, write, net, etc.)
	Resource string     // The specific resource (file path, hostname, etc.)
}

// String returns a human-readable description of the permission descriptor.
func (pd PermissionDescriptor) String() string {
	if pd.Resource != "" {
		return fmt.Sprintf("%s access to '%s'", pd.Name, pd.Resource)
	}
	return fmt.Sprintf("%s access", pd.Name)
}

// Manager handles permission checks and interactive prompts.
// It maintains allow lists for each permission type and manages
// prompt caching for the current session.
//
// The manager is typically initialized from CLI flags and accessed
// globally throughout the runtime.
type Manager struct {
	allowRead     *[]string                  // Allowed read paths (nil = denied, empty = all)
	allowWrite    *[]string                  // Allowed write paths
	allowNet      *[]string                  // Allowed network hosts
	allowEnv      *[]string                  // Allowed environment variables
	allowRun      *[]string                  // Allowed programs to execute
	promptMode    bool                       // Whether to prompt for permissions
	prompter      Prompter                   // Interface for prompting user
	promptCache   map[string]PermissionState // Cache of user responses
	promptCacheMu sync.RWMutex               // Protects promptCache
}

// globalManager is the singleton permission manager instance.
var globalManager *Manager

// IsTerminal checks if stdin is connected to a terminal.
// This determines whether interactive prompts are available.
func IsTerminal() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// NewManager creates a new permission manager with default settings.
// Prompt mode is automatically enabled if stdin is a terminal.
func NewManager() *Manager {
	return &Manager{
		allowRead:   nil,
		allowWrite:  nil,
		allowNet:    nil,
		allowEnv:    nil,
		allowRun:    nil,
		promptMode:  IsTerminal(),
		prompter:    NewStdioPrompter(),
		promptCache: make(map[string]PermissionState),
	}
}

// SetGlobalManager sets the global permission manager instance.
// This should be called once during runtime initialization.
func SetGlobalManager(m *Manager) {
	globalManager = m
}

// GetManager returns the global permission manager, creating one if needed.
func GetManager() *Manager {
	if globalManager == nil {
		globalManager = NewManager()
	}
	return globalManager
}

// GrantAll grants all permission types without restriction.
// This is equivalent to --allow-all and should only be used in development.
func (m *Manager) GrantAll() {
	m.allowRead = &[]string{}
	m.allowWrite = &[]string{}
	m.allowNet = &[]string{}
	m.allowEnv = &[]string{}
	m.allowRun = &[]string{}
}

// GrantRead grants read permission for the specified paths.
// If paths is empty, all read access is granted.
// If paths contains specific paths, only those paths (and their subdirectories) are allowed.
func (m *Manager) GrantRead(paths []string) {
	cp := append([]string(nil), paths...)
	m.allowRead = &cp
}

// GrantWrite grants write permission for the specified paths.
// If paths is empty, all write access is granted.
func (m *Manager) GrantWrite(paths []string) {
	cp := append([]string(nil), paths...)
	m.allowWrite = &cp
}

// GrantNet grants network permission for the specified hosts.
// Supports wildcards (*.example.com) and port specifications (localhost:3000).
// If hosts is empty, all network access is granted.
func (m *Manager) GrantNet(hosts []string) {
	cp := append([]string(nil), hosts...)
	m.allowNet = &cp
}

// GrantEnv grants environment variable access for the specified variables.
// If vars is empty, all environment variables are accessible.
func (m *Manager) GrantEnv(vars []string) {
	cp := append([]string(nil), vars...)
	m.allowEnv = &cp
}

// GrantRun grants permission to execute the specified programs.
// If programs is empty, all program execution is allowed.
func (m *Manager) GrantRun(programs []string) {
	cp := append([]string(nil), programs...)
	m.allowRun = &cp
}

// SetPromptMode enables or disables interactive permission prompts.
// When enabled and stdin is a terminal, users are prompted for missing permissions.
func (m *Manager) SetPromptMode(enabled bool) {
	m.promptMode = enabled
}

// Check verifies if a permission is granted for the specified resource.
// Returns true if the permission is explicitly granted, false otherwise.
// This method does NOT trigger interactive prompts.
func (m *Manager) Check(perm Permission, resource string) bool {
	switch perm {
	case PermissionRead:
		return m.checkPermission(m.allowRead, resource, matchPath)
	case PermissionWrite:
		return m.checkPermission(m.allowWrite, resource, matchPath)
	case PermissionNet:
		return m.checkPermission(m.allowNet, resource, matchHost)
	case PermissionEnv:
		return m.checkPermission(m.allowEnv, resource, matchExact)
	case PermissionRun:
		return m.checkPermission(m.allowRun, resource, matchExact)
	default:
		return false
	}
}

// checkPermission is the internal implementation for permission checking.
// It uses a matcher function to compare allowed resources with the requested resource.
func (m *Manager) checkPermission(allowList *[]string, resource string, matcher func(string, string) bool) bool {
	if allowList == nil {
		return false
	}

	if len(*allowList) == 0 {
		return true
	}

	for _, allowed := range *allowList {
		if matcher(allowed, resource) {
			return true
		}
	}

	return false
}

// matchPath checks if a requested path is allowed based on an allowed path.
// Supports directory hierarchies: if /home/user is allowed, /home/user/file.txt is also allowed.
// Prevents directory traversal attacks by checking for ".." in relative paths.
func matchPath(allowedPath, requestedPath string) bool {
	allowed, err := filepath.Abs(filepath.Clean(allowedPath))
	if err != nil {
		return false
	}

	requested, err := filepath.Abs(filepath.Clean(requestedPath))
	if err != nil {
		return false
	}

	if requested == allowed {
		return true
	}

	rel, err := filepath.Rel(allowed, requested)
	if err != nil {
		return false
	}

	// prevent escape
	if strings.HasPrefix(rel, "..") {
		return false
	}

	return true
}

// splitHostPort parses a host:port string, handling IPv6 addresses correctly.
// Supports formats: localhost:3000, example.com:443, [::1]:8080, ::1
func splitHostPort(h string) (host, port string) {
	// Bracketed IPv6 with optional port: [::1]:3000 or [::1]
	if strings.HasPrefix(h, "[") {
		if i := strings.Index(h, "]"); i != -1 {
			base := h[1:i]
			rest := h[i+1:]
			if strings.HasPrefix(rest, ":") {
				return base, rest[1:]
			}
			return base, ""
		}
	}

	// Raw IPv6 without brackets should not be split by colon
	if strings.Count(h, ":") >= 2 && !strings.HasPrefix(h, "[") && !strings.Contains(h, "]") {
		return h, ""
	}

	// host:port
	if i := strings.LastIndex(h, ":"); i > 0 && i < len(h)-1 && !strings.Contains(h[i+1:], ":") {
		return h[:i], h[i+1:]
	}

	return h, ""
}

// normalizeLoopback converts various loopback addresses to a canonical form.
func normalizeLoopback(host string) string {
	switch host {
	case "localhost", "127.0.0.1", "::1":
		return "localhost"
	}
	return host
}

// normalizePort removes default ports (80, 443) for matching purposes.
func normalizePort(port string) string {
	if port == "80" || port == "443" {
		return ""
	}
	return port
}

// matchHost checks if a requested host is allowed based on an allowed host pattern.
// Supports wildcards (*.example.com), localhost with any port, and specific host:port combinations.
func matchHost(allowedHost, requestedHost string) bool {
	aHost, aPort := splitHostPort(allowedHost)
	rHost, rPort := splitHostPort(requestedHost)

	aHost = normalizeLoopback(aHost)
	rHost = normalizeLoopback(rHost)
	aPort = normalizePort(aPort)
	rPort = normalizePort(rPort)

	if aHost == "localhost" {
		if rHost == "localhost" {
			if aPort == "" { // --allow-net=localhost
				return true // any port on loopback
			}
			return aPort == rPort // --allow-net=localhost:3000
		}
		return false
	}

	if strings.HasPrefix(aHost, "*.") {
		domain := strings.TrimPrefix(aHost, "*.")
		hostMatches := rHost == domain || strings.HasSuffix(rHost, "."+domain)
		if !hostMatches {
			return false
		}
		if aPort == "" {
			return true
		}
		return aPort == rPort
	}

	if aHost == rHost {
		if aPort == "" {
			return rPort == ""
		} // strict behavior for non-local hosts
		return aPort == rPort
	}

	return false
}

// matchExact performs exact string matching (used for env vars and program names).
func matchExact(allowed, requested string) bool {
	return allowed == requested
}

// Query determines the current state of a permission.
// Returns StateGranted if allowed, StatePrompt if prompting is enabled, or StateDenied otherwise.
func (m *Manager) Query(desc PermissionDescriptor) PermissionState {
	if m.Check(desc.Name, desc.Resource) {
		return StateGranted
	}

	if m.promptMode {
		return StatePrompt
	}

	return StateDenied
}

// ErrorMessage generates a helpful error message for permission denials.
// Includes examples of how to grant the required permission via CLI flags.
func (m *Manager) ErrorMessage(perm Permission, resource string) string {
	var flag string
	var example string

	switch perm {
	case PermissionRead:
		flag = "--allow-read"
		if resource != "" {
			example = fmt.Sprintf("  dougless --allow-read=%s script.js", resource)
		} else {
			example = "  dougless --allow-read script.js"
		}
	case PermissionWrite:
		flag = "--allow-write"
		if resource != "" {
			example = fmt.Sprintf("  dougless --allow-write=%s script.js", resource)
		} else {
			example = "  dougless --allow-write script.js"
		}
	case PermissionNet:
		flag = "--allow-net"
		if resource != "" {
			example = fmt.Sprintf("  dougless --allow-net=%s script.js", resource)
		} else {
			example = "  dougless --allow-net script.js"
		}
	case PermissionEnv:
		flag = "--allow-env"
		if resource != "" {
			example = fmt.Sprintf("  dougless --allow-env=%s script.js", resource)
		} else {
			example = "  dougless --allow-env script.js"
		}
	case PermissionRun:
		flag = "--allow-run"
		if resource != "" {
			example = fmt.Sprintf("  dougless --allow-run=%s script.js", resource)
		} else {
			example = "  dougless --allow-run script.js"
		}
	}

	desc := PermissionDescriptor{Name: perm, Resource: resource}

	msg := fmt.Sprintf("Permission denied: %s\n\n", desc)
	msg += fmt.Sprintf("Run your script with:\n%s\n\n", example)
	msg += fmt.Sprintf("Or grant all %s access:\n  dougless %s script.js\n", perm, flag)
	msg += fmt.Sprintf("\nFor dev, use:\n  dougless --allow-all script.js")

	return msg
}

// cacheKey generates a unique key for permission caching.
func cacheKey(perm Permission, resource string) string {
	return fmt.Sprintf("%s:%s", perm, resource)
}

// SetPrompter replaces the default prompter with a custom implementation.
// Useful for testing or alternative UI implementations.
func (m *Manager) SetPrompter(p Prompter) {
	m.prompter = p
}

// ClearPromptCache clears all cached prompt responses.
// Useful for testing or when starting a new interactive session.
func (m *Manager) ClearPromptCache() {
	m.promptCacheMu.Lock()
	defer m.promptCacheMu.Unlock()
	m.promptCache = make(map[string]PermissionState)
}

// CheckWithPrompt checks a permission and prompts the user if needed.
// If the permission is not granted and prompt mode is enabled, the user is prompted.
// Responses can be cached for the session if the user selects "always".
//
// This is the primary method used by runtime operations to check permissions.
func (m *Manager) CheckWithPrompt(ctx context.Context, perm Permission, resource string) bool {
	if m.Check(perm, resource) {
		return true
	}

	if !m.promptMode {
		return false
	}

	key := cacheKey(perm, resource)
	m.promptCacheMu.RLock()
	if state, exists := m.promptCache[key]; exists {
		m.promptCacheMu.RUnlock()
		return state == StateGranted
	}
	m.promptCacheMu.RUnlock()

	desc := PermissionDescriptor{Name: perm, Resource: resource}
	response, err := m.prompter.Prompt(ctx, desc)
	if err != nil {
		return false
	}

	if response.Permanent {
		state := StateDenied
		if response.Granted {
			state = StateGranted
		}

		m.promptCacheMu.Lock()
		m.promptCache[key] = state
		m.promptCacheMu.Unlock()
	}

	return response.Granted
}
