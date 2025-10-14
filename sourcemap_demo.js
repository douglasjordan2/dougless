// Source Map Demonstration
// This script shows how source maps help with debugging

// ES6+ arrow function with template literal
const calculate = (a, b) => {
  const sum = a + b;
  const product = a * b;
  
  console.log(`Sum: ${sum}`);
  console.log(`Product: ${product}`);
  
  // Intentional error on line 12!
  return thisVariableDoesNotExist;
};

calculate(5, 10);
