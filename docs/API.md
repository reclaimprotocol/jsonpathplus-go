# JSONPath-Plus Go API Documentation

## Overview

This is a comprehensive JSONPath implementation with extended JSONPath-Plus features, providing a modular architecture and advanced query capabilities.

> [!IMPORTANT]
> **v2.0.0 Breaking Change**: Default indexing for object properties now includes the entire key-value pair as a single range.

## Core Functions

### `Query(path string, input interface{}) ([]Result, error)`

**Primary query function** - Executes a JSONPath query against JSON string or parsed data.

**Parameters:**
- `path` - JSONPath expression (e.g., "$.users[*].name")
- `input` - JSON string or parsed data structure

**Returns:**
- `[]Result` - Array of query results with metadata
- `error` - Error if query fails

### `QueryWithOptions(path string, input interface{}, options *Options) ([]Result, error)`

Executes a JSONPath query with custom options.

## Result Structure

```go
type Result struct {
    Value            interface{} // The actual value
    Path             string      // JSONPath to this element  
    Parent           interface{} // Reference to parent object/array
    ParentProperty   string      // Property name or array index in parent
    Index            int         // Position in result set
    Start            int         // Starting character position (key)
    End              int         // Ending character position (value)
    Length           int         // Total length (key to value)
    OriginalIndex    int         // Same as Start (for backward compatibility)
}
```

## Production Engine

### `NewEngine() (*JSONPathEngine, error)`

Creates a production-ready JSONPath engine.

### Engine Methods

- `Query(path, input)` - Query JSON string or parsed data
- `QueryWithOptions(path, input, options)` - Query with options
- `Close()` - Cleanup resources (no-op)

## Configuration

### `Options`

```go
type Options struct {
    Root interface{} // Root object for $ references in filters
}
```

## Error Types

- `JSONPathError` - JSONPath parsing/execution errors
- `ValidationError` - Configuration validation errors
- `PathLengthError` - Path too long errors
- `RecursionLimitError` - Recursion depth exceeded