package main

import (
	"fmt"
	"github.com/reclaimprotocol/jsonpathplus-go/internal/filters"
	"github.com/reclaimprotocol/jsonpathplus-go/pkg/types"
)

func main() {
	fmt.Println("=== Testing Filter Evaluation Directly ===")
	
	// Create a mock context like what would be created during recursive descent
	filterEval := filters.NewFilterEvaluator()
	
	// Test case 1: A value with property name "price"
	fmt.Println("\n1. Testing context with property name 'price'")
	ctx1 := types.NewContext(
		nil,           // root
		8.95,          // current (the price value)
		map[string]interface{}{"title": "Book1", "price": 8.95}, // parent (the book object)
		"price",       // parentProperty (this is what @property should return)
		"$.store.book[0].price", // path
		0,             // index
	)
	
	result1 := filterEval.EvaluateFilter("@property === 'price'", ctx1)
	fmt.Printf("Filter result for price property: %t (expected: true)\n", result1)
	fmt.Printf("Context property name: %s\n", ctx1.GetPropertyName())
	
	// Test case 2: A value with property name "title"
	fmt.Println("\n2. Testing context with property name 'title'")
	ctx2 := types.NewContext(
		nil,           // root
		"Book1",       // current (the title value)
		map[string]interface{}{"title": "Book1", "price": 8.95}, // parent (the book object)
		"title",       // parentProperty (this is what @property should return)
		"$.store.book[0].title", // path
		0,             // index
	)
	
	result2 := filterEval.EvaluateFilter("@property === 'price'", ctx2)
	fmt.Printf("Filter result for title property: %t (expected: false)\n", result2)
	fmt.Printf("Context property name: %s\n", ctx2.GetPropertyName())
	
	// Test case 3: Test the title filter
	fmt.Println("\n3. Testing title filter on title context")
	result3 := filterEval.EvaluateFilter("@property === 'title'", ctx2)
	fmt.Printf("Filter result for title === title: %t (expected: true)\n", result3)
	
	// Test case 4: Test the title filter on price context
	fmt.Println("\n4. Testing title filter on price context")
	result4 := filterEval.EvaluateFilter("@property === 'title'", ctx1)
	fmt.Printf("Filter result for title === title on price: %t (expected: false)\n", result4)
	
	// Test case 5: Test some other filters for comparison
	fmt.Println("\n5. Testing other filters")
	result5 := filterEval.EvaluateFilter("@ > 5", ctx1) // Value-based filter
	fmt.Printf("Filter @ > 5 on price 8.95: %t (expected: true)\n", result5)
}