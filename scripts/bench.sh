#!/bin/bash
#
# Benchmark runner for Dougless Runtime
# Runs comprehensive benchmarks and saves results for comparison

set -e

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RESULTS_DIR="bench_results"
RESULT_FILE="$RESULTS_DIR/bench_$TIMESTAMP.txt"

# Create results directory if it doesn't exist
mkdir -p "$RESULTS_DIR"

echo "🚀 Running Dougless Runtime Benchmarks"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Results will be saved to: $RESULT_FILE"
echo ""

# Run benchmarks with memory stats
echo "Running Runtime benchmarks..."
go test -bench=. -benchmem -benchtime=1s ./internal/runtime | tee -a "$RESULT_FILE"
echo ""

echo "Running Event Loop benchmarks..."
go test -bench=. -benchmem -benchtime=1s ./internal/event | tee -a "$RESULT_FILE"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Benchmarks complete!"
echo ""
echo "Results saved to: $RESULT_FILE"
echo ""

# Check if this is the first benchmark
BASELINE="$RESULTS_DIR/baseline.txt"
if [ ! -f "$BASELINE" ]; then
    echo "📊 This is your first benchmark run!"
    echo "Saving as baseline for future comparisons..."
    cp "$RESULT_FILE" "$BASELINE"
    echo "Baseline saved to: $BASELINE"
else
    echo "📊 Comparing with baseline..."
    echo ""
    
    # Check if benchstat is installed
    if command -v benchstat &> /dev/null; then
        benchstat "$BASELINE" "$RESULT_FILE"
    else
        echo "⚠️  benchstat not installed. Install it for statistical comparison:"
        echo "   go install golang.org/x/perf/cmd/benchstat@latest"
        echo ""
        echo "For now, here's a simple comparison:"
        echo ""
        echo "BASELINE (previous):"
        echo "-------------------"
        head -20 "$BASELINE"
        echo ""
        echo "CURRENT (now):"
        echo "-------------"
        head -20 "$RESULT_FILE"
    fi
fi

echo ""
echo "💡 Tips:"
echo "  - Run 'cat $RESULT_FILE' to see full results"
echo "  - Run 'benchstat baseline.txt $RESULT_FILE' to compare"
echo "  - Run 'scripts/bench.sh' again to track progress"
echo "  - See docs/benchmarking.md for detailed guidance"
