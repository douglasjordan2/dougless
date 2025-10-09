package event

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestNewLoop verifies that a new event loop is created with proper initialization
func TestNewLoop(t *testing.T) {
	// The simplest test - just create a loop and check it's not nil
	loop := NewLoop()
	
	if loop == nil {
		t.Fatal("NewLoop() returned nil")
	}
	
	// Check that internal structures are initialized
	if loop.tasks == nil {
		t.Error("tasks channel should be initialized")
	}
	
	if loop.timers == nil {
		t.Error("timers map should be initialized")
	}
	
	if loop.ctx == nil {
		t.Error("context should be initialized")
	}
}

// TestScheduleTask tests that tasks can be scheduled and executed
func TestScheduleTask(t *testing.T) {
	// t.Run() creates a subtest - great for organizing related tests
	t.Run("immediate task execution", func(t *testing.T) {
		loop := NewLoop()
		go loop.Run()
		defer loop.Stop()
		
		// Use a channel to track if the task executed
		executed := make(chan bool, 1)
		
		task := &Task{
			ID: "test-1",
			Callback: func() {
				executed <- true
			},
			Delay: 0,
		}
		
		loop.ScheduleTask(task)
		
		// Wait for task execution with timeout to prevent hanging tests
		select {
		case <-executed:
			// Success! Task was executed
		case <-time.After(1 * time.Second):
			t.Fatal("task did not execute within timeout")
		}
		
		loop.Wait()
	})
	
	t.Run("delayed task execution", func(t *testing.T) {
		loop := NewLoop()
		go loop.Run()
		defer loop.Stop()
		
		executed := make(chan bool, 1)
		startTime := time.Now()
		delay := 100 * time.Millisecond
		
		task := &Task{
			ID: "test-delayed",
			Callback: func() {
				executed <- true
			},
			Delay: delay,
		}
		
		loop.ScheduleTask(task)
		
		select {
		case <-executed:
			elapsed := time.Since(startTime)
			// Verify the delay was respected (with some tolerance)
			if elapsed < delay {
				t.Errorf("task executed too early: %v < %v", elapsed, delay)
			}
			// Allow up to 50ms tolerance for timing
			if elapsed > delay+50*time.Millisecond {
				t.Errorf("task executed too late: %v > %v", elapsed, delay+50*time.Millisecond)
			}
		case <-time.After(1 * time.Second):
			t.Fatal("delayed task did not execute within timeout")
		}
		
		loop.Wait()
	})
	
	t.Run("multiple tasks execution", func(t *testing.T) {
		loop := NewLoop()
		go loop.Run()
		defer loop.Stop()
		
		// Use WaitGroup to track multiple tasks
		var wg sync.WaitGroup
		taskCount := 10
		wg.Add(taskCount)
		
		for i := 0; i < taskCount; i++ {
			task := &Task{
				ID: "test-multi-" + string(rune(i)),
				Callback: func() {
					wg.Done()
				},
				Delay: 0,
			}
			loop.ScheduleTask(task)
		}
		
		// Wait for all tasks with timeout
		done := make(chan bool)
		go func() {
			wg.Wait()
			done <- true
		}()
		
		select {
		case <-done:
			// All tasks completed successfully
		case <-time.After(2 * time.Second):
			t.Fatal("not all tasks completed within timeout")
		}
		
		loop.Wait()
	})
}

// TestIntervalTask tests recurring interval tasks
func TestIntervalTask(t *testing.T) {
	loop := NewLoop()
	go loop.Run()
	defer loop.Stop()
	
	executionCount := 0
	var mu sync.Mutex
	
	task := &Task{
		ID: "interval-test",
		Callback: func() {
			mu.Lock()
			executionCount++
			mu.Unlock()
		},
		Delay:    50 * time.Millisecond,
		Interval: true,
	}
	
	loop.ScheduleTask(task)
	
	// Let it run for enough time to execute multiple times
	time.Sleep(200 * time.Millisecond)
	
	// Stop the interval
	loop.ClearTimer("interval-test")
	
	// Give it a moment to process the stop
	time.Sleep(50 * time.Millisecond)
	
	mu.Lock()
	count := executionCount
	mu.Unlock()
	
	// Should have executed 3-4 times (200ms / 50ms with some tolerance)
	if count < 2 {
		t.Errorf("interval executed too few times: %d", count)
	}
	if count > 5 {
		t.Errorf("interval executed too many times: %d", count)
	}
	
	// Verify it stopped - wait and check count doesn't increase
	time.Sleep(100 * time.Millisecond)
	mu.Lock()
	finalCount := executionCount
	mu.Unlock()
	
	if finalCount != count {
		t.Errorf("interval continued after clear: %d != %d", finalCount, count)
	}
}

// TestClearTimer tests timer cancellation
func TestClearTimer(t *testing.T) {
	t.Run("clear timeout before execution", func(t *testing.T) {
		loop := NewLoop()
		go loop.Run()
		defer loop.Stop()
		
		executed := false
		var mu sync.Mutex
		
		task := &Task{
			ID: "clear-test",
			Callback: func() {
				mu.Lock()
				executed = true
				mu.Unlock()
			},
			Delay: 100 * time.Millisecond,
		}
		
		loop.ScheduleTask(task)
		
		// Clear the timer before it executes
		time.Sleep(20 * time.Millisecond)
		loop.ClearTimer("clear-test")
		
		// Wait longer than the original delay
		time.Sleep(150 * time.Millisecond)
		
		mu.Lock()
		defer mu.Unlock()
		
		if executed {
			t.Error("task should not have executed after being cleared")
		}
	})
	
	t.Run("clear non-existent timer", func(t *testing.T) {
		loop := NewLoop()
		go loop.Run()
		defer loop.Stop()
		
		// This should not panic or cause issues
		loop.ClearTimer("non-existent-timer")
		
		// If we get here without panic, test passes
	})
	
	t.Run("clear already executed timer", func(t *testing.T) {
		loop := NewLoop()
		go loop.Run()
		defer loop.Stop()
		
		done := make(chan bool)
		
		task := &Task{
			ID: "already-executed",
			Callback: func() {
				done <- true
			},
			Delay: 10 * time.Millisecond,
		}
		
		loop.ScheduleTask(task)
		
		// Wait for execution
		<-done
		
		// Try to clear after execution - should not cause issues
		loop.ClearTimer("already-executed")
	})
}

// TestStopLoop tests graceful shutdown
func TestStopLoop(t *testing.T) {
	t.Run("stop before any tasks", func(t *testing.T) {
		loop := NewLoop()
		go loop.Run()
		
		// Stop immediately
		loop.Stop()
		
		// Should be safe to call Stop multiple times
		loop.Stop()
		loop.Stop()
	})
	
	t.Run("stop with pending tasks", func(t *testing.T) {
		loop := NewLoop()
		go loop.Run()
		
		// Schedule some delayed tasks
		for i := 0; i < 5; i++ {
			task := &Task{
				ID: fmt.Sprintf("pending-%d", i),
				Callback: func() {
					time.Sleep(10 * time.Millisecond)
				},
				Delay: 100 * time.Millisecond,
			}
			loop.ScheduleTask(task)
		}
		
		// Check timer count before stop
		loop.mu.Lock()
		beforeCount := len(loop.timers)
		loop.mu.Unlock()
		t.Logf("Timers before stop: %d", beforeCount)
		
		// Give Run() goroutine time to start
		time.Sleep(10 * time.Millisecond)
		
		// Stop before tasks execute
		loop.Stop()
		
		// Give the Stop() method time to complete
		time.Sleep(10 * time.Millisecond)
		
		// Verify timers are cleaned up
		loop.mu.Lock()
		timerCount := len(loop.timers)
		loop.mu.Unlock()
		
		if timerCount != 0 {
			t.Errorf("timers not cleaned up after stop: %d remaining", timerCount)
		}
	})
}

// TestConcurrentOperations tests thread safety
func TestConcurrentOperations(t *testing.T) {
	loop := NewLoop()
	go loop.Run()
	defer loop.Stop()
	
	// Launch multiple goroutines scheduling tasks concurrently
	var wg sync.WaitGroup
	concurrency := 50
	
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			task := &Task{
				ID: "concurrent-" + string(rune(id)),
				Callback: func() {
					// Simulate work
					time.Sleep(1 * time.Millisecond)
				},
				Delay: time.Duration(id%10) * time.Millisecond,
			}
			
			loop.ScheduleTask(task)
		}(i)
	}
	
	wg.Wait()
	loop.Wait()
	
	// If we get here without race conditions or deadlocks, test passes
}

// BenchmarkScheduleTask measures the performance of task scheduling
func BenchmarkScheduleTask(b *testing.B) {
	loop := NewLoop()
	go loop.Run()
	defer loop.Stop()
	
	// Reset timer to exclude setup time
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		done := make(chan bool, 1)
		task := &Task{
			ID: "bench-task",
			Callback: func() {
				done <- true
			},
			Delay: 0,
		}
		
		loop.ScheduleTask(task)
		<-done
	}
}

// BenchmarkTimerScheduling measures delayed task scheduling performance
func BenchmarkTimerScheduling(b *testing.B) {
	loop := NewLoop()
	go loop.Run()
	defer loop.Stop()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		task := &Task{
			ID:    "bench-delayed",
			Callback: func() {},
			Delay: 10 * time.Millisecond,
		}
		
		loop.ScheduleTask(task)
	}
	
	loop.Wait()
}
