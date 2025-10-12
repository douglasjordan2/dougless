package modules

import (
	"sync"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/douglasjordan2/dougless/internal/event"
)

func TestTimersSetTimeout(t *testing.T) {
	vm := goja.New()
	eventLoop := event.NewLoop()
	timers := NewTimers(eventLoop)

	// Set timers as global functions
	timerObj := timers.Export(vm).ToObject(vm)
	vm.Set("setTimeout", timerObj.Get("setTimeout"))
	vm.Set("clearTimeout", timerObj.Get("clearTimeout"))

	go eventLoop.Run()
	defer eventLoop.Stop()

	t.Run("basic timeout execution", func(t *testing.T) {
		executed := false
		vm.Set("callback", func() {
			executed = true
		})

		_, err := vm.RunString(`setTimeout(callback, 50)`)
		if err != nil {
			t.Fatalf("Failed to set timeout: %v", err)
		}

		time.Sleep(100 * time.Millisecond)
		eventLoop.Wait()

		if !executed {
			t.Error("Timeout callback was not executed")
		}
	})

	t.Run("timeout with zero delay", func(t *testing.T) {
		executed := false
		vm.Set("callback", func() {
			executed = true
		})

		_, err := vm.RunString(`setTimeout(callback, 0)`)
		if err != nil {
			t.Fatalf("Failed to set timeout: %v", err)
		}

		time.Sleep(50 * time.Millisecond)
		eventLoop.Wait()

		if !executed {
			t.Error("Zero-delay timeout was not executed")
		}
	})

	t.Run("timeout without delay argument", func(t *testing.T) {
		executed := false
		vm.Set("callback", func() {
			executed = true
		})

		_, err := vm.RunString(`setTimeout(callback)`)
		if err != nil {
			t.Fatalf("Failed to set timeout: %v", err)
		}

		time.Sleep(50 * time.Millisecond)
		eventLoop.Wait()

		if !executed {
			t.Error("No-delay timeout was not executed")
		}
	})

	t.Run("returns timer ID", func(t *testing.T) {
		vm.Set("callback", func() {})

		val, err := vm.RunString(`setTimeout(callback, 100)`)
		if err != nil {
			t.Fatalf("Failed to set timeout: %v", err)
		}

		timerID := val.String()
		if timerID == "" {
			t.Error("setTimeout should return a timer ID")
		}
	})
}

func TestTimersSetInterval(t *testing.T) {
	vm := goja.New()
	eventLoop := event.NewLoop()
	timers := NewTimers(eventLoop)

	timerObj := timers.Export(vm).ToObject(vm)
	vm.Set("setInterval", timerObj.Get("setInterval"))
	vm.Set("clearInterval", timerObj.Get("clearInterval"))

	go eventLoop.Run()
	defer eventLoop.Stop()

	t.Run("interval execution", func(t *testing.T) {
		counter := 0
		var mu sync.Mutex

		vm.Set("callback", func() {
			mu.Lock()
			counter++
			mu.Unlock()
		})

		val, err := vm.RunString(`setInterval(callback, 50)`)
		if err != nil {
			t.Fatalf("Failed to set interval: %v", err)
		}

		// Let it run for ~250ms (should execute ~5 times)
		time.Sleep(250 * time.Millisecond)

		// Clear the interval
		timerID := val.String()
		vm.Set("timerID", timerID)
		_, err = vm.RunString(`clearInterval(timerID)`)
		if err != nil {
			t.Fatalf("Failed to clear interval: %v", err)
		}

		time.Sleep(100 * time.Millisecond)
		eventLoop.Wait()

		mu.Lock()
		finalCount := counter
		mu.Unlock()

		if finalCount < 3 {
			t.Errorf("Expected at least 3 interval executions, got %d", finalCount)
		}
		if finalCount > 7 {
			t.Errorf("Expected at most 7 interval executions, got %d", finalCount)
		}
	})

	t.Run("returns timer ID", func(t *testing.T) {
		vm.Set("callback", func() {})

		val, err := vm.RunString(`setInterval(callback, 100)`)
		if err != nil {
			t.Fatalf("Failed to set interval: %v", err)
		}

		timerID := val.String()
		if timerID == "" {
			t.Error("setInterval should return a timer ID")
		}

		// Clean up
		vm.Set("timerID", timerID)
		vm.RunString(`clearInterval(timerID)`)
	})
}

func TestTimersClearTimeout(t *testing.T) {
	vm := goja.New()
	eventLoop := event.NewLoop()
	timers := NewTimers(eventLoop)

	timerObj := timers.Export(vm).ToObject(vm)
	vm.Set("setTimeout", timerObj.Get("setTimeout"))
	vm.Set("clearTimeout", timerObj.Get("clearTimeout"))

	go eventLoop.Run()
	defer eventLoop.Stop()

	t.Run("cancel before execution", func(t *testing.T) {
		executed := false
		vm.Set("callback", func() {
			executed = true
		})

		val, err := vm.RunString(`setTimeout(callback, 100)`)
		if err != nil {
			t.Fatalf("Failed to set timeout: %v", err)
		}

		// Cancel immediately
		timerID := val.String()
		vm.Set("timerID", timerID)
		_, err = vm.RunString(`clearTimeout(timerID)`)
		if err != nil {
			t.Fatalf("Failed to clear timeout: %v", err)
		}

		time.Sleep(150 * time.Millisecond)
		eventLoop.Wait()

		if executed {
			t.Error("Cancelled timeout should not have executed")
		}
	})

	t.Run("clear non-existent timer", func(t *testing.T) {
		// Should not panic
		_, err := vm.RunString(`clearTimeout("fake-id")`)
		if err != nil {
			t.Fatalf("Clearing non-existent timer should not error: %v", err)
		}
	})

	t.Run("clear without argument", func(t *testing.T) {
		// Should not panic
		_, err := vm.RunString(`clearTimeout()`)
		if err != nil {
			t.Fatalf("Clearing without argument should not error: %v", err)
		}
	})

	t.Run("double clear", func(t *testing.T) {
		vm.Set("callback", func() {})

		val, err := vm.RunString(`setTimeout(callback, 100)`)
		if err != nil {
			t.Fatalf("Failed to set timeout: %v", err)
		}

		timerID := val.String()
		vm.Set("timerID", timerID)

		// Clear once
		_, err = vm.RunString(`clearTimeout(timerID)`)
		if err != nil {
			t.Fatalf("Failed first clear: %v", err)
		}

		// Clear again - should not panic
		_, err = vm.RunString(`clearTimeout(timerID)`)
		if err != nil {
			t.Fatalf("Second clear should not error: %v", err)
		}
	})
}

func TestTimersClearInterval(t *testing.T) {
	vm := goja.New()
	eventLoop := event.NewLoop()
	timers := NewTimers(eventLoop)

	timerObj := timers.Export(vm).ToObject(vm)
	vm.Set("setInterval", timerObj.Get("setInterval"))
	vm.Set("clearInterval", timerObj.Get("clearInterval"))

	go eventLoop.Run()
	defer eventLoop.Stop()

	t.Run("stop recurring interval", func(t *testing.T) {
		counter := 0
		var mu sync.Mutex

		vm.Set("callback", func() {
			mu.Lock()
			counter++
			mu.Unlock()
		})

		val, err := vm.RunString(`setInterval(callback, 50)`)
		if err != nil {
			t.Fatalf("Failed to set interval: %v", err)
		}

		// Let it run a bit
		time.Sleep(150 * time.Millisecond)

		// Clear the interval
		timerID := val.String()
		vm.Set("timerID", timerID)
		_, err = vm.RunString(`clearInterval(timerID)`)
		if err != nil {
			t.Fatalf("Failed to clear interval: %v", err)
		}

		mu.Lock()
		countAfterClear := counter
		mu.Unlock()

		// Wait more time
		time.Sleep(150 * time.Millisecond)
		eventLoop.Wait()

		mu.Lock()
		finalCount := counter
		mu.Unlock()

		// Counter should not have increased after clearing
		if finalCount != countAfterClear {
			t.Errorf("Counter increased after clearing interval: before=%d, after=%d", countAfterClear, finalCount)
		}
	})
}

func TestTimersMultipleConcurrent(t *testing.T) {
	vm := goja.New()
	eventLoop := event.NewLoop()
	timers := NewTimers(eventLoop)

	timerObj := timers.Export(vm).ToObject(vm)
	vm.Set("setTimeout", timerObj.Get("setTimeout"))

	go eventLoop.Run()
	defer eventLoop.Stop()

	t.Run("multiple timeouts execute independently", func(t *testing.T) {
		count1 := 0
		count2 := 0
		count3 := 0

		vm.Set("callback1", func() { count1++ })
		vm.Set("callback2", func() { count2++ })
		vm.Set("callback3", func() { count3++ })

		_, err := vm.RunString(`
			setTimeout(callback1, 50);
			setTimeout(callback2, 100);
			setTimeout(callback3, 150);
		`)
		if err != nil {
			t.Fatalf("Failed to set multiple timeouts: %v", err)
		}

		time.Sleep(200 * time.Millisecond)
		eventLoop.Wait()

		if count1 != 1 || count2 != 1 || count3 != 1 {
			t.Errorf("Expected all callbacks to execute once, got: %d, %d, %d", count1, count2, count3)
		}
	})
}

func TestTimersErrorHandling(t *testing.T) {
	vm := goja.New()
	eventLoop := event.NewLoop()
	timers := NewTimers(eventLoop)

	timerObj := timers.Export(vm).ToObject(vm)
	vm.Set("setTimeout", timerObj.Get("setTimeout"))

	go eventLoop.Run()
	defer eventLoop.Stop()

	t.Run("non-function callback", func(t *testing.T) {
		_, err := vm.RunString(`setTimeout("not a function", 100)`)
		if err == nil {
			t.Error("Expected error when passing non-function callback")
		}
	})

	t.Run("missing callback", func(t *testing.T) {
		_, err := vm.RunString(`setTimeout()`)
		if err == nil {
			t.Error("Expected error when missing callback")
		}
	})
}

func TestTimersIntegration(t *testing.T) {
	vm := goja.New()
	eventLoop := event.NewLoop()
	timers := NewTimers(eventLoop)

	timerObj := timers.Export(vm).ToObject(vm)
	vm.Set("setTimeout", timerObj.Get("setTimeout"))
	vm.Set("setInterval", timerObj.Get("setInterval"))
	vm.Set("clearTimeout", timerObj.Get("clearTimeout"))
	vm.Set("clearInterval", timerObj.Get("clearInterval"))

	go eventLoop.Run()
	defer eventLoop.Stop()

	t.Run("nested timers", func(t *testing.T) {
		results := make([]string, 0)
		var mu sync.Mutex

		vm.Set("addResult", func(msg string) {
			mu.Lock()
			results = append(results, msg)
			mu.Unlock()
		})

		_, err := vm.RunString(`
			setTimeout(function() {
				addResult("outer");
				setTimeout(function() {
					addResult("inner");
				}, 50);
			}, 50);
		`)
		if err != nil {
			t.Fatalf("Failed to execute nested timers: %v", err)
		}

		time.Sleep(200 * time.Millisecond)
		eventLoop.Wait()

		mu.Lock()
		defer mu.Unlock()

		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(results))
		}
		if len(results) >= 2 {
			if results[0] != "outer" {
				t.Errorf("Expected first result to be 'outer', got %q", results[0])
			}
			if results[1] != "inner" {
				t.Errorf("Expected second result to be 'inner', got %q", results[1])
			}
		}
	})

	t.Run("timeout and interval together", func(t *testing.T) {
		timeoutExecuted := false
		intervalCount := 0
		var mu sync.Mutex

		vm.Set("timeoutCallback", func() {
			mu.Lock()
			timeoutExecuted = true
			mu.Unlock()
		})

		vm.Set("intervalCallback", func() {
			mu.Lock()
			intervalCount++
			mu.Unlock()
		})

		val, err := vm.RunString(`
			setTimeout(timeoutCallback, 100);
			setInterval(intervalCallback, 50);
		`)
		if err != nil {
			t.Fatalf("Failed to execute: %v", err)
		}

		time.Sleep(200 * time.Millisecond)

		// Clear interval
		vm.Set("intervalID", val.String())
		vm.RunString(`clearInterval(intervalID)`)

		time.Sleep(50 * time.Millisecond)
		eventLoop.Wait()

		mu.Lock()
		defer mu.Unlock()

		if !timeoutExecuted {
			t.Error("Timeout should have executed")
		}
		if intervalCount < 2 {
			t.Errorf("Expected at least 2 interval executions, got %d", intervalCount)
		}
	})
}
