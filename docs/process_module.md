# Process Module

The `process` module provides Node.js-compatible process-level operations and system information. It is available globally without requiring imports.

## Table of Contents
- [Properties](#properties)
- [Methods](#methods)
- [Events](#events)
- [Examples](#examples)

## Properties

### process.argv
**Type:** `Array<string>`

Array of command-line arguments. The first element is the runtime executable path, the second is the script path, and subsequent elements are script arguments.

```javascript
// Running: ./dougless script.js arg1 arg2
console.log(process.argv);
// Output: ['./dougless', 'script.js', 'arg1', 'arg2']
```

### process.env
**Type:** `Object<string, string>`

Object containing all environment variables as key-value pairs.

```javascript
console.log(process.env.HOME);    // '/home/username'
console.log(process.env.PATH);    // System PATH
console.log(process.env.USER);    // Current username
```

### process.pid
**Type:** `number`

The process ID (PID) of the current runtime instance.

```javascript
console.log('Process ID:', process.pid);
```

### process.platform
**Type:** `string`

The operating system platform: `'linux'`, `'darwin'` (macOS), `'windows'`, etc.

```javascript
console.log('Platform:', process.platform);

if (process.platform === 'linux') {
  console.log('Running on Linux');
}
```

### process.arch
**Type:** `string`

The CPU architecture: `'amd64'`, `'arm64'`, `'386'`, etc.

```javascript
console.log('Architecture:', process.arch);
```

### process.version
**Type:** `string`

The Dougless runtime version string.

```javascript
console.log('Dougless version:', process.version);
// Output: 'v0.8.0'
```

## Methods

### process.exit([code])
**Parameters:**
- `code` (optional): Exit code, defaults to 0

Immediately terminates the process with the specified exit code. Calls any registered exit handlers before exiting.

```javascript
// Exit with success
process.exit(0);

// Exit with error
process.exit(1);

// Exit with default code (0)
process.exit();
```

### process.cwd()
**Returns:** `string`

Returns the current working directory as an absolute path.

```javascript
const currentDir = process.cwd();
console.log('Working directory:', currentDir);
```

### process.chdir(directory)
**Parameters:**
- `directory`: Path to change to (absolute or relative)

Changes the current working directory to the specified path.

```javascript
console.log('Before:', process.cwd());
process.chdir('/tmp');
console.log('After:', process.cwd());
```

### process.on(event, callback)
**Parameters:**
- `event`: Event name (`'exit'`, `'SIGINT'`, `'SIGTERM'`, `'SIGHUP'`)
- `callback`: Function to call when event occurs

Registers an event handler for process events.

```javascript
// Exit event - called before process exits
process.on('exit', function(code) {
  console.log('Exiting with code:', code);
});

// SIGINT - Ctrl+C
process.on('SIGINT', function(signal) {
  console.log('Received', signal);
  process.exit(0);
});

// SIGTERM - termination signal
process.on('SIGTERM', function(signal) {
  console.log('Terminating:', signal);
  cleanup();
  process.exit(0);
});
```

## Events

### 'exit'
Emitted when the process is about to exit. The callback receives the exit code.

**Callback signature:** `function(code)`

```javascript
process.on('exit', function(code) {
  console.log('Process exiting with code:', code);
  // Cleanup operations here
});
```

### 'SIGINT'
Emitted when the process receives a SIGINT signal (Ctrl+C in terminal).

**Callback signature:** `function(signal)`

```javascript
process.on('SIGINT', function(signal) {
  console.log('Received interrupt signal');
  // Graceful shutdown
  process.exit(0);
});
```

### 'SIGTERM'
Emitted when the process receives a SIGTERM signal (termination request).

**Callback signature:** `function(signal)`

```javascript
process.on('SIGTERM', function(signal) {
  console.log('Received termination signal');
  cleanup();
  process.exit(0);
});
```

### 'SIGHUP'
Emitted when the process receives a SIGHUP signal (terminal hangup).

**Callback signature:** `function(signal)`

```javascript
process.on('SIGHUP', function(signal) {
  console.log('Terminal disconnected');
  // Reload configuration or restart
});
```

## Examples

### Basic Process Information
```javascript
console.log('Process Info:');
console.log('- PID:', process.pid);
console.log('- Platform:', process.platform);
console.log('- Architecture:', process.arch);
console.log('- Version:', process.version);
console.log('- Directory:', process.cwd());
```

### Command-line Arguments
```javascript
// script.js
console.log('Script:', process.argv[1]);
console.log('Arguments:', process.argv.slice(2));

// Running: ./dougless script.js --port 3000 --host localhost
// Output:
// Script: script.js
// Arguments: ['--port', '3000', '--host', 'localhost']
```

### Environment Variables
```javascript
// Check if running in production
if (process.env.NODE_ENV === 'production') {
  console.log('Production mode');
} else {
  console.log('Development mode');
}

// Get API key from environment
const apiKey = process.env.API_KEY;
if (!apiKey) {
  console.error('API_KEY environment variable not set');
  process.exit(1);
}
```

### Graceful Shutdown
```javascript
let server = http.createServer(function(req, res) {
  res.writeHead(200);
  res.end('Hello World');
});

server.listen(3000, function() {
  console.log('Server running on port 3000');
});

// Handle Ctrl+C
process.on('SIGINT', function() {
  console.log('Shutting down gracefully...');
  server.close();
  process.exit(0);
});

// Handle termination signal
process.on('SIGTERM', function() {
  console.log('Received SIGTERM, shutting down...');
  server.close();
  process.exit(0);
});
```

### Exit Handler for Cleanup
```javascript
const tempFiles = [];

function createTempFile(name) {
  tempFiles.push(name);
  files.write(name, 'temp data');
}

// Cleanup on exit
process.on('exit', function(code) {
  console.log('Cleaning up temp files...');
  tempFiles.forEach(function(file) {
    files.rm(file);
  });
});

createTempFile('/tmp/test1.txt');
createTempFile('/tmp/test2.txt');
```

### Directory Management
```javascript
console.log('Starting directory:', process.cwd());

// Change to subdirectory
process.chdir('src');
console.log('Changed to:', process.cwd());

// Go back up
process.chdir('..');
console.log('Back to:', process.cwd());
```

## Node.js Compatibility

The Dougless `process` module provides a subset of Node.js's `process` API. Currently implemented:

- ✅ `process.argv`
- ✅ `process.env`
- ✅ `process.exit()`
- ✅ `process.cwd()`
- ✅ `process.chdir()`
- ✅ `process.pid`
- ✅ `process.platform`
- ✅ `process.arch`
- ✅ `process.version`
- ✅ `process.on('exit')`
- ✅ Signal handling (SIGINT, SIGTERM, SIGHUP)

Not yet implemented (planned for future phases):
- ⏳ `process.stdin`, `process.stdout`, `process.stderr`
- ⏳ `process.nextTick()`
- ⏳ `process.uptime()`
- ⏳ `process.memoryUsage()`
- ⏳ `process.hrtime()`
- ⏳ Additional signal handlers

## See Also

- [Examples](../examples/process_demo.js)
- [Runtime Documentation](../WARP.md)
- [Node.js Process Documentation](https://nodejs.org/api/process.html)
