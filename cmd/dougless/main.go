// Package main provides the Dougless Runtime CLI executable.
//
// Dougless is a custom JavaScript runtime built in Go, designed to execute
// modern JavaScript (ES6+) with automatic transpilation to ES5. It supports
// two modes of operation:
//
//  1. REPL Mode (no arguments): Interactive JavaScript shell
//  2. Script Mode: Execute JavaScript files
//
// Usage:
//
//	dougless [flags] [script.js]
//
// Flags:
//
//	--allow-read[=path]       Grant read access (optionally to specific paths)
//	--allow-write[=path]      Grant write access (optionally to specific paths)
//	--allow-net[=host]        Grant network access (optionally to specific hosts)
//	--allow-env[=var]         Grant environment variable access
//	--allow-run[=program]     Grant subprocess execution access
//	--allow-all               Grant all permissions (for development)
//
// Examples:
//
//	# Start interactive REPL
//	dougless
//
//	# Execute a script
//	dougless script.js
//
//	# Execute with specific permissions
//	dougless --allow-read=/tmp --allow-net=api.example.com script.js
//
//	# Execute with all permissions
//	dougless --allow-all script.js
package main

import (
	"fmt"
	"os"

	"github.com/douglasjordan2/dougless/internal/runtime"
	"github.com/douglasjordan2/dougless/internal/repl"
	"github.com/douglasjordan2/dougless/internal/permissions"
)

func main() {
  permManager, remainingArgs, err := permissions.ParseFlags(os.Args[1:])
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
    os.Exit(1)
  }

  permissions.SetGlobalManager(permManager)

	rt := runtime.New()

  // go into repl mode if no args
  if len(remainingArgs) == 0 {
    r := repl.New(rt, os.Stdin, os.Stdout)
    if err := r.Run(); err != nil {
      fmt.Fprintf(os.Stderr, "REPL Error: %v\n", err)
      os.Exit(1)
    }
    return
  }

  // the only other accepted arg are .js files
	scriptPath := remainingArgs[0]
	if err := rt.ExecuteFile(scriptPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
