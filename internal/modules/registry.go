package modules

import (
	"github.com/dop251/goja"
)

// Module represents a built-in module
type Module interface {
	Export(vm *goja.Runtime) goja.Value
}

// Registry manages built-in modules
type Registry struct {
	modules map[string]Module
}

// NewRegistry creates a new module registry
func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
}

// Register registers a module with the given name
func (r *Registry) Register(name string, module Module) {
	r.modules[name] = module
}

// Get retrieves a module by name
func (r *Registry) Get(name string) Module {
	return r.modules[name]
}

// List returns all registered module names
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.modules))
	for name := range r.modules {
		names = append(names, name)
	}
	return names
}
