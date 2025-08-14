package main

import (
	"encoding/json"
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	data := `{"store":{"book":[{"title":"Book0"},{"title":"Book1"},{"title":"Book2"}]}}`
	
	fmt.Println("=== Testing Go Implementation ===")
	
	fmt.Println("\n1. Testing $..book")
	results1, _ := jp.Query("$..book", data)
	fmt.Printf("Count: %d\n", len(results1))
	for i, r := range results1 {
		fmt.Printf("  [%d] %s: %T\n", i, r.Path, r.Value)
	}
	
	fmt.Println("\n2. Testing $..book[*]")
	results2, _ := jp.Query("$..book[*]", data)
	fmt.Printf("Count: %d\n", len(results2))
	for i, r := range results2 {
		valueStr, _ := json.Marshal(r.Value)
		fmt.Printf("  [%d] %s: %s\n", i, r.Path, string(valueStr))
	}
	
	fmt.Println("\n3. Testing $..book[?(@property !== 0)]")
	results3, _ := jp.Query("$..book[?(@property !== 0)]", data)
	fmt.Printf("Count: %d\n", len(results3))
	for i, r := range results3 {
		valueStr, _ := json.Marshal(r.Value)
		fmt.Printf("  [%d] %s: %s\n", i, r.Path, string(valueStr))
	}
}