package main

import (
	"fmt"
	"encoding/json"
	
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func printAST(node interface{}, indent string) {
	bytes, _ := json.MarshalIndent(node, indent, "  ")
	fmt.Println(string(bytes))
}

func main() {
	query := `$.orders[?(@.items[?(@.product === 'laptop')])]`
	fmt.Printf("Parsing: %s\n", query)
	
	jsonpath, err := jp.New(query)
	if err != nil {
		fmt.Printf("Parse Error: %v\n", err)
	} else {
		fmt.Println("AST:")
		printAST(jsonpath.AST(), "")
	}
}