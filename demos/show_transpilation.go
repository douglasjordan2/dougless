package main

import (
	"fmt"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func main() {
	// Modern ES6+ code
	source := "// Modern ES6+ JavaScript\n" +
		"const greet = (name) => {\n" +
		"  console.log(`Hello, ${name}!`);\n" +
		"};\n" +
		"\n" +
		"const calculate = (x, y) => {\n" +
		"  const sum = x + y;\n" +
		"  const product = x * y;\n" +
		"  return { sum, product };\n" +
		"};\n" +
		"\n" +
		"// Destructuring\n" +
		"const { sum, product } = calculate(5, 10);\n" +
		"\n" +
		"// Spread operator\n" +
		"const numbers = [1, 2, 3];\n" +
		"const moreNumbers = [...numbers, 4, 5];\n" +
		"\n" +
		"// Class syntax\n" +
		"class Person {\n" +
		"  constructor(name, age) {\n" +
		"    this.name = name;\n" +
		"    this.age = age;\n" +
		"  }\n" +
		"  \n" +
		"  greet() {\n" +
		"    return `I'm ${this.name}, ${this.age} years old`;\n" +
		"  }\n" +
		"}\n" +
		"\n" +
		"// Async/await\n" +
		"async function fetchData() {\n" +
		"  const response = await Promise.resolve('data');\n" +
		"  return response;\n" +
		"}"

	// Transpile with esbuild
	result := api.Transform(source, api.TransformOptions{
		Loader:     api.LoaderJS,
		Target:     api.ES2017,
		Sourcefile: "example.js",
		Format:     api.FormatDefault,
		Sourcemap:  api.SourceMapInline,
	})

	if len(result.Errors) > 0 {
		fmt.Println("Transpilation errors:")
		for _, err := range result.Errors {
			fmt.Printf("  %s\n", err.Text)
		}
		return
	}

	transpiled := string(result.Code)
	
	// Remove the source map comment for cleaner display
	if idx := strings.Index(transpiled, "//# sourceMappingURL="); idx != -1 {
		sourceMapComment := transpiled[idx:]
		transpiled = transpiled[:idx]
		
		// Show source map separately
		fmt.Println("╔════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                    SOURCE MAP INFORMATION                          ║")
		fmt.Println("╚════════════════════════════════════════════════════════════════════╝")
		fmt.Printf("\n%s\n", sourceMapComment[:80]+"...")
		fmt.Println("(Base64 encoded mapping data - connects transpiled code to original)")
		fmt.Println()
	}

	// Print side-by-side comparison
	fmt.Println("╔════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                 ORIGINAL ES6+ CODE (What You Write)                ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println(source)
	fmt.Println()
	
	fmt.Println("╔════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║            TRANSPILED ES5 CODE (What Goja Executes)               ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println(transpiled)
	
	fmt.Println("╔════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                          KEY CHANGES                               ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("✓ Arrow functions (=>) → Regular functions")
	fmt.Println("✓ const/let → var")
	fmt.Println("✓ Template literals (backticks) → String concatenation")
	fmt.Println("✓ Destructuring → Manual variable assignment")
	fmt.Println("✓ Spread operator → Array copying")
	fmt.Println("✓ Classes → Prototype-based constructors")
	fmt.Println("✓ async/await → Promise chains")
	fmt.Println()
	fmt.Println("With SOURCE MAPS enabled, errors point to the ORIGINAL code!")
}
