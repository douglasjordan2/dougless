package permissions

import (
  "context"
  "fmt"
  "os"
  "path/filepath"
  "strings"
  "sync"
)

type Permission string

const (
  PermissionRead  Permission = "read"
  PermissionWrite Permission = "write"
  PermissionNet   Permission = "net"
  PermissionEnv   Permission = "env"
  PermissionRun   Permission = "run"
)

type PermissionState string

const (
  StateGranted PermissionState = "granted"
  StateDenied  PermissionState = "denied"
  StatePrompt  PermissionState = "prompt"
)

type PermissionDescriptor struct {
  Name     Permission
  Resource string
}

func (pd PermissionDescriptor) String() string {
  if pd.Resource != "" {
    return fmt.Sprintf("%s access to '%s'", pd.Name, pd.Resource)
  }
  return fmt.Sprintf("%s access", pd.Name)
}

type Manager struct {
  allowRead     *[]string
  allowWrite    *[]string
  allowNet      *[]string
  allowEnv      *[]string
  allowRun      *[]string
  promptMode    bool
  prompter      Prompter
  promptCache   map[string]PermissionState
  promptCacheMu sync.RWMutex
}

var globalManager *Manager

func IsTerminal() bool {
  fileInfo, err := os.Stdin.Stat()
  if err != nil {
    return false
  }
  return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

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

func SetGlobalManager(m *Manager) {
  globalManager = m
}

func GetManager() *Manager {
  if globalManager == nil {
    globalManager = NewManager()
  }
  return globalManager
}

func (m *Manager) GrantAll() {
  m.allowRead  = &[]string{}
  m.allowWrite = &[]string{}
  m.allowNet   = &[]string{}
  m.allowEnv   = &[]string{}
  m.allowRun   = &[]string{}
}

// each of the Grant* methods below work on the {paths,hosts,vars,etc} specified
// if none passed, access is granted to all {paths,hosts,vars,etc}

func (m*Manager) GrantRead(paths []string) {
  cp := append([]string(nil), paths...)
  m.allowRead = &cp
}

func (m*Manager) GrantWrite(paths []string) {
  cp := append([]string(nil), paths...)
  m.allowWrite = &cp
}

func (m*Manager) GrantNet(hosts []string) {
  cp := append([]string(nil), hosts...)
  m.allowNet = &cp
}

func (m*Manager) GrantEnv(vars []string) {
  cp := append([]string(nil), vars...)
  m.allowEnv = &cp
}

func (m*Manager) GrantRun(programs []string) {
  cp := append([]string(nil), programs...)
  m.allowRun = &cp
}

func (m *Manager) SetPromptMode(enabled bool) {
  m.promptMode = enabled
}

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

func normalizeLoopback(host string) string {
  switch host {
  case "localhost", "127.0.0.1", "::1":
    return "localhost"
  }
  return host
}

func normalizePort(port string) string {
  if port == "80" || port == "443" {
    return ""
  }
  return port
}

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
    if !hostMatches { return false }
    if aPort == "" { return true }
    return aPort == rPort
  }

  if aHost == rHost {
    if aPort == "" { return rPort == "" } // strict behavior for non-local hosts
    return aPort == rPort
  }

  return false
}

func matchExact(allowed, requested string) bool {
  return allowed == requested
}

func (m *Manager) Query(desc PermissionDescriptor) PermissionState {
  if m.Check(desc.Name, desc.Resource) {
    return StateGranted
  }

  if m.promptMode {
    return StatePrompt
  }

  return StateDenied
}

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

func cacheKey(perm Permission, resource string) string {
  return fmt.Sprintf("%s:%s", perm, resource)
}

func (m *Manager) SetPrompter(p Prompter) {
  m.prompter = p
}

func (m *Manager) ClearPromptCache() {
  m.promptCacheMu.Lock()
  defer m.promptCacheMu.Unlock()
  m.promptCache = make(map[string]PermissionState)
}

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
