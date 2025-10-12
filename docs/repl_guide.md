# REPL Guide

## Overview

Dougless Runtime includes an interactive REPL (Read-Eval-Print Loop) that allows you to experiment with JavaScript code in real-time. The REPL supports multi-line input, maintains state between commands, and provides helpful utilities for interactive development.

## Starting the REPL

Simply run `dougless` without any arguments:

```bash
./dougless
```

You'll see the welcome banner:
```
Dougless Runtime REPL v0.1.0
Type JavaScript code to evaluate, or .help for commands

> 
```

## Basic Usage

### Simple Expressions

```javascript
> 2 + 2
4

> "Hello, " + "World!"
Hello, World!

> Math.PI * 2
6.283185307179586
```

### Variable Assignment

The REPL maintains state between evaluations:

```javascript
> let x = 10
undefined

> x * 5
50

> const name = "Dougless"
undefined

> "Hello, " + name
Hello, Dougless
```

### Multi-line Input

The REPL automatically detects incomplete input and prompts for continuation:

```javascript
> function greet(name) {
...   return "Hello, " + name;
... }
undefined

> greet("World")
Hello, World!
```

Multi-line detection works for:
- Function definitions
- Object literals
- Array literals
- Any unclosed brackets, braces, or parentheses

### Using Built-in Features

#### Console Operations
```javascript
> console.log("Testing", 123)
Testing 123
undefined

> console.time("test")
undefined

> for (let i = 0; i < 1000000; i++) {}
undefined

> console.timeEnd("test")
test: 12.345ms
undefined
```

#### Timers (with event loop)
```javascript
> setTimeout(function() { console.log("Delayed!"); }, 1000)
[object Object]

> // After 1 second:
Delayed!
```

Note: The REPL will wait for all scheduled timers to complete before returning to the prompt.

## Special Commands

The REPL supports several special commands starting with a dot (`.`):

### `.help`
Display available commands:
```
> .help
Available commands:
  .help   - Show this help message
  .exit   - Exit the REPL (or Ctrl+D)
  .quit   - Same as .exit
  .clear  - Clear the screen
```

### `.exit` or `.quit`
Exit the REPL:
```
> .exit
Goodbye!
```

You can also use `Ctrl+D` (EOF) to exit.

### `.clear`
Clear the terminal screen:
```
> .clear
```

This uses ANSI escape codes to clear the screen and reset the cursor.

## Tips & Tricks

### 1. Use `let` or `const` for persistent variables
Variables declared with `let` or `const` persist across REPL evaluations:
```javascript
> const config = { port: 3000 }
> config.port
3000
```

### 2. Inspect objects
Just type the variable name to see its value:
```javascript
> const obj = { a: 1, b: 2, c: 3 }
> obj
[object Object]
```

### 3. Test functions interactively
Define functions and test them immediately:
```javascript
> function fibonacci(n) {
...   if (n <= 1) return n;
...   return fibonacci(n - 1) + fibonacci(n - 2);
... }

> fibonacci(10)
55
```

### 4. Use console.table for data
```javascript
> const data = [
...   { name: "Alice", age: 30 },
...   { name: "Bob", age: 25 }
... ]

> console.table(data)
┌─────────┬──────┐
│  name   │ age  │
├─────────┼──────┤
│  Alice  │  30  │
│   Bob   │  25  │
└─────────┴──────┘
```

### 5. Experiment with timers
Test async behavior:
```javascript
> let count = 0
> const interval = setInterval(function() {
...   count++;
...   console.log("Count:", count);
...   if (count === 3) clearInterval(interval);
... }, 500)
```

## Error Handling

The REPL displays JavaScript errors clearly:

```javascript
> unknownVariable
Error: ReferenceError: unknownVariable is not defined

> JSON.parse("{invalid}")
Error: SyntaxError: Unexpected token i in JSON at position 1
```

Errors don't crash the REPL - it continues running normally.

## Limitations

### ES5.1 Syntax Only
The REPL uses the Goja JavaScript engine which supports ES5.1:

**Supported:**
- `let` and `const` declarations
- `function` keyword
- Traditional `for` loops
- Object and array literals

**Not Supported (yet):**
- Arrow functions (`=>`)
- Template literals
- `async`/`await`
- Destructuring
- Classes

ES6+ support is planned through transpilation in future phases.

### No Module Loading in REPL
Currently, `require()` is not fully functional in the REPL. Module support is planned for Phase 2.

## Keyboard Shortcuts

- **Enter** - Evaluate current line (or continue multi-line input)
- **Ctrl+D** - Exit the REPL (EOF)
- **Ctrl+C** - Cancel current input (if implemented)

## Comparing to Node.js REPL

If you're familiar with Node.js, the Dougless REPL is similar but simpler:

| Feature | Node.js REPL | Dougless REPL |
|---------|--------------|---------------|
| Basic evaluation | ✅ | ✅ |
| Multi-line input | ✅ | ✅ |
| State persistence | ✅ | ✅ |
| Special commands | ✅ | ✅ (subset) |
| Tab completion | ✅ | ❌ (future) |
| History (up/down) | ✅ | ❌ (future) |
| ES6+ syntax | ✅ | ❌ (future) |
| require() in REPL | ✅ | ⏳ (Phase 2) |

## Examples

### Example 1: Quick Math
```javascript
> const radius = 5
> const area = Math.PI * radius * radius
> console.log("Area:", area)
Area: 78.53981633974483
```

### Example 2: String Manipulation
```javascript
> const message = "hello world"
> message.toUpperCase()
HELLO WORLD
> message.split(" ")
hello,world
```

### Example 3: Timer Demo
```javascript
> console.time("total")
> setTimeout(function() {
...   console.log("Step 1");
...   setTimeout(function() {
...     console.log("Step 2");
...     console.timeEnd("total");
...   }, 500);
... }, 500)
// After 500ms: Step 1
// After 1000ms: Step 2
// total: 1002.345ms
```

## Troubleshooting

### REPL doesn't start
Make sure you're running `dougless` without any arguments:
```bash
./dougless          # ✅ Starts REPL
./dougless file.js  # ❌ Runs file instead
```

### Multi-line input stuck
If the REPL keeps showing `...` prompts, you have unclosed brackets. Either:
1. Close all brackets/braces/parentheses
2. Press `Ctrl+C` to cancel (if implemented)
3. Press `Ctrl+D` to exit and restart

### Error: "undefined is not a function"
You might be trying to use an ES6+ feature. Remember, Dougless currently supports ES5.1 syntax only.

---

**Happy coding with Dougless REPL!** 🚀
