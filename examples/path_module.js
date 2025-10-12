// ===================================
// Dougless Runtime - Path Module
// ===================================
// Demonstrates the path module (CommonJS require)

const path = require('path');

console.log('=== Path Module Demo ===\n');

// Test 1: path.join()
console.log('1. path.join():');
const joined = path.join('foo', 'bar', 'baz.txt');
console.log('   path.join("foo", "bar", "baz.txt") =>', joined);

const joined2 = path.join('/home', 'user', 'documents', 'file.txt');
console.log('   path.join("/home", "user", "documents", "file.txt") =>', joined2);
console.log('');

// Test 2: path.dirname()
console.log('2. path.dirname():');
const dir = path.dirname('/home/user/file.txt');
console.log('   path.dirname("/home/user/file.txt") =>', dir);
console.log('');

// Test 3: path.basename()
console.log('3. path.basename():');
const base = path.basename('/home/user/file.txt');
console.log('   path.basename("/home/user/file.txt") =>', base);

const baseNoExt = path.basename('/home/user/file.txt', '.txt');
console.log('   path.basename("/home/user/file.txt", ".txt") =>', baseNoExt);
console.log('');

// Test 4: path.extname()
console.log('4. path.extname():');
const ext = path.extname('file.txt');
console.log('   path.extname("file.txt") =>', ext);

const ext2 = path.extname('archive.tar.gz');
console.log('   path.extname("archive.tar.gz") =>', ext2);
console.log('');

// Test 5: path.resolve()
console.log('5. path.resolve():');
const absolute = path.resolve('examples', 'test-path.js');
console.log('   path.resolve("examples", "test-path.js") =>', absolute);
console.log('');

// Test 6: path.sep
console.log('6. path.sep (path separator):');
console.log('   path.sep =>', path.sep);

console.log('\n✅ All path tests complete!');
