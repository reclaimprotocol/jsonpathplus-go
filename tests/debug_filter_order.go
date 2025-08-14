package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Copy the relevant filter patterns
func testFilterPatterns() {
	expr := `.items[0].product === "laptop"`
	
	fmt.Printf("Testing expression: %q\n", expr)
	fmt.Println(strings.Repeat("=", 50))
	
	// Test comparison filter regex (tried before array wildcard)
	comparisonRe := regexp.MustCompile(`\.([a-zA-Z_]\w*(?:\.[a-zA-Z_]\w*)*)\s*(===|!==|<=|>=|==|!=|<|>)\s*(.+)`)
	compMatches := comparisonRe.FindStringSubmatch(expr)
	
	fmt.Printf("Comparison filter regex:\n")
	fmt.Printf("Pattern: %s\n", comparisonRe.String())
	if len(compMatches) == 0 {
		fmt.Println("❌ No matches")
	} else {
		fmt.Printf("✅ Matches: %d\n", len(compMatches))
		for i, match := range compMatches {
			fmt.Printf("  [%d]: %q\n", i, match)
		}
	}
	
	fmt.Println()
	
	// Test array wildcard filter regex
	arrayRe := regexp.MustCompile(`\.([a-zA-Z_]\w*)\[\*\]\.([a-zA-Z_]\w*(?:\.[a-zA-Z_]\w*)*)\s*(===|!==|<=|>=|==|!=|<|>)\s*(.+)`)
	arrayMatches := arrayRe.FindStringSubmatch(expr)
	
	fmt.Printf("Array wildcard filter regex:\n")
	fmt.Printf("Pattern: %s\n", arrayRe.String())
	if len(arrayMatches) == 0 {
		fmt.Println("❌ No matches")
	} else {
		fmt.Printf("✅ Matches: %d\n", len(arrayMatches))
		for i, match := range arrayMatches {
			fmt.Printf("  [%d]: %q\n", i, match)
		}
	}
	
	fmt.Println()
	
	// Test with actual array wildcard expression
	wildcardExpr := `.items[*].product === "laptop"`
	fmt.Printf("Testing wildcard expression: %q\n", wildcardExpr)
	
	compMatches2 := comparisonRe.FindStringSubmatch(wildcardExpr)
	fmt.Printf("Comparison filter on wildcard: %s\n", func() string {
		if len(compMatches2) == 0 { return "❌ No matches" }
		return "✅ Matches"
	}())
	
	arrayMatches2 := arrayRe.FindStringSubmatch(wildcardExpr)
	fmt.Printf("Array wildcard filter on wildcard: %s\n", func() string {
		if len(arrayMatches2) == 0 { return "❌ No matches" }
		return "✅ Matches"
	}())
}

func main() {
	testFilterPatterns()
}