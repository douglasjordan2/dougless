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

4. **CLI Entry Point** (`cmd/dougless/main.go`)
   - Simple command-line interface
   - Script file execution
   - Error handling and reporting

### Module Architecture

Each built-in module implements the `Module` interface:
```go
type Module interface {
    Export(vm *goja.Runtime) goja.Value
}
```

Modules are registered in the registry and accessed via `require()`. Current placeholder modules:
- **fs**: File system operations (readFile, writeFile, readdir)
- **http**: HTTP client/server functionality
- **path**: Path manipulation utilities

### Event Loop Design

The event loop uses Go's concurrency primitives:
- **Channels** for task queuing
- **Goroutines** for parallel I/O operations
- **Context** for cancellation and cleanup
- **WaitGroup** for synchronization

Tasks can be scheduled with delays (timers) or executed immediately. The loop continues until all tasks complete or the context is cancelled.

## Development Status

Currently in **Phase 1** (Foundation) - Nearly Complete:

### Completed ✅
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

### In Progress ⏳
- Error handling improvements (stack traces, uncaught exceptions)
- Testing infrastructure (unit, integration, benchmarks)

### Up Next (Phase 2)
- Path module implementation
- File system operations (sync and async)
- Module system enhancements

## Key Implementation Files

### Runtime Core
- **Runtime initialization**: `internal/runtime/runtime.go:26-43` (New() function)
- **Script execution**: `internal/runtime/runtime.go:46-71` (ExecuteFile/Execute)
- **Module loading**: `internal/runtime/runtime.go:230-243` (requireFunction)

### Event Loop
- **Event loop core**: `internal/event/loop.go:40-59` (Run() method)
- **Task scheduling**: `internal/event/loop.go:86-93` (ScheduleTask)
- **Delayed task handling**: `internal/event/loop.go:96-115` (scheduleDelayedTask)
- **Timer cancellation**: `internal/event/loop.go:118-127` (ClearTimer)

### Timer System
- **Timer helper**: `internal/runtime/runtime.go:167-198` (delayHelper)
- **setTimeout/setInterval**: `internal/runtime/runtime.go:201-207` (setTimeout/setInterval)
- **Clear functions**: `internal/runtime/runtime.go:209-227` (clearHelper/clearTimeout/clearInterval)

### Console Operations
- **Basic console**: `internal/runtime/runtime.go:101-128` (consoleLog/Error/Warn)
- **Timer console**: `internal/runtime/runtime.go:130-165` (consoleTime/timeEnd)
- **Table console**: `internal/runtime/runtime.go:170-279` (consoleTable/helpers)

## Future Development Phases

The project follows an 8-phase development plan:
1. Foundation (current)
2. File System & Modules
3. Networking & HTTP
4. WebSockets & Real-time
5. Advanced Async & Promises
6. Crypto & Security
7. Process & System Integration
8. Performance & Optimization

See `docs/project_plan.md` for detailed phase descriptions and milestones.

## Transpilation Strategy

The runtime will support ES6+ through build-time transpilation to ES5:
- Primary tool: esbuild (Go native integration)
- Alternative options: Babel, SWC, TypeScript compiler
- Implementation approaches: on-the-fly, build-time, or hybrid

See `docs/transpilation_strategy.md` for complete strategy details.

## Performance Targets

- Startup Time: < 100ms for basic scripts
- Memory Usage: < 50MB for typical applications
- HTTP Throughput: > 10,000 requests/second
- File I/O: Comparable to Node.js performance

## Important Considerations

1. **Goja Limitations**: The engine supports ES5.1, not ES6+. Modern syntax requires transpilation.
2. **Module System**: CommonJS-style, not ES6 modules (import/export).
3. **Async Pattern**: Currently using callbacks; Promises/async-await planned for Phase 5.
4. **WebSocket Focus**: Core design goal is supporting real-time WebSocket applications.
5. **Plugin System**: Custom plugin architecture planned for framework extensibility.
