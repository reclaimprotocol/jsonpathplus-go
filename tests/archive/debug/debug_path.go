package jsonpathplus

import (
	"fmt"
	"strings"
)

// extractParentPropertyFromPath extracts the parent property name from a JSONPath
func extractParentPropertyFromPath(path string) string {
	fmt.Printf("  Extracting from path: '%s'\n", path)
	
	// For @parentProperty, we want the property that led to the parent of the current item
	// Examples:
	// "$.users.1['name']" -> "users" (parent of parent of 'name')
	// "$.store.book[0]['title']" -> "store" (parent of parent of 'title')  
	// "$.users['1']" -> "" (filtering users object directly)
	
	if strings.Contains(path, "['") {
		fmt.Printf("  Found bracket notation\n")
		// Path has bracket notation for current property like $.users.1['name']
		// Remove the bracket part to get the parent path
		lastBracket := strings.LastIndex(path, "['")
		fmt.Printf("  Last bracket at: %d\n", lastBracket)
		if lastBracket > 0 {
			parentPath := path[:lastBracket] // "$.users.1"
			fmt.Printf("  Parent path: '%s'\n", parentPath)
			
			// Now get the parent of this parent
			// Find the second-to-last property
			if strings.Contains(parentPath, ".") {
				parts := strings.Split(parentPath, ".")
				fmt.Printf("  Parts: %v\n", parts)
				if len(parts) >= 3 { // $, intermediate, property
					// Get the second-to-last part (parent of parent)
					result := parts[len(parts)-2]
					fmt.Printf("  Extracted: '%s'\n", result)
					return result
				}
			}
		}
	}
	
	fmt.Printf("  No extraction possible\n")
	return "" // Default for cases we can't parse
}

func debug_pathMain() {
	testPaths := []string{
		"$.users.1['name']",
		"$.store.book[0]['title']", 
		"$.users['1']",
		"$.users.1",
	}
	
	for _, path := range testPaths {
		fmt.Printf("Path: %s\n", path)
		result := extractParentPropertyFromPath(path)
		fmt.Printf("Result: '%s'\n\n", result)
	}
}