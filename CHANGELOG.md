# Changelog

All notable changes to Dougless Runtime will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Recent Updates - October 18, 2024 (Latest)

#### Completed - Phase 7 (Crypto & Security) ✅
- **Global Crypto API**
  - `crypto.createHash(algorithm)` - Hash functions (MD5, SHA-1, SHA-256, SHA-512)
  - `crypto.createHmac(algorithm, key)` - HMAC for message authentication
  - `crypto.random(size, [encoding])` - Cryptographically secure random bytes
  - `crypto.randomBytes(size, [encoding])` - Alias for random()
  - `crypto.uuid()` - Generate UUID v4
  - `crypto.timingSafeEqual(a, b)` - Timing-safe equality comparison
  - Multiple encoding support: hex, base64, raw
  - Node.js-compatible API design
  - Global access (no require() needed)
  
- **Hash Algorithms**
  - MD5 - For checksums (not recommended for security)
  - SHA-1 - Legacy support
  - SHA-256 - Recommended for most use cases
  - SHA-512 - Maximum security with larger output
  
- **Security Features**
  - Timing-safe comparison prevents timing attacks
  - HMAC for webhook signature verification
  - Cryptographically secure random number generation
  - UUID v4 generation for unique identifiers
  
- **Implementation Details**
  - `internal/modules/crypto.go` - Complete crypto module
  - Uses Go's crypto standard library (crypto/md5, crypto/sha256, etc.)
  - 34 comprehensive unit tests with 100% pass rate
  - Real-world examples: password hashing, API signing, webhook verification
  
- **New Files**
  - `internal/modules/crypto.go` - Crypto module implementation
  - `internal/modules/crypto_test.go` - Full test suite with 34 tests
  - `docs/crypto_api.md` - Complete API documentation with examples
  - `examples/crypto_demo.js` - Comprehensive usage examples
  
- **Documentation**
  - Complete API reference with all methods
  - Real-world examples for common use cases
  - Security best practices and guidelines
  - Node.js compatibility notes

### Previous Updates - October 17, 2024

#### Completed - Config-Based Permissions System ✅
- **`.douglessrc` Configuration File**
  - Project-centric permission model using JSON config files
  - Replaces CLI flag-based permission approach
  - Simple, readable JSON format for defining permissions
  - Per-project permission configuration
  - Config discovery starting from script directory
  
- **Permission Types Supported**
  - `read` - File system read access paths
  - `write` - File system write access paths
  - `net` - Network access hosts
  - `env` - Environment variable access
  - `run` - Subprocess execution permissions
  
- **Configuration Format**
  ```json
  {
    "permissions": {
      "read": ["./examples", "/tmp"],
      "write": ["/tmp"],
      "net": ["api.example.com"],
      "env": ["API_KEY"],
      "run": ["git"]
    }
  }
  ```
  
- **Implementation Details**
  - `internal/permissions/config.go` - Config loading and parsing
  - Automatic config file discovery from script directory
  - Clean JSON structure for easy editing
  - Integration with existing permission manager
  - Error handling for missing or malformed configs
  
- **Interactive Two-Prompt Save-to-Config ✅**
  - Two-step developer workflow when prompting in terminal:
    1) `Allow? (y/n)` → grant or deny for current operation
    2) If granted: `Save to .douglessrc? (y/n)` → persist permission to config
  - Saves via `SavePermissionToConfig()` in `internal/permissions/config.go`
  - Wired into `Manager.CheckWithPrompt()` in `internal/permissions/permissions.go`
  - Runtime sets config path for saves in `internal/runtime/runtime.go`
  - Non-blocking: save failures log a warning and do not affect the grant
  
- **Benefits Over CLI Flags**
  - Cleaner, more maintainable permission definitions
  - Version-controlled permissions with project code
  - No need to remember complex CLI flag combinations
  - Easier to share and document project permissions
  - Project-centric configuration approach

### Previous Updates - October 17, 2024

#### Completed - Phase 6 (WebSockets & Real-time) ✅
- **WebSocket Server Implementation**
  - `server.websocket(path, callbacks)` - Add WebSocket endpoint to HTTP server
  - Real-time bidirectional communication
  - Connection state management (readyState: CONNECTING, OPEN, CLOSING, CLOSED)
  - Browser-compatible API matching WebSocket specification
  - Thread-safe message sending with mutex protection
  - Event callbacks: `open`, `message`, `close`, `error`
  - Broadcasting to multiple connected clients
  - Proper connection lifecycle management
  - `ws.send(data)` - Send messages to client
  - `ws.close()` - Gracefully close connection
  
- **Implementation Details**
  - Built on gorilla/websocket library
  - Integration with existing HTTP server
  - Automatic upgrade from HTTP to WebSocket protocol
  - Concurrent connection handling with goroutines
  - Clean separation between HTTP and WebSocket handlers
  - Error handling for connection failures and invalid frames
  
- **New Example Files**
  - `examples/websocket_simple.js` - Basic WebSocket echo server
  - `examples/websocket_server.js` - Full-featured WebSocket server
  - `examples/websocket_chat.js` - Multi-client chat application with broadcasting
  
- **Architecture Improvements**
  - WebSocket module integrated with HTTP module
  - Connection objects with proper state machine
  - Mutex-protected write operations for thread safety
  - Event-driven callback system matching browser API

### Recent Updates - October 16, 2024

#### Added - File System Promise Support 🎉
- **Promise API for File Operations**
  - All file methods now support both callbacks AND promises
  - When callback is omitted, methods return a Promise
  - Full async/await support without wrapping
  - Backward compatible - existing callback code unchanged
  
- **Updated Methods**
  - `files.read(path)` - Returns `Promise<string | string[] | null>`
  - `files.write(path, content)` - Returns `Promise<void>`
  - `files.write(path)` - Returns `Promise<void>` (for directories)
  - `files.rm(path)` - Returns `Promise<void>`
  
- **Usage Examples**:
  ```javascript
  // Callback style (still works)
  files.read('data.txt', (err, content) => {
    if (!err) console.log(content);
  });
  
  // Promise style
  files.read('data.txt')
    .then(content => console.log(content))
    .catch(err => console.error(err));
  
  // Async/await style (cleanest!)
  const content = await files.read('data.txt');
  await files.write('output.txt', content);
  await files.rm('temp.txt');
  ```

- **Implementation Details**
  - Promise creation integrated with existing event loop
  - Promises reject with error strings, resolve with data
  - No breaking changes to existing API
  - Documentation updated across all files

- **Improved UX: Optional Content Parameter**
  - `files.write(path)` now creates empty files (like `touch` command)
  - Content parameter is fully optional - defaults to empty string
  - **Consistent API**: Both files and directories can be created with just a path
  - **Touch-like behavior**: Create empty files or truncate existing ones
  - **Cleaner code**: `files.write('file.txt')` instead of `files.write('file.txt', '')`
  
- **New Example Files**
  - `examples/files_promise.js` - Comprehensive demonstration of promise-based file operations
    - Shows callback, promise, and async/await patterns side by side
    - Includes parallel operations with Promise.all()
    - 10+ practical examples including error handling
  - `examples/files_touch.js` - NEW! Demonstrates empty file creation
    - Touch-like behavior for creating/truncating files
    - Parallel file creation with Promise.all()
    - API comparison showing improved ergonomics

### Recent Updates - October 15, 2024

#### Changed - File System API Simplification ⚡
- **Unified `files` API** - Simplified from 8 methods to 3 smart methods
  - **Breaking Change**: Renamed global from `file` → `files`
  - Convention-based routing using path patterns
  - 62% reduction in API surface while maintaining full functionality
  
- **New `files.read(path, callback)` Method**
  - **Trailing `/`**: Reads directory → returns `string[]` of names
  - **No trailing `/`**: Reads file → returns `string` content
  - **File doesn't exist**: Returns `null` (doubles as exists check)
  - Replaces: `file.read()`, `file.readdir()`, `file.exists()`
  
- **New `files.write(path, [content], callback)` Method**
  - **2 args** (path with `/`): Creates directory recursively
  - **3 args** (path + content): Writes file (auto-creates parent dirs)
  - Smart detection based on path conventions
  - Replaces: `file.write()`, `file.mkdir()`
  
- **New `files.rm(path, callback)` Method**
  - **Unified removal**: Deletes files OR directories recursively
  - Uses `os.RemoveAll()` under the hood
  - Handles non-existent paths gracefully
  - Replaces: `file.unlink()`, `file.rmdir()`
  
- **Removed Methods**: `file.stat()` (may return in future if needed)

- **Migration Path**:
  ```javascript
  // OLD → NEW
  file.read(path, cb)           → files.read(path, cb)
  file.write(path, data, cb)    → files.write(path, data, cb)
  file.readdir(path, cb)        → files.read(path + '/', cb)
  file.mkdir(path, cb)          → files.write(path + '/', cb)
  file.rmdir(path, cb)          → files.rm(path, cb)
  file.unlink(path, cb)         → files.rm(path, cb)
  file.exists(path, cb)         → files.read(path, cb) // null check
  ```

### Recent Updates - October 15, 2024 (Continued)

#### Cleanup & Code Quality
- **Test Suite Consolidation**
  - Merged duplicate test files into unified test suites
  - Removed temporary debug scripts and redundant test files
  - Deleted 7 standalone test files (test-debug.js, test-order.js, etc.)
  - Cleaned up examples/ directory structure
  - Renamed `path_module.js` to `path_examples.js` for consistency
  - Created comprehensive `sourcemap_examples.js` example
  
- **Documentation Cleanup**
  - Removed 3 temporary documentation files (768 lines total):
    - DOCUMENTATION_STATUS.md (comprehensive audit completed)
    - SOURCEMAP_EXPLANATION.md (integrated into code examples)
    - TODO_EFFORT_ANALYSIS.md (effort estimation no longer needed)
  - Updated TODO.md to reflect completed tasks
  - Consolidated documentation into permanent locations
  
- **Code Formatting & LSP**
  - Fixed duplicate `main` function declarations in demos/
  - Added build tags (`//go:build ignore`) to demo files
  - Ran `go fmt ./...` on entire codebase
  - Verified `go vet ./...` passes with zero warnings
  - Verified `go build ./...` succeeds
  - All LSP errors resolved
  
- **Statistics**
  - Total lines removed: 3,752 (mostly redundant test code and docs)
  - Total lines modified/added: 2,971 (formatting and consolidation)
  - Net reduction: ~780 lines of cleaner, more maintainable code
  - Zero skipped tests, zero TODO comments, zero LSP errors

### Recent Updates - October 15, 2024 (Earlier)

#### Completed - Phase 5 (Promises & ES6+) ✅
- **Promise.allSettled() Implementation**
  - ES2020-compliant `Promise.allSettled()` static method
  - Always resolves (never rejects), even when all promises reject
  - Returns array of result objects with `status` and `value`/`reason` properties
  - Handles mixed fulfilled and rejected promises
  - Properly handles non-promise values (wraps as fulfilled)
  - Waits for all promises to settle regardless of outcome
  - Full test coverage with 7 comprehensive tests
  - Example file: `examples/test-promise-allsettled.js`

- **Promise.any() Validation**
  - Already implemented and working (ES2021 feature)
  - Returns first fulfilled promise, ignoring rejections
  - Throws AggregateError when all promises reject
  - Properly handles non-promise values
  - Example file: `examples/test-promise-any.js`
  - **Phase 5 is now 100% COMPLETE** ✅

#### Fixed
- **Promise Reuse Bug**
  - Fixed panic when reusing settled promise instances
  - Test examples updated to create fresh promise arrays
  - Error: "slice bounds out of range [:-1]" resolved

### Previous Updates - October 14, 2024

#### Enhanced
- **Source Map Support**
  - Enabled inline source maps for transpiled code
  - Error messages now reference original source code line numbers
  - Automatic mapping from transpiled ES5 to original ES6+ code
  - Edge case handling: disabled for empty scripts to prevent parsing errors
  - Improved debugging experience with accurate stack traces
  - Smart variable name preservation in error messages
  - Benefits all ES6+ features: arrow functions, template literals, destructuring, etc.

- **Path Module - Now Global API**
  - `path` object now available globally (no `require()` needed)
  - Consistent with `file` and `http` global API design
  - Backward compatible: `require('path')` still works
  - Usage: `path.join('a', 'b')` instead of `require('path').join('a', 'b')`
  - All methods accessible: join, resolve, dirname, basename, extname, sep
  - Unified global API pattern across all core modules

- **Module System Fix**
  - Added `require()` function to global scope (was previously missing)
  - Ensures backward compatibility for existing code
  - Modules can be accessed both globally and via `require()`

#### Fixed
- **Promise Error Propagation**
  - Fixed thenable detection for returned promises (Promise/A+ compliance)
  - Now properly detects promise-like objects with `.then()` method
  - Fixed: returning `Promise.reject()` from `.then()` now chains correctly
  - Fixed: rejections now propagate through `.then()` without error handlers
  - All rejection flows now reach `.catch()` handlers as expected
  - Resolves issue where nested promise rejections would silently fail

#### Documentation
- **Comprehensive Package Documentation** (Enhanced!)
  - Added godoc comments to all core packages and modules
  - **cmd/dougless**: Full CLI documentation with usage examples
  - **internal/repl**: Complete REPL API documentation
  - **internal/runtime**: Detailed runtime architecture documentation
  - **internal/event**: Event loop concurrency model documented
  - **internal/modules**: Complete module API documentation (NEW!)
    - console.go: Console API with all methods documented
    - timers.go: Timer system with event loop integration details
    - path.go: Cross-platform path manipulation utilities
    - file.go: Async file operations with permission details
    - http.go: HTTP client/server + WebSocket support
    - promise.go: Promise/A+ implementation with spec compliance notes
  - **internal/permissions**: Security system documentation (NEW!)
    - permissions.go: 450+ lines covering core permission manager
    - parser.go: CLI flag parsing with usage examples
    - prompt.go: Interactive prompt behavior and responses
  - **1,200+ lines of documentation added**
  - All exported types, functions, and methods documented
  - JavaScript usage examples provided for all APIs
  - Permission requirements documented where applicable
  - Thread-safety and concurrency notes included
  - Follows Go documentation conventions throughout

#### Testing & Quality
- Runtime module test coverage improved: 75.0% → 77.2%
- All existing tests pass with new global APIs
- Added comprehensive test suite for path global functionality
- Verified backward compatibility with `require()` system
- **Promise test suite: 18/18 tests passing** (previously 16/18)
- Removed all `t.Skip()` calls - no skipped tests remaining
- Removed all TODO comments from codebase

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

## [Phase 2] - October 2024 - COMPLETE ✅ **[UPDATED OCT 15, 2024]**

### Path Module (Global API)
- **Unique Dougless Feature**: `path` object available globally (no require needed)
- **Full Implementation**
  - `path.join()` - Join path segments with OS-specific separator
  - `path.resolve()` - Resolve paths to absolute paths
  - `path.dirname()` - Extract directory name from path
  - `path.basename()` - Extract filename with optional extension removal
  - `path.extname()` - Get file extension
  - `path.sep` - OS-specific path separator constant
  - Cross-platform compatibility (Windows/Unix)
  - Backward compatible: `require('path')` still supported

### File System Module (Global API) - **SIMPLIFIED OCT 15, 2024**
- **Unique Dougless Feature**: `files` object available globally (no require needed)
- **Simplified 3-Method API** (convention-based)
  - `files.read(path, callback)` - Read file OR list directory (trailing `/`)
    - Returns `null` (not error) when file doesn't exist
    - Replaces: read, readdir, exists
  - `files.write(path, [content], callback)` - Write file OR create directory
    - Auto-creates parent directories for file writes
    - Replaces: write, mkdir
  - `files.rm(path, callback)` - Remove file OR directory (recursive)
    - Idempotent - no error if path doesn't exist
    - Replaces: unlink, rmdir
- **62% Reduction**: 8 methods → 3 methods
- **Event Loop Integration** - All operations properly scheduled on event loop
- **Smart Defaults** - Convention over configuration design

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

## [Phase 5] - October 2024 - COMPLETE ✅

### Promises Implementation
- **Native Promise Support**
  - Full Promise/A+ compliant implementation
  - `new Promise(executor)` constructor with resolve/reject callbacks
  - `promise.then(onFulfilled, onRejected)` - Chain promises
  - `promise.catch(onRejected)` - Error handling
  - Promise chaining with automatic value propagation
  - State management (Pending, Fulfilled, Rejected)
  - Thread-safe operations with mutex protection
  - Event loop integration for async resolution

- **Promise Static Methods**
  - `Promise.resolve(value)` - Create resolved promise
  - `Promise.reject(reason)` - Create rejected promise
  - `Promise.all(promises)` - Wait for all promises to resolve
  - `Promise.race(promises)` - Wait for first promise to settle
  - `Promise.allSettled(promises)` - Wait for all promises to settle (resolve or reject)
  - `Promise.any(promises)` - Wait for first promise to fulfill

### ES6+ Transpilation Support
- **esbuild Integration**
  - Automatic transpilation of ES6+ syntax to ES5 for Goja compatibility
  - Support for modern JavaScript features:
    - Arrow functions (`=>`)
    - Template literals (\`string ${expr}\`)
    - Destructuring (`const {a, b} = obj`)
    - Spread operator (`...args`)
    - `let` and `const` declarations
    - Classes and class inheritance
    - async/await (transpiled to Promise chains)
  - Target: ES2017 (for async/await support)
  - Error reporting with line numbers and source context
  - Warning display for non-fatal issues
  - Seamless integration - developers write modern JS, runtime handles transpilation

### Architecture Improvements
- Promise state machine with proper lifecycle management
- Handler queuing for pending promises
- Automatic promise chaining with result unwrapping
- Integration with existing event loop for microtask scheduling
- esbuild API integration for build-time transpilation
- Source transformation pipeline in runtime execution

### Examples Added
- `examples/test-promise.js` - Basic promise creation and chaining
- `examples/test-promise-all.js` - Promise.all() usage patterns
- `examples/test-promise-race.js` - Promise.race() and competitive scenarios

### Testing & Quality Assurance
- **Comprehensive Promise Test Suite**
  - Basic promise resolution/rejection
  - Promise chaining behavior
  - Error propagation through chains
  - Promise.all() success and failure cases
  - Promise.race() with various timing scenarios
  - Promise.allSettled() with mixed results
  - Promise.any() with rejection handling
  - Edge cases and error conditions
- **Full Test Coverage** for promise module
- All transpilation features validated with example scripts

## [Phase 4] - October 2024 - COMPLETE ✅

### Security & Permissions System
- **Context-Aware Permission Management**
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
- **WebSocket Support** (Phase 6)
  - WebSocket client
  - WebSocket server
  - Real-time bidirectional communication
  - Connection management and broadcasting

- **Crypto & Security** (Phase 7)
  - Cryptographic operations
  - Hashing algorithms
  - Encryption/decryption
  - Secure random generation

- **Process & System Integration** (Phase 8)
  - Child process spawning
  - Environment variable access
  - System information queries
  - Signal handling

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
