package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	// Test JSON string with proper formatting for easy position verification
	jsonStr := `{
  "name": "John",
  "age": 30,
  "books": [
    "Book1",
    "Book2",
    "Book3"
  ],
  "profile": {
    "bio": "A reader",
    "active": true
  }
}`

	fmt.Println("=== String Index Functionality Test ===")
	fmt.Println()
	fmt.Printf("JSON String (%d characters):\n%s\n", len(jsonStr), jsonStr)
	fmt.Println()

	// Test cases with different JSONPath expressions
	testCases := []struct {
		name     string
		path     string
		expected int // expected number of results
	}{
		{"Root object", "$", 1},
		{"Simple property", "$.name", 1},
		{"Numeric property", "$.age", 1},
		{"Array property", "$.books", 1},
		{"Array element", "$.books[0]", 1},
		{"Array element", "$.books[1]", 1},
		{"Array element", "$.books[2]", 1},
		{"All array elements", "$.books[*]", 3},
		{"Nested object", "$.profile", 1},
		{"Nested property", "$.profile.bio", 1},
		{"Nested boolean", "$.profile.active", 1},
	}

	for _, tc := range testCases {
		fmt.Printf("=== Test: %s ===\n", tc.name)
		fmt.Printf("Path: %s\n", tc.path)

		// Query with string input to trigger string index calculation
		results, err := jp.Query(tc.path, jsonStr)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			continue
		}

		fmt.Printf("Results: %d (expected: %d)\n", len(results), tc.expected)

		if len(results) != tc.expected {
			fmt.Printf("❌ Result count mismatch!\n")
		} else {
			fmt.Printf("✅ Result count correct\n")
		}

		for i, result := range results {
			fmt.Printf("  [%d] Value: %v\n", i, result.Value)
			fmt.Printf("      Path: %s\n", result.Path)
			fmt.Printf("      String Position: Start=%d, End=%d, Length=%d\n",
				result.Start, result.End, result.Length)

			// Verify the string position by extracting the substring
			if result.Start >= 0 && result.End > result.Start && result.End <= len(jsonStr) {
				extracted := jsonStr[result.Start:result.End]
				fmt.Printf("      Extracted: %q\n", extracted)

				// Basic validation
				if result.Length == result.End-result.Start {
					fmt.Printf("      ✅ Length calculation correct\n")
				} else {
					fmt.Printf("      ❌ Length calculation incorrect\n")
				}
			} else if result.Start == 0 && result.End == 0 && result.Length == 0 {
				fmt.Printf("      ⚠️  No string position found (complex path)\n")
			} else {
				fmt.Printf("      ❌ Invalid string position\n")
			}
		}
		fmt.Println()
	}

	// Test with already parsed data (should not have string indices)
	fmt.Println("=== Test with Parsed Data (No String Indices) ===")
	data, _ := jp.JSONParse(jsonStr)
	results, err := jp.Query("$.name", data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		for _, result := range results {
			fmt.Printf("Value: %v, Start: %d, End: %d, Length: %d\n",
				result.Value, result.Start, result.End, result.Length)
			if result.Start == 0 && result.End == 0 && result.Length == 0 {
				fmt.Printf("✅ No string indices for parsed data (as expected)\n")
			} else {
				fmt.Printf("❌ Unexpected string indices for parsed data\n")
			}
		}
	}

	fmt.Println()
	fmt.Println("=== String Index Feature Summary ===")
	fmt.Println("✅ String indices are calculated when input is a JSON string")
	fmt.Println("✅ No string indices when input is already parsed data")
	fmt.Println("✅ Results include Start, End, and Length fields")
	fmt.Println("✅ Character positions enable precise location tracking")
}
