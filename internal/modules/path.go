package modules

import (
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
)

// Path provides file path manipulation utilities for JavaScript.
// All methods use the OS-specific path separator and follow OS path conventions.
//
// Available globally in JavaScript as the 'path' object, or via require('path').
//
// Example usage:
//
//	path.join('foo', 'bar', 'baz.txt')  // 'foo/bar/baz.txt'
//	path.dirname('/home/user/file.txt')  // '/home/user'
//	path.extname('file.txt')  // '.txt'
type Path struct {
	vm *goja.Runtime // JavaScript runtime instance
}

// NewPath creates a new Path module instance.
func NewPath() *Path {
	return &Path{}
}

// Export creates and returns the path JavaScript object with all path methods.
func (p *Path) Export(vm *goja.Runtime) goja.Value {
	p.vm = vm
	obj := vm.NewObject()

	obj.Set("join", p.join)
	obj.Set("resolve", p.resolve)
	obj.Set("dirname", p.dirname)
	obj.Set("basename", p.basename)
	obj.Set("extname", p.extname)
	obj.Set("sep", filepath.Separator)

	return obj
}

// argToStr converts all function call arguments to strings.
// Helper function used by path methods that accept multiple string arguments.
func argToStr(call goja.FunctionCall) []string {
	parts := make([]string, len(call.Arguments))
	for i, arg := range call.Arguments {
		parts[i] = arg.String()
	}
	return parts
}

// join implements path.join() - joins path segments using the OS path separator.
//
// JavaScript usage:
//
//	path.join('foo', 'bar', 'baz')  // 'foo/bar/baz' on Unix
//	path.join('/usr', 'local', 'bin')  // '/usr/local/bin'
func (p *Path) join(call goja.FunctionCall) goja.Value {
	parts := argToStr(call)
	result := filepath.Join(parts...)

	return p.vm.ToValue(result)
}

// resolve implements path.resolve() - resolves path segments to an absolute path.
// The resulting path is normalized and resolved relative to the current working directory.
//
// JavaScript usage:
//
//	path.resolve('foo', 'bar')  // '/current/working/dir/foo/bar'
//	path.resolve('/home', 'user', 'file.txt')  // '/home/user/file.txt'
//
// Panics if the path cannot be resolved to an absolute path.
func (p *Path) resolve(call goja.FunctionCall) goja.Value {
	parts := argToStr(call)
	joined := filepath.Join(parts...)

	// convert to absolute path
	absolute, err := filepath.Abs(joined)
	if err != nil {
		panic(p.vm.NewGoError(err))
	}

	return p.vm.ToValue(absolute)
}

// dirname implements path.dirname() - returns the directory portion of a path.
// Trailing path separators are ignored.
//
// JavaScript usage:
//
//	path.dirname('/home/user/file.txt')  // '/home/user'
//	path.dirname('src/index.js')  // 'src'
//
// Returns an empty string if no path argument is provided.
func (p *Path) dirname(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 1 {
		return p.vm.ToValue("")
	}

	path := call.Arguments[0].String()
	dir := filepath.Dir(path)

	return p.vm.ToValue(dir)
}

// basename implements path.basename() - returns the last element of a path.
// Optionally removes a specified extension.
//
// JavaScript usage:
//
//	path.basename('/home/user/file.txt')  // 'file.txt'
//	path.basename('/home/user/file.txt', '.txt')  // 'file'
//
// Returns an empty string if no path argument is provided.
func (p *Path) basename(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 1 {
		return p.vm.ToValue("")
	}

	path := call.Arguments[0].String()

	if len(call.Arguments) >= 2 {
		ext := call.Arguments[1].String()
		base := filepath.Base(path)
		return p.vm.ToValue(strings.TrimSuffix(base, ext))
	}

	return p.vm.ToValue(filepath.Base(path))
}

// extname implements path.extname() - returns the file extension including the dot.
//
// JavaScript usage:
//
//	path.extname('file.txt')  // '.txt'
//	path.extname('archive.tar.gz')  // '.gz'
//	path.extname('README')  // ''
//
// Returns an empty string if there's no extension or no path argument is provided.
func (p *Path) extname(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 1 {
		return p.vm.ToValue("")
	}

	path := call.Arguments[0].String()
	ext := filepath.Ext(path)

	return p.vm.ToValue(ext)
}
