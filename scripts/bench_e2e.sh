#!/bin/bash
#
# End-to-End Benchmarks for Dougless Runtime
# Measures realistic user-facing performance
#
# Tests:
# 1. Actual startup time (process creation to completion)
# 2. HTTP server throughput (requests per second)
# 3. File I/O performance (read/write operations)

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "🚀 Dougless Runtime - End-to-End Benchmarks"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Build the runtime if needed
if [ ! -f "$DOUGLESS_BIN" ]; then
    echo "📦 Building dougless runtime..."
    (cd "$PROJECT_ROOT" && go build -o dougless cmd/dougless/main.go)
    echo "✅ Build complete"
    echo ""
fi

# Get absolute path to project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DOUGLESS_BIN="$PROJECT_ROOT/dougless"

# Create temp directory for benchmark scripts
BENCH_DIR=$(mktemp -d)
trap "rm -rf $BENCH_DIR" EXIT

# Create permissive .douglessrc for benchmarks
cat > "$BENCH_DIR/.douglessrc" << EOF
{
  "permissions": {
    "read": ["*", "/tmp/*", "$BENCH_DIR/*"],
    "write": ["*", "/tmp/*", "$BENCH_DIR/*"],
    "net": ["*"],
    "env": ["*"],
    "run": ["*"]
  }
}
EOF

echo "📊 Running benchmarks..."
echo ""

# =============================================================================
# 1. STARTUP TIME BENCHMARK
# =============================================================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "1️⃣  STARTUP TIME (Goal: < 100ms)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Create simple hello world script
cat > "$BENCH_DIR/hello.js" << 'EOF'
console.log("Hello, World!");
EOF

echo "Testing: Simple 'Hello World' script"
echo "Iterations: 100"
echo ""

total_time=0
for i in {1..100}; do
    start=$(date +%s%N)
    (cd "$BENCH_DIR" && "$DOUGLESS_BIN" hello.js) > /dev/null 2>&1
    end=$(date +%s%N)
    elapsed=$((end - start))
    total_time=$((total_time + elapsed))
done

avg_time=$((total_time / 100))
avg_ms=$(echo "scale=2; $avg_time / 1000000" | bc)

echo -e "Average startup time: ${GREEN}${avg_ms}ms${NC}"

if (( $(echo "$avg_ms < 100" | bc -l) )); then
    echo -e "Status: ${GREEN}✅ PASS${NC} (under 100ms goal)"
else
    echo -e "Status: ${RED}❌ FAIL${NC} (over 100ms goal)"
fi
echo ""

# =============================================================================
# 2. HTTP SERVER THROUGHPUT BENCHMARK
# =============================================================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "2️⃣  HTTP THROUGHPUT (Goal: > 10,000 req/sec)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if we have a benchmarking tool
if ! command -v ab &> /dev/null && ! command -v wrk &> /dev/null; then
    echo -e "${YELLOW}⚠️  No HTTP benchmarking tool found (ab or wrk)${NC}"
    echo "   Install one to test HTTP throughput:"
    echo "   - Apache Bench: sudo apt install apache2-utils"
    echo "   - wrk: sudo apt install wrk"
    echo ""
    echo "Skipping HTTP benchmark..."
    echo ""
else
    # Create HTTP server script
    cat > "$BENCH_DIR/server.js" << 'EOF'
const server = http.createServer((req, res) => {
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('Hello, World!');
});

server.listen(3456, () => {
    console.log('Server ready');
});
EOF

    echo "Starting HTTP server on port 3456..."
    (cd "$BENCH_DIR" && "$DOUGLESS_BIN" --allow-all server.js) > /dev/null 2>&1 &
    SERVER_PID=$!
    
    # Wait for server to start
    sleep 3
    
    if command -v wrk &> /dev/null; then
        echo "Using wrk for benchmarking..."
        echo "Test: 10 seconds, 2 threads, 100 connections"
        echo ""
        
        wrk -t2 -c100 -d10s http://localhost:3456 > "$BENCH_DIR/wrk_results.txt" 2>&1
        
        req_per_sec=$(grep "Requests/sec:" "$BENCH_DIR/wrk_results.txt" | awk '{print $2}')
        echo "Results:"
        cat "$BENCH_DIR/wrk_results.txt"
        echo ""
        echo -e "Throughput: ${GREEN}${req_per_sec} req/sec${NC}"
        
        if (( $(echo "$req_per_sec > 10000" | bc -l) )); then
            echo -e "Status: ${GREEN}✅ PASS${NC} (over 10k req/sec goal)"
        else
            echo -e "Status: ${YELLOW}⚠️  BELOW GOAL${NC} (under 10k req/sec)"
        fi
    else
        echo "Using Apache Bench for benchmarking..."
        echo "Test: 10,000 requests, 100 concurrent"
        echo ""
        
        ab -n 10000 -c 100 -q http://localhost:3456/ > "$BENCH_DIR/ab_results.txt" 2>&1
        
        req_per_sec=$(grep "Requests per second:" "$BENCH_DIR/ab_results.txt" | awk '{print $4}')
        echo "Results:"
        grep "Requests per second:" "$BENCH_DIR/ab_results.txt"
        grep "Time per request:" "$BENCH_DIR/ab_results.txt"
        grep "Transfer rate:" "$BENCH_DIR/ab_results.txt"
        echo ""
        echo -e "Throughput: ${GREEN}${req_per_sec} req/sec${NC}"
        
        if (( $(echo "$req_per_sec > 10000" | bc -l) )); then
            echo -e "Status: ${GREEN}✅ PASS${NC} (over 10k req/sec goal)"
        else
            echo -e "Status: ${YELLOW}⚠️  BELOW GOAL${NC} (under 10k req/sec)"
        fi
    fi
    
    # Cleanup
    kill $SERVER_PID 2>/dev/null || true
    wait $SERVER_PID 2>/dev/null || true
fi

echo ""

# =============================================================================
# 3. FILE I/O BENCHMARK
# =============================================================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "3️⃣  FILE I/O (Goal: Comparable to Node.js)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Create file I/O benchmark script (using promises for accurate timing)
cat > "$BENCH_DIR/file_bench.js" << 'EOF'
const iterations = 100;
const testFile = '/tmp/dougless_bench_test.txt';
const testData = 'Hello, Dougless! '.repeat(100); // ~1.7KB

async function benchWrite() {
    const start = Date.now();
    for (let i = 0; i < iterations; i++) {
        await files.write(testFile, testData);
    }
    const end = Date.now();
    console.log('write_ops: ' + (end - start) + 'ms for ' + iterations + ' operations');
}

async function benchRead() {
    const start = Date.now();
    for (let i = 0; i < iterations; i++) {
        await files.read(testFile);
    }
    const end = Date.now();
    console.log('read_ops: ' + (end - start) + 'ms for ' + iterations + ' operations');
}

async function main() {
    await benchWrite();
    await benchRead();
    await files.rm(testFile);
    console.log('File I/O benchmark complete');
}

main();
EOF

echo "Testing: 100 write + 100 read operations (~1.7KB each)"
echo ""

(cd "$BENCH_DIR" && "$DOUGLESS_BIN" --allow-all file_bench.js) 2>&1 | tee "$BENCH_DIR/file_results.txt"

write_time=$(grep "write_ops:" "$BENCH_DIR/file_results.txt" | cut -d' ' -f2)
read_time=$(grep "read_ops:" "$BENCH_DIR/file_results.txt" | cut -d' ' -f2)

echo ""
echo -e "Write performance: ${GREEN}${write_time}${NC}"
echo -e "Read performance: ${GREEN}${read_time}${NC}"
echo ""

# =============================================================================
# MEMORY USAGE BENCHMARK
# =============================================================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "4️⃣  MEMORY USAGE (Goal: < 50MB for typical apps)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Create typical application
cat > "$BENCH_DIR/typical_app.js" << 'EOF'
// Typical application: HTTP server + file operations + timers

const data = [];
for (let i = 0; i < 1000; i++) {
    data.push({ id: i, name: `Item ${i}` });
}

const server = http.createServer((req, res) => {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ message: 'Hello', data: data.slice(0, 10) }));
});

server.listen(3457, () => {
    console.log('Memory test server ready');
});

// Keep alive for measurement
setTimeout(() => {}, 5000);
EOF

echo "Starting typical application..."
(cd "$BENCH_DIR" && "$DOUGLESS_BIN" --allow-all typical_app.js) > "$BENCH_DIR/app_output.txt" 2>&1 &
APP_PID=$!

# Wait for server to start
sleep 3

# Get memory usage
if command -v ps &> /dev/null; then
    if ps -p $APP_PID > /dev/null 2>&1; then
        mem_kb=$(ps -o rss= -p $APP_PID 2>/dev/null || echo "0")
        mem_mb=$(echo "scale=2; $mem_kb / 1024" | bc)
        
        echo -e "Memory usage: ${GREEN}${mem_mb}MB${NC}"
        
        if (( $(echo "$mem_mb < 50" | bc -l) )); then
            echo -e "Status: ${GREEN}✅ PASS${NC} (under 50MB goal)"
        else
            echo -e "Status: ${YELLOW}⚠️  OVER GOAL${NC} (over 50MB)"
        fi
    else
        echo -e "${RED}⚠️  Process died - check $BENCH_DIR/app_output.txt${NC}"
        cat "$BENCH_DIR/app_output.txt"
    fi
else
    echo -e "${YELLOW}⚠️  ps command not available, cannot measure memory${NC}"
fi

# Cleanup
kill $APP_PID 2>/dev/null || true
wait $APP_PID 2>/dev/null || true

echo ""

# =============================================================================
# SUMMARY
# =============================================================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 SUMMARY"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Performance Targets vs Actual:"
echo ""
echo "  Startup Time:   < 100ms     →  ${avg_ms}ms"
echo "  Memory Usage:   < 50MB       →  ${mem_mb}MB"
if [ -n "$req_per_sec" ]; then
    echo "  HTTP Throughput: > 10k req/s →  ${req_per_sec} req/s"
fi
echo ""
echo "For detailed micro-benchmarks, run: ./scripts/bench.sh"
echo "For profiling, see: docs/benchmarking.md"
echo ""
