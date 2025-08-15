package jsonpathplus

import (
	"fmt"

	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func test_lengthMain() {
	data := `{"data":[42,"hello",true,null,{"key":"value"},[1,2,3]]}`

	query := `$.data[?(@.length > 3)]`
	fmt.Printf("Testing: %s\n", query)

	results, err := jp.Query(query, data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Count: %d\n", len(results))
		for i, result := range results {
			fmt.Printf("[%d] Value: %v\n", i, result.Value)
		}
	}
}
