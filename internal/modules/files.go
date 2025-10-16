package modules

import (
	"context"
	"os"
  "path/filepath"
	"time"

	"github.com/dop251/goja"

	"github.com/douglasjordan2/dougless/internal/event"
	"github.com/douglasjordan2/dougless/internal/permissions"
)

// Files provides a simplified, convention-based file system API.
//
// Unlike traditional file APIs with dozens of methods, Files uses 3 smart methods
// with path-based conventions to handle all file system operations:
//   - files.read(path, [callback]) - Read files or list directories
//   - files.write(path, [content], [callback]) - Write files or create directories
//   - files.rm(path, [callback]) - Remove files or directories
//
// All methods support both callbacks and promises:
//   - With callback: Traditional error-first callback pattern
//   - Without callback: Returns a Promise for use with .then() or async/await
//
// Path conventions:
//   - Trailing '/' indicates directory operations
//   - No trailing '/' indicates file operations
//   - Parent directories are created automatically for file writes
//   - Content is optional - omit to create empty file (like 'touch' command)
//
// Example usage in JavaScript:
//
//	// Callback style - Read a file
//	files.read('data.txt', (err, content) => {
//	    if (content === null) console.log('File does not exist');
//	});
//
//	// Promise style - Read a file
//	const content = await files.read('data.txt');
//	if (content === null) console.log('File does not exist');
//
//	// List directory (trailing slash)
//	files.read('src/', (err, fileNames) => {
//	    console.log('Files:', fileNames);
//	});
//
//	// Write file with content
//	await files.write('data/output.txt', 'Hello');
//
//	// Create empty file (like 'touch')
//	await files.write('empty.txt');
//
//	// Create directory
//	await files.write('new-dir/');
//
//	// Remove file or directory
//	await files.rm('old.txt');
type Files struct {
	vm        *goja.Runtime // JavaScript runtime instance
	eventLoop *event.Loop   // Event loop for async task scheduling
}

// NewFiles creates a new Files instance with the given event loop.
// The event loop is used to schedule all async file operations.
func NewFiles(eventLoop *event.Loop) *Files {
	return &Files{
		eventLoop: eventLoop,
	}
}

// Export creates and returns the files object for use in JavaScript.
// The returned object provides three methods with dual callback/promise APIs:
//   - read(path, [callback]) - Smart read for files or directories
//   - write(path, [content], [callback]) - Smart write for files or directories
//   - rm(path, [callback]) - Unified removal for files or directories
//
// When callback is omitted, methods return Promises compatible with async/await.
func (fs *Files) Export(vm *goja.Runtime) goja.Value {
	fs.vm = vm
	obj := vm.NewObject()

	obj.Set("read", fs.read)
	obj.Set("write", fs.write)
	obj.Set("rm", fs.rm)

	return obj
}

func dirCheck(dest string) bool { 
  return len(dest) > 0 && dest[len(dest)-1] == '/' 
}

func (fs *Files) doRead(ctx context.Context, dest string) (goja.Value, goja.Value) {
  mgr := permissions.GetManager()
  canRead := permissions.PermissionRead
  if !mgr.CheckWithPrompt(ctx, canRead, dest) {
    errMsg := mgr.ErrorMessage(canRead, dest)
    return fs.vm.ToValue(errMsg), goja.Undefined()
  }

  isDir := dirCheck(dest)

  _, statErr := os.Stat(dest)
  if os.IsNotExist(statErr) {
    // Path doesn't exist - return null data (not an error)
    return goja.Null(), goja.Null()
  }

  var fileData []byte
  var dirData []os.DirEntry
  var err error

  if isDir {
    dirData, err = os.ReadDir(dest)
  } else {
    fileData, err = os.ReadFile(dest)
  }

  var errArg, dataArg goja.Value
  if err != nil {
    errArg = fs.vm.ToValue(err.Error())
    dataArg = goja.Undefined()
  } else {
    errArg = goja.Null()
    if isDir {
      names := make([]string, len(dirData))
      for i, entry := range dirData {
        names[i] = entry.Name()
      }
      dataArg = fs.vm.ToValue(names)
    } else {
      dataArg = fs.vm.ToValue(string(fileData))
    }
  }

  return errArg, dataArg
}

// read implements the files.read() method.
//
// Behavior:
//   - Trailing '/': Read directory, returns string[] of filenames
//   - No trailing '/': Read file, returns string content or null if doesn't exist
//
// Parameters:
//   - path (string): Path to file or directory
//   - callback (function, optional): Callback function (err, data)
//
// Return Value:
//   - With callback: undefined (result passed to callback)
//   - Without callback: Promise<string | string[] | null>
//
// JavaScript Usage:
//   // Callback style
//   files.read('data.txt', (err, content) => { ... });
//
//   // Promise style
//   const content = await files.read('data.txt');
//
// Returns null (not error) when file doesn't exist, perfect for existence checks.
// Requires PermissionRead for the specified path.
func (fs *Files) read(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 1 {
		panic(fs.vm.NewTypeError("read requires a file or directory path"))
	}

	dest := call.Arguments[0].String()

	var callback goja.Callable
	var ok bool
	if len(call.Arguments) > 1 {
		callback, ok = goja.AssertFunction(call.Arguments[1])
	}
	if !ok {
    promise := &Promise{
      vm:          fs.vm,
      eventLoop:   fs.eventLoop,
      state:       PromisePending,
      onFulfilled: []goja.Callable{},
      onRejected:  []goja.Callable{},
    }

    fs.eventLoop.ScheduleTask(&event.Task{
      Callback: func() {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        errArg, dataArg := fs.doRead(ctx, dest)

        if errArg != goja.Null() && !goja.IsNull(errArg) {
          promise.reject(errArg)
        } else {
          promise.resolve(dataArg)
        }
      },
    })

    return CreatePromiseObject(fs.vm, promise)
	}

	fs.eventLoop.ScheduleTask(&event.Task{
		Callback: func() {
      ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
      defer cancel()

      errArg, dataArg := fs.doRead(ctx, dest)
      callback(goja.Undefined(), errArg, dataArg)
		},
	})

	return goja.Undefined()
}

func (fs *Files) doWrite(ctx context.Context, dest string, data string) (goja.Value) {
  mgr := permissions.GetManager()
  canWrite := permissions.PermissionWrite
  if !mgr.CheckWithPrompt(ctx, canWrite, dest) {
    errMsg := mgr.ErrorMessage(canWrite, dest)
    return fs.vm.ToValue(errMsg)
  }

  isDir := dirCheck(dest)

  var err error
  if isDir {
    err = os.MkdirAll(dest, 0755) 
  } else {
    // Create parent directories if needed
    if mkdirErr := os.MkdirAll(filepath.Dir(dest), 0755); mkdirErr != nil {
      return fs.vm.ToValue(mkdirErr.Error())
    }
    err = os.WriteFile(dest, []byte(data), 0644)
  }

  var errArg goja.Value
  if err != nil {
    errArg = fs.vm.ToValue(err.Error())
  } else {
    errArg = goja.Null()
  }

  return errArg
}

// write implements the files.write() method.
//
// Behavior:
//   - Trailing '/': Create directory recursively
//   - No trailing '/': Write file (empty if no content provided)
//
// Parameters:
//   - path (string): Path to file or directory
//   - content (string, optional): Data to write (defaults to empty string for files)
//   - callback (function, optional): Callback function (err)
//
// Return Value:
//   - With callback: undefined (result passed to callback)
//   - Without callback: Promise<void>
//
// JavaScript Usage:
//   // Write file with content
//   await files.write('output.txt', 'data');
//   files.write('output.txt', 'data', (err) => { ... });
//
//   // Create empty file (like 'touch')
//   await files.write('empty.txt');
//   files.write('empty.txt', (err) => { ... });
//
//   // Create directory
//   await files.write('new-dir/');
//   files.write('new-dir/', (err) => { ... });
//
// Parent directories are created automatically for file writes.
// Requires PermissionWrite for the specified path.
func (fs *Files) write(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 1 {
		panic(fs.vm.NewTypeError("write requires at least a path"))
	}

	dest := call.Arguments[0].String()
  isDir := dirCheck(dest)
  data := ""  // Default to empty string for files

  var callback goja.Callable
  var ok bool
  
  if isDir {
    // Directory creation: write('path/', [callback])
    if len(call.Arguments) > 1 {
      callback, ok = goja.AssertFunction(call.Arguments[1])
    }
  } else {
    // File write: write('path', [content], [callback])
    // Check argument 1: could be content (string) or callback (function)
    if len(call.Arguments) > 1 {
      // Try to get as string first (content)
      if !goja.IsUndefined(call.Arguments[1]) && !goja.IsNull(call.Arguments[1]) {
        callback, ok = goja.AssertFunction(call.Arguments[1])
        if !ok {
          // It's content, not a callback
          data = call.Arguments[1].String()
          // Check for callback in argument 2
          if len(call.Arguments) > 2 {
            callback, ok = goja.AssertFunction(call.Arguments[2])
          }
        }
      }
    }
  }

  if !ok {
    // Promise path
    promise := &Promise{
      vm:          fs.vm,
      eventLoop:   fs.eventLoop,
      state:       PromisePending,
      onFulfilled: []goja.Callable{},
      onRejected:  []goja.Callable{},
    }

    fs.eventLoop.ScheduleTask(&event.Task{
      Callback: func() {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        errArg := fs.doWrite(ctx, dest, data)

        if errArg != goja.Null() && !goja.IsNull(errArg) {
          promise.reject(errArg)
        } else {
          promise.resolve(goja.Null())
        }
      },
    })

    return CreatePromiseObject(fs.vm, promise)
  }

  // Callback path
	fs.eventLoop.ScheduleTask(&event.Task{
		Callback: func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

      errArg := fs.doWrite(ctx, dest, data)
			callback(goja.Undefined(), errArg)
		},
	})

	return goja.Undefined()
}

func (fs *Files) doRm(ctx context.Context, path string) (goja.Value) {
  mgr := permissions.GetManager()
  canWrite := permissions.PermissionWrite
  if !mgr.CheckWithPrompt(ctx, canWrite, path) {
    errMsg := mgr.ErrorMessage(canWrite, path)
    return fs.vm.ToValue(errMsg)
  }

  err := os.RemoveAll(path)

  var errArg goja.Value
  if err != nil {
    errArg = fs.vm.ToValue(err.Error())
  } else {
    errArg = goja.Null()
  }

  return errArg
}

// rm implements the files.rm() method.
//
// Unified removal operation that works on both files and directories.
// Removes directories recursively (including contents).
// Idempotent - succeeds even if path doesn't exist.
//
// Parameters:
//   - path (string): Path to file or directory to remove
//   - callback (function, optional): Callback function (err)
//
// Return Value:
//   - With callback: undefined (result passed to callback)
//   - Without callback: Promise<void>
//
// JavaScript Usage:
//   // Callback style
//   files.rm('temp.txt', (err) => { ... });
//
//   // Promise style
//   await files.rm('temp.txt');
//
// Uses os.RemoveAll() under the hood for recursive removal.
// Requires PermissionWrite for the specified path.
func (fs *Files) rm(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 1 {
		panic(fs.vm.NewTypeError("rm requires a path"))
	}

	path := call.Arguments[0].String()
	var callback goja.Callable
	var ok bool
	if len(call.Arguments) > 1 {
		callback, ok = goja.AssertFunction(call.Arguments[1])
	}
	if !ok {
    // Promise path
    promise := &Promise{
      vm:          fs.vm,
      eventLoop:   fs.eventLoop,
      state:       PromisePending,
      onFulfilled: []goja.Callable{},
      onRejected:  []goja.Callable{},
    }

    fs.eventLoop.ScheduleTask(&event.Task{
      Callback: func() {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        errArg := fs.doRm(ctx, path)

        if errArg != goja.Null() && !goja.IsNull(errArg) {
          promise.reject(errArg)
        } else {
          promise.resolve(goja.Null())
        }
      },
    })

    return CreatePromiseObject(fs.vm, promise)
	}

  // Callback path
	fs.eventLoop.ScheduleTask(&event.Task{
		Callback: func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

      errArg := fs.doRm(ctx, path)

			callback(goja.Undefined(), errArg)
		},
	})

	return goja.Undefined()
}
