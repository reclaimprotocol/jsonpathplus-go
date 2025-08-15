package jsonpathplus

import (
	"fmt"

	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func test_wildcard_filterMain() {
	data := `{
		"orders": [
			{"id": 1, "items": [{"product": "laptop", "qty": 1}, {"product": "mouse", "qty": 2}]},
			{"id": 2, "items": [{"product": "phone", "qty": 1}]},
			{"id": 3, "items": [{"product": "laptop", "qty": 1}]}
		]
	}`

	query := `$.orders[?(@.items[*].product === 'laptop')]`
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
