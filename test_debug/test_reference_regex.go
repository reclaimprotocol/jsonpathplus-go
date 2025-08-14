package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
	"regexp"
)

func main() {
	// Test Go regex directly
	fmt.Println("=== Testing Go regex directly ===")
	pattern := "(?i)TION$"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	testStrings := []string{"reference", "fiction", "action", "science"}
	for _, str := range testStrings {
		matches := regex.MatchString(str)
		fmt.Printf("'%s' matches /%s/: %t\n", str, "TION$", matches)
	}

	// Test with JSONPath
	fmt.Println("\n=== Testing with JSONPath ===")
	jsonData := `{
		"categories": ["reference", "fiction", "action", "science"]
	}`

	results, err := jp.Query("$.categories[?(@.match(/TION$/i))]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Found %d matches:\n", len(results))
	for i, r := range results {
		fmt.Printf("  [%d] '%s'\n", i, r.Value)
	}

	// Test each individually
	fmt.Println("\n=== Testing each category individually ===")
	for _, cat := range testStrings {
		jsonData := fmt.Sprintf(`{"test": "%s"}`, cat)
		results, err := jp.Query("$.test[?(@.match(/TION$/i))]", jsonData)
		if err != nil {
			fmt.Printf("Error for '%s': %v\n", cat, err)
			continue
		}
		fmt.Printf("'%s': %d matches\n", cat, len(results))
	}
}
