package modules

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/dop251/goja"
)

type Process struct {
	vm      *goja.Runtime
  runtime RuntimeKeepAlive
	argv    []string
	onExit  []func(int)
}

func NewProcess(argv []string) *Process {
	return &Process{
		argv:   argv,
		onExit: make([]func(int), 0),
	}
}

func (p *Process) SetRuntime(rt RuntimeKeepAlive) {
  p.runtime = rt
}

func (p *Process) Export(vm *goja.Runtime) goja.Value {
	p.vm = vm
	return vm.ToValue(p.createProcessAPI())
}

func (p *Process) createProcessAPI() map[string]any {
	return map[string]any{
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

func (p *Process) getEnv() map[string]string {
	envMap := make(map[string]string)
	for _, e := range os.Environ() {
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

func (p *Process) exit(call goja.FunctionCall) goja.Value {
	code := 0
	if len(call.Arguments) > 0 {
		code = int(call.Argument(0).ToInteger())
	}

	for _, handler := range p.onExit {
		handler(code)
	}

	os.Exit(code)
	return goja.Undefined()
}

func (p *Process) cwd(call goja.FunctionCall) goja.Value {
	dir, err := os.Getwd()
	if err != nil {
		panic(p.vm.NewGoError(fmt.Errorf("failed to get current directory: %w", err)))
	}
	return p.vm.ToValue(dir)
}

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

func (p *Process) getPlatform() string {
	return runtime.GOOS
}

func (p *Process) getArch() string {
	return runtime.GOARCH
}

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
		p.onExit = append(p.onExit, func(code int) {
			callback(goja.Undefined(), p.vm.ToValue(code))
		})

	case "SIGINT":
		p.setupSignalHandler(syscall.SIGINT, callback)

	case "SIGTERM":
		p.setupSignalHandler(syscall.SIGTERM, callback)

	case "SIGHUP":
		p.setupSignalHandler(syscall.SIGHUP, callback)

	default:
		panic(p.vm.NewTypeError(fmt.Sprintf("unsupported event: %s", event)))
	}

	return goja.Undefined()
}

func (p *Process) setupSignalHandler(sig os.Signal, callback goja.Callable) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, sig)

	// Signal handlers should NOT use KeepAlive - they're passive listeners
	// that should not prevent the runtime from exiting when all real work is done
	go func() {
		for range sigChan {
			callback(goja.Undefined(), p.vm.ToValue(sig.String()))
		}
	}()
}
