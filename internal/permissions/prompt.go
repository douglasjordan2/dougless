package permissions

import (
  "bufio"
  "context"
  "fmt"
  "os"
  "strings"
  "sync"
)

type PromptResponse struct {
  Granted   bool
  Permanent bool
}

type Prompter interface {
  Prompt(ctx context.Context, desc PermissionDescriptor) (PromptResponse, error)
}

type StdioPrompter struct {
  mu sync.Mutex // Serialize prompts to prevent concurrent stdin reads
}

func NewStdioPrompter() *StdioPrompter {
  return &StdioPrompter{}
}

func (p *StdioPrompter) Prompt(ctx context.Context, desc PermissionDescriptor) (PromptResponse, error) {
  // Serialize prompts to prevent multiple concurrent stdin reads
  p.mu.Lock()
  defer p.mu.Unlock()

  responseChan := make(chan PromptResponse, 1)
  errorChan := make(chan error, 1)

  go func() {
    fmt.Fprintf(os.Stderr, "\n⚠️  Permission request: %s\n", desc)
    fmt.Fprintf(os.Stderr, "Allow? (y/n/always): ")

    // Create a fresh reader for each prompt to avoid buffering issues
    reader := bufio.NewReader(os.Stdin)
    line, err := reader.ReadString('\n')
    if err != nil {
      errorChan <- fmt.Errorf("failed to read response: %w", err)
      return
    }

    response := strings.TrimSpace(strings.ToLower(line))

    switch response {
    case "y", "yes":
      fmt.Fprintln(os.Stderr, "✓ Granted temporarily")
      responseChan <- PromptResponse{Granted: true, Permanent: false}
    case "a", "always":
      fmt.Fprintln(os.Stderr, "✓ Granted permanently (this session)")
      responseChan <- PromptResponse{Granted: true, Permanent: true}
    default:
      fmt.Fprintln(os.Stderr, "✗ Permission denied")
      responseChan <- PromptResponse{Granted: false, Permanent: false}
    }
  }()

  select {
  case response := <-responseChan:
    return response, nil
  case err := <-errorChan:
    return PromptResponse{Granted: false}, err
  case <-ctx.Done():
    fmt.Fprintln(os.Stderr, "\n⏱️  Timeout - permission denied")
    return PromptResponse{Granted: false}, ctx.Err()
  }
}
