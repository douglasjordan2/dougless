//go:build ignore
// +build ignore

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
		"}"

	fmt.Println("╔════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                 ORIGINAL ES6+ CODE (What You Write)                ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println(source)
	fmt.Println()

	// Transpile to ES5 (strict old-school)
	resultES5 := api.Transform(source, api.TransformOptions{
		Loader:     api.LoaderJS,
		Target:     api.ES5, // Pure ES5!
		Sourcefile: "example.js",
		Format:     api.FormatDefault,
		Sourcemap:  api.SourceMapNone,
	})

	if len(resultES5.Errors) > 0 {
		fmt.Println("ES5 Transpilation errors:")
		for _, err := range resultES5.Errors {
			fmt.Printf("  %s\n", err.Text)
		}
	} else {
		transpiledES5 := string(resultES5.Code)

		fmt.Println("╔════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║              TRANSPILED TO ES5 (Pure Old-School JS)                ║")
		fmt.Println("╚════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println(transpiledES5)
		fmt.Println()
	}

	// Transpile to ES2017 (what Dougless uses)
	resultES2017 := api.Transform(source, api.TransformOptions{
		Loader:     api.LoaderJS,
		Target:     api.ES2017, // Modern ES2017
		Sourcefile: "example.js",
		Format:     api.FormatDefault,
		Sourcemap:  api.SourceMapInline,
	})

	if len(resultES2017.Errors) > 0 {
		fmt.Println("ES2017 Transpilation errors:")
		for _, err := range resultES2017.Errors {
			fmt.Printf("  %s\n", err.Text)
		}
	} else {
		transpiledES2017 := string(resultES2017.Code)

		// Remove source map for display
		if idx := strings.Index(transpiledES2017, "//# sourceMappingURL="); idx != -1 {
			sourceMapComment := transpiledES2017[idx:]
			transpiledES2017 = transpiledES2017[:idx]

			fmt.Println("╔════════════════════════════════════════════════════════════════════╗")
			fmt.Println("║         TRANSPILED TO ES2017 (What Dougless Actually Uses)        ║")
			fmt.Println("╚════════════════════════════════════════════════════════════════════╝")
			fmt.Println()
			fmt.Println(transpiledES2017)

			fmt.Println("╔════════════════════════════════════════════════════════════════════╗")
			fmt.Println("║                    SOURCE MAP (First 200 chars)                    ║")
			fmt.Println("╚════════════════════════════════════════════════════════════════════╝")
			fmt.Println()
			if len(sourceMapComment) > 200 {
				fmt.Println(sourceMapComment[:200] + "...")
			} else {
				fmt.Println(sourceMapComment)
			}
			fmt.Println()
			fmt.Println("This base64 data contains mappings like:")
			fmt.Println("  'Line 15, Col 10 in transpiled code → Line 3, Col 8 in original'")
			fmt.Println()
		}
	}

	fmt.Println("╔════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                        KEY OBSERVATIONS                            ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("ES5 Target (strict backward compatibility):")
	fmt.Println("  ✓ Arrow functions → Regular function expressions")
	fmt.Println("  ✓ const/let → var")
	fmt.Println("  ✓ Template literals → String concatenation with +")
	fmt.Println("  ✓ Destructuring → Manual property access")
	fmt.Println("  ✓ Spread operator → Array.prototype.concat")
	fmt.Println("  ✓ Classes → Prototype-based constructors")
	fmt.Println()
	fmt.Println("ES2017 Target (Dougless Runtime):")
	fmt.Println("  ✓ Keeps most modern syntax (arrow functions, classes, etc.)")
	fmt.Println("  ✓ Async/await supported natively")
	fmt.Println("  ✓ Minimal changes = faster transpilation")
	fmt.Println("  ✓ Source maps still crucial for property renaming & line shifts")
	fmt.Println()
	fmt.Println("🎯 With source maps: When an error occurs at line X in transpiled")
	fmt.Println("   code, it's automatically mapped back to line Y in your original!")
}
