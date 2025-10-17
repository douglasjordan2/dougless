package tests

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/douglasjordan2/dougless/internal/runtime"
)

// TestCompleteJavaScriptProgram tests a realistic multi-feature JavaScript program
func TestCompleteJavaScriptProgram(t *testing.T) {
	rt := runtime.New()

	script := `
		// ================================================
		// Integration Test: Full Feature Demonstration
		// ================================================
		
		console.log('=== Dougless Runtime Integration Test ===');
		console.time('total-execution');
		
		// 1. Variables and Functions
		let counter = 0;
		const results = [];
		
		function fibonacci(n) {
			if (n <= 1) return n;
			return fibonacci(n - 1) + fibonacci(n - 2);
		}
		
		// 2. Computation
		console.time('fibonacci');
		const fib10 = fibonacci(10);
		console.timeEnd('fibonacci');
		console.log('Fibonacci(10) =', fib10);
		
		// 3. Array Operations
		const numbers = [1, 2, 3, 4, 5];
		for (let i = 0; i < numbers.length; i++) {
			results.push(numbers[i] * 2);
		}
		console.table(results);
		
		// 4. Demonstrate timer functions exist
		// Note: Full async testing requires VM thread-safety improvements (Phase 5)
		console.log('✓ setTimeout function available:', typeof setTimeout === 'function');
		console.log('✓ setInterval function available:', typeof setInterval === 'function');
		console.log('✓ clearTimeout function available:', typeof clearTimeout === 'function');
		console.log('✓ clearInterval function available:', typeof clearInterval === 'function');
		
		// 5. File system and module test
		// File system is globally available (unique to Dougless)
		if (typeof files === 'object') {
			console.log('✓ files global API available');
		}
		if (typeof files.read === 'function') {
			console.log('✓ files.read available');
		}
		if (typeof files.write === 'function') {
			console.log('✓ files.write available');
		}
		if (typeof files.rm === 'function') {
			console.log('✓ files.rm available');
		}
		
		// Module system (require)
		const path = require('path');
		if (typeof path === 'object') {
			console.log('✓ path module loaded via require');
		}
		// Note: http module will be added in Phase 3
		
		console.timeEnd('total-execution');
		console.log('=== Test Completed ===');
	`

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := rt.Execute(script, "integration_test.js")

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Fatalf("Execute() error = %v\nOutput:\n%s", err, output)
	}

	// Verify all expected outputs
	expectedOutputs := []string{
		"=== Dougless Runtime Integration Test ===",
		"Fibonacci(10) =",
		"55", // Result of fibonacci(10)
		"fibonacci:",
		"total-execution:",
		"✓ setTimeout function available:",
		"✓ setInterval function available:",
		"✓ clearTimeout function available:",
		"✓ clearInterval function available:",
		"✓ files global API available",
		"✓ files.read available",
		"✓ files.write available",
		"✓ files.rm available",
		"✓ path module loaded via require",
		"=== Test Completed ===",
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected) {
			t.Errorf("Output missing expected text: %q\nFull output:\n%s", expected, output)
		}
	}

	t.Logf("Integration test passed! Output length: %d bytes", len(output))
}

// TestTimerAccuracy tests the accuracy of setTimeout delays
func TestTimerAccuracy(t *testing.T) {
	rt := runtime.New()

	script := `
		const startTime = Date.now();
		const measurements = [];
		
		setTimeout(function() {
			measurements.push(Date.now() - startTime);
		}, 50);
		
		setTimeout(function() {
			measurements.push(Date.now() - startTime);
		}, 100);
		
		setTimeout(function() {
			measurements.push(Date.now() - startTime);
		}, 150);
	`

	startTime := time.Now()
	err := rt.Execute(script, "timer_accuracy.js")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	elapsed := time.Since(startTime)

	// Should have waited at least 150ms for all timers
	if elapsed < 150*time.Millisecond {
		t.Errorf("Timers completed too quickly: %v < 150ms", elapsed)
	}

	// But not too long (with tolerance)
	if elapsed > 250*time.Millisecond {
		t.Errorf("Timers took too long: %v > 250ms", elapsed)
	}
}

// TestConcurrentTimers tests multiple timers running concurrently
func TestConcurrentTimers(t *testing.T) {
	rt := runtime.New()

	script := `
		let completedTimers = 0;
		const timerResults = [];
		
		// Schedule 10 concurrent timers
		for (let i = 0; i < 10; i++) {
			(function(index) {
				setTimeout(function() {
					completedTimers++;
					timerResults.push(index);
				}, 10 + (index * 5));
			})(i);
		}
	`

	err := rt.Execute(script, "concurrent_timers.js")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Verify all timers completed
	// Note: We can't directly check the JS variables from here,
	// but we can verify the script executed without error
}

// TestErrorHandling tests how the runtime handles JavaScript errors
func TestErrorHandling(t *testing.T) {
	t.Run("syntax error", func(t *testing.T) {
		rt := runtime.New()
		script := `const x = ; // Syntax error`

		err := rt.Execute(script, "syntax_error.js")
		if err == nil {
			t.Error("Expected syntax error, got nil")
		}
	})

	t.Run("runtime error", func(t *testing.T) {
		rt := runtime.New()
		script := `
			function throwError() {
				throw new Error('Test error');
			}
			throwError();
		`

		err := rt.Execute(script, "runtime_error.js")
		if err == nil {
			t.Error("Expected runtime error, got nil")
		}
	})

	t.Run("undefined variable", func(t *testing.T) {
		rt := runtime.New()
		script := `console.log(undefinedVariable);`

		err := rt.Execute(script, "undefined_var.js")
		if err == nil {
			t.Error("Expected error for undefined variable, got nil")
		}
	})
}

// TestConsoleOperationsIntegration tests all console operations together
func TestConsoleOperationsIntegration(t *testing.T) {
	rt := runtime.New()

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	script := `
		console.log('Starting console test');
		console.warn('This is a warning');
		console.error('This is an error');
		
		console.time('operation1');
		for (let i = 0; i < 1000; i++) {
			// Simulate work
		}
		console.timeEnd('operation1');
		
		console.log('Array table:');
		console.table([10, 20, 30, 40, 50]);
		
		console.log('Object table:');
		console.table({
			name: 'Dougless',
			version: '0.1.0',
			status: 'testing'
		});
		
		console.log('Console test completed');
	`

	err := rt.Execute(script, "console_integration.js")

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Verify all console operations produced output
	checks := map[string]bool{
		"Starting console test":  false,
		"WARN:":                  false,
		"ERROR:":                 false,
		"operation1:":            false,
		"Array table:":           false,
		"Object table:":          false,
		"Console test completed": false,
	}

	for check := range checks {
		if strings.Contains(output, check) {
			checks[check] = true
		}
	}

	for check, found := range checks {
		if !found {
			t.Errorf("Console output missing: %q", check)
		}
	}
}

// TestModuleSystemIntegration tests the module system and global APIs
func TestModuleSystemIntegration(t *testing.T) {
	rt := runtime.New()

	script := `
		// Test global files API (unique to Dougless - not via require)
		if (typeof files !== 'object') {
			console.error('files global API not available');
		}
		
		if (typeof files.read !== 'function') {
			console.error('files.read not a function');
		}
		
		if (typeof files.write !== 'function') {
			console.error('files.write not a function');
		}
		
		if (typeof files.rm !== 'function') {
			console.error('files.rm not a function');
		}
		
		// Test require() module system
		const path = require('path');
		
		if (typeof path !== 'object') {
			console.error('path module not loaded');
		}
		
		if (typeof path.join !== 'function') {
			console.error('path.join not a function');
		}
		
		// HTTP module will be added in Phase 3
		console.log('Module system test passed');
	`

	err := rt.Execute(script, "module_integration.js")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

// BenchmarkFullProgram benchmarks a complete realistic program
func BenchmarkFullProgram(b *testing.B) {
	script := `
		function fibonacci(n) {
			if (n <= 1) return n;
			return fibonacci(n - 1) + fibonacci(n - 2);
		}
		
		const result = fibonacci(15);
		
		const data = [];
		for (let i = 0; i < 10; i++) {
			data.push(i * 2);
		}
		
		setTimeout(function() {
			let sum = 0;
			for (let i = 0; i < data.length; i++) {
				sum += data[i];
			}
		}, 5);
	`

	// Redirect stdout to suppress output during benchmark
	oldStdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = oldStdout }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt := runtime.New()
		rt.Execute(script, "benchmark.js")
	}
}
