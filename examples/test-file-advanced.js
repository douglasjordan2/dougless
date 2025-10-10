console.log('=== Testing Advanced File Operations ===\n');

// Test 1: Create a directory
console.log('1. Creating directory "test-dir"...');
file.mkdir('test-dir', function(err) {
    if (err) {
        console.error('   ✗ Error:', err);
        return;
    }
    console.log('   ✓ Directory created\n');
    
    // Test 2: Write a file in the directory
    console.log('2. Writing file "test-dir/data.txt"...');
    file.write('test-dir/data.txt', 'Hello Dougless!', function(err) {
        if (err) {
            console.error('   ✗ Error:', err);
            return;
        }
        console.log('   ✓ File written\n');
        
        // Test 3: Get file stats
        console.log('3. Getting file stats...');
        file.stat('test-dir/data.txt', function(err, stats) {
            if (err) {
                console.error('   ✗ Error:', err);
            } else {
                console.log('   ✓ Size:', stats.size, 'bytes');
                console.log('   ✓ Is file:', stats.isFile);
                console.log('   ✓ Is directory:', stats.isDirectory);
                console.log('   ✓ Name:', stats.name);
            }
            console.log('');
            
            // Test 4: List directory contents
            console.log('4. Listing directory contents...');
            file.readdir('test-dir', function(err, files) {
                if (err) {
                    console.error('   ✗ Error:', err);
                } else {
                    console.log('   ✓ Files:', files);
                }
                console.log('');
                
                // Test 5: Delete the file
                console.log('5. Deleting file...');
                file.unlink('test-dir/data.txt', function(err) {
                    if (err) {
                        console.error('   ✗ Error:', err);
                    } else {
                        console.log('   ✓ File deleted\n');
                        
                        // Test 6: Remove the directory
                        console.log('6. Removing directory...');
                        file.rmdir('test-dir', function(err) {
                            if (err) {
                                console.error('   ✗ Error:', err);
                            } else {
                                console.log('   ✓ Directory removed\n');
                                console.log('✅ All tests complete!');
                            }
                        });
                    }
                });
            });
        });
    });
});
