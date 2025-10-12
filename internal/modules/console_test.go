package modules

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dop251/goja"
)

// Helper function to capture stdout
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestConsoleLog(t *testing.T) {
	vm := goja.New()
	console := NewConsole()
	vm.Set("console", console.Export(vm))

	tests := []struct {
		name     string
		script   string
		expected string
	}{
		{
			name:     "single string argument",
			script:   `console.log("hello")`,
			expected: "hello\n",
		},
		{
			name:     "multiple arguments",
			script:   `console.log("Count:", 42, "items")`,
			expected: "Count: 42 items\n",
		},
		{
			name:     "number argument",
			script:   `console.log(123)`,
			expected: "123\n",
		},
		{
			name:     "boolean argument",
			script:   `console.log(true)`,
			expected: "true\n",
		},
		{
			name:     "object argument",
			script:   `console.log({name: "test", value: 42})`,
			expected: "map[name:test value:42]\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				_, err := vm.RunString(tt.script)
				if err != nil {
					t.Fatalf("Script execution failed: %v", err)
				}
			})

			if output != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, output)
			}
		})
	}
}

func TestConsoleError(t *testing.T) {
	vm := goja.New()
	console := NewConsole()
	vm.Set("console", console.Export(vm))

	output := captureOutput(func() {
		_, err := vm.RunString(`console.error("Something went wrong")`)
		if err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}
	})

	expected := "ERROR: Something went wrong\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestConsoleWarn(t *testing.T) {
	vm := goja.New()
	console := NewConsole()
	vm.Set("console", console.Export(vm))

	output := captureOutput(func() {
		_, err := vm.RunString(`console.warn("This is a warning")`)
		if err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}
	})

	expected := "WARN: This is a warning\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestConsoleTime(t *testing.T) {
	vm := goja.New()
	console := NewConsole()
	vm.Set("console", console.Export(vm))

	t.Run("measure elapsed time", func(t *testing.T) {
		output := captureOutput(func() {
			_, err := vm.RunString(`
				console.time("test-timer");
				// Simulate some work
				var sum = 0;
				for (var i = 0; i < 1000; i++) {
					sum += i;
				}
				console.timeEnd("test-timer");
			`)
			if err != nil {
				t.Fatalf("Script execution failed: %v", err)
			}
		})

		// Should output something like "test-timer: 0.123ms"
		if !strings.HasPrefix(output, "test-timer:") {
			t.Errorf("Expected output to start with 'test-timer:', got %q", output)
		}
		if !strings.HasSuffix(strings.TrimSpace(output), "ms") {
			t.Errorf("Expected output to end with 'ms', got %q", output)
		}
	})

	t.Run("default label", func(t *testing.T) {
		output := captureOutput(func() {
			_, err := vm.RunString(`
				console.time();
				console.timeEnd();
			`)
			if err != nil {
				t.Fatalf("Script execution failed: %v", err)
			}
		})

		if !strings.HasPrefix(output, "default:") {
			t.Errorf("Expected output to start with 'default:', got %q", output)
		}
	})

	t.Run("timeEnd without time", func(t *testing.T) {
		output := captureOutput(func() {
			_, err := vm.RunString(`console.timeEnd("nonexistent")`)
			if err != nil {
				t.Fatalf("Script execution failed: %v", err)
			}
		})

		expected := "Warning: No such label 'nonexistent' for console.timeEnd()\n"
		if output != expected {
			t.Errorf("Expected %q, got %q", expected, output)
		}
	})

	t.Run("multiple timers concurrently", func(t *testing.T) {
		output := captureOutput(func() {
			_, err := vm.RunString(`
				console.time("timer1");
				console.time("timer2");
				console.timeEnd("timer1");
				console.timeEnd("timer2");
			`)
			if err != nil {
				t.Fatalf("Script execution failed: %v", err)
			}
		})

		lines := strings.Split(strings.TrimSpace(output), "\n")
		if len(lines) != 2 {
			t.Errorf("Expected 2 timer outputs, got %d", len(lines))
		}
		if !strings.HasPrefix(lines[0], "timer1:") {
			t.Errorf("Expected first line to start with 'timer1:', got %q", lines[0])
		}
		if !strings.HasPrefix(lines[1], "timer2:") {
			t.Errorf("Expected second line to start with 'timer2:', got %q", lines[1])
		}
	})
}

func TestConsoleTable(t *testing.T) {
	vm := goja.New()
	console := NewConsole()
	vm.Set("console", console.Export(vm))

	t.Run("array table", func(t *testing.T) {
		output := captureOutput(func() {
			_, err := vm.RunString(`console.table([1, 2, 3])`)
			if err != nil {
				t.Fatalf("Script execution failed: %v", err)
			}
		})

		// Should contain table characters and the values
		if !strings.Contains(output, "┌") || !strings.Contains(output, "└") {
			t.Errorf("Expected table borders in output, got %q", output)
		}
		if !strings.Contains(output, "(index)") {
			t.Errorf("Expected '(index)' header in output, got %q", output)
		}
		if !strings.Contains(output, "1") || !strings.Contains(output, "2") || !strings.Contains(output, "3") {
			t.Errorf("Expected values 1, 2, 3 in output, got %q", output)
		}
	})

	t.Run("object table", func(t *testing.T) {
		output := captureOutput(func() {
			_, err := vm.RunString(`console.table({name: "test", age: 42})`)
			if err != nil {
				t.Fatalf("Script execution failed: %v", err)
			}
		})

		// Should contain table characters and the keys/values
		if !strings.Contains(output, "┌") || !strings.Contains(output, "└") {
			t.Errorf("Expected table borders in output, got %q", output)
		}
		if !strings.Contains(output, "name") || !strings.Contains(output, "test") {
			t.Errorf("Expected 'name' and 'test' in output, got %q", output)
		}
		if !strings.Contains(output, "age") || !strings.Contains(output, "42") {
			t.Errorf("Expected 'age' and '42' in output, got %q", output)
		}
	})

	t.Run("empty array", func(t *testing.T) {
		output := captureOutput(func() {
			_, err := vm.RunString(`console.table([])`)
			if err != nil {
				t.Fatalf("Script execution failed: %v", err)
			}
		})

		// Empty array should produce no output
		if output != "" {
			t.Errorf("Expected empty output for empty array, got %q", output)
		}
	})

	t.Run("unsupported type fallback", func(t *testing.T) {
		output := captureOutput(func() {
			_, err := vm.RunString(`console.table(42)`)
			if err != nil {
				t.Fatalf("Script execution failed: %v", err)
			}
		})

		// Should fallback to regular print
		if !strings.Contains(output, "42") {
			t.Errorf("Expected '42' in output, got %q", output)
		}
	})
}

func TestConsoleThreadSafety(t *testing.T) {
	console := NewConsole()

	// Test that Console struct is thread-safe for time/timeEnd
	// Note: Goja VM itself is not thread-safe, so we only test the Console internals
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			// Create a new VM for each goroutine (Goja requirement)
			vm := goja.New()
			vm.Set("console", console.Export(vm))

			label := fmt.Sprintf("timer-%d", id)
			script := fmt.Sprintf(`
				console.time("%s");
				console.timeEnd("%s");
			`, label, label)

			captureOutput(func() {
				vm.RunString(script)
			})

			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Test timed out waiting for goroutines")
		}
	}
}

func TestRepeatChar(t *testing.T) {
	tests := []struct {
		name     string
		char     rune
		count    int
		expected string
	}{
		{
			name:     "single character",
			char:     'a',
			count:    1,
			expected: "a",
		},
		{
			name:     "multiple characters",
			char:     '-',
			count:    5,
			expected: "-----",
		},
		{
			name:     "zero count",
			char:     'x',
			count:    0,
			expected: "",
		},
		{
			name:     "unicode character",
			char:     '─',
			count:    3,
			expected: "───",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repeatChar(tt.char, tt.count)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
