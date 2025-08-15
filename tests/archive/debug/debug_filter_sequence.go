package jsonpathplus

import (
	"fmt"
	"regexp"
	"strings"
)

func testDirectComparisonFilter(expr string) (bool, bool) {
	fmt.Printf("Testing direct comparison filter on: %q\n", expr)
	re := regexp.MustCompile(`^\s*(===|!==|<=|>=|==|!=|<|>)\s*(.+)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 3 {
		fmt.Printf("  ❌ No matches for direct comparison\n")
		return false, false
	}
	fmt.Printf("  ✅ Matches direct comparison pattern\n")
	return true, true
}

func testComparisonFilter(expr string) (bool, bool) {
	fmt.Printf("Testing comparison filter on: %q\n", expr)
	
	// Check for array access exclusion
	if strings.Contains(expr, "[") {
		fmt.Printf("  ❌ Contains '[' - excluded\n")
		return false, false
	}
	
	re := regexp.MustCompile(`\.([a-zA-Z_]\w*(?:\.[a-zA-Z_]\w*)*)\s*(===|!==|<=|>=|==|!=|<|>)\s*(.+)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 4 {
		fmt.Printf("  ❌ No matches for comparison\n")
		return false, false
	}
	fmt.Printf("  ✅ Matches comparison pattern\n")
	return true, true
}

func testArrayWildcardFilter(expr string) (bool, bool) {
	fmt.Printf("Testing array wildcard filter on: %q\n", expr)
	re := regexp.MustCompile(`\.([a-zA-Z_]\w*)\[\*\]\.([a-zA-Z_]\w*(?:\.[a-zA-Z_]\w*)*)\s*(===|!==|<=|>=|==|!=|<|>)\s*(.+)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 5 {
		fmt.Printf("  ❌ No matches for array wildcard\n")
		return false, false
	}
	fmt.Printf("  ✅ Matches array wildcard pattern\n")
	return true, true
}

func simulateFilterOrder(expr string) {
	fmt.Printf("=== SIMULATING FILTER ORDER ===\n")
	fmt.Printf("Expression: %q\n", expr)
	fmt.Println(strings.Repeat("-", 40))
	
	// Order from the actual code:
	
	// tryDirectComparisonFilter
	if result, ok := testDirectComparisonFilter(expr); ok {
		fmt.Printf("🔴 STOPPED at direct comparison filter (result: %t)\n", result)
		return
	}
	
	// tryArrayWildcardFilter  
	if result, ok := testArrayWildcardFilter(expr); ok {
		fmt.Printf("🟢 STOPPED at array wildcard filter (result: %t)\n", result)
		return
	}
	
	// tryComparisonFilter
	if result, ok := testComparisonFilter(expr); ok {
		fmt.Printf("🟡 STOPPED at comparison filter (result: %t)\n", result)
		return
	}
	
	fmt.Printf("🔴 No filter matched\n")
}

func debug_filter_sequenceMain() {
	expressions := []string{
		`.items[*].product === "laptop"`,
		`.items[0].product === "laptop"`,
		`.name === "Alice"`, // Should work with comparison filter
	}
	
	for _, expr := range expressions {
		simulateFilterOrder(expr)
		fmt.Println()
	}
}