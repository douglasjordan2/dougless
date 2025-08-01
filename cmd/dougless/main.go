package main

import (
	"fmt"
	"os"

	"github.com/douglasjordan/dougless-runtime/internal/runtime"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: dougless <script.js>")
		os.Exit(1)
	}

	scriptPath := os.Args[1]
	
	rt := runtime.New()
	if err := rt.ExecuteFile(scriptPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
