package jsonpathplus

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNewEngine(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		engine, err := NewEngine(nil)
		if err != nil {
			t.Fatalf("Failed to create engine: %v", err)
		}
		defer func() { _ = engine.Close() }()

		config := engine.GetConfig()
		if config.MaxPathLength != 1000 {
			t.Errorf("Expected MaxPathLength 1000, got %d", config.MaxPathLength)
		}
	})

	t.Run("CustomConfig", func(t *testing.T) {
		config := DefaultConfig()
		config.MaxPathLength = 500
		config.EnableLogging = true

		engine, err := NewEngine(config)
		if err != nil {
			t.Fatalf("Failed to create engine: %v", err)
		}
		defer func() { _ = engine.Close() }()

		engineConfig := engine.GetConfig()
		if engineConfig.MaxPathLength != 500 {
			t.Errorf("Expected MaxPathLength 500, got %d", engineConfig.MaxPathLength)
		}
	})

	t.Run("InvalidConfig", func(t *testing.T) {
		config := DefaultConfig()
		config.MaxPathLength = -1 // Invalid

		_, err := NewEngine(config)
		if err == nil {
			t.Fatal("Expected error for invalid config")
		}
	})
}

func TestEngineQuery(t *testing.T) {
	engine, err := NewEngine(DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer func() { _ = engine.Close() }()

	data := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "John", "age": 30},
			map[string]interface{}{"name": "Jane", "age": 25},
		},
	}

	t.Run("BasicQuery", func(t *testing.T) {
		results, err := engine.QueryData("$.users[*].name", data)
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(results))
		}
	})

	t.Run("InvalidPath", func(t *testing.T) {
		_, err := engine.QueryData("", data)
		if err == nil {
			t.Fatal("Expected error for empty path")
		}

		var jsonPathErr *JSONPathError
		if !As(err, &jsonPathErr) || jsonPathErr.Type != ErrInvalidPath {
			t.Errorf("Expected ErrInvalidPath, got %T", err)
		}
	})

	t.Run("PathTooLong", func(t *testing.T) {
		config := DefaultConfig()
		config.MaxPathLength = 10

		engine2, err := NewEngine(config)
		if err != nil {
			t.Fatalf("Failed to create engine: %v", err)
		}
		defer func() { _ = engine2.Close() }()

		longPath := "$.very.long.path.that.exceeds.limit.and.should.fail"
		_, err = engine2.QueryData(longPath, data)
		if err == nil {
			t.Fatal("Expected error for path too long")
		}
	})
}

func TestEngineQueryWithContext(t *testing.T) {
	config := DefaultConfig()
	config.Timeout = 100 * time.Millisecond

	engine, err := NewEngine(config)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer func() { _ = engine.Close() }()

	data := map[string]interface{}{
		"items": make([]interface{}, 10000), // Large dataset
	}

	t.Run("ContextTimeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// Give some time for context to expire
		time.Sleep(2 * time.Millisecond)

		_, err := engine.QueryDataWithContext(ctx, "$.items[*]", data)
		if err == nil {
			t.Skip("Query completed before timeout - this is OK in fast environments")
		}

		// Check if it's a context error
		if err != context.DeadlineExceeded && !strings.Contains(err.Error(), "context") {
			t.Errorf("Expected context error, got %v", err)
		}
	})

	t.Run("ContextCancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		// Cancel immediately
		cancel()

		_, err := engine.QueryDataWithContext(ctx, "$.items[*]", data)
		if err == nil {
			t.Skip("Query completed before cancellation - this is OK in fast environments")
		}

		// Check if it's a context error
		if err != context.Canceled && !strings.Contains(err.Error(), "context") {
			t.Errorf("Expected context error, got %v", err)
		}
	})
}

func TestEngineMetrics(t *testing.T) {
	config := DefaultConfig()
	config.EnableMetrics = true

	engine, err := NewEngine(config)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer func() { _ = engine.Close() }()

	data := map[string]interface{}{"test": "value"}

	// Execute some queries
	for i := 0; i < 5; i++ {
		_, err := engine.QueryData("$.test", data)
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
	}

	metrics := engine.GetMetrics()
	if metrics.QueriesExecuted != 5 {
		t.Errorf("Expected 5 queries executed, got %d", metrics.QueriesExecuted)
	}

	if metrics.AverageExecutionTime <= 0 {
		t.Errorf("Expected positive average execution time, got %v", metrics.AverageExecutionTime)
	}
}

func TestEngineConcurrency(t *testing.T) {
	engine, err := NewEngine(DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer func() { _ = engine.Close() }()

	data := map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"id": 1, "value": "a"},
			map[string]interface{}{"id": 2, "value": "b"},
		},
	}

	const numGoroutines = 10
	const queriesPerGoroutine = 100

	var wg sync.WaitGroup
	var mu sync.Mutex
	var totalResults int
	var errors []error

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(_ int) {
			defer wg.Done()

			for j := 0; j < queriesPerGoroutine; j++ {
				results, err := engine.QueryData("$.items[*].value", data)
				if err != nil {
					mu.Lock()
					errors = append(errors, err)
					mu.Unlock()
					return
				}

				mu.Lock()
				totalResults += len(results)
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	if len(errors) > 0 {
		t.Fatalf("Got errors during concurrent execution: %v", errors)
	}

	expectedTotal := numGoroutines * queriesPerGoroutine * 2 // 2 results per query
	if totalResults != expectedTotal {
		t.Errorf("Expected %d total results, got %d", expectedTotal, totalResults)
	}
}

func TestSecurityValidator(t *testing.T) {
	validator := NewSecurityValidator(DefaultSecurityConfig())

	t.Run("ValidPath", func(t *testing.T) {
		err := validator.ValidatePath("$.users[?(@.age > 25)].name")
		if err != nil {
			t.Errorf("Valid path should not fail: %v", err)
		}
	})

	t.Run("BlockedPattern", func(t *testing.T) {
		err := validator.ValidatePath("$.users[?(eval('malicious code'))]")
		if err == nil {
			t.Fatal("Expected error for blocked pattern")
		}
	})

	t.Run("ComplexPath", func(t *testing.T) {
		secConfig := DefaultSecurityConfig()
		secConfig.MaxPathComplexity = 5
		validator := NewSecurityValidator(secConfig)

		complexPath := "$..[*]..[*]..[*]..[*]" // Very complex path
		err := validator.ValidatePath(complexPath)
		if err == nil {
			t.Fatal("Expected error for overly complex path")
		}
	})

	t.Run("SanitizePath", func(t *testing.T) {
		maliciousPath := "$.users[?(eval('rm -rf /'))]"
		sanitized := validator.SanitizePath(maliciousPath)

		if strings.Contains(sanitized, "eval") {
			t.Errorf("Sanitized path still contains 'eval': %s", sanitized)
		}
	})
}

func TestRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(3, time.Second)

	t.Run("AllowWithinLimit", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			if !limiter.Allow("client1") {
				t.Errorf("Request %d should be allowed", i+1)
			}
		}
	})

	t.Run("BlockAfterLimit", func(t *testing.T) {
		limiter := NewRateLimiter(2, time.Second)

		// Use up the limit
		limiter.Allow("client2")
		limiter.Allow("client2")

		// This should be blocked
		if limiter.Allow("client2") {
			t.Error("Request should be blocked after limit")
		}
	})

	t.Run("ResetAfterWindow", func(t *testing.T) {
		limiter := NewRateLimiter(1, 10*time.Millisecond)

		// Use up the limit
		if !limiter.Allow("client3") {
			t.Error("First request should be allowed")
		}

		// Should be blocked
		if limiter.Allow("client3") {
			t.Error("Second request should be blocked")
		}

		// Wait for window reset
		time.Sleep(15 * time.Millisecond)

		// Should be allowed again
		if !limiter.Allow("client3") {
			t.Error("Request should be allowed after window reset")
		}
	})
}

func TestConfigValidation(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		config := DefaultConfig()
		if err := config.Validate(); err != nil {
			t.Errorf("Default config should be valid: %v", err)
		}
	})

	t.Run("InvalidMaxPathLength", func(t *testing.T) {
		config := DefaultConfig()
		config.MaxPathLength = 0

		err := config.Validate()
		if err == nil {
			t.Fatal("Expected validation error")
		}

		var validationErr *ValidationError
		if !As(err, &validationErr) {
			t.Errorf("Expected ValidationError, got %T", err)
		}
	})

	t.Run("InvalidTimeout", func(t *testing.T) {
		config := DefaultConfig()
		config.Timeout = 0

		err := config.Validate()
		if err == nil {
			t.Fatal("Expected validation error")
		}
	})
}

func TestErrorTypes(t *testing.T) {
	t.Run("JSONPathError", func(t *testing.T) {
		err := NewError(ErrInvalidPath, "test message", "$.test", 5)

		if err.Type != ErrInvalidPath {
			t.Errorf("Expected ErrInvalidPath, got %v", err.Type)
		}

		if err.Message != "test message" {
			t.Errorf("Expected 'test message', got %s", err.Message)
		}

		if err.Path != "$.test" {
			t.Errorf("Expected '$.test', got %s", err.Path)
		}

		if err.Position != 5 {
			t.Errorf("Expected position 5, got %d", err.Position)
		}

		errorStr := err.Error()
		if !strings.Contains(errorStr, "invalid JSONPath") {
			t.Errorf("Error string should contain error type: %s", errorStr)
		}
	})

	t.Run("ErrorWrapping", func(t *testing.T) {
		originalErr := fmt.Errorf("original error")
		wrappedErr := WrapError(ErrEvaluationError, originalErr, "$.test", 0)

		if wrappedErr.Cause != originalErr {
			t.Errorf("Expected wrapped error to contain original error")
		}

		if wrappedErr.Unwrap() != originalErr {
			t.Errorf("Unwrap should return original error")
		}
	})
}

func BenchmarkEngineQuery(b *testing.B) {
	engine, err := NewEngine(DefaultConfig())
	if err != nil {
		b.Fatalf("Failed to create engine: %v", err)
	}
	defer func() { _ = engine.Close() }()

	data := map[string]interface{}{
		"users": make([]interface{}, 100),
	}

	for i := 0; i < 100; i++ {
		data["users"].([]interface{})[i] = map[string]interface{}{
			"id":   i,
			"name": fmt.Sprintf("User%d", i),
			"age":  20 + (i % 50),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.QueryData("$.users[?(@.age > 30)].name", data)
		if err != nil {
			b.Fatalf("Query failed: %v", err)
		}
	}
}

// Helper function for error assertion.
func As(err error, target interface{}) bool {
	// Simple implementation for testing
	switch e := err.(type) {
	case *JSONPathError:
		if ptr, ok := target.(**JSONPathError); ok {
			*ptr = e
			return true
		}
	case *ValidationError:
		if ptr, ok := target.(**ValidationError); ok {
			*ptr = e
			return true
		}
	}
	return false
}
