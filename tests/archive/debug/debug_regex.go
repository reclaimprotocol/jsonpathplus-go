package jsonpathplus

import (
	"fmt"
	"regexp"
)

func debug_regexMain() {
	// Test the nested filter regex
	expr := `@.items[?(@.product === "laptop")]`

	patterns := []string{
		`@\.([a-zA-Z_]\w*)(\[[^\[\]]*\])?\[\\?\(([^)]+)\)\]`, // Current pattern
		`@\.([a-zA-Z_]\w*)\[\\?\(([^)]+)\)\]`,                // Simplified pattern
		`@\.(\w+)\[\?\(([^)]+)\)\]`,                          // Even simpler
	}

	fmt.Printf("Testing expression: %s\n", expr)
	fmt.Printf("==================================================\n")

	for i, pattern := range patterns {
		fmt.Printf("\nPattern %d: %s\n", i+1, pattern)
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(expr)

		if len(matches) == 0 {
			fmt.Printf("❌ No matches\n")
		} else {
			fmt.Printf("✅ Matches: %d\n", len(matches))
			for j, match := range matches {
				fmt.Printf("  [%d]: %q\n", j, match)
			}
		}
	}

	// Let's also test manually building the pattern piece by piece
	fmt.Printf("\n==================================================\n")
	fmt.Printf("Manual pattern building:\n")

	testParts := []string{
		`@\.`,            // @.
		`([a-zA-Z_]\w*)`, // items
		`\[`,             // [
		`\?\(`,           // ?(
		`([^)]+)`,        // @.product === "laptop"
		`\)\]`,           // )]
	}

	fullPattern := ""
	for i, part := range testParts {
		fullPattern += part
		fmt.Printf("Step %d: %s -> %s\n", i+1, part, fullPattern)
		re := regexp.MustCompile(fullPattern)
		if re.MatchString(expr) {
			fmt.Printf("  ✅ Matches so far\n")
		} else {
			fmt.Printf("  ❌ No match\n")
		}
	}
}
