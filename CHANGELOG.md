# Changelog

All notable changes to Dougless Runtime will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Phase 1 - Foundation (In Progress)

#### Added - December 2024
- **Event Loop Implementation**
  - Complete async task scheduling using Go channels and goroutines
  - Timer management system with proper lifecycle handling
  - WaitGroup-based synchronization for graceful shutdown
  - Context-based cancellation support

- **Timer Functions**
  - `setTimeout(callback, delay)` - Schedule one-time delayed execution
  - `setInterval(callback, delay)` - Schedule recurring execution at intervals
  - `clearTimeout(timerID)` - Cancel pending timeout
  - `clearInterval(timerID)` - Stop active interval
  - UUID-based timer ID tracking for reliable cancellation
  - Proper cleanup preventing memory leaks

- **Enhanced Console Operations**
  - `console.log(...)` - Standard output logging
  - `console.error(...)` - Error message logging
  - `console.warn(...)` - Warning message logging
  - `console.time(label)` - Start performance timer with optional label
  - `console.timeEnd(label)` - End timer and display duration in milliseconds
  - `console.table(data)` - Display arrays and objects in formatted table
  - Thread-safe timer tracking with mutex protection

- **Module System**
  - CommonJS-style `require()` function
  - Module registry for built-in modules
  - Placeholder implementations for `fs`, `http`, and `path` modules

- **Core Infrastructure**
  - Goja JavaScript engine integration (ES5.1 support)
  - CLI tool for executing JavaScript files
  - Basic error handling and reporting

- **Examples and Tests**
  - `examples/interval_test.js` - Timer system demonstration
  - `examples/timer_edge_cases.js` - Edge case testing for timers
  - `examples/console_test.js` - Console enhancements demonstration

#### Documentation
- Comprehensive README with feature status
- Detailed project plan with 8-phase roadmap
- WARP.md for AI-assisted development guidance
- Transpilation strategy document for ES6+ support planning

### Technical Details

**Event Loop Architecture:**
- Single-threaded event loop with Go goroutine pool
- Non-blocking I/O operations
- Task queue with buffered channel (100 capacity)
- Timer queue using Go's time.AfterFunc

**Performance Characteristics:**
- Startup time: ~150ms for basic scripts (target: <100ms)
- Memory usage: Minimal for current feature set
- Timer accuracy: Sub-millisecond precision
- Graceful shutdown with proper resource cleanup

## [0.0.1] - Initial Commit

### Added
- Basic project structure
- Go module initialization
- Initial Goja integration
- Placeholder runtime implementation

---

## Upcoming (Phase 2)

### Planned Features
- **File System Module (`fs`)**
  - Synchronous operations: readFileSync, writeFileSync, existsSync
  - Asynchronous operations: readFile, writeFile with callbacks
  
- **Path Module (`path`)**
  - path.join, path.resolve, path.dirname, path.basename
  - Cross-platform path handling
  
- **Error Handling**
  - Stack trace preservation
  - Uncaught exception handling
  
- **Testing Infrastructure**
  - Unit tests for core components
  - Integration tests for JavaScript execution
  - Benchmark suite for performance tracking
