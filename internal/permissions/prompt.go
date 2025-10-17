// Package permissions prompter handles interactive permission prompts.
// Displays permission requests to users and captures their responses.
package permissions

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
)

// PromptResponse represents the user's response to a permission prompt.
type PromptResponse struct {
	Granted      bool // Whether the permission was granted
	SaveToConfig bool // Whether to write to .douglessrc
}

// Prompter is the interface for prompting users for permissions.
// Implementations can provide different UI experiences (terminal, GUI, test mocks).
type Prompter interface {
	Prompt(ctx context.Context, desc PermissionDescriptor) (PromptResponse, error)
}

// StdioPrompter implements interactive terminal prompts via stdin/stderr.
// All prompts are serialized to prevent concurrent stdin reads.
type StdioPrompter struct {
	mu sync.Mutex // Serialize prompts to prevent concurrent stdin reads
}

// NewStdioPrompter creates a new terminal prompter.
func NewStdioPrompter() *StdioPrompter {
	return &StdioPrompter{}
}

// Prompt displays a permission request and waits for user input.
// Accepts responses: y/yes (grant), or any other (deny).
// If granted, prompts whether to save to .douglessrc.
// Respects context cancellation and timeouts.
//
// The prompt is displayed on stderr to avoid interfering with program output.
func (p *StdioPrompter) Prompt(ctx context.Context, desc PermissionDescriptor) (PromptResponse, error) {
	// Serialize prompts to prevent multiple concurrent stdin reads
	p.mu.Lock()
	defer p.mu.Unlock()

	responseChan := make(chan PromptResponse, 1)
	errorChan := make(chan error, 1)

	go func() {
		fmt.Fprintf(os.Stderr, "\n⚠️  Permission request: %s\n", desc)
		fmt.Fprintf(os.Stderr, "Allow? (y/n): ")

		// Create a fresh reader for each prompt to avoid buffering issues
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			errorChan <- fmt.Errorf("failed to read response: %w", err)
			return
		}

		response := strings.TrimSpace(strings.ToLower(line))

		// If denied, return early
		if response != "y" && response != "yes" {
			fmt.Fprintln(os.Stderr, "✗ Permission denied")
			responseChan <- PromptResponse{Granted: false, SaveToConfig: false}
			return
		}

		// User said yes - ask about config
		fmt.Fprintf(os.Stderr, "Save to .douglessrc? (y/n): ")

		line, err = reader.ReadString('\n')
		if err != nil {
			// Grant for session even if second read fails
			fmt.Fprintln(os.Stderr, "✓ Granted for this session")
			responseChan <- PromptResponse{Granted: true, SaveToConfig: false}
			return
		}

		saveResponse := strings.TrimSpace(strings.ToLower(line))
		saveToConfig := (saveResponse == "y" || saveResponse == "yes")

		if saveToConfig {
			fmt.Fprintln(os.Stderr, "✓ Granted and saved to .douglessrc")
		} else {
			fmt.Fprintln(os.Stderr, "✓ Granted for this session only")
		}

		responseChan <- PromptResponse{
			Granted:      true,
			SaveToConfig: saveToConfig,
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
