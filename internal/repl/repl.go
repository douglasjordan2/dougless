package repl

import (
  "fmt"
  "io"
  "strings"

  "github.com/dop251/goja"
  "github.com/peterh/liner"

  "github.com/douglasjordan2/dougless/internal/runtime"
)

type REPL struct {
  runtime *runtime.Runtime
  line *liner.State
  writer  io.Writer
}

func New(rt *runtime.Runtime, reader io.Reader, writer io.Writer) *REPL {
  line := liner.NewLiner()
  line.SetCtrlCAborts(true)

  return &REPL{
    runtime: rt,
    line:    line,
    writer:  writer,
  }
}

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

func (r *REPL) printWelcome() {
  fmt.Fprintln(r.writer, "Dougless Runtime REPL v0.1.0")
  fmt.Fprintln(r.writer, "type some JS, use `.help`, or quit with `.exit`")
  fmt.Fprintln(r.writer, "")
}

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

func (r *REPL) printHelp() {
  fmt.Fprintln(r.writer, "Available commands:")
  fmt.Fprintln(r.writer, "  .help   - Show this help message")
	fmt.Fprintln(r.writer, "  .exit   - Exit the REPL (or Ctrl+D)")
	fmt.Fprintln(r.writer, "  .quit   - Same as .exit")
	fmt.Fprintln(r.writer, "  .clear  - Clear the screen")
	fmt.Fprintln(r.writer, "")
}

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
