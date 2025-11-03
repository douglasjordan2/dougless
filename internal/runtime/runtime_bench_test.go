package runtime

import (
	"os"
	"testing"
)

// BenchmarkRuntimeCreation measures the cost of creating a new Runtime instance
func BenchmarkRuntimeCreation(b *testing.B) {
	argv := []string{"dougless", "test.js"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(argv)
	}
}

// BenchmarkSimpleExecution measures execution of a trivial script
func BenchmarkSimpleExecution(b *testing.B) {
	script := `console.log("Hello, World!");`
	argv := []string{"dougless", "test.js"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt := New(argv)
		if err := rt.Execute(script, "bench.js"); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkTranspilation measures the cost of ES6+ transpilation
func BenchmarkTranspilation(b *testing.B) {
	source := `
		const greeting = (name) => {
			return "Hello, " + name + "!";
		};
		
		const result = greeting("World");
		console.log(result);
	`
	argv := []string{"dougless", "test.js"}
	rt := New(argv)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := rt.transpile(source, "bench.js")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkLargeScriptExecution measures execution of a larger script
func BenchmarkLargeScriptExecution(b *testing.B) {
	script := `
		function fibonacci(n) {
			if (n <= 1) return n;
			return fibonacci(n - 1) + fibonacci(n - 2);
		}
		
		const result = fibonacci(10);
		console.log(result);
	`
	argv := []string{"dougless", "test.js"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt := New(argv)
		if err := rt.Execute(script, "bench.js"); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkES6Features measures execution of ES6+ features
func BenchmarkES6Features(b *testing.B) {
	script := `
		const arr = [1, 2, 3, 4, 5];
		const doubled = arr.map(x => x * 2);
		const sum = doubled.reduce((a, b) => a + b, 0);
		
		const obj = { a: 1, b: 2, c: 3 };
		const { a, b, ...rest } = obj;
		
		const template = "Sum is " + sum;
		console.log(template);
	`
	argv := []string{"dougless", "test.js"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt := New(argv)
		if err := rt.Execute(script, "bench.js"); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkAsyncAwait measures async/await performance
func BenchmarkAsyncAwait(b *testing.B) {
	script := `
		async function delay(ms) {
			return new Promise(resolve => setTimeout(resolve, ms));
		}
		
		async function main() {
			await delay(1);
			console.log("Done");
		}
		
		main();
	`
	argv := []string{"dougless", "test.js"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt := New(argv)
		if err := rt.Execute(script, "bench.js"); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkModuleRequire measures require() overhead
func BenchmarkModuleRequire(b *testing.B) {
	script := `
		const path1 = require('path');
		const path2 = require('path'); // Should hit cache
		console.log(path1.sep);
	`
	argv := []string{"dougless", "test.js"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt := New(argv)
		if err := rt.Execute(script, "bench.js"); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConsoleLogNew measures console.log overhead with new runtime per iteration
func BenchmarkConsoleLogNew(b *testing.B) {
	script := `console.log("test");`
	argv := []string{"dougless", "test.js"}
	
	// Redirect stdout to suppress output
	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt := New(argv)
		if err := rt.Execute(script, "bench.js"); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSetTimeoutNew measures setTimeout overhead with new runtime per iteration
func BenchmarkSetTimeoutNew(b *testing.B) {
	script := `setTimeout(() => {}, 0);`
	argv := []string{"dougless", "test.js"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt := New(argv)
		if err := rt.Execute(script, "bench.js"); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPromiseCreation measures Promise creation overhead
func BenchmarkPromiseCreation(b *testing.B) {
	script := `new Promise((resolve) => resolve(42));`
	argv := []string{"dougless", "test.js"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt := New(argv)
		if err := rt.Execute(script, "bench.js"); err != nil {
			b.Fatal(err)
		}
	}
}
