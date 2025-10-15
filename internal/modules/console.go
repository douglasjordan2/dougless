package modules

import (
	"fmt"
	"sync"
	"time"

	"github.com/dop251/goja"
)

// Console provides debugging and logging functionality for JavaScript.
// It implements the standard console API with methods for output, warnings,
// errors, performance timing, and table formatting.
//
// Available globally in JavaScript as the 'console' object.
type Console struct {
	vm       *goja.Runtime        // JavaScript runtime instance
	timers   map[string]time.Time // Performance timers for console.time/timeEnd
	timersMu sync.Mutex           // Protects timers map
}

// NewConsole creates a new Console instance.
func NewConsole() *Console {
	return &Console{
		timers: make(map[string]time.Time),
	}
}

// Export creates and returns the console JavaScript object with all methods.
// This is called when the console global is initialized.
func (c *Console) Export(vm *goja.Runtime) goja.Value {
	c.vm = vm
	obj := vm.NewObject()

	obj.Set("log", c.consoleLog)
	obj.Set("error", c.consoleError)
	obj.Set("warn", c.consoleWarn)
	obj.Set("time", c.consoleTime)
	obj.Set("timeEnd", c.consoleTimeEnd)
	obj.Set("table", c.consoleTable)

	return obj
}

// consoleLog implements console.log() - outputs messages to stdout.
// Accepts multiple arguments of any type.
//
// JavaScript usage:
//
//	console.log('Hello', 'World', 123, {foo: 'bar'});
func (c *Console) consoleLog(call goja.FunctionCall) goja.Value {
	args := make([]any, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg.Export()
	}
	fmt.Println(args...)
	return goja.Undefined()
}

// consoleError implements console.error() - outputs error messages with ERROR prefix.
// Accepts multiple arguments of any type.
//
// JavaScript usage:
//
//	console.error('Something went wrong:', error);
func (c *Console) consoleError(call goja.FunctionCall) goja.Value {
	args := make([]any, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg.Export()
	}
	fmt.Print("ERROR: ")
	fmt.Println(args...)
	return goja.Undefined()
}

// consoleWarn implements console.warn() - outputs warning messages with WARN prefix.
// Accepts multiple arguments of any type.
//
// JavaScript usage:
//
//	console.warn('Deprecated function used');
func (c *Console) consoleWarn(call goja.FunctionCall) goja.Value {
	args := make([]any, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg.Export()
	}
	fmt.Print("WARN: ")
	fmt.Println(args...)
	return goja.Undefined()
}

// consoleTime implements console.time() - starts a performance timer.
// The timer is identified by an optional label (defaults to "default").
// Use console.timeEnd() with the same label to measure elapsed time.
//
// JavaScript usage:
//
//	console.time('myTimer');
//	// ... some operations ...
//	console.timeEnd('myTimer');  // Outputs: myTimer: 123.456ms
func (c *Console) consoleTime(call goja.FunctionCall) goja.Value {
	label := "default"
	if len(call.Arguments) > 0 {
		label = call.Arguments[0].String()
	}

	c.timersMu.Lock()
	c.timers[label] = time.Now()
	c.timersMu.Unlock()

	return goja.Undefined()
}

// consoleTimeEnd implements console.timeEnd() - stops a timer and displays elapsed time.
// Must be paired with a console.time() call using the same label.
// Outputs the duration in milliseconds with 3 decimal precision.
//
// If no timer exists with the given label, outputs a warning.
func (c *Console) consoleTimeEnd(call goja.FunctionCall) goja.Value {
	label := "default"
	if len(call.Arguments) > 0 {
		label = call.Arguments[0].String()
	}

	c.timersMu.Lock()
	startTime, exists := c.timers[label]
	if exists {
		delete(c.timers, label)
	}
	c.timersMu.Unlock()

	if !exists {
		fmt.Printf("Warning: No such label '%s' for console.timeEnd()\n", label)
		return goja.Undefined()
	}

	duration := time.Since(startTime)
	fmt.Printf("%s: %.3fms\n", label, float64(duration.Microseconds())/1000.0)

	return goja.Undefined()
}

// consoleTable implements console.table() - displays data in a formatted table.
// Supports arrays and objects. Data is displayed with borders and proper alignment.
//
// JavaScript usage:
//
//	console.table([1, 2, 3, 4]);  // Array table with index
//	console.table({name: 'Alice', age: 30});  // Object table with keys
func (c *Console) consoleTable(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		return goja.Undefined()
	}

	data := call.Arguments[0].Export()

	// Handle different data types
	switch v := data.(type) {
	case []any:
		c.printArrayTable(v)
	case map[string]any:
		c.printObjectTable(v)
	default:
		// Fallback to regular log for unsupported types
		fmt.Println(data)
	}

	return goja.Undefined()
}

// printArrayTable formats and prints an array as a table.
// Each element gets an index column and a value column.
func (c *Console) printArrayTable(data []any) {
	if len(data) == 0 {
		return
	}

	// Calculate column width based on content
	maxWidth := 36
	for _, item := range data {
		str := fmt.Sprintf("%v", item)
		if len(str) > maxWidth {
			maxWidth = len(str)
		}
	}
	if maxWidth > 60 {
		maxWidth = 60
	}

	// Print table header
	fmt.Println("┌─────────┬" + repeatChar('─', maxWidth+2) + "┐")
	fmt.Printf("│ (index) │ %-*s │\n", maxWidth, "Values")
	fmt.Println("├─────────┼" + repeatChar('─', maxWidth+2) + "┤")

	// Print table rows
	for i, item := range data {
		valueStr := fmt.Sprintf("%v", item)
		if len(valueStr) > maxWidth {
			valueStr = valueStr[:maxWidth-3] + "..."
		}
		fmt.Printf("│ %-7d │ %-*s │\n", i, maxWidth, valueStr)
	}

	// Print table footer
	fmt.Println("└─────────┴" + repeatChar('─', maxWidth+2) + "┘")
}

// printObjectTable formats and prints an object as a table.
// Keys are displayed in the first column, values in the second.
func (c *Console) printObjectTable(data map[string]any) {
	if len(data) == 0 {
		return
	}

	// Calculate column widths
	maxKeyWidth := 10
	maxValWidth := 24
	for key, value := range data {
		if len(key) > maxKeyWidth {
			maxKeyWidth = len(key)
		}
		valStr := fmt.Sprintf("%v", value)
		if len(valStr) > maxValWidth {
			maxValWidth = len(valStr)
		}
	}
	if maxKeyWidth > 30 {
		maxKeyWidth = 30
	}
	if maxValWidth > 50 {
		maxValWidth = 50
	}

	// Print table header
	fmt.Println("┌" + repeatChar('─', maxKeyWidth+2) + "┬" + repeatChar('─', maxValWidth+2) + "┐")
	fmt.Printf("│ %-*s │ %-*s │\n", maxKeyWidth, "(index)", maxValWidth, "Values")
	fmt.Println("├" + repeatChar('─', maxKeyWidth+2) + "┼" + repeatChar('─', maxValWidth+2) + "┤")

	// Print table rows
	for key, value := range data {
		keyStr := key
		if len(keyStr) > maxKeyWidth {
			keyStr = keyStr[:maxKeyWidth-3] + "..."
		}
		valueStr := fmt.Sprintf("%v", value)
		if len(valueStr) > maxValWidth {
			valueStr = valueStr[:maxValWidth-3] + "..."
		}
		fmt.Printf("│ %-*s │ %-*s │\n", maxKeyWidth, keyStr, maxValWidth, valueStr)
	}

	// Print table footer
	fmt.Println("└" + repeatChar('─', maxKeyWidth+2) + "┴" + repeatChar('─', maxValWidth+2) + "┘")
}

// Helper function to repeat a character n times
func repeatChar(char rune, count int) string {
	result := make([]rune, count)
	for i := 0; i < count; i++ {
		result[i] = char
	}
	return string(result)
}
