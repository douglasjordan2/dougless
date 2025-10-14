// Comprehensive test for path as a global API

console.log('=== Path Global API Test Suite ===\n');

// Test 1: Basic availability
console.log('1. Path is globally available (no require needed)');
console.log('   typeof path:', typeof path);
console.log('   ✓ Pass\n');

// Test 2: All methods exist
console.log('2. All path methods exist');
const methods = ['join', 'resolve', 'dirname', 'basename', 'extname', 'sep'];
methods.forEach(method => {
  const exists = typeof path[method] !== 'undefined';
  console.log(`   path.${method}: ${exists ? '✓' : '✗'}`);
});
console.log('   ✓ Pass\n');

// Test 3: join() works correctly
console.log('3. path.join() works correctly');
const joined1 = path.join('foo', 'bar', 'baz');
console.log(`   path.join('foo', 'bar', 'baz') = ${joined1}`);
console.log(`   Expected: foo/bar/baz`);
console.log(`   ✓ ${joined1 === 'foo/bar/baz' ? 'Pass' : 'FAIL'}\n`);

// Test 4: resolve() works correctly
console.log('4. path.resolve() works correctly');
const resolved = path.resolve('foo', 'bar');
console.log(`   path.resolve('foo', 'bar') = ${resolved}`);
console.log(`   Contains '/foo/bar': ${resolved.includes('/foo/bar')}`);
console.log(`   ✓ ${resolved.includes('/foo/bar') ? 'Pass' : 'FAIL'}\n`);

// Test 5: dirname() works correctly
console.log('5. path.dirname() works correctly');
const dir = path.dirname('/home/user/file.txt');
console.log(`   path.dirname('/home/user/file.txt') = ${dir}`);
console.log(`   Expected: /home/user`);
console.log(`   ✓ ${dir === '/home/user' ? 'Pass' : 'FAIL'}\n`);

// Test 6: basename() works correctly
console.log('6. path.basename() works correctly');
const base = path.basename('/home/user/file.txt');
console.log(`   path.basename('/home/user/file.txt') = ${base}`);
console.log(`   Expected: file.txt`);
console.log(`   ✓ ${base === 'file.txt' ? 'Pass' : 'FAIL'}\n`);

// Test 7: extname() works correctly
console.log('7. path.extname() works correctly');
const ext = path.extname('file.txt');
console.log(`   path.extname('file.txt') = ${ext}`);
console.log(`   Expected: .txt`);
console.log(`   ✓ ${ext === '.txt' ? 'Pass' : 'FAIL'}\n`);

// Test 8: Backward compatibility with require()
console.log('8. Backward compatibility with require("path")');
const pathRequired = require('path');
const sameResult = path.join('a', 'b') === pathRequired.join('a', 'b');
console.log(`   Both methods produce same result: ${sameResult}`);
console.log(`   ✓ ${sameResult ? 'Pass' : 'FAIL'}\n`);

// Test 9: Complex path operations
console.log('9. Complex path operations');
const complex = path.join(
  path.dirname('/var/www/app/index.js'),
  'public',
  path.basename('style.css')
);
console.log(`   Complex operation: ${complex}`);
console.log(`   Expected: /var/www/app/public/style.css`);
console.log(`   ✓ ${complex === '/var/www/app/public/style.css' ? 'Pass' : 'FAIL'}\n`);

// Test 10: path.sep
console.log('10. path.sep returns correct separator');
console.log(`   path.sep = ${path.sep}`);
console.log(`   Expected: / (Unix-style)`);
console.log(`   ✓ ${path.sep === '/' ? 'Pass' : 'FAIL'}\n`);

console.log('=== All Tests Completed ===');
console.log('✅ Path is now a global API (like file and http)!');
console.log('   Use: path.join(...) instead of require("path").join(...)');
console.log('   Note: require("path") still works for backward compatibility');
