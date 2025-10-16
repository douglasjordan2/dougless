// files_promise.js - Demonstration of the files API using promises and async/await

console.log('=== Dougless Files API - Promise/Async Examples ===\n');

// Using async/await for cleaner code
async function runExamples() {
    try {
        // Example 1: Read a file
        console.log('1. Reading a file with await:');
        const content = await files.read('examples/hello.js');
        
        if (content === null) {
            console.log('  File does not exist');
        } else {
            console.log('  File length:', content.length, 'bytes');
            console.log('  First 50 chars:', content.substring(0, 250));
        }
        console.log('');
        
        // Example 2: Read a directory (note the trailing slash)
        console.log('2. Listing directory contents with await:');
        const fileNames = await files.read('examples/');
        
        console.log('  Files in examples/:', fileNames.length);
        fileNames.slice(0, 5).forEach(name => {
            console.log('    -', name);
        });
        if (fileNames.length > 5) {
            console.log('    ... and', fileNames.length - 5, 'more');
        }
        console.log('');
        
        // Example 3: Check if a file exists (null = doesn't exist)
        console.log('3. Checking if file exists with await:');
        const data = await files.read('nonexistent.txt');
        
        if (data === null) {
            console.log('  ✓ File does not exist (got null as expected)');
        } else {
            console.log('  × File exists unexpectedly');
        }
        console.log('');
        
        // Example 4: Write a file
        console.log('4. Writing a file with await:');
        const testContent = 'Hello from Dougless!\nCreated at: ' + new Date().toISOString();
        await files.write('/tmp/dougless-test-promise.txt', testContent);
        console.log('  ✓ File written successfully');
        
        // Verify by reading it back
        const written = await files.read('/tmp/dougless-test-promise.txt');
        console.log('  ✓ Verified - content length:', written.length);
        console.log('');
        
        // Example 5: Write to nested path (auto-creates parent dirs)
        console.log('5. Writing with auto-created parent directories:');
        await files.write('/tmp/dougless-promise/nested/deep/file.txt', 'Nested content');
        console.log('  ✓ Created /tmp/dougless-promise/nested/deep/ automatically');
        console.log('');
        
        // Example 6: Create an empty directory
        console.log('6. Creating an empty directory:');
        await files.write('/tmp/dougless-promise/empty-dir/');
        console.log('  ✓ Directory created');
        console.log('');
        
        // Example 7: Remove a file
        console.log('7. Removing a file:');
        await files.rm('/tmp/dougless-test-promise.txt');
        console.log('  ✓ File deleted');
        console.log('');
        
        // Example 8: Remove directory recursively
        console.log('8. Removing directory recursively:');
        await files.rm('/tmp/dougless-promise/');
        console.log('  ✓ Directory and all contents removed');
        console.log('');
        
        // Example 9: Idempotent removal
        console.log('9. Removing non-existent path (idempotent):');
        await files.rm('/tmp/does-not-exist-promise.txt');
        console.log('  ✓ No error even though file doesn\'t exist');
        console.log('');
        
        // Example 10: Parallel operations with Promise.all()
        console.log('10. Parallel file operations with Promise.all():');
        await Promise.all([
            files.write('/tmp/parallel-1.txt', 'File 1'),
            files.write('/tmp/parallel-2.txt', 'File 2'),
            files.write('/tmp/parallel-3.txt', 'File 3')
        ]);
        console.log('  ✓ Three files written in parallel');
        
        const [content1, content2, content3] = await Promise.all([
            files.read('/tmp/parallel-1.txt'),
            files.read('/tmp/parallel-2.txt'),
            files.read('/tmp/parallel-3.txt')
        ]);
        console.log('  ✓ Three files read in parallel');
        console.log('  Contents:', content1, '|', content2, '|', content3);
        
        // Cleanup
        await Promise.all([
            files.rm('/tmp/parallel-1.txt'),
            files.rm('/tmp/parallel-2.txt'),
            files.rm('/tmp/parallel-3.txt')
        ]);
        console.log('  ✓ Cleanup complete');
        console.log('');
        
        console.log('=== All examples complete! ===');
        
    } catch (err) {
        console.error('Error:', err);
    }
}

// Run the async examples
runExamples();
