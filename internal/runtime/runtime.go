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
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"

	"github.com/douglasjordan2/dougless/internal/event"
	"github.com/douglasjordan2/dougless/internal/modules"
)

// Runtime represents the JavaScript execution environment.
// It coordinates the Goja VM, event loop, and module system to provide
// a complete JavaScript runtime with async capabilities.
type Runtime struct {
	vm        *goja.Runtime        // Goja JavaScript VM (ES5.1)
	eventLoop *event.Loop          // Event loop for async operations
	modules   *modules.Registry    // Registry of loadable modules
	timers    map[string]time.Time // Timer tracking for console.time()
	timersMu  sync.Mutex           // Protects timers map
}

// New creates and initializes a new Runtime instance.
// It sets up the Goja VM, event loop, module registry, and all global APIs.
//
// The returned runtime is ready to execute JavaScript code via Execute or ExecuteFile.
//
// Example:
//
//	rt := runtime.New()
//	if err := rt.ExecuteFile("script.js"); err != nil {
//	    log.Fatal(err)
//	}
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

	rt.initializeGlobals()
	rt.initializeModules()

	return rt
}

// ExecuteFile reads and executes a JavaScript file.
// The source code is automatically transpiled from ES6+ to ES5 before execution.
//
// Returns an error if the file cannot be read, transpilation fails, or
// execution encounters a JavaScript error.
func (rt *Runtime) ExecuteFile(filename string) error {
	source, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return rt.Execute(string(source), filename)
}

// Execute runs JavaScript code from a string.
//
// The source code is:
//  1. Transpiled from ES6+ to ES5 using esbuild
//  2. Executed in the Goja VM
//  3. All async operations are run through the event loop
//
// The filename parameter is used for error messages and source maps.
//
// Returns an error if transpilation or execution fails.
func (rt *Runtime) Execute(source, filename string) error {
	// Start the event loop in a separate goroutine
	go rt.eventLoop.Run()
	defer rt.eventLoop.Stop()

	transpiledCode, err := rt.transpile(source, filename)
	if err != nil {
		return fmt.Errorf("transpilation error: %w", err)
	}

	_, err = rt.vm.RunScript(filename, transpiledCode)
	if err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	rt.eventLoop.Wait()

	return nil
}

// transpile converts ES6+ JavaScript to ES5 using esbuild.
//
// Features:
//   - Transpiles to ES2017 (supports async/await natively)
//   - Generates inline source maps for accurate error reporting
//   - Handles edge case of empty scripts (which produce invalid source maps)
//   - Reports warnings to stderr
//
// Returns the transpiled JavaScript code or an error if transpilation fails.
func (rt *Runtime) transpile(source, filename string) (string, error) {
	// Use inline source maps for better debugging, but only when there's actual code
	// Empty/whitespace-only scripts produce invalid source maps that Goja can't parse
	sourcemap := api.SourceMapInline
	if len(source) == 0 {
		sourcemap = api.SourceMapNone
	}

	result := api.Transform(source, api.TransformOptions{
		Loader:     api.LoaderJS,
		Target:     api.ES2017,
		Sourcefile: filename,
		Format:     api.FormatDefault,
		Sourcemap:  sourcemap,
	})

	if len(result.Errors) > 0 {
		// Return the first error with details
		err := result.Errors[0]
		return "", fmt.Errorf("%s:%d:%d: %s",
			err.Location.File,
			err.Location.Line,
			err.Location.Column,
			err.Text,
		)
	}

	if len(result.Warnings) > 0 {
		for _, warning := range result.Warnings {
			fmt.Fprintf(os.Stderr, "Warning: %s:%d:%d: %s\n",
				warning.Location.File,
				warning.Location.Line,
				warning.Location.Column,
				warning.Text,
			)
		}
	}

	return string(result.Code), nil
}

// initializeGlobals sets up all global objects and functions available in JavaScript.
//
// Global APIs (no require needed):
//   - console (log, error, warn, time, timeEnd, table)
//   - setTimeout, setInterval, clearTimeout, clearInterval
//   - path (join, resolve, dirname, basename, extname, sep)
//   - file (read, write, readdir, exists, mkdir, rmdir, unlink, stat)
//   - http (get, post, createServer)
//   - Promise (constructor, resolve, reject, all, race)
//   - require() (for CommonJS-style module loading)
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

	// Path
	path := modules.NewPath()
	rt.vm.Set("path", path.Export(rt.vm))

	// File system
	fileSystem := modules.NewFileSystem(rt.eventLoop)
	rt.vm.Set("file", fileSystem.Export(rt.vm))

	// HTTP
	httpClient := modules.NewHTTP(rt.eventLoop)
	rt.vm.Set("http", httpClient.Export(rt.vm))

	// Promise
	modules.SetupPromise(rt.vm, rt.eventLoop)

	// require() function for module loading
	rt.vm.Set("require", rt.requireFunction)
}

// initializeModules registers built-in modules that can be loaded via require().
//
// Currently registered modules:
//   - path: File path manipulation utilities
//
// Note: Most core APIs (file, http, console) are available globally and
// don't need to be required.
func (rt *Runtime) initializeModules() {
	rt.modules.Register("path", modules.NewPath())
}

// requireFunction implements the CommonJS require() function.
// It loads and returns built-in modules registered in the module registry.
//
// Usage in JavaScript:
//
//	const path = require('path');
//	path.join('foo', 'bar');  // 'foo/bar'
//
// Panics with a TypeError if no module name is provided, or a GoError if
// the module doesn't exist.
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

// Evaluate executes JavaScript code and returns the result.
// This is used primarily by the REPL for interactive evaluation.
//
// Unlike Execute(), this method:
//   - Does NOT transpile the code (assumes ES5)
//   - Does NOT start/stop the event loop
//   - Returns the evaluation result directly
//
// Returns the result value and any error that occurred during execution.
func (r *Runtime) Evaluate(code string) (goja.Value, error) {
	return r.vm.RunString(code)
}
