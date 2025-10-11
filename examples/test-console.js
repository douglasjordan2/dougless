// Test console enhancements
console.log('=== Console Enhancement Tests ===\n');

// Test 1: console.time/timeEnd with default label
console.log('Test 1: Timer with default label');
console.time();
let sum = 0;
for (let i = 0; i < 1000000; i++) {
  sum += i;
}
console.timeEnd();

// Test 2: console.time/timeEnd with custom label
console.log('\nTest 2: Timer with custom label');
console.time('calculation');
let product = 1;
for (let i = 1; i < 100; i++) {
  product = (product * i) % 1000000;
}
console.timeEnd('calculation');

// Test 3: Multiple concurrent timers
console.log('\nTest 3: Multiple concurrent timers');
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

// Test 4: console.table with array
setTimeout(function() {
  console.log('\nTest 4: Table with array');
  console.table([10, 20, 30, 40, 50]);
}, 200);

// Test 5: console.table with strings
setTimeout(function() {
  console.log('\nTest 5: Table with strings');
  console.table(['apple', 'banana', 'cherry', 'date', 'elderberry']);
}, 250);

// Test 6: console.table with object
setTimeout(function() {
  console.log('\nTest 6: Table with object');
  console.table({
    name: 'Dougless Runtime',
    version: '0.1.0',
    language: 'Go',
    engine: 'Goja',
    status: 'In Development'
  });
}, 300);

// Test 7: console.table with empty data
setTimeout(function() {
  console.log('\nTest 7: Table with empty array (should print nothing)');
  console.table([]);
}, 350);

// Test 8: Calling timeEnd without starting timer
setTimeout(function() {
  console.log('\nTest 8: TimeEnd without time (should show warning)');
  console.timeEnd('nonexistent');
}, 400);

// Test 9: All tests complete
setTimeout(function() {
  console.log('\n=== All Console Tests Complete ===');
}, 450);
