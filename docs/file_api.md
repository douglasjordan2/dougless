# File API Guide

## Overview

Dougless provides a unique global `file` API for file system operations. Unlike Node.js which requires `require('fs')`, the `file` object is always available globally.

## Why Global?

**Dougless Philosophy**: File operations are so fundamental to runtime operations that they should be as accessible as `console`. This makes Dougless code cleaner and more intuitive.

**Comparison:**

```javascript
// Node.js
const fs = require('fs');
fs.readFile('data.txt', callback);

// Dougless
file.read('data.txt', callback);  // No require!
```

---

## File Operations

### `file.read(path, callback)`

Read the entire contents of a file asynchronously.

**Parameters:**
- `path` (string) - Path to the file
- `callback` (function) - Callback function `(err, data)`

**Example:**
```javascript
file.read('data.txt', function(err, data) {
    if (err) {
        console.error('Error reading file:', err);
        return;
    }
    console.log('File contents:', data);
});
```

---

### `file.write(path, data, callback)`

Write data to a file asynchronously. Creates the file if it doesn't exist, overwrites if it does.

**Parameters:**
- `path` (string) - Path to the file
- `data` (string) - Data to write
- `callback` (function) - Callback function `(err)`

**Example:**
```javascript
file.write('output.txt', 'Hello Dougless!', function(err) {
    if (err) {
        console.error('Error writing file:', err);
        return;
    }
    console.log('File written successfully');
});
```

---

### `file.readdir(path, callback)`

Read the contents of a directory.

**Parameters:**
- `path` (string) - Path to the directory
- `callback` (function) - Callback function `(err, files)`
  - `files` is an array of filenames (strings)

**Example:**
```javascript
file.readdir('.', function(err, files) {
    if (err) {
        console.error('Error reading directory:', err);
        return;
    }
    console.log('Files:', files);
    // Files: ["file1.txt", "file2.txt", "subdir"]
});
```

---

### `file.exists(path, callback)`

Check if a file or directory exists.

**Parameters:**
- `path` (string) - Path to check
- `callback` (function) - Callback function `(exists)`
  - `exists` is a boolean (no error parameter)

**Example:**
```javascript
file.exists('data.txt', function(exists) {
    if (exists) {
        console.log('File exists!');
    } else {
        console.log('File does not exist');
    }
});
```

**Note**: Unlike other file operations, `exists()` doesn't pass an error - just a boolean.

---

### `file.mkdir(path, callback)`

Create a new directory.

**Parameters:**
- `path` (string) - Path for the new directory
- `callback` (function) - Callback function `(err)`

**Permissions**: Creates with `0755` (rwxr-xr-x)

**Example:**
```javascript
file.mkdir('new-folder', function(err) {
    if (err) {
        console.error('Error creating directory:', err);
        return;
    }
    console.log('Directory created');
});
```

---

### `file.rmdir(path, callback)`

Remove an empty directory.

**Parameters:**
- `path` (string) - Path to the directory
- `callback` (function) - Callback function `(err)`

**Note**: Directory must be empty, or an error will occur.

**Example:**
```javascript
file.rmdir('old-folder', function(err) {
    if (err) {
        console.error('Error removing directory:', err);
        return;
    }
    console.log('Directory removed');
});
```

---

### `file.unlink(path, callback)`

Delete a file.

**Parameters:**
- `path` (string) - Path to the file
- `callback` (function) - Callback function `(err)`

**Example:**
```javascript
file.unlink('temp.txt', function(err) {
    if (err) {
        console.error('Error deleting file:', err);
        return;
    }
    console.log('File deleted');
});
```

---

### `file.stat(path, callback)`

Get information about a file or directory.

**Parameters:**
- `path` (string) - Path to the file/directory
- `callback` (function) - Callback function `(err, stats)`
  - `stats` object contains:
    - `size` (number) - Size in bytes
    - `isDirectory` (boolean) - True if directory
    - `isFile` (boolean) - True if regular file
    - `modified` (number) - Unix timestamp of last modification
    - `name` (string) - Base name of the file/directory

**Example:**
```javascript
file.stat('data.txt', function(err, stats) {
    if (err) {
        console.error('Error:', err);
        return;
    }
    
    console.log('Size:', stats.size, 'bytes');
    console.log('Is file:', stats.isFile);
    console.log('Is directory:', stats.isDirectory);
    console.log('Modified:', new Date(stats.modified * 1000));
    console.log('Name:', stats.name);
});
```

---

## Complete Examples

### Example 1: Read and Process File

```javascript
file.read('input.txt', function(err, data) {
    if (err) {
        console.error('Cannot read file:', err);
        return;
    }
    
    // Process the data
    var processed = data.toUpperCase();
    
    // Write to output
    file.write('output.txt', processed, function(err) {
        if (err) {
            console.error('Cannot write file:', err);
        } else {
            console.log('Processing complete!');
        }
    });
});
```

### Example 2: Create Directory Structure

```javascript
// Create a directory
file.mkdir('project', function(err) {
    if (err) {
        console.error('Error:', err);
        return;
    }
    
    // Create subdirectories
    file.mkdir('project/src', function(err) {
        if (err) console.error(err);
    });
    
    file.mkdir('project/docs', function(err) {
        if (err) console.error(err);
    });
    
    // Create files
    file.write('project/README.md', '# My Project', function(err) {
        if (err) console.error(err);
        else console.log('Project structure created!');
    });
});
```

### Example 3: Directory Cleanup

```javascript
// List directory contents
file.readdir('temp', function(err, files) {
    if (err) {
        console.error('Error:', err);
        return;
    }
    
    // Delete each file
    var remaining = files.length;
    files.forEach(function(filename) {
        file.unlink('temp/' + filename, function(err) {
            if (err) console.error('Error deleting', filename, err);
            remaining--;
            
            // When all files are deleted, remove the directory
            if (remaining === 0) {
                file.rmdir('temp', function(err) {
                    if (err) console.error('Error removing dir:', err);
                    else console.log('Cleanup complete!');
                });
            }
        });
    });
});
```

### Example 4: Check Before Writing

```javascript
file.exists('config.json', function(exists) {
    if (exists) {
        console.log('Config exists, backing up...');
        file.read('config.json', function(err, data) {
            if (!err) {
                file.write('config.json.backup', data, function(err) {
                    if (!err) console.log('Backup created');
                });
            }
        });
    } else {
        console.log('Creating new config...');
        var defaultConfig = '{"version": "1.0"}';
        file.write('config.json', defaultConfig, function(err) {
            if (!err) console.log('Config created');
        });
    }
});
```

---

## Error Handling

All operations except `exists()` follow the Node.js error-first callback pattern:

```javascript
function callback(err, result) {
    if (err) {
        // Error occurred
        console.error('Error:', err);
        return;
    }
    // Success - use result
    console.log('Result:', result);
}
```

**Common Errors:**
- File not found
- Permission denied
- Directory not empty (for `rmdir`)
- Path already exists (for `mkdir`)

---

## Limitations

### Current Limitations (Phase 2)
- No synchronous operations (all are async)
- No streaming (reads entire file into memory)
- No append mode (write always overwrites)
- No file permissions control beyond default
- `rmdir` only removes empty directories

### Future Enhancements (Phase 3+)
- Streaming support for large files
- Append operations
- File watching
- Advanced permissions control
- Recursive directory operations

---

## Comparison with Node.js

| Feature | Node.js | Dougless |
|---------|---------|----------|
| **Import** | `require('fs')` | Global (no require) |
| **Read** | `fs.readFile()` | `file.read()` |
| **Write** | `fs.writeFile()` | `file.write()` |
| **List dir** | `fs.readdir()` | `file.readdir()` |
| **Exists** | `fs.exists()` (deprecated) | `file.exists()` |
| **Make dir** | `fs.mkdir()` | `file.mkdir()` |
| **Remove dir** | `fs.rmdir()` | `file.rmdir()` |
| **Delete file** | `fs.unlink()` | `file.unlink()` |
| **File info** | `fs.stat()` | `file.stat()` |

---

## Best Practices

1. **Always handle errors** - Check the error parameter in callbacks
2. **Use stat before operations** - Check file type before reading/deleting
3. **Sequential operations** - Nest callbacks for dependent operations
4. **Check exists before create** - Avoid overwriting existing files/dirs
5. **Clean up temp files** - Always delete temporary files when done

---

**Happy file handling with Dougless!** 📁
