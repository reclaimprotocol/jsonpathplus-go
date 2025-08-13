package jsonpathplus

import (
	"time"
)

// Configuration constants
const (
	// Default limits
	DefaultMaxPathLength     = 1000
	DefaultMaxRecursionDepth = 100
	DefaultMaxResultCount    = 10000
	DefaultTimeout           = 30 * time.Second
	DefaultMaxMemoryUsage    = 100 * 1024 * 1024 // 100MB

	// Production limits
	ProductionMaxPathLength     = 500
	ProductionMaxRecursionDepth = 50
	ProductionMaxResultCount    = 1000
	ProductionTimeout           = 5 * time.Second
	ProductionMaxMemoryUsage    = 50 * 1024 * 1024 // 50MB
)

// Config holds configuration options for JSONPath operations.
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

	// StrictMode enforces stricter validation
	StrictMode bool

	// AllowUnsafeOperations enables potentially unsafe operations
	AllowUnsafeOperations bool

	// MaxMemoryUsage limits memory usage in bytes (0 = no limit)
	MaxMemoryUsage int64

	// EnableMetrics enables performance metrics collection
	EnableMetrics bool
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		MaxPathLength:         DefaultMaxPathLength,
		MaxRecursionDepth:     DefaultMaxRecursionDepth,
		MaxResultCount:        DefaultMaxResultCount,
		Timeout:               DefaultTimeout,
		EnableLogging:         false,
		StrictMode:            false,
		AllowUnsafeOperations: false,
		MaxMemoryUsage:        DefaultMaxMemoryUsage,
		EnableMetrics:         false,
	}
}

// ProductionConfig returns a configuration suitable for production use.
func ProductionConfig() *Config {
	config := DefaultConfig()
	config.MaxPathLength = ProductionMaxPathLength
	config.MaxRecursionDepth = ProductionMaxRecursionDepth
	config.MaxResultCount = ProductionMaxResultCount
	config.Timeout = ProductionTimeout
	config.StrictMode = true
	config.AllowUnsafeOperations = false
	config.MaxMemoryUsage = ProductionMaxMemoryUsage
	config.EnableMetrics = true
	return config
}

// Validate validates the configuration.
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

	if c.MaxMemoryUsage < 0 {
		return &ValidationError{
			Field:   "MaxMemoryUsage",
			Value:   c.MaxMemoryUsage,
			Message: "must be >= 0",
		}
	}

	return nil
}

// Clone creates a copy of the configuration.
func (c *Config) Clone() *Config {
	clone := *c
	return &clone
}
