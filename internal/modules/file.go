package modules

import (
	"context"
	"os"
	"time"

	"github.com/dop251/goja"

	"github.com/douglasjordan2/dougless/internal/event"
	"github.com/douglasjordan2/dougless/internal/permissions"
)

// FileSystem provides asynchronous file system operations for JavaScript.
// All operations are scheduled on the event loop for non-blocking execution.
// Requires appropriate permissions (read/write) based on the operation.
//
// Available globally in JavaScript as the 'file' object (unique to Dougless).
//
// Example usage:
//
//	file.read('data.txt', (err, content) => console.log(content));
//	file.write('output.txt', 'Hello', (err) => console.log('Done'));
type FileSystem struct {
	vm        *goja.Runtime // JavaScript runtime instance
	eventLoop *event.Loop   // Event loop for async task scheduling
}

// NewFileSystem creates a new FileSystem instance with the given event loop.
func NewFileSystem(eventLoop *event.Loop) *FileSystem {
	return &FileSystem{
		eventLoop: eventLoop,
	}
}

// Export creates and returns the file system JavaScript object with all file methods.
func (fs *FileSystem) Export(vm *goja.Runtime) goja.Value {
	fs.vm = vm
	obj := vm.NewObject()

	obj.Set("read", fs.read)
	obj.Set("write", fs.write)
	obj.Set("readdir", fs.readdir)
	obj.Set("exists", fs.exists)
	obj.Set("mkdir", fs.mkdir)
	obj.Set("rmdir", fs.rmdir)
	obj.Set("unlink", fs.unlink)
	obj.Set("stat", fs.stat)

	return obj
}

// read reads the entire contents of a file asynchronously.
// The operation is scheduled on the event loop and requires read permission.
//
// Parameters:
//   - filename (string): The path to the file to read
//   - callback (function): Called with (thisArg, error, data) after completion
//
// If permission is denied or an error occurs, the callback receives an error message
// and the data is undefined. On success, error is null and data contains the file content as a string.
//
// Example:
//
//	file.read('config.json', function(thisArg, err, data) {
//	  if (err) {
//	    console.error('Failed to read file:', err);
//	  } else {
//	    console.log('File contents:', data);
//	  }
//	});
func (fs *FileSystem) read(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 2 {
		panic(fs.vm.NewTypeError("read requires a file and a callback"))
	}

	filename := call.Arguments[0].String()
	callback, ok := goja.AssertFunction(call.Arguments[1])
	if !ok {
		panic(fs.vm.NewTypeError("second argument must be a callback function"))
	}

	fs.eventLoop.ScheduleTask(&event.Task{
		Callback: func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			mgr := permissions.GetManager()
			canRead := permissions.PermissionRead
			if !mgr.CheckWithPrompt(ctx, canRead, filename) {
				errMsg := mgr.ErrorMessage(canRead, filename)
				callback(goja.Undefined(), fs.vm.ToValue(errMsg), goja.Undefined())
				return
			}

			// goroutine reads the file
			data, err := os.ReadFile(filename)

			var errArg, dataArg goja.Value
			if err != nil {
				errArg = fs.vm.ToValue(err.Error())
				dataArg = goja.Undefined()
			} else {
				errArg = goja.Null()
				dataArg = fs.vm.ToValue(string(data))
			}

			callback(goja.Undefined(), errArg, dataArg)
		},
	})

	return goja.Undefined()
}

// write writes data to a file asynchronously, creating or overwriting the file.
// The operation is scheduled on the event loop and requires write permission.
//
// Parameters:
//   - filename (string): The path to the file to write
//   - data (string): The content to write to the file
//   - callback (function): Called with (thisArg, error) after completion
//
// The file is created with permissions 0644 (rw-r--r--) if it doesn't exist.
// If permission is denied or an error occurs, the callback receives an error message.
// On success, the error argument is null.
//
// Example:
//
//	file.write('output.txt', 'Hello World', function(thisArg, err) {
//	  if (err) {
//	    console.error('Failed to write file:', err);
//	  } else {
//	    console.log('File written successfully');
//	  }
//	});
func (fs *FileSystem) write(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 3 {
		panic(fs.vm.NewTypeError("write requires filename, data, and callback"))
	}

	filename := call.Arguments[0].String()
	data := call.Arguments[1].String()
	callback, ok := goja.AssertFunction(call.Arguments[2])
	if !ok {
		panic(fs.vm.NewTypeError("third argument must be a callback function"))
	}

	fs.eventLoop.ScheduleTask(&event.Task{
		Callback: func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			mgr := permissions.GetManager()
			canWrite := permissions.PermissionWrite
			if !mgr.CheckWithPrompt(ctx, canWrite, filename) {
				errMsg := mgr.ErrorMessage(canWrite, filename)
				callback(goja.Undefined(), fs.vm.ToValue(errMsg))
				return
			}

			err := os.WriteFile(filename, []byte(data), 0644)

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

// readdir reads the contents of a directory asynchronously.
// The operation is scheduled on the event loop and requires read permission.
//
// Parameters:
//   - path (string): The path to the directory to read
//   - callback (function): Called with (thisArg, error, entries) after completion
//
// If permission is denied or an error occurs, the callback receives an error message
// and entries is undefined. On success, error is null and entries is an array of
// filename strings (not including '.' and '..').
//
// Example:
//
//	file.readdir('.', function(thisArg, err, files) {
//	  if (err) {
//	    console.error('Failed to read directory:', err);
//	  } else {
//	    console.log('Files:', files);
//	  }
//	});
func (fs *FileSystem) readdir(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 2 {
		panic(fs.vm.NewTypeError("readdir requires a path and a callback"))
	}

	dirPath := call.Arguments[0].String()
	callback, ok := goja.AssertFunction(call.Arguments[1])
	if !ok {
		panic(fs.vm.NewTypeError("second argument must be a callback function"))
	}

	fs.eventLoop.ScheduleTask(&event.Task{
		Callback: func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			mgr := permissions.GetManager()
			canRead := permissions.PermissionRead
			if !mgr.CheckWithPrompt(ctx, canRead, dirPath) {
				errMsg := mgr.ErrorMessage(canRead, dirPath)
				callback(goja.Undefined(), fs.vm.ToValue(errMsg), goja.Undefined())
				return
			}

			entries, err := os.ReadDir(dirPath)

			var errArg, dataArg goja.Value
			if err != nil {
				errArg = fs.vm.ToValue(err.Error())
				dataArg = goja.Undefined()
			} else {
				names := make([]string, len(entries))
				for i, entry := range entries {
					names[i] = entry.Name()
				}
				errArg = goja.Null()
				dataArg = fs.vm.ToValue(names)
			}

			callback(goja.Undefined(), errArg, dataArg)
		},
	})

	return goja.Undefined()
}

// exists checks whether a file or directory exists at the specified path.
// The operation is scheduled on the event loop and requires read permission.
//
// Parameters:
//   - path (string): The path to check for existence
//   - callback (function): Called with (thisArg, exists) after completion
//
// If permission is denied, the callback receives false. Otherwise, it receives
// true if the path exists or false if it doesn't.
//
// Example:
//
//	file.exists('config.json', function(thisArg, exists) {
//	  if (exists) {
//	    console.log('File exists');
//	  } else {
//	    console.log('File not found');
//	  }
//	});
func (fs *FileSystem) exists(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 2 {
		panic(fs.vm.NewTypeError("exists requires a path and a callback"))
	}

	path := call.Arguments[0].String()
	callback, ok := goja.AssertFunction(call.Arguments[1])
	if !ok {
		panic(fs.vm.NewTypeError("second argument must be a callback"))
	}

	fs.eventLoop.ScheduleTask(&event.Task{
		Callback: func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			mgr := permissions.GetManager()
			canRead := permissions.PermissionRead
			if !mgr.CheckWithPrompt(ctx, canRead, path) {
				callback(goja.Undefined(), fs.vm.ToValue(false))
				return
			}

			_, err := os.Stat(path)
			exists := err == nil

			callback(goja.Undefined(), fs.vm.ToValue(exists))
		},
	})

	return goja.Undefined()
}

// mkdir creates a new directory at the specified path.
// It schedules the operation asynchronously via the event loop.
// Requires write permission for the target path.
//
// Parameters:
//   - path (string): The path where the directory should be created
//   - callback (function): Called with (thisArg, error) after completion
//
// The directory is created with permissions 0755 (rwxr-xr-x).
// If permission is denied, the callback receives an error message.
// On success, the error argument is null.
//
// Example:
//
//	file.mkdir('mydir', function(thisArg, err) {
//	  if (err) {
//	    console.error('Failed to create directory:', err);
//	  } else {
//	    console.log('Directory created successfully');
//	  }
//	});
func (fs *FileSystem) mkdir(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 2 {
		panic(fs.vm.NewTypeError("mkdir requires a path and a callback"))
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

			err := os.Mkdir(path, 0755)

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

// rmdir removes an empty directory at the specified path.
// It schedules the operation asynchronously via the event loop.
// Requires write permission for the target path.
//
// Parameters:
//   - path (string): The path to the directory to remove
//   - callback (function): Called with (thisArg, error) after completion
//
// The directory must be empty for removal to succeed.
// If permission is denied, the callback receives an error message.
// On success, the error argument is null.
//
// Example:
//
//	file.rmdir('mydir', function(thisArg, err) {
//	  if (err) {
//	    console.error('Failed to remove directory:', err);
//	  } else {
//	    console.log('Directory removed successfully');
//	  }
//	});
func (fs *FileSystem) rmdir(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 2 {
		panic(fs.vm.NewTypeError("rmdir requires a path and a callback"))
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

			err := os.Remove(path)

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

// unlink removes a file at the specified path.
// It schedules the operation asynchronously via the event loop.
// Requires write permission for the target path.
//
// Parameters:
//   - path (string): The path to the file to remove
//   - callback (function): Called with (thisArg, error) after completion
//
// This function removes files, not directories. Use rmdir to remove directories.
// If permission is denied, the callback receives an error message.
// On success, the error argument is null.
//
// Example:
//
//	file.unlink('test.txt', function(thisArg, err) {
//	  if (err) {
//	    console.error('Failed to delete file:', err);
//	  } else {
//	    console.log('File deleted successfully');
//	  }
//	});
func (fs *FileSystem) unlink(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 2 {
		panic(fs.vm.NewTypeError("unlink requires a path and a callback"))
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

			err := os.Remove(path)

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

// stat retrieves metadata about a file or directory at the specified path.
// It schedules the operation asynchronously via the event loop.
// Requires read permission for the target path.
//
// Parameters:
//   - path (string): The path to the file or directory
//   - callback (function): Called with (thisArg, error, statObject) after completion
//
// The statObject contains the following properties:
//   - size (number): File size in bytes
//   - isDirectory (boolean): True if the path is a directory
//   - isFile (boolean): True if the path is a file
//   - modified (number): Last modification time as Unix timestamp
//   - name (string): Base name of the file or directory
//
// If permission is denied or an error occurs, the callback receives an error message
// and the statObject is undefined. On success, the error argument is null.
//
// Example:
//
//	file.stat('test.txt', function(thisArg, err, stat) {
//	  if (err) {
//	    console.error('Failed to stat file:', err);
//	  } else {
//	    console.log('File size:', stat.size, 'bytes');
//	    console.log('Is directory:', stat.isDirectory);
//	    console.log('Modified:', new Date(stat.modified * 1000));
//	  }
//	});
func (fs *FileSystem) stat(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 2 {
		panic(fs.vm.NewTypeError("stat requires a path and a callback"))
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
			canRead := permissions.PermissionRead
			if !mgr.CheckWithPrompt(ctx, canRead, path) {
				errMsg := mgr.ErrorMessage(canRead, path)
				callback(goja.Undefined(), fs.vm.ToValue(errMsg), goja.Undefined())
				return
			}

			info, err := os.Stat(path)

			var errArg, dataArg goja.Value
			if err != nil {
				errArg = fs.vm.ToValue(err.Error())
				dataArg = goja.Undefined()
			} else {
				// create stat object with file information
				statObj := fs.vm.NewObject()
				statObj.Set("size", info.Size())
				statObj.Set("isDirectory", info.IsDir())
				statObj.Set("isFile", !info.IsDir())
				statObj.Set("modified", info.ModTime().Unix())
				statObj.Set("name", info.Name())

				errArg = goja.Null()
				dataArg = statObj
			}

			callback(goja.Undefined(), errArg, dataArg)
		},
	})

	return goja.Undefined()
}
