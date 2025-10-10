package modules

import (
  "os"

  "github.com/dop251/goja"
  "github.com/douglasjordan2/dougless/internal/event"
)

type FileSystem struct {
  vm        *goja.Runtime
  eventLoop *event.Loop
}

func NewFileSystem(eventLoop *event.Loop) *FileSystem {
  return &FileSystem{
    eventLoop: eventLoop,
  }
}

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

func (fs *FileSystem) read(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 2 {
    panic(fs.vm.ToValue("read requires a file and a callback"))
  }

  filename := call.Arguments[0].String()
  callback, ok := goja.AssertFunction(call.Arguments[1])
  if !ok {
    panic(fs.vm.ToValue("second argument must be a callback function"))
  }

  fs.eventLoop.ScheduleTask(&event.Task{
    Callback: func() {
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

// write file async
func (fs *FileSystem) write(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 3 {
    panic(fs.vm.ToValue("write requires filename, data, and callback"))
  }

  filename := call.Arguments[0].String()
  data := call.Arguments[1].String()
  callback, ok := goja.AssertFunction(call.Arguments[2])
  if !ok {
    panic(fs.vm.ToValue("third argument must be a callback function"))
  }

  fs.eventLoop.ScheduleTask(&event.Task{
    Callback: func() {
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

// read directory contents async
func (fs *FileSystem) readdir(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 2 {
    panic(fs.vm.ToValue("readdir requires a path and a callback"))
  }

  dirPath := call.Arguments[0].String()
  callback, ok := goja.AssertFunction(call.Arguments[1])
  if !ok {
    panic(fs.vm.ToValue("second argument must be a callback function"))
  }

  fs.eventLoop.ScheduleTask(&event.Task{
    Callback: func() {
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

func (fs *FileSystem) exists(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 2 {
    panic(fs.vm.ToValue("exists requires a path and a callback"))
  }

  path := call.Arguments[0].String()
  callback, ok := goja.AssertFunction(call.Arguments[1])
  if !ok {
    panic(fs.vm.ToValue("second argument must be a callback"))
  }

  fs.eventLoop.ScheduleTask(&event.Task{
    Callback: func() {
      _, err := os.Stat(path)
      exists := err == nil

      callback(goja.Undefined(), fs.vm.ToValue(exists))
    },
  })

  return goja.Undefined()
}

func (fs *FileSystem) mkdir(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 2 {
    panic(fs.vm.ToValue("mkdir requires a path and a callback"))
  }

  path := call.Arguments[0].String()
  callback, ok := goja.AssertFunction(call.Arguments[1])
  if !ok {
    panic(fs.vm.ToValue("second argument must be a callback function"))
  }

  fs.eventLoop.ScheduleTask(&event.Task{
    Callback: func() {
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

func (fs *FileSystem) rmdir(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 2 {
    panic(fs.vm.ToValue("rmdir requires a path and a callback"))
  }

  path := call.Arguments[0].String()
  callback, ok := goja.AssertFunction(call.Arguments[1])
  if !ok {
    panic(fs.vm.ToValue("second argument must be a callback function"))
  }

  fs.eventLoop.ScheduleTask(&event.Task{
    Callback: func() {
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

func (fs *FileSystem) unlink(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 2 {
    panic(fs.vm.ToValue("unlink requires a path and a callback"))
  }

  path := call.Arguments[0].String()
  callback, ok := goja.AssertFunction(call.Arguments[1])
  if !ok {
    panic(fs.vm.ToValue("second argument must be a callback function"))
  }

  fs.eventLoop.ScheduleTask(&event.Task{
    Callback: func() {
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

func (fs *FileSystem) stat(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 2 {
    panic(fs.vm.ToValue("stat requires a path and a callback"))
  }

  path := call.Arguments[0].String()
  callback, ok := goja.AssertFunction(call.Arguments[1])
  if !ok {
    panic(fs.vm.ToValue("second argument must be a callback function"))
  }

  fs.eventLoop.ScheduleTask(&event.Task{
    Callback: func() {
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
