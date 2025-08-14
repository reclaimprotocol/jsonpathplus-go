package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{"category": "fiction"},
				{"category": "action"}
			]
		}
	}`
	
	// Create a debug version to understand paths during evaluation
	fmt.Println("=== Path Analysis ===")
	fmt.Println("We need to understand what paths look like during wildcard evaluation")
	
	// Test simple cases first
	fmt.Println()
	fmt.Println("1. $.store.book (should end with 'book')")
	results1, err := jp.Query("$.store.book", jsonData)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		for _, r := range results1 {
			fmt.Printf("   Path: '%s' (ends with ']': %t)\n", r.Path, len(r.Path) > 0 && r.Path[len(r.Path)-1] == ']')
		}
	}
	
	fmt.Println()
	fmt.Println("2. $.store.book[0] (should end with ']')")
	results2, err := jp.Query("$.store.book[0]", jsonData)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		for _, r := range results2 {
			fmt.Printf("   Path: '%s' (ends with ']': %t)\n", r.Path, len(r.Path) > 0 && r.Path[len(r.Path)-1] == ']')
		}
	}
	
	fmt.Println()
	fmt.Println("The issue might be that by the time we reach the wildcard evaluation,")
	fmt.Println("we've already lost the context of whether it came from .* or [*]")
}