package main

import (
	"fmt"
	
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	data := `{"users":{"1":{"name":"Alice"},"2":{"name":"Bob"}}}`
	
	// First, see what paths we generate
	fmt.Println("=== Getting paths for $.users[*] ===")
	results, _ := jp.Query("$.users[*]", data)
	for i, result := range results {
		fmt.Printf("[%d] Path: %s, Value: %v\n", i, result.Path, result.Value)
	}
	
	// Now test the path filter
	queries := []string{
		`$.users[?(@path === "$['users']['1']")]`,
		`$.users[?(@path === "$.users['1']")]`,
		`$.users[?(@path === "$.users.1")]`,
	}
	
	for _, query := range queries {
		fmt.Printf("\n=== Testing: %s ===\n", query)
		results, err := jp.Query(query, data)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}
		
		fmt.Printf("Count: %d\n", len(results))
		for i, result := range results {
			fmt.Printf("[%d] Path: %s\n", i, result.Path)
		}
	}
}