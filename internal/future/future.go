package future

import (
  "fmt"
  "sync"
)

type Future struct {
  value any
  err   error
  ready chan struct{}
  once  sync.Once
}

func NewFuture(fn func() (any, error)) *Future {
  f := &Future{
    ready: make(chan struct{}),
  }

  go func() {
    f.value, f.err = fn()
    f.once.Do(func() {
      close(f.ready)
    })
  }()

  return f
}

func (f *Future) Get() (any, error) {
  <-f.ready
  return f.value, f.err
}

func (f *Future) MustGet() any {
  val, err := f.Get()
  if err != nil {
    panic(fmt.Sprintf("future failed: %v", err))
  }
  return val
}
