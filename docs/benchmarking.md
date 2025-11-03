# Benchmarking Guide

This guide explains how to run and interpret benchmarks for the Dougless runtime.

## Quick Start

```bash
# Run all benchmarks
go test -bench=. ./...

# Run benchmarks with memory stats
go test -bench=. -benchmem ./...

# Run specific benchmark suite
go test -bench=. -benchmem ./internal/runtime
go test -bench=. -benchmem ./internal/event
go test -bench=. -benchmem ./internal/modules

# Run specific benchmark
go test -bench=BenchmarkRuntimeCreation -benchmem ./internal/runtime

# Run benchmarks multiple times for accuracy
go test -bench=. -benchmem -count=5 ./...

# Save benchmark results to file
go test -bench=. -benchmem ./... > bench_results.txt

# Compare benchmark results (requires benchstat: go install golang.org/x/perf/cmd/benchstat@latest)
go test -bench=. -benchmem ./... > new.txt
benchstat old.txt new.txt
```

## Performance Targets

Based on the project goals outlined in WARP.md:

### Runtime Targets
- **Startup Time**: < 100ms for basic scripts
- **Memory Usage**: < 50MB for typical applications
- **HTTP Throughput**: > 10,000 requests/second
- **File I/O**: Comparable to Node.js performance

### Specific Benchmarks

#### Runtime (`internal/runtime/runtime_bench_test.go`)
- `BenchmarkRuntimeCreation` - Should be < 10ms
- `BenchmarkSimpleExecution` - Should be < 100ms (includes startup)
- `BenchmarkTranspilation` - Measures ES6+ → ES5 overhead
- `BenchmarkLargeScriptExecution` - Tests with compute-heavy scripts
- `BenchmarkES6Features` - Validates modern JS performance
- `BenchmarkAsyncAwait` - Measures Promise/async overhead
- `BenchmarkModuleRequire` - Tests module loading performance
- `BenchmarkConsoleLog` - Baseline console output cost
- `BenchmarkTimerCreation` - setTimeout/setInterval overhead
- `BenchmarkPromiseCreation` - Promise constructor cost

#### Event Loop (`internal/event/loop_bench_test.go`)
- `BenchmarkLoopCreation` - Event loop initialization cost
- `BenchmarkTaskScheduling` - Immediate task overhead
- `BenchmarkDelayedTaskScheduling` - Timer-based task cost
- `BenchmarkParallelTaskScheduling` - Concurrent scheduling
- `BenchmarkTimerCancellation` - clearTimeout performance
- `BenchmarkMultipleTimers` - Scales with 10/100/1000 timers
- `BenchmarkKeepAlive` - KeepAlive/Done overhead
- `BenchmarkTaskThroughput` - Maximum tasks/second

#### File I/O (`internal/modules/files_bench_test.go`)
- `BenchmarkFileRead` - Read performance at 1KB/10KB/100KB/1MB
- `BenchmarkFileWrite` - Write performance at various sizes
- `BenchmarkFileRemove` - File deletion speed
- `BenchmarkDirectoryRemoval` - Recursive removal with 10/100/1000 files
- `BenchmarkPromiseBasedRead` - Promise API overhead
- `BenchmarkConcurrentReads` - Parallel read performance

## Interpreting Results

### Benchmark Output Format
```
BenchmarkRuntimeCreation-8    1000    1234567 ns/op    12345 B/op    123 allocs/op
│                         │    │       │                │             └─ allocations per op
│                         │    │       │                └─ bytes allocated per op
│                         │    │       └─ nanoseconds per operation
│                         │    └─ number of iterations
│                         └─ GOMAXPROCS value (CPU cores used)
```

### What to Look For

1. **High ns/op** - Slow operations that need optimization
2. **High B/op** - Memory-intensive operations (potential GC pressure)
3. **High allocs/op** - Many allocations (consider pooling/reuse)
4. **Low iterations** - Test is too slow (benchmark may time out)

### Optimization Priorities

Focus optimization efforts on:
1. High-frequency operations (e.g., task scheduling, console.log)
2. Startup path operations (e.g., runtime creation, module loading)
3. Operations with high allocation counts
4. Bottlenecks identified through profiling

## Profiling

### CPU Profile
```bash
# Generate CPU profile
go test -bench=BenchmarkSimpleExecution -cpuprofile=cpu.prof ./internal/runtime

# Analyze with pprof
go tool pprof cpu.prof
# Commands: top, list, web
```

### Memory Profile
```bash
# Generate memory profile
go test -bench=BenchmarkSimpleExecution -memprofile=mem.prof ./internal/runtime

# Analyze with pprof
go tool pprof mem.prof
# Commands: top, list, web
```

### Trace Analysis
```bash
# Generate execution trace
go test -bench=BenchmarkSimpleExecution -trace=trace.out ./internal/runtime

# View trace
go tool trace trace.out
```

## Continuous Benchmarking

To track performance over time:

1. **Baseline**: Run benchmarks on main branch
   ```bash
   git checkout main
   go test -bench=. -benchmem ./... > baseline.txt
   ```

2. **Compare**: Run benchmarks on feature branch
   ```bash
   git checkout feature-branch
   go test -bench=. -benchmem ./... > feature.txt
   benchstat baseline.txt feature.txt
   ```

3. **Interpret**: Look for significant changes (> 5-10%)
   - ✅ Improvements show negative percentages (e.g., -15%)
   - ⚠️ Regressions show positive percentages (e.g., +20%)
   - ~ Small changes (< 5%) may be noise

## Benchmark Best Practices

1. **Disable CPU throttling**: Set CPU governor to performance mode
   ```bash
   # Linux
   echo performance | sudo tee /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor
   ```

2. **Close background apps**: Minimize interference
3. **Run multiple times**: Use `-count=5` or higher for statistical significance
4. **Use benchstat**: Compare results scientifically
5. **Profile before optimizing**: Don't guess, measure

## Adding New Benchmarks

When adding benchmarks:

1. Name them descriptively: `Benchmark<Operation><Variant>`
2. Use `b.ResetTimer()` to exclude setup time
3. Use `b.StopTimer()`/`b.StartTimer()` to exclude cleanup
4. Always check for errors (don't `b.Fatal()` before `b.ResetTimer()`)
5. Document what the benchmark measures
6. Consider testing multiple scales (small/medium/large inputs)

Example:
```go
func BenchmarkNewFeature(b *testing.B) {
    // Setup (excluded from timing)
    setup := createTestSetup()
    
    b.ResetTimer() // Start timing here
    for i := 0; i < b.N; i++ {
        result := doOperation(setup)
        if result == nil {
            b.Fatal("operation failed")
        }
    }
}
```

## Known Performance Characteristics

Based on current implementation:

- **Goja VM Creation**: ~5-10ms (relatively expensive)
- **Transpilation**: ~10-50ms depending on code size (uses esbuild)
- **Event Loop Overhead**: ~1-2µs per task
- **Module Caching**: First require() is expensive, subsequent are fast
- **File I/O**: Async overhead adds ~100µs per operation

## Future Benchmark Additions

Planned benchmarks for Phase 9:

- HTTP client request throughput
- WebSocket message handling rate
- Crypto operations (hash/hmac/random)
- Process API overhead
- Memory leak detection tests
- Long-running stability tests
