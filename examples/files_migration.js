// files_migration.js - Migration examples from old file API to new files API
//
// This file shows side-by-side comparisons of old vs new patterns.
// Note: Old API patterns are commented out since they no longer work.

console.log('=== Files API Migration Examples ===\n');
console.log('Showing new `files` API patterns that replace old `file` API\n');

// MIGRATION 1: Reading a file
console.log('1. Reading a file:');
// OLD: file.read('data.txt', callback)
// NEW: files.read('data.txt', callback)  <- Same, just renamed!
console.log('   files.read(\'data.txt\', callback)  // Same pattern!\n');

// MIGRATION 2: Listing directory
console.log('2. Listing directory:');
// OLD: file.readdir('src', callback)
// NEW: files.read('src/', callback)  <- Add trailing slash
console.log('   files.read(\'src/\', callback)  // Note the trailing slash!\n');

// MIGRATION 3: Creating directory
console.log('3. Creating directory:');
// OLD: file.mkdir('newdir', callback)
// NEW: files.write('newdir/', callback)  <- Trailing slash, no content
console.log('   files.write(\'newdir/\', callback)  // Trailing slash = directory\n');

// MIGRATION 4: Checking if file exists
console.log('4. Checking file existence:');
// OLD: file.exists('data.txt', function(exists) { ... })
// NEW: files.read('data.txt', function(err, data) { if (data === null) ... })
console.log('   files.read(\'data.txt\', function(err, data) {');
console.log('     if (data === null) console.log("doesn\'t exist");');
console.log('   })  // null means file doesn\'t exist\n');

// MIGRATION 5: Deleting file
console.log('5. Deleting a file:');
// OLD: file.unlink('temp.txt', callback)
// NEW: files.rm('temp.txt', callback)
console.log('   files.rm(\'temp.txt\', callback)  // Works for files or dirs!\n');

// MIGRATION 6: Removing directory
console.log('6. Removing a directory:');
// OLD: file.rmdir('olddir', callback)  <- Only empty dirs
// NEW: files.rm('olddir/', callback)  <- Recursive by default!
console.log('   files.rm(\'olddir/\', callback)  // Recursive automatically!\n');

// MIGRATION 7: Getting file info
console.log('7. Getting file metadata:');
// OLD: file.stat('data.txt', callback)
// NEW: Not available in simplified API (removed)
console.log('   [REMOVED] - stat() not available in new API\n');

// PRACTICAL EXAMPLE: Old vs New
console.log('=== PRACTICAL COMPARISON ===\n');

console.log('OLD WAY - Creating project structure:');
console.log('  file.mkdir(\'project\', function(err) {');
console.log('    if (err) return;');
console.log('    file.mkdir(\'project/src\', function(err) {');
console.log('      if (err) return;');
console.log('      file.write(\'project/src/app.js\', code, callback);');
console.log('    });');
console.log('  });\n');

console.log('NEW WAY - Parent dirs created automatically:');
console.log('  files.write(\'project/src/app.js\', code, callback);');
console.log('  // That\'s it! No mkdir needed!\n');

// Live demonstration
console.log('=== LIVE DEMONSTRATION ===\n');

files.write('/tmp/demo/deep/nested/file.txt', 'Hello from new API!', function(err) {
    if (err) {
        console.error('Error:', err);
        return;
    }
    
    console.log('✓ Created /tmp/demo/deep/nested/file.txt');
    console.log('  (all parent directories created automatically)\n');
    
    // Read it back
    files.read('/tmp/demo/deep/nested/file.txt', function(err, content) {
        if (err) {
            console.error('Error:', err);
            return;
        }
        
        console.log('✓ Read back content:', content);
        console.log('');
        
        // List the directory
        files.read('/tmp/demo/deep/nested/', function(err, fileNames) {
            if (err) {
                console.error('Error:', err);
                return;
            }
            
            console.log('✓ Directory contains:', fileNames);
            console.log('');
            
            // Check if another file exists
            files.read('/tmp/demo/nonexistent.txt', function(err, data) {
                if (err) {
                    console.error('Error:', err);
                    return;
                }
                
                if (data === null) {
                    console.log('✓ Nonexistent file correctly returns null');
                }
                console.log('');
                
                // Clean up (recursive removal)
                files.rm('/tmp/demo/', function(err) {
                    if (err) {
                        console.error('Error:', err);
                        return;
                    }
                    
                    console.log('✓ Cleaned up entire /tmp/demo/ directory');
                    console.log('');
                    console.log('=== Migration demonstration complete! ===');
                });
            });
        });
    });
});
