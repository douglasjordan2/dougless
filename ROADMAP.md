# Dougless Runtime - Development Roadmap

## Current Status

**Phases 1-6 COMPLETE! ✅**
- Phase 1: Foundation ✅
- Phase 2: File System & Modules ✅  
- Phase 3: Networking & HTTP ✅
- Phase 4: Security & Permissions ✅
- Phase 5: Promises & ES6+ ✅ **NEWLY COMPLETED** (Oct 15, 2024)
- Phase 6: WebSockets & Real-time ✅

All core async features, promises, and ES6+ transpilation are fully implemented, tested, and validated.

## Recently Completed

- ✅ **Config-Based Permissions System** - Project-centric permission model (Oct 17, 2025)
  - `.douglessrc` JSON configuration files
  - Per-project permission definitions (read, write, net, env, run)
  - Automatic config discovery from script directory
  - Cleaner alternative to CLI flags
  - Version-controlled permissions with project code
- ✅ **Phase 5: Promises & ES6+** - Full Promise/A+ implementation with all static methods (Oct 15, 2024)
  - Promise.all(), Promise.race(), Promise.allSettled(), Promise.any()
  - ES6+ transpilation with esbuild (arrow functions, async/await, classes, etc.)
  - Inline source maps for accurate error reporting
- ✅ **Interactive Permissions System** - Context-aware prompts with session caching
- ✅ **Phase 6: WebSockets & Real-time** - Full WebSocket server implementation

## Development Phases

### Phase 1: Foundation ✅ **COMPLETE**

#### Core Infrastructure
- ✅ Basic project structure and Go module setup
- ✅ Core runtime with Goja integration
- ✅ Event loop with proper async operation handling
- ✅ Module registry system with CommonJS-style require()
- ✅ Placeholder implementations for fs, http, and path modules

#### Timer System
- ✅ `setTimeout()` - Schedule one-time delayed execution
- ✅ `setInterval()` - Schedule recurring execution
- ✅ `clearTimeout()` - Cancel pending timeouts
- ✅ `clearInterval()` - Cancel active intervals
- ✅ Proper WaitGroup management for graceful shutdown

#### Console Operations
- ✅ `console.log()`, `console.error()`, `console.warn()` - Standard output
- ✅ `console.time()` / `console.timeEnd()` - Performance measurement
- ✅ `console.table()` - Structured data visualization with table formatting

#### REPL (Interactive Shell)
- ✅ Interactive JavaScript evaluation
- ✅ Multi-line input support (automatic detection)
- ✅ State preservation between commands
- ✅ Special commands (`.help`, `.exit`, `.clear`)
- ✅ Proper error handling and display

---

### Phase 2: File System & Modules ✅ **COMPLETE**

#### Path Module
- ✅ `path.join()` - Join path segments
- ✅ `path.resolve()` - Resolve absolute paths
- ✅ `path.dirname()` - Get directory name
- ✅ `path.basename()` - Get file name
- ✅ `path.extname()` - Get file extension
- ✅ `path.sep` - OS-specific path separator

#### File Module (Unique Global API)
- ✅ `file.read()` - Read file contents
- ✅ `file.write()` - Write data to file
- ✅ `file.readdir()` - List directory contents
- ✅ `file.exists()` - Check if path exists
- ✅ `file.mkdir()` - Create directory
- ✅ `file.rmdir()` - Remove directory
- ✅ `file.unlink()` - Delete file
- ✅ `file.stat()` - Get file/directory information
- ✅ Global access (no `require()` needed!)

---

### Phase 3: Networking & HTTP ✅ **COMPLETE**

#### HTTP Module (Unique Global API)
- ✅ `http.get()` - Make HTTP GET requests with callbacks
- ✅ `http.post()` - Make HTTP POST requests with JSON payload
- ✅ `http.createServer()` - Create HTTP server
- ✅ Server request/response handling
- ✅ Custom header support (`setHeader()`)
- ✅ Response status codes and body content
- ✅ Multiple header values support
- ✅ Global access (no `require()` needed!)

---

### Phase 4: WebSockets & Real-time ✅ **COMPLETE**

#### WebSocket Module
- ✅ `server.websocket(path, callbacks)` - Add WebSocket endpoint to server
- ✅ Real-time bidirectional communication
- ✅ Connection state management (`readyState` property)
- ✅ Browser-compatible API (CONNECTING, OPEN, CLOSING, CLOSED)
- ✅ Thread-safe message sending with mutex protection
- ✅ Broadcasting to multiple clients
- ✅ Event callbacks: `open`, `message`, `close`, `error`
- ✅ `ws.send()`, `ws.close()` methods on connection object

---

### Security & Permissions ✅ **COMPLETE**

#### Permissions System
- ✅ Interactive permission prompting with context awareness
- ✅ CLI flags for explicit permission granting (`--allow-read`, `--allow-write`, `--allow-net`)
- ✅ Context-aware defaults (interactive vs non-interactive)
- ✅ Session-based permission caching for "always" grants
- ✅ Thread-safe permission checking with concurrent access support
- ✅ Clear, actionable error messages
- ✅ Path canonicalization and security (prevents directory traversal)
- ✅ Host matching with wildcard and port support
- ✅ 30-second timeout for interactive prompts

---

### Testing & Quality ✅

- ✅ **Tests passing** (unit + integration)
- ✅ **~75% code coverage** across all packages
- ✅ **Benchmark suite** for performance tracking
- ✅ **Race condition testing** (thread-safe event loop)
- ✅ Full test coverage for file system and path modules
- ✅ WebSocket examples and documentation

---

## Permissions Enhancements

While the core permissions system is complete, these enhancements will improve usability and developer experience:

### Testing & Validation
- ⏳ Unit tests for `CheckWithPrompt` with mock prompter
- ⏳ Integration tests for permission caching behavior
- ⏳ Tests for context timeout scenarios
- ⏳ Tests for concurrent permission checks
- ⏳ Tests for `--prompt` and `--no-prompt` flags

### CLI Improvements
- ⏳ `--help` flag showing all available commands and options
- ⏳ Permission flags documentation in help text
- ⏳ Usage examples in help output
- ⏳ Version information (`--version` flag)

### Configuration File Support (Config-First Permission Model) ✅ **COMPLETE**

**Vision**: Deprecate CLI flags in favor of a cleaner, more project-centric permission model using configuration files.

#### Production Mode
- ✅ Permissions defined in `.douglessrc` JSON configuration file
- ✅ Config file discovery from script directory
- ✅ Clear error messages for missing or malformed configs
- ✅ JSON parsing with validation

#### Configuration File Format
- ✅ `.douglessrc` - Primary config format (JSON)
- ✅ JSON schema for storing default permissions:
  ```json
  {
    "permissions": {
      "read": ["/data", "./config"],
      "write": ["./output", "./logs"],
      "net": ["api.example.com", "localhost:3000"],
      "env": ["API_KEY", "DATABASE_URL"],
      "run": ["git", "npm"]
    }
  }
  ```
- ✅ Per-project permission profiles
- ✅ Config file validation and error reporting

#### Future Enhancements
- ⏳ Two-step interactive prompt flow to build `.douglessrc` during development
- ⏳ `.douglessrc.json` - Alternative explicit JSON extension support
- ⏳ Cascading config (global `~/.douglessrc` → project `.douglessrc`)
- ⏳ `dougless init` command to generate config template
- ⏳ Comments support in config (use JSONC parser)
- ⏳ CLI flag deprecation warnings

### Advanced Features
- ⏳ Persistent permission cache across runs
- ⏳ Permission audit logging (track what was accessed)
- ⏳ Color-coded terminal output for prompts
- ⏳ Wildcard patterns in permission grants (`--allow-read=/tmp/*.txt`)
- ⏳ Permission profiles (dev, test, prod)

---

## Phase 5: Promises & ES6+ ✅ **COMPLETE**

### Promises & Async/Await
- ✅ Promise constructor and basic Promise operations
- ✅ `Promise.resolve()` and `Promise.reject()`
- ✅ `Promise.all()`, `Promise.race()` - fully implemented and tested!
- ✅ `Promise.allSettled()`, `Promise.any()` - **COMPLETE** (Oct 15, 2024)
- ✅ async/await syntax support (ES6+ transpilation with esbuild)
- ✅ ES6+ transpilation (arrow functions, destructuring, classes, etc.)
- ⏳ Promise-based versions of file operations (Future enhancement)
- ⏳ Promise-based versions of HTTP operations (Future enhancement)

### Event Emitter Pattern (Future Enhancement)
- ⏳ `EventEmitter` class
- ⏳ `on()`, `once()`, `emit()`, `removeListener()`
- ⏳ Integration with HTTP server events

---

## Next Up: Phase 7 - Crypto & Security

### Cryptographic Functions
- ⏳ Hash functions (SHA-256, SHA-512, MD5)
- ⏳ HMAC for message authentication
- ⏳ Random number generation (cryptographically secure)
- ⏳ Base64 encoding/decoding
- ⏳ UUID generation

### Additional Security
- ⏳ HTTPS support
- ⏳ TLS/SSL certificate handling
- ⏳ Secure WebSocket (WSS) support

---

## Phase 8: Process & System Integration

### Process Management
- ⏳ `process.exit()` - Graceful shutdown
- ⏳ `process.env` - Environment variable access (with permissions)
- ⏳ `process.argv` - Command-line arguments
- ⏳ `process.cwd()` - Current working directory
- ⏳ Signal handling (SIGINT, SIGTERM, etc.)

### Subprocess Execution
- ⏳ `exec()` - Execute shell commands (with permissions)
- ⏳ `spawn()` - Spawn child processes
- ⏳ Stream handling for child process I/O

### OS Integration
- ⏳ Platform detection
- ⏳ System information queries
- ⏳ User/group information

---

## Phase 9: Performance & Optimization

### Runtime Optimizations
- ⏳ Object pooling for reduced allocations
- ⏳ JIT-style code caching
- ⏳ Lazy loading of modules
- ⏳ Memory profiling and optimization

### Build Optimizations
- ⏳ Binary size reduction
- ⏳ Startup time improvements
- ⏳ Static linking options

### Benchmarking
- ⏳ Comprehensive performance benchmarks vs Node.js
- ⏳ HTTP server throughput benchmarks
- ⏳ File I/O performance benchmarks
- ⏳ Memory usage profiling

---

## Post Phase 9: Package Manager

### Core Package Management
- 📦 Dependency resolution and installation (`dougless install <package>`)
- 📦 Package manifest (`dougless.json`) with version management
- 📦 Lock file for reproducible builds (`dougless-lock.json`)
- 📦 Support for npm registry compatibility
- 📦 Local module cache and `dougless_modules/` directory
- 📦 Enhanced `require()` to support external packages

### Package Manager Features
- 📦 Semantic versioning support
- 📦 Development vs production dependencies
- 📦 Script running (`dougless run <script>`)
- 📦 Package publishing capabilities
- 📦 Dependency auditing and security checks

---

## Future Considerations

### ES6+ Support
- 🎯 Transpilation pipeline (esbuild, Babel, or SWC)
- 🎯 Arrow functions, destructuring, template literals
- 🎯 Classes and modules (import/export)
- 🎯 Async iterators and generators
- 🎯 Optional chaining and nullish coalescing

### Advanced Features
- 🔮 Worker threads for parallel execution
- 🔮 Streaming APIs (ReadableStream, WritableStream)
- 🔮 Buffer implementation for binary data
- 🔮 Native addon support
- 🔮 Plugin system for framework extensibility

### Framework Foundation
- 🌐 Routing system
- 🌐 Middleware architecture
- 🌐 Template engine integration
- 🌐 Database adapters
- 🌐 Session management
- 🌐 Authentication/authorization utilities

---

## Phase ∞: Become Dependency-Free (Long-term Vision)

**Goal**: Build everything from scratch to maximize learning and understanding.

This is a distant-future, aspirational phase focused on replacing all external dependencies with custom implementations. While we'll use libraries like Goja and esbuild in the short-term for pragmatic development, the ultimate learning goal is to understand these systems deeply enough to build them ourselves.

### Custom JavaScript Engine
**Replace**: Goja  
**Why**: Deep understanding of JavaScript execution, memory management, and runtime internals

#### Core Components
- 🔬 **Lexer**: Tokenize JavaScript source code
  - Character stream processing
  - Token recognition (keywords, identifiers, operators, literals)
  - Position tracking for error reporting
  - Unicode support

- 🔬 **Parser**: Build Abstract Syntax Tree (AST)
  - Recursive descent parsing or Pratt parsing
  - ES5.1 specification compliance initially
  - Operator precedence handling
  - Error recovery and reporting
  - AST node types for all JavaScript constructs

- 🔬 **Bytecode Compiler**: Transform AST to executable bytecode
  - Instruction set design
  - Register or stack-based VM architecture
  - Optimization passes (constant folding, dead code elimination)
  - Symbol table and scope management

- 🔬 **Virtual Machine**: Execute bytecode
  - Instruction dispatch loop
  - Value representation (tagged pointers, NaN boxing, etc.)
  - Memory management and garbage collection
  - Call stack and exception handling
  - Native function integration

- 🔬 **Garbage Collector**: Automatic memory management
  - Mark-and-sweep algorithm
  - Generational collection for performance
  - Incremental/concurrent collection
  - Weak references and finalizers

- 🔬 **Runtime Library**: Built-in JavaScript objects
  - Object, Array, String, Number, Boolean
  - Function, Date, RegExp, Error
  - Math, JSON, console
  - Prototype chain implementation

**Learning Resources**:
- "Crafting Interpreters" by Robert Nystrom
- "Engineering a Compiler" by Cooper & Torczon
- V8 design documentation
- SpiderMonkey source code study

---

### Custom AST Transformer & Transpiler
**Replace**: esbuild  
**Why**: Master code transformation, understand compilation pipeline, enable custom optimizations

#### Transpilation Pipeline
- 🔬 **Parser**: ES2017+ to AST (can reuse from custom engine)
- 🔬 **AST Visitors**: Pattern-based tree traversal
  - Visitor pattern implementation
  - Transform async/await → promises
  - Transform arrow functions → regular functions
  - Transform let/const → var with proper scoping
  - Transform classes → constructor functions
  - Transform template literals → string concatenation
  - Transform destructuring → explicit assignments

- 🔬 **Scope Analyzer**: Track variable bindings
  - Lexical scope tracking
  - Variable hoisting
  - Closure detection
  - Conflict resolution

- 🔬 **Code Generator**: AST back to JavaScript
  - ES5.1 compliant output
  - Source map generation
  - Readable output formatting
  - Optimization opportunities

- 🔬 **Optimization Passes**:
  - Dead code elimination
  - Constant propagation
  - Function inlining
  - Tail call optimization

**Learning Resources**:
- "Compilers: Principles, Techniques, and Tools" (Dragon Book)
- Babel plugin development docs
- AST Explorer for experimentation
- esbuild source code analysis

---

### Custom WebSocket Implementation
**Replace**: gorilla/websocket  
**Why**: Understand network protocols, frame parsing, and real-time communication

#### WebSocket Protocol
- 🔬 **HTTP Upgrade Handling**:
  - Parse upgrade headers
  - Validate WebSocket handshake
  - Generate Sec-WebSocket-Accept key
  - Subprotocol negotiation

- 🔬 **Frame Parser**:
  - Bit-level frame structure parsing
  - Opcode handling (text, binary, close, ping, pong)
  - Masking/unmasking implementation
  - Fragmented message reassembly
  - Control frame handling

- 🔬 **Connection Management**:
  - Connection state machine (CONNECTING, OPEN, CLOSING, CLOSED)
  - Ping/pong keep-alive
  - Graceful close handshake
  - Error handling and recovery

- 🔬 **Message Queue**:
  - Buffered send/receive
  - Backpressure handling
  - Priority queuing

**Learning Resources**:
- RFC 6455 (WebSocket Protocol)
- "TCP/IP Illustrated" by W. Richard Stevens
- Wireshark for protocol analysis
- gorilla/websocket source code study

---

### Custom HTTP Client/Server
**Replace**: net/http (Go standard library)  
**Why**: Master HTTP protocol, connection handling, and server architecture

#### HTTP Implementation
- 🔬 **Request Parser**:
  - HTTP/1.1 request line parsing
  - Header parsing and validation
  - Chunked transfer encoding
  - Content-Length handling
  - URL parsing and query string extraction

- 🔬 **Response Builder**:
  - Status line generation
  - Header formatting
  - Body encoding
  - Chunked responses

- 🔬 **Connection Management**:
  - Keep-alive support
  - Connection pooling
  - Timeout handling
  - Concurrent connection limits

- 🔬 **TLS/SSL**:
  - Certificate validation
  - Encryption/decryption
  - Handshake protocol
  - Session resumption

**Learning Resources**:
- RFC 2616 (HTTP/1.1) and RFC 7540 (HTTP/2)
- "HTTP: The Definitive Guide"
- Go net/http source code
- nginx architecture study

---

### Custom Event Loop
**Replace**: Current event loop (keep but enhance)  
**Why**: Understand async I/O, concurrency patterns, and scheduling

#### Advanced Event Loop
- 🔬 **I/O Multiplexing**:
  - epoll/kqueue integration (platform-specific)
  - Non-blocking I/O
  - Edge-triggered vs level-triggered events

- 🔬 **Task Scheduling**:
  - Priority queues for task management
  - Microtask vs macrotask distinction
  - Fair scheduling algorithms
  - CPU affinity for worker threads

- 🔬 **Timer Management**:
  - Efficient timer wheel or heap
  - Microsecond precision
  - Timer coalescing for efficiency

**Learning Resources**:
- libuv documentation and source
- "The Art of Multiprocessor Programming"
- Linux epoll/BSD kqueue man pages
- Node.js event loop deep dive

---

### Development Approach

**Phase 1: Study & Design** (Per Component)
- 📚 Read specifications and RFCs
- 📚 Study existing implementations
- 📚 Design custom architecture
- 📚 Create proof-of-concept prototypes

**Phase 2: Implement & Test**
- 🔨 Build minimal viable version
- 🔨 Comprehensive unit tests
- 🔨 Integration with existing runtime
- 🔨 Performance benchmarking

**Phase 3: Optimize & Refine**
- ⚡ Profile and identify bottlenecks
- ⚡ Optimize hot paths
- ⚡ Memory usage improvements
- ⚡ Documentation and examples

**Phase 4: Production Hardening**
- 🛡️ Edge case handling
- 🛡️ Security auditing
- 🛡️ Stress testing
- 🛡️ Real-world validation

---

### Why This Matters

**Learning Goals**:
- 🎓 **Deep System Understanding**: Know how things work at the lowest level
- 🎓 **Problem-Solving Skills**: Face and solve complex engineering challenges
- 🎓 **Performance Intuition**: Understand why things are fast or slow
- 🎓 **Debugging Mastery**: Fix issues in code you fully understand
- 🎓 **Architecture Expertise**: Design large, complex systems from scratch

**Practical Benefits**:
- 🚀 **Custom Optimizations**: Optimize for specific use cases
- 🚀 **No External Dependencies**: Complete control and no supply chain risks
- 🚀 **Tailored Features**: Add exactly what Dougless needs
- 🚀 **Educational Tool**: Serve as learning material for others
- 🚀 **Zero Bloat**: Only include what's necessary

**Timeline**: This is a multi-year journey, tackled one component at a time as learning projects. There's no rush - the goal is deep understanding, not quick completion.

---

## Performance Targets

### Current Goals
- **Startup Time**: < 100ms for basic scripts
- **Memory Usage**: < 50MB for typical applications  
- **HTTP Throughput**: > 10,000 requests/second
- **File I/O**: Comparable to Node.js performance

### Long-term Goals
- Match or exceed Node.js performance for common operations
- Sub-10ms cold start for serverless deployments
- Minimal memory footprint for embedded use cases
- Native speed for hot-path operations

---

## Success Metrics

### Phase Completion Criteria
- ✅ All features implemented and documented
- ✅ Unit tests passing with >70% coverage
- ✅ Integration tests validating real-world usage
- ✅ Performance benchmarks meeting targets
- ✅ Example programs demonstrating capabilities
- ✅ Documentation complete and accurate

---

*Last Updated: October 15, 2024 - Phase 5 (Promises & ES6+) COMPLETE ✅ - All Promise static methods implemented and tested*
