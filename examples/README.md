# Dougless Runtime - Example Programs

This directory contains comprehensive examples demonstrating all features of Dougless Runtime.

## Quick Start

Run any example with:
```bash
./dougless examples/<filename>.js
```

For example:
```bash
./dougless examples/hello.js
```

---

## Examples Overview

### 1. `hello.js` - **START HERE!**
**Your first Dougless program**

A welcoming introduction that demonstrates:
- Basic console logging
- Unique global APIs (file, http)
- Module system with `require()`
- Async operations with timers

**Run time:** < 1 second

```bash
./dougless examples/hello.js
```

---

### 2. `console_features.js`
**All console operations**

Demonstrates every console method available in Dougless:
- `console.log()`, `console.warn()`, `console.error()`
- `console.time()` / `console.timeEnd()` for performance measurement
- `console.table()` for structured data visualization
- Multiple concurrent timers
- Edge cases and error handling

**Run time:** ~500ms

```bash
./dougless examples/console_features.js
```

**Key Features:**
- 10 different console demonstrations
- Performance timing examples
- Table formatting with arrays and objects
- Timer management best practices

---

### 3. `timers.js`
**Complete timer system demonstration**

Shows all timer functions and patterns:
- `setTimeout()` with various delays
- `setInterval()` for recurring execution
- `clearTimeout()` and `clearInterval()`
- Multiple timers running simultaneously
- Nested timers
- Edge cases (double-clearing, fake IDs)

**Run time:** ~6 seconds

```bash
./dougless examples/timers.js
```

**Key Features:**
- 11 different timer scenarios
- Execution order demonstrations
- Timer cancellation patterns
- Safety testing (no crashes!)

---

### 4. `files_basic.js`
**Simplified files API (no require needed!)**

Comprehensive file system operations using callback style:
- `files.read(path, callback)` - Read files or list directories
- `files.write(path, content, callback)` - Write files or create directories
- `files.rm(path, callback)` - Remove files or directories recursively

**Run time:** ~200ms  
**Side effects:** Creates and deletes temporary files/directories (auto-cleanup)

```bash
./dougless examples/files_basic.js
```

**What it does:**
1. Reads a file and checks if it exists (null = doesn't exist)
2. Lists directory contents (trailing `/`)
3. Writes files with auto-created parent directories
4. Creates empty directories
5. Removes files and directories recursively
6. Demonstrates idempotent removal

**Unique Features:**
- The `files` API is **globally available** - no `require('fs')` needed!
- **Convention-based**: trailing `/` for directories
- **3 methods** instead of 8+ traditional methods
- **Smart defaults**: auto-creates parent directories, null for existence checks

---

### 4b. `files_promise.js` ✨ **NEW!**
**Simplified files API using promises and async/await**

Same file operations as `files_basic.js` but using modern promise syntax:
- Callback style → `files.read(path, callback)`
- Promise style → `await files.read(path)`

**Run time:** ~200ms  
**Side effects:** Creates and deletes temporary files/directories (auto-cleanup)

```bash
./dougless examples/files_promise.js
```

**What it does:**
1. All file operations using async/await (much cleaner!)
2. Demonstrates Promise.all() for parallel operations
3. Shows both sequential and parallel file processing
4. Error handling with try/catch blocks

**Key Features:**
- ✅ **Cleaner code** - async/await instead of nested callbacks
- ✅ **Parallel operations** - Promise.all() for concurrent file I/O
- ✅ **Better error handling** - try/catch instead of error-first callbacks
- ✅ **Same API** - just omit the callback to get a Promise!

---

### 5. `http_demo.js`
**HTTP client and server (no require needed!)**

Full HTTP demonstration in 4 parts:
1. **HTTP Server** - Create a server with multiple routes
2. **HTTP Client (GET)** - Make GET requests
3. **HTTP Client (POST)** - Send POST requests with JSON
4. **External APIs** - Request data from public APIs

**Run time:** ~5 seconds, then keeps running  
**Port:** 3000  
**Note:** Press Ctrl+C to stop the server

```bash
./dougless examples/http_demo.js
```

**Features:**
- Server with routing (/, /health, /api/data, /api/echo)
- GET and POST requests
- JSON request/response handling
- External API integration
- Error handling (404s, invalid JSON)
- Both server and client in one file!

**Test the server:**
```bash
# In another terminal while the example is running:
curl http://localhost:3000
curl http://localhost:3000/health
curl -X POST http://localhost:3000/api/echo -d '{"test":"data"}'
```

**Unique Feature:** The `http` API is **globally available** - no `require('http')` needed!

---

### 6. `path_examples.js`
**Path manipulation (Global API)**

Comprehensive path module demonstration:
- `path.join()` - Join path segments
- `path.dirname()` - Get directory name
- `path.basename()` - Get file name (with optional extension removal)
- `path.extname()` - Get file extension
- `path.resolve()` - Resolve absolute paths
- `path.sep` - OS-specific path separator
- Complex path operations
- Backward compatibility with `require('path')`

**Run time:** < 100ms

```bash
./dougless examples/path_examples.js
```

**Unique Feature:** The `path` API is **globally available** - no `require()` needed!  
**Note:** `require('path')` still works for backward compatibility.

---

### 7. `sourcemap_examples.js`
**ES6+ transpilation with source maps**

Demonstrates how Dougless handles modern JavaScript:
- Arrow functions transpiled to ES5
- Template literals
- Error messages with accurate line numbers
- Source map support for debugging

**Run time:** < 100ms

```bash
./dougless examples/sourcemap_examples.js
```

**Key Feature:** Errors reference original source lines, not transpiled code!

---

### 8-12. Promise Examples
**ES6+ Promises with deterministic FIFO ordering**

- `test-promise.js` - Basic promise operations, chaining, error handling
- `test-promise-all.js` - `Promise.all()` - wait for all promises
- `test-promise-race.js` - `Promise.race()` - first to settle wins
- `test-promise-any.js` - `Promise.any()` - first to fulfill wins
- `test-promise-allsettled.js` - `Promise.allSettled()` - wait for all, never rejects

```bash
./dougless examples/test-promise.js
./dougless examples/test-promise-all.js
./dougless examples/test-promise-race.js
./dougless examples/test-promise-any.js
./dougless examples/test-promise-allsettled.js
```

**Note:** All promise methods include timing tests and edge cases.

---

### 13-15. WebSocket Examples
**Real-time bidirectional communication**

- `websocket_simple.js` - Basic WebSocket server
- `websocket_server.js` - WebSocket with message broadcasting
- `websocket_chat.js` - Multi-client chat application

```bash
./dougless examples/websocket_simple.js
```

**Key Feature:** WebSocket support integrated with HTTP server!

---

### 16. `test_permissions.js`
**Permission system demonstration**

Shows the interactive permission model:
- File read/write permissions
- Network access permissions
- Interactive prompts
- Permission caching

```bash
./dougless examples/test_permissions.js
```

---

## Example Categories

### Basic Features
- `hello.js` - Introduction to Dougless
- `console_features.js` - Console operations
- `timers.js` - setTimeout & setInterval

### Global APIs (Unique to Dougless)
- `files_basic.js` - File system callbacks (global `files`)
- `files_promise.js` - File system promises (global `files`) ✨ **NEW!**
- `http_demo.js` - HTTP client/server (global `http`)
- `path_examples.js` - Path manipulation (global `path`)

### ES6+ & Modern JavaScript
- `sourcemap_examples.js` - Transpilation and source maps
- `test-promise*.js` - Promise/A+ implementation (5 files)

### Real-time Communication
- `websocket_*.js` - WebSocket examples (3 files)

### Security
- `test_permissions.js` - Permission system

---

## Running Multiple Examples

```bash
# Run basic examples in sequence
./dougless examples/hello.js
./dougless examples/console_features.js
./dougless examples/timers.js
./dougless examples/file_operations.js
./dougless examples/path_examples.js
./dougless examples/sourcemap_examples.js

# Promise examples
./dougless examples/test-promise.js
./dougless examples/test-promise-all.js
./dougless examples/test-promise-race.js

# HTTP demo last (it keeps running)
./dougless examples/http_demo.js
```

---

## Learning Path

**Recommended order for learning Dougless:**

1. **`hello.js`** - Get familiar with the basics
2. **`console_features.js`** - Learn debugging and output
3. **`timers.js`** - Understand async operations
4. **`path_examples.js`** - Work with file paths
5. **`files_basic.js`** - Read/write files (callback style)
6. **`files_promise.js`** - Read/write files (async/await style) ✨
7. **`sourcemap_examples.js`** - ES6+ and transpilation
8. **`test-promise.js`** - Modern async with Promises
9. **`http_demo.js`** - Build web applications
10. **`websocket_simple.js`** - Real-time communication

---

## Key Dougless Features

### Global APIs (No Require!)
Unlike Node.js, Dougless makes common APIs globally available:

```javascript
// Node.js way
const fs = require('fs');
const http = require('http');

// Dougless way (callback style)
files.read('data.txt', (err, content) => { ... });  // Already global!
http.get(url, callback);                            // Already global!

// Dougless way (promise style) ✨
const content = await files.read('data.txt');  // No require, no wrapping!
const response = await http.get(url);          // Coming soon!
```

### Event Loop
All async operations use the event loop:
- File I/O
- HTTP requests
- Timers (setTimeout, setInterval)

### ES5.1 Support
Currently supports ES5.1 syntax via Goja engine:
- ✅ `const` and `let`
- ✅ Traditional functions
- ✅ Objects and arrays
- ❌ Arrow functions (coming with transpilation)
- ❌ async/await (Phase 5 roadmap)

---

## Getting Help

- **Documentation:** See `/docs` directory
  - [File API Guide](../docs/file_api.md)
  - [HTTP API Guide](../docs/http_api.md)
  - [REPL Guide](../docs/repl_guide.md)
- **Interactive Mode:** Run `./dougless` without arguments for REPL

---

## Example Output Expectations

All examples include:
- ✓ Success indicators
- ✗ Error indicators (when demonstrating error handling)
- Clear section headers
- Descriptive messages

Examples are designed to:
- Run without external dependencies (except HTTP demo uses a public API)
- Clean up after themselves (file operations)
- Be self-contained and educational
- Demonstrate best practices

---

**Happy coding with Dougless!** 🚀
