package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{
					"category": "reference",
					"title": "Book 1"
				},
				{
					"category": "fiction", 
					"title": "Book 2"
				}
			]
		}
	}`
	
	fmt.Println("=== Testing $..book.* (all book properties) ===")
	results1, err := jp.Query("$..book.*", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found %d properties\n", len(results1))
		for i, r := range results1 {
			fmt.Printf("  [%d] %v (path: %s)\n", i, r.Value, r.Path)
		}
	}
	
	fmt.Println("\n=== Testing $.store.book[0] access ===")
	results2, err := jp.Query("$.store.book[0]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found %d books\n", len(results2))
		for i, r := range results2 {
			fmt.Printf("  [%d] %v (path: %s)\n", i, r.Value, r.Path)
		}
	}
	
	fmt.Println("\n=== Testing $.store.book[0].* (first book properties) ===")
	results3, err := jp.Query("$.store.book[0].*", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found %d properties\n", len(results3))
		for i, r := range results3 {
			fmt.Printf("  [%d] %v (path: %s)\n", i, r.Value, r.Path)
		}
	}
	
	fmt.Println("\n=== Testing @property filter on first book ===")
	results4, err := jp.Query("$.store.book[0].*[?(@property === \"category\")]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found %d category properties\n", len(results4))
		for i, r := range results4 {
			fmt.Printf("  [%d] %v (path: %s)\n", i, r.Value, r.Path)
		}
	}
	
	fmt.Println("\n=== Testing simple regex match ===")
	results5, err := jp.Query("$.store.book[0].category[?(@.match(/reference/))]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found %d matching categories\n", len(results5))
		for i, r := range results5 {
			fmt.Printf("  [%d] %v\n", i, r.Value)
		}
	}
}