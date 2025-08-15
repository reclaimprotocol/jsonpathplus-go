package jsonpathplus

import (
	"fmt"
	"regexp"
	"strings"
	"github.com/reclaimprotocol/jsonpathplus-go/pkg/utils"
	"encoding/json"
)

// Simulate the array wildcard filter step by step
func debugArrayWildcardFilter(expr string, current interface{}) {
	fmt.Printf("=== DEBUG ARRAY WILDCARD FILTER ===\n")
	fmt.Printf("Expression: %q\n", expr)
	fmt.Printf("Current object: %T\n", current)
	
	// Print the current object
	if jsonBytes, err := json.MarshalIndent(current, "", "  "); err == nil {
		fmt.Printf("Current value:\n%s\n", string(jsonBytes))
	}
	
	// Test the regex
	re := regexp.MustCompile(`\.([a-zA-Z_]\w*)\[\*\]\.([a-zA-Z_]\w*(?:\.[a-zA-Z_]\w*)*)\s*(===|!==|<=|>=|==|!=|<|>)\s*(.+)`)
	matches := re.FindStringSubmatch(expr)
	
	fmt.Printf("Regex matches: %d\n", len(matches))
	if len(matches) != 5 {
		fmt.Printf("❌ Expected 5 matches, got %d\n", len(matches))
		return
	}
	
	for i, match := range matches {
		fmt.Printf("  [%d]: %q\n", i, match)
	}
	
	arrayProperty := matches[1]     // e.g., "items"
	subProperty := matches[2]       // e.g., "product"
	operator := matches[3]          // e.g., "==="
	valueStr := strings.TrimSpace(matches[4]) // e.g., "'laptop'"
	
	fmt.Printf("\nExtracted values:\n")
	fmt.Printf("  arrayProperty: %q\n", arrayProperty)
	fmt.Printf("  subProperty: %q\n", subProperty)
	fmt.Printf("  operator: %q\n", operator)
	fmt.Printf("  valueStr: %q\n", valueStr)
	
	// Get the array from current object
	obj, ok := current.(map[string]interface{})
	if !ok {
		fmt.Printf("❌ Current is not a map[string]interface{}, it's %T\n", current)
		return
	}
	
	fmt.Printf("✓ Current is a map with keys: %v\n", func() []string {
		keys := make([]string, 0, len(obj))
		for k := range obj {
			keys = append(keys, k)
		}
		return keys
	}())
	
	arrayValue, exists := obj[arrayProperty]
	if !exists {
		fmt.Printf("❌ Property %q not found in object\n", arrayProperty)
		return
	}
	
	fmt.Printf("✓ Found property %q: %T\n", arrayProperty, arrayValue)
	
	arr, ok := arrayValue.([]interface{})
	if !ok {
		fmt.Printf("❌ Property %q is not an array, it's %T\n", arrayProperty, arrayValue)
		return
	}
	
	fmt.Printf("✓ Property %q is an array with %d elements\n", arrayProperty, len(arr))
	
	parsedValue := utils.ParseValue(valueStr)
	fmt.Printf("Parsed value: %v (%T)\n", parsedValue, parsedValue)
	
	// Check if any element in the array matches the condition
	fmt.Printf("\nChecking array elements:\n")
	for i, item := range arr {
		fmt.Printf("  Element %d: %T\n", i, item)
		if jsonBytes, err := json.MarshalIndent(item, "    ", "  "); err == nil {
			fmt.Printf("    Value: %s\n", string(jsonBytes))
		}
		
		subValue := utils.GetPropertyValue(item, subProperty)
		fmt.Printf("    %s = %v (%T)\n", subProperty, subValue, subValue)
		
		match := utils.CompareValues(subValue, operator, parsedValue)
		fmt.Printf("    %v %s %v = %t\n", subValue, operator, parsedValue, match)
		
		if match {
			fmt.Printf("✓ Found match in element %d\n", i)
			return
		}
	}
	
	fmt.Printf("❌ No matches found in any array element\n")
}

func debug_array_step_by_stepMain() {
	// Test data
	orderJSON := `{"id": "ORD001", "items": [{"product": "laptop", "price": 999.99}, {"product": "mouse", "price": 29.99}]}`
	
	var order map[string]interface{}
	json.Unmarshal([]byte(orderJSON), &order)
	
	expr := `.items[*].product === "laptop"`
	
	debugArrayWildcardFilter(expr, order)
}