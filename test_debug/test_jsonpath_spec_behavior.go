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
				},
				{
					"category": "action",
					"title": "Book 3"
				}
			]
		}
	}`

	fmt.Println("=== JSONPath Specification Analysis ===")
	fmt.Println()

	fmt.Println("According to JSONPath spec, these should be equivalent to XPath:")
	fmt.Println()

	fmt.Println("1. $..book (recursive descent to find 'book')")
	results1, err := jp.Query("$..book", jsonData)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Results: %d\n", len(results1))
		for i, r := range results1 {
			fmt.Printf("   [%d] Array with %d books (path: %s)\n", i, len(r.Value.([]interface{})), r.Path)
		}
	}

	fmt.Println()
	fmt.Println("2. $..book.* (all properties of books found recursively)")
	fmt.Println("   XPath equivalent: //book/*")
	fmt.Println("   Should return individual property values, not book objects")
	results2, err := jp.Query("$..book.*", jsonData)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Current Results: %d\n", len(results2))
		for i, r := range results2 {
			fmt.Printf("   [%d] %v (type: %T, path: %s)\n", i, r.Value, r.Value, r.Path)
		}
		fmt.Println("   ^ Currently returns book objects, but should return individual properties")
	}

	fmt.Println()
	fmt.Println("3. What $..book.* SHOULD return (based on XPath //book/*):")
	fmt.Println("   Should be equivalent to: $.store.book[*].*")
	results3, err := jp.Query("$.store.book[*].*", jsonData)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Expected Results: %d\n", len(results3))
		for i, r := range results3 {
			fmt.Printf("   [%d] '%v' (type: %T, path: %s)\n", i, r.Value, r.Value, r.Path)
		}
		fmt.Println("   ^ These are individual property values")
	}

	fmt.Println()
	fmt.Println("=== Issue Analysis ===")
	fmt.Println("The current implementation of $..book.* appears to have an issue:")
	fmt.Println("- It returns book objects instead of individual properties")
	fmt.Println("- This makes the original query fail because there are no 'category' properties")
	fmt.Println("  at the book object level that can be filtered")
	fmt.Println()
	fmt.Println("To make the original query work as specified:")
	fmt.Println("$..book.*[?(@property === \"category\" && @.match(/TION$/i))]")
	fmt.Println("The $..book.* part needs to return individual property values,")
	fmt.Println("not book objects.")
}
