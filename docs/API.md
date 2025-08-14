# JSONPath-Plus Go API Documentation

## Overview

This is a comprehensive JSONPath implementation with extended JSONPath-Plus features, providing a modular architecture and advanced query capabilities.

## Core Functions

### `Query(path string, input interface{}) ([]Result, error)`

**Primary query function** - Executes a JSONPath query against JSON string or parsed data.

**Parameters:**
- `path` - JSONPath expression (e.g., "$.users[*].name")
- `input` - JSON string or parsed data structure

**Returns:**
- `[]Result` - Array of query results with metadata
- `error` - Error if query fails

**Example:**
```go
// With JSON string
results, err := jp.Query("$.users[*].name", `{"users":[{"name":"Alice"}]}`)

// With parsed data
var data interface{}
json.Unmarshal([]byte(jsonStr), &data)
results, err := jp.Query("$.users[*].name", data)
```

**Parameters:**
- `path` - JSONPath expression
- `data` - Parsed JSON data as interface{}

**Returns:**
- `[]Result` - Array of query results
- `error` - Error if query fails

## Result Structure

```go
type Result struct {
    Value            interface{} // The actual value
    Path             string      // JSONPath to this element  
    Parent           interface{} // Reference to parent object/array
    ParentProperty   string      // Property name or array index in parent
    Index            int         // Position in result set
    OriginalIndex    int         // Character position in original JSON string
    Length           int         // Length of the element in the JSON string
}
```

## Production Engine

### `NewEngine(config *Config) (*JSONPathEngine, error)`

Creates a production-ready JSONPath engine.

### Engine Methods

- `QueryData(path, data)` - Query parsed data
- `QueryDataWithContext(ctx, path, data)` - Query with context
- `GetMetrics()` - Performance metrics
  
- `Close()` - Cleanup resources

## Configuration

### `DefaultConfig() *Config`

Returns default configuration suitable for development.

### `ProductionConfig() *Config`

Returns configuration optimized for production use with security features enabled.

## Error Types

- `JSONPathError` - JSONPath parsing/execution errors
- `ValidationError` - Configuration validation errors
- `PathLengthError` - Path too long errors
- `RecursionLimitError` - Recursion depth exceeded