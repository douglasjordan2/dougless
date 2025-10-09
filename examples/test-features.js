console.log('Testing Dougless Runtime Features...\n');

// Test console.time
console.time('script-execution');

// Test setTimeout
console.log('1. Setting up a timeout...');
setTimeout(function() {
    console.log('   ✓ Timeout executed after 1 second');
}, 1000);

// Test setInterval
console.log('2. Setting up an interval...');
let count = 0;
const interval = setInterval(function() {
    count++;
    console.log('   ✓ Interval tick:', count);
    if (count === 3) {
        clearInterval(interval);
        console.log('   ✓ Interval cleared\n');
        console.timeEnd('script-execution');
    }
}, 500);

console.log('3. Event loop is running...\n');
