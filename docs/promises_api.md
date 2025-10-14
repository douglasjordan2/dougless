# Promises API Guide

## Overview

Dougless provides a native Promise implementation that is fully compliant with the Promise/A+ specification. The `Promise` object is available globally, allowing you to write modern asynchronous JavaScript code with ease.

## Why Global?

**Dougless Philosophy**: Promises are fundamental to modern JavaScript and async programming. They should be as accessible as `console` or `setTimeout`. This makes Dougless code cleaner and eliminates the need for polyfills.

**ES6+ Compatibility**: Combined with our esbuild transpilation, you can use modern `async/await` syntax, which is automatically converted to Promise chains.

---

## Promise Basics

### Creating a Promise

```javascript
const promise = new Promise(function(resolve, reject) {
    // Async operation
    setTimeout(function() {
        const success = true;
        if (success) {
            resolve('Operation succeeded!');
        } else {
            reject('Operation failed!');
        }
    }, 1000);
});
```

### Promise States

A Promise can be in one of three states:

1. **Pending** - Initial state, neither fulfilled nor rejected
2. **Fulfilled** - The operation completed successfully
3. **Rejected** - The operation failed

Once a Promise is fulfilled or rejected, it becomes **settled** and cannot change states.

---

## Promise Methods

### `promise.then(onFulfilled, onRejected)`

Attach callbacks to handle the fulfilled or rejected state.

**Parameters:**
- `onFulfilled` (function) - Called when promise is fulfilled `(value) => {}`
- `onRejected` (function) - Optional, called when promise is rejected `(reason) => {}`

**Returns:** A new Promise for chaining

**Example:**
```javascript
promise.then(
    function(value) {
        console.log('Success:', value);
    },
    function(error) {
        console.error('Error:', error);
    }
);
```

---

### `promise.catch(onRejected)`

Attach an error handler. Equivalent to `.then(null, onRejected)`.

**Parameters:**
- `onRejected` (function) - Called when promise is rejected `(reason) => {}`

**Returns:** A new Promise for chaining

**Example:**
```javascript
promise
    .then(function(value) {
        console.log('Success:', value);
    })
    .catch(function(error) {
        console.error('Error:', error);
    });
```

---

## Promise Chaining

Promises can be chained to perform sequential async operations:

```javascript
new Promise(function(resolve, reject) {
    resolve(1);
})
.then(function(value) {
    console.log(value); // 1
    return value + 1;
})
.then(function(value) {
    console.log(value); // 2
    return value + 1;
})
.then(function(value) {
    console.log(value); // 3
});
```

### Returning Promises in Chains

If you return a Promise from a `.then()` handler, the next `.then()` will wait for that Promise to settle:

```javascript
file.read('config.json', function(err, data) {
    if (err) throw err;
    console.log('Config loaded');
})
// Using Promises
new Promise(function(resolve, reject) {
    file.read('config.json', function(err, data) {
        if (err) reject(err);
        else resolve(data);
    });
})
.then(function(data) {
    console.log('Config loaded');
    return new Promise(function(resolve, reject) {
        file.read('users.json', function(err, data) {
            if (err) reject(err);
            else resolve(data);
        });
    });
})
.then(function(data) {
    console.log('Users loaded');
})
.catch(function(err) {
    console.error('Error:', err);
});
```

---

## Static Methods

### `Promise.resolve(value)`

Create a Promise that is immediately resolved with the given value.

**Parameters:**
- `value` - The value to resolve with

**Returns:** A fulfilled Promise

**Example:**
```javascript
Promise.resolve('Hello')
    .then(function(value) {
        console.log(value); // 'Hello'
    });

// Useful for converting values to Promises
const promise = Promise.resolve(42);
```

---

### `Promise.reject(reason)`

Create a Promise that is immediately rejected with the given reason.

**Parameters:**
- `reason` - The rejection reason (typically an error)

**Returns:** A rejected Promise

**Example:**
```javascript
Promise.reject(new Error('Something went wrong'))
    .catch(function(error) {
        console.error(error.message); // 'Something went wrong'
    });
```

---

### `Promise.all(promises)`

Wait for all promises to be fulfilled, or reject if any promise rejects.

**Parameters:**
- `promises` (array) - Array of Promises to wait for

**Returns:** A Promise that resolves with an array of all fulfilled values, or rejects with the first rejection reason

**Behavior:**
- Resolves when **all** promises are fulfilled
- Rejects immediately when **any** promise rejects
- Results are in the same order as input promises

**Example:**
```javascript
const promises = [
    Promise.resolve(1),
    Promise.resolve(2),
    Promise.resolve(3)
];

Promise.all(promises).then(function(values) {
    console.log(values); // [1, 2, 3]
});

// Practical example with file operations
Promise.all([
    new Promise(function(resolve, reject) {
        file.read('file1.txt', function(err, data) {
            if (err) reject(err);
            else resolve(data);
        });
    }),
    new Promise(function(resolve, reject) {
        file.read('file2.txt', function(err, data) {
            if (err) reject(err);
            else resolve(data);
        });
    })
]).then(function(files) {
    console.log('File 1:', files[0]);
    console.log('File 2:', files[1]);
}).catch(function(err) {
    console.error('Failed to load files:', err);
});
```

---

### `Promise.race(promises)`

Race multiple promises - resolve or reject with the first promise that settles.

**Parameters:**
- `promises` (array) - Array of Promises to race

**Returns:** A Promise that settles the same way as the first settled promise

**Behavior:**
- Resolves/rejects as soon as the **first** promise settles
- Only the first result matters; other promises are ignored

**Example:**
```javascript
const fast = new Promise(function(resolve) {
    setTimeout(function() { resolve('Fast'); }, 100);
});

const slow = new Promise(function(resolve) {
    setTimeout(function() { resolve('Slow'); }, 1000);
});

Promise.race([fast, slow]).then(function(value) {
    console.log(value); // 'Fast' (after 100ms)
});

// Practical example: timeout pattern
function withTimeout(promise, ms) {
    const timeout = new Promise(function(resolve, reject) {
        setTimeout(function() {
            reject(new Error('Timeout after ' + ms + 'ms'));
        }, ms);
    });
    
    return Promise.race([promise, timeout]);
}

withTimeout(
    new Promise(function(resolve) {
        setTimeout(function() { resolve('Done'); }, 5000);
    }),
    2000
).catch(function(err) {
    console.error(err.message); // 'Timeout after 2000ms'
});
```

---

### (TODO) `Promise.allSettled(promises)`

Wait for all promises to settle (either fulfilled or rejected).

**Parameters:**
- `promises` (array) - Array of Promises to wait for

**Returns:** A Promise that resolves with an array of result objects

**Result Object:**
- `{ status: 'fulfilled', value: <value> }` for fulfilled promises
- `{ status: 'rejected', reason: <reason> }` for rejected promises

**Behavior:**
- Never rejects
- Always waits for all promises to settle
- Returns status and value/reason for each promise

**Example:**
```javascript
const promises = [
    Promise.resolve(1),
    Promise.reject('Error'),
    Promise.resolve(3)
];

Promise.allSettled(promises).then(function(results) {
    results.forEach(function(result, index) {
        if (result.status === 'fulfilled') {
            console.log('Promise', index, 'fulfilled:', result.value);
        } else {
            console.log('Promise', index, 'rejected:', result.reason);
        }
    });
});
// Output:
// Promise 0 fulfilled: 1
// Promise 1 rejected: Error
// Promise 2 fulfilled: 3
```

---

### (TODO) `Promise.any(promises)`

Wait for the first promise to fulfill. Reject only if all promises reject.

**Parameters:**
- `promises` (array) - Array of Promises to wait for

**Returns:** A Promise that fulfills with the first fulfilled value, or rejects if all promises reject

**Behavior:**
- Resolves with the **first fulfilled** value
- Rejects only if **all** promises reject
- Ignores rejections as long as at least one promise fulfills

**Example:**
```javascript
const promises = [
    Promise.reject('Error 1'),
    new Promise(function(resolve) {
        setTimeout(function() { resolve('Success'); }, 100);
    }),
    Promise.reject('Error 2')
];

Promise.any(promises).then(function(value) {
    console.log(value); // 'Success'
});

// If all reject
Promise.any([
    Promise.reject('Error 1'),
    Promise.reject('Error 2')
]).catch(function(err) {
    console.error('All promises rejected');
});
```

---

## Async/Await

Dougless supports `async/await` syntax through automatic transpilation with esbuild!

### Async Functions

```javascript
async function fetchData() {
    const response = await new Promise(function(resolve, reject) {
        http.get('https://api.example.com/data', function(err, res) {
            if (err) reject(err);
            else resolve(res);
        });
    });
    
    return JSON.parse(response.body);
}

fetchData().then(function(data) {
    console.log('Data:', data);
}).catch(function(err) {
    console.error('Error:', err);
});
```

### Await Keyword

Use `await` to wait for promises inside async functions:

```javascript
async function loadFiles() {
    try {
        const file1 = await new Promise(function(resolve, reject) {
            file.read('file1.txt', function(err, data) {
                if (err) reject(err);
                else resolve(data);
            });
        });
        
        const file2 = await new Promise(function(resolve, reject) {
            file.read('file2.txt', function(err, data) {
                if (err) reject(err);
                else resolve(data);
            });
        });
        
        console.log('File 1:', file1);
        console.log('File 2:', file2);
    } catch (err) {
        console.error('Error loading files:', err);
    }
}

loadFiles();
```

---

## Complete Examples

### Example 1: Sequential File Processing

```javascript
// Read a file, process it, then write the result
new Promise(function(resolve, reject) {
    file.read('input.txt', function(err, data) {
        if (err) reject(err);
        else resolve(data);
    });
})
.then(function(data) {
    // Process the data
    const processed = data.toUpperCase();
    return new Promise(function(resolve, reject) {
        file.write('output.txt', processed, function(err) {
            if (err) reject(err);
            else resolve(processed);
        });
    });
})
.then(function(result) {
    console.log('Processing complete!');
    console.log('Result:', result);
})
.catch(function(err) {
    console.error('Error:', err);
});
```

### Example 2: Parallel HTTP Requests

```javascript
// Fetch multiple APIs in parallel
const apis = [
    'https://api.example.com/users',
    'https://api.example.com/posts',
    'https://api.example.com/comments'
];

const requests = apis.map(function(url) {
    return new Promise(function(resolve, reject) {
        http.get(url, function(err, response) {
            if (err) reject(err);
            else resolve(JSON.parse(response.body));
        });
    });
});

Promise.all(requests)
    .then(function(results) {
        console.log('Users:', results[0]);
        console.log('Posts:', results[1]);
        console.log('Comments:', results[2]);
    })
    .catch(function(err) {
        console.error('Failed to fetch data:', err);
    });
```

### Example 3: Retry Logic with Promises

```javascript
function retry(fn, maxAttempts) {
    return new Promise(function(resolve, reject) {
        let attempts = 0;
        
        function attempt() {
            attempts++;
            fn().then(resolve).catch(function(err) {
                if (attempts >= maxAttempts) {
                    reject(err);
                } else {
                    console.log('Attempt', attempts, 'failed, retrying...');
                    setTimeout(attempt, 1000);
                }
            });
        }
        
        attempt();
    });
}

// Usage
retry(
    function() {
        return new Promise(function(resolve, reject) {
            http.get('https://unreliable-api.com/data', function(err, res) {
                if (err) reject(err);
                else resolve(res);
            });
        });
    },
    3 // Max 3 attempts
).then(function(response) {
    console.log('Success:', response.body);
}).catch(function(err) {
    console.error('Failed after 3 attempts:', err);
});
```

### Example 4: Promise-based Wrapper for Callbacks

```javascript
// Helper to convert callback-based APIs to Promises
function promisify(fn) {
    return function() {
        const args = Array.prototype.slice.call(arguments);
        return new Promise(function(resolve, reject) {
            args.push(function(err, result) {
                if (err) reject(err);
                else resolve(result);
            });
            fn.apply(null, args);
        });
    };
}

// Create promise-based versions
const readFile = promisify(file.read);
const writeFile = promisify(file.write);

// Use with async/await
async function copyFile(source, dest) {
    try {
        const data = await readFile(source);
        await writeFile(dest, data);
        console.log('File copied successfully!');
    } catch (err) {
        console.error('Copy failed:', err);
    }
}

copyFile('source.txt', 'destination.txt');
```

### Example 5: Waterfall Pattern

```javascript
// Execute promises in sequence, passing results forward
function waterfall(tasks) {
    return tasks.reduce(function(promise, task) {
        return promise.then(task);
    }, Promise.resolve());
}

// Usage
waterfall([
    function() {
        console.log('Step 1');
        return Promise.resolve('Result 1');
    },
    function(prev) {
        console.log('Step 2, got:', prev);
        return Promise.resolve('Result 2');
    },
    function(prev) {
        console.log('Step 3, got:', prev);
        return Promise.resolve('Result 3');
    }
]).then(function(final) {
    console.log('Final result:', final);
});
```

---

## Error Handling Best Practices

### Always Use .catch()

```javascript
// Good - errors are handled
promise
    .then(function(value) {
        return processValue(value);
    })
    .catch(function(err) {
        console.error('Error:', err);
    });

// Bad - unhandled rejection
promise.then(function(value) {
    return processValue(value);
});
```

### Catching Errors in Chains

```javascript
promise
    .then(function(value) {
        // If this throws, .catch() will handle it
        throw new Error('Processing failed');
    })
    .then(function(value) {
        // This won't run
        console.log(value);
    })
    .catch(function(err) {
        console.error('Caught:', err.message); // 'Processing failed'
    });
```

### Try/Catch with Async/Await

```javascript
async function safeOperation() {
    try {
        const result = await riskyOperation();
        return result;
    } catch (err) {
        console.error('Operation failed:', err);
        return null; // Return default value
    }
}
```

---

## Performance Tips

1. **Use Promise.all() for Parallel Operations**: When operations don't depend on each other, run them in parallel
   ```javascript
   // Parallel (faster)
   Promise.all([op1(), op2(), op3()]);
   
   // Sequential (slower)
   await op1();
   await op2();
   await op3();
   ```

2. **Avoid Creating Unnecessary Promises**: Don't wrap values that are already promises
   ```javascript
   // Bad
   return new Promise(resolve => resolve(existingPromise));
   
   // Good
   return existingPromise;
   ```

3. **Use Promise.allSettled() When You Need All Results**: Even if some operations fail
   ```javascript
   const results = await Promise.allSettled([
       fetchUser(),
       fetchPosts(),
       fetchComments() // Even if this fails, we get the others
   ]);
   ```

---

## Integration with Event Loop

Promises in Dougless are fully integrated with the event loop:

- Promise handlers are executed asynchronously
- `.then()` and `.catch()` callbacks are scheduled on the event loop
- The event loop waits for all pending promises before exiting
- Thread-safe state management with mutex protection

This means your script will automatically wait for all promises to settle before terminating!

---

## See Also

- [HTTP API Guide](http_api.md) - Using HTTP with promises
- [File API Guide](file_api.md) - File operations with promises
- [REPL Guide](repl_guide.md) - Testing promises interactively
- [Testing Guide](testing_guide.md) - Writing tests with promises
