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

## 🎯 Real-World Examples

### String Position Tracking (v2 Default)

One of the most powerful features is getting the exact character ranges in the original JSON.

```go
jsonStr := `{"id": 123, "name": "test"}`
results, _ := jp.Query("$.name", jsonStr)

for _, result := range results {
    fmt.Printf("Value: %v\n", result.Value)      // "test"
    fmt.Printf("Start: %d\n", result.Start)      // Position of "name" key
    fmt.Printf("End: %d\n", result.End)          // End position of "test" value
    fmt.Printf("Length: %d\n", result.Length)    // Total range length
    fmt.Printf("Path: %s\n", result.Path)        // "$.name"
}
```

> [!NOTE]
> In v2.0.0, the indices for object properties cover the **entire range** from the start of the key to the end of the value. For array elements, they cover the element itself.

## 🚀 Performance Tips

1. **Reuse engines**: `jp.NewEngine()` results are thread-safe and can be reused.
2. **Specific paths**: Avoid `$..` if you know the exact depth.
3. **Internal imports**: When importing internal packages in your own projects (though usually not recommended), ensure you track the version correctly.

## 📚 Additional Resources

- **[API Documentation](API.md)** - Complete API reference
- **[Examples](../cmd/)** - Real-world usage examples  
- **[JSONPath Specification](https://goessner.net/articles/JsonPath/)** - Original JSONPath spec
- **[JSONPath-Plus](https://github.com/JSONPath-Plus/JSONPath)** - JavaScript reference implementation

---

**🤖 Generated with [Claude Code](https://claude.ai/code)**