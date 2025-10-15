// Package repl provides an interactive Read-Eval-Print Loop for the Dougless runtime.
//
// The REPL allows users to interactively execute JavaScript code, with features including:
//   - Multi-line input support with automatic bracket/brace detection
//   - Command history (up/down arrows)
//   - Special commands (.help, .exit, .clear)
//   - State preservation between evaluations
//   - Proper error display with Goja exception handling
//
// Example usage:
//
//	rt := runtime.New()
//	r := repl.New(rt, os.Stdin, os.Stdout)
//	if err := r.Run(); err != nil {
//	    log.Fatal(err)
//	}
package repl

import (
	"fmt"
	"io"
	"strings"

	"github.com/dop251/goja"
	"github.com/peterh/liner"

	"github.com/douglasjordan2/dougless/internal/runtime"
)

// REPL represents an interactive JavaScript shell.
// It maintains state between evaluations and supports multi-line input.
type REPL struct {
	runtime *runtime.Runtime // JavaScript runtime for code execution
	line    *liner.State     // Liner instance for input handling and history
	writer  io.Writer        // Output writer for results and messages
}

// New creates a new REPL instance with the given runtime and I/O streams.
//
// The reader parameter is included for API compatibility but currently unused
// as liner handles input directly from the terminal.
//
// Example:
//
//	rt := runtime.New()
//	repl := repl.New(rt, os.Stdin, os.Stdout)
func New(rt *runtime.Runtime, reader io.Reader, writer io.Writer) *REPL {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)

	return &REPL{
		runtime: rt,
		line:    line,
		writer:  writer,
	}
}

// isIncompleteInput detects if the user's input is incomplete (unmatched brackets).
// This enables multi-line input support by checking for unclosed:
//   - Braces: { }
//   - Brackets: [ ]
//   - Parentheses: ( )
//
// Returns true if there are unmatched opening delimiters.
func (r *REPL) isIncompleteInput(input string) bool {
	input = strings.TrimSpace(input)
	if input == "" {
		return false
	}

	openBraces := strings.Count(input, "{") - strings.Count(input, "}")
	openBrackets := strings.Count(input, "[") - strings.Count(input, "]")
	openParens := strings.Count(input, "(") - strings.Count(input, ")")

	return openBraces > 0 || openBrackets > 0 || openParens > 0
}

// printWelcome displays the welcome message when the REPL starts.
func (r *REPL) printWelcome() {
	fmt.Fprintln(r.writer, "Dougless Runtime REPL v0.1.0")
	fmt.Fprintln(r.writer, "type some JS, use `.help`, or quit with `.exit`")
	fmt.Fprintln(r.writer, "")
}

// handleCommand processes REPL special commands (those starting with '.').
// Returns true if the REPL should exit, false to continue.
//
// Supported commands:
//
//	.exit, .quit - Exit the REPL
//	.help        - Display help message
//	.clear       - Clear the screen
func (r *REPL) handleCommand(cmd string) bool {
	switch cmd {
	case ".exit", ".quit":
		fmt.Fprintln(r.writer, "see ya")
		return true
	case ".help":
		r.printHelp()
		return false
	case ".clear":
		fmt.Fprint(r.writer, "\033[H\033[2J")
		return false
	default:
		fmt.Fprintf(r.writer, "Unknown command: %s (type .help for available commands)\n", cmd)
		return false
	}
}

// printHelp displays the help message with available commands.
func (r *REPL) printHelp() {
	fmt.Fprintln(r.writer, "Available commands:")
	fmt.Fprintln(r.writer, "  .help   - Show this help message")
	fmt.Fprintln(r.writer, "  .exit   - Exit the REPL (or Ctrl+D)")
	fmt.Fprintln(r.writer, "  .quit   - Same as .exit")
	fmt.Fprintln(r.writer, "  .clear  - Clear the screen")
	fmt.Fprintln(r.writer, "")
}

// Run starts the REPL loop and processes user input until exit.
//
// The REPL:
//  1. Displays a welcome message
//  2. Prompts for input (> for single-line, ... for multi-line)
//  3. Evaluates JavaScript code
//  4. Prints results or errors
//  5. Repeats until .exit command or EOF (Ctrl+D)
//
// Returns an error if there's a problem with I/O, or nil on normal exit.
func (r *REPL) Run() error {
	defer r.line.Close()

	r.printWelcome()

	var multilineBuffer strings.Builder
	inMultiline := false

	for {
		prompt := "> "
		if inMultiline {
			prompt = "... "
		}

		line, err := r.line.Prompt(prompt)
		if err != nil {
			if err == liner.ErrPromptAborted {
				fmt.Fprintln(r.writer, "\nsee ya")
				return nil
			}

			if err == io.EOF {
				fmt.Fprintln(r.writer, "\nsee ya")
				return nil
			}
			return err
		}

		line = strings.TrimSpace(line)

		if !inMultiline && strings.HasPrefix(line, ".") {
			if r.handleCommand(line) {
				return nil
			}

			continue
		}

		if inMultiline {
			multilineBuffer.WriteString(line)
			multilineBuffer.WriteString("\n")
		} else {
			multilineBuffer.WriteString(line)
		}

		currentInput := multilineBuffer.String()

		if r.isIncompleteInput(currentInput) {
			inMultiline = true
			continue
		}

		if !inMultiline && line != "" {
			r.line.AppendHistory(line)
		}

		result, err := r.runtime.Evaluate(currentInput)
		if err != nil {
			if jsErr, ok := err.(*goja.Exception); ok {
				fmt.Fprintf(r.writer, "Error: %s\n", jsErr.String())
			} else {
				fmt.Fprintf(r.writer, "Error: %v\n", err)
			}
		} else if result != nil && !goja.IsUndefined(result) {
			fmt.Fprintf(r.writer, "%v\n", result)
		}

		multilineBuffer.Reset()
		inMultiline = false
	}
}
