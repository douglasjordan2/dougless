package runtime

import (
	"fmt"
	"os"
	"time"
  "sync"

	"github.com/dop251/goja"
	"github.com/google/uuid"

	"github.com/douglasjordan2/dougless/internal/event"
	"github.com/douglasjordan2/dougless/internal/modules"
)

// Runtime represents the JavaScript execution environment
type Runtime struct {
	vm           *goja.Runtime
	eventLoop    *event.Loop
	modules      *modules.Registry
  timers       map[string]time.Time
  timersMu     sync.Mutex
}

// New creates a new runtime instance
func New() *Runtime {
	vm := goja.New()
	eventLoop := event.NewLoop()
	moduleRegistry := modules.NewRegistry()

	rt := &Runtime{
		vm:        vm,
		eventLoop: eventLoop,
		modules:   moduleRegistry,
    timers:    make(map[string]time.Time),
	}

	// Initialize built-in modules and globals
	rt.initializeGlobals()
	rt.initializeModules()

	return rt
}

// ExecuteFile executes a JavaScript file
func (rt *Runtime) ExecuteFile(filename string) error {
	source, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return rt.Execute(string(source), filename)
}

// Execute runs JavaScript code
func (rt *Runtime) Execute(source, filename string) error {
	// Start the event loop in a separate goroutine
	go rt.eventLoop.Run()
	defer rt.eventLoop.Stop()

	// Execute the script
	_, err := rt.vm.RunScript(filename, source)
	if err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	// Wait for all async operations to complete
	rt.eventLoop.Wait()

	return nil
}

// initializeGlobals sets up global objects and functions
func (rt *Runtime) initializeGlobals() {
	// Console object
	console := rt.vm.NewObject()
	console.Set("log", rt.consoleLog)
	console.Set("error", rt.consoleError)
	console.Set("warn", rt.consoleWarn)
	console.Set("time", rt.consoleTime)
	console.Set("timeEnd", rt.consoleTimeEnd)
	console.Set("table", rt.consoleTable)
	rt.vm.Set("console", console)

	// setTimeout and setInterval (basic implementation)
	rt.vm.Set("setTimeout", rt.setTimeout)
	rt.vm.Set("setInterval", rt.setInterval)
	rt.vm.Set("clearTimeout", rt.clearTimeout)
	rt.vm.Set("clearInterval", rt.clearInterval)
}

// initializeModules registers built-in modules
func (rt *Runtime) initializeModules() {
	// Register built-in modules
	rt.modules.Register("fs", modules.NewFileSystem())
	rt.modules.Register("http", modules.NewHTTP())
	rt.modules.Register("path", modules.NewPath())
	
	// Set up require function
	rt.vm.Set("require", rt.requireFunction)
}

// Console functions
func (rt *Runtime) consoleLog(call goja.FunctionCall) goja.Value {
	args := make([]any, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg.Export()
	}
	fmt.Println(args...)
	return goja.Undefined()
}

func (rt *Runtime) consoleError(call goja.FunctionCall) goja.Value {
	args := make([]any, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg.Export()
	}
	fmt.Print("ERROR: ")
	fmt.Println(args...)
	return goja.Undefined()
}

func (rt *Runtime) consoleWarn(call goja.FunctionCall) goja.Value {
	args := make([]any, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg.Export()
	}
	fmt.Print("WARN: ")
	fmt.Println(args...)
	return goja.Undefined()
}

func (rt *Runtime) consoleTime(call goja.FunctionCall) goja.Value {
  label := "default"
  if len(call.Arguments) > 0 {
    label = call.Arguments[0].String()
  }

  rt.timersMu.Lock()
  rt.timers[label] = time.Now()
  rt.timersMu.Unlock()

  return goja.Undefined()
}

func (rt *Runtime) consoleTimeEnd(call goja.FunctionCall) goja.Value {
  label := "default"
  if len(call.Arguments) > 0 {
    label = call.Arguments[0].String()
  }

  rt.timersMu.Lock()
  startTime, exists := rt.timers[label]
  if exists {
    delete(rt.timers, label)
  }
  rt.timersMu.Unlock()

  if !exists {
    fmt.Printf("Warning: No such label '%s' for console.timeEnd()\n", label)
    return goja.Undefined()
  }

  duration := time.Since(startTime)
  fmt.Printf("%s: %.3fms\n", label, float64(duration.Microseconds())/1000.0)

  return goja.Undefined()
}

func (rt *Runtime) consoleTable(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) == 0 {
    return goja.Undefined()
  }

  data := call.Arguments[0].Export()

  // Handle different data types
  switch v := data.(type) {
  case []any:
    rt.printArrayTable(v)
  case map[string]any:
    rt.printObjectTable(v)
  default:
    // Fallback to regular log for unsupported types
    fmt.Println(data)
  }

  return goja.Undefined()
}

func (rt *Runtime) printArrayTable(data []any) {
  if len(data) == 0 {
    return
  }

  // Calculate column width based on content
  maxWidth := 36
  for _, item := range data {
    str := fmt.Sprintf("%v", item)
    if len(str) > maxWidth {
      maxWidth = len(str)
    }
  }
  if maxWidth > 60 {
    maxWidth = 60
  }

  // Print table header
  fmt.Println("┌─────────┬" + repeatChar('─', maxWidth+2) + "┐")
  fmt.Printf("│ (index) │ %-*s │\n", maxWidth, "Values")
  fmt.Println("├─────────┼" + repeatChar('─', maxWidth+2) + "┤")

  // Print table rows
  for i, item := range data {
    valueStr := fmt.Sprintf("%v", item)
    if len(valueStr) > maxWidth {
      valueStr = valueStr[:maxWidth-3] + "..."
    }
    fmt.Printf("│ %-7d │ %-*s │\n", i, maxWidth, valueStr)
  }

  // Print table footer
  fmt.Println("└─────────┴" + repeatChar('─', maxWidth+2) + "┘")
}

func (rt *Runtime) printObjectTable(data map[string]any) {
  if len(data) == 0 {
    return
  }

  // Calculate column widths
  maxKeyWidth := 10
  maxValWidth := 24
  for key, value := range data {
    if len(key) > maxKeyWidth {
      maxKeyWidth = len(key)
    }
    valStr := fmt.Sprintf("%v", value)
    if len(valStr) > maxValWidth {
      maxValWidth = len(valStr)
    }
  }
  if maxKeyWidth > 30 {
    maxKeyWidth = 30
  }
  if maxValWidth > 50 {
    maxValWidth = 50
  }

  // Print table header
  fmt.Println("┌" + repeatChar('─', maxKeyWidth+2) + "┬" + repeatChar('─', maxValWidth+2) + "┐")
  fmt.Printf("│ %-*s │ %-*s │\n", maxKeyWidth, "(index)", maxValWidth, "Values")
  fmt.Println("├" + repeatChar('─', maxKeyWidth+2) + "┼" + repeatChar('─', maxValWidth+2) + "┤")

  // Print table rows
  for key, value := range data {
    keyStr := key
    if len(keyStr) > maxKeyWidth {
      keyStr = keyStr[:maxKeyWidth-3] + "..."
    }
    valueStr := fmt.Sprintf("%v", value)
    if len(valueStr) > maxValWidth {
      valueStr = valueStr[:maxValWidth-3] + "..."
    }
    fmt.Printf("│ %-*s │ %-*s │\n", maxKeyWidth, keyStr, maxValWidth, valueStr)
  }

  // Print table footer
  fmt.Println("└" + repeatChar('─', maxKeyWidth+2) + "┴" + repeatChar('─', maxValWidth+2) + "┘")
}

// Helper function to repeat a character n times
func repeatChar(char rune, count int) string {
  result := make([]rune, count)
  for i := 0; i < count; i++ {
    result[i] = char
  }
  return string(result)
}

func (rt *Runtime) delayHelper(call goja.FunctionCall, isInterval bool) goja.Value {
  callback, ok := goja.AssertFunction(call.Arguments[0])
  if !ok {
    panic(rt.vm.NewTypeError("callback must be a function"))
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

  rt.eventLoop.ScheduleTask(task)

	return rt.vm.ToValue(timerID)
}

// Timer functions (placeholder implementations)
func (rt *Runtime) setTimeout(call goja.FunctionCall) goja.Value {
  return rt.delayHelper(call, false)
}

func (rt *Runtime) setInterval(call goja.FunctionCall) goja.Value {
  return rt.delayHelper(call, true)
}

func (rt *Runtime) clearHelper(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) == 0 {
    return goja.Undefined()
  }

  timerID := call.Arguments[0].String()

  rt.eventLoop.ClearTimer(timerID)

	return goja.Undefined()
}

func (rt *Runtime) clearTimeout(call goja.FunctionCall) goja.Value {
  return rt.clearHelper(call)
}

func (rt *Runtime) clearInterval(call goja.FunctionCall) goja.Value {
  return rt.clearHelper(call)
}

// requireFunction implements the require() function for modules
func (rt *Runtime) requireFunction(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(rt.vm.NewTypeError("require() missing module name"))
	}

	moduleName := call.Arguments[0].String()
	module := rt.modules.Get(moduleName)
	
	if module == nil {
		panic(rt.vm.NewGoError(fmt.Errorf("Cannot find module '%s'", moduleName)))
	}

	return module.Export(rt.vm)
}

func (r *Runtime) Evaluate(code string) (goja.Value, error) {
  return r.vm.RunString(code)
}
