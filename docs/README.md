# JSONPath Plus Go

A production-ready Go implementation of JSONPath with **enhanced string character position tracking**. This library provides all the functionality of [JSONPath-Plus](https://github.com/JSONPath-Plus/JSONPath) plus the unique ability to return the **exact character ranges** of elements in the original JSON string.

## 🚀 **Key Features**

### **Enhanced String Character Position Tracking** ⭐
- **Returns character ranges in original JSON string by default**
- Object properties: `{"id": 123, "name": "test"}` → `$.id` returns range `[1, 10)` (covers `"id": 123`)
- Array elements: `["a","b","c"]` → `$[1]` returns position `5` (covers `"b"`)
- Preserves whitespace and formatting perfectly

### **Complete JSONPath Support**
- Root (`$`) and current (`@`) operators
- Dot notation (`.property`) and bracket notation (`['property']`)
- Wildcards (`*`) and recursive descent (`..`)
- Array slicing (`[start:end:step]`) with negative indices
- Filter expressions (`[?(@.price < 10)]`) with complex operators
- Union operator (`['prop1','prop2']`)
- **No transformation** or **whitespace changes**

### **Production Features**
- Minimal dependencies (only standard library)
- Thread-safe concurrent operations
- LRU caching for compiled expressions
- Comprehensive error handling
- Security validation and rate limiting

## 📦 **Installation**

```bash
go get github.com/reclaimprotocol/jsonpathplus-go
```

## 🎯 **Quick Start**

### **Basic Query (with positions)**

```go
package main

import (
    "fmt"
    jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
    jsonStr := `{"id":123,"name":"test","active":true}`
    
    // Query with string position tracking (default v2)
    results, _ := jp.Query("$.name", jsonStr)
    
    for _, result := range results {
        fmt.Printf("Property '%s' found at range [%d:%d], length %d\n", 
            "name", result.Start, result.End, result.Length)
    }
}
```

## 📋 **API Reference**

### **Core Functions**

```go
// Standard JSONPath query (with string indices)
func Query(path string, data interface{}) ([]Result, error)

// Query with options
func QueryWithOptions(path string, data interface{}, options *Options) ([]Result, error)
```

## 🧪 **Testing**

```bash
# Run all tests
go test ./...

# Run benchmarks
go test -bench=.
```

## 📄 **License**

MIT License - see LICENSE file for details.

---

## 🌟 **What Makes This Special**

In v2.0.0, this library is the **only Go JSONPath implementation** that providing **exact character ranges** for object properties (key-to-value) by default. 

**Traditional approach:**  
Returns just the value or a simple index.

**Our approach:**  
`Start=1, End=10, Length=9` (covers `"id": 123` precisely)

Perfect for applications requiring **precise JSON element location tracking**! 🎯