/*
Package jsonpathplus provides a production-ready JSONPath implementation for Go
with comprehensive features including security, metrics, and resource limits.

# Overview

JSONPath Plus Go is a robust, high-performance JSONPath library that implements
the full JSONPath specification with additional production-ready features:

- Full JSONPath syntax support ($.., *, filters, slicing, etc.)
- Original index preservation
- Thread-safe concurrent access
- Security validation and sandboxing
- Performance metrics and monitoring
- Resource limits (memory, time, recursion)
- Structured logging
- Context-aware cancellation
- Rate limiting

# Quick Start

	package main

	import (
		"fmt"
		"log"
		jp "github.com/reclaimprotocol/jsonpathplus-go"
	)

	func main() {
		// Create engine with default configuration
		engine, err := jp.NewEngine(jp.DefaultConfig())
		if err != nil {
			log.Fatal(err)
		}
		defer engine.Close()

		// Parse JSON data
		data, err := jp.JSONParse(`{
			"users": [
				{"name": "John", "age": 30},
				{"name": "Jane", "age": 25}
			]
		}`)
		if err != nil {
			log.Fatal(err)
		}

		// Query with JSONPath
		results, err := engine.Query("$.users[?(@.age > 25)].name", data)
		if err != nil {
			log.Fatal(err)
		}

		for _, result := range results {
			fmt.Printf("Name: %s, Path: %s\n", result.Value, result.Path)
		}
	}

# Production Configuration

For production environments, use the production configuration:

	config := jp.ProductionConfig()
	config.MaxResultCount = 500
	config.Timeout = 3 * time.Second
	config.EnableMetrics = true

	engine, err := jp.NewEngine(config)
	if err != nil {
		log.Fatal(err)
	}

# Security Features

The library includes comprehensive security features:

	// Create secure configuration
	config := jp.DefaultConfig()
	config.StrictMode = true
	config.AllowUnsafeOperations = false

	// Add security validator
	securityConfig := jp.DefaultSecurityConfig()
	validator := jp.NewSecurityValidator(securityConfig)

	// Validate paths before execution
	if err := validator.ValidatePath(jsonPath); err != nil {
		log.Printf("Unsafe path: %v", err)
		return
	}

# Performance Monitoring

Monitor performance with built-in metrics:

	engine, _ := jp.NewEngine(jp.DefaultConfig())

	// Execute queries...

	metrics := engine.GetMetrics()
	fmt.Printf("Queries: %d, Avg time: %v\n",
		metrics.QueriesExecuted,
		metrics.AverageExecutionTime)

# JSONPath Syntax Support

The library supports the full JSONPath specification:

	$                    - Root object
	@                    - Current object (in filters)
	.property            - Child property
	['property']         - Bracket notation
	[index]              - Array index
	[start:end]          - Array slice
	[start:end:step]     - Array slice with step
	*                    - Wildcard (all properties/elements)
	..                   - Recursive descent
	[?(@.price < 10)]    - Filter expression
	['prop1','prop2']    - Union of properties
	[-1]                 - Negative array index

# Filter Expressions

Comprehensive filter expression support:

	// Comparison operators
	$.store.book[?(@.price < 10)]           // Less than
	$.store.book[?(@.price <= 10)]          // Less than or equal
	$.store.book[?(@.price > 10)]           // Greater than
	$.store.book[?(@.price >= 10)]          // Greater than or equal
	$.store.book[?(@.price == 10)]          // Equal
	$.store.book[?(@.price != 10)]          // Not equal

	// Property existence
	$.store.book[?(@.isbn)]                 // Has isbn property

	// String matching
	$.store.book[?(@.category == 'fiction')]

# Error Handling

The library provides detailed error types for better error handling:

	results, err := engine.Query(path, data)
	if err != nil {
		var jsonPathErr *jp.JSONPathError
		if errors.As(err, &jsonPathErr) {
			switch jsonPathErr.Type {
			case jp.ErrInvalidPath:
				// Handle invalid path
			case jp.ErrEvaluationError:
				// Handle evaluation error
			case jp.ErrRecursionLimit:
				// Handle recursion limit
			}
		}
	}

# Thread Safety

All operations are thread-safe and support concurrent access:

	engine, _ := jp.NewEngine(jp.DefaultConfig())

	// Safe to use from multiple goroutines
	go func() {
		results, _ := engine.Query("$.users[*].name", data)
		// Process results...
	}()

	go func() {
		results, _ := engine.Query("$.products[*].price", data)
		// Process results...
	}()

# Resource Limits

Configure resource limits to prevent abuse:

	config := jp.DefaultConfig()
	config.MaxRecursionDepth = 50          // Limit recursion
	config.MaxResultCount = 1000           // Limit results
	config.MaxMemoryUsage = 50 * 1024 * 1024  // 50MB limit
	config.Timeout = 5 * time.Second       // 5 second timeout

	engine, _ := jp.NewEngine(config)

# Custom Logging

Integrate with your logging system:

	type MyLogger struct {
		logger *logrus.Logger
	}

	func (l *MyLogger) Debug(msg string, fields ...jp.Field) {
		l.logger.WithFields(convertFields(fields)).Debug(msg)
	}
	// ... implement other methods

	engine, _ := jp.NewEngine(jp.DefaultConfig())
	engine.SetLogger(&MyLogger{logger: logrus.New()})

# Rate Limiting

Implement rate limiting for API endpoints:

	limiter := jp.NewRateLimiter(100, time.Minute) // 100 requests per minute

	func handleJSONPathQuery(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		if !limiter.Allow(clientIP) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Process query...
	}

# Best Practices

1. Use production configuration for production deployments
2. Enable metrics and monitoring
3. Set appropriate resource limits
4. Validate paths in strict mode
5. Use context for cancellation in long-running operations
6. Monitor performance metrics and optimize accordingly
7. Implement proper error handling
8. Use rate limiting for public APIs
9. Regular security audits of allowed expressions
10. Monitor memory usage in high-throughput scenarios

# Performance

The library is optimized for performance with:

- Optimized JSONPath expression parsing
- Efficient AST evaluation
- Memory pool reuse
- Concurrent-safe operations
- Minimal allocations in hot paths

Benchmark results on modern hardware:
- Simple path queries: ~0.67 μs/op
- Recursive queries: ~1.7 μs/op
- Filter expressions: ~348 μs/op

# Compatibility

- Go 1.19+ required
- Thread-safe for concurrent use
- Compatible with standard library JSON types
- No external dependencies (only standard library)

# License

MIT License - see LICENSE file for details.

# Contributing

Contributions are welcome! Please see CONTRIBUTING.md for guidelines.

# Support

For issues, feature requests, or questions:
- GitHub Issues: https://github.com/reclaimprotocol/jsonpathplus-go/issues
- Documentation: https://pkg.go.dev/github.com/reclaimprotocol/jsonpathplus-go
*/
package jsonpathplus
