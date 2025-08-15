package jsonpathplus

import (
	"fmt"
	"regexp"
	"strings"
)

// Copy of cleanFilterExpression function to test
func cleanFilterExpression(expr string) string {
	expr = strings.TrimSpace(expr)
	// Don't remove @ symbols for context-based expressions
	// Only remove leading @ for simple property references like "@.field"
	if strings.HasPrefix(expr, "@.") {
		expr = strings.TrimPrefix(expr, "@")
	}
	// Preserve context keywords starting with @
	if strings.HasPrefix(expr, "@parentProperty") || strings.HasPrefix(expr, "@parent") || strings.HasPrefix(expr, "@property") || strings.HasPrefix(expr, "@path") {
		return strings.TrimSpace(expr)
	}
	// Allow direct-value comparisons like "@=== 'x'" by stripping leading @
	if strings.HasPrefix(expr, "@") {
		expr = strings.TrimPrefix(expr, "@")
	}
	// Simplify by unconditionally trimming the "@ " prefix if present
	expr = strings.TrimPrefix(expr, "@ ")
	return strings.TrimSpace(expr)
}

func debug_filter_exprMain() {
	originalExpr := `@.items[*].product === "laptop"`

	// Test what happens to the expression after cleaning
	cleanedExpr := cleanFilterExpression(originalExpr)
	fmt.Printf("Original: %q\n", originalExpr)
	fmt.Printf("Cleaned:  %q\n", cleanedExpr)

	// Test the array wildcard regex
	re := regexp.MustCompile(`\.([a-zA-Z_]\w*)\[\*\]\.([a-zA-Z_]\w*(?:\.[a-zA-Z_]\w*)*)\s*(===|!==|<=|>=|==|!=|<|>)\s*(.+)`)
	matches := re.FindStringSubmatch(cleanedExpr)

	fmt.Printf("\nArray wildcard regex test:\n")
	fmt.Printf("Pattern: %s\n", re.String())
	if len(matches) == 0 {
		fmt.Println("❌ No matches")
	} else {
		fmt.Printf("✅ Matches: %d\n", len(matches))
		for i, match := range matches {
			fmt.Printf("  [%d]: %q\n", i, match)
		}
	}

	// Test nested filter cleaning
	fmt.Printf("\n" + strings.Repeat("=", 50) + "\n")

	nestedOriginal := `@.items[?(@.product === "laptop")]`
	nestedCleaned := cleanFilterExpression(nestedOriginal)
	fmt.Printf("Nested Original: %q\n", nestedOriginal)
	fmt.Printf("Nested Cleaned:  %q\n", nestedCleaned)

	// Test nested filter regex
	nestedRe := regexp.MustCompile(`\.(\w+)\[\?\(([^)]+)\)\]`)
	nestedMatches := nestedRe.FindStringSubmatch(nestedCleaned)

	fmt.Printf("\nNested filter regex test:\n")
	fmt.Printf("Pattern: %s\n", nestedRe.String())
	if len(nestedMatches) == 0 {
		fmt.Println("❌ No matches")
	} else {
		fmt.Printf("✅ Matches: %d\n", len(nestedMatches))
		for i, match := range nestedMatches {
			fmt.Printf("  [%d]: %q\n", i, match)
		}
	}
}
