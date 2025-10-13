# Changelog

All notable changes to Dougless Runtime will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Phase 1 - Foundation (COMPLETE)

#### Added - October 2024
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

- **REPL (Interactive Shell)**
  - Interactive JavaScript evaluation mode
  - Multi-line input support with automatic bracket/brace detection
  - State preservation between commands
  - Special commands: `.help`, `.exit`, `.quit`, `.clear`
  - Dual-mode CLI: REPL when no arguments, file execution when file provided
  - Proper error display with Goja exception handling

- **Core Infrastructure**
  - Goja JavaScript engine integration (ES5.1 support)
  - CLI tool for executing JavaScript files
  - `Runtime.Evaluate()` method for REPL support
  - Basic error handling and reporting

- **Examples and Tests**
  - `examples/interval_test.js` - Timer system demonstration
  - `examples/timer_edge_cases.js` - Edge case testing for timers
  - `examples/console_test.js` - Console enhancements demonstration
  - `examples/test-features.js` - Comprehensive feature testing
  - Unit tests for core components
  - Integration tests for JavaScript execution

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

## [Phase 2] - October 2024 - COMPLETE ✅

### Path Module
- **Full Implementation**
  - `path.join()` - Join path segments with OS-specific separator
  - `path.resolve()` - Resolve paths to absolute paths
  - `path.dirname()` - Extract directory name from path
  - `path.basename()` - Extract filename with optional extension removal
  - `path.extname()` - Get file extension
  - `path.sep` - OS-specific path separator constant
  - Cross-platform compatibility (Windows/Unix)

### File System Module (Global API)
- **Unique Dougless Feature**: `file` object available globally (no require needed)
- **Async Operations** (callback-based)
  - `file.read(path, callback)` - Read file contents
  - `file.write(path, data, callback)` - Write data to file
  - `file.readdir(path, callback)` - List directory contents
  - `file.exists(path, callback)` - Check if file/directory exists
  - `file.mkdir(path, callback)` - Create directory
  - `file.rmdir(path, callback)` - Remove empty directory
  - `file.unlink(path, callback)` - Delete file
  - `file.stat(path, callback)` - Get file/directory information (size, timestamps, type)
- **Event Loop Integration** - All operations properly scheduled on event loop
- **Simplified API** - Shorter method names (`read` vs `readFile`)

### Architecture Improvements
- Modules can access event loop for async operations
- Clean separation between global APIs and require-based modules
- FileSystem module pattern for async callback handling

### Examples Added
- `examples/test-path.js` - Path module demonstration
- `examples/test-file.js` - Basic file operations
- `examples/test-file-advanced.js` - Complete file system workflow

### Testing & Quality Assurance
- **Full Test Coverage** for file system module (4/4 tests passing)
- **Unit Tests** for path module operations
- **Integration Tests** updated to reflect global `file` API
- **Test Suite Status**: 25/25 tests passing (~75% code coverage)
- Fixed import path issues in test files
- Updated test expectations to match global `file` API design
- All Phase 1 & 2 functionality fully tested and validated

## [Phase 3] - October 2024 - COMPLETE ✅

### HTTP Module (Global API)
- **Unique Dougless Feature**: `http` object available globally (no require needed)
- **HTTP Client**
  - `http.get(url, callback)` - Make HTTP GET requests
  - `http.post(url, payload, callback)` - Make HTTP POST requests with JSON payload
  - Custom content-type support
  - Async operations integrated with event loop
  - Response object with status, statusCode, body, and headers
  
- **HTTP Server**
  - `http.createServer(requestHandler)` - Create HTTP server
  - `server.listen(port, callback)` - Start server on specified port
  - Request object with method, url, headers, and body
  - Response object with:
    - `setHeader(name, value)` - Set response headers
    - `statusCode` property - Set HTTP status code
    - `end(data)` - Send response and close connection
  - Multiple header values support (arrays for headers like Set-Cookie)
  - Background goroutine execution for non-blocking server operation
  - Error logging to stderr for server issues

### Architecture Improvements
- HTTP module integrated with event loop for async operations
- Clean separation between client and server functionality
- Proper error handling with Go's net/http package
- Global API design consistent with `file` module philosophy

### Examples Added
- `examples/http_server.js` - Full HTTP server with client requests
- `examples/simple_server.js` - Simple HTTP server with keep-alive

### Testing & Validation
- Manual testing with curl validated GET/POST requests
- Server successfully handles concurrent requests
- Custom headers properly set and received
- Request body parsing working correctly
- Multiple header values handled correctly

## [Phase 4] - October 2024 - IN PROGRESS 🚧

### Security & Permissions System
- **Context-Aware Permission Management** (Deno-inspired)
  - Runtime permission checks for file, network, and environment access
  - User prompts for permission granting in interactive mode
  - CLI flags for explicit permission grants:
    - `--allow-read[=path]` - Grant read access (optionally to specific paths)
    - `--allow-write[=path]` - Grant write access (optionally to specific paths)
    - `--allow-net[=host]` - Grant network access (optionally to specific hosts)
    - `--allow-env[=var]` - Grant environment variable access
    - `--allow-run[=program]` - Grant subprocess execution access
    - `--allow-all` - Grant all permissions (for development)
  - Interactive permission prompts with options:
    - `y/yes` - Grant temporarily (one-time)
    - `n/no` - Deny access
    - `a/always` - Grant permanently for session
  - Permission caching for repeated access
  - Automatic terminal detection for prompt mode
  
- **Path-Based Permissions**
  - Hierarchical path matching (parent directory grants access to children)
  - Prevents path traversal escapes (../ attacks)
  - Absolute path resolution and normalization
  
- **Network Permissions**
  - Host and port-based granular control
  - Wildcard domain support (`*.example.com`)
  - Localhost/loopback normalization (127.0.0.1, ::1, localhost)
  - Port defaulting for standard HTTP/HTTPS
  
- **Permission Error Messages**
  - Clear, actionable error messages
  - Suggested command-line flags for fixing permission issues
  - Examples for both specific and broad permission grants

### Architecture Improvements
- Global permission manager with thread-safe operations
- CLI flag parsing with support for comma-separated values
- Permission descriptor system for structured permission requests
- Mutex-protected prompt handling for concurrent operations
- Context-based timeout support for permission prompts

### Examples Added
- `examples/test_permissions.js` - Concurrent permission testing
- `examples/test_permissions_sequential.js` - Sequential interactive testing

### Documentation
- Moved roadmap from README.md to dedicated ROADMAP.md
- Updated WARP.md with permission system details
- Comprehensive permission documentation in README

### Bug Fixes
- Fixed stdin buffering issues in permission prompts
- Create fresh bufio.Reader for each prompt to avoid state issues
- Improved error handling in permission checks

## Upcoming

### Planned Features
- **WebSocket Support**
  - WebSocket client
  - WebSocket server
  - Real-time bidirectional communication
  - Connection management and broadcasting

## Future Phases

### Package Manager (Post Phase 4)
- **Package Management System** (npm/bun-style)
  - `dougless install <package>` - Install packages from npm registry
  - `dougless install` - Install dependencies from `dougless.json`
  - `dougless uninstall <package>` - Remove packages
  - `dougless update` - Update package versions
  - Package manifest (`dougless.json`) for dependency tracking
  - Lock file (`dougless-lock.json`) for reproducible builds
  - Semantic versioning support (`^`, `~`, `>=`, etc.)
  - Dependency resolution with conflict handling
  - Local module cache (`~/.dougless/cache/`)
  - `dougless_modules/` directory structure
  - Enhanced `require()` for external package resolution
  - Integration with npm registry API
