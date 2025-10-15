// ===================================
// Dougless Runtime - Path Module Examples
// ===================================
// Demonstrates the path module as a global API
// Note: require('path') still works for backward compatibility

console.log('=== Path Module Examples ===\n');

// ==========================================
// PART 1: Global Access (Recommended)
// ==========================================
console.log('1. Path is globally available (no require needed)');
console.log('   typeof path:', typeof path);
console.log('');

// ==========================================
// PART 2: All Methods
// ==========================================
console.log('2. path.join() - Join path segments');
const joined = path.join('foo', 'bar', 'baz.txt');
console.log('   path.join("foo", "bar", "baz.txt"):', joined);

const joined2 = path.join('/home', 'user', 'documents', 'file.txt');
console.log('   path.join("/home", "user", "documents", "file.txt"):', joined2);
console.log('');

console.log('3. path.resolve() - Resolve to absolute path');
const resolved = path.resolve('examples', 'test.js');
console.log('   path.resolve("examples", "test.js"):', resolved);
console.log('   (Result depends on current working directory)');
console.log('');

console.log('4. path.dirname() - Get directory name');
const dir = path.dirname('/home/user/file.txt');
console.log('   path.dirname("/home/user/file.txt"):', dir);
console.log('');

console.log('5. path.basename() - Get file name');
const base = path.basename('/home/user/file.txt');
console.log('   path.basename("/home/user/file.txt"):', base);

const baseNoExt = path.basename('/home/user/file.txt', '.txt');
console.log('   path.basename("/home/user/file.txt", ".txt"):', baseNoExt);
console.log('');

console.log('6. path.extname() - Get file extension');
const ext = path.extname('file.txt');
console.log('   path.extname("file.txt"):', ext);

const ext2 = path.extname('archive.tar.gz');
console.log('   path.extname("archive.tar.gz"):', ext2);
console.log('');

console.log('7. path.sep - Path separator');
console.log('   path.sep:', path.sep);
console.log('   (Unix: "/", Windows: "\\\\")');
console.log('');

// ==========================================
// PART 3: Complex Example
// ==========================================
console.log('8. Complex path operations');
const complex = path.join(
  path.dirname('/var/www/app/index.js'),
  'public',
  path.basename('style.css')
);
console.log('   Building path from components:', complex);
console.log('');

// ==========================================
// PART 4: Backward Compatibility
// ==========================================
console.log('9. Backward compatibility with require()');
const pathRequired = require('path');
const sameResult = path.join('a', 'b') === pathRequired.join('a', 'b');
console.log('   require("path") still works:', sameResult);
console.log('   Both methods produce same result');
console.log('');

console.log('✅ All path examples complete!');
console.log('\n💡 TIP: Use "path.method()" directly - no require() needed!');
