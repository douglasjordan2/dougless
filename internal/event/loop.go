// Package event implements a non-blocking event loop for asynchronous operations.
//
// The event loop enables async JavaScript features like setTimeout, setInterval,
// and Promise resolution. It uses Go's concurrency primitives (channels, goroutines)
// to provide non-blocking I/O and task scheduling.
//
// Key features:
//   - Task queue with buffered channel (100 capacity)
//   - Timer management using Go's time.AfterFunc
//   - Graceful shutdown with WaitGroup synchronization
//   - Context-based cancellation support
//
// Example usage:
//
//	loop := event.NewLoop()
//	go loop.Run()
//	defer loop.Stop()
//
//	loop.ScheduleTask(&event.Task{
//	    Callback: func() { fmt.Println("Hello") },
//	    Delay: 100 * time.Millisecond,
//	})
//
//	loop.Wait()  // Wait for all tasks to complete
package event

import (
	"context"
	"sync"
	"time"
)

// Task represents a unit of work to be executed in the event loop.
type Task struct {
	ID       string        // Unique identifier for the task (used for timer cancellation)
	Callback func()        // Function to execute
	Delay    time.Duration // Delay before execution (0 for immediate)
	Interval bool          // If true, task repeats at Delay intervals
}

// Loop represents the event loop that processes async tasks.
// It manages a queue of tasks and scheduled timers, executing them
// sequentially to maintain FIFO order for immediate tasks.
type Loop struct {
	tasks   chan *Task             // Buffered channel for task queue (capacity: 100)
	timers  map[string]*time.Timer // Map of timer IDs to Go timers
	ctx     context.Context        // Context for cancellation
	cancel  context.CancelFunc     // Function to cancel the context
	wg      sync.WaitGroup         // Tracks pending tasks for graceful shutdown
	mu      sync.RWMutex           // Protects timers map and running flag
	execMu  sync.Mutex             // Ensures sequential execution of tasks
	running bool                   // Indicates if the loop is currently running
}

// NewLoop creates and initializes a new event loop.
// The loop must be started with Run() before it can process tasks.
//
// Example:
//
//	loop := event.NewLoop()
//	go loop.Run()
func NewLoop() *Loop {
	ctx, cancel := context.WithCancel(context.Background())
	return &Loop{
		tasks:  make(chan *Task, 100),
		timers: make(map[string]*time.Timer),
		ctx:    ctx,
		cancel: cancel,
	}
}

// NewLoopWithContext creates a new event loop with a custom parent context.
// The loop will be cancelled when the parent context is cancelled.
//
// This is useful for implementing timeout or cancellation policies:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	loop := event.NewLoopWithContext(ctx)
func NewLoopWithContext(ctx context.Context) *Loop {
	loopCtx, cancel := context.WithCancel(ctx)
	return &Loop{
		tasks:  make(chan *Task, 100),
		timers: make(map[string]*time.Timer),
		ctx:    loopCtx,
		cancel: cancel,
	}
}

// Run starts the event loop and begins processing tasks.
// This method blocks until Stop() is called or the context is cancelled.
//
// Run should typically be called in a separate goroutine:
//
//	go loop.Run()
//
// The loop is safe to call multiple times; subsequent calls are no-ops.
func (l *Loop) Run() {
	l.mu.Lock()
	if l.running {
		l.mu.Unlock()
		return
	}
	l.running = true
	l.mu.Unlock()

	for {
		select {
		case <-l.ctx.Done():
			return
		case task := <-l.tasks:
			if task != nil {
				l.executeTask(task)
			}
		}
	}
}

// Stop gracefully shuts down the event loop.
//
// It:
//  1. Cancels the loop context
//  2. Stops all pending timers
//  3. Marks tasks as complete to unblock Wait()
//
// After Stop() is called, no new tasks will be processed.
func (l *Loop) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.running {
		return
	}

	l.cancel()
	l.running = false

	// Clean up timers and decrement WaitGroup for each pending timer
	for _, timer := range l.timers {
		timer.Stop()
		l.wg.Done() // Decrement for the pending task that won't execute
	}
	l.timers = make(map[string]*time.Timer)
}

// Wait blocks until all pending tasks have completed.
// This is typically called after the main script execution finishes
// to ensure all async operations (timers, promises) complete.
//
// Example:
//
//	loop.ScheduleTask(task)
//	loop.Wait()  // Blocks until task completes
func (l *Loop) Wait() {
	l.wg.Wait()
}

// KeepAlive increments the task counter and returns a function to decrement it.
// This is useful for keeping the event loop alive during async operations
// that don't go through ScheduleTask (e.g., HTTP servers).
//
// Example:
//
//	done := loop.KeepAlive()
//	defer done()
//	// ... perform async work ...
func (l *Loop) KeepAlive() func() {
	l.wg.Add(1)
	return func() {
		l.wg.Done()
	}
}

// ScheduleTask queues a task for execution on the event loop.
//
// Tasks with Delay > 0 are scheduled using Go timers and executed after the delay.
// Tasks with Delay == 0 are executed as soon as the loop processes them.
//
// All tasks are tracked in the WaitGroup until they complete.
//
// Example:
//
//	loop.ScheduleTask(&event.Task{
//	    Callback: func() { fmt.Println("Delayed") },
//	    Delay: 1 * time.Second,
//	})
func (l *Loop) ScheduleTask(task *Task) {
	l.wg.Add(1) // track pending task when it's scheduled to account for delayed tasks (ex: setTimeout)
	if task.Delay > 0 {
		l.scheduleDelayedTask(task)
	} else {
		l.tasks <- task
	}
}

// scheduleDelayedTask schedules a task to execute after a delay using Go's time.AfterFunc.
// For interval tasks (task.Interval == true), the task is automatically rescheduled after execution.
func (l *Loop) scheduleDelayedTask(task *Task) {
	timer := time.AfterFunc(task.Delay, func() {
		l.tasks <- task

		if task.Interval {
			// add to waitgroup and reschedule for intervals
			l.wg.Add(1)
			l.scheduleDelayedTask(task)
		} else {
			// Remove from timers map for one-time tasks
			l.mu.Lock()
			delete(l.timers, task.ID)
			l.mu.Unlock()
		}
	})

	l.mu.Lock()
	l.timers[task.ID] = timer
	l.mu.Unlock()
}

// ClearTimer cancels a scheduled timer by its ID.
// This is used to implement clearTimeout() and clearInterval().
//
// If the timer doesn't exist, this is a no-op.
func (l *Loop) ClearTimer(id string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if timer, exists := l.timers[id]; exists {
		timer.Stop()
		delete(l.timers, id)
		l.wg.Done()
	}
}

// executeTask executes a task's callback in a goroutine and waits for completion.
// This ensures FIFO ordering: tasks are executed in the exact order they're dequeued,
// and the event loop doesn't process the next task until the current one finishes.
//
// This prevents Goja VM reentrancy issues while maintaining deterministic ordering.
func (l *Loop) executeTask(task *Task) {
	done := make(chan struct{})
	go func() {
		defer l.wg.Done()
		defer close(done)
		task.Callback()
	}()
	// Wait for task to complete before returning to Run() loop
	<-done
}
