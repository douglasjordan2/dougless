console.log('=== Testing Dougless File System ===\n');

// Test 1: Write a file
console.log('1. Writing file...');
file.write('test-output.txt', 'Hello from Dougless!', function(err) {
    if (err) {
        console.error('Write error:', err);
    } else {
        console.log('   ✓ File written\n');
        
        // Test 2: Read it back
        console.log('2. Reading file...');
        file.read('test-output.txt', function(err, data) {
            if (err) {
                console.error('Read error:', err);
            } else {
                console.log('   ✓ Content:', data);
            }
        });
    }
});
