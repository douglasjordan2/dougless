package modules

import (
	"fmt"
  "os"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/google/uuid"
)

type RuntimeKeepAlive interface {
	KeepAlive() func()
}

type Timers struct {
	vm      *goja.Runtime
  timers  map[string]chan struct{}
  mu      sync.Mutex
  runtime RuntimeKeepAlive
}

func NewTimers() *Timers {
	return &Timers{
		timers: make(map[string]chan struct{}),
	}
}

func (t *Timers) SetRuntime(rt RuntimeKeepAlive) {
	t.runtime = rt
}

func (t *Timers) Export(vm *goja.Runtime) goja.Value {
	t.vm = vm
	obj := vm.NewObject()

	obj.Set("setTimeout", t.setTimeout)
	obj.Set("setInterval", t.setInterval)
	obj.Set("clearTimeout", t.clearTimeout)
	obj.Set("clearInterval", t.clearInterval)

	return obj
}

func timerHelper(t *Timers, call goja.FunctionCall) (fn goja.Callable, ms int64, timerID string, done func(), cancel chan struct{}) {
  if len(call.Arguments) < 2 {
		panic(t.vm.NewTypeError("timer setters require at least 2 arguments"))
	}

  fn, ok := goja.AssertFunction(call.Arguments[0])
  if !ok {
    panic(t.vm.NewTypeError("First argument must be a function"))
  }

  ms = call.Arguments[1].ToInteger()

  timerID = uuid.New().String()
  cancel = make(chan struct{})

  t.mu.Lock()
  t.timers[timerID] = cancel
  t.mu.Unlock()

  done = t.runtime.KeepAlive()

  return fn, ms, timerID, done, cancel
}

func (t *Timers) setTimeout(call goja.FunctionCall) goja.Value {
  fn, ms, timerID, done, cancel := timerHelper(t, call)

  go func() {
    defer done()

    select {
    case <-time.After(time.Duration(ms) * time.Millisecond):
      // execute callback in vm
      if _, err := fn(nil, call.Arguments[2:]...); err != nil {
        fmt.Fprintf(os.Stderr, "setTimeout callback error: %v\n", err)
      }
      
      // cleanup
      t.mu.Lock()
      delete(t.timers, timerID)
      t.mu.Unlock()

    case <-cancel:
      t.mu.Lock()
      delete(t.timers, timerID)
      t.mu.Unlock()
      return
    }
  }()

  return t.vm.ToValue(timerID)
}

func (t *Timers) setInterval(call goja.FunctionCall) goja.Value {
  fn, ms, timerID, done, cancel := timerHelper(t, call)

  go func() {
    defer done()
    ticker := time.NewTicker(time.Duration(ms) * time.Millisecond)
    defer ticker.Stop()

    for {
      select {
      case <-ticker.C:
        if _, err := fn(nil, call.Arguments[2:]...); err != nil {
          fmt.Fprintf(os.Stderr, "setInterval callback error: %v\n", err)
        }
      case <-cancel:
        t.mu.Lock()
        delete(t.timers, timerID)
        t.mu.Unlock()
        return
      }
    }
  }()

  return t.vm.ToValue(timerID)
}

func (t *Timers) clearTimeout(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 1 {
    return goja.Undefined()
  }

  timerID := call.Arguments[0].String()

  t.mu.Lock()
  cancel, ok := t.timers[timerID]
  if ok {
    close(cancel)
    delete(t.timers, timerID)
  }
  t.mu.Unlock()

  return goja.Undefined()
}

func (t *Timers) clearInterval(call goja.FunctionCall) goja.Value {
  return t.clearTimeout(call)
}
