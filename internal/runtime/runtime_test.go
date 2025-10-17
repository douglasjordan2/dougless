package runtime

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestNew verifies runtime initialization
func TestNew(t *testing.T) {
	rt := New()

	if rt == nil {
		t.Fatal("New() returned nil")
	}

	if rt.vm == nil {
		t.Error("JavaScript VM should be initialized")
	}

	if rt.eventLoop == nil {
		t.Error("Event loop should be initialized")
	}

	if rt.modules == nil {
		t.Error("Module registry should be initialized")
	}

	if rt.timers == nil {
		t.Error("Timers map should be initialized")
	}
}

// TestExecuteBasicJavaScript tests basic JavaScript execution
func TestExecuteBasicJavaScript(t *testing.T) {
	tests := []struct {
		name   string
		script string
		want   error
	}{
		{
			name:   "empty script",
			script: "",
			want:   nil,
		},
		{
			name:   "simple variable",
			script: "var x = 42;",
			want:   nil,
		},
		{
			name:   "function declaration",
			script: "function add(a, b) { return a + b; }",
			want:   nil,
		},
		{
			name:   "expression",
			script: "1 + 1;",
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := New()
			err := rt.Execute(tt.script, "test.js")

			if (err != nil) != (tt.want != nil) {
				t.Errorf("Execute() error = %v, want %v", err, tt.want)
			}
		})
	}
}

// TestConsoleLog tests console.log functionality
func TestConsoleLog(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rt := New()
	script := `console.log("Hello, World!");`

	err := rt.Execute(script, "test.js")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "Hello, World!") {
		t.Errorf("console.log output = %q, want to contain %q", output, "Hello, World!")
	}
}

// TestConsoleError tests console.error functionality
func TestConsoleError(t *testing.T) {
	// Capture stdout (console.error also writes to stdout in our implementation)
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rt := New()
	script := `console.error("Error message");`

	err := rt.Execute(script, "test.js")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "ERROR:") {
		t.Errorf("console.error output = %q, want to contain %q", output, "ERROR:")
	}

	if !strings.Contains(output, "Error message") {
		t.Errorf("console.error output = %q, want to contain %q", output, "Error message")
	}
}

// TestConsoleWarn tests console.warn functionality
func TestConsoleWarn(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rt := New()
	script := `console.warn("Warning message");`

	err := rt.Execute(script, "test.js")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "WARN:") {
		t.Errorf("console.warn output = %q, want to contain %q", output, "WARN:")
	}

	if !strings.Contains(output, "Warning message") {
		t.Errorf("console.warn output = %q, want to contain %q", output, "Warning message")
	}
}

// TestConsoleTime tests console.time/timeEnd functionality
func TestConsoleTime(t *testing.T) {
	t.Run("basic timer", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		rt := New()
		script := `
			console.time('test');
			// Simulate some work
			for (let i = 0; i < 1000; i++) {}
			console.timeEnd('test');
		`

		err := rt.Execute(script, "test.js")
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		if !strings.Contains(output, "test:") {
			t.Errorf("console.time output = %q, want to contain %q", output, "test:")
		}

		if !strings.Contains(output, "ms") {
			t.Errorf("console.time output = %q, want to contain %q", output, "ms")
		}
	})

	t.Run("default label", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		rt := New()
		script := `
			console.time();
			console.timeEnd();
		`

		err := rt.Execute(script, "test.js")
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		if !strings.Contains(output, "default:") {
			t.Errorf("console.time output = %q, want to contain %q", output, "default:")
		}
	})

	t.Run("timeEnd without time", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		rt := New()
		script := `console.timeEnd('nonexistent');`

		err := rt.Execute(script, "test.js")
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		if !strings.Contains(output, "Warning") {
			t.Errorf("console.timeEnd output = %q, want to contain %q", output, "Warning")
		}
	})
}

// TestConsoleTable tests console.table functionality
func TestConsoleTable(t *testing.T) {
	t.Run("array data", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		rt := New()
		script := `console.table([1, 2, 3]);`

		err := rt.Execute(script, "test.js")
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		// Check for table borders (box drawing characters)
		if !strings.Contains(output, "┌") || !strings.Contains(output, "│") {
			t.Errorf("console.table output = %q, want to contain table formatting", output)
		}

		if !strings.Contains(output, "(index)") {
			t.Errorf("console.table output = %q, want to contain %q", output, "(index)")
		}
	})

	t.Run("object data", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		rt := New()
		script := `console.table({name: 'Doug', age: 30});`

		err := rt.Execute(script, "test.js")
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		if !strings.Contains(output, "┌") || !strings.Contains(output, "│") {
			t.Errorf("console.table output = %q, want to contain table formatting", output)
		}
	})
}

// TestSetTimeout tests setTimeout functionality
func TestSetTimeout(t *testing.T) {
	t.Run("basic timeout", func(t *testing.T) {
		rt := New()

		// Use a shared variable to track execution
		executed := false
		var mu sync.Mutex

		script := `
			let executed = false;
			setTimeout(function() {
				executed = true;
			}, 50);
		`

		err := rt.Execute(script, "test.js")
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		// Check that callback was executed
		result, err := rt.vm.RunString("executed")
		if err != nil {
			t.Fatalf("Failed to get result: %v", err)
		}

		mu.Lock()
		executed = result.ToBoolean()
		mu.Unlock()

		if !executed {
			t.Error("setTimeout callback should have executed")
		}
	})

	t.Run("timeout with delay", func(t *testing.T) {
		rt := New()
		startTime := time.Now()

		script := `
			let done = false;
			setTimeout(function() {
				done = true;
			}, 100);
		`

		err := rt.Execute(script, "test.js")
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		elapsed := time.Since(startTime)

		// Should have waited at least 100ms
		if elapsed < 100*time.Millisecond {
			t.Errorf("setTimeout executed too quickly: %v < 100ms", elapsed)
		}

		// But not too long (with 100ms tolerance)
		if elapsed > 250*time.Millisecond {
			t.Errorf("setTimeout executed too slowly: %v > 250ms", elapsed)
		}
	})
}

// TestSetInterval tests setInterval functionality
func TestSetInterval(t *testing.T) {
	rt := New()

	script := `
		let count = 0;
		const intervalId = setInterval(function() {
			count++;
			if (count >= 3) {
				clearInterval(intervalId);
			}
		}, 50);
	`

	err := rt.Execute(script, "test.js")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Verify the count increased
	result, err := rt.vm.RunString("count")
	if err != nil {
		t.Fatalf("Failed to get count: %v", err)
	}

	count := result.ToInteger()
	if count < 3 {
		t.Errorf("setInterval should have executed at least 3 times, got %d", count)
	}
}

// TestClearTimeout tests clearTimeout functionality
func TestClearTimeout(t *testing.T) {
	rt := New()

	script := `
		let executed = false;
		const timeoutId = setTimeout(function() {
			executed = true;
		}, 100);
		clearTimeout(timeoutId);
	`

	err := rt.Execute(script, "test.js")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Wait a bit to ensure it doesn't execute
	time.Sleep(150 * time.Millisecond)

	result, err := rt.vm.RunString("executed")
	if err != nil {
		t.Fatalf("Failed to get result: %v", err)
	}

	if result.ToBoolean() {
		t.Error("setTimeout callback should not have executed after clearTimeout")
	}
}

// TestModuleRequire tests the require() function
func TestModuleRequire(t *testing.T) {
	rt := New()

	t.Run("global files API available", func(t *testing.T) {
		script := `
			// File system is globally available (not via require)
			const filesAvailable = typeof files === 'object';
		`

		err := rt.Execute(script, "test.js")
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		// Verify the global files object is available
		result, err := rt.vm.RunString("filesAvailable")
		if err != nil {
			t.Fatalf("Failed to check files availability: %v", err)
		}

		if !result.ToBoolean() {
			t.Error("global 'files' object should be available")
		}
	})

	t.Run("files API has expected methods", func(t *testing.T) {
		script := `
			const hasRead = typeof files.read === 'function';
			const hasWrite = typeof files.write === 'function';
			const hasRm = typeof files.rm === 'function';
		`

		err := rt.Execute(script, "test.js")
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		// Check each method
		hasRead, _ := rt.vm.RunString("hasRead")
		hasWrite, _ := rt.vm.RunString("hasWrite")
		hasRm, _ := rt.vm.RunString("hasRm")

		if !hasRead.ToBoolean() {
			t.Error("files.read should be a function")
		}
		if !hasWrite.ToBoolean() {
			t.Error("files.write should be a function")
		}
		if !hasRm.ToBoolean() {
			t.Error("files.rm should be a function")
		}
	})

	t.Run("require path module", func(t *testing.T) {
		script := `
			const path = require('path');
		`

		err := rt.Execute(script, "test.js")
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		result, err := rt.vm.RunString("typeof path")
		if err != nil {
			t.Fatalf("Failed to check path type: %v", err)
		}

		if result.String() != "object" {
			t.Errorf("require('path') should return an object, got %s", result.String())
		}
	})
}

// TestComplexScript tests a more realistic JavaScript program
func TestComplexScript(t *testing.T) {
	rt := New()

	script := `
		// Test variables and functions
		const x = 10;
		const y = 20;
		
		function add(a, b) {
			return a + b;
		}
		
		const result = add(x, y);
		
		// Test console
		console.log('Result:', result);
		
		// Test timeout
		let asyncDone = false;
		setTimeout(function() {
			asyncDone = true;
			console.log('Async completed');
		}, 50);
	`

	err := rt.Execute(script, "test.js")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Verify result
	result, err := rt.vm.RunString("result")
	if err != nil {
		t.Fatalf("Failed to get result: %v", err)
	}

	if result.ToInteger() != 30 {
		t.Errorf("result = %d, want 30", result.ToInteger())
	}

	// Verify async completion
	asyncResult, err := rt.vm.RunString("asyncDone")
	if err != nil {
		t.Fatalf("Failed to get asyncDone: %v", err)
	}

	if !asyncResult.ToBoolean() {
		t.Error("async operation should have completed")
	}
}

// BenchmarkExecuteSimpleScript measures execution performance
func BenchmarkExecuteSimpleScript(b *testing.B) {
	script := "const x = 42; const y = x * 2;"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt := New()
		rt.Execute(script, "bench.js")
	}
}

// BenchmarkConsoleLog measures console.log performance
func BenchmarkConsoleLog(b *testing.B) {
	// Redirect stdout to /dev/null for benchmark
	oldStdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = oldStdout }()

	rt := New()
	script := `console.log('benchmark test');`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt.Execute(script, "bench.js")
	}
}

// BenchmarkSetTimeout measures setTimeout performance
func BenchmarkSetTimeout(b *testing.B) {
	script := `setTimeout(function() {}, 1);`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt := New()
		rt.Execute(script, "bench.js")
	}
}
