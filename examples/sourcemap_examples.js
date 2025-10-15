// ===================================
// Dougless Runtime - Source Map Examples
// ===================================
// Demonstrates how source maps work with transpiled ES6+ code

console.log('=== Source Map Examples ===\n');

// This arrow function will be transpiled to ES5
const greet = (name) => {
  console.log(`Hello, ${name}!`);
};

greet('World');

// Intentional error to show source mapping
console.log('\nTrying to cause an error for source map demo:');
try {
  const obj = { foo: 'bar' };
  console.log(obj.baz.qux); // This will throw
} catch (e) {
  console.log('Error caught:', e.message);
  console.log('(Check that error line numbers match this source file)');
}

console.log('\n✅ Source map example complete!');
console.log('💡 Errors should reference original line numbers, not transpiled code');
