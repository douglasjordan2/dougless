# Dougless Runtime

A custom JavaScript runtime built in Go, designed to serve as the foundation for a custom full-stack framework powered by WebSockets. Built with the Goja engine and featuring ES6+ transpilation capabilities.

## Overview

Dougless Runtime is a custom runtime designed with the end goal of serving a custom full-stack framework powered by WebSockets. It's not inherently compatible with everything Node.js supports as it represents a new paradigm. This includes a custom system for building plugins to extend the framework. Built on top of the Goja JavaScript engine, Dougless Runtime provides a clean, maintainable codebase with excellent JavaScript interoperability.

In addition to the Goja engine, we are introducing a build-time tool that will compile ES6+ into ES5, enhancing compatibility and performance.

For more information on how esbuild integrates with Go, visit [esbuild Go API](https://pkg.go.dev/github.com/evanw/esbuild/pkg/api).

## Features (Planned)

- 🚀 **High-performance JavaScript execution** using Goja (pure Go)
- 📁 **File I/O operations** with both sync and async APIs
- 🌐 **HTTP client and server** support
- 🔌 **WebSocket** implementation for real-time applications
- ⚡ **Event loop** with proper async operation handling
- 📦 **CommonJS module system** with built-in modules
- 🔒 **Crypto utilities** and security features
- 🛠️ **Process management** and system integration

## Current Status

This project is in early development (Phase 1 - Foundation). Currently implemented:

### Core Infrastructure ✅
- ✅ Basic project structure and Go module setup
- ✅ Core runtime with Goja integration
- ✅ Event loop with proper async operation handling
- ✅ Module registry system with CommonJS-style require()
- ✅ Placeholder implementations for fs, http, and path modules

### Timer System ✅
- ✅ `setTimeout()` - Schedule one-time delayed execution
- ✅ `setInterval()` - Schedule recurring execution
- ✅ `clearTimeout()` - Cancel pending timeouts
- ✅ `clearInterval()` - Cancel active intervals
- ✅ Proper WaitGroup management for graceful shutdown

### Console Operations ✅
- ✅ `console.log()`, `console.error()`, `console.warn()` - Standard output
- ✅ `console.time()` / `console.timeEnd()` - Performance measurement
- ✅ `console.table()` - Structured data visualization with table formatting

### REPL (Interactive Shell) ✅
- ✅ Interactive JavaScript evaluation
- ✅ Multi-line input support (automatic detection)
- ✅ State preservation between commands
- ✅ Special commands (`.help`, `.exit`, `.clear`)
- ✅ Proper error handling and display

### Next Up (Phase 2)
- ⏳ File system operations (fs module)
- ⏳ Path manipulation utilities (path module)
- ⏳ Enhanced error handling with stack traces
- ⏳ Unit and integration tests

## Quick Start

### Prerequisites
- Go 1.21 or later

### Installation
```bash
git clone https://github.com/douglasjordan2/dougless.git
cd dougless
go mod tidy
```

### Build and Run
```bash
# Build the runtime
go build -o dougless cmd/dougless/main.go

# Start interactive REPL mode
./dougless

# Or run a JavaScript file
./dougless examples/hello.js
```

## Project Structure

```
dougless-runtime/
├── cmd/dougless/           # CLI entry point
├── internal/
│   ├── runtime/           # Core runtime logic
│   ├── repl/              # Interactive REPL implementation
│   ├── modules/           # Built-in modules (fs, http, path, etc.)
│   ├── event/             # Event loop implementation
│   └── bindings/          # Go-JS bindings and utilities
├── pkg/api/               # Public API (if needed as library)
├── examples/              # Example JavaScript programs
├── tests/                 # Test suite
└── docs/                  # Documentation
```

## Documentation

### Planning & Architecture
- **[Project Plan](docs/project_plan.md)** - Comprehensive development roadmap with 8 phases, technical architecture details, and success metrics
- **[Transpilation Strategy](docs/transpilation_strategy.md)** - Strategy for supporting modern ES6+ JavaScript syntax through transpilation to ES5
- **[REPL Guide](docs/repl_guide.md)** - Complete guide to using the interactive REPL shell
- **[Changelog](CHANGELOG.md)** - Detailed history of changes, features, and improvements

### Development Phases
1. **Foundation** (Current) - Basic runtime with console operations and timers
2. **File System & Modules** - File I/O and robust module system
3. **Networking & HTTP** - HTTP client/server capabilities
4. **WebSockets & Real-time** - WebSocket implementation
5. **Advanced Async & Promises** - Promise support and async/await
6. **Crypto & Security** - Cryptographic functions and security features
7. **Process & System Integration** - System-level operations
8. **Performance & Optimization** - Production-ready optimizations

## Technology Stack

### Core Dependencies
- **[Goja](https://github.com/dop251/goja)** - Pure Go JavaScript engine (ES5.1)
- **Go standard library** - For system operations, networking, and crypto

### Potential Future Dependencies
- **[esbuild Go API](https://pkg.go.dev/github.com/evanw/esbuild/pkg/api)** - For ES6+ transpilation
- **[gorilla/websocket](https://github.com/gorilla/websocket)** - WebSocket implementation

## Inspiration & References

### Similar Projects
- **[Node.js](https://nodejs.org/)** - The gold standard for JavaScript runtimes
- **[Deno](https://deno.land/)** - Modern JavaScript/TypeScript runtime
- **[Bun](https://bun.sh/)** - Fast all-in-one JavaScript runtime

### Technical Resources
- **[Goja Documentation](https://github.com/dop251/goja)** - JavaScript engine documentation
- **[Node.js API Reference](https://nodejs.org/api/)** - API compatibility reference
- **[Node.js Event Loop Guide](https://nodejs.org/en/docs/guides/event-loop-timers-and-nexttick/)** - Event loop implementation guidance
- **[WebSocket RFC 6455](https://tools.ietf.org/html/rfc6455)** - WebSocket protocol specification
- **[CommonJS Specification](http://wiki.commonjs.org/wiki/Modules/1.1)** - Module system specification

### Performance References
- **[V8 Performance Tips](https://v8.dev/docs/turbofan)** - JavaScript optimization insights
- **[Go Performance Best Practices](https://github.com/dgryski/go-perfbook)** - Go optimization techniques

## Development

### Running Tests
```bash
go test ./...
```

### Building for Different Platforms
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o dougless-linux cmd/dougless/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o dougless-macos cmd/dougless/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o dougless-windows.exe cmd/dougless/main.go
```

## Contributing

This project is in early development. Contributions, ideas, and feedback are welcome! Please see the [Project Plan](docs/project_plan.md) for current development priorities.

## License

[MIT License](LICENSE) (to be added)

## Goals

### Performance Targets
- **Startup Time**: < 100ms for basic scripts
- **Memory Usage**: < 50MB for typical applications  
- **HTTP Throughput**: > 10,000 requests/second
- **File I/O**: Comparable to Node.js performance

### Framework Goals
- Serve as the foundation for a WebSocket-powered full-stack framework
- Provide a custom plugin system for extending framework capabilities
- Support ES6+ JavaScript through build-time transpilation to ES5
- Create a new paradigm for web development focused on real-time communication

---

**Note**: This project is under active development. APIs and features are subject to change as we progress through the development phases.
