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
	t.Skip("TODO: Promise chaining with nested Promise.reject() needs work")
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
	t.Skip("TODO: Error propagation through .then() without handler needs work")
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
