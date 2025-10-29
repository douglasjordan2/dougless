# Dougless Runtime

> *"JavaScript is really quite nice."* — Ryan Dahl

## Overview

Dougless Runtime is a custom runtime designed with the end goal of serving a custom full-stack framework. Built on top of the Goja JavaScript engine (ES5.1), Dougless Runtime provides a clean, maintainable codebase with excellent JavaScript interoperability and a unique globals-first API design.

## Features

- 🚀 **High-performance JavaScript execution** using Goja (pure Go, ES5.1)
- ✨ **ES6+ Support** - Arrow functions, async/await, classes, and more via esbuild transpilation
- 🎯 **Native Promises** - Full Promise/A+ implementation with all static methods
- 🔒 **Security-first permissions** - Interactive prompts with context-aware defaults
- 🔐 **Cryptographic operations** - Built-in hash, HMAC, random, and UUID generation
- 🖥️ **Process & system integration** - Environment vars, signals, exit handling
- ✅ **File I/O operations** with async callback APIs
- ✅ **HTTP client and server** support
- 🌐 **Global-first API** - core functionality available without require()
- ⚡ **Event loop** with proper async operation handling
- 📦 **CommonJS module system** for additional modules

## 🔒 Security & Permissions

Dougless implements a comprehensive permission system that addresses security concerns while maintaining developer experience.

### Config-Based Permissions (Recommended)

Dougless uses a **project-centric permission model** with `.douglessrc` configuration files. This approach is cleaner than CLI flags and allows permissions to be version-controlled with your code.

#### Basic Usage

Create a `.douglessrc` file in your project directory:

```json
{
  "permissions": {
    "read": ["./examples", "/tmp"],
    "write": ["/tmp", "./output"],
    "net": ["api.example.com", "localhost:3000"],
    "env": ["API_KEY", "DATABASE_URL"],
    "run": ["git", "npm"]
  }
}
```

Then run your script:

```bash
./dougless script.js
# Automatically discovers and loads .douglessrc from current directory
```

#### Permission Types

- **`read`** - File system read access (files/directories)
- **`write`** - File system write access
- **`net`** - Network access (HTTP/WebSocket hosts)
- **`env`** - Environment variable access
- **`run`** - Subprocess execution permissions

#### Benefits

- ✅ **Version controlled** - Commit `.douglessrc` with your code
- ✅ **Team friendly** - Everyone uses same permissions
- ✅ **Self-documenting** - Clear list of what your app needs
- ✅ **No flag hell** - No complex CLI commands to remember
- ✅ **Project-centric** - Config files over command-line flags

### Interactive Permission Prompts (Fallback)

Dougless features **context-aware permission prompting** that balances security with usability:

#### Development Mode (Interactive Terminal)
When running scripts interactively, Dougless prompts for permissions as needed:

```bash
./dougless script.js

# When script tries to read a file:
⚠️  Permission request: read access to '/data/config.json'
Allow? (y/n): y
Save to .douglessrc? (y/n): y
✓ Granted and saved to .douglessrc

# Second access to same file - no prompt (cached)
# Different file - prompts again
```

**Prompt flow:**
1. `y` or `yes` to grant (or any other key to deny)
2. If granted: `y` or `yes` to save to `.douglessrc` (persist across runs), or `n` for session-only

#### Production/CI Mode (Non-Interactive)
Automatically uses **strict deny-by-default** in non-interactive environments:

```bash
echo "file.read('/etc/passwd')" | ./dougless
# Error: Permission denied - no prompts in non-interactive mode
```

### Explicit Permission Flags

For production deployments and fine-grained control:

```bash
# Grant specific file access
./dougless --allow-read=/data script.js

# Grant all read access
./dougless --allow-read script.js

# Grant network access to specific host
./dougless --allow-net=api.example.com script.js

# Multiple permissions
./dougless --allow-read=/data --allow-net=api.example.com script.js

# Grant all permissions (for trusted scripts)
./dougless --allow-all script.js

# Force strict mode even in interactive terminal
./dougless --no-prompt script.js
```

### Permission Types

- **`--allow-read[=<paths>]`** - File system read access
  - No path = allow all reads
  - With path = allow specific path and subdirectories
- **`--allow-write[=<paths>]`** - File system write access
- **`--allow-net[=<hosts>]`** - Network access (HTTP/WebSocket)
  - Supports wildcards: `*.example.com`
  - Port-specific: `localhost:3000`
- **`--allow-env[=<vars>]`** - Environment variable access (future)
- **`--allow-run[=<programs>]`** - Subprocess execution (future)
- **`--allow-all`** or **`-A`** - Grant all permissions

### Clear Error Messages

When permission is denied, Dougless provides actionable guidance:

```
Permission denied: read access to '/tmp/config.json'

Run your script with:
  dougless --allow-read=/tmp/config.json script.js

Or grant all read access:
  dougless --allow-read script.js

For dev, use:
  dougless --allow-all script.js

Or interactive mode:
  dougless --prompt script.js
```

### Smart Defaults

- ✅ **Interactive terminal** → Automatic prompt mode (convenient for dev)
- ✅ **Non-interactive** → Strict deny-by-default (secure for CI/production)
- ✅ **Context-aware** → Detects environment automatically
- ✅ **Session-based caching** → "always" grants persist for script lifetime
- ✅ **30-second timeout** → Auto-deny if no response to prompt
- ✅ **Thread-safe** → Concurrent permission checks handled correctly

### Security Benefits

1. **Prevent unauthorized file access** - Scripts can't read sensitive files without permission
2. **Control network access** - Prevent scripts from making unexpected HTTP requests
3. **Audit script behavior** - Interactive prompts reveal what scripts are trying to do
4. **Safe defaults** - Non-interactive environments are secure by default
5. **No silent failures** - Clear error messages guide proper usage

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

Dougless has a unique API with several distinctive design choices:

### Global File System Access
Dougless provides the global `files` object with a simplified, convention-based API that supports both callbacks and promises:

```javascript
// No require needed!

// Callback style
files.read('data.txt', function(err, data) {
    if (err) {
        console.error('Error:', err);
    } else {
        console.log('Content:', data);
    }
});

// Promise style
files.read('data.txt')
    .then(data => console.log('Content:', data))
    .catch(err => console.error('Error:', err));

// Async/await style (cleanest!)
async function readFile() {
    try {
        const data = await files.read('data.txt');
        console.log('Content:', data);
    } catch (err) {
        console.error('Error:', err);
    }
}

// Write a file (auto-creates parent directories)
await files.write('output.txt', 'Hello Dougless!');

// Read a directory (note the trailing slash)
const fileNames = await files.read('src/');
console.log('Files:', fileNames);

// Create a directory
await files.write('new-dir/');
console.log('Directory created');

// Delete a file or directory
await files.rm('old-file.txt');
console.log('Deleted');
```

### Convention-Based API Design
- **3 methods** instead of 8: `files.read()`, `files.write()`, `files.rm()`
- **Trailing `/`** indicates directory operations
- **Automatic parent directory creation** for file writes
- **Unified removal** for files and directories
- **Smart null handling** - missing files return `null` instead of error
- **Dual API** - supports both callbacks and promises/async-await

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

### Native Promises & Async/Await
Dougless has full Promise/A+ support built-in, with modern async/await syntax. All async file operations support promises automatically:

```javascript
// Promises are built into file operations!
async function processFiles() {
  try {
    // Direct async/await - no wrapping needed!
    const data1 = await files.read('file1.txt');
    const data2 = await files.read('file2.txt');
    console.log('Files loaded:', data1, data2);
  } catch (err) {
    console.error('Error:', err);
  }
}

processFiles();

// All Promise methods available
Promise.all([
  files.read('a.txt'),
  files.read('b.txt')
])
  .then(files => console.log('All files:', files))
  .catch(err => console.error('Failed:', err));

// Or create your own promises
const myPromise = new Promise((resolve, reject) => {
  setTimeout(() => resolve('Done!'), 1000);
});

await myPromise; // 'Done!'
```

### ES6+ Modern Syntax
Write modern JavaScript with automatic transpilation:

```javascript
// Arrow functions, template literals, destructuring
const users = ['Alice', 'Bob', 'Charlie'];
const greetings = users.map(user => `Hello, ${user}!`);

// Classes and inheritance
class Person {
  constructor(name) {
    this.name = name;
  }
  
  greet() {
    return `Hi, I'm ${this.name}`;
  }
}

const person = new Person('Douglas');
console.log(person.greet()); // "Hi, I'm Douglas"

// Async/await for clean async code
async function fetchData() {
  const response = await fetch('https://api.example.com/data');
  const data = await response.json();
  return data;
}
```

### Always Available Globals
```javascript
console.log('Logging');           // ✅ Built-in
files.read('file.txt', callback); // ✅ Built-in
http.get('http://...', callback); // ✅ Built-in
crypto.createHash('sha256');      // ✅ Built-in
process.exit(0);                  // ✅ Built-in
setTimeout(callback, 1000);       // ✅ Built-in
Promise.resolve(value);           // ✅ Built-in

const path = require('path');     // Module system still available
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

- **[ROADMAP.md](ROADMAP.md)** - Development phases, implementation status, and future plans
- **[Permissions Config Guide](docs/permissions_config.md)** - Complete guide to config-based permissions with `.douglessrc`
- **[REPL Guide](docs/repl_guide.md)** - Complete guide to using the interactive REPL shell
- **[Promises API Guide](docs/promises_api.md)** - Complete reference for Promises and async/await
- **[File API Guide](docs/file_api.md)** - Complete reference for the global `files` API
- **[HTTP API Guide](docs/http_api.md)** - Complete reference for the global `http` API
- **[Crypto API Guide](docs/crypto_api.md)** - Complete reference for the global `crypto` API
- **[Process Module Guide](docs/process_module.md)** - Complete reference for the global `process` object
- **[Changelog](CHANGELOG.md)** - Detailed history of changes and features

## Technology Stack

### Core Dependencies
- **[Goja](https://github.com/dop251/goja)** - Pure Go JavaScript engine (ES5.1)
- **[esbuild](https://esbuild.github.io/)** - Ultra-fast ES6+ to ES5 transpilation ✅
- **Go standard library** - For system operations, networking, and crypto

### Current Dependencies
- **[gorilla/websocket](https://github.com/gorilla/websocket)** - WebSocket implementation (Phase 6)

## Inspiration & References

### Similar Projects
- **[Node.js](https://nodejs.org/)** - The gold standard for JavaScript runtimes
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

📋 **See [ROADMAP.md](ROADMAP.md) for detailed implementation status and future plans.**

---

**Note**: This project is under active development. APIs and features are subject to change as we progress through the development phases.

