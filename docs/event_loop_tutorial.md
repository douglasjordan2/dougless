# Building an Event Loop from Scratch

## Introduction: What is an Event Loop?

The event loop is the heart of JavaScript's asynchronous execution model. It's what allows JavaScript to be non-blocking despite being single-threaded. Think of it as a waiter in a restaurant:
- Takes orders (tasks) from customers (your code)
- Puts them in a queue
- Processes them one at a time
- Can handle "delayed" orders (timers)
- Keeps working until all orders are complete

## Part 1: Understanding the Problem

JavaScript appears to do multiple things at once:
```javascript
console.log("First");
setTimeout(() => console.log("Second"), 0);
console.log("Third");
// Output: First, Third, Second
```

How does "Second" print last even with 0 delay? The event loop!

## Part 2: Core Concepts

### 2.1 The Task Queue
A queue of functions waiting to be executed. Think FIFO (First In, First Out).

### 2.2 The Call Stack
Where JavaScript actually executes code. The event loop only adds tasks when this is empty.

### 2.3 Timers
Tasks that should run after a delay. These need special handling.

### 2.4 Microtasks vs Macrotasks
For now, we'll focus on macrotasks (setTimeout, setInterval). Microtasks (Promises) come later.

## Part 3: Building Blocks in Go

### 3.1 Basic Task Structure
First, think about what a "task" needs:
```go
type Task struct {
    ID       string        // Unique identifier
    Callback func()        // The function to run
    Delay    time.Duration // How long to wait (0 for immediate)
}
```

### 3.2 Channels: Go's Secret Weapon
Go channels are perfect for queues:
```go
taskQueue := make(chan *Task, 100) // Buffered channel holds up to 100 tasks
```

### 3.3 The Main Loop Pattern
```go
for {
    select {
    case task := <-taskQueue:
        // Execute the task
        task.Callback()
    case <-stopSignal:
        // Shutdown gracefully
        return
    }
}
```

## Part 4: Step-by-Step Implementation Guide

### Step 1: Start Simple - Immediate Tasks Only

Create a minimal event loop that can:
1. Accept tasks with no delay
2. Execute them in order
3. Shut down cleanly

Questions to consider:
- How will you start the loop?
- How will you add tasks?
- How will you stop it?

Try this first before moving on!

### Step 2: Add Timer Support

Now the interesting part! For delayed tasks:
1. You can't just sleep - that blocks everything
2. You need to track when each timer should fire
3. Multiple timers might be running at once

Hint: `time.AfterFunc` is your friend in Go:
```go
timer := time.AfterFunc(delay, func() {
    // This runs after 'delay'
    // But in a different goroutine!
})
```

Challenge: How do you get the callback back to your main loop?

### Step 3: Cancellable Timers

`clearTimeout` needs to cancel timers:
1. Keep a map of timer IDs to timer objects
2. Make sure to clean up after timers fire
3. Handle thread-safety (multiple goroutines!)

Consider: What happens if you try to cancel an already-fired timer?

### Step 4: Wait for Completion

Your runtime needs to know when all work is done:
1. Track how many tasks are pending
2. Use sync.WaitGroup or similar
3. Important: Timers count as pending work!

### Step 5: Thread Safety & JavaScript Context

This is the tricky part! Remember:
1. Goja (your JS engine) is NOT thread-safe
2. Timer callbacks fire in Go goroutines
3. You must execute JS code on the main thread

Solution approach:
- Timer fires in goroutine → sends task to channel → main loop executes it

## Part 5: Testing Your Implementation

### Test 1: Basic Execution
```javascript
console.log("1");
setTimeout(() => console.log("2"), 100);
console.log("3");
// Should print: 1, 3, (wait 100ms), 2
```

### Test 2: Multiple Timers
```javascript
setTimeout(() => console.log("A"), 300);
setTimeout(() => console.log("B"), 100);
setTimeout(() => console.log("C"), 200);
// Should print: B, C, A (in that order)
```

### Test 3: Zero Delay
```javascript
setTimeout(() => console.log("Async"), 0);
console.log("Sync");
// Should print: Sync, Async
```

### Test 4: Clear Timer
```javascript
const id = setTimeout(() => console.log("Never"), 1000);
clearTimeout(id);
// Should print nothing
```

## Part 6: Common Pitfalls

1. **Goroutine Leaks**: Forgetting to stop timers on shutdown
2. **Race Conditions**: Multiple goroutines accessing timer map
3. **Busy Waiting**: Using a tight loop instead of channels
4. **Blocking the Loop**: Doing heavy work in callbacks
5. **Context Loss**: JavaScript 'this' binding issues

## Part 7: Advanced Concepts (After Basic Implementation)

Once your basic loop works, consider:
1. **setInterval**: Repeating timers
2. **Task Priorities**: Some tasks might be more important
3. **Error Handling**: What if a callback panics?
4. **Performance**: Minimize allocations, optimize timer management
5. **Debugging**: How can you inspect what's queued?

## Implementation Exercise Plan

### Day 1: Build Basic Loop
- [ ] Create Task struct
- [ ] Implement basic loop with channel
- [ ] Add Start(), Stop(), and ScheduleTask() methods
- [ ] Test with immediate tasks only

### Day 2: Add Timer Support
- [ ] Implement delayed task scheduling
- [ ] Add timer tracking map
- [ ] Handle timer cleanup
- [ ] Test multiple timers

### Day 3: Integration
- [ ] Connect to your Runtime
- [ ] Implement setTimeout in JavaScript
- [ ] Add clearTimeout support
- [ ] Test with real JavaScript code

### Day 4: Polish
- [ ] Add proper error handling
- [ ] Implement setInterval
- [ ] Add debugging capabilities
- [ ] Performance optimization

## Debugging Tips

Add logging to understand flow:
```go
log.Printf("[EventLoop] Task scheduled: %s, delay: %v", task.ID, task.Delay)
log.Printf("[EventLoop] Executing task: %s", task.ID)
log.Printf("[EventLoop] Timer fired: %s", task.ID)
```

## Key Questions to Think About

1. Why does JavaScript need an event loop instead of just using threads?
2. What happens if a callback schedules another timer?
3. How would you handle recursive setTimeout calls?
4. What's the minimum useful delay? (Hint: system timer resolution)
5. How would you implement setImmediate vs setTimeout(fn, 0)?

## Resources

- [Node.js Event Loop Docs](https://nodejs.org/en/docs/guides/event-loop-timers-and-nexttick/)
- [MDN Event Loop](https://developer.mozilla.org/en-US/docs/Web/JavaScript/EventLoop)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Channels in Go](https://gobyexample.com/channels)

## Your First Goal

Start by building a loop that can:
```go
loop := NewEventLoop()
loop.Start()
loop.ScheduleTask(&Task{
    Callback: func() { fmt.Println("Hello from event loop!") },
})
loop.Stop()
```

Once this works, you truly understand event loops!

---

Remember: The best way to learn is to build it wrong first, understand why it's wrong, then fix it. Don't be afraid to experiment!
