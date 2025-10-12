package modules

import (
  "fmt"
	"time"

	"github.com/dop251/goja"
	"github.com/google/uuid"

	"github.com/douglasjordan2/dougless/internal/event"
)

type Timers struct {
  vm *goja.Runtime
	eventLoop    *event.Loop
}

func NewTimers(eventLoop *event.Loop) *Timers {
  return &Timers{
		eventLoop: eventLoop,
  }
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

func (t *Timers) delayHelper(call goja.FunctionCall, isInterval bool) goja.Value {
  if len(call.Arguments) == 0 {
    panic(t.vm.NewTypeError("setTimeout/setInterval requires a callback function"))
  }

  callback, ok := goja.AssertFunction(call.Arguments[0])
  if !ok {
    panic(t.vm.NewTypeError("callback must be a function"))
  }

  cb := func() {
    _, err := callback(goja.Undefined())
    if err != nil {
      fmt.Printf("Timer callback error: %v\n", err)
    }
  }

  var delayMs int64 = 0
  if len(call.Arguments) > 1 {
    delayMs = call.Arguments[1].ToInteger()
  }
  delay := time.Duration(delayMs) * time.Millisecond

  timerID := uuid.New().String()

  task := &event.Task{
    ID: timerID,
    Callback: cb,
    Delay: delay,
    Interval: isInterval,
  }

  t.eventLoop.ScheduleTask(task)

	return t.vm.ToValue(timerID)
}

// Timer functions (placeholder implementations)
func (t *Timers) setTimeout(call goja.FunctionCall) goja.Value {
  return t.delayHelper(call, false)
}

func (t *Timers) setInterval(call goja.FunctionCall) goja.Value {
  return t.delayHelper(call, true)
}

func (t *Timers) clearHelper(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) == 0 {
    return goja.Undefined()
  }

  timerID := call.Arguments[0].String()

  t.eventLoop.ClearTimer(timerID)

	return goja.Undefined()
}

func (t *Timers) clearTimeout(call goja.FunctionCall) goja.Value {
  return t.clearHelper(call)
}

func (t *Timers) clearInterval(call goja.FunctionCall) goja.Value {
  return t.clearHelper(call)
}
