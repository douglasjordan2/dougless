// files_basic.js - Basic demonstrations of the simplified files API

console.log('=== Dougless Files API - Basic Examples ===\n');

// Example 1: Read a file
console.log('1. Reading a file:');
files.read('examples/hello.js', function(err, content) {
    if (err) {
        console.error('  Error:', err);
        return;
    }
    
    if (content === null) {
        console.log('  File does not exist');
    } else {
        console.log('  File length:', content.length, 'bytes');
        console.log('  First 50 chars:', content.substring(0, 50));
    }
    console.log('');
    
    // Example 2: Read a directory (note the trailing slash)
    console.log('2. Listing directory contents:');
    files.read('examples/', function(err, fileNames) {
        if (err) {
            console.error('  Error:', err);
            return;
        }
        
        console.log('  Files in examples/:', fileNames.length);
        fileNames.slice(0, 5).forEach(function(name) {
            console.log('    -', name);
        });
        if (fileNames.length > 5) {
            console.log('    ... and', fileNames.length - 5, 'more');
        }
        console.log('');
        
        // Example 3: Check if a file exists (null = doesn't exist)
        console.log('3. Checking if file exists:');
        files.read('nonexistent.txt', function(err, data) {
            if (err) {
                console.error('  Error:', err);
            } else if (data === null) {
                console.log('  ✓ File does not exist (got null as expected)');
            } else {
                console.log('  × File exists unexpectedly');
            }
            console.log('');
            
            // Example 4: Write a file
            console.log('4. Writing a file:');
            const testContent = 'Hello from Dougless!\nCreated at: ' + new Date().toISOString();
            files.write('/tmp/dougless-test.txt', testContent, function(err) {
                if (err) {
                    console.error('  Error:', err);
                    return;
                }
                
                console.log('  ✓ File written successfully');
                
                // Verify by reading it back
                files.read('/tmp/dougless-test.txt', function(err, content) {
                    if (!err && content) {
                        console.log('  ✓ Verified - content length:', content.length);
                    }
                    console.log('');
                    
                    // Example 5: Write to nested path (auto-creates parent dirs)
                    console.log('5. Writing with auto-created parent directories:');
                    files.write('/tmp/dougless-test/nested/deep/file.txt', 'Nested content', function(err) {
                        if (err) {
                            console.error('  Error:', err);
                            return;
                        }
                        
                        console.log('  ✓ Created /tmp/dougless-test/nested/deep/ automatically');
                        console.log('');
                        
                        // Example 6: Create an empty directory
                        console.log('6. Creating an empty directory:');
                        files.write('/tmp/dougless-test/empty-dir/', function(err) {
                            if (err) {
                                console.error('  Error:', err);
                                return;
                            }
                            
                            console.log('  ✓ Directory created');
                            console.log('');
                            
                            // Example 7: Remove a file
                            console.log('7. Removing a file:');
                            files.rm('/tmp/dougless-test.txt', function(err) {
                                if (err) {
                                    console.error('  Error:', err);
                                    return;
                                }
                                
                                console.log('  ✓ File deleted');
                                console.log('');
                                
                                // Example 8: Remove directory recursively
                                console.log('8. Removing directory recursively:');
                                files.rm('/tmp/dougless-test/', function(err) {
                                    if (err) {
                                        console.error('  Error:', err);
                                        return;
                                    }
                                    
                                    console.log('  ✓ Directory and all contents removed');
                                    console.log('');
                                    
                                    // Example 9: Idempotent removal
                                    console.log('9. Removing non-existent path (idempotent):');
                                    files.rm('/tmp/does-not-exist.txt', function(err) {
                                        if (err) {
                                            console.error('  Error:', err);
                                        } else {
                                            console.log('  ✓ No error even though file doesn\'t exist');
                                        }
                                        console.log('');
                                        console.log('=== All examples complete! ===');
                                    });
                                });
                            });
                        });
                    });
                });
            });
        });
    });
});
