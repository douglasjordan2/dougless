package modules

import (
  "sync"

  "github.com/dop251/goja"

  "github.com/douglasjordan2/dougless/internal/event"
)

type PromiseState int

const (
  PromisePending PromiseState = iota
  PromiseFulfilled
  PromiseRejected
)

type Promise struct {
  vm          *goja.Runtime
  eventLoop   *event.Loop
  state       PromiseState
  value       goja.Value
  reason      goja.Value
  onFulfilled []goja.Callable
  onRejected  []goja.Callable
  mu          sync.Mutex
}

func NewPromise(vm *goja.Runtime, eventLoop *event.Loop, executor goja.Callable) *Promise {
  p := &Promise{
    vm:          vm,
    eventLoop:   eventLoop,
    state:       PromisePending,
    onFulfilled: []goja.Callable{},
    onRejected:  []goja.Callable{},
  }

  resolve := func(call goja.FunctionCall) goja.Value {
    p.resolve(call.Argument(0))
    return goja.Undefined()
  }

  reject := func(call goja.FunctionCall) goja.Value {
    p.reject(call.Argument(0))
    return goja.Undefined()
  }

  _, err := executor(goja.Undefined(), vm.ToValue(resolve), vm.ToValue(reject))
  if err != nil {
    p.reject(vm.ToValue(err.Error()))
  }

  return p
}

func (p *Promise) resolve(value goja.Value) {
  p.mu.Lock()
  defer p.mu.Unlock()

  if p.state != PromisePending {
    return // already settled
  }
  p.state = PromiseFulfilled
  p.value = value

  for _, handler := range p.onFulfilled {
    h := handler // capture for closure
    v := value   // capture value
    p.eventLoop.ScheduleTask(&event.Task{
      Callback: func() {
        h(goja.Undefined(), v)
      },
    })
  }

  p.onFulfilled = nil
  p.onRejected = nil
}

func (p *Promise) reject(reason goja.Value) {
  p.mu.Lock()
  defer p.mu.Unlock()

  if p.state != PromisePending {
    return // already settled
  }
  p.state = PromiseRejected
  p.reason = reason

  for _, handler := range p.onRejected {
    h := handler // capture for closure
    r := reason  // capture reason
    p.eventLoop.ScheduleTask(&event.Task{
      Callback: func() {
        h(goja.Undefined(), r)
      },
    })
  }

  p.onFulfilled = nil
  p.onRejected = nil
}

func (p *Promise) Then(onFulfilled, onRejected goja.Callable) *Promise {
  newPromise := &Promise{
    vm:          p.vm,
    eventLoop:   p.eventLoop,
    state:       PromisePending,
    onFulfilled: []goja.Callable{},
    onRejected:  []goja.Callable{},
  }

  fulfilledWrapper := func(call goja.FunctionCall) goja.Value {
    if onFulfilled == nil {
      newPromise.resolve(call.Argument(0))
      return goja.Undefined()
    }

    result, err := onFulfilled(goja.Undefined(), call.Argument(0))
    if err != nil {
      newPromise.reject(p.vm.ToValue(err.Error()))
    } else {
      if resultPromise, ok := result.Export().(*Promise); ok {
        // Chain the promises
        resolveFn := func(call goja.FunctionCall) goja.Value {
          newPromise.resolve(call.Argument(0))
          return goja.Undefined()
        }
        rejectFn := func(call goja.FunctionCall) goja.Value {
          newPromise.reject(call.Argument(0))
          return goja.Undefined()
        }
        
        resolveCallable, _ := goja.AssertFunction(p.vm.ToValue(resolveFn))
        rejectCallable, _ := goja.AssertFunction(p.vm.ToValue(rejectFn))
        
        resultPromise.Then(resolveCallable, rejectCallable)
      } else {
        newPromise.resolve(result)
      }
    }

    return goja.Undefined()
  }

  rejectedWrapper := func(call goja.FunctionCall) goja.Value {
    if onRejected == nil {
      newPromise.reject(call.Argument(0))
      return goja.Undefined()
    }

    result, err := onRejected(goja.Undefined(), call.Argument(0))
    if err != nil {
      newPromise.reject(p.vm.ToValue(err.Error()))
    } else {
      newPromise.resolve(result)
    }

    return goja.Undefined()
  }

  p.mu.Lock()
  defer p.mu.Unlock()

  if p.state == PromisePending {
    // Always attach fulfilledWrapper to propagate resolution through the chain
    wrappedFulfilled, _ := goja.AssertFunction(p.vm.ToValue(fulfilledWrapper))
    p.onFulfilled = append(p.onFulfilled, wrappedFulfilled)
    
    // Always attach rejectedWrapper to propagate rejection through the chain
    wrappedRejected, _ := goja.AssertFunction(p.vm.ToValue(rejectedWrapper))
    p.onRejected = append(p.onRejected, wrappedRejected)
  } else if p.state == PromiseFulfilled {
    if onFulfilled != nil {
      val := p.value
      p.eventLoop.ScheduleTask(&event.Task{
        Callback: func() {
          fulfilledWrapper(goja.FunctionCall{Arguments: []goja.Value{val}})
        },
      })
    }
  } else if p.state == PromiseRejected {
    if onRejected != nil {
      rsn := p.reason
      p.eventLoop.ScheduleTask(&event.Task{
        Callback: func() {
          rejectedWrapper(goja.FunctionCall{Arguments: []goja.Value{rsn}})
        },
      })
    }
  }

  return newPromise
}

func (p *Promise) Catch(onRejected goja.Callable) *Promise {
  return p.Then(nil, onRejected)
}

func SetupPromise(vm *goja.Runtime, eventLoop *event.Loop) {
  // Create the Promise constructor
  promiseConstructor := func(call goja.ConstructorCall) *goja.Object {
    executor, ok := goja.AssertFunction(call.Argument(0))
    if !ok {
      panic(vm.NewTypeError("Promise executor must be a function"))
    }

    promise := NewPromise(vm, eventLoop, executor)
    obj := vm.NewObject()

    obj.Set("then", func(call goja.FunctionCall) goja.Value {
      onFulfilled, _ := goja.AssertFunction(call.Argument(0))
      onRejected, _ := goja.AssertFunction(call.Argument(1))
      newPromise := promise.Then(onFulfilled, onRejected)
      
      // Create a new object for the returned promise
      newObj := vm.NewObject()
      newObj.Set("then", func(call goja.FunctionCall) goja.Value {
        onF, _ := goja.AssertFunction(call.Argument(0))
        onR, _ := goja.AssertFunction(call.Argument(1))
        return vm.ToValue(newPromise.Then(onF, onR))
      })
      newObj.Set("catch", func(call goja.FunctionCall) goja.Value {
        onR, _ := goja.AssertFunction(call.Argument(0))
        return vm.ToValue(newPromise.Catch(onR))
      })
      return newObj
    })

    obj.Set("catch", func(call goja.FunctionCall) goja.Value {
      onRejected, _ := goja.AssertFunction(call.Argument(0))
      newPromise := promise.Catch(onRejected)
      
      // Create a new object for the returned promise
      newObj := vm.NewObject()
      newObj.Set("then", func(call goja.FunctionCall) goja.Value {
        onF, _ := goja.AssertFunction(call.Argument(0))
        onR, _ := goja.AssertFunction(call.Argument(1))
        return vm.ToValue(newPromise.Then(onF, onR))
      })
      newObj.Set("catch", func(call goja.FunctionCall) goja.Value {
        onR, _ := goja.AssertFunction(call.Argument(0))
        return vm.ToValue(newPromise.Catch(onR))
      })
      return newObj
    })

    return obj
  }

  // Helper to create a promise object with methods
  var createPromiseObject func(promise *Promise) goja.Value
  createPromiseObject = func(promise *Promise) goja.Value {
    obj := vm.NewObject()
    obj.Set("then", func(call goja.FunctionCall) goja.Value {
      onFulfilled, _ := goja.AssertFunction(call.Argument(0))
      onRejected, _ := goja.AssertFunction(call.Argument(1))
      newPromise := promise.Then(onFulfilled, onRejected)
      return createPromiseObject(newPromise)
    })
    obj.Set("catch", func(call goja.FunctionCall) goja.Value {
      onRejected, _ := goja.AssertFunction(call.Argument(0))
      newPromise := promise.Catch(onRejected)
      return createPromiseObject(newPromise)
    })
    return obj
  }

  // Set up the Promise constructor function
  promiseFunc := vm.ToValue(promiseConstructor)
  promiseFuncObj := promiseFunc.ToObject(vm)

  // Promise.resolve()
  promiseFuncObj.Set("resolve", func(call goja.FunctionCall) goja.Value {
    value := call.Argument(0)
    promise := &Promise{
      vm:          vm,
      eventLoop:   eventLoop,
      state:       PromiseFulfilled,
      value:       value,
      onFulfilled: []goja.Callable{},
      onRejected:  []goja.Callable{},
    }
    return createPromiseObject(promise)
  })

  // Promise.reject()
  promiseFuncObj.Set("reject", func(call goja.FunctionCall) goja.Value {
    reason := call.Argument(0)
    promise := &Promise{
      vm:          vm,
      eventLoop:   eventLoop,
      state:       PromiseRejected,
      reason:      reason,
      onFulfilled: []goja.Callable{},
      onRejected:  []goja.Callable{},
    }
    return createPromiseObject(promise)
  })

  promiseFuncObj.Set("all", func(call goja.FunctionCall) goja.Value {
    promisesArg := call.Argument(0)

    if promisesArg.ExportType() == nil {
      panic(vm.NewTypeError("Promise.all requires an iterable"))
    }

    promisesObj := promisesArg.ToObject(vm)
    lengthVal := promisesObj.Get("length") 
    if lengthVal == nil {
      panic(vm.NewTypeError("Promise.all requires an array"))
    }

    length := int(lengthVal.ToInteger())

    if length == 0 {
      emptyPromise := &Promise{
        vm:        vm,
        eventLoop: eventLoop,
        state:     PromiseFulfilled,
        value:     vm.ToValue([]interface{}{}),
      }
      return createPromiseObject(emptyPromise)
    }

    allPromise := &Promise{
      vm:          vm,
      eventLoop:   eventLoop,
      state:       PromisePending,
      onFulfilled: []goja.Callable{},
      onRejected:  []goja.Callable{},
    }

    results := make([]goja.Value, length)
    var remaining = length
    var mu sync.Mutex
    var rejected = false

    for i := 0; i < length; i++ {
      index := i // capture for closure
      promiseVal := promisesObj.Get(string(rune('0' + i)))

      if promiseVal.ExportType() == nil || promiseVal.ToObject(vm).Get("then") == nil {
        // not a promise, treat as resolved value
        mu.Lock()
        results[index] = promiseVal
        remaining--
        if remaining == 0 && !rejected {
          allPromise.resolve(vm.ToValue(results))
        }
        mu.Unlock()
        continue
      }

      // this is a promise, attach handlers
      promiseObj := promiseVal.ToObject(vm)
      thenFunc, ok := goja.AssertFunction(promiseObj.Get("then"))
      if !ok {
        continue
      }

      successHandler := func(call goja.FunctionCall) goja.Value {
        mu.Lock()
        defer mu.Unlock()

        if rejected {
          return goja.Undefined()
        }

        results[index] = call.Argument(0)
        remaining--

        if remaining == 0 {
          allPromise.resolve(vm.ToValue(results))
        }

        return goja.Undefined()
      }

      errorHandler := func(call goja.FunctionCall) goja.Value {
        mu.Lock()
        defer mu.Unlock()

        if !rejected {
          rejected = true
          allPromise.reject(call.Argument(0))
        }

        return goja.Undefined()
      }

      thenFunc(goja.Undefined(), vm.ToValue(successHandler), vm.ToValue(errorHandler))
    }

    return createPromiseObject(allPromise)
  })

  promiseFuncObj.Set("race", func(call goja.FunctionCall) goja.Value {
    promisesArg := call.Argument(0)

    if promisesArg.ExportType() == nil {
      panic(vm.NewTypeError("Promise.race requires an iterable"))
    }

    promisesObj := promisesArg.ToObject(vm)
    lengthVal := promisesObj.Get("length")
    if lengthVal == nil {
      panic(vm.NewTypeError("Promise.race requires an array"))
    }

    length := int(lengthVal.ToInteger())

    racePromise := &Promise{
      vm:          vm,
      eventLoop:   eventLoop,
      state:       PromisePending,
      onFulfilled: []goja.Callable{},
      onRejected:  []goja.Callable{},
    }

    if length == 0 { // empty array returns a forever-pending promise
      return createPromiseObject(racePromise)
    }

    var mu sync.Mutex
    var settled = false

    // Iterate in reverse order so first promise wins when multiple are already resolved
    // (due to goroutine scheduling, last-scheduled task often executes first)
    for i := length - 1; i >= 0; i-- {
      promiseVal := promisesObj.Get(string(rune('0' + i)))

      if promiseVal.ExportType() == nil || promiseVal.ToObject(vm).Get("then") == nil {
        // not a promise, treat as resolved. race is won immediately
        if !settled {
          settled = true
          racePromise.resolve(promiseVal)
        }
        return createPromiseObject(racePromise)
      }

      // this is a promise
      promiseObj := promiseVal.ToObject(vm)
      thenFunc, ok := goja.AssertFunction(promiseObj.Get("then"))
      if !ok {
        continue
      }

      successHandler := func(call goja.FunctionCall) goja.Value {
        mu.Lock()
        defer mu.Unlock()

        if !settled {
          settled = true
          racePromise.resolve(call.Argument(0))
        }

        return goja.Undefined()
      }

      errorHandler := func(call goja.FunctionCall) goja.Value {
        mu.Lock()
        defer mu.Unlock()

        if !settled {
          settled = true
          racePromise.reject(call.Argument(0))
        }

        return goja.Undefined()
      }

      thenFunc(goja.Undefined(), vm.ToValue(successHandler), vm.ToValue(errorHandler))
    }

    return createPromiseObject(racePromise)
  })

  vm.Set("Promise", promiseFuncObj)
}
