package jsonpathplus

import (
	"time"
)

// Config holds configuration options for JSONPath operations
type Config struct {
	// MaxPathLength limits the maximum length of JSONPath expressions
	MaxPathLength int
	
	// MaxRecursionDepth limits the depth of recursive descent operations
	MaxRecursionDepth int
	
	// MaxResultCount limits the number of results returned
	MaxResultCount int
	
	// Timeout for query execution
	Timeout time.Duration
	
	// EnableLogging enables debug logging
	EnableLogging bool
	
	// CacheSize for compiled JSONPath expressions (0 = no cache)
	CacheSize int
	
	// StrictMode enforces stricter validation
	StrictMode bool
	
	// AllowUnsafeOperations enables potentially unsafe operations
	AllowUnsafeOperations bool
	
	// MaxMemoryUsage limits memory usage in bytes (0 = no limit)
	MaxMemoryUsage int64
	
	// EnableMetrics enables performance metrics collection
	EnableMetrics bool
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		MaxPathLength:         1000,
		MaxRecursionDepth:     100,
		MaxResultCount:        10000,
		Timeout:               30 * time.Second,
		EnableLogging:         false,
		CacheSize:             100,
		StrictMode:            false,
		AllowUnsafeOperations: false,
		MaxMemoryUsage:        100 * 1024 * 1024, // 100MB
		EnableMetrics:         false,
	}
}

// ProductionConfig returns a configuration suitable for production use
func ProductionConfig() *Config {
	config := DefaultConfig()
	config.MaxPathLength = 500
	config.MaxRecursionDepth = 50
	config.MaxResultCount = 1000
	config.Timeout = 5 * time.Second
	config.StrictMode = true
	config.AllowUnsafeOperations = false
	config.MaxMemoryUsage = 50 * 1024 * 1024 // 50MB
	config.EnableMetrics = true
	return config
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.MaxPathLength <= 0 {
		return &ValidationError{
			Field:   "MaxPathLength",
			Value:   c.MaxPathLength,
			Message: "must be greater than 0",
		}
	}
	
	if c.MaxRecursionDepth <= 0 {
		return &ValidationError{
			Field:   "MaxRecursionDepth",
			Value:   c.MaxRecursionDepth,
			Message: "must be greater than 0",
		}
	}
	
	if c.MaxResultCount <= 0 {
		return &ValidationError{
			Field:   "MaxResultCount",
			Value:   c.MaxResultCount,
			Message: "must be greater than 0",
		}
	}
	
	if c.Timeout <= 0 {
		return &ValidationError{
			Field:   "Timeout",
			Value:   c.Timeout,
			Message: "must be greater than 0",
		}
	}
	
	if c.CacheSize < 0 {
		return &ValidationError{
			Field:   "CacheSize",
			Value:   c.CacheSize,
			Message: "must be >= 0",
		}
	}
	
	if c.MaxMemoryUsage < 0 {
		return &ValidationError{
			Field:   "MaxMemoryUsage",
			Value:   c.MaxMemoryUsage,
			Message: "must be >= 0",
		}
	}
	
	return nil
}

// Clone creates a copy of the configuration
func (c *Config) Clone() *Config {
	clone := *c
	return &clone
}