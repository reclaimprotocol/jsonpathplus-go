package main

import (
	"fmt"
	
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	data := `{"store":{"book":[{"title":"Book1"},{"title":"Book2"}],"bicycle":{"color":"red"}}}`
	
	query := "$.store.book[?(@parent.bicycle)]"
	fmt.Printf("Testing: %s\n", query)
	
	results, err := jp.Query(query, data)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	
	fmt.Printf("Count: %d\n", len(results))
	for i, result := range results {
		fmt.Printf("[%d] Value: %v\n", i, result.Value)
	}
}