// Comprehensive timer edge case tests
console.log('=== Timer Edge Case Tests ===\n');

// Test 1: Multiple intervals
console.log('Test 1: Multiple intervals running simultaneously');
let count1 = 0;
let count2 = 0;

const interval1 = setInterval(function() {
  count1++;
  console.log('  Interval 1 tick:', count1);
  if (count1 >= 3) {
    clearInterval(interval1);
    console.log('  Interval 1 stopped');
  }
}, 300);

const interval2 = setInterval(function() {
  count2++;
  console.log('  Interval 2 tick:', count2);
  if (count2 >= 3) {
    clearInterval(interval2);
    console.log('  Interval 2 stopped');
  }
}, 500);

// Test 2: Clearing a timer multiple times (should not crash)
setTimeout(function() {
  console.log('\nTest 2: Double-clear test');
  const timer = setTimeout(function() {
    console.log('  This should never print');
  }, 5000);
  
  clearTimeout(timer);
  console.log('  Timer cleared once');
  
  clearTimeout(timer);
  console.log('  Timer cleared again (should not crash)');
  
  // Try to clear non-existent timer
  clearTimeout('fake-timer-id-123');
  console.log('  Cleared fake timer ID (should not crash)');
}, 1000);

// Test 3: setTimeout that completes before being cleared
setTimeout(function() {
  console.log('\nTest 3: Timeout completes normally');
}, 1500);

// Test 4: Mix of setTimeout and setInterval
setTimeout(function() {
  console.log('\nTest 4: All timers completed successfully!');
}, 2000);

console.log('\nAll tests scheduled...\n');
