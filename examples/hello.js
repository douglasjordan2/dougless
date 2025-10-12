// =======================================
// Dougless Runtime - Hello World
// =======================================
// Your first Dougless program!
// This demonstrates basic features and unique globals

console.log('\n=== Welcome to Dougless Runtime! ===\n');

// 1. Basic console logging
console.log('Hello from Dougless Runtime!');
console.log('This is a custom JavaScript runtime built with Go.\n');

// 2. Unique global APIs (no require needed!)
console.log('Checking global APIs...');
console.log('  ✓ console:', typeof console);
console.log('  ✓ file:', typeof file);
console.log('  ✓ http:', typeof http);
console.log('  ✓ setTimeout:', typeof setTimeout);
console.log('  ✓ setInterval:', typeof setInterval);
console.log('');

// 3. Module system (CommonJS)
console.log('Testing module system...');
const path = require('path');
console.log('  ✓ path module loaded');
console.log('  ✓ path.join("a", "b", "c") =', path.join('a', 'b', 'c'));
console.log('');

// 4. Async operations with timers
console.log('Testing async operations...');
setTimeout(function() {
    console.log('  ✓ setTimeout works!');
    console.log('');
    console.log('=== All systems operational! ===');
    console.log('\nTry other examples:');
    console.log('  - console_features.js - Console operations');
    console.log('  - timers.js - setTimeout & setInterval');
    console.log('  - file_operations.js - File system API');
    console.log('  - http_demo.js - HTTP client & server');
    console.log('  - path_module.js - Path manipulation\n');
}, 100);
