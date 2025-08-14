package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

// ProductionExample demonstrates production-ready usage
func main() {
	fmt.Println("JSONPath Plus - Production Example")
	fmt.Println("==================================")

	// 1. Create production-ready engine
	engine := createProductionEngine()
	defer engine.Close()

	// 2. Demonstrate basic usage
	basicUsage(engine)

	// 3. Demonstrate security features
	securityExample(engine)

	// 4. Demonstrate performance monitoring
	performanceMonitoring(engine)

	// 5. Demonstrate rate limiting
	rateLimitingExample()

	// 6. Demonstrate HTTP server integration
	fmt.Println("\n6. HTTP Server Integration:")
	fmt.Println("Starting HTTP server on :8080...")
	fmt.Println("Try: curl 'http://localhost:8080/query?path=$.users[*].name'")

	server := createHTTPServer(engine)
	log.Fatal(server.ListenAndServe())
}

// createProductionEngine creates a production-ready JSONPath engine
func createProductionEngine() *jp.JSONPathEngine {
	fmt.Println("\n1. Creating Production Engine:")

	// Use production configuration
	config := jp.ProductionConfig()
	config.EnableLogging = true
	config.EnableMetrics = true

	engine, err := jp.NewEngine(config)
	if err != nil {
		log.Fatalf("Failed to create engine: %v", err)
	}

	// Set custom logger (optional)
	logger := jp.NewDefaultLogger(jp.LogLevelInfo)
	engine.SetLogger(logger)

	fmt.Printf("✓ Engine created with production config\n")
	fmt.Printf("  - Max path length: %d\n", config.MaxPathLength)
	fmt.Printf("  - Max recursion depth: %d\n", config.MaxRecursionDepth)
	fmt.Printf("  - Timeout: %v\n", config.Timeout)
	fmt.Printf("  - Metrics enabled: %v\n", config.EnableMetrics)

	return engine
}

// basicUsage demonstrates basic JSONPath operations
func basicUsage(engine *jp.JSONPathEngine) {
	fmt.Println("\n2. Basic Usage:")

	// Sample data
	jsonData := `{
		"company": {
			"name": "TechCorp",
			"departments": [
				{
					"name": "Engineering",
					"employees": [
						{"name": "Alice", "role": "Senior Developer", "salary": 95000},
						{"name": "Bob", "role": "DevOps Engineer", "salary": 85000}
					]
				},
				{
					"name": "Sales",
					"employees": [
						{"name": "Charlie", "role": "Account Manager", "salary": 75000},
						{"name": "Diana", "role": "Sales Director", "salary": 110000}
					]
				}
			]
		}
	}`

	data, err := jp.JSONParse(jsonData)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Execute various queries
	queries := []string{
		"$.company.name",
		"$.company.departments[*].name",
		"$..employees[*].name",
		"$..employees[?(@.salary > 80000)].name",
		"$..employees[?(@.role == 'Senior Developer')]",
	}

	for _, query := range queries {
		fmt.Printf("\nQuery: %s\n", query)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		results, err := engine.QueryDataWithContext(ctx, query, data)
		cancel()

		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}

		for i, result := range results {
			fmt.Printf("  [%d] Value: %v, Path: %s, OriginalIndex: %d\n",
				i, result.Value, result.Path, result.OriginalIndex)
		}
	}
}

// securityExample demonstrates security features
func securityExample(engine *jp.JSONPathEngine) {
	fmt.Println("\n3. Security Features:")

	// Create security validator
	secConfig := jp.DefaultSecurityConfig()
	validator := jp.NewSecurityValidator(secConfig)

	// Test various paths for security
	testPaths := []string{
		"$.users[*].name",               // Safe
		"$.users[?(@.age > 21)]",        // Safe
		"$.users[?(eval('malicious'))]", // Unsafe - contains eval
		"$.users[?(@.__proto__)]",       // Unsafe - prototype pollution
		"$.users[?(require('fs'))]",     // Unsafe - require statement
		"$" + string(make([]byte, 600)), // Unsafe - too long
	}

	for _, path := range testPaths {
		displayPath := path
		if len(displayPath) > 50 {
			displayPath = displayPath[:50] + "..."
		}

		err := validator.ValidatePath(path)
		if err != nil {
			fmt.Printf("  ❌ BLOCKED: %s - %v\n", displayPath, err)
		} else {
			fmt.Printf("  ✅ ALLOWED: %s\n", displayPath)
		}
	}

	// Demonstrate path sanitization
	maliciousPath := "$.users[?(eval('rm -rf /') && @.age > 21)]"
	sanitizedPath := validator.SanitizePath(maliciousPath)
	fmt.Printf("\nPath sanitization:\n")
	fmt.Printf("  Original: %s\n", maliciousPath)
	fmt.Printf("  Sanitized: %s\n", sanitizedPath)
}

// performanceMonitoring demonstrates performance monitoring
func performanceMonitoring(engine *jp.JSONPathEngine) {
	fmt.Println("\n4. Performance Monitoring:")

	// Create test data
	data := map[string]interface{}{
		"items": make([]interface{}, 1000),
	}

	for i := 0; i < 1000; i++ {
		data["items"].([]interface{})[i] = map[string]interface{}{
			"id":    i,
			"value": fmt.Sprintf("item_%d", i),
			"score": float64(i % 100),
		}
	}

	// Execute some queries for metrics
	queries := []string{
		"$.items[*].id",
		"$.items[?(@.score > 50)].value",
		"$..score",
	}

	start := time.Now()
	for i, query := range queries {
		for j := 0; j < 10; j++ {
			results, err := engine.QueryData(query, data)
			if err != nil {
				fmt.Printf("  Query %d failed: %v\n", i+1, err)
				continue
			}
			fmt.Printf("  Query %d.%d: %d results\n", i+1, j+1, len(results))
		}
	}
	totalTime := time.Since(start)

	// Get metrics
	metrics := engine.GetMetrics()

	fmt.Printf("\nPerformance Metrics:\n")
	fmt.Printf("  Total execution time: %v\n", totalTime)
	fmt.Printf("  Queries executed: %d\n", metrics.QueriesExecuted)
	fmt.Printf("  Average execution time: %v\n", metrics.AverageExecutionTime)
	fmt.Printf("  Error count: %d\n", metrics.ErrorCount)
	fmt.Printf("  Memory usage: %d bytes\n", metrics.MemoryUsage)

	// Production setup focuses on security and monitoring
}

// rateLimitingExample demonstrates rate limiting
func rateLimitingExample() {
	fmt.Println("\n5. Rate Limiting:")

	// Create rate limiter: 3 requests per second
	limiter := jp.NewRateLimiter(3, time.Second)

	client := "client_123"

	fmt.Printf("Testing rate limiter (3 requests/second):\n")
	for i := 1; i <= 5; i++ {
		allowed := limiter.Allow(client)
		status := "✅ ALLOWED"
		if !allowed {
			status = "❌ BLOCKED"
		}
		fmt.Printf("  Request %d: %s\n", i, status)

		if i == 3 {
			fmt.Printf("  Waiting 1 second for rate limit reset...\n")
			time.Sleep(time.Second)
		}
	}
}

// HTTPHandler handles JSONPath queries via HTTP
type HTTPHandler struct {
	engine    *jp.JSONPathEngine
	validator *jp.SecurityValidator
	limiter   *jp.RateLimiter
}

// createHTTPServer creates an HTTP server with JSONPath endpoints
func createHTTPServer(engine *jp.JSONPathEngine) *http.Server {
	handler := &HTTPHandler{
		engine:    engine,
		validator: jp.NewSecurityValidator(jp.DefaultSecurityConfig()),
		limiter:   jp.NewRateLimiter(10, time.Minute), // 10 requests per minute
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/query", handler.handleQuery)
	mux.HandleFunc("/metrics", handler.handleMetrics)
	mux.HandleFunc("/health", handler.handleHealth)

	return &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}
}

// handleQuery handles JSONPath query requests
func (h *HTTPHandler) handleQuery(w http.ResponseWriter, r *http.Request) {
	// Apply rate limiting
	clientIP := r.RemoteAddr
	if !h.limiter.Allow(clientIP) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// Get query parameters
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Missing 'path' parameter", http.StatusBadRequest)
		return
	}

	// Validate path for security
	if err := h.validator.ValidatePath(path); err != nil {
		http.Error(w, fmt.Sprintf("Invalid path: %v", err), http.StatusBadRequest)
		return
	}

	// Sample data for demo
	data := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "John", "age": 30, "city": "New York"},
			map[string]interface{}{"name": "Jane", "age": 25, "city": "San Francisco"},
			map[string]interface{}{"name": "Bob", "age": 35, "city": "Chicago"},
		},
		"products": []interface{}{
			map[string]interface{}{"name": "Laptop", "price": 999.99, "category": "Electronics"},
			map[string]interface{}{"name": "Book", "price": 29.99, "category": "Education"},
		},
	}

	// Execute query with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	results, err := h.engine.QueryDataWithContext(ctx, path, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Return results as JSON
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"path":    path,
		"count":   len(results),
		"results": results,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// handleMetrics returns engine metrics
func (h *HTTPHandler) handleMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := h.engine.GetMetrics()

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"metrics": metrics,
	}

	json.NewEncoder(w).Encode(response)
}

// handleHealth returns health status
func (h *HTTPHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":  "healthy",
		"version": jp.Version,
		"uptime":  time.Since(time.Now()).String(), // Placeholder
	}

	json.NewEncoder(w).Encode(response)
}
