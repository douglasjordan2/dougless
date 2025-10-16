# Implementation Plan: Simplified `files` API

## Overview
Consolidate 8 file operations into 3 smart methods using convention-based routing.

---

## API Design

### `files.read(path, [callback])`
**Behavior based on path:**
- Trailing `/`: Read directory, return `string[]`
- No trailing `/`: Read file, return `string` content
- File doesn't exist: Return `null` (doubles as exists check)

**Signatures:**
```js
// Callback style
files.read(path: string, callback: (err, data) => void)
// data: string | string[] | null

// Promise style (callback omitted)
files.read(path: string): Promise<string | string[] | null>
```

---

### `files.write(path, [content], [callback])`
**Behavior based on parameters:**
- 2 args with trailing `/`: Create directory(ies) recursively
- 3 args: Write file (create parent dirs if needed)

**Path conventions:**
- Trailing `/`: Directory creation
- No trailing `/`: File write

**Signatures:**
```js
// Callback style - file
files.write(path: string, content: string, callback: (err) => void)

// Callback style - directory
files.write(path: string, callback: (err) => void)

// Promise style - file (callback omitted)
files.write(path: string, content: string): Promise<void>

// Promise style - directory (callback omitted)
files.write(path: string): Promise<void>
```

---

### `files.rm(path, [callback])`
**Behavior:**
- Unified removal using `os.RemoveAll()` (recursive, handles files + dirs)
- Gracefully handle non-existent paths (no error)

**Signatures:**
```js
// Callback style
files.rm(path: string, callback: (err) => void)

// Promise style (callback omitted)
files.rm(path: string): Promise<void>
```

---

## Implementation Steps

### Phase 1: Core Implementation
1. **Rename module** (`internal/modules/file.go` → `files.go`)
   - Rename `FileSystem` struct → `Files`
   - Update constructor `NewFileSystem` → `NewFiles`

2. **Implement `files.read()`**
   - Check if path ends with `/`
   - If yes: Call `os.ReadDir()`, return names array
   - If no: Call `os.Stat()` to check existence
     - Directory without slash: Error with helpful message
     - File: Call `os.ReadFile()`, return content
     - Not exists: Return `null` (not error)

3. **Implement `files.write()`**
   - Parameter detection logic:
     ```go
     hasExtension := strings.Contains(filepath.Base(path), ".")
     endsWithSlash := strings.HasSuffix(path, "/")
     isNullContent := content is nil/undefined
     ```
   - Decision tree:
     - `endsWithSlash && !isNullContent`: Error
     - `endsWithSlash || isNullContent`: `os.MkdirAll(path, 0755)`
     - `hasExtension || !isNullContent`: `os.MkdirAll(parentDir)` then `os.WriteFile(path)`

4. **Implement `files.rm()`**
   - Use `os.RemoveAll()` (handles files, empty dirs, and recursive dirs)
   - Check `os.IsNotExist(err)` → return success (idempotent)

### Phase 2: Error Enhancement
5. **Add intelligent error messages**
   ```go
   func enhanceError(path string, operation string, err error) string {
       if operation == "write" && !hasExtension(path) && !strings.HasSuffix(path, "/") {
           return fmt.Sprintf("%s\nHint: Add trailing slash for directory: '%s/'", err, path)
       }
       // ... more cases
   }
   ```

6. **Add path validation**
   - Check for empty paths
   - Validate no double slashes (`//`)
   - Warn about unusual patterns

### Phase 3: Integration
7. **Update `internal/runtime/runtime.go`**
   - Change global from `file` → `files`
   - Update initialization: `vm.Set("files", filesModule.Export(vm))`

8. **Update permissions integration**
   - `files.read()`: `PermissionRead`
   - `files.write()`: `PermissionWrite`
   - `files.rm()`: `PermissionWrite`

### Phase 4: Testing
9. **Unit tests** (`internal/modules/files_test.go`)
   - Read file: normal, missing, no permission
   - Read dir: with/without slash, empty dir
   - Write file: new, overwrite, nested paths
   - Write dir: single, nested, already exists
   - Remove: file, dir, recursive, non-existent
   - Edge cases: dotfiles, no extension, special chars

10. **Integration tests**
    - Create test fixtures
    - Run example scripts
    - Verify error messages

### Phase 5: Documentation & Examples
11. **Update WARP.md**
    - New API reference
    - Remove old method mentions
    - Update code snippets

12. **Create examples**
    - `examples/files_basic.js`: read/write files
    - `examples/files_dirs.js`: directory operations
    - `examples/files_advanced.js`: nested paths, error handling

13. **Migration guide** (`docs/migration_file_to_files.md`)
    - Old → New API mapping
    - Breaking changes list
    - Codemod examples

---

## Migration Path

### Breaking Changes
- `file` → `files` (global name change)
- All old methods removed: `readdir`, `mkdir`, `rmdir`, `unlink`, `exists`, `stat`

### Old → New Mapping
```js
// Old                              // New
file.read(path, cb)                 files.read(path, cb)
file.write(path, data, cb)          files.write(path, data, cb)
file.readdir(path, cb)              files.read(path + '/', cb)
file.mkdir(path, cb)                files.write(path + '/', null, cb)
file.rmdir(path, cb)                files.rm(path, cb)
file.unlink(path, cb)               files.rm(path, cb)
file.exists(path, cb)               files.read(path, cb) // check if null
file.stat(path, cb)                 // ❌ Removed (add if needed)
```

---

## Edge Cases to Handle

1. **Extensionless files** (README, Makefile)
   - Default: Treat as file if content provided
   - Error: Helpful message if ambiguous

2. **Hidden files/dirs** (`.git`, `.env`)
   - Work normally, respect conventions

3. **Trailing slash on file write**
   - Error: "Cannot write content to directory path"

4. **Nested path creation**
   - `files.write('a/b/c/file.js', 'content', cb)` → create a/, b/, c/ automatically

5. **Permissions on nested creates**
   - Check write permission on parent, not every level

6. **Empty content vs null**
   - `""` (empty string): Write empty file
   - `null`/`undefined`: Create directory

---

## Testing Checklist

### Automated Tests
- [ ] Read existing file
- [ ] Read missing file (returns null)
- [ ] Read directory with slash
- [ ] Read directory without slash (error)
- [ ] Write file with extension
- [ ] Write file without extension
- [ ] Write nested file (auto-create dirs)
- [ ] Write directory with slash
- [ ] Write directory with null content
- [ ] Remove file
- [ ] Remove directory (recursive)
- [ ] Remove non-existent (no error)
- [ ] Permission denials (all operations)
- [ ] Error message clarity

### Manual Testing
- [ ] Run existing examples with new API
- [ ] Test REPL interactive usage
- [ ] Verify error messages are helpful
- [ ] Check performance (no regression)

---

## Timeline Estimate

- **Phase 1-2**: 4-6 hours (core + errors)
- **Phase 3**: 1-2 hours (integration)
- **Phase 4**: 3-4 hours (testing)
- **Phase 5**: 2-3 hours (docs/examples)

**Total**: ~10-15 hours of focused work

---

## Open Questions

1. **Keep `stat()` as `files.stat()`?** Or remove entirely?
2. **Async version names?** `files.readAsync()` or assume all are async?
3. **Sync versions?** `files.readSync()` for framework use cases?
4. **Return values on write?** Currently `(err)`, could return `(err, bytesWritten)`

---

## Promise Support

All file operations now support both callback and promise-based APIs:

```javascript
// Callback style
files.read('data.txt', (err, content) => {
  if (err) console.error(err);
  else console.log(content);
});

// Promise style
files.read('data.txt')
  .then(content => console.log(content))
  .catch(err => console.error(err));

// Async/await style
async function readFile() {
  try {
    const content = await files.read('data.txt');
    console.log(content);
  } catch (err) {
    console.error(err);
  }
}
```

**Implementation Details:**
- When callback is omitted, operations return a Promise
- Promise resolves with data (read) or null (write/rm)
- Promise rejects with error message (string)
- Fully integrated with event loop for proper async handling

---

## Notes

Created: October 15, 2024
Updated: October 16, 2024 - Added promise support
Status: ✅ Complete - All features implemented and tested
