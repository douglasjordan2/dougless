package permissions

import (
  "fmt"
  "os"
  "strings"
)

func ParseFlags(args []string) (*Manager, []string, error) {
  manager := NewManager()
  remainingArgs := []string{}
  allowAll := false

  for i := 0; i < len(args); i++ {
    arg := args[i]

    if strings.HasPrefix(arg, "--allow-all") || arg == "-A" {
      allowAll = true
    } else if strings.HasPrefix(arg, "--allow-read") {
      paths, err := parsePermissionValue(arg, "--allow-read")
      if err != nil {
        return nil, nil, err
      }
      manager.GrantRead(paths)
    } else if strings.HasPrefix(arg, "--allow-write") {
      paths, err := parsePermissionValue(arg, "--allow-write")
      if err != nil {
        return nil, nil, err
      }
      manager.GrantWrite(paths)
    } else if strings.HasPrefix(arg, "--allow-net") {
      hosts, err := parsePermissionValue(arg, "--allow-net")
      if err != nil {
        return nil, nil, err
      }
      manager.GrantNet(hosts)
    } else if strings.HasPrefix(arg, "--allow-env") {
      vars, err := parsePermissionValue(arg, "--allow-env")
      if err != nil {
        return nil, nil, err
      }
      manager.GrantEnv(vars)
    } else if strings.HasPrefix(arg, "--allow-run") {
      programs, err := parsePermissionValue(arg, "--allow-run")
      if err != nil {
        return nil, nil, err
      }
      manager.GrantRun(programs)
    } else if arg == "--prompt" {
      manager.SetPromptMode(true)
    } else if arg == "--no-prompt" {
      manager.SetPromptMode(false)
    } else {
      remainingArgs = append(remainingArgs, arg)
    }
  }

  if allowAll {
    fmt.Fprintln(os.Stderr, "⚠️  WARNING: Running with --allow-all grants full system access")
		fmt.Fprintln(os.Stderr, "   This is convenient for development but NOT recommended for production.")
		fmt.Fprintln(os.Stderr, "   Consider using specific permissions: --allow-read, --allow-write, --allow-net")
    manager.GrantAll()
  }

  return manager, remainingArgs, nil
}

func parsePermissionValue(arg, flagName string) ([]string, error) {
  if !strings.Contains(arg, "=") {
    return []string{}, nil // nothing specified means allow all
  }

  parts := strings.SplitN(arg, "=", 2)
  if len(parts) != 2 {
    return nil, fmt.Errorf("invalid flag format: %s", arg)
  }

  value := parts[1]

  if value == "" {
    return nil, fmt.Errorf("%s requires a value or omit '=' to allow all", flagName)
  }

  values := strings.Split(value, ",")

  for i, v := range values {
    values[i] = strings.TrimSpace(v)
  }

  return values, nil
}
