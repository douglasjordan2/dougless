# Dougless Runtime - TODO List

Generated: 2025-10-14  
Last Updated: 2025-10-14 23:19 UTC

## ✅ Completed Tasks

### 3. Clear TODOs from Code Comments ✅ (Oct 14, 2025)
- [x] Enabled source maps in transpiler
- [x] Fixed Promise error propagation
- [x] **Result:** Zero TODO comments in code

### 4. Remove t.Skip Calls ✅ (Oct 14, 2025)
- [x] Fixed `TestPromiseErrorPropagation`
- [x] Fixed `TestPromiseThenWithoutErrorHandler`
- [x] **Result:** Zero skipped tests

### 5. Complete Package Documentation ✅ (Oct 14, 2025)
- [x] Documented all internal/modules files (9 files)
  - console.go, timers.go, path.go, file.go, http.go, promise.go
- [x] Documented all internal/permissions files (3 files)
  - permissions.go, parser.go, prompt.go
- [x] Added package-level, type, constructor, and method comments
- [x] Included JavaScript usage examples for all APIs
- [x] **Result:** 1,200+ lines of professional-grade documentation

### 7. Clean LSP Errors ✅ (Oct 14, 2025)
- [x] Verified with `go vet ./...` - clean
- [x] **Result:** No LSP errors

---

## 🚧 Remaining Tasks

## Project Cleanup & Improvement Tasks

### 1. Add Missing Tests
- [ ] Identify modules and functions without test coverage
- [ ] Add comprehensive unit tests for:
  - Runtime initialization
  - Transpilation
  - Module system
  - Event loop edge cases
  - Promise implementation
  - File operations
  - HTTP operations
  - Permission system
- [ ] Run `go test -cover ./...` to check coverage

### 1.5. Complete Promise Implementation ✅ (COMPLETE - Oct 15, 2024)
**Status:** ✅ ALL METHODS IMPLEMENTED - PHASE 5 COMPLETE
- [x] `Promise.all()` - ✅ Working
- [x] `Promise.race()` - ✅ Working
- [x] `Promise.any()` - ✅ Working (was already implemented)
- [x] `Promise.allSettled()` - ✅ Working (newly implemented Oct 15)
- [x] Comprehensive test coverage (18/18 tests passing)
- [x] Example files for all methods
- [x] ES6+ transpilation with esbuild (async/await, arrow functions, etc.)
- [x] CHANGELOG.md updated
- [x] ROADMAP.md updated
- [x] **Phase 5 (Promises & ES6+) is now 100% COMPLETE**

### 1.6. Unify file.read() and file.readdir()
- [ ] Consider making `file.read()` smart enough to detect files vs directories
- [ ] Options:
  - Return consistent object structure: `{ type: 'file', content: '...' }` or `{ type: 'directory', entries: [...] }`
  - Keep both methods but make one an alias
  - Use `file.stat()` internally to determine the type
- [ ] Ensure backward compatibility or document breaking changes
- [ ] Update tests and examples if implemented

### 2. Fix Issues from Current Tests (MAINTENANCE)
- [x] All tests currently passing ✅
- [ ] Monitor for future test failures during development

### 3. Clear TODOs from Code Comments
- [ ] Search for TODO comments: `grep -r "TODO" --include="*.go" .`
- [ ] Resolve each TODO item by implementing the feature or removing if no longer relevant
- [ ] Document decisions for removed TODOs

### 4. Remove t.Skip Calls from Tests
- [ ] Find all t.Skip() calls: `grep -r "t.Skip" --include="*_test.go" .`
- [ ] Either implement the skipped tests or remove them if no longer needed
- [ ] Document rationale for any that must remain skipped

### 5. Analyze File Sizes and Folder Structure
- [ ] Review codebase organization
- [ ] Identify large files that could benefit from splitting (>500 lines)
- [ ] Ensure logical separation of concerns
- [ ] Consider refactoring into smaller, focused modules where appropriate
- [ ] Check with: `find . -name "*.go" -type f -exec wc -l {} + | sort -rn | head -20`

### 6. Ensure Package Documentation Comments are Complete ✅ (COMPLETE)
- [x] Reviewed all packages and exported functions/types
- [x] Added godoc comments following Go conventions
  - All comments start with the name of the documented item
  - Package comments added to all relevant packages
- [x] JavaScript usage examples provided for all APIs
- [x] **Status:** All core modules and packages fully documented
- [ ] Future: Add doc.go files for complex packages if needed
- [ ] Future: Consider adding godoc examples for testing

### 7. Clean LSP Errors
- [ ] Run `go vet ./...` and check for warnings
- [ ] Fix type issues
- [ ] Remove unused variables
- [ ] Correct incorrect imports
- [ ] Address any other static analysis warnings
- [ ] Run `go fmt ./...` to ensure formatting consistency

### 8. Security Sweep of Codebase
- [ ] Review code for security issues:
  - Input validation
  - Path traversal vulnerabilities
  - Permission bypasses
  - Unsafe operations
  - Error information leakage
- [ ] Ensure proper handling of user-supplied data throughout the codebase
- [ ] Review permission system for edge cases
- [ ] Check file operations for security issues
- [ ] Verify network operations are properly secured

## Quick Commands

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Find TODO comments
grep -r "TODO" --include="*.go" .

# Find t.Skip calls
grep -r "t.Skip" --include="*_test.go" .

# Check file sizes
find . -name "*.go" -type f -exec wc -l {} + | sort -rn | head -20

# Run go vet
go vet ./...

# Format code
go fmt ./...

# Check documentation
go doc -all ./internal/runtime
```

## Notes

- Maintain test coverage above 80%
- Follow Go conventions and best practices
- Document any architectural decisions
- Keep security as a top priority
