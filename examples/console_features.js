// ====================================
// Dougless Runtime - Console Features
// ====================================
// Demonstrates all console operations available in Dougless

console.log('=== Console Features Demo ===\n');

// 1. Basic logging with multiple arguments
console.log('1. Basic Logging');
console.log('   Simple message');
console.log('   Multiple', 'arguments', 'work', 'too!');
console.log('   Numbers:', 42, 'and strings:', 'hello');
console.log('');

// 2. Warning and error messages
console.log('2. Warning and Error Messages');
console.warn('   This is a warning message');
console.error('   This is an error message');
console.log('');

// 3. Performance timing - default label
console.log('3. Performance Timing (Default Label)');
console.time();
let sum = 0;
for (let i = 0; i < 1000000; i++) {
    sum += i;
}
console.timeEnd();
console.log('   Sum calculated:', sum);
console.log('');

// 4. Performance timing - custom label
console.log('4. Performance Timing (Custom Label)');
console.time('fibonacci-calculation');
function fibonacci(n) {
    if (n <= 1) return n;
    return fibonacci(n - 1) + fibonacci(n - 2);
}
const fib15 = fibonacci(15);
console.timeEnd('fibonacci-calculation');
console.log('   Fibonacci(15) =', fib15);
console.log('');

// 5. Multiple concurrent timers
console.log('5. Multiple Concurrent Timers');
console.time('timer1');
console.time('timer2');
console.time('timer3');

setTimeout(function() {
    console.timeEnd('timer1');
}, 50);

setTimeout(function() {
    console.timeEnd('timer2');
}, 100);

setTimeout(function() {
    console.timeEnd('timer3');
}, 150);

// 6. console.table with arrays
setTimeout(function() {
    console.log('\n6. Table Display - Array of Numbers');
    console.table([10, 20, 30, 40, 50]);
}, 200);

// 7. console.table with array of strings
setTimeout(function() {
    console.log('\n7. Table Display - Array of Strings');
    console.table(['apple', 'banana', 'cherry', 'date', 'elderberry']);
}, 250);

// 8. console.table with objects
setTimeout(function() {
    console.log('\n8. Table Display - Object Properties');
    console.table({
        name: 'Dougless Runtime',
        version: '0.1.0',
        language: 'Go + JavaScript',
        engine: 'Goja (ES5.1)',
        status: 'Active Development',
        phase: 'Phase 3 Complete'
    });
}, 300);

// 9. Handling timer that doesn't exist
setTimeout(function() {
    console.log('\n9. Edge Case - Timer Not Started');
    console.timeEnd('nonexistent-timer');
}, 350);

// 10. Empty table
setTimeout(function() {
    console.log('\n10. Edge Case - Empty Array Table');
    console.table([]);
    console.log('   (Empty table produces no output)');
}, 400);

// Final message
setTimeout(function() {
    console.log('\n=== Console Demo Complete ===');
    console.log('All console features demonstrated successfully!');
}, 450);
