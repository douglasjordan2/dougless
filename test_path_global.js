// Test 1: Path should be available globally
console.log('Test 1: Global path object exists');
console.log('typeof path:', typeof path);
console.log('path object:', path);
console.log('');

// Test 2: Path methods should work
console.log('Test 2: Path methods work');
const joined = path.join('foo', 'bar', 'baz.txt');
console.log('path.join("foo", "bar", "baz.txt"):', joined);

const resolved = path.resolve('foo', 'bar');
console.log('path.resolve("foo", "bar"):', resolved);

const dir = path.dirname('/home/user/file.txt');
console.log('path.dirname("/home/user/file.txt"):', dir);

const base = path.basename('/home/user/file.txt');
console.log('path.basename("/home/user/file.txt"):', base);

const ext = path.extname('file.txt');
console.log('path.extname("file.txt"):', ext);
console.log('');

// Test 3: require('path') should still work for backward compatibility
console.log('Test 3: require("path") still works');
const pathRequired = require('path');
console.log('typeof pathRequired:', typeof pathRequired);

const joinedReq = pathRequired.join('a', 'b', 'c');
console.log('pathRequired.join("a", "b", "c"):', joinedReq);
console.log('');

// Test 4: Both should reference the same functionality
console.log('Test 4: Comparison');
console.log('Global path.join === require path.join:', path.join('x', 'y') === pathRequired.join('x', 'y'));
console.log('');

console.log('✅ All path global tests passed!');
