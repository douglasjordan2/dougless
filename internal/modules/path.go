package modules

import (
  "path/filepath"
  "strings"

  "github.com/dop251/goja"
)

type Path struct {
  vm *goja.Runtime
}

func NewPath() *Path {
  return &Path{}
}

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

func argToStr(call goja.FunctionCall) []string {
  // convert all args to strings
  parts := make([]string, len(call.Arguments))
  for i, arg := range call.Arguments {
    parts[i] = arg.String()
  }

  return parts
}

// join path segments together
func (p *Path) join(call goja.FunctionCall) goja.Value {
  parts := argToStr(call)
  result := filepath.Join(parts...)

  return p.vm.ToValue(result)
}

// resolve an absolute path
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

// returns directory part of path
func (p *Path) dirname(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 1 {
    return p.vm.ToValue("")
  }

  path := call.Arguments[0].String()
  dir := filepath.Dir(path)

  return p.vm.ToValue(dir)
}

// returns last element of a path
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

// returns the file extension
func (p *Path) extname(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 1 {
    return p.vm.ToValue("")
  }

  path := call.Arguments[0].String()
  ext := filepath.Ext(path)

  return p.vm.ToValue(ext)
}
