package modules

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/dop251/goja"

	"github.com/douglasjordan2/dougless/internal/permissions"
)

type Files struct {
	vm      *goja.Runtime
  runtime RuntimeKeepAlive
}

func NewFiles() *Files {
	return &Files{}
}

func (fs *Files) SetRuntime(rt RuntimeKeepAlive) {
  fs.runtime = rt
}

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
			runtime:     fs.runtime,
			state:       PromisePending,
			onFulfilled: []goja.Callable{},
			onRejected:  []goja.Callable{},
		}

    done := fs.runtime.KeepAlive()
    go func() {
      defer done()

      ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
      defer cancel()

      errArg, dataArg := fs.doRead(ctx, dest)

      if errArg != goja.Null() && !goja.IsNull(errArg) {
        promise.reject(errArg)
      } else {
        promise.resolve(dataArg)
      }
		}()

		return CreatePromiseObject(fs.vm, promise)
	}

  done := fs.runtime.KeepAlive()
  go func() {
    defer done()
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    errArg, dataArg := fs.doRead(ctx, dest)
    callback(goja.Undefined(), errArg, dataArg)
	}()

	return goja.Undefined()
}

func (fs *Files) doWrite(ctx context.Context, dest string, data string) goja.Value {
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

func (fs *Files) write(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 1 {
		panic(fs.vm.NewTypeError("write requires at least a path"))
	}

	dest := call.Arguments[0].String()
	isDir := dirCheck(dest)
	data := "" // Default to empty string for files

	var callback goja.Callable
	var ok bool

	if isDir {
		if len(call.Arguments) > 1 {
			callback, ok = goja.AssertFunction(call.Arguments[1])
		}
	} else {
		if len(call.Arguments) > 1 {
			if !goja.IsUndefined(call.Arguments[1]) && !goja.IsNull(call.Arguments[1]) {
				callback, ok = goja.AssertFunction(call.Arguments[1])
				if !ok {
					data = call.Arguments[1].String()
					if len(call.Arguments) > 2 {
						callback, ok = goja.AssertFunction(call.Arguments[2])
					}
				}
			}
		}
	}

	if !ok {
		promise := &Promise{
			vm:          fs.vm,
			runtime:     fs.runtime,
			state:       PromisePending,
			onFulfilled: []goja.Callable{},
			onRejected:  []goja.Callable{},
		}

    done := fs.runtime.KeepAlive()
    go func() {
      defer done()
      ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
      defer cancel()

      errArg := fs.doWrite(ctx, dest, data)

      if errArg != goja.Null() && !goja.IsNull(errArg) {
        promise.reject(errArg)
      } else {
        promise.resolve(goja.Null())
      }
		}()

		return CreatePromiseObject(fs.vm, promise)
	}

  done := fs.runtime.KeepAlive()
  go func() {
    defer done()
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    errArg := fs.doWrite(ctx, dest, data)
    callback(goja.Undefined(), errArg)
	}()

	return goja.Undefined()
}

func (fs *Files) doRm(ctx context.Context, path string) goja.Value {
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
			runtime:     fs.runtime,
			state:       PromisePending,
			onFulfilled: []goja.Callable{},
			onRejected:  []goja.Callable{},
		}

    done := fs.runtime.KeepAlive()
    go func() {
      defer done()
      ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
      defer cancel()

      errArg := fs.doRm(ctx, path)

      if errArg != goja.Null() && !goja.IsNull(errArg) {
        promise.reject(errArg)
      } else {
        promise.resolve(goja.Null())
      }
		}()

		return CreatePromiseObject(fs.vm, promise)
	}

  done := fs.runtime.KeepAlive()
  go func() {
    defer done()
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    errArg := fs.doRm(ctx, path)

    callback(goja.Undefined(), errArg)
	}()

	return goja.Undefined()
}
