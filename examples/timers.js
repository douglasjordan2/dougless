// =================================
// Dougless Runtime - Timer Features
// =================================
// Demonstrates setTimeout, setInterval, and timer management

console.log('=== Timer Features Demo ===\n');

// 1. Basic setTimeout - immediate (0ms delay)
console.log('1. setTimeout with 0ms delay');
setTimeout(function() {
    console.log('   ✓ Immediate timeout executed');
}, 0);

// 2. setTimeout with delay
console.log('2. setTimeout with 100ms delay');
setTimeout(function() {
    console.log('   ✓ 100ms timeout executed');
}, 100);

// 3. setTimeout with longer delay
console.log('3. setTimeout with 500ms delay');
setTimeout(function() {
    console.log('   ✓ 500ms timeout executed');
}, 500);

// 4. Multiple timeouts (testing execution order)
console.log('4. Multiple timeouts - order test');
setTimeout(function() {
    console.log('   ✓ Should execute THIRD (300ms)');
}, 300);

setTimeout(function() {
    console.log('   ✓ Should execute FIRST (50ms)');
}, 50);

setTimeout(function() {
    console.log('   ✓ Should execute SECOND (150ms)');
}, 150);

// 5. setTimeout without delay argument (defaults to 0)
console.log('5. setTimeout without delay argument');
setTimeout(function() {
    console.log('   ✓ No-delay timeout executed');
});

// 6. Nested setTimeout
console.log('6. Nested setTimeout');
setTimeout(function() {
    console.log('   ✓ Outer timeout executed (200ms)');
    setTimeout(function() {
        console.log('   ✓ Inner timeout executed (100ms after outer)');
    }, 100);
}, 200);

// 7. Basic setInterval
console.log('\n7. Basic setInterval - 3 ticks at 400ms intervals');
let intervalCount = 0;
const basicInterval = setInterval(function() {
    intervalCount++;
    console.log('   ✓ Interval tick:', intervalCount);
    
    if (intervalCount >= 3) {
        clearInterval(basicInterval);
        console.log('   ✓ Interval cleared after 3 ticks');
    }
}, 400);

// 8. Multiple intervals running simultaneously
setTimeout(function() {
    console.log('\n8. Multiple intervals simultaneously');
    
    let count1 = 0;
    let count2 = 0;
    
    const interval1 = setInterval(function() {
        count1++;
        console.log('   Interval 1 tick:', count1);
        if (count1 >= 3) {
            clearInterval(interval1);
            console.log('   Interval 1 stopped');
        }
    }, 350);
    
    const interval2 = setInterval(function() {
        count2++;
        console.log('   Interval 2 tick:', count2);
        if (count2 >= 3) {
            clearInterval(interval2);
            console.log('   Interval 2 stopped');
        }
    }, 550);
}, 2000);

// 9. clearTimeout - canceling a timer before it executes
setTimeout(function() {
    console.log('\n9. clearTimeout - canceling before execution');
    
    const timer = setTimeout(function() {
        console.log('   ✗ This should never print!');
    }, 5000);
    
    clearTimeout(timer);
    console.log('   ✓ Timer canceled successfully');
}, 4500);

// 10. Double-clear safety test (should not crash)
setTimeout(function() {
    console.log('\n10. Edge case - double-clearing a timer');
    
    const timer = setTimeout(function() {
        console.log('   ✗ This should never execute');
    }, 5000);
    
    clearTimeout(timer);
    console.log('   ✓ First clear successful');
    
    clearTimeout(timer);
    console.log('   ✓ Second clear (no crash!)');
    
    // Try clearing a fake ID
    clearTimeout('fake-timer-id-123');
    console.log('   ✓ Clearing fake ID (no crash!)');
}, 5000);

// 11. Timer that completes normally (not canceled)
setTimeout(function() {
    console.log('\n11. Timer completing normally (not canceled)');
    setTimeout(function() {
        console.log('   ✓ Timer completed successfully');
    }, 100);
}, 5500);

// Final completion message
setTimeout(function() {
    console.log('\n=== Timer Demo Complete ===');
    console.log('All timer features demonstrated successfully!');
}, 6000);
