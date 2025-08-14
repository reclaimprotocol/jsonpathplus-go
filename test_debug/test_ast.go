package main

import (
	"fmt"
	"github.com/reclaimprotocol/jsonpathplus-go/internal/parser"
	"github.com/reclaimprotocol/jsonpathplus-go/pkg/types"
)

func printAST(node *types.AstNode, depth int) {
	if node == nil {
		return
	}
	
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	
	fmt.Printf("%sType: %s, Value: %s\n", indent, node.Type, node.Value)
	
	for _, child := range node.Children {
		printAST(child, depth+1)
	}
}

func main() {
	p := parser.NewParser()
	
	paths := []string{"$..book", "$..book[*]", "$.store.book[*]"}
	
	for _, path := range paths {
		fmt.Printf("=== Parsing: %s ===\n", path)
		ast, err := p.Parse(path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			printAST(ast, 0)
		}
		fmt.Println()
	}
}