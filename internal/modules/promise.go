package modules

import (
	"strconv"
	"sync"

	"github.com/dop251/goja"
)

type PromiseState int

const (
	PromisePending   PromiseState = iota
	PromiseFulfilled
	PromiseRejected
)

type Promise struct {
	vm          *goja.Runtime
  runtime     RuntimeKeepAlive
	state       PromiseState
	value       goja.Value
	reason      goja.Value
	onFulfilled []goja.Callable
	onRejected  []goja.Callable
	mu          sync.Mutex
}

func NewPromise(vm *goja.Runtime, executor goja.Callable) *Promise {
	p := &Promise{
		vm:          vm,
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

func (p *Promise) SetRuntime(rt RuntimeKeepAlive) {
  p.runtime = rt
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

    done := p.runtime.KeepAlive()
    go func() {
      defer done()
      h(goja.Undefined(), v)
		}()
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
    done := p.runtime.KeepAlive()
    go func() {
      defer done()
      h(goja.Undefined(), r)
		}()
	}

	p.onFulfilled = nil
	p.onRejected = nil
}

func (p *Promise) Then(onFulfilled, onRejected goja.Callable) *Promise {
	newPromise := &Promise{
		vm:          p.vm,
    runtime:     p.runtime,
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
			return goja.Undefined()
		}

		if result != nil && !goja.IsUndefined(result) && !goja.IsNull(result) {
			resultObj := result.ToObject(p.vm)
			if resultObj != nil {
				thenMethod := resultObj.Get("then")
				if thenMethod != nil && !goja.IsUndefined(thenMethod) && !goja.IsNull(thenMethod) {
					if thenFunc, ok := goja.AssertFunction(thenMethod); ok {
						resolveFn := func(call goja.FunctionCall) goja.Value {
							newPromise.resolve(call.Argument(0))
							return goja.Undefined()
						}
						rejectFn := func(call goja.FunctionCall) goja.Value {
							newPromise.reject(call.Argument(0))
							return goja.Undefined()
						}

						thenFunc(result, p.vm.ToValue(resolveFn), p.vm.ToValue(rejectFn))
						return goja.Undefined()
					}
				}
			}
		}

		newPromise.resolve(result)
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

  switch p.state {
  case PromisePending:
    wrappedFulfilled, _ := goja.AssertFunction(p.vm.ToValue(fulfilledWrapper))
    p.onFulfilled = append(p.onFulfilled, wrappedFulfilled)

    wrappedRejected, _ := goja.AssertFunction(p.vm.ToValue(rejectedWrapper))
    p.onRejected = append(p.onRejected, wrappedRejected)
  case PromiseFulfilled:
    val := p.value
    fulfilledWrapper(goja.FunctionCall{Arguments: []goja.Value{val}})
  case PromiseRejected:
    rsn := p.reason
    rejectedWrapper(goja.FunctionCall{Arguments: []goja.Value{rsn}})
  }

	return newPromise
}

func (p *Promise) Catch(onRejected goja.Callable) *Promise {
	return p.Then(nil, onRejected)
}

func CreatePromiseObject(vm *goja.Runtime, promise *Promise) goja.Value {
	obj := vm.NewObject()

	obj.Set("then", func(call goja.FunctionCall) goja.Value {
		onFulfilled, _ := goja.AssertFunction(call.Argument(0))
		onRejected, _ := goja.AssertFunction(call.Argument(1))
		newPromise := promise.Then(onFulfilled, onRejected)
		return CreatePromiseObject(vm, newPromise)
	})

	obj.Set("catch", func(call goja.FunctionCall) goja.Value {
		onRejected, _ := goja.AssertFunction(call.Argument(0))
		newPromise := promise.Catch(onRejected)
		return CreatePromiseObject(vm, newPromise)
	})

	return obj
}

func SetupPromise(vm *goja.Runtime, rt RuntimeKeepAlive) {
	promiseConstructor := func(call goja.ConstructorCall) *goja.Object {
		executor, ok := goja.AssertFunction(call.Argument(0))
		if !ok {
			panic(vm.NewTypeError("Promise executor must be a function"))
		}

		promise := NewPromise(vm, executor)
    promise.SetRuntime(rt)

		obj := vm.NewObject()
		obj.Set("then", func(call goja.FunctionCall) goja.Value {
			onFulfilled, _ := goja.AssertFunction(call.Argument(0))
			onRejected, _ := goja.AssertFunction(call.Argument(1))
			newPromise := promise.Then(onFulfilled, onRejected)

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

	promiseFunc := vm.ToValue(promiseConstructor)
	promiseFuncObj := promiseFunc.ToObject(vm)

	promiseFuncObj.Set("resolve", func(call goja.FunctionCall) goja.Value {
		value := call.Argument(0)
		promise := &Promise{
			vm:          vm,
			runtime:     rt,
			state:       PromiseFulfilled,
			value:       value,
			onFulfilled: []goja.Callable{},
			onRejected:  []goja.Callable{},
		}
		return CreatePromiseObject(vm, promise)
	})

	promiseFuncObj.Set("reject", func(call goja.FunctionCall) goja.Value {
		reason := call.Argument(0)
		promise := &Promise{
			vm:          vm,
			runtime:     rt,
			state:       PromiseRejected,
			reason:      reason,
			onFulfilled: []goja.Callable{},
			onRejected:  []goja.Callable{},
		}
		return CreatePromiseObject(vm, promise)
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
				runtime:   rt,
				state:     PromiseFulfilled,
				value:     vm.ToValue([]any{}),
			}
			return CreatePromiseObject(vm, emptyPromise)
		}

		allPromise := &Promise{
			vm:          vm,
			runtime:     rt,
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
			promiseVal := promisesObj.Get(strconv.Itoa(i))

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

		return CreatePromiseObject(vm, allPromise)
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
		runtime:     rt,
		state:       PromisePending,
		onFulfilled: []goja.Callable{},
		onRejected:  []goja.Callable{},
	}

		if length == 0 { // empty array returns a forever-pending promise
			return CreatePromiseObject(vm, racePromise)
		}

		var mu sync.Mutex
		var settled = false

		for i := 0; i < length; i++ {
			promiseVal := promisesObj.Get(strconv.Itoa(i))

			if promiseVal.ExportType() == nil || promiseVal.ToObject(vm).Get("then") == nil {
				// not a promise, treat as resolved. race is won immediately
				if !settled {
					settled = true
					racePromise.resolve(promiseVal)
				}
				return CreatePromiseObject(vm, racePromise)
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

		return CreatePromiseObject(vm, racePromise)
	})

	promiseFuncObj.Set("any", func(call goja.FunctionCall) goja.Value {
		promisesArg := call.Argument(0)

		if promisesArg.ExportType() == nil {
			panic(vm.NewTypeError("Promise.any requires an iterable"))
		}

		promisesObj := promisesArg.ToObject(vm)
		lengthVal := promisesObj.Get("length")
		if lengthVal == nil {
			panic(vm.NewTypeError("Promise.any requires an array"))
		}

		length := int(lengthVal.ToInteger())

	anyPromise := &Promise{
		vm:          vm,
		runtime:     rt,
		state:       PromisePending,
		onFulfilled: []goja.Callable{},
		onRejected:  []goja.Callable{},
	}

		if length == 0 {
			aggregateError := vm.NewObject()
			aggregateError.Set("name", "AggregateError")
			aggregateError.Set("message", "All promises were rejected")
			aggregateError.Set("errors", vm.ToValue([]goja.Value{}))
			anyPromise.reject(aggregateError)
			return CreatePromiseObject(vm, anyPromise)
		}

		errors := make([]goja.Value, length)
		var mu sync.Mutex
		var settled = false
		var remaining = length

		for i := 0; i < length; i++ {
			index := i // closure stuff again
			promiseVal := promisesObj.Get(strconv.Itoa(i))

			if promiseVal.ExportType() == nil || promiseVal.ToObject(vm).Get("then") == nil {
				// not a promise, return immediately
				mu.Lock()
				anyPromise.resolve(promiseVal)
				mu.Unlock()

				return CreatePromiseObject(vm, anyPromise)
			}

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
					anyPromise.resolve(call.Argument(0))
				}

				return goja.Undefined()
			}

			errorHandler := func(call goja.FunctionCall) goja.Value {
				mu.Lock()
				defer mu.Unlock()

				if settled {
					return goja.Undefined()
				}

				errors[index] = call.Argument(0)
				remaining--

				if remaining == 0 {
					aggregateError := vm.NewObject()
					aggregateError.Set("name", "AggregateError")
					aggregateError.Set("message", "All promises were rejected")
					aggregateError.Set("errors", vm.ToValue(errors))
					anyPromise.reject(aggregateError)
				}

				return goja.Undefined()
			}

			thenFunc(goja.Undefined(), vm.ToValue(successHandler), vm.ToValue(errorHandler))
		}

		return CreatePromiseObject(vm, anyPromise)
	})

	promiseFuncObj.Set("allSettled", func(call goja.FunctionCall) goja.Value {
		promisesArg := call.Argument(0)

		if promisesArg.ExportType() == nil {
			panic(vm.NewTypeError("Promise.allSettled requires an iterable"))
		}

		promisesObj := promisesArg.ToObject(vm)
		lengthVal := promisesObj.Get("length")
		if lengthVal == nil {
			panic(vm.NewTypeError("Promise.allSettled requires an array"))
		}

		length := int(lengthVal.ToInteger())

	allSettledPromise := &Promise{
		vm:          vm,
		runtime:     rt,
		state:       PromisePending,
		onFulfilled: []goja.Callable{},
		onRejected:  []goja.Callable{},
	}

		if length == 0 {
			allSettledPromise.resolve(vm.ToValue([]goja.Value{}))
			return CreatePromiseObject(vm, allSettledPromise)
		}

		results := make([]goja.Value, length)
		var mu sync.Mutex
		var remaining = length

		for i := 0; i < length; i++ {
			index := i
			promiseVal := promisesObj.Get(strconv.Itoa(i))

			if promiseVal.ExportType() == nil || promiseVal.ToObject(vm).Get("then") == nil {
				// not a promise
				mu.Lock()

				resultObj := vm.NewObject()
				resultObj.Set("status", "fulfilled")
				resultObj.Set("value", promiseVal)
				results[index] = resultObj

				remaining--
				if remaining == 0 {
					allSettledPromise.resolve(vm.ToValue(results))
				}

				mu.Unlock()

				continue
			}

			promiseObj := promiseVal.ToObject(vm)
			thenFunc, ok := goja.AssertFunction(promiseObj.Get("then"))

			if !ok {
				continue
			}

			successHandler := func(call goja.FunctionCall) goja.Value {
				mu.Lock()
				defer mu.Unlock()

				resultObj := vm.NewObject()
				resultObj.Set("status", "fulfilled")
				resultObj.Set("value", call.Argument(0))
				results[index] = resultObj

				remaining--
				if remaining == 0 {
					allSettledPromise.resolve(vm.ToValue(results))
				}

				return goja.Undefined()
			}

			errorHandler := func(call goja.FunctionCall) goja.Value {
				mu.Lock()
				defer mu.Unlock()

				resultObj := vm.NewObject()
				resultObj.Set("status", "rejected")
				resultObj.Set("reason", call.Argument(0))
				results[index] = resultObj

				remaining--
				if remaining == 0 {
					allSettledPromise.resolve(vm.ToValue(results))
				}

				return goja.Undefined()
			}

			thenFunc(goja.Undefined(), vm.ToValue(successHandler), vm.ToValue(errorHandler))
		}

		return CreatePromiseObject(vm, allSettledPromise)
	})

	vm.Set("Promise", promiseFuncObj)
}
