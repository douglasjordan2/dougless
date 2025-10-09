package event

import (
	"context"
	"sync"
	"time"
)

// Task represents a task to be executed in the event loop
type Task struct {
	ID       string
	Callback func()
	Delay    time.Duration
	Interval bool
}

// Loop represents the event loop
type Loop struct {
	tasks   chan *Task
	timers  map[string]*time.Timer
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	mu      sync.RWMutex
	running bool
}

// NewLoop creates a new event loop
func NewLoop() *Loop {
	ctx, cancel := context.WithCancel(context.Background())
	return &Loop{
		tasks:  make(chan *Task, 100),
		timers: make(map[string]*time.Timer),
		ctx:    ctx,
		cancel: cancel,
	}
}

// Run starts the event loop
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

// Stop stops the event loop
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

// Wait waits for all pending tasks to complete
func (l *Loop) Wait() {
	l.wg.Wait()
}

// ScheduleTask schedules a task to be executed
func (l *Loop) ScheduleTask(task *Task) {
  l.wg.Add(1) // track pending task when it's scheduled to account for delayed tasks (ex: setTimeout)
	if task.Delay > 0 {
		l.scheduleDelayedTask(task)
	} else {
		l.tasks <- task
	}
}

// scheduleDelayedTask schedules a task with a delay
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

// ClearTimer clears a scheduled timer
func (l *Loop) ClearTimer(id string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	if timer, exists := l.timers[id]; exists {
		timer.Stop()
		delete(l.timers, id)
    l.wg.Done()
	}
}

// executeTask executes a task
func (l *Loop) executeTask(task *Task) {
	go func() {
		defer l.wg.Done()
		task.Callback()
	}()
}
