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
