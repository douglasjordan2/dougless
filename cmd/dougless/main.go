package main

import (
	"fmt"
	"os"

	"github.com/douglasjordan2/dougless/internal/runtime"
	"github.com/douglasjordan2/dougless/internal/repl"
)

func main() {
	rt := runtime.New()

  // go into repl mode if no args
  if len(os.Args) < 2 {
    r := repl.New(rt, os.Stdin, os.Stdout)
    if err := r.Run(); err != nil {
      fmt.Fprintf(os.Stderr, "REPL Error: %v\n", err)
      os.Exit(1)
    }
    return
  }

  // the only other accepted arg are .js files
	scriptPath := os.Args[1]
	if err := rt.ExecuteFile(scriptPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
