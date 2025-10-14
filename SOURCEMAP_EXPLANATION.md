# Source Maps - Visual Explanation

## What You Just Saw

### Original Code (What YOU Write)
```javascript
// Modern ES6+ JavaScript
const greet = (name) => {
  console.log(`Hello, ${name}!`);
};

const calculate = (x, y) => {
  const sum = x + y;
  const product = x * y;
  return { sum, product };
};

// Destructuring
const { sum, product } = calculate(5, 10);

// Spread operator
const numbers = [1, 2, 3];
const moreNumbers = [...numbers, 4, 5];

// Class syntax
class Person {
  constructor(name, age) {
    this.name = name;
    this.age = age;
  }
  
  greet() {
    return `I'm ${this.name}, ${this.age} years old`;
  }
}
```

### Transpiled Code (What Goja EXECUTES)
```javascript
const greet = (name) => {
  console.log(`Hello, ${name}!`);
};
const calculate = (x, y) => {
  const sum2 = x + y;           // ← NOTICE: 'sum' became 'sum2'!
  const product2 = x * y;       // ← NOTICE: 'product' became 'product2'!
  return { sum: sum2, product: product2 };  // ← Shorthand expanded!
};
const { sum, product } = calculate(5, 10);
const numbers = [1, 2, 3];
const moreNumbers = [...numbers, 4, 5];
class Person {
  constructor(name, age) {
    this.name = name;
    this.age = age;
  }
  greet() {
    return `I'm ${this.name}, ${this.age} years old`;
  }
}
```

---

## The Problem Without Source Maps

Imagine you have this bug:

```javascript
const calculate = (x, y) => {
  const sum = x + y;
  const product = x * y;
  console.log(sum2);  // ← BUG! Should be 'sum' not 'sum2'
  return { sum, product };
};
```

**WITHOUT Source Maps:**
```
Error: sum2 is not defined
  at calculate (example.js:6:15)
```

But wait... in the **transpiled** code, line 6 might be completely different because:
- Variable names changed (`sum` → `sum2`)
- Object shorthand was expanded
- Line numbers shifted

You'd be staring at:
```javascript
const sum2 = x + y;
const product2 = x * y;
console.log(sum2);  // Line 6 in transpiled code
return { sum: sum2, product: product2 };
```

And thinking "But sum2 DOES exist! What's going on?!" 😵

---

**WITH Source Maps:**
```
Error: sum2 is not defined
  at calculate (example.js:4:15)
     → Original source: console.log(sum2);
```

Now you see line **4** in YOUR original code, and you immediately spot the typo!

---

## How Source Maps Work

The source map contains a lookup table (simplified):

```
Transpiled Code         →    Original Code
─────────────────────────────────────────────
Line 5, Column 10       →    Line 7, Column 10
Line 6, Column 15       →    Line 8, Column 15
Line 9, Column 5        →    Line 10, Column 5
```

When Goja throws an error at "Line 6, Column 15" in the transpiled code, it checks the source map and reports "Line 8, Column 15" in your original code!

---

## The Source Map Data

Here's what that base64 string contains (decoded):

```json
{
  "version": 3,
  "sources": ["example.js"],
  "sourcesContent": ["// Your original code here..."],
  "mappings": "AAAA,MAAM,KAAK,GAAG...",
  "names": ["greet", "name", "console", "log", "calculate", "sum", "product"]
}
```

### Key Parts:

1. **`sources`**: The original filename(s)
2. **`sourcesContent`**: Your actual original code (embedded!)
3. **`mappings`**: Encoded line/column mappings using VLQ format
4. **`names`**: Variable/function names for better debugging

---

## Why Dougless Needs This

Even though Dougless targets ES2017 (which supports most modern syntax), esbuild still makes changes:

✅ **Variable Renaming** (to avoid conflicts)
   - `sum` → `sum2`
   - `product` → `product2`

✅ **Object Shorthand Expansion**
   - `{ sum, product }` → `{ sum: sum2, product: product2 }`

✅ **Line Number Shifts** (from formatting/optimization)

Without source maps, you'd be debugging the *transformed* code, not your code!

---

## Real-World Example

Run this in Dougless to see source maps in action:

```bash
# Create a test file
cat > test_error.js << 'EOF'
const calculate = (x, y) => {
  const result = x + y;
  console.log(resultTypo);  // Intentional error!
  return result;
};

calculate(5, 10);
EOF

# Run it
./dougless test_error.js
```

**Output with source maps:**
```
Error: ReferenceError: resultTypo is not defined
  at calculate (test_error.js:3:15)
```

Line 3! That's exactly where you wrote it in your source file! 🎯

**Output without source maps:**
```
Error: ReferenceError: resultTypo is not defined
  at calculate (test_error.js:4:37)
```

Line 4? Column 37? But your file only has 3 lines of code in that function... 🤔

---

## Summary

| Feature | Without Source Maps | With Source Maps |
|---------|-------------------|------------------|
| **Error Line Numbers** | Points to transpiled code | Points to YOUR code |
| **Variable Names** | Shows renamed variables | Shows original names |
| **Debugging Experience** | Confusing | Intuitive |
| **Stack Traces** | Hard to follow | Easy to follow |

**Enabled = Happy developers! 🎉**
