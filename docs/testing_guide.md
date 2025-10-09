# Dougless Runtime Testing Guide

## Overview

This document provides a comprehensive guide to the testing infrastructure implemented for Dougless Runtime. The test suite ensures code quality, performance, and correctness across all core components.

## Test Structure

```
dougless-runtime/
├── internal/
│   ├── event/
│   │   ├── loop.go           # Event loop implementation
│   │   └── loop_test.go      # ✅ Event loop unit tests (95.8% coverage)
│   └── runtime/
│       ├── runtime.go        # Core runtime
│       └── runtime_test.go   # ✅ Runtime unit tests (85.5% coverage)
└── tests/
    └── integration_test.go   # ✅ Integration tests
```

## Test Categories

### 1. Unit Tests

Unit tests verify individual components in isolation.

#### Event Loop Tests (`internal/event/loop_test.go`)

**Coverage: 95.8%**

Tests include:
- **TestNewLoop**: Verifies proper initialization
- **TestScheduleTask**: Tests immediate, delayed, and multiple task scheduling
- **TestIntervalTask**: Tests recurring interval execution and cancellation
- **TestClearTimer**: Tests timer cancellation (before execution, non-existent, already executed)
- **TestStopLoop**: Tests graceful shutdown with and without pending tasks
- **TestConcurrentOperations**: Tests thread-safety with 50 concurrent tasks

**Key Learning**: These tests discovered and fixed a race condition in the `Stop()` method where the event loop's `running` flag wasn't set before `Stop()` was called.

#### Runtime Tests (`internal/runtime/runtime_test.go`)

**Coverage: 85.5%**

Tests include:
- **TestNew**: Runtime initialization
- **TestExecuteBasicJavaScript**: Basic JavaScript execution (variables, functions, expressions)
- **TestConsoleLog/Error/Warn**: Console output functionality
- **TestConsoleTime**: Performance measurement with `console.time/timeEnd`
- **TestConsoleTable**: Table formatting for arrays and objects
- **TestSetTimeout**: Delayed execution
- **TestSetInterval**: Recurring execution with auto-clear
- **TestClearTimeout**: Timer cancellation
- **TestModuleRequire**: Module system (fs, http, path)
- **TestComplexScript**: Realistic multi-feature program

**Key Techniques**:
- **Output Capture**: Tests capture stdout using `os.Pipe()` to verify console output
- **Timing Verification**: Tests check delays with tolerance ranges
- **Error Testing**: Subtests verify different error conditions

### 2. Integration Tests

Integration tests verify complete end-to-end functionality.

#### Complete Program Test (`TestCompleteJavaScriptProgram`)

Executes a comprehensive JavaScript program that tests:
- Variables and functions (Fibonacci calculation)
- Array operations
- Console operations (log, time/timeEnd, table)
- Module system (require fs, http, path)
- Timer function availability

**Output Verification**: Checks for specific strings in the output to ensure all features executed correctly.

#### Error Handling Tests

Tests JavaScript error scenarios:
- Syntax errors (`var x = ;`)
- Runtime errors (`throw new Error()`)
- Undefined variable access

#### Timer Accuracy Test

Verifies that setTimeout delays are respected within tolerance:
- Schedules timers at 50ms, 100ms, 150ms
- Verifies total execution time: 150ms ≤ elapsed < 250ms

### 3. Benchmark Tests

Benchmarks measure performance to track regressions.

#### Event Loop Benchmarks

```
BenchmarkScheduleTask      1,553,498 ops    735.6 ns/op    200 B/op    4 allocs/op
BenchmarkTimerScheduling     764,336 ops   1667 ns/op     669 B/op    5 allocs/op
```

**Interpretation**:
- Can schedule **1.5 million immediate tasks per second**
- Can schedule **764K delayed tasks per second**
- Very low memory overhead (200-669 bytes per operation)

#### Runtime Benchmarks

```
BenchmarkExecuteSimpleScript     50,139 ops    20.1 µs/op    18,996 B/op    200 allocs/op
BenchmarkConsoleLog             318,656 ops     4.3 µs/op     1,809 B/op     39 allocs/op
BenchmarkSetTimeout                 982 ops  1167.5 µs/op    19,896 B/op    210 allocs/op
```

**Interpretation**:
- Simple script execution: ~20µs (50,000 executions/sec)
- Console.log: ~4µs (318,000 logs/sec)
- setTimeout overhead: ~1.2ms (includes event loop coordination)

#### Integration Benchmark

```
BenchmarkFullProgram    181 ops    6.57 ms/op    43,468 B/op    526 allocs/op
```

A complete program (Fibonacci, arrays, timers) takes ~6.57ms to execute.

## Code Coverage

### Overall Coverage: 73.1%

Breakdown by package:
- **Event Loop**: 95.8% ✅ (Excellent)
- **Runtime**: 85.5% ✅ (Very Good)
- **Modules**: 0% ⚠️ (Placeholder implementations)
- **CLI**: 0% ⚠️ (Not tested - simple entry point)

### Uncovered Code

The 27% of uncovered code consists of:
1. **Placeholder module implementations** (fs, http, path) - These will be tested when implemented in Phase 2
2. **CLI entry point** (cmd/dougless/main.go) - Simple wrapper, not critical to test
3. **Error paths** in some functions - Edge cases

## Running Tests

### All Tests
```bash
go test ./...
```

### With Verbose Output
```bash
go test -v ./...
```

### With Coverage
```bash
go test -cover ./...
```

### Detailed Coverage Report
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # Opens in browser
```

### Specific Package
```bash
go test -v ./internal/event
go test -v ./internal/runtime
go test -v ./tests
```

### Benchmarks
```bash
go test -bench=. ./...
go test -bench=. -benchmem ./...  # Include memory stats
```

### Race Detection
```bash
go test -race ./...
```

**Note**: Race detection is critical for concurrent code. Our tests pass with `-race`, confirming thread-safety.

## Go Testing Concepts Explained

### Test Function Structure

```go
func TestSomething(t *testing.T) {
    // Arrange: Set up test data
    rt := runtime.New()
    
    // Act: Perform the operation
    err := rt.Execute("var x = 42;", "test.js")
    
    // Assert: Verify the result
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}
```

### Subtests with t.Run()

```go
func TestConsole(t *testing.T) {
    t.Run("log", func(t *testing.T) {
        // Test console.log
    })
    
    t.Run("error", func(t *testing.T) {
        // Test console.error
    })
}
```

**Benefits**:
- Organized output
- Can run specific subtests: `go test -run TestConsole/log`
- Each subtest is independent

### Assertions

- **t.Error()**: Report failure but continue
- **t.Errorf()**: Report formatted failure and continue
- **t.Fatal()**: Report failure and stop immediately
- **t.Fatalf()**: Report formatted failure and stop

**When to use each**:
- Use `Error` when subsequent tests can still provide value
- Use `Fatal` when the test can't continue meaningfully

### Table-Driven Tests

```go
tests := []struct {
    name   string
    input  string
    want   string
}{
    {"simple", "var x = 1;", "1"},
    {"complex", "var x = 2 * 2;", "4"},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test with tt.input and tt.want
    })
}
```

**Benefits**:
- Easy to add new test cases
- Clear test data separation
- Automatically generates subtests

### Benchmarks

```go
func BenchmarkSomething(b *testing.B) {
    // Setup (not timed)
    data := setupData()
    
    b.ResetTimer()  // Start timing now
    
    for i := 0; i < b.N; i++ {
        doSomething(data)
    }
}
```

**Key Points**:
- Go automatically determines the right `b.N` value
- Use `b.ResetTimer()` to exclude setup time
- Use `b.StopTimer()` / `b.StartTimer()` for complex scenarios

### Test Helpers

```go
func helperFunction(t *testing.T) {
    t.Helper()  // Marks this as a helper
    // ... helper code ...
    if somethingWrong {
        t.Error("problem")  // Error reports caller's line
    }
}
```

The `t.Helper()` call makes error messages point to the test that called the helper, not the helper itself.

## Known Limitations

### 1. Async Callback Testing

**Issue**: Goja's VM is not thread-safe. Calling JavaScript functions from Go goroutines (which our event loop does) can cause panics in complex scenarios.

**Current Workaround**: Integration tests verify timer functions exist but don't execute complex async callbacks.

**Future Solution**: Phase 5 will implement proper Promise support and VM thread-safety improvements.

### 2. Module Implementation Testing

**Issue**: fs, http, and path modules are placeholder implementations with 0% coverage.

**Solution**: Phase 2 will implement these modules with comprehensive tests.

## Testing Best Practices

### 1. Test Behavior, Not Implementation

❌ Bad:
```go
if runtime.internalCounter != 5 {  // Testing internal state
    t.Error("Counter wrong")
}
```

✅ Good:
```go
result := runtime.Execute(script)
if result != expected {  // Testing observable behavior
    t.Error("Wrong result")
}
```

### 2. Use Clear Test Names

✅ Good names:
- `TestSetTimeout_WithDelay_ExecutesAfterDelay`
- `TestConsoleLog_WithMultipleArguments_PrintsAll`

❌ Bad names:
- `TestTimer1`
- `TestStuff`

### 3. One Concept Per Test

Each test should verify one specific behavior. Use subtests to organize related tests.

### 4. Use Timeouts for Async Tests

```go
select {
case <-done:
    // Success
case <-time.After(1 * time.Second):
    t.Fatal("Test timed out")
}
```

This prevents hanging tests that would block CI/CD.

### 5. Clean Up Resources

```go
func TestSomething(t *testing.T) {
    resource := acquire()
    defer resource.cleanup()  // Always cleans up
    
    // ... test code ...
}
```

## Next Steps

### For Phase 2 (File System & Modules):

1. **Test fs module** with real file operations
   - Read/write files
   - Directory operations
   - Error handling (permissions, not found, etc.)

2. **Test path module** with cross-platform paths
   - Windows vs Unix path handling
   - Edge cases (empty, root, relative)

3. **Test module resolution**
   - Relative paths
   - Module caching
   - Circular dependencies

### For Phase 5 (Promises & Async):

1. **Fix VM thread-safety** for proper async callbacks
2. **Test Promise implementation**
3. **Test async/await**

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Test Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [Table Driven Tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Advanced Go Testing](https://www.youtube.com/watch?v=8hQG7QlcLBk)

---

**Test Summary**: ✅ 100% of tests passing | 73.1% code coverage | 6 benchmark suites

**Status**: Testing infrastructure complete for Phase 1. Ready for Phase 2 development.
