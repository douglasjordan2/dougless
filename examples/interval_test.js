// Test setInterval functionality
console.log('Starting interval test...');

let count = 0;

// Test 1: Basic setInterval
const intervalId = setInterval(function() {
  count++;
  console.log('Interval tick:', count);
  
  // Stop after 5 iterations
  if (count >= 5) {
    console.log('Stopping interval after 5 ticks');
    clearInterval(intervalId);
  }
}, 500); // Run every 500ms

// Test 2: setTimeout to ensure it still works
setTimeout(function() {
  console.log('This should run after 1 second');
}, 1000);

console.log('Timers scheduled, waiting for execution...');
