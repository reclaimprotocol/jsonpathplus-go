# Examples

## Basic Usage

```go
package main

import (
    "fmt"
    "log"
    jp "github.com/reclaimprotocol/jsonpathplus-go/v2"
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
        fmt.Printf("Title: %s at position %d (Length: %d)\n", 
            result.Value, result.Start, result.Length)
    }
}
```

## String Position Tracking (v2 Default)

```go
jsonStr := `{"users":[{"name":"Alice"},{"name":"Bob"}]}`

results, err := jp.Query("$.users[*].name", jsonStr)
// Result contains (v2 behavior):
// - Value: "test" 
// - Start: 12 (character position of "name" key)
// - End: 26 (end character position of "test" value)
// - Length: 14 (covers `"name": "test"`)
// - Path: "$.name"
```

## Production Engine

```go
package main

import (
    jp "github.com/reclaimprotocol/jsonpathplus-go/v2"
)

func main() {
    // Create production engine
    engine, err := jp.NewEngine()
    if err != nil {
        panic(err)
    }
    defer engine.Close()
    
    // Thread-safe queries
    results, err := engine.Query("$.users[*]", data)
    if err != nil {
        panic(err)
    }
}
```

## Advanced Queries

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