package modules

import (
	"strconv"
	"sync"

	"github.com/dop251/goja"

	"github.com/douglasjordan2/dougless/internal/event"
)

// PromiseState represents the current state of a Promise.
type PromiseState int

// Promise states as defined by the Promise/A+ specification.
const (
	PromisePending   PromiseState = iota // Promise is pending (initial state)
	PromiseFulfilled                     // Promise has been resolved with a value
	PromiseRejected                      // Promise has been rejected with a reason
)

// Promise represents a Promise/A+ compliant promise implementation.
// Promises provide a way to handle asynchronous operations with chainable .then() and .catch() methods.
// All promise handlers are executed asynchronously on the event loop.
//
// Available globally in JavaScript as the 'Promise' constructor.
//
// Example usage:
//
//	const p = new Promise((resolve, reject) => {
//	  setTimeout(() => resolve(42), 1000);
//	});
//	p.then(value => console.log('Got:', value));
type Promise struct {
	vm          *goja.Runtime   // JavaScript runtime instance
	eventLoop   *event.Loop     // Event loop for async handler execution
	state       PromiseState    // Current promise state (pending, fulfilled, or rejected)
	value       goja.Value      // Resolved value (when fulfilled)
	reason      goja.Value      // Rejection reason (when rejected)
	onFulfilled []goja.Callable // Handlers to call when promise is fulfilled
	onRejected  []goja.Callable // Handlers to call when promise is rejected
	mu          sync.Mutex      // Protects state changes and handler lists
}

// NewPromise creates a new Promise instance and executes the executor function.
// The executor is called immediately with resolve and reject callback functions.
//
// Parameters:
//   - vm: JavaScript runtime instance
//   - eventLoop: Event loop for scheduling async handlers
//   - executor: Function called with (resolve, reject) callbacks
//
// If the executor throws an error, the promise is automatically rejected.
// This function is used internally by the Promise constructor in JavaScript.
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

// resolve transitions the promise from pending to fulfilled state.
// All registered fulfillment handlers are scheduled on the event loop.
// If the promise is already settled, this is a no-op (per Promise/A+ spec).
//
// This method is thread-safe and idempotent.
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

// reject transitions the promise from pending to rejected state.
// All registered rejection handlers are scheduled on the event loop.
// If the promise is already settled, this is a no-op (per Promise/A+ spec).
//
// This method is thread-safe and idempotent.
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

// Then implements promise chaining by attaching fulfillment and rejection handlers.
// Returns a new Promise that will be resolved/rejected based on the handler's return value.
//
// Parameters:
//   - onFulfilled: Called when the promise is fulfilled (can be nil)
//   - onRejected: Called when the promise is rejected (can be nil)
//
// Behavior (per Promise/A+ spec):
//   - If onFulfilled returns a value, the new promise is fulfilled with that value
//   - If onFulfilled returns a thenable (has .then method), it's chained
//   - If onFulfilled throws, the new promise is rejected
//   - If onFulfilled is nil, the value propagates to the next .then()
//   - Similar rules apply for onRejected
//
// All handlers are executed asynchronously on the event loop.
//
// Example:
//
//	promise.then(
//	  value => console.log('Success:', value),
//	  error => console.error('Error:', error)
//	);
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
			return goja.Undefined()
		}

		// Check if result is a thenable (has a .then method)
		if result != nil && !goja.IsUndefined(result) && !goja.IsNull(result) {
			resultObj := result.ToObject(p.vm)
			if resultObj != nil {
				thenMethod := resultObj.Get("then")
				if thenMethod != nil && !goja.IsUndefined(thenMethod) && !goja.IsNull(thenMethod) {
					if thenFunc, ok := goja.AssertFunction(thenMethod); ok {
						// It's thenable - chain it
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

		// Not a promise, just resolve with value
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
		// Always attach fulfilledWrapper to propagate resolution through the chain
		wrappedFulfilled, _ := goja.AssertFunction(p.vm.ToValue(fulfilledWrapper))
		p.onFulfilled = append(p.onFulfilled, wrappedFulfilled)

		// Always attach rejectedWrapper to propagate rejection through the chain
		wrappedRejected, _ := goja.AssertFunction(p.vm.ToValue(rejectedWrapper))
		p.onRejected = append(p.onRejected, wrappedRejected)
	case PromiseFulfilled:
		// Always propagate value (fulfilledWrapper handles nil handler)
		val := p.value
		p.eventLoop.ScheduleTask(&event.Task{
			Callback: func() {
				fulfilledWrapper(goja.FunctionCall{Arguments: []goja.Value{val}})
			},
		})
	case PromiseRejected:
		// Always propagate rejection (rejectedWrapper handles nil handler)
		rsn := p.reason
		p.eventLoop.ScheduleTask(&event.Task{
			Callback: func() {
				rejectedWrapper(goja.FunctionCall{Arguments: []goja.Value{rsn}})
			},
		})
	}

	return newPromise
}

// Catch is a convenience method for handling promise rejections.
// It's equivalent to calling .then(nil, onRejected).
//
// Parameters:
//   - onRejected: Called when the promise is rejected
//
// Returns a new Promise that resolves to the return value of onRejected.
//
// Example:
//
//	promise.catch(error => console.error('Error:', error));
func (p *Promise) Catch(onRejected goja.Callable) *Promise {
	return p.Then(nil, onRejected)
}

// CreatePromiseObject wraps a Promise struct into a JavaScript object
// with .then() and .catch() methods for use in JavaScript.
// This is used by modules that want to return Promises directly.
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

// SetupPromise initializes the global Promise constructor in JavaScript.
// This function is called once during runtime initialization.
//
// It creates the Promise constructor and adds static methods:
//   - Promise.resolve(value): Creates a fulfilled promise
//   - Promise.reject(reason): Creates a rejected promise
//   - Promise.all(promises): Waits for all promises to fulfill
//   - Promise.race(promises): Resolves/rejects with the first settled promise
//   - Promise.allSettled(promises): Waits for all promises to settle
//   - Promise.any(promises): Resolves with the first fulfilled promise
//
// The Promise constructor is made available globally, following ECMAScript standards.
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
		return CreatePromiseObject(vm, promise)
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
				eventLoop: eventLoop,
				state:     PromiseFulfilled,
				value:     vm.ToValue([]any{}),
			}
      return CreatePromiseObject(vm, emptyPromise)
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
			eventLoop:   eventLoop,
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
			eventLoop:   eventLoop,
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
			eventLoop:   eventLoop,
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
