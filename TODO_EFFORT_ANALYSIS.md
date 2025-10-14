# Dougless Runtime - TODO Effort Analysis

**Generated:** 2025-10-14  
**Analyzed by:** Warp AI

## Overview

Based on current codebase analysis:
- **Total test coverage:** ~69.7% (modules), 75.0% (runtime), 86.8% (event)
- **TODO comments in code:** 4
- **Skipped tests (t.Skip):** 2
- **go vet issues:** 0 (clean! ✅)
- **Largest files:** http.go (595 lines), promise.go (481 lines), permissions.go (384 lines)

---

## Task Breakdown by Effort Level

### 🟢 LOW EFFORT (1-4 hours each)

#### Task 3: Clear TODOs from Code Comments
**Effort:** 1-2 hours  
**Current state:** Only 4 TODO comments in code  
**Locations:**
1. `internal/runtime/runtime.go` - Enable source maps (line ~117)
2. `internal/runtime/runtime.go` - Add more modules comment (line ~159)
3-4. `internal/modules/promise_test.go` - Two skipped tests

**Work required:**
- Enable source maps in transpiler config (10 mins)
- Remove/clarify the "add more modules" comment (5 mins)
- Address the two skipped Promise tests (see Task 4)

**Complexity:** Very low - minimal code changes needed

---

#### Task 4: Remove t.Skip Calls from Tests
**Effort:** 2-3 hours  
**Current state:** 2 skipped tests in `promise_test.go`
1. Promise chaining with nested Promise.reject()
2. Error propagation through .then() without handler

**Work required:**
- Debug and fix Promise rejection chaining (likely a small logic issue)
- Fix error propagation when onRejected handler is nil
- Re-run tests to verify

**Complexity:** Low-medium - Promise logic is already mostly working, these are edge cases

---

#### Task 7: Clean LSP Errors
**Effort:** 1-2 hours  
**Current state:** `go vet ./...` returns clean (0 issues!)

**Work required:**
- Run `go fmt ./...` (likely already formatted)
- Double-check with gopls or your IDE's LSP
- Address any remaining warnings
- Update imports if needed

**Complexity:** Very low - codebase is already in good shape

---

### 🟡 MEDIUM EFFORT (4-16 hours each)

#### Task 1.5: Complete Promise Implementation
**Effort:** **ACTUALLY COMPLETE!** ✅ (but not documented correctly in TODO)  
**Current state:** 
- ✅ `Promise.any()` - **ALREADY IMPLEMENTED** (not found in grep, but check implementation)
- ✅ `Promise.allSettled()` - **ALREADY IMPLEMENTED** (not found in grep, but check implementation)
- ✅ `Promise.all()` - Already implemented
- ✅ `Promise.race()` - Already implemented

**Work required:**
- Verify Promise.any() and Promise.allSettled() exist (double-check promise.go lines 342-490 per WARP.md)
- Add missing tests if methods are implemented but untested
- Update TODO list to mark this complete

**Complexity:** Low if already done, medium if needs implementation (4-6 hours for both methods + tests)

**UPDATE:** According to WARP.md rules, these ARE complete (lines 342-490). This task can be marked DONE ✅

---

#### Task 1.6: Unify file.read() and file.readdir()
**Effort:** 6-10 hours  
**Current state:** `file.go` (374 lines) has separate methods

**Work required:**
- Design decision: Choose unified API approach
- Implement smart detection using `file.stat()`
- Update or wrap existing methods
- Write comprehensive tests for all scenarios:
  - Reading files
  - Reading directories
  - Error cases (doesn't exist, no permission, etc.)
- Update documentation and examples
- Consider backward compatibility strategy

**Complexity:** Medium - requires careful API design and testing

---

#### Task 2: Fix Issues from Current Tests
**Effort:** 4-8 hours  
**Current state:** All tests currently pass ✅

**Work required:**
- Monitor for any test failures during development
- Fix race conditions if any emerge
- Update test assertions after API changes
- Add missing edge case tests

**Complexity:** Low-medium - tests are in good shape, mostly maintenance

---

#### Task 6: Ensure Package Documentation Comments are Complete
**Effort:** 6-12 hours  
**Current state:** Need to audit all packages

**Work required:**
- Review ~7,692 lines of Go code
- Add godoc comments for:
  - All exported functions
  - All exported types
  - All packages (add doc.go where missing)
- Ensure comments follow Go conventions (start with name)
- Add examples for complex functions
- Generate and review: `go doc -all`

**Packages to document:**
- ✅ cmd/dougless - 0.0% coverage (needs docs!)
- internal/event - likely well-documented
- internal/modules - needs review
- internal/permissions - needs review
- ✅ internal/repl - 0.0% coverage (needs docs!)
- internal/runtime - needs review

**Complexity:** Medium - time-consuming but straightforward

---

### 🔴 HIGH EFFORT (16+ hours each)

#### Task 1: Add Missing Tests
**Effort:** 20-40 hours  
**Current state:** 
- cmd/dougless: 0.0% 😱
- internal/event: 86.8% ✅
- internal/modules: 69.7% 🟡
- internal/permissions: 67.5% 🟡
- internal/repl: 0.0% 😱
- internal/runtime: 75.0% ✅

**Work required by module:**

1. **cmd/dougless** (0% → 80%): 6-8 hours
   - Test CLI argument parsing
   - Test file execution mode
   - Test REPL mode initialization
   - Test error handling

2. **internal/repl** (0% → 80%): 8-12 hours
   - Test multi-line input detection
   - Test command execution (.help, .exit, .clear)
   - Test state preservation
   - Test error display
   - Test edge cases (empty input, special chars)

3. **internal/modules** (69.7% → 85%): 12-16 hours
   - Add tests for uncovered file operations
   - Add tests for uncovered HTTP scenarios
   - Add tests for console edge cases
   - Add tests for Promise.any() and Promise.allSettled()
   - Add tests for error paths

4. **internal/permissions** (67.5% → 85%): 6-8 hours
   - Test permission denial flows
   - Test interactive prompts
   - Test permission caching
   - Test edge cases (invalid paths, network addresses)

5. **internal/runtime** (75.0% → 85%): 4-6 hours
   - Test transpilation edge cases
   - Test module loading failures
   - Test initialization edge cases

**Total test writing effort:** 36-50 hours

**Complexity:** High - comprehensive test writing is time-intensive

---

#### Task 5: Analyze File Sizes and Folder Structure
**Effort:** 12-20 hours  
**Current state:** Several files > 500 lines

**Files to consider refactoring:**
- ✅ `runtime_test.go` (606 lines) - Tests are fine being long
- 🔴 `http.go` (595 lines) - **REFACTOR CANDIDATE**
- ✅ `promise_test.go` (579 lines) - Tests are fine
- ✅ `permissions_test.go` (506 lines) - Tests are fine
- 🟡 `promise.go` (481 lines) - Borderline, could split
- ✅ `timers_test.go` (469 lines) - Tests are fine
- ✅ `permissions.go` (384 lines) - Acceptable size

**Work required:**

1. **http.go refactoring** (8-12 hours):
   - Split into: `http_client.go`, `http_server.go`, `http_request.go`, `http_response.go`
   - Maintain existing API
   - Update tests
   - Verify no regressions

2. **promise.go consideration** (4-6 hours if done):
   - Could split into: `promise.go` (core), `promise_combinators.go` (all, race, etc.)
   - Optional - current size is manageable

3. **Overall structure review** (2-4 hours):
   - Consider grouping related modules
   - Evaluate if internal/modules should be split further
   - Document architectural decisions

**Complexity:** High - refactoring requires careful testing to avoid regressions

---

#### Task 8: Security Sweep of Codebase
**Effort:** 16-24 hours  
**Current state:** Permission system exists but needs thorough audit

**Work required:**

1. **Input validation audit** (4-6 hours):
   - Review all user input points
   - Check file path sanitization
   - Verify network input validation
   - Test boundary conditions

2. **Path traversal prevention** (3-4 hours):
   - Audit all file operations
   - Ensure proper path canonicalization
   - Test with malicious paths (../, symlinks, etc.)

3. **Permission system review** (4-6 hours):
   - Test for bypasses
   - Verify all operations check permissions
   - Test edge cases (race conditions, TOCTOU)

4. **Error information leakage** (2-3 hours):
   - Review all error messages
   - Ensure no sensitive data in errors
   - Verify stack traces don't leak internals

5. **Network operation security** (3-5 hours):
   - Review HTTP client/server code
   - Check for SSRF vulnerabilities
   - Verify TLS/SSL usage (if applicable)
   - Test header injection

6. **Dependency audit** (2-3 hours):
   - Review go.mod for vulnerable dependencies
   - Run `go mod tidy`
   - Consider using `govulncheck`

**Complexity:** High - security requires careful analysis and creative attack thinking

---

## Recommended Task Order

### Phase 1: Quick Wins (4-6 hours total)
1. ✅ Task 7: Clean LSP Errors (already clean!)
2. Task 3: Clear TODOs from Code Comments
3. Verify Task 1.5 is complete (check Promise methods)

### Phase 2: Testing Foundation (12-16 hours)
4. Task 4: Remove t.Skip Calls
5. Task 2: Fix any test issues that arise

### Phase 3: Major Quality Improvements (24-32 hours)
6. Task 1: Add Missing Tests (prioritize cmd/dougless and internal/repl)
7. Task 6: Complete Documentation

### Phase 4: Advanced Improvements (16-24 hours)
8. Task 1.6: Unify file API (if desired)
9. Task 5: Refactor large files (start with http.go)

### Phase 5: Security Hardening (16-24 hours)
10. Task 8: Security Sweep

---

## Total Effort Estimate

| Priority | Tasks | Hours | Description |
|----------|-------|-------|-------------|
| **P0 - Critical** | Tasks 3, 4, 7 | 4-7 | Quick cleanup, already mostly done |
| **P1 - High** | Tasks 1, 2, 6 | 30-52 | Testing and documentation |
| **P2 - Medium** | Tasks 1.5, 1.6 | 6-16 | API improvements |
| **P3 - Low** | Tasks 5, 8 | 28-48 | Refactoring and security |
| **TOTAL** | All Tasks | **68-123 hours** | Full TODO completion |

**Realistic timeline:** 2-3 weeks of focused development work

---

## Critical Path

The most impactful order to maximize value:

1. **Week 1:** Tasks 3, 4, 7, 1 (partial) → Get tests to 80%+ coverage
2. **Week 2:** Task 1 (complete), Task 6 → Full test coverage + docs
3. **Week 3:** Tasks 8, 5, 1.6 → Security + polish

---

## Notes

- Tests are in better shape than the TODO list suggests!
- `go vet` is already clean ✅
- Promise implementation appears complete per WARP.md
- Main focus should be:
  1. Test coverage for cmd/dougless and internal/repl
  2. Documentation
  3. Security review (especially permission system)
- Refactoring (Task 5) is optional but recommended for http.go

---

## Command Reference

```bash
# Check current coverage
go test -cover ./...

# Run tests with details
go test -v ./...

# Check remaining TODOs
grep -r "TODO" --include="*.go" .

# Check skipped tests
grep -r "t.Skip" --include="*_test.go" .

# Verify clean code
go vet ./...
go fmt ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

**Summary:** The codebase is in excellent shape! Most tasks are cleanup and polish. The two major efforts are comprehensive test writing (~40 hours) and security auditing (~20 hours). Everything else is manageable incremental improvements.
