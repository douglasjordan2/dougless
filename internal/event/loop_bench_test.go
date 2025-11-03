package event

import (
	"sync"
	"testing"
	"time"
)

// BenchmarkLoopCreation measures the cost of creating a new event loop
func BenchmarkLoopCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewLoop()
	}
}

// BenchmarkTaskScheduling measures immediate task scheduling overhead
func BenchmarkTaskScheduling(b *testing.B) {
	loop := NewLoop()
	go loop.Run()
	defer loop.Stop()
	
	var wg sync.WaitGroup
	task := &Task{
		Callback: func() { wg.Done() },
		Delay:    0,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		loop.ScheduleTask(task)
		wg.Wait()
	}
}

// BenchmarkDelayedTaskScheduling measures delayed task scheduling overhead
func BenchmarkDelayedTaskScheduling(b *testing.B) {
	loop := NewLoop()
	go loop.Run()
	defer loop.Stop()
	
	var wg sync.WaitGroup
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		task := &Task{
			ID:       "timer",
			Callback: func() { wg.Done() },
			Delay:    1 * time.Millisecond,
		}
		loop.ScheduleTask(task)
		wg.Wait()
	}
}

// BenchmarkParallelTaskScheduling measures concurrent task scheduling
func BenchmarkParallelTaskScheduling(b *testing.B) {
	loop := NewLoop()
	go loop.Run()
	defer loop.Stop()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var wg sync.WaitGroup
			wg.Add(1)
			task := &Task{
				Callback: func() { wg.Done() },
				Delay:    0,
			}
			loop.ScheduleTask(task)
			wg.Wait()
		}
	})
}

// BenchmarkTimerCancellation measures clearTimeout overhead
func BenchmarkTimerCancellation(b *testing.B) {
	loop := NewLoop()
	go loop.Run()
	defer loop.Stop()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		timerID := "timer-cancel-test"
		task := &Task{
			ID:       timerID,
			Callback: func() {},
			Delay:    100 * time.Millisecond,
		}
		loop.ScheduleTask(task)
		loop.ClearTimer(timerID)
	}
}

// BenchmarkMultipleTimers measures performance with many concurrent timers
func BenchmarkMultipleTimers(b *testing.B) {
	counts := []int{10, 100, 1000}
	
	for _, count := range counts {
		b.Run(string(rune(count))+"_timers", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				loop := NewLoop()
				go loop.Run()
				var wg sync.WaitGroup
				
				b.StartTimer()
				for j := 0; j < count; j++ {
					wg.Add(1)
					task := &Task{
						Callback: func() { wg.Done() },
						Delay:    1 * time.Millisecond,
					}
					loop.ScheduleTask(task)
				}
				wg.Wait()
				b.StopTimer()
				loop.Stop()
			}
		})
	}
}

// BenchmarkKeepAlive measures KeepAlive/Done overhead
func BenchmarkKeepAlive(b *testing.B) {
	loop := NewLoop()
	go loop.Run()
	defer loop.Stop()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		done := loop.KeepAlive()
		done()
	}
}

// BenchmarkTaskThroughput measures maximum task processing rate
func BenchmarkTaskThroughput(b *testing.B) {
	loop := NewLoop()
	go loop.Run()
	defer loop.Stop()
	
	var counter int64
	task := &Task{
		Callback: func() { counter++ },
		Delay:    0,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loop.ScheduleTask(task)
	}
	loop.Wait()
	
	b.ReportMetric(float64(counter)/b.Elapsed().Seconds(), "tasks/sec")
}
