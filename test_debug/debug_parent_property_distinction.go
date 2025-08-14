package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	// Test data based on the failing tests
	jsonData := `{
		"company": {
			"departments": {
				"engineering": {
					"employees": [
						{"name": "Alice", "level": "senior"},
						{"name": "Bob", "level": "junior"}
					],
					"manager": "Eve"
				},
				"sales": {
					"employees": [
						{"name": "Charlie", "level": "senior"}
					],
					"manager": "Dave"
				}
			}
		},
		"store": {
			"book": [
				{"title": "Book1", "author": "Author1", "category": "fiction", "price": 8.95},
				{"title": "Book2", "author": "Author2", "category": "fiction", "price": 12.99}
			]
		}
	}`

	fmt.Println("=== Testing @property vs @parentProperty distinction ===")
	
	// Test 1: The failing case - array elements with @parentProperty
	fmt.Println("\n1. Array elements with @parentProperty: $.company.departments.engineering.employees[?(@parentProperty === 'employees')]")
	results1, err := jp.Query("$.company.departments.engineering.employees[?(@parentProperty === 'employees')]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d (expected: 2)\n", len(results1))
	for i, result := range results1 {
		fmt.Printf("  Result %d: %s -> %v\n", i, result.Path, result.Value)
	}
	
	// Test 2: Array elements with @property 
	fmt.Println("\n2. Array elements with @property: $.company.departments.engineering.employees[?(@property === '0')]")
	results2, err := jp.Query("$.company.departments.engineering.employees[?(@property === '0')]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d (expected: 1)\n", len(results2))
	for i, result := range results2 {
		fmt.Printf("  Result %d: %s -> %v\n", i, result.Path, result.Value)
	}
	
	// Test 3: The book case that was working - object properties with @parentProperty  
	fmt.Println("\n3. Object properties with @parentProperty: $.store.book.*[?(@parentProperty === '0')]")
	results3, err := jp.Query("$.store.book.*[?(@parentProperty === '0')]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d (expected: properties of first book)\n", len(results3))
	for i, result := range results3 {
		fmt.Printf("  Result %d: %s -> %v\n", i, result.Path, result.Value)
	}
	
	// Test 4: Check what we get without filters
	fmt.Println("\n4. All engineering employees: $.company.departments.engineering.employees")
	results4, err := jp.Query("$.company.departments.engineering.employees", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d\n", len(results4))
	for i, result := range results4 {
		fmt.Printf("  Result %d: %s -> %v\n", i, result.Path, result.Value)
	}
}