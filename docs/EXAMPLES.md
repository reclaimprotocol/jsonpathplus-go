# Examples

## Basic Usage

```go
package main

import (
    "fmt"
    "log"
    jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
    jsonStr := `{
        "store": {
            "book": [
                {"title": "Book 1", "price": 10.50},
                {"title": "Book 2", "price": 15.99}
            ]
        }
    }`
    
    // Get all book titles
    results, err := jp.Query("$.store.book[*].title", jsonStr)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, result := range results {
        fmt.Printf("Title: %s at position %d\n", 
            result.Value, result.OriginalIndex)
    }
}
```

## String Position Tracking

```go
jsonStr := `{"users":[{"name":"Alice"},{"name":"Bob"}]}`

results, err := jp.Query("$.users[*].name", jsonStr)
for _, result := range results {
    fmt.Printf("Name: %s\n", result.Value)
    fmt.Printf("  Character position: %d\n", result.OriginalIndex) 
    fmt.Printf("  Length in JSON: %d\n", result.Length)
    fmt.Printf("  JSONPath: %s\n", result.Path)
}

// Output:
// Name: Alice
//   Character position: 11
//   Length in JSON: 6
//   JSONPath: $.users[0].name
// Name: Bob  
//   Character position: 32
//   Length in JSON: 5
//   JSONPath: $.users[1].name
```

## Production Engine

```go
package main

import (
    jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
    // Create production engine
    config := jp.ProductionConfig()
    config.EnableMetrics = true
    
    engine, err := jp.NewEngine(config)
    if err != nil {
        panic(err)
    }
    defer engine.Close()
    
    // Thread-safe queries
    results, err := engine.QueryData("$.users[*]", data)
    if err != nil {
        panic(err)
    }
    
    // Check performance metrics
    metrics := engine.GetMetrics()
    fmt.Printf("Queries executed: %d\n", metrics.QueriesExecuted)
    fmt.Printf("Average time: %v\n", metrics.AverageExecutionTime)
    
}
```

## Complex Queries

```go
jsonStr := `{
    "products": [
        {"name": "Laptop", "price": 999.99, "category": "electronics"},
        {"name": "Book", "price": 19.99, "category": "books"},
        {"name": "Phone", "price": 599.99, "category": "electronics"}
    ]
}`

// Filter by price
results, _ := jp.Query("$.products[?(@.price > 500)]", jsonStr)

// Recursive search
results, _ := jp.Query("$..price", jsonStr)

// Array slicing
results, _ := jp.Query("$.products[0:2]", jsonStr)

// Union operator
results, _ := jp.Query("$.products[*]['name','price']", jsonStr)
```

## Error Handling

```go
results, err := jp.Query("$.invalid..path", jsonStr)
if err != nil {
    if jpErr, ok := err.(*jp.JSONPathError); ok {
        fmt.Printf("JSONPath error: %s at position %d\n", 
            jpErr.Message, jpErr.Position)
    }
}
```