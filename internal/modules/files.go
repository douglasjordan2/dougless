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
//   - files.read(path, callback) - Read files or list directories
//   - files.write(path, [content], callback) - Write files or create directories
//   - files.rm(path, callback) - Remove files or directories
//
// Path conventions:
//   - Trailing '/' indicates directory operations
//   - No trailing '/' indicates file operations
//   - Parent directories are created automatically for file writes
//
// Example usage in JavaScript:
//
//	// Read a file
//	files.read('data.txt', (err, content) => {
//	    if (content === null) console.log('File does not exist');
//	});
//
//	// List directory
//	files.read('src/', (err, fileNames) => {
//	    console.log('Files:', fileNames);
//	});
//
//	// Write file (auto-creates parent dirs)
//	files.write('data/output.txt', 'Hello', (err) => {});
//
//	// Create directory
//	files.write('new-dir/', (err) => {});
//
//	// Remove file or directory
//	files.rm('old.txt', (err) => {});
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

// dirCheck determines if a path represents a directory based on trailing slash.
// Returns true if the path ends with '/', false otherwise.
func dirCheck(dest string) bool { 
  return len(dest) > 0 && dest[len(dest)-1] == '/' 
}

// Export creates and returns the files object for use in JavaScript.
// The returned object provides three methods:
//   - read(path, callback) - Smart read for files or directories
//   - write(path, [content], callback) - Smart write for files or directories
//   - rm(path, callback) - Unified removal for files or directories
func (fs *Files) Export(vm *goja.Runtime) goja.Value {
	fs.vm = vm
	obj := vm.NewObject()

	obj.Set("read", fs.read)
	obj.Set("write", fs.write)
	obj.Set("rm", fs.rm)

	return obj
}

// read implements the files.read() method.
//
// Behavior:
//   - Trailing '/': Read directory, returns string[] of filenames
//   - No trailing '/': Read file, returns string content or null if doesn't exist
//
// Parameters:
//   - path (string): Path to file or directory
//   - callback (function): Callback function (err, data)
//
// Returns null (not error) when file doesn't exist, perfect for existence checks.
// Requires PermissionRead for the specified path.
func (fs *Files) read(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 2 {
		panic(fs.vm.NewTypeError("read requires a file or directory and a callback"))
	}

	dest := call.Arguments[0].String()
	callback, ok := goja.AssertFunction(call.Arguments[1])
	if !ok {
		panic(fs.vm.NewTypeError("second argument must be a callback function"))
	}

  isDir := dirCheck(dest)

	fs.eventLoop.ScheduleTask(&event.Task{
		Callback: func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			mgr := permissions.GetManager()
			canRead := permissions.PermissionRead
			if !mgr.CheckWithPrompt(ctx, canRead, dest) {
				errMsg := mgr.ErrorMessage(canRead, dest)
				callback(goja.Undefined(), fs.vm.ToValue(errMsg), goja.Undefined())
				return
			}

      _, statErr := os.Stat(dest)
      if os.IsNotExist(statErr) {
        // Path doesn't exist - return null data (not an error)
        callback(goja.Undefined(), goja.Null(), goja.Null())
        return
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

      callback(goja.Undefined(), errArg, dataArg)
		},
	})

	return goja.Undefined()
}

// write implements the files.write() method.
//
// Behavior:
//   - 2 args (path with '/'): Create directory recursively
//   - 3 args (path + content): Write file, auto-creating parent directories
//
// Parameters:
//   - path (string): Path to file or directory
//   - content (string, optional): Data to write (omit for directory creation)
//   - callback (function): Callback function (err)
//
// Parent directories are created automatically for file writes.
// Requires PermissionWrite for the specified path.
func (fs *Files) write(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 2 {
		panic(fs.vm.NewTypeError("write requires path ending with / for directories (2 args) or file path with data (3 args)"))
	}

	dest := call.Arguments[0].String()
  isDir := dirCheck(dest)
  data := ""

  var callback goja.Callable
  var ok bool
  if isDir {
    callback, ok = goja.AssertFunction(call.Arguments[1])
    if !ok {
      panic(fs.vm.NewTypeError("if creating a directory, second argument must be a callback"))
    }
  } else {
    if len(call.Arguments) < 3 {
      panic(fs.vm.NewTypeError("write to file requires 3 arguments: path, content, callback"))
    }
    callback, ok = goja.AssertFunction(call.Arguments[2])
    if !ok {
      panic(fs.vm.NewTypeError("if writing to a file, third argument must be a callback"))
    }

    data = call.Arguments[1].String()
  }

	fs.eventLoop.ScheduleTask(&event.Task{
		Callback: func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			mgr := permissions.GetManager()
			canWrite := permissions.PermissionWrite
			if !mgr.CheckWithPrompt(ctx, canWrite, dest) {
				errMsg := mgr.ErrorMessage(canWrite, dest)
				callback(goja.Undefined(), fs.vm.ToValue(errMsg))
				return
			}

      var err error
      if isDir {
        err = os.MkdirAll(dest, 0755) 
      } else {
        // Create parent directories if needed
        if mkdirErr := os.MkdirAll(filepath.Dir(dest), 0755); mkdirErr != nil {
          callback(goja.Undefined(), fs.vm.ToValue(mkdirErr.Error()))
          return
        }
        err = os.WriteFile(dest, []byte(data), 0644)
      }

			var errArg goja.Value
			if err != nil {
				errArg = fs.vm.ToValue(err.Error())
			} else {
				errArg = goja.Null()
			}

			callback(goja.Undefined(), errArg)
		},
	})

	return goja.Undefined()
}

// rm implements the files.rm() method.
//
// Unified removal operation that works on both files and directories.
// Removes directories recursively (including contents).
// Idempotent - succeeds even if path doesn't exist.
//
// Parameters:
//   - path (string): Path to file or directory to remove
//   - callback (function): Callback function (err)
//
// Uses os.RemoveAll() under the hood for recursive removal.
// Requires PermissionWrite for the specified path.
func (fs *Files) rm(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 2 {
		panic(fs.vm.NewTypeError("rm requires a path and a callback"))
	}

	path := call.Arguments[0].String()
	callback, ok := goja.AssertFunction(call.Arguments[1])
	if !ok {
		panic(fs.vm.NewTypeError("second argument must be a callback function"))
	}

	fs.eventLoop.ScheduleTask(&event.Task{
		Callback: func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			mgr := permissions.GetManager()
			canWrite := permissions.PermissionWrite
			if !mgr.CheckWithPrompt(ctx, canWrite, path) {
				errMsg := mgr.ErrorMessage(canWrite, path)
				callback(goja.Undefined(), fs.vm.ToValue(errMsg))
				return
			}

			err := os.RemoveAll(path)

			var errArg goja.Value
			if err != nil {
				errArg = fs.vm.ToValue(err.Error())
			} else {
				errArg = goja.Null()
			}

			callback(goja.Undefined(), errArg)
		},
	})

	return goja.Undefined()
}
