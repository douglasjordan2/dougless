# Dougless Runtime Examples

This directory contains comprehensive test examples for the Dougless runtime.

## Running Examples

Build the runtime first:
```bash
go build -o dougless ./cmd/dougless
```

Then run any example with appropriate permissions:

### Timers
Tests `setTimeout`, `setInterval`, `clearTimeout`, and `clearInterval`.
```bash
./dougless examples/timers.js
```

### Promises
Tests Promise creation, chaining, rejection, and static methods (`Promise.all`, `Promise.race`, `Promise.any`, `Promise.allSettled`).
```bash
./dougless examples/promises.js
```

### File System
Tests file read/write, directory operations, and cleanup. Requires file system permissions.
```bash
./dougless --allow-read --allow-write examples/files.js
```

### HTTP
Tests HTTP GET/POST requests and HTTP server creation. Requires network permissions.
```bash
./dougless --allow-net examples/http.js
```

### Process
Tests the process module including environment variables, command-line arguments, and system information.
```bash
./dougless examples/process.js
```

### Crypto
Tests cryptographic hash functions (MD5, SHA1, SHA256, SHA512).
```bash
./dougless examples/crypto.js
```

## Testing All Features

Run all examples:
```bash
./dougless --allow-all examples/timers.js
./dougless --allow-all examples/promises.js
./dougless --allow-all examples/files.js
./dougless --allow-all examples/http.js
./dougless --allow-all examples/process.js
./dougless --allow-all examples/crypto.js
```

## Notes

- The file system example creates temporary files in `/tmp/` and cleans them up
- The HTTP example starts a server on port 8080 temporarily
- All examples test the event loop deprecation by ensuring async operations complete correctly
