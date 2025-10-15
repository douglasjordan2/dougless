package modules

import (
  "fmt"
	"time"

	"github.com/dop251/goja"
	"github.com/google/uuid"

	"github.com/douglasjordan2/dougless/internal/event"
)

// Timers provides setTimeout/setInterval functionality for JavaScript.
// All timers are scheduled on the event loop for non-blocking execution.
//
// Available globally in JavaScript as setTimeout(), setInterval(), clearTimeout(), and clearInterval().
type Timers struct {
  vm *goja.Runtime   // JavaScript runtime instance
	eventLoop    *event.Loop  // Event loop for async task scheduling
}

// NewTimers creates a new Timers instance with the given event loop.
func NewTimers(eventLoop *event.Loop) *Timers {
  return &Timers{
		eventLoop: eventLoop,
  }
}

// Export creates and returns the timers JavaScript object with all timer methods.
func (t *Timers) Export(vm *goja.Runtime) goja.Value {
  t.vm = vm
  obj := vm.NewObject()

	obj.Set("setTimeout", t.setTimeout)
	obj.Set("setInterval", t.setInterval)
	obj.Set("clearTimeout", t.clearTimeout)
	obj.Set("clearInterval", t.clearInterval)

  return obj
}

// delayHelper is the shared implementation for setTimeout and setInterval.
// It validates the callback, calculates the delay, creates a unique timer ID,
// and schedules the task on the event loop.
//
// Parameters:
//   - callback: Function to execute
//   - delay: Milliseconds to wait before execution (optional, defaults to 0)
//   - isInterval: If true, task repeats at the specified interval
//
// Returns a unique timer ID string for use with clearTimeout/clearInterval.
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

// setTimeout implements setTimeout() - executes a function after a delay.
// Returns a timer ID that can be used with clearTimeout().
//
// JavaScript usage:
//
//	const id = setTimeout(() => console.log('Delayed'), 1000);
func (t *Timers) setTimeout(call goja.FunctionCall) goja.Value {
  return t.delayHelper(call, false)
}

// setInterval implements setInterval() - repeatedly executes a function at intervals.
// Returns a timer ID that can be used with clearInterval().
//
// JavaScript usage:
//
//	const id = setInterval(() => console.log('Tick'), 1000);
func (t *Timers) setInterval(call goja.FunctionCall) goja.Value {
  return t.delayHelper(call, true)
}

// clearHelper is the shared implementation for clearTimeout and clearInterval.
// It cancels a scheduled timer by ID. If the timer doesn't exist, this is a no-op.
func (t *Timers) clearHelper(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) == 0 {
    return goja.Undefined()
  }

  timerID := call.Arguments[0].String()

  t.eventLoop.ClearTimer(timerID)

	return goja.Undefined()
}

// clearTimeout implements clearTimeout() - cancels a timer created with setTimeout().
//
// JavaScript usage:
//
//	const id = setTimeout(() => console.log('Never runs'), 1000);
//	clearTimeout(id);
func (t *Timers) clearTimeout(call goja.FunctionCall) goja.Value {
  return t.clearHelper(call)
}

// clearInterval implements clearInterval() - stops a repeating timer created with setInterval().
//
// JavaScript usage:
//
//	const id = setInterval(() => console.log('Tick'), 1000);
//	clearInterval(id);  // Stops the interval
func (t *Timers) clearInterval(call goja.FunctionCall) goja.Value {
  return t.clearHelper(call)
}
