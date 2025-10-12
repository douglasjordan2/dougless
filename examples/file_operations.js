// =========================================
// Dougless Runtime - File System Features
// =========================================
// Demonstrates the global file API (no require needed!)

console.log('=== File System Operations Demo ===\n');

// 1. Write a simple file
console.log('1. Writing a simple file...');
file.write('hello.txt', 'Hello from Dougless Runtime!', function(err) {
    if (err) {
        console.error('   ✗ Write error:', err);
        return;
    }
    console.log('   ✓ File written: hello.txt\n');
    
    // 2. Read the file back
    console.log('2. Reading the file back...');
    file.read('hello.txt', function(err, data) {
        if (err) {
            console.error('   ✗ Read error:', err);
            return;
        }
        console.log('   ✓ Content:', data);
        console.log('');
        
        // 3. Check if file exists
        console.log('3. Checking if file exists...');
        file.exists('hello.txt', function(exists) {
            console.log('   ✓ hello.txt exists:', exists);
            
            file.exists('nonexistent.txt', function(exists) {
                console.log('   ✓ nonexistent.txt exists:', exists);
                console.log('');
                
                // 4. Create a directory
                console.log('4. Creating directory "test-dir"...');
                file.mkdir('test-dir', function(err) {
                    if (err) {
                        console.error('   ✗ mkdir error:', err);
                        return;
                    }
                    console.log('   ✓ Directory created\n');
                    
                    // 5. Write multiple files in the directory
                    console.log('5. Writing multiple files in directory...');
                    
                    file.write('test-dir/file1.txt', 'Content 1', function(err) {
                        if (err) console.error('   ✗ Error:', err);
                        else console.log('   ✓ file1.txt created');
                    });
                    
                    file.write('test-dir/file2.txt', 'Content 2', function(err) {
                        if (err) console.error('   ✗ Error:', err);
                        else console.log('   ✓ file2.txt created');
                    });
                    
                    file.write('test-dir/file3.txt', 'Content 3', function(err) {
                        if (err) console.error('   ✗ Error:', err);
                        else console.log('   ✓ file3.txt created');
                        
                        // Small delay to ensure all files are written
                        setTimeout(function() {
                            console.log('');
                            
                            // 6. List directory contents
                            console.log('6. Listing directory contents...');
                            file.readdir('test-dir', function(err, files) {
                                if (err) {
                                    console.error('   ✗ readdir error:', err);
                                    return;
                                }
                                console.log('   ✓ Files in test-dir:', files);
                                console.log('');
                                
                                // 7. Get file stats
                                console.log('7. Getting file statistics...');
                                file.stat('test-dir/file1.txt', function(err, stats) {
                                    if (err) {
                                        console.error('   ✗ stat error:', err);
                                        return;
                                    }
                                    console.log('   ✓ File stats:');
                                    console.log('      - Name:', stats.name);
                                    console.log('      - Size:', stats.size, 'bytes');
                                    console.log('      - Is file:', stats.isFile);
                                    console.log('      - Is directory:', stats.isDirectory);
                                    console.log('      - Modified:', new Date(stats.modified * 1000).toISOString());
                                    console.log('');
                                    
                                    // 8. Read and process file
                                    console.log('8. Reading and processing file...');
                                    file.read('test-dir/file1.txt', function(err, data) {
                                        if (err) {
                                            console.error('   ✗ Read error:', err);
                                            return;
                                        }
                                        
                                        const processed = data.toUpperCase();
                                        console.log('   ✓ Original:', data);
                                        console.log('   ✓ Processed:', processed);
                                        console.log('');
                                        
                                        // 9. Write processed data
                                        console.log('9. Writing processed data...');
                                        file.write('test-dir/file1-upper.txt', processed, function(err) {
                                            if (err) {
                                                console.error('   ✗ Write error:', err);
                                                return;
                                            }
                                            console.log('   ✓ Processed file created\n');
                                            
                                            // 10. Cleanup - delete files
                                            console.log('10. Cleaning up - deleting files...');
                                            
                                            file.unlink('test-dir/file1.txt', function(err) {
                                                if (err) console.error('   ✗ Delete error:', err);
                                                else console.log('   ✓ file1.txt deleted');
                                            });
                                            
                                            file.unlink('test-dir/file2.txt', function(err) {
                                                if (err) console.error('   ✗ Delete error:', err);
                                                else console.log('   ✓ file2.txt deleted');
                                            });
                                            
                                            file.unlink('test-dir/file3.txt', function(err) {
                                                if (err) console.error('   ✗ Delete error:', err);
                                                else console.log('   ✓ file3.txt deleted');
                                            });
                                            
                                            file.unlink('test-dir/file1-upper.txt', function(err) {
                                                if (err) console.error('   ✗ Delete error:', err);
                                                else console.log('   ✓ file1-upper.txt deleted');
                                                
                                                // Small delay to ensure all deletes complete
                                                setTimeout(function() {
                                                    console.log('');
                                                    
                                                    // 11. Remove directory
                                                    console.log('11. Removing directory...');
                                                    file.rmdir('test-dir', function(err) {
                                                        if (err) {
                                                            console.error('   ✗ rmdir error:', err);
                                                            return;
                                                        }
                                                        console.log('   ✓ Directory removed\n');
                                                        
                                                        // 12. Final cleanup
                                                        console.log('12. Final cleanup...');
                                                        file.unlink('hello.txt', function(err) {
                                                            if (err) {
                                                                console.error('   ✗ Delete error:', err);
                                                            } else {
                                                                console.log('   ✓ hello.txt deleted');
                                                            }
                                                            
                                                            console.log('\n=== File Operations Demo Complete ===');
                                                            console.log('All file operations demonstrated successfully!');
                                                        });
                                                    });
                                                }, 100);
                                            });
                                        });
                                    });
                                });
                            });
                        }, 100);
                    });
                });
            });
        });
    });
});
