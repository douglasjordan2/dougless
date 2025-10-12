# Dougless Runtime

A custom JavaScript runtime built in Go, designed to eventually serve as the foundation for a custom full-stack framework (but that's a half-baked idea tbh).

## Overview

Dougless Runtime is a custom runtime designed with the end goal of serving a custom full-stack framework. Built on top of the Goja JavaScript engine (ES5.1), Dougless Runtime provides a clean, maintainable codebase with excellent JavaScript interoperability and a unique globals-first API design.

## Features

- 🚀 **High-performance JavaScript execution** using Goja (pure Go, ES5.1)
- ✅ **File I/O operations** with async callback APIs
- ✅ **HTTP client and server** support
- 🌐 **Global-first API** - core functionality available without require()
- ⚡ **Event loop** with proper async operation handling
- 📦 **CommonJS module system** for additional modules

### Planned
- 🔌 **WebSocket** implementation for real-time applications
- 📦 **Package manager** - npm/bun-style dependency management (`dougless install`)
- 🔒 **Crypto utilities** and security features
- 🛠️ **Process management** and system integration
- 🎯 **ES6+ support** through transpilation (future phases)

## Current Status

**Phase 1 (Foundation), Phase 2 (File System & Modules), and Phase 3 (Networking & HTTP) are COMPLETE! ✅**

All features are fully implemented, tested, and validated.

Currently implemented:

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

### Path Module ✅
- ✅ `path.join()` - Join path segments
- ✅ `path.resolve()` - Resolve absolute paths
- ✅ `path.dirname()` - Get directory name
- ✅ `path.basename()` - Get file name
- ✅ `path.extname()` - Get file extension
- ✅ `path.sep` - OS-specific path separator

### File Module ✅ (Unique Global API)
- ✅ `file.read()` - Read file contents
- ✅ `file.write()` - Write data to file
- ✅ `file.readdir()` - List directory contents
- ✅ `file.exists()` - Check if path exists
- ✅ `file.mkdir()` - Create directory
- ✅ `file.rmdir()` - Remove directory
- ✅ `file.unlink()` - Delete file
- ✅ `file.stat()` - Get file/directory information
- ✅ Global access (no `require()` needed!)

### HTTP Module ✅ (Unique Global API)
- ✅ `http.get()` - Make HTTP GET requests with callbacks
- ✅ `http.post()` - Make HTTP POST requests with JSON payload
- ✅ `http.createServer()` - Create HTTP server
- ✅ Server request/response handling
- ✅ Custom header support (`setHeader()`)
- ✅ Response status codes and body content
- ✅ Multiple header values support
- ✅ Global access (no `require()` needed!)

### Testing & Quality ✅
- ✅ **25/25 tests passing** (unit + integration)
- ✅ **~75% code coverage** across all packages
- ✅ **Benchmark suite** for performance tracking
- ✅ **Race condition testing** (thread-safe event loop)
- ✅ Full test coverage for file system and path modules

### Next Up (Phase 4)
- ⏳ WebSocket client and server
- ⏳ Real-time bidirectional communication
- ⏳ Connection management and broadcasting

### Future Features
- 📦 **Package Manager** (Post Phase 4)
  - Dependency resolution and installation (`dougless install <package>`)
  - Package manifest (`dougless.json`) with version management
  - Lock file for reproducible builds (`dougless-lock.json`)
  - Support for npm registry compatibility
  - Local module cache and `dougless_modules/` directory
  - Enhanced `require()` to support external packages

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

## Unique Dougless Features

Dougless has a unique API that sets it apart from Node.js, Deno, and Bun:

### Global File System Access
Unlike Node.js which requires `const fs = require('fs')`, Dougless provides the `file` object globally:

```javascript
// No require needed!
file.read('data.txt', function(err, data) {
    if (err) {
        console.error('Error:', err);
    } else {
        console.log('Content:', data);
    }
});

file.write('output.txt', 'Hello Dougless!', function(err) {
    if (err) console.error(err);
});
```

### Simplified Method Names
- `file.read()` instead of `fs.readFile()`
- `file.write()` instead of `fs.writeFile()`
- Clean, intuitive API design

### Global HTTP Access
Unlike Node.js which requires `const http = require('http')`, Dougless provides the `http` object globally:

```javascript
// Create an HTTP server - no require needed!
const server = http.createServer((req, res) => {
  res.setHeader('Content-Type', 'application/json')
  res.statusCode = 200
  res.end(JSON.stringify({ message: 'Hello from Dougless!' }))
})

server.listen(3000, () => {
  console.log('Server running on port 3000')
})

// Make HTTP requests - also global!
http.get('http://api.example.com/data', (err, response) => {
  if (!err) {
    console.log('Response:', response.body)
  }
})
```

### Always Available Globals
```javascript
console.log('Logging');          // ✅ Built-in
file.read('file.txt', callback); // ✅ Built-in
http.get('http://...', callback);// ✅ Built-in
setTimeout(callback, 1000);      // ✅ Built-in

const path = require('path');    // Module system still available
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
- **[REPL Guide](docs/repl_guide.md)** - Complete guide to using the interactive REPL shell
- **[File API Guide](docs/file_api.md)** - Complete reference for the global `file` API with examples
- **[HTTP API Guide](docs/http_api.md)** - Complete reference for the global `http` API with examples
- **[HTTP Design](docs/http_design.md)** - HTTP module design and implementation details
- **[Changelog](CHANGELOG.md)** - Detailed history of changes, features, and improvements

### Development Phases
1. **Foundation** ✅ - Basic runtime with console operations and timers
2. **File System & Modules** ✅ - File I/O and robust module system
3. **Networking & HTTP** ✅ - HTTP client/server capabilities
4. **WebSockets & Real-time** (Current) - WebSocket implementation
5. **Advanced Async & Promises** - Promise support and async/await
6. **Crypto & Security** - Cryptographic functions and security features
7. **Process & System Integration** - System-level operations
8. **Performance & Optimization** - Production-ready optimizations

## Technology Stack

### Core Dependencies
- **[Goja](https://github.com/dop251/goja)** - Pure Go JavaScript engine (ES5.1)
- **Go standard library** - For system operations, networking, and crypto

### Planned Dependencies
- **[gorilla/websocket](https://github.com/gorilla/websocket)** - WebSocket implementation (Phase 4)

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
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage report
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

**Current Test Status**: ✅ 25/25 passing | ~75% coverage

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
- Serve as the foundation for a full-stack framework
- Provide a custom plugin system for extending framework capabilities
- Create a new paradigm for web development focused on real-time communication
- Globals-first API design for simplicity and developer experience

---

**Note**: This project is under active development. APIs and features are subject to change as we progress through the development phases.
