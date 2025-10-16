// files_touch.js - Demonstrate creating empty files (like 'touch' command)

console.log('=== Empty File Creation (Touch-like Behavior) ===\n');

async function demo() {
    try {
        // Create an empty file - content is optional!
        console.log('1. Creating empty file:');
        await files.write('/tmp/empty-file.txt');
        console.log('   ✓ Created /tmp/empty-file.txt (0 bytes)');
        
        // Verify it's empty
        const content = await files.read('/tmp/empty-file.txt');
        console.log('   ✓ Verified: file is empty (length:', content.length, 'bytes)\n');
        
        // Write content to existing file
        console.log('2. Adding content to empty file:');
        await files.write('/tmp/empty-file.txt', 'Now it has content!');
        const newContent = await files.read('/tmp/empty-file.txt');
        console.log('   ✓ Content added:', newContent, '\n');
        
        // Truncate file back to empty
        console.log('3. Truncating file (back to empty):');
        await files.write('/tmp/empty-file.txt');
        const truncated = await files.read('/tmp/empty-file.txt');
        console.log('   ✓ Truncated: file is empty again (length:', truncated.length, 'bytes)\n');
        
        // Touch-like: create or update timestamp
        console.log('4. Touch multiple files at once:');
        await Promise.all([
            files.write('/tmp/file1.txt'),
            files.write('/tmp/file2.txt'),
            files.write('/tmp/file3.txt')
        ]);
        console.log('   ✓ Created 3 empty files in parallel\n');
        
        // Compare with traditional approach
        console.log('5. API Comparison:');
        console.log('   Traditional:  files.write(path, "", callback)');
        console.log('   Improved:     files.write(path)  // Much cleaner!');
        console.log('   Traditional:  files.write(path, "data", callback)');
        console.log('   Improved:     files.write(path, "data")  // Same!\n');
        
        // Cleanup
        await Promise.all([
            files.rm('/tmp/empty-file.txt'),
            files.rm('/tmp/file1.txt'),
            files.rm('/tmp/file2.txt'),
            files.rm('/tmp/file3.txt')
        ]);
        console.log('✓ Cleanup complete!');
        
        console.log('\n=== Demo Complete ===');
        console.log('Key benefits:');
        console.log('  • Consistent API - content always optional');
        console.log('  • Touch-like behavior - create or truncate files');
        console.log('  • Cleaner code - no empty strings needed');
        console.log('  • Works with both callbacks and promises');
        
    } catch (err) {
        console.error('Error:', err);
    }
}

demo();
