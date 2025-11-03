// Test setTimeout and setInterval
console.log('Testing timers...');

let count = 0;

// Test setTimeout
setTimeout(() => {
  console.log('setTimeout: This runs after 100ms');
}, 100);

// Test setInterval
const intervalId = setInterval(() => {
  count++;
  console.log(`setInterval: Count is ${count}`);
  
  if (count === 3) {
    clearInterval(intervalId);
    console.log('setInterval: Cleared after 3 iterations');
  }
}, 50);

// Test clearTimeout
const timeoutId = setTimeout(() => {
  console.log('This should NOT print');
}, 200);
clearTimeout(timeoutId);
console.log('Cleared timeout before it could run');

console.log('Timers set, waiting for completion...');
