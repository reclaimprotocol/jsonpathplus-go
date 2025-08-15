# JSONPath-Plus Go Usage Guide

This comprehensive guide covers all features and best practices for using the JSONPath-Plus Go library.

## 🚀 Quick Start

### Installation

```bash
go get github.com/reclaimprotocol/jsonpathplus-go
```

### Basic Usage

```go
package main

import (
    "fmt"
    jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
    // JSON string input
    jsonStr := `{
        "store": {
            "book": [
                {"title": "The Great Gatsby", "price": 12.99},
                {"title": "To Kill a Mockingbird", "price": 8.99}
            ]
        }
    }`
    
    // Query all book titles
    results, err := jp.Query("$.store.book[*].title", jsonStr)
    if err != nil {
        panic(err)
    }
    
    for _, result := range results {
        fmt.Printf("Title: %v\n", result.Value)
    }
}
```

## 📖 JSONPath Syntax Reference

### Basic Selectors

| Syntax | Description | Example |
|--------|-------------|---------|
| `$` | Root element | `$` |
| `.property` | Property access | `$.store` |
| `['property']` | Bracket notation | `$['store']` |
| `[index]` | Array index | `$.book[0]` |
| `[*]` | All elements | `$.book[*]` |

### Advanced Selectors

| Syntax | Description | Example |
|--------|-------------|---------|
| `..property` | Recursive descent | `$..price` |
| `..*` | All descendants | `$..*` |
| `[start:end]` | Array slice | `$.book[0:2]` |
| `[0,2,4]` | Union | `$.book[0,2,4]` |
| `[?(@.property)]` | Filter | `$.book[?(@.price < 10)]` |

## 🔍 Filter Expressions

### Property Filters

```go
// Filter by property existence
jp.Query("$.book[?(@.isbn)]", jsonStr)

// Filter by property value
jp.Query("$.book[?(@.category === 'fiction')]", jsonStr)

// Filter by property name
jp.Query("$..*[?(@property === 'price')]", jsonStr)
```

### Comparison Operators

```go
// Equality
jp.Query("$.book[?(@.price === 12.99)]", jsonStr)

// Inequality  
jp.Query("$.book[?(@.price !== 8.99)]", jsonStr)

// Comparison
jp.Query("$.book[?(@.price > 10)]", jsonStr)
jp.Query("$.book[?(@.price <= 15)]", jsonStr)
```

### Logical Operators

```go
// AND
jp.Query("$.book[?(@.price > 5 && @.price < 15)]", jsonStr)

// OR
jp.Query("$.book[?(@.category === 'fiction' || @.category === 'science')]", jsonStr)

// NOT
jp.Query("$.book[?(!@.outOfStock)]", jsonStr)
```

### Advanced Filter Features

```go
// Parent property access
jp.Query("$.store.book[?(@parent.bicycle)]", jsonStr)

// Parent property name
jp.Query("$.users.1[?(@parentProperty === '1')]", jsonStr)

// Path filters
jp.Query("$.users[?(@path === \"$['users'][1]\")]", jsonStr)

// Length property
jp.Query("$.data[?(@.length > 3)]", jsonStr) // Note: Throws error for null values
```

## 🎯 Real-World Examples

### E-commerce Store Query

```go
storeJSON := `{
    "store": {
        "book": [
            {
                "category": "reference",
                "author": "Nigel Rees",
                "title": "Sayings of the Century",
                "price": 8.95
            },
            {
                "category": "fiction",
                "author": "Evelyn Waugh", 
                "title": "Sword of Honour",
                "price": 12.99
            }
        ],
        "bicycle": {
            "color": "red",
            "price": 19.95
        }
    }
}`

// All authors
authors, _ := jp.Query("$..author", storeJSON)

// Books cheaper than $10
cheapBooks, _ := jp.Query("$..book[?(@.price < 10)]", storeJSON)

// All prices (books and bicycle)
allPrices, _ := jp.Query("$..*[?(@property === 'price')]", storeJSON)

// Fiction books
fictionBooks, _ := jp.Query("$..book[?(@.category === 'fiction')]", storeJSON)
```

### API Response Processing

```go
apiResponse := `{
    "users": [
        {"id": 1, "name": "Alice", "active": true, "age": 30},
        {"id": 2, "name": "Bob", "active": false, "age": 25},
        {"id": 3, "name": "Charlie", "active": true, "age": 35}
    ],
    "meta": {
        "total": 3,
        "page": 1
    }
}`

// Active users
activeUsers, _ := jp.Query("$.users[?(@.active === true)]", apiResponse)

// Users over 30
seniorUsers, _ := jp.Query("$.users[?(@.age > 30)]", apiResponse)

// User names only
userNames, _ := jp.Query("$.users[*].name", apiResponse)

// First two users
firstTwo, _ := jp.Query("$.users[0:2]", apiResponse)
```

## 🔧 Advanced Features

### String Position Tracking

```go
jsonStr := `{"id": 123, "name": "test"}`
results, _ := jp.Query("$.name", jsonStr)

for _, result := range results {
    fmt.Printf("Value: %v\n", result.Value)           // "test"
    fmt.Printf("Position: %d\n", result.OriginalIndex) // Character position in JSON
    fmt.Printf("Length: %d\n", result.Length)         // Length of the key/value
    fmt.Printf("Path: %s\n", result.Path)             // "$.name"
}
```

### Error Handling

```go
// Handle parsing errors
results, err := jp.Query("$.invalid.[syntax", jsonStr)
if err != nil {
    fmt.Printf("Parse error: %v\n", err)
}

// Handle runtime errors (like null.length)
results, err := jp.Query("$.data[?(@.length > 3)]", jsonWithNull)
if err != nil {
    fmt.Printf("Runtime error: %v\n", err) // "Cannot read properties of null (reading 'length')"
}
```

### Working with Parsed Data

```go
// Parse JSON first, then query
var data interface{}
json.Unmarshal([]byte(jsonStr), &data)

results, err := jp.QueryData("$.store.book[*]", data)
```

## 🚀 Performance Tips

### Best Practices

1. **Reuse JSONPath objects** for repeated queries:
```go
jsonpath, err := jp.New("$.store.book[*].price")
if err != nil {
    return err
}

// Reuse for multiple JSON documents
results1, _ := jsonpath.QueryString(json1)
results2, _ := jsonpath.QueryString(json2)
```

2. **Use specific paths** instead of broad recursive descent:
```go
// Faster
jp.Query("$.store.book[*].price", jsonStr)

// Slower
jp.Query("$..price", jsonStr)
```

3. **Filter early** in the query:
```go
// Better: Filter first, then process
jp.Query("$.users[?(@.active)].details", jsonStr)

// Worse: Get all, then filter
jp.Query("$.users[*].details", jsonStr) // Then filter in Go
```

### Benchmarking

```bash
# Run performance benchmarks
go test -bench=. ./...

# Profile memory usage
go test -bench=. -memprofile=mem.prof ./...

# CPU profiling
go test -bench=. -cpuprofile=cpu.prof ./...
```

## 🐛 Debugging

### Debug Mode

```go
// Enable debug logging (if available)
jp.SetDebugMode(true)

// Query with detailed output
results, err := jp.Query("$.complex[*].query", jsonStr)
```

### Common Issues

1. **Empty Results**
   - Check JSON structure matches query
   - Verify property names (case-sensitive)
   - Use bracket notation for special characters

2. **Parse Errors**
   - Validate JSONPath syntax
   - Check for proper quoting in filters
   - Verify bracket matching

3. **Runtime Errors**
   - Handle null values in filters
   - Check array bounds for slice operations
   - Validate filter expressions

### Testing Compatibility

```bash
# Run compatibility tests against JavaScript
cd tests && node compare.js

# Test specific query
echo '$.your.query' > tests/temp_query.txt
echo '{"your": {"data": "here"}}' > tests/temp_data.json
cd tests && node compare.js
```

## 📚 Additional Resources

- **[API Documentation](API.md)** - Complete API reference
- **[Examples](../cmd/)** - Real-world usage examples  
- **[Performance Benchmarks](BENCHMARKS.md)** - Performance analysis
- **[JSONPath Specification](https://goessner.net/articles/JsonPath/)** - Original JSONPath spec
- **[JSONPath-Plus](https://github.com/JSONPath-Plus/JSONPath)** - JavaScript reference implementation

## 🤝 Contributing

Found a compatibility issue or have suggestions? Please:

1. **Test against JavaScript**: Use `tests/compare.js` to verify behavior
2. **File an issue**: Include the query, data, and expected vs actual results
3. **Submit a PR**: Add test cases for new features or bug fixes

---

**🤖 Generated with [Claude Code](https://claude.ai/code)**