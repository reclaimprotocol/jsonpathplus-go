package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
	"github.com/reclaimprotocol/jsonpathplus-go/pkg/types"
)

func printAST(node *types.AstNode, indent string) {
	if node == nil {
		return
	}
	fmt.Printf("%s%s: %s\n", indent, node.Type, node.Value)
	for _, child := range node.Children {
		printAST(child, indent+"  ")
	}
}

func main() {
	fmt.Println("=== AST Analysis for Filter Issues ===")
	
	// Test case 1: Working case
	fmt.Println("\n1. Working case: $.store.book[*][?(@property === 'price')]")
	ast1, err := jp.Parse("$.store.book[*][?(@property === 'price')]")
	if err != nil {
		fmt.Printf("❌ Parse error: %v\n", err)
	} else {
		printAST(ast1, "")
	}
	
	// Test case 2: Broken case
	fmt.Println("\n2. Broken case: $..*[?(@property === 'price')]")
	ast2, err := jp.Parse("$..*[?(@property === 'price')]")
	if err != nil {
		fmt.Printf("❌ Parse error: %v\n", err)
	} else {
		printAST(ast2, "")
	}
	
	// Test case 3: Simple recursive descent
	fmt.Println("\n3. Simple recursive descent: $..*")
	ast3, err := jp.Parse("$..*")
	if err != nil {
		fmt.Printf("❌ Parse error: %v\n", err)
	} else {
		printAST(ast3, "")
	}
	
	// Test case 4: Just the filter
	fmt.Println("\n4. Just the filter: [?(@property === 'price')]")
	ast4, err := jp.Parse("$[?(@property === 'price')]")
	if err != nil {
		fmt.Printf("❌ Parse error: %v\n", err)
	} else {
		printAST(ast4, "")
	}
}