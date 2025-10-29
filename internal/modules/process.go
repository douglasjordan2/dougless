package modules

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/dop251/goja"
)

// Process provides process-level operations and information.
// This module gives JavaScript access to environment variables, command-line
// arguments, working directory, and process lifecycle management.
type Process struct {
	vm     *goja.Runtime
	argv   []string
	onExit []func(int)
}

// NewProcess creates a new Process module with the given command-line arguments.
// The argv parameter should include the runtime executable and script path.
func NewProcess(argv []string) *Process {
	return &Process{
		argv:   argv,
		onExit: make([]func(int), 0),
	}
}

// Export creates the global process object available in JavaScript.
// The process object provides Node.js-like process APIs including:
//   - process.env: Environment variables
//   - process.argv: Command-line arguments
//   - process.cwd(): Current working directory
//   - process.exit(): Exit the process
//   - process.on(): Signal handling
func (p *Process) Export(vm *goja.Runtime) goja.Value {
	p.vm = vm
	return vm.ToValue(p.createProcessAPI())
}

// createProcessAPI constructs the process object with all its properties and methods.
func (p *Process) createProcessAPI() map[string]interface{} {
	return map[string]interface{}{
		"env":      p.getEnv(),
		"argv":     p.argv,
		"exit":     p.exit,
		"cwd":      p.cwd,
		"chdir":    p.chdir,
		"pid":      os.Getpid(),
		"platform": p.getPlatform(),
		"arch":     p.getArch(),
		"version":  "v0.8.0", // Dougless runtime version
		"on":       p.on,
	}
}

// getEnv returns all environment variables as a JavaScript object.
// Keys are variable names, values are their string values.
func (p *Process) getEnv() map[string]string {
	envMap := make(map[string]string)
	for _, e := range os.Environ() {
		// Parse KEY=VALUE format
		for i := 0; i < len(e); i++ {
			if e[i] == '=' {
				key := e[:i]
				value := e[i+1:]
				envMap[key] = value
				break
			}
		}
	}
	return envMap
}

// exit terminates the process with the specified exit code.
// Calls any registered exit handlers before exiting.
func (p *Process) exit(call goja.FunctionCall) goja.Value {
	code := 0
	if len(call.Arguments) > 0 {
		code = int(call.Argument(0).ToInteger())
	}

	// Call exit handlers
	for _, handler := range p.onExit {
		handler(code)
	}

	os.Exit(code)
	return goja.Undefined()
}

// cwd returns the current working directory as a string.
func (p *Process) cwd(call goja.FunctionCall) goja.Value {
	dir, err := os.Getwd()
	if err != nil {
		panic(p.vm.NewGoError(fmt.Errorf("failed to get current directory: %w", err)))
	}
	return p.vm.ToValue(dir)
}

// chdir changes the current working directory to the specified path.
func (p *Process) chdir(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 1 {
		panic(p.vm.NewTypeError("chdir requires a directory path argument"))
	}

	dir := call.Argument(0).String()
	err := os.Chdir(dir)
	if err != nil {
		panic(p.vm.NewGoError(fmt.Errorf("failed to change directory: %w", err)))
	}

	return goja.Undefined()
}

// getPlatform returns the operating system platform (linux, darwin, windows, etc.)
func (p *Process) getPlatform() string {
	// runtime.GOOS returns the target OS: linux, darwin, windows, etc.
	return runtime.GOOS
}

// getArch returns the CPU architecture (amd64, arm64, etc.)
func (p *Process) getArch() string {
	// runtime.GOARCH returns the target architecture
	return runtime.GOARCH
}

// on registers event handlers for process events.
// Supports:
//   - 'exit': Called before process exits
//   - 'SIGINT': Ctrl+C signal
//   - 'SIGTERM': Termination signal
//   - 'SIGHUP': Hangup signal
func (p *Process) on(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 2 {
		panic(p.vm.NewTypeError("on requires an event name and callback function"))
	}

	event := call.Argument(0).String()
	callback, ok := goja.AssertFunction(call.Argument(1))
	if !ok {
		panic(p.vm.NewTypeError("second argument must be a function"))
	}

	switch event {
	case "exit":
		// Register exit handler
		p.onExit = append(p.onExit, func(code int) {
			callback(goja.Undefined(), p.vm.ToValue(code))
		})

	case "SIGINT":
		// Handle Ctrl+C
		p.setupSignalHandler(syscall.SIGINT, callback)

	case "SIGTERM":
		// Handle termination signal
		p.setupSignalHandler(syscall.SIGTERM, callback)

	case "SIGHUP":
		// Handle hangup signal
		p.setupSignalHandler(syscall.SIGHUP, callback)

	default:
		panic(p.vm.NewTypeError(fmt.Sprintf("unsupported event: %s", event)))
	}

	return goja.Undefined()
}

// setupSignalHandler creates a goroutine to listen for OS signals and invoke the callback.
func (p *Process) setupSignalHandler(sig os.Signal, callback goja.Callable) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, sig)

	go func() {
		for range sigChan {
			// Call the JavaScript callback when signal is received
			callback(goja.Undefined(), p.vm.ToValue(sig.String()))
		}
	}()
}