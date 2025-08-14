package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"arr": [1, 2, 3, 4, 5]
	}`

	// Test simple chaining cases
	fmt.Println("=== Testing simple array operations ===")

	// Test slice
	results1, err := jp.Query("$.arr[0:3]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("First 3 elements: %d results\n", len(results1))
	for i, r := range results1 {
		fmt.Printf("  [%d] %v\n", i, r.Value)
	}

	// Test simple double bracket
	fmt.Println("\n=== Testing double bracket syntax ===")
	results2, err := jp.Query("$.arr[0][0]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("arr[0][0]: %d results\n", len(results2))
		for i, r := range results2 {
			fmt.Printf("  [%d] %v\n", i, r.Value)
		}
	}

	// Test more meaningful case
	jsonData2 := `{
		"data": [
			[10, 20, 30],
			[40, 50, 60],
			[70, 80, 90]
		]
	}`

	fmt.Println("\n=== Testing with nested arrays ===")
	results3, err := jp.Query("$.data[0:2]", jsonData2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("First 2 arrays: %d results\n", len(results3))
	for i, r := range results3 {
		fmt.Printf("  [%d] %v\n", i, r.Value)
	}

	// Test chained operation on nested arrays
	fmt.Println("\n=== Testing chained operation ===")
	results4, err := jp.Query("$.data[0:2][0]", jsonData2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("First element of first 2 arrays: %d results\n", len(results4))
	for i, r := range results4 {
		fmt.Printf("  [%d] %v\n", i, r.Value)
	}
}
