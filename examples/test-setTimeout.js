// Test script for setTimeout implementation
console.log("=== setTimeout Test Suite ===");
console.log("Starting tests...\n");

// Test 1: Basic setTimeout with 0ms delay
console.log("Test 1: setTimeout with 0ms delay");
setTimeout(function() {
    console.log("✓ Test 1 passed: Immediate timeout executed");
}, 0);

// Test 2: setTimeout with 100ms delay
console.log("Test 2: setTimeout with 100ms delay");
setTimeout(function() {
    console.log("✓ Test 2 passed: 100ms timeout executed");
}, 100);

// Test 3: setTimeout with 500ms delay
console.log("Test 3: setTimeout with 500ms delay");
setTimeout(function() {
    console.log("✓ Test 3 passed: 500ms timeout executed");
}, 500);

// Test 4: Multiple setTimeout calls (testing order)
console.log("Test 4: Multiple timeouts (order test)");
setTimeout(function() {
    console.log("✓ Should execute THIRD (300ms)");
}, 300);

setTimeout(function() {
    console.log("✓ Should execute FIRST (50ms)");
}, 50);

setTimeout(function() {
    console.log("✓ Should execute SECOND (150ms)");
}, 150);

// Test 5: setTimeout without delay argument (should default to 0)
console.log("Test 5: setTimeout without delay argument");
setTimeout(function() {
    console.log("✓ Test 5 passed: No-delay timeout executed");
});

// Test 6: Nested setTimeout
console.log("Test 6: Nested setTimeout");
setTimeout(function() {
    console.log("✓ Test 6a: Outer timeout executed (200ms)");
    setTimeout(function() {
        console.log("✓ Test 6b: Inner timeout executed (100ms after outer)");
    }, 100);
}, 200);

console.log("\nAll setTimeout calls scheduled. Waiting for execution...\n");
