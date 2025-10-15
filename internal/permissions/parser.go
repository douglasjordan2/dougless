// Package permissions parser handles CLI flag parsing for permission grants.
// Supports flags like --allow-read, --allow-net, --allow-all with optional values.
package permissions

import (
	"fmt"
	"os"
	"strings"
)

// ParseFlags parses command-line arguments and extracts permission flags.
// Returns a configured Manager, remaining non-permission arguments, and any parse errors.
//
// Supported flags:
//
//	--allow-all or -A: Grant all permissions (warns about security implications)
//	--allow-read[=paths]: Grant read permission (comma-separated paths, empty = all)
//	--allow-write[=paths]: Grant write permission
//	--allow-net[=hosts]: Grant network permission (supports wildcards and ports)
//	--allow-env[=vars]: Grant environment variable access
//	--allow-run[=programs]: Grant program execution permission
//	--prompt: Force enable interactive prompts
//	--no-prompt: Disable interactive prompts
//
// Examples:
//
//	dougless --allow-read script.js                    (all read access)
//	dougless --allow-read=.,/tmp script.js             (specific paths)
//	dougless --allow-net=localhost:3000 script.js      (specific host:port)
//	dougless --allow-all script.js                     (all permissions)
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

// parsePermissionValue extracts values from permission flags.
// Handles flags in formats: --flag (allow all), --flag= (error), --flag=val1,val2
// Returns a slice of trimmed values or an error for invalid formats.
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
