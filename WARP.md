# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

Dougless Runtime is a custom JavaScript runtime built in Go, designed to serve as the foundation for a custom full-stack framework powered by WebSockets. It uses the Goja JavaScript engine for ES5.1 execution with planned ES6+ transpilation support.

## Common Development Commands

### Building the Runtime
```bash
# Build the runtime executable
go build -o dougless cmd/dougless/main.go

# Build with optimizations
go build -ldflags="-s -w" -o dougless cmd/dougless/main.go

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o dougless-linux cmd/dougless/main.go
GOOS=darwin GOARCH=amd64 go build -o dougless-macos cmd/dougless/main.go
GOOS=windows GOARCH=amd64 go build -o dougless-windows.exe cmd/dougless/main.go
```

### Running JavaScript Files
```bash
# Run a JavaScript file
./dougless examples/hello.js

# Run with Go directly (without building)
go run cmd/dougless/main.go examples/hello.js

# Start REPL mode (interactive JavaScript shell)
./dougless

# REPL commands
# .help   - Show available commands
# .exit   - Exit the REPL
# .clear  - Clear the screen
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/runtime
go test ./internal/modules
go test ./internal/event

# Run benchmarks
go test -bench=. ./...
```

### Dependency Management
```bash
# Download dependencies
go mod download

# Update dependencies
go mod tidy

# Vendor dependencies
go mod vendor

# Verify dependencies
go mod verify
```

### Development Workflow
```bash
# Format code
go fmt ./...

# Run linter (install: go install golang.org/x/lint/golint@latest)
golint ./...

# Check for issues with go vet
go vet ./...

# Generate documentation
go doc -all ./internal/runtime
```

## Code Architecture

### Core Components

The runtime consists of several interconnected components that work together to execute JavaScript:

1. **Runtime Core** (`internal/runtime/runtime.go`)
   - Manages the Goja VM instance
   - Coordinates between event loop and module system
   - Provides global objects (console, timers)
   - Handles script execution lifecycle

2. **Event Loop** (`internal/event/loop.go`)
   - Implements non-blocking async operations using Go channels
   - Manages timer scheduling (setTimeout/setInterval)
   - Handles task prioritization and execution
   - Provides graceful shutdown capabilities

3. **Module System** (`internal/modules/`)
   - Registry pattern for built-in modules
   - CommonJS-style require() implementation
   - Module caching to prevent re-execution
   - Placeholder implementations for fs, http, and path modules

4. **REPL** (`internal/repl/repl.go`)
   - Interactive JavaScript shell
   - Multi-line input support (automatic detection)
   - Special commands (.help, .exit, .clear)
   - Maintains state between evaluations
   - Proper error handling and display

5. **CLI Entry Point** (`cmd/dougless/main.go`)
   - Dual-mode operation: REPL or file execution
   - REPL mode when no arguments provided
   - Script file execution mode with file argument
   - Error handling and reporting

### Module Architecture

Each built-in module implements the `Module` interface:
```go
type Module interface {
    Export(vm *goja.Runtime) goja.Value
}
```

Modules are registered in the registry and accessed via `require()`. Current modules:
- **path**: Path manipulation utilities (require-able)
- **file**: File system operations (global API)
- **http**: HTTP client/server functionality (global API)
- **promise**: Promise/A+ implementation (global API)

### Promise System

The Promise implementation provides full Promise/A+ compliance:
- **Promise Constructor**: `new Promise(executor)` with resolve/reject callbacks
- **Promise Methods**: `.then()`, `.catch()` for chaining and error handling
- **Static Methods**: `Promise.resolve()`, `Promise.reject()`, `Promise.all()`, `Promise.race()`, `Promise.allSettled()`, `Promise.any()`
- **State Management**: Proper state transitions (Pending → Fulfilled/Rejected)
- **Event Loop Integration**: Handler execution scheduled via event loop
- **Thread Safety**: Mutex-protected state management

### Event Loop Design

The event loop uses Go's concurrency primitives:
- **Channels** for task queuing
- **Goroutines** for parallel I/O operations
- **Context** for cancellation and cleanup
- **WaitGroup** for synchronization

Tasks can be scheduled with delays (timers) or executed immediately. The loop continues until all tasks complete or the context is cancelled.

## Development Status

**Phases 1-6 Status:**
- ✅ Phase 1: Foundation - COMPLETE
- ✅ Phase 2: File System & Modules - COMPLETE  
- ✅ Phase 3: Networking & HTTP - COMPLETE
- ✅ Phase 4: Security & Permissions - COMPLETE
- ✅ Phase 5: Promises & ES6+ - **COMPLETE** (Oct 15, 2024)
- ✅ Phase 6: WebSockets & Real-time - COMPLETE

### Phase 1 (Foundation) - COMPLETE ✅
- ✅ Basic project structure and Go module setup
- ✅ Goja JavaScript engine integration
- ✅ Event loop with async operation handling
  - Task scheduling via Go channels
  - WaitGroup synchronization for graceful shutdown
  - Context-based cancellation
- ✅ Module registry system with CommonJS require()
- ✅ Timer system (setTimeout/setInterval/clearTimeout/clearInterval)
  - One-time and recurring timers
  - Proper cleanup and WaitGroup management
  - UUID-based timer ID tracking
- ✅ Enhanced console operations
  - console.log/error/warn for output
  - console.time/timeEnd for performance measurement
  - console.table for structured data visualization
- ✅ Error handling improvements (stack traces, uncaught exceptions)
- ✅ Testing infrastructure (unit, integration, benchmarks)
- ✅ Interactive REPL (Read-Eval-Print Loop)
  - Multi-line input support
  - State preservation between commands
  - Special REPL commands (.help, .exit, .clear)
  - Proper error display

### Phase 2 (File System & Modules) - COMPLETE ✅
- ✅ Path module with full functionality (join, resolve, dirname, basename, extname)
- ✅ File system module with async operations (read, write, readdir, exists, mkdir, rmdir, unlink, stat)
- ✅ Unique global `file` API (no require needed)
- ✅ Event loop integration for async file operations

### Phase 3 (Networking & HTTP) - COMPLETE ✅
- ✅ HTTP client implementation (GET, POST with callbacks)
- ✅ HTTP server capabilities (createServer, listen)
- ✅ Request/response handling with headers and status codes
- ✅ Multiple header values support
- ✅ Unique global `http` API (no require needed)
- ✅ Event loop integration for async HTTP operations

### Phase 4 (Security & Permissions) - COMPLETE ✅
- ✅ Initial permission system implementation
- ✅ CLI flags for permission grants (--allow-read, --allow-write, --allow-net, etc.) **[DEPRECATED]**
- ✅ Interactive permission prompts in terminal mode
- ✅ Path-based and network-based granular controls
- ✅ Permission caching and session management
- ✅ Clear error messages with actionable suggestions

**Planned Evolution**: The CLI flag-based approach will be deprecated in favor of a config-first model:
- **Production**: Permissions defined exclusively in `.douglessrc` or `.douglessrc.json`
- **Development**: Two-step interactive prompt flow:
  1. First prompt: `Allow [operation]? yes/no` - grants permission for current session
  2. If 'yes': `Add to .douglessrc? yes/no` - optionally persists to config file
  3. This allows developers to build their `.douglessrc` incrementally during development
- **Goal**: Distinguish Dougless from Deno with a cleaner, more project-centric permission model

### Phase 5 (Promises & ES6+) - COMPLETE ✅
- ✅ Full Promise/A+ implementation
  - Promise constructor with resolve/reject
  - Promise chaining with .then() and .catch()
  - Static methods: resolve, reject, all, race, allSettled, any
- ✅ ES6+ Transpilation with esbuild
  - Arrow functions, template literals, destructuring
  - let/const declarations, spread operator
  - Classes and inheritance
  - async/await (transpiled to Promises)
  - Automatic source transformation
- ✅ Event loop integration for promise resolution
- ✅ Thread-safe promise state management

### Phase 6 (WebSockets & Real-time) - COMPLETE ✅
- ✅ WebSocket server integration (`server.websocket(path, callbacks)`)
- ✅ Real-time bidirectional communication
- ✅ Connection state management (readyState: CONNECTING, OPEN, CLOSING, CLOSED)
- ✅ Event callbacks: open, message, close, error
- ✅ Thread-safe message sending with mutex protection
- ✅ Broadcasting to multiple clients
- ✅ Working examples: websocket_simple.js, websocket_server.js, websocket_chat.js

### Phase 5 Complete! ✅ (October 15, 2024)
**Status:** ALL features implemented and tested
- ✅ Promise constructor and basic operations
- ✅ Promise.resolve() and Promise.reject()
- ✅ Promise.all() - fully implemented and tested
- ✅ Promise.race() - fully implemented and tested
- ✅ Promise.allSettled() - **NEWLY IMPLEMENTED** (Oct 15, 2024)
- ✅ Promise.any() - fully implemented and tested
- ✅ ES6+ transpilation with esbuild
- ✅ Full test coverage for all Promise methods
- ✅ Example files demonstrating all features

## Key Implementation Files

### Runtime Core
- **Runtime initialization**: `internal/runtime/runtime.go:42-60` (New() function)
- **Script execution**: `internal/runtime/runtime.go:63-94` (ExecuteFile/Execute)
- **Transpilation**: `internal/runtime/runtime.go:96-131` (transpile function)
- **Module loading**: `internal/runtime/runtime.go:169-183` (requireFunction)
- **Global initialization**: `internal/runtime/runtime.go:133-157` (initializeGlobals)

### Event Loop
- **Event loop core**: `internal/event/loop.go` (Run() method)
- **Task scheduling**: `internal/event/loop.go` (ScheduleTask)
- **Delayed task handling**: `internal/event/loop.go` (scheduleDelayedTask)
- **Timer cancellation**: `internal/event/loop.go` (ClearTimer)

### Promise System
- **Promise constructor**: `internal/modules/promise.go:30-55` (NewPromise)
- **Promise resolution**: `internal/modules/promise.go:57-79` (resolve method)
- **Promise rejection**: `internal/modules/promise.go:81-103` (reject method)
- **Promise chaining**: `internal/modules/promise.go:105-195` (Then method)
- **Promise catch**: `internal/modules/promise.go:197-199` (Catch method)
- **Promise.all**: `internal/modules/promise.go:201-280` (PromiseAll)
- **Promise.race**: `internal/modules/promise.go:282-340` (PromiseRace)
- **Promise.allSettled**: `internal/modules/promise.go:342-420` (PromiseAllSettled)
- **Promise.any**: `internal/modules/promise.go:422-490` (PromiseAny)
- **Setup function**: `internal/modules/promise.go:492-520` (SetupPromise)

### Transpilation
- **ES6+ Transpilation**: `internal/runtime/runtime.go:96-131` (transpile function)
  - Uses esbuild API for transformation
  - Target: ES2017 for async/await support
  - Handles errors and warnings with line-accurate reporting

## Future Development Phases

The project follows a multi-phase development plan:
1. Foundation ✅
2. File System & Modules ✅
3. Networking & HTTP ✅
4. Security & Permissions ✅
5. Promises & ES6+ ✅ (COMPLETE - Oct 15, 2024)
6. WebSockets & Real-time ✅
7. Crypto & Security (Next)
8. Process & System Integration
9. Performance & Optimization

### Post Phase 6: Package Manager
After completing the core runtime phases, a package management system is planned:
- **Package Installation**: `dougless install <package>` - npm-style package installation
- **Dependency Management**: `dougless.json` manifest with `dougless-lock.json` for reproducibility
- **Registry Integration**: Compatible with npm registry for package downloads
- **Module Resolution**: Enhanced `require()` to support `dougless_modules/` directory
- **Semver Support**: Version range resolution (`^`, `~`, `>=`, etc.)
- **Dependency Tree**: Recursive dependency resolution with conflict handling
- **Local Cache**: Global package cache at `~/.dougless/cache/`
- **CLI Commands**: install, uninstall, update, list

See `docs/project_plan.md` for detailed phase descriptions and milestones.

## Transpilation Strategy (IMPLEMENTED ✅)

The runtime now supports ES6+ through automatic transpilation to ES5:
- **Implementation**: esbuild (Go native integration)
- **Target**: ES2017 for async/await support
- **Approach**: On-the-fly transpilation during script execution
- **Features Supported**:
  - Arrow functions and template literals
  - let/const declarations
  - Destructuring and spread operator
  - Classes and inheritance
  - async/await (converted to Promise chains)
- **Error Handling**: Line-accurate error reporting with source context
- **Warnings**: Display non-fatal transpilation warnings to stderr

See `docs/transpilation_strategy.md` for complete strategy details and `internal/runtime/runtime.go:96-131` for implementation.

## Performance Targets

- Startup Time: < 100ms for basic scripts
- Memory Usage: < 50MB for typical applications
- HTTP Throughput: > 10,000 requests/second
- File I/O: Comparable to Node.js performance

## Important Considerations

1. **ES6+ Support**: Modern JavaScript syntax is now fully supported via esbuild transpilation.
2. **Module System**: CommonJS-style, not ES6 modules (import/export) - ES modules planned for future.
3. **Async Pattern**: ✅ Promises and async/await now fully supported!
4. **Security Model**: Config-first permission system using `.douglessrc`/`.douglessrc.json`. Development mode uses two-step prompts (yes/no → add to config?) to build permissions incrementally. CLI flags are deprecated.
5. **WebSocket Focus**: Core design goal is supporting real-time WebSocket applications.
6. **Plugin System**: Custom plugin architecture planned for framework extensibility.
7. **Global-First Design**: Core APIs (file, http, Promise) are available globally without require() - a unique Dougless feature.
