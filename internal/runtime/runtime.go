package runtime

import (
	"fmt"
	"io/ioutil"

	"github.com/dop251/goja"
	"github.com/douglasjordan/dougless-runtime/internal/event"
	"github.com/douglasjordan/dougless-runtime/internal/modules"
)

// Runtime represents the JavaScript execution environment
type Runtime struct {
	vm        *goja.Runtime
	eventLoop *event.Loop
	modules   *modules.Registry
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
	}

	// Initialize built-in modules and globals
	rt.initializeGlobals()
	rt.initializeModules()

	return rt
}

// ExecuteFile executes a JavaScript file
func (rt *Runtime) ExecuteFile(filename string) error {
	source, err := ioutil.ReadFile(filename)
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
	args := make([]interface{}, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg.Export()
	}
	fmt.Println(args...)
	return goja.Undefined()
}

func (rt *Runtime) consoleError(call goja.FunctionCall) goja.Value {
	args := make([]interface{}, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg.Export()
	}
	fmt.Println("ERROR:", args...)
	return goja.Undefined()
}

func (rt *Runtime) consoleWarn(call goja.FunctionCall) goja.Value {
	args := make([]interface{}, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg.Export()
	}
	fmt.Println("WARN:", args...)
	return goja.Undefined()
}

// Timer functions (placeholder implementations)
func (rt *Runtime) setTimeout(call goja.FunctionCall) goja.Value {
	// TODO: Implement proper setTimeout with event loop
	return rt.vm.ToValue(1)
}

func (rt *Runtime) setInterval(call goja.FunctionCall) goja.Value {
	// TODO: Implement proper setInterval with event loop
	return rt.vm.ToValue(1)
}

func (rt *Runtime) clearTimeout(call goja.FunctionCall) goja.Value {
	// TODO: Implement clearTimeout
	return goja.Undefined()
}

func (rt *Runtime) clearInterval(call goja.FunctionCall) goja.Value {
	// TODO: Implement clearInterval
	return goja.Undefined()
}

// requireFunction implements the require() function for modules
func (rt *Runtime) requireFunction(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(rt.vm.NewTypeError("require() missing module name"))
	}

	moduleName := call.Arguments[0].String()
	module := rt.modules.Get(moduleName)
	
	if module == nil {
		panic(rt.vm.NewError(fmt.Sprintf("Cannot find module '%s'", moduleName)))
	}

	return module.Export(rt.vm)
}
