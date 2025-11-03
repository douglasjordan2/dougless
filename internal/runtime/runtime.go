package runtime

import (
	"fmt"
	"os"
	"sync"

	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"

	"github.com/douglasjordan2/dougless/internal/modules"
	"github.com/douglasjordan2/dougless/internal/permissions"
)

type Runtime struct {
	vm        *goja.Runtime
	modules   *modules.Registry
	config    *permissions.Config
  wg        sync.WaitGroup // track pending i/o
}

func New(argv []string) *Runtime {
	vm := goja.New()
	moduleRegistry := modules.NewRegistry()

	var config *permissions.Config
	var configPath string
	foundPath, err := permissions.FindConfig(".")
	if err == nil {
		configPath = foundPath
		config, err = permissions.LoadConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to load .douglessrc: %v\n", err)
		}
	}

	rt := &Runtime{
		vm:        vm,
		modules:   moduleRegistry,
		config:    config,
	}

	permManager := permissions.GetManager()
	if config != nil {
		permManager.SetConfig(config)
	}
	if configPath != "" {
		permManager.SetConfigPath(configPath)
	}

	rt.initializeGlobals(argv)
	rt.initializeModules()

	return rt
}

func (rt *Runtime) ExecuteFile(filename string) error {
	source, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return rt.Execute(string(source), filename)
}

func (rt *Runtime) Execute(source, filename string) error {
	transpiledCode, err := rt.transpile(source, filename)
	if err != nil {
		return fmt.Errorf("transpilation error: %w", err)
	}

	_, err = rt.vm.RunScript(filename, transpiledCode)
	if err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

  rt.wg.Wait() // wait for pending futures

	return nil
}

func (rt *Runtime) KeepAlive() func() {
  rt.wg.Add(1)
  return func() {
    rt.wg.Done()
  }
}

func (rt *Runtime) transpile(source, filename string) (string, error) {
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

func (rt *Runtime) initializeGlobals(argv []string) {
	console := modules.NewConsole()
	rt.vm.Set("console", console.Export(rt.vm))

	timers := modules.NewTimers()
	timers.SetRuntime(rt)
	timerObj := timers.Export(rt.vm).ToObject(rt.vm)
	rt.vm.Set("setTimeout", timerObj.Get("setTimeout"))
	rt.vm.Set("setInterval", timerObj.Get("setInterval"))
	rt.vm.Set("clearTimeout", timerObj.Get("clearTimeout"))
	rt.vm.Set("clearInterval", timerObj.Get("clearInterval"))

	path := modules.NewPath()
	rt.vm.Set("path", path.Export(rt.vm))

	files := modules.NewFiles()
  files.SetRuntime(rt)
  rt.vm.Set("files", files.Export(rt.vm))

	httpClient := modules.NewHTTP(rt.vm)
  httpClient.SetRuntime(rt)
  rt.vm.Set("http", httpClient.Export(rt.vm))

	modules.SetupPromise(rt.vm, rt)

	cryptoModule := modules.NewCrypto()
	rt.vm.Set("crypto", cryptoModule.Export(rt.vm))

	processModule := modules.NewProcess(argv)
  processModule.SetRuntime(rt)
  rt.vm.Set("process", processModule.Export(rt.vm))

	rt.vm.Set("require", rt.requireFunction)
}

func (rt *Runtime) initializeModules() {
	rt.modules.Register("path", modules.NewPath())
}

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
