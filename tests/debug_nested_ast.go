package main

import (
	"encoding/json"
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func printAST(node interface{}, indent string) {
	bytes, _ := json.MarshalIndent(node, indent, "  ")
	fmt.Println(string(bytes))
}

func main() {
	// Test simple filter (works)
	simpleQuery := `$.orders[?(@.id === "ORD001")]`
	fmt.Printf("=== SIMPLE FILTER AST ===\n")
	fmt.Printf("Query: %s\n", simpleQuery)
	
	simpleJsonPath, err := jp.New(simpleQuery)
	if err != nil {
		fmt.Printf("Parse Error: %v\n", err)
	} else {
		fmt.Println("AST:")
		printAST(simpleJsonPath.AST(), "")
	}

	// Test nested filter (fails)
	nestedQuery := `$.orders[?(@.items[?(@.product === "laptop")])]`
	fmt.Printf("\n=== NESTED FILTER AST ===\n")
	fmt.Printf("Query: %s\n", nestedQuery)
	
	nestedJsonPath, err := jp.New(nestedQuery)
	if err != nil {
		fmt.Printf("Parse Error: %v\n", err)
	} else {
		fmt.Println("AST:")
		printAST(nestedJsonPath.AST(), "")
	}
	
	// Let's also test what the filter string looks like
	fmt.Printf("\n=== FILTER STRINGS ===\n")
	if nestedJsonPath != nil && nestedJsonPath.AST() != nil {
		fmt.Printf("Root AST Type: %v\n", nestedJsonPath.AST().Type)
		if len(nestedJsonPath.AST().Children) > 0 {
			fmt.Printf("First Child Type: %v\n", nestedJsonPath.AST().Children[0].Type)
			if len(nestedJsonPath.AST().Children[0].Children) > 0 {
				filterNode := nestedJsonPath.AST().Children[0].Children[0]
				fmt.Printf("Filter Node Type: %v\n", filterNode.Type)
				fmt.Printf("Filter Value: %q\n", filterNode.Value)
			}
		}
	}
}