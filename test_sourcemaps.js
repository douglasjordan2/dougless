// Test script to verify source maps work with ES6+ syntax
const greet = (name) => {
  console.log(`Hello, ${name}!`);
};

greet('World');

// This will cause an error to test source map debugging
// Uncomment to test:
// nonExistentFunction();
