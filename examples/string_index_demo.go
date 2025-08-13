package main

import (
	"fmt"
	"strings"
	jp "jsonpathplus-go"
)

func main() {
	fmt.Println("🧪 JSONPath String Index Functionality Demo")
	fmt.Println("===========================================")
	fmt.Println()

	// Test 1: Basic Property Key Positions
	fmt.Println("📋 Test 1: Basic Property Key Positions")
	testBasicPropertyPositions()
	fmt.Println()

	// Test 2: Array Element Positions
	fmt.Println("📋 Test 2: Array Element Positions")
	testArrayElementPositions()
	fmt.Println()

	// Test 3: Nested Object Positions
	fmt.Println("📋 Test 3: Nested Object Positions")
	testNestedObjectPositions()
	fmt.Println()

	// Test 4: Whitespace Preservation
	fmt.Println("📋 Test 4: Whitespace Preservation")
	testWhitespacePreservation()
	fmt.Println()

	// Test 5: Complex Nested Structures
	fmt.Println("📋 Test 5: Complex Nested Structures")
	testComplexNestedStructures()
	fmt.Println()

	// Test 6: Edge Cases
	fmt.Println("📋 Test 6: Edge Cases")
	testEdgeCases()
	fmt.Println()

	fmt.Println("🎉 String Index Functionality Demo Complete!")
}

func testBasicPropertyPositions() {
	jsonStr := `{"id":123,"name":"test","active":true}`
	fmt.Printf("JSON: %s\n", jsonStr)
	showPositions(jsonStr)
	
	tests := []struct {
		query    string
		expected int
		desc     string
	}{
		{"$.id", 1, "Property 'id' key position"},
		{"$.name", 10, "Property 'name' key position"},
		{"$.active", 24, "Property 'active' key position"},
	}

	for _, test := range tests {
		results, err := jp.QueryWithStringIndex(test.query, jsonStr)
		if err != nil {
			fmt.Printf("❌ %s failed: %v\n", test.desc, err)
			continue
		}
		
		if len(results) == 0 {
			fmt.Printf("❌ %s: No results\n", test.desc)
			continue
		}
		
		result := results[0]
		if result.OriginalIndex == test.expected {
			fmt.Printf("✅ %s: Expected %d, Got %d ('%c')\n", 
				test.desc, test.expected, result.OriginalIndex, jsonStr[result.OriginalIndex])
		} else {
			fmt.Printf("❌ %s: Expected %d, Got %d\n", 
				test.desc, test.expected, result.OriginalIndex)
		}
	}
}

func testArrayElementPositions() {
	jsonStr := `["first","second","third"]`
	fmt.Printf("JSON: %s\n", jsonStr)
	showPositions(jsonStr)
	
	tests := []struct {
		query    string
		desc     string
	}{
		{"$[0]", "First array element"},
		{"$[1]", "Second array element"}, 
		{"$[2]", "Third array element"},
		{"$[*]", "All array elements"},
	}

	for _, test := range tests {
		results, err := jp.QueryWithStringIndex(test.query, jsonStr)
		if err != nil {
			fmt.Printf("❌ %s failed: %v\n", test.desc, err)
			continue
		}
		
		if len(results) == 0 {
			fmt.Printf("❌ %s: No results\n", test.desc)
			continue
		}
		
		fmt.Printf("📍 %s:\n", test.desc)
		for i, result := range results {
			if result.OriginalIndex < len(jsonStr) {
				fmt.Printf("   Result %d: Value='%v', Index=%d ('%c')\n", 
					i, result.Value, result.OriginalIndex, jsonStr[result.OriginalIndex])
			} else {
				fmt.Printf("   Result %d: Value='%v', Index=%d (out of bounds)\n", 
					i, result.Value, result.OriginalIndex)
			}
		}
	}
}

func testNestedObjectPositions() {
	jsonStr := `{"user":{"name":"john","age":30},"status":"active"}`
	fmt.Printf("JSON: %s\n", jsonStr)
	showPositions(jsonStr)
	
	tests := []struct {
		query    string
		desc     string
	}{
		{"$.user", "Nested object"},
		{"$.user.name", "Nested property 'name'"},
		{"$.user.age", "Nested property 'age'"},
		{"$.status", "Root level 'status'"},
	}

	for _, test := range tests {
		results, err := jp.QueryWithStringIndex(test.query, jsonStr)
		if err != nil {
			fmt.Printf("❌ %s failed: %v\n", test.desc, err)
			continue
		}
		
		if len(results) == 0 {
			fmt.Printf("❌ %s: No results\n", test.desc)
			continue
		}
		
		result := results[0]
		if result.OriginalIndex < len(jsonStr) {
			fmt.Printf("✅ %s: Value='%v', Index=%d ('%c')\n", 
				test.desc, result.Value, result.OriginalIndex, jsonStr[result.OriginalIndex])
		} else {
			fmt.Printf("❌ %s: Index %d out of bounds\n", test.desc, result.OriginalIndex)
		}
	}
}

func testWhitespacePreservation() {
	jsonStr := `{
  "id": 123,
  "data": {
    "name": "test",
    "values": [1, 2, 3]
  }
}`
	fmt.Printf("JSON:\n%s\n", jsonStr)
	
	tests := []struct {
		query string
		desc  string
	}{
		{"$.id", "Property 'id' with whitespace"},
		{"$.data", "Property 'data' with whitespace"},
		{"$.data.name", "Nested property 'name' with whitespace"},
		{"$.data.values", "Nested array 'values' with whitespace"},
	}

	for _, test := range tests {
		results, err := jp.QueryWithStringIndex(test.query, jsonStr)
		if err != nil {
			fmt.Printf("❌ %s failed: %v\n", test.desc, err)
			continue
		}
		
		if len(results) == 0 {
			fmt.Printf("❌ %s: No results\n", test.desc)
			continue
		}
		
		result := results[0]
		if result.OriginalIndex < len(jsonStr) {
			// Show context around the position
			start := max(0, result.OriginalIndex-5)
			end := min(len(jsonStr), result.OriginalIndex+10)
			context := jsonStr[start:end]
			context = strings.ReplaceAll(context, "\n", "\\n")
			context = strings.ReplaceAll(context, "  ", "·")
			
			fmt.Printf("✅ %s: Index=%d, Context='%s'\n", 
				test.desc, result.OriginalIndex, context)
		} else {
			fmt.Printf("❌ %s: Index %d out of bounds\n", test.desc, result.OriginalIndex)
		}
	}
}

func testComplexNestedStructures() {
	jsonStr := `{"company":{"departments":[{"name":"eng","employees":[{"name":"alice","id":1}]}],"founded":2020}}`
	fmt.Printf("JSON: %s\n", jsonStr)
	showPositions(jsonStr)
	
	tests := []struct {
		query string
		desc  string
	}{
		{"$.company", "Root company object"},
		{"$.company.departments", "Departments array"},
		{"$.company.departments[0]", "First department"},
		{"$.company.departments[0].name", "Department name"},
		{"$.company.departments[0].employees", "Employees array"},
		{"$.company.departments[0].employees[0]", "First employee"},
		{"$.company.departments[0].employees[0].name", "Employee name"},
		{"$.company.departments[0].employees[0].id", "Employee id"},
		{"$.company.founded", "Company founded year"},
	}

	for _, test := range tests {
		results, err := jp.QueryWithStringIndex(test.query, jsonStr)
		if err != nil {
			fmt.Printf("❌ %s failed: %v\n", test.desc, err)
			continue
		}
		
		if len(results) == 0 {
			fmt.Printf("❌ %s: No results\n", test.desc)
			continue
		}
		
		result := results[0]
		if result.OriginalIndex < len(jsonStr) && result.OriginalIndex >= 0 {
			char := jsonStr[result.OriginalIndex]
			fmt.Printf("📍 %s: Value='%v', Index=%d ('%c')\n", 
				test.desc, result.Value, result.OriginalIndex, char)
		} else {
			fmt.Printf("❌ %s: Invalid index %d\n", test.desc, result.OriginalIndex)
		}
	}
}

func testEdgeCases() {
	tests := []struct {
		jsonStr string
		query   string
		desc    string
	}{
		{`{}`, "$", "Empty object"},
		{`[]`, "$", "Empty array"},
		{`{"":123}`, `$[""]`, "Empty string key"},
		{`{"key with spaces":true}`, `$["key with spaces"]`, "Key with spaces"},
		{`{"special\"chars":1}`, `$["special\"chars"]`, "Key with escape chars"},
		{`null`, "$", "Null value"},
		{`"simple string"`, "$", "Simple string"},
		{`123`, "$", "Simple number"},
	}

	for _, test := range tests {
		fmt.Printf("Testing: %s\n", test.desc)
		fmt.Printf("JSON: %s\n", test.jsonStr)
		
		results, err := jp.QueryWithStringIndex(test.query, test.jsonStr)
		if err != nil {
			fmt.Printf("❌ %s failed: %v\n", test.desc, err)
			continue
		}
		
		if len(results) == 0 {
			fmt.Printf("❌ %s: No results\n", test.desc)
			continue
		}
		
		result := results[0]
		if result.OriginalIndex >= 0 && result.OriginalIndex < len(test.jsonStr) {
			char := test.jsonStr[result.OriginalIndex]
			fmt.Printf("✅ %s: Value='%v', Index=%d ('%c')\n", 
				test.desc, result.Value, result.OriginalIndex, char)
		} else {
			fmt.Printf("📍 %s: Value='%v', Index=%d\n", 
				test.desc, result.Value, result.OriginalIndex)
		}
		fmt.Println()
	}
}

// Helper function to show character positions in a string
func showPositions(jsonStr string) {
	if len(jsonStr) > 80 {
		fmt.Printf("Positions (first 80 chars): ")
		for i := 0; i < min(80, len(jsonStr)); i++ {
			fmt.Printf("%d", i%10)
		}
		fmt.Println("...")
	} else {
		fmt.Printf("Positions: ")
		for i := 0; i < len(jsonStr); i++ {
			fmt.Printf("%d", i%10)
		}
		fmt.Println()
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}