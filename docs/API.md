# API Documentation

## Core Functions

### `Query(path string, jsonStr string) ([]Result, error)`

Executes a JSONPath query with string character position tracking.

**Parameters:**
- `path` - JSONPath expression (e.g., "$.users[*].name")
- `jsonStr` - JSON string to query

**Returns:**
- `[]Result` - Array of query results with position information
- `error` - Error if query fails

### `QueryData(path string, data interface{}) ([]Result, error)`

Executes a JSONPath query on parsed data (legacy support).

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
- `GetCacheStats()` - Cache statistics  
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