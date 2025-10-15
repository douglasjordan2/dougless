package modules

import (
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/douglasjordan2/dougless/internal/event"
)

func setupTestEnvironment() (*goja.Runtime, *event.Loop) {
	vm := goja.New()
	loop := event.NewLoop()
	
	go loop.Run()
	
	SetupPromise(vm, loop)
	
	return vm, loop
}

func TestPromiseResolve(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	script := `
		var result = null;
		Promise.resolve("test value").then(function(val) {
			result = val;
		});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(50 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("result")
	if result.String() != "test value" {
		t.Errorf("Expected 'test value', got '%s'", result.String())
	}
}

func TestPromiseReject(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	script := `
		var error = null;
		Promise.reject("test error").catch(function(err) {
			error = err;
		});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(50 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("error")
	if result.String() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", result.String())
	}
}

func TestPromiseChaining(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	script := `
		var result = null;
		Promise.resolve(1)
			.then(function(val) { return val + 1; })
			.then(function(val) { return val + 1; })
			.then(function(val) { result = val; });
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(50 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("result")
	if result.ToInteger() != 3 {
		t.Errorf("Expected 3, got %d", result.ToInteger())
	}
}

func TestPromiseAllResolved(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	script := `
		var result = null;
		Promise.all([
			Promise.resolve(1),
			Promise.resolve(2),
			Promise.resolve(3)
		]).then(function(values) {
			result = values;
		});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(50 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("result")
	if result.ExportType() == nil {
		t.Fatal("Result is nil")
	}
	
	arr := result.Export().([]goja.Value)
	if len(arr) != 3 {
		t.Errorf("Expected array of length 3, got %d", len(arr))
	}
}

func TestPromiseAllRejected(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	script := `
		var error = null;
		Promise.all([
			Promise.resolve(1),
			Promise.reject("failed"),
			Promise.resolve(3)
		]).catch(function(err) {
			error = err;
		});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(50 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("error")
	if result.String() != "failed" {
		t.Errorf("Expected 'failed', got '%s'", result.String())
	}
}

func TestPromiseAllEmpty(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	script := `
		var result = null;
		Promise.all([]).then(function(values) {
			result = values;
		});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(50 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("result")
	arr := result.Export().([]interface{})
	if len(arr) != 0 {
		t.Errorf("Expected empty array, got length %d", len(arr))
	}
}

func TestPromiseRaceResolved(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	script := `
		var result = null;
		Promise.race([
			Promise.resolve("first"),
			Promise.resolve("second")
		]).then(function(val) {
			result = val;
		});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(50 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("result")
	if result.String() != "first" {
		t.Errorf("Expected 'first', got '%s'", result.String())
	}
}

func TestPromiseRaceRejected(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	script := `
		var error = null;
		Promise.race([
			Promise.reject("error"),
			Promise.resolve("success")
		]).catch(function(err) {
			error = err;
		});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(50 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("error")
	if result.String() != "error" {
		t.Errorf("Expected 'error', got '%s'", result.String())
	}
}

func TestPromiseRaceEmpty(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	script := `
		var resolved = false;
		Promise.race([]).then(function() {
			resolved = true;
		});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(50 * time.Millisecond)
	loop.Wait()
	
	// Empty race should never resolve
	resolved := vm.Get("resolved").ToBoolean()
	if resolved {
		t.Error("Empty Promise.race should never resolve")
	}
}

func TestPromiseAsyncResolution(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	// Setup setTimeout for async test
	timers := NewTimers(loop)
	timerObj := timers.Export(vm).ToObject(vm)
	vm.Set("setTimeout", timerObj.Get("setTimeout"))
	
	script := `
		var result = null;
		new Promise(function(resolve) {
			setTimeout(function() {
				resolve("async value");
			}, 50);
		}).then(function(val) {
			result = val;
		});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(150 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("result")
	if result.String() != "async value" {
		t.Errorf("Expected 'async value', got '%s'", result.String())
	}
}

func TestPromiseAsyncRejection(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	timers := NewTimers(loop)
	timerObj := timers.Export(vm).ToObject(vm)
	vm.Set("setTimeout", timerObj.Get("setTimeout"))
	
	script := `
		var error = null;
		new Promise(function(resolve, reject) {
			setTimeout(function() {
				reject("async error");
			}, 50);
		}).catch(function(err) {
			error = err;
		});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(150 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("error")
	if result.String() != "async error" {
		t.Errorf("Expected 'async error', got '%s'", result.String())
	}
}

func TestPromiseRaceAsyncMultiple(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	timers := NewTimers(loop)
	timerObj := timers.Export(vm).ToObject(vm)
	vm.Set("setTimeout", timerObj.Get("setTimeout"))
	
	script := `
		var result = null;
		Promise.race([
			new Promise(function(resolve) {
				setTimeout(function() { resolve("slow"); }, 100);
			}),
			new Promise(function(resolve) {
				setTimeout(function() { resolve("fast"); }, 50);
			})
		]).then(function(val) {
			result = val;
		});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(200 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("result")
	if result.String() != "fast" {
		t.Errorf("Expected 'fast', got '%s'", result.String())
	}
}

func TestPromiseRaceAsyncRejection(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	timers := NewTimers(loop)
	timerObj := timers.Export(vm).ToObject(vm)
	vm.Set("setTimeout", timerObj.Get("setTimeout"))
	
	script := `
		var error = null;
		Promise.race([
			new Promise(function(resolve) {
				setTimeout(function() { resolve("success"); }, 100);
			}),
			new Promise(function(resolve, reject) {
				setTimeout(function() { reject("fast error"); }, 50);
			})
		]).catch(function(err) {
			error = err;
		});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(200 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("error")
	if result.String() != "fast error" {
		t.Errorf("Expected 'fast error', got '%s'", result.String())
	}
}

func TestPromiseWithNonPromiseValues(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	script := `
		var result = null;
		Promise.all([
			42,
			"string",
			Promise.resolve("promise"),
			true
		]).then(function(values) {
			result = values;
		});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(50 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("result")
	arr := result.Export().([]goja.Value)
	
	if len(arr) != 4 {
		t.Errorf("Expected 4 elements, got %d", len(arr))
	}
	
	if arr[0].ToInteger() != 42 {
		t.Errorf("Expected 42, got %v", arr[0].ToInteger())
	}
}

func TestPromiseRaceWithNonPromiseValue(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	script := `
		var result = null;
		Promise.race([
			"instant",
			Promise.resolve("delayed")
		]).then(function(val) {
			result = val;
		});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(50 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("result")
	if result.String() != "instant" {
		t.Errorf("Expected 'instant', got '%s'", result.String())
	}
}

func TestPromiseErrorPropagation(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	script := `
		var error = null;
		Promise.resolve(1)
			.then(function(val) { 
				return Promise.reject("propagated error");
			})
			.then(function(val) {
				// This should not execute
				return val;
			})
			.catch(function(err) {
				error = err;
			});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(150 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("error")
	if result.String() != "propagated error" {
		t.Errorf("Expected 'propagated error', got '%s'", result.String())
	}
}

func TestPromiseThenWithoutErrorHandler(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	script := `
		var error = null;
		Promise.reject("error")
			.then(function(val) {
				// No error handler here
				return val;
			})
			.catch(function(err) {
				error = err;
			});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(150 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("error")
	if result.String() != "error" {
		t.Errorf("Expected 'error', got '%s'", result.String())
	}
}

func TestPromiseAllWithMixedTimings(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()
	
	timers := NewTimers(loop)
	timerObj := timers.Export(vm).ToObject(vm)
	vm.Set("setTimeout", timerObj.Get("setTimeout"))
	
	script := `
		var result = null;
		Promise.all([
			new Promise(function(resolve) {
				setTimeout(function() { resolve("slow"); }, 100);
			}),
			Promise.resolve("instant"),
			new Promise(function(resolve) {
				setTimeout(function() { resolve("medium"); }, 50);
			})
		]).then(function(values) {
			result = values;
		});
	`
	
	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	time.Sleep(200 * time.Millisecond)
	loop.Wait()
	
	result := vm.Get("result")
	arr := result.Export().([]goja.Value)
	
	if len(arr) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(arr))
	}
	
	// Order should be preserved
	if arr[0].String() != "slow" || arr[1].String() != "instant" || arr[2].String() != "medium" {
		t.Errorf("Order not preserved: got [%s, %s, %s]", arr[0].String(), arr[1].String(), arr[2].String())
	}
}

func TestPromiseAllSettledAllFulfilled(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()

	script := `
		var result = null;
		Promise.allSettled([
			Promise.resolve(1),
			Promise.resolve(2),
			Promise.resolve(3)
		]).then(function(values) {
			result = values;
		});
	`

	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	loop.Wait()

	result := vm.Get("result")
	if result.ExportType() == nil {
		t.Fatal("Result is nil")
	}

	arr := result.Export().([]goja.Value)
	if len(arr) != 3 {
		t.Errorf("Expected array of length 3, got %d", len(arr))
	}

	// Check first result
	obj0 := arr[0].ToObject(vm)
	if obj0.Get("status").String() != "fulfilled" {
		t.Errorf("Expected status 'fulfilled', got '%s'", obj0.Get("status").String())
	}
	if obj0.Get("value").ToInteger() != 1 {
		t.Errorf("Expected value 1, got %d", obj0.Get("value").ToInteger())
	}
}

func TestPromiseAllSettledAllRejected(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()

	script := `
		var result = null;
		Promise.allSettled([
			Promise.reject("error1"),
			Promise.reject("error2"),
			Promise.reject("error3")
		]).then(function(values) {
			result = values;
		});
	`

	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	loop.Wait()

	result := vm.Get("result")
	if result.ExportType() == nil {
		t.Fatal("Result is nil")
	}

	arr := result.Export().([]goja.Value)
	if len(arr) != 3 {
		t.Errorf("Expected array of length 3, got %d", len(arr))
	}

	// Check first result
	obj0 := arr[0].ToObject(vm)
	if obj0.Get("status").String() != "rejected" {
		t.Errorf("Expected status 'rejected', got '%s'", obj0.Get("status").String())
	}
	if obj0.Get("reason").String() != "error1" {
		t.Errorf("Expected reason 'error1', got '%s'", obj0.Get("reason").String())
	}
}

func TestPromiseAllSettledMixed(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()

	script := `
		var result = null;
		Promise.allSettled([
			Promise.resolve(42),
			Promise.reject("failed"),
			Promise.resolve("success")
		]).then(function(values) {
			result = values;
		});
	`

	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	loop.Wait()

	result := vm.Get("result")
	arr := result.Export().([]goja.Value)

	if len(arr) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(arr))
	}

	// First: fulfilled
	obj0 := arr[0].ToObject(vm)
	if obj0.Get("status").String() != "fulfilled" {
		t.Errorf("Expected status 'fulfilled', got '%s'", obj0.Get("status").String())
	}
	if obj0.Get("value").ToInteger() != 42 {
		t.Errorf("Expected value 42, got %d", obj0.Get("value").ToInteger())
	}

	// Second: rejected
	obj1 := arr[1].ToObject(vm)
	if obj1.Get("status").String() != "rejected" {
		t.Errorf("Expected status 'rejected', got '%s'", obj1.Get("status").String())
	}
	if obj1.Get("reason").String() != "failed" {
		t.Errorf("Expected reason 'failed', got '%s'", obj1.Get("reason").String())
	}

	// Third: fulfilled
	obj2 := arr[2].ToObject(vm)
	if obj2.Get("status").String() != "fulfilled" {
		t.Errorf("Expected status 'fulfilled', got '%s'", obj2.Get("status").String())
	}
	if obj2.Get("value").String() != "success" {
		t.Errorf("Expected value 'success', got '%s'", obj2.Get("value").String())
	}
}

func TestPromiseAllSettledEmpty(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()

	script := `
		var result = null;
		Promise.allSettled([]).then(function(values) {
			result = values;
		});
	`

	_, err := vm.RunString(script)
	if err != nil {
			t.Fatalf("Script execution failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	loop.Wait()

	result := vm.Get("result")
	arr := result.Export().([]goja.Value)
	if len(arr) != 0 {
		t.Errorf("Expected empty array, got length %d", len(arr))
	}
}

func TestPromiseAllSettledWithNonPromises(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()

	script := `
		var result = null;
		Promise.allSettled([
			42,
			"string",
			Promise.resolve("promise"),
			true
		]).then(function(values) {
			result = values;
		});
	`

	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	loop.Wait()

	result := vm.Get("result")
	arr := result.Export().([]goja.Value)

	if len(arr) != 4 {
		t.Errorf("Expected 4 elements, got %d", len(arr))
	}

	// Check that all non-promises are fulfilled
	obj0 := arr[0].ToObject(vm)
	if obj0.Get("status").String() != "fulfilled" {
		t.Errorf("Expected status 'fulfilled', got '%s'", obj0.Get("status").String())
	}
	if obj0.Get("value").ToInteger() != 42 {
		t.Errorf("Expected value 42, got %d", obj0.Get("value").ToInteger())
	}
}

func TestPromiseAllSettledNeverRejects(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()

	script := `
		var rejected = false;
		var resolved = false;
		Promise.allSettled([
			Promise.reject("error1"),
			Promise.reject("error2")
		]).then(function(values) {
			resolved = true;
		}).catch(function(err) {
			rejected = true;
		});
	`

	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	loop.Wait()

	resolved := vm.Get("resolved").ToBoolean()
	rejected := vm.Get("rejected").ToBoolean()

	if !resolved {
		t.Error("Promise.allSettled should resolve even when all promises reject")
	}
	if rejected {
		t.Error("Promise.allSettled should never reject")
	}
}

func TestPromiseAllSettledWithMixedTimings(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()

	timers := NewTimers(loop)
	timerObj := timers.Export(vm).ToObject(vm)
	vm.Set("setTimeout", timerObj.Get("setTimeout"))

	script := `
		var result = null;
		Promise.allSettled([
			new Promise(function(resolve) {
				setTimeout(function() { resolve("slow"); }, 100);
			}),
			Promise.resolve("instant"),
			new Promise(function(resolve, reject) {
				setTimeout(function() { reject("medium error"); }, 50);
			})
		]).then(function(values) {
			result = values;
		});
	`

	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	time.Sleep(200 * time.Millisecond)
	loop.Wait()

	result := vm.Get("result")
	arr := result.Export().([]goja.Value)

	if len(arr) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(arr))
	}

	// Order should be preserved
	obj0 := arr[0].ToObject(vm)
	if obj0.Get("status").String() != "fulfilled" || obj0.Get("value").String() != "slow" {
		t.Errorf("First promise result incorrect")
	}

	obj1 := arr[1].ToObject(vm)
	if obj1.Get("status").String() != "fulfilled" || obj1.Get("value").String() != "instant" {
		t.Errorf("Second promise result incorrect")
	}

	obj2 := arr[2].ToObject(vm)
	if obj2.Get("status").String() != "rejected" || obj2.Get("reason").String() != "medium error" {
		t.Errorf("Third promise result incorrect")
	}
}

func TestPromiseAnyFirstFulfilled(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()

	script := `
		var result = null;
		Promise.any([
			Promise.reject("error 1"),
			Promise.resolve("success")
		]).then(function(val) {
			result = val;
		});
	`

	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	loop.Wait()

	result := vm.Get("result")
	if result.String() != "success" {
		t.Errorf("Expected 'success', got '%s'", result.String())
	}
}

func TestPromiseAnyAllRejected(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()

	script := `
		var error = null;
		Promise.any([
			Promise.reject("error 1"),
			Promise.reject("error 2"),
			Promise.reject("error 3")
		]).catch(function(err) {
			error = err;
		});
	`

	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	loop.Wait()

	result := vm.Get("error")
	if result.ExportType() == nil {
		t.Fatal("Error is nil")
	}

	errObj := result.ToObject(vm)
	if errObj.Get("name").String() != "AggregateError" {
		t.Errorf("Expected AggregateError, got '%s'", errObj.Get("name").String())
	}
}

func TestPromiseAnyEmpty(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()

	script := `
		var error = null;
		Promise.any([]).catch(function(err) {
			error = err;
		});
	`

	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	loop.Wait()

	result := vm.Get("error")
	errObj := result.ToObject(vm)
	if errObj.Get("name").String() != "AggregateError" {
		t.Errorf("Expected AggregateError for empty array, got '%s'", errObj.Get("name").String())
	}
}

func TestPromiseAnyWithNonPromises(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()

	script := `
		var result = null;
		Promise.any([
			Promise.reject("error"),
			42,
			Promise.resolve("too late")
		]).then(function(val) {
			result = val;
		});
	`

	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	loop.Wait()

	result := vm.Get("result")
	if result.ToInteger() != 42 {
		t.Errorf("Expected 42, got %d", result.ToInteger())
	}
}

func TestPromiseAnyMixedTiming(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()

	timers := NewTimers(loop)
	timerObj := timers.Export(vm).ToObject(vm)
	vm.Set("setTimeout", timerObj.Get("setTimeout"))

	script := `
		var result = null;
		Promise.any([
			new Promise(function(resolve, reject) {
				setTimeout(function() { reject("error"); }, 200);
			}),
			new Promise(function(resolve) {
				setTimeout(function() { resolve("fast success"); }, 100);
			}),
			new Promise(function(resolve) {
				setTimeout(function() { resolve("slow success"); }, 300);
			})
		]).then(function(val) {
			result = val;
		});
	`

	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	time.Sleep(400 * time.Millisecond)
	loop.Wait()

	result := vm.Get("result")
	if result.String() != "fast success" {
		t.Errorf("Expected 'fast success', got '%s'", result.String())
	}
}

func TestPromiseAnyIgnoresRejections(t *testing.T) {
	vm, loop := setupTestEnvironment()
	defer loop.Stop()

	script := `
		var result = null;
		Promise.any([
			Promise.reject("error 1"),
			Promise.reject("error 2"),
			Promise.resolve("the one success"),
			Promise.reject("error 3")
		]).then(function(val) {
			result = val;
		});
	`

	_, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	loop.Wait()

	result := vm.Get("result")
	if result.String() != "the one success" {
		t.Errorf("Expected 'the one success', got '%s'", result.String())
	}
}
