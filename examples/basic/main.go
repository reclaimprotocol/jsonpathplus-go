package main

import (
	"encoding/json"
	"fmt"
	"log"
	
	jp "jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{
					"category": "reference",
					"author": "Nigel Rees",
					"title": "Sayings of the Century",
					"price": 8.95
				},
				{
					"category": "fiction",
					"author": "Evelyn Waugh",
					"title": "Sword of Honour",
					"price": 12.99
				},
				{
					"category": "fiction",
					"author": "Herman Melville",
					"title": "Moby Dick",
					"isbn": "0-553-21311-3",
					"price": 8.99
				},
				{
					"category": "fiction",
					"author": "J. R. R. Tolkien",
					"title": "The Lord of the Rings",
					"isbn": "0-395-19395-8",
					"price": 22.99
				}
			],
			"bicycle": {
				"color": "red",
				"price": 19.95
			}
		},
		"expensive": 10
	}`
	
	fmt.Println("JSONPath Examples with String Position Tracking")
	fmt.Println("================================================")
	
	example1(jsonData)
	example2(jsonData)
	example3(jsonData)
	example4(jsonData)
	example5(jsonData)
	example6()
}

func example1(jsonStr string) {
	fmt.Println("\n1. Get all book authors:")
	results, err := jp.Query("$.store.book[*].author", jsonStr)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	
	for _, result := range results {
		fmt.Printf("   Value: %v, Path: %s, Index: %d, OriginalIndex: %d\n", 
			result.Value, result.Path, result.Index, result.OriginalIndex)
	}
}

func example2(jsonStr string) {
	fmt.Println("\n2. Get all prices recursively:")
	results, err := jp.Query("$..price", jsonStr)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	
	for _, result := range results {
		fmt.Printf("   Value: %.2f, Path: %s, OriginalIndex: %d\n", 
			result.Value, result.Path, result.OriginalIndex)
	}
}

func example3(jsonStr string) {
	fmt.Println("\n3. Filter books with price < 10:")
	results, err := jp.Query("$.store.book[?(@.price < 10)]", jsonStr)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	
	for _, result := range results {
		if book, ok := result.Value.(map[string]interface{}); ok {
			fmt.Printf("   Title: %v, Price: %v, OriginalIndex: %d\n", 
				book["title"], book["price"], result.OriginalIndex)
		}
	}
}

func example4(jsonStr string) {
	fmt.Println("\n4. Get books using slice notation [0:2]:")
	results, err := jp.Query("$.store.book[0:2]", jsonStr)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	
	for _, result := range results {
		if book, ok := result.Value.(map[string]interface{}); ok {
			fmt.Printf("   Title: %v, OriginalIndex: %d\n", 
				book["title"], result.OriginalIndex)
		}
	}
}

func example5(jsonStr string) {
	fmt.Println("\n5. Get last book using negative index:")
	results, err := jp.Query("$.store.book[-1]", jsonStr)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	
	for _, result := range results {
		if book, ok := result.Value.(map[string]interface{}); ok {
			fmt.Printf("   Title: %v, Author: %v, OriginalIndex: %d\n", 
				book["title"], book["author"], result.OriginalIndex)
		}
	}
}

func example6() {
	fmt.Println("\n6. Working with arrays and preserving indices:")
	
	arrayData := `{
		"numbers": [10, 20, 30, 40, 50],
		"nested": {
			"items": ["a", "b", "c", "d"]
		}
	}`
	
	results, err := jp.Query("$.numbers[*]", arrayData)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	
	fmt.Println("   Numbers array:")
	for _, result := range results {
		fmt.Printf("     Value: %v, OriginalIndex: %d, Path: %s\n", 
			result.Value, result.OriginalIndex, result.Path)
	}
	
	results, err = jp.Query("$.nested.items[*]", arrayData)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	
	fmt.Println("   Nested items array:")
	for _, result := range results {
		fmt.Printf("     Value: %v, OriginalIndex: %d, Path: %s\n", 
			result.Value, result.OriginalIndex, result.Path)
	}
}

func prettyPrint(v interface{}) string {
	bytes, _ := json.MarshalIndent(v, "", "  ")
	return string(bytes)
}