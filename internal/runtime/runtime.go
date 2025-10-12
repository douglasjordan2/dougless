// Package runtime implements the core JavaScript execution environment for Dougless.
//
// The runtime orchestrates the Goja VM, event loop, and module system to provide
// a Node.js-style JavaScript runtime with support for async operations, timers,
// file I/O, HTTP, and WebSockets.
//
// Key components:
//   - JavaScript execution using the Goja engine (ES5.1 with planned ES6+ transpilation)
//   - Non-blocking event loop for async operations
//   - CommonJS-style module system with require()
//   - Global APIs (console, setTimeout, file, http) available without imports
//   - REPL support for interactive development
//
// Example usage:
//
//	rt := runtime.New()
//	err := rt.ExecuteFile("script.js")
package runtime

import (
	"fmt"
	"os"
	"time"
  "sync"

	"github.com/dop251/goja"

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
	console := modules.NewConsole()
	rt.vm.Set("console", console.Export(rt.vm))

	// Timers
  timers := modules.NewTimers(rt.eventLoop)
  timerObj := timers.Export(rt.vm).ToObject(rt.vm)
  rt.vm.Set("setTimeout", timerObj.Get("setTimeout"))
  rt.vm.Set("setInterval", timerObj.Get("setInterval"))
  rt.vm.Set("clearTimeout", timerObj.Get("clearTimeout"))
  rt.vm.Set("clearInterval", timerObj.Get("clearInterval"))

	// File system
	fileSystem := modules.NewFileSystem(rt.eventLoop)
	rt.vm.Set("file", fileSystem.Export(rt.vm))

  // HTTP
  httpClient := modules.NewHTTP(rt.eventLoop)
  rt.vm.Set("http", httpClient.Export(rt.vm))
}

// initializeModules registers built-in modules
func (rt *Runtime) initializeModules() {
	// Register built-in modules (for require() support)
	rt.modules.Register("path", modules.NewPath())
	// TODO: Add more modules here (http, crypto, etc.)
	
	// Set up require function
	rt.vm.Set("require", rt.requireFunction)
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
