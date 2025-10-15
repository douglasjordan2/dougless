// Package modules provides the module system and built-in modules for Dougless.
//
// This package implements:
//   - Module registry for CommonJS-style require()
//   - Global APIs (console, timers, file, http, path, Promise)
//   - Module interface for extensibility
//
// Built-in modules can be accessed either:
//  1. Globally (e.g., console.log, file.read, http.get)
//  2. Via require() (e.g., require('path'))
//
// Example:
//
//	// Create and register a custom module
//	registry := modules.NewRegistry()
//	registry.Register("mymodule", customModule)
//
//	// In JavaScript:
//	const mymodule = require('mymodule');
package modules

import (
	"github.com/dop251/goja"
)

// Module represents a built-in module that can be loaded via require().
// All modules must implement the Export method, which returns a Goja value
// (typically an object) containing the module's API.
type Module interface {
	// Export creates and returns the module's JavaScript API object.
	// The vm parameter is the Goja runtime instance.
	Export(vm *goja.Runtime) goja.Value
}

// Registry manages the collection of built-in modules available via require().
// It provides a simple name-to-module mapping for module resolution.
type Registry struct {
	modules map[string]Module  // Map of module names to Module implementations
}

// NewRegistry creates and initializes a new module registry.
// The registry starts empty; modules must be registered via Register().
func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
}

// Register adds a module to the registry with the specified name.
// This makes the module available via require(name) in JavaScript.
//
// Example:
//
//	registry.Register("path", modules.NewPath())
//	// Now in JS: const path = require('path');
func (r *Registry) Register(name string, module Module) {
	r.modules[name] = module
}

// Get retrieves a module by name from the registry.
// Returns nil if the module doesn't exist.
//
// This is used internally by the require() function to resolve modules.
func (r *Registry) Get(name string) Module {
	return r.modules[name]
}

// List returns a slice of all registered module names.
// Useful for debugging or displaying available modules.
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.modules))
	for name := range r.modules {
		names = append(names, name)
	}
	return names
}
