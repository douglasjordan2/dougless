// files_advanced.js - Advanced patterns with the files API

console.log('=== Dougless Files API - Advanced Examples ===\n');

// Example 1: Config file with default fallback
console.log('1. Config file pattern (with default):');
function loadConfig(callback) {
    files.read('config.json', function(err, data) {
        if (err) {
            callback(err, null);
            return;
        }
        
        if (data === null) {
            // File doesn't exist, use default
            const defaultConfig = { version: '1.0', debug: false };
            console.log('  Using default config:', JSON.stringify(defaultConfig));
            callback(null, defaultConfig);
        } else {
            try {
                const config = JSON.parse(data);
                console.log('  Loaded config from file:', JSON.stringify(config));
                callback(null, config);
            } catch (parseErr) {
                callback('Invalid JSON in config file: ' + parseErr, null);
            }
        }
    });
}

loadConfig(function(err, config) {
    if (err) {
        console.error('  Error:', err);
        return;
    }
    
    console.log('  Config ready!\n');
    
    // Example 2: Build tool pattern (clean + build)
    console.log('2. Build tool pattern (clean then build):');
    const distDir = '/tmp/dougless-build/';
    
    // Clean first
    files.rm(distDir, function(err) {
        if (err) {
            console.error('  Clean failed:', err);
            return;
        }
        
        console.log('  ✓ Cleaned dist/');
        
        // Build outputs (directories auto-created)
        let pending = 3;
        let errors = [];
        
        function checkComplete() {
            pending--;
            if (pending === 0) {
                if (errors.length > 0) {
                    console.error('  Build failed:', errors);
                } else {
                    console.log('  ✓ Build complete!\n');
                    
                    // Example 3: Directory tree traversal
                    console.log('3. Listing build output:');
                    files.read(distDir, function(err, items) {
                        if (err) {
                            console.error('  Error:', err);
                            return;
                        }
                        
                        console.log('  Build artifacts:');
                        items.forEach(function(item) {
                            console.log('    dist/' + item);
                        });
                        console.log('');
                        
                        // Example 4: Backup pattern
                        console.log('4. Backup pattern:');
                        const importantFile = distDir + 'app.js';
                        files.read(importantFile, function(err, content) {
                            if (err || content === null) {
                                console.error('  Source file not found');
                                return;
                            }
                            
                            const backupFile = importantFile + '.backup';
                            files.write(backupFile, content, function(err) {
                                if (err) {
                                    console.error('  Backup failed:', err);
                                } else {
                                    console.log('  ✓ Created backup:', backupFile);
                                }
                                console.log('');
                                
                                // Example 5: Batch operations
                                console.log('5. Batch file operations:');
                                const filesToCreate = [
                                    { path: '/tmp/dougless-batch/a.txt', content: 'File A' },
                                    { path: '/tmp/dougless-batch/b.txt', content: 'File B' },
                                    { path: '/tmp/dougless-batch/c.txt', content: 'File C' }
                                ];
                                
                                let batchPending = filesToCreate.length;
                                let created = 0;
                                
                                filesToCreate.forEach(function(file) {
                                    files.write(file.path, file.content, function(err) {
                                        if (!err) created++;
                                        batchPending--;
                                        
                                        if (batchPending === 0) {
                                            console.log('  ✓ Created', created, 'of', filesToCreate.length, 'files');
                                            
                                            // Clean up
                                            files.rm('/tmp/dougless-batch/', function() {
                                                files.rm(distDir, function() {
                                                    console.log('');
                                                    console.log('=== All advanced examples complete! ===');
                                                });
                                            });
                                        }
                                    });
                                });
                            });
                        });
                    });
                }
            }
        }
        
        files.write(distDir + 'js/app.js', '// App code', function(err) {
            if (err) errors.push(err);
            else console.log('    ✓ Built js/app.js');
            checkComplete();
        });
        
        files.write(distDir + 'css/styles.css', '/* Styles */', function(err) {
            if (err) errors.push(err);
            else console.log('    ✓ Built css/styles.css');
            checkComplete();
        });
        
        files.write(distDir + 'index.html', '<html>...</html>', function(err) {
            if (err) errors.push(err);
            else console.log('    ✓ Built index.html');
            checkComplete();
        });
    });
});
