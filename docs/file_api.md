# Files API Guide

## Overview

Dougless provides a unique global `files` API for file system operations. Unlike Node.js which requires `require('fs')` and has dozens of methods, Dougless uses **3 smart methods** with **convention-based routing**.

## Why This Design?

**Dougless Philosophy**: 
1. **Global Access** - File operations are fundamental, should be as accessible as `console`
2. **Convention Over Configuration** - Use path patterns (trailing `/`) instead of separate methods
3. **Simplicity** - 3 methods instead of 8+ reduces cognitive load
4. **Smart Defaults** - Auto-create parent directories, graceful null handling

**Comparison:**

```javascript
// Node.js - requires import and multiple methods
const fs = require('fs');
fs.readFile('data.txt', callback);
fs.readdir('src', callback);
fs.mkdir('dir', callback);
fs.unlink('file.txt', callback);

// Dougless - global, 3 methods, convention-based
files.read('data.txt', callback);     // Read file
files.read('src/', callback);         // Read directory (trailing /)
files.write('dir/', callback);        // Create directory
files.rm('file.txt', callback);       // Delete anything
```

---

## Core Methods

### `files.read(path, [callback])`

**Smart read operation** - behavior depends on the path:

**Parameters:**
- `path` (string) - Path to file or directory
  - **No trailing `/`**: Read file contents
  - **Trailing `/`**: List directory contents
- `callback` (function, optional) - Callback function `(err, data)`
  - For files: `data` is `string` content (or `null` if file doesn't exist)
  - For directories: `data` is `string[]` array of filenames
  - **If omitted**: Returns a Promise

**Return Value:**
- Without callback: `Promise<string | string[] | null>`
- With callback: `undefined`

**Examples:**

```javascript
// Callback style - Read a file
files.read('data.txt', function(err, content) {
    if (err) {
        console.error('Error:', err);
        return;
    }
    if (content === null) {
        console.log('File does not exist');
    } else {
        console.log('File contents:', content);
    }
});

// Promise style - Read a file
files.read('data.txt')
    .then(content => {
        if (content === null) {
            console.log('File does not exist');
        } else {
            console.log('File contents:', content);
        }
    })
    .catch(err => console.error('Error:', err));

// Async/await style - Read a file
async function readFile() {
    try {
        const content = await files.read('data.txt');
        if (content === null) {
            console.log('File does not exist');
        } else {
            console.log('File contents:', content);
        }
    } catch (err) {
        console.error('Error:', err);
    }
}

// Read a directory (note the trailing slash)
files.read('src/', function(err, fileNames) {
    if (err) {
        console.error('Error:', err);
        return;
    }
    console.log('Files in src/:', fileNames);
    // ["app.js", "utils.js", "config.json"]
});

// Promise style - Read a directory
const fileNames = await files.read('src/');
console.log('Files in src/:', fileNames);

// Check if file exists (null = doesn't exist)
files.read('config.json', function(err, data) {
    if (data === null) {
        console.log('Config file missing - creating default...');
        files.write('config.json', '{}', function(err) {
            if (!err) console.log('Created!');
        });
    }
});

// Async/await - Check if file exists
if (await files.read('config.json') === null) {
    console.log('Config file missing - creating default...');
    await files.write('config.json', '{}');
    console.log('Created!');
}
```

**Key Features:**
- Returns `null` (not error) when file doesn't exist - perfect for existence checks
- Trailing `/` convention makes directory operations explicit
- Single method replaces: `file.read()`, `file.readdir()`, `file.exists()`

---

### `files.write(path, [content], [callback])`

**Smart write operation** - behavior depends on arguments and path:

**Parameters:**
- `path` (string) - Path to file or directory
- `content` (string, optional) - Data to write (omit for directory creation)
- `callback` (function, optional) - Callback function `(err)`
  - **If omitted**: Returns a Promise

**Modes:**
- **2 args** with trailing `/`: Create directory
- **3 args**: Write file (auto-creates parent directories)

**Return Value:**
- Without callback: `Promise<void>`
- With callback: `undefined`

**Examples:**

```javascript
// Callback style - Write a file
files.write('output.txt', 'Hello Dougless!', function(err) {
    if (err) {
        console.error('Error:', err);
        return;
    }
    console.log('File written successfully');
});

// Promise style - Write a file
files.write('output.txt', 'Hello Dougless!')
    .then(() => console.log('File written successfully'))
    .catch(err => console.error('Error:', err));

// Async/await style - Write a file
try {
    await files.write('output.txt', 'Hello Dougless!');
    console.log('File written successfully');
} catch (err) {
    console.error('Error:', err);
}

// Write to nested path (auto-creates parent dirs)
files.write('data/users/profile.json', '{"name":"Alice"}', function(err) {
    if (!err) console.log('Created data/ and users/ directories automatically!');
});

// Async/await - Write to nested path
await files.write('data/users/profile.json', '{"name":"Alice"}');
console.log('Created data/ and users/ directories automatically!');

// Create a directory
files.write('project/', function(err) {
    if (err) {
        console.error('Error:', err);
        return;
    }
    console.log('Directory created');
});

// Async/await - Create a directory
await files.write('project/');
console.log('Directory created');

// Create nested directories
files.write('src/components/buttons/', function(err) {
    if (!err) console.log('All directories created!');
});

// Async/await - Create nested directories
await files.write('src/components/buttons/');
console.log('All directories created!');
```

**Key Features:**
- Automatically creates parent directories for file writes
- Trailing `/` convention for directories
- Single method replaces: `file.write()`, `file.mkdir()`
- No need to manually create directory structure

---

### `files.rm(path, [callback])`

**Unified removal** - deletes files or directories (recursively).

**Parameters:**
- `path` (string) - Path to file or directory to remove
- `callback` (function, optional) - Callback function `(err)`
  - **If omitted**: Returns a Promise

**Return Value:**
- Without callback: `Promise<void>`
- With callback: `undefined`

**Examples:**

```javascript
// Callback style - Delete a file
files.rm('temp.txt', function(err) {
    if (err) {
        console.error('Error:', err);
        return;
    }
    console.log('File deleted');
});

// Promise style - Delete a file
files.rm('temp.txt')
    .then(() => console.log('File deleted'))
    .catch(err => console.error('Error:', err));

// Async/await style - Delete a file
try {
    await files.rm('temp.txt');
    console.log('File deleted');
} catch (err) {
    console.error('Error:', err);
}

// Delete a directory (recursively, even if not empty)
files.rm('old-project/', function(err) {
    if (!err) console.log('Directory and all contents removed');
});

// Async/await - Delete a directory
await files.rm('old-project/');
console.log('Directory and all contents removed');

// Idempotent - no error if path doesn't exist
files.rm('maybe-exists.txt', function(err) {
    // Will succeed even if file doesn't exist
    if (!err) console.log('Removed (or was already gone)');
});

// Async/await - Idempotent removal
await files.rm('maybe-exists.txt');
console.log('Removed (or was already gone)');
```

**Key Features:**
- Works on files AND directories (no separate `rmdir`)
- Recursive deletion - removes directories with contents
- Idempotent - gracefully handles non-existent paths
- Single method replaces: `file.unlink()`, `file.rmdir()`

---

## Complete Examples

### Example 1: Read and Process File

**Callback Style:**
```javascript
files.read('input.txt', function(err, data) {
    if (err) {
        console.error('Cannot read file:', err);
        return;
    }
    
    if (data === null) {
        console.error('File does not exist');
        return;
    }
    
    // Process the data
    const processed = data.toUpperCase();
    
    // Write to output (auto-creates parent dirs if needed)
    files.write('output.txt', processed, function(err) {
        if (err) {
            console.error('Cannot write file:', err);
        } else {
            console.log('Processing complete!');
        }
    });
});
```

**Async/Await Style (Much Cleaner!):**
```javascript
async function processFile() {
    try {
        const data = await files.read('input.txt');
        
        if (data === null) {
            console.error('File does not exist');
            return;
        }
        
        // Process the data
        const processed = data.toUpperCase();
        
        // Write to output (auto-creates parent dirs if needed)
        await files.write('output.txt', processed);
        console.log('Processing complete!');
    } catch (err) {
        console.error('Error:', err);
    }
}

processFile();
```

### Example 2: Create Directory Structure (Simplified!)

```javascript
// Old way: multiple mkdir calls, manual nesting
// NEW WAY: Just write files, directories auto-created!

files.write('project/src/app.js', 'console.log("Hello");', function(err) {
    if (!err) console.log('Created project/src/ and app.js!');
});

files.write('project/docs/README.md', '# My Project', function(err) {
    if (!err) console.log('Created project/docs/ and README.md!');
});

// Or create empty directories explicitly
files.write('project/tests/', function(err) {
    if (!err) console.log('Created tests/ directory');
});
```

### Example 3: Directory Cleanup (Much Simpler!)

```javascript
// Old way: list files, delete each, then remove directory
// NEW WAY: Just remove the directory (recursive)

files.rm('temp/', function(err) {
    if (err) {
        console.error('Error:', err);
    } else {
        console.log('Cleanup complete! (removed directory and all contents)');
    }
});

// Or list files first if you need to
files.read('temp/', function(err, fileNames) {
    if (err) {
        console.error('Error:', err);
        return;
    }
    
    console.log('About to delete:', fileNames);
    
    // Remove the entire directory
    files.rm('temp/', function(err) {
        if (!err) console.log('Deleted!');
    });
});
```

### Example 4: Check Before Writing

```javascript
files.read('config.json', function(err, data) {
    if (err) {
        console.error('Error:', err);
        return;
    }
    
    if (data !== null) {
        // File exists, back it up
        console.log('Config exists, backing up...');
        files.write('config.json.backup', data, function(err) {
            if (!err) console.log('Backup created');
        });
    } else {
        // File doesn't exist, create default
        console.log('Creating new config...');
        const defaultConfig = '{"version": "1.0"}';
        files.write('config.json', defaultConfig, function(err) {
            if (!err) console.log('Config created');
        });
    }
});
```

### Example 5: Build Tool Pattern

```javascript
// Clean build directory
files.rm('dist/', function(err) {
    if (err) {
        console.error('Clean failed:', err);
        return;
    }
    
    // Build outputs (directories auto-created)
    files.write('dist/js/bundle.js', '/* bundled code */', function(err) {
        if (!err) console.log('Built JS');
    });
    
    files.write('dist/css/styles.css', '/* styles */', function(err) {
        if (!err) console.log('Built CSS');
    });
    
    files.write('dist/index.html', '<html>...</html>', function(err) {
        if (!err) console.log('Built HTML');
    });
});
```

---

## Error Handling

All operations follow the error-first callback pattern:

```javascript
function callback(err, result) {
    if (err) {
        // Error occurred
        console.error('Error:', err);
        return;
    }
    
    // Success - handle result
    // Note: result can be null for files.read() if file doesn't exist
    if (result === null) {
        console.log('File does not exist');
    } else {
        console.log('Result:', result);
    }
}
```

**Common Errors:**
- Permission denied (requires `--allow-read` or `--allow-write`)
- I/O errors (disk full, read errors, etc.)
- Invalid paths

**Note on `files.read()` Behavior:**
- Returns `null` (not an error) when file doesn't exist
- Only returns error for actual I/O failures or permission issues

---

## Limitations

### Current Limitations
- No synchronous operations (all are async)
- No streaming (reads entire file into memory)
- No append mode (write always overwrites)
- No file permissions control beyond default (0644 files, 0755 directories)
- No `stat()` method for file metadata (may return in future)

### Future Enhancements
- Streaming support for large files
- Append operations (`files.append()`)
- File watching (`files.watch()`)
- Advanced permissions control
- Optional `files.stat()` for metadata queries

---

## Comparison with Node.js

| Feature | Node.js | Dougless |
|---------|---------|----------|
| **Import** | `require('fs')` | Global (no require) |
| **Read file** | `fs.readFile()` | `files.read(path, cb)` or `await files.read(path)` |
| **Write file** | `fs.writeFile()` | `files.write(path, data, cb)` or `await files.write(path, data)` |
| **List dir** | `fs.readdir()` | `files.read(path + '/', cb)` or `await files.read(path + '/')` |
| **Exists** | `fs.exists()` (deprecated) | `files.read()` returns `null` |
| **Make dir** | `fs.mkdir()` | `files.write(path + '/', cb)` or `await files.write(path + '/')` |
| **Make dirs (recursive)** | `fs.mkdir({recursive: true})` | **Auto!** `files.write()` |
| **Remove dir** | `fs.rmdir()` | `files.rm(path, cb)` or `await files.rm(path)` |
| **Remove recursive** | `fs.rm({recursive: true})` | **Default!** `files.rm()` |
| **Delete file** | `fs.unlink()` | `files.rm(path, cb)` or `await files.rm(path)` |
| **File info** | `fs.stat()` | Not available (removed) |
| **Promises** | `require('fs/promises')` | **Built-in!** Just omit callback |
| **Method count** | 50+ methods | **3 methods** |

---

## Promise Support Details

All file operations return promises when the callback parameter is omitted:

**How it works:**
- Callback provided → Standard callback behavior, returns `undefined`
- Callback omitted → Returns a Promise
- Promise resolves with data (for `read()`) or `null` (for `write()`/`rm()`)
- Promise rejects with error message string

**Benefits:**
- ✅ **Cleaner code** - async/await is more readable than nested callbacks
- ✅ **Error handling** - try/catch blocks instead of error-first callbacks
- ✅ **Composition** - Use Promise.all() for parallel operations
- ✅ **No breaking changes** - Existing callback code continues to work

---

## Best Practices

1. **Prefer async/await** - Cleaner and easier to read than callbacks
2. **Always handle errors** - Use try/catch with async/await or .catch() with promises
3. **Check for null** - `files.read()` returns `null` when file doesn't exist
4. **Use trailing `/` for directories** - Makes intent explicit: `files.read('src/')` vs `files.read('src')`
5. **No need to create dirs first** - `files.write()` auto-creates parent directories
6. **Use `files.rm()` for everything** - Works on files and directories, recursive by default
7. **Parallel operations** - Use `Promise.all()` to read/write multiple files concurrently
8. **Check exists before create** - Avoid overwriting existing files/dirs
9. **Clean up temp files** - Always delete temporary files when done

---

**Happy file handling with Dougless!** 📁
