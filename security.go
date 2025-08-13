package jsonpathplus

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Security constants
const (
	DefaultMaxPathComplexity = 50
	DefaultMaxExecutionTime  = 5 * time.Second
	ComplexityPerChar        = 10
	ComplexityRecursive      = 3
	ComplexityFilter         = 4
	MaxRecursionDepthLimit   = 10
)

// SecurityConfig defines security-related configuration.
type SecurityConfig struct {
	// MaxPathComplexity limits the complexity of JSONPath expressions
	MaxPathComplexity int

	// AllowedFunctions lists allowed functions in filter expressions
	AllowedFunctions []string

	// BlockedPatterns contains patterns that are not allowed
	BlockedPatterns []*regexp.Regexp

	// MaxExecutionTime limits execution time per query
	MaxExecutionTime time.Duration

	// EnableSandbox enables sandboxed execution
	EnableSandbox bool

	// AllowNetworkAccess allows network access in expressions
	AllowNetworkAccess bool

	// AllowFileAccess allows file access in expressions
	AllowFileAccess bool
}

// DefaultSecurityConfig returns secure default configuration.
func DefaultSecurityConfig() *SecurityConfig {
	// Common dangerous patterns to block
	blockedPatterns := []*regexp.Regexp{
		regexp.MustCompile(`eval\s*\(`),
		regexp.MustCompile(`exec\s*\(`),
		regexp.MustCompile(`system\s*\(`),
		regexp.MustCompile(`require\s*\(`),
		regexp.MustCompile(`import\s+`),
		regexp.MustCompile(`__.*__`),      // Python dunder methods
		regexp.MustCompile(`\$\{.*\}`),    // Template injection
		regexp.MustCompile(`<!--.*-->`),   // Comment injection
		regexp.MustCompile(`javascript:`), // JavaScript URLs
		regexp.MustCompile(`data:`),       // Data URLs
		regexp.MustCompile(`file:`),       // File URLs
		regexp.MustCompile(`ftp:`),        // FTP URLs
	}

	return &SecurityConfig{
		MaxPathComplexity:  DefaultMaxPathComplexity,
		AllowedFunctions:   []string{"length", "size", "keys", "values", "type"},
		BlockedPatterns:    blockedPatterns,
		MaxExecutionTime:   DefaultMaxExecutionTime,
		EnableSandbox:      true,
		AllowNetworkAccess: false,
		AllowFileAccess:    false,
	}
}

// SecurityValidator validates JSONPath expressions for security issues.
type SecurityValidator struct {
	config *SecurityConfig
}

// NewSecurityValidator creates a new security validator.
func NewSecurityValidator(config *SecurityConfig) *SecurityValidator {
	if config == nil {
		config = DefaultSecurityConfig()
	}

	return &SecurityValidator{
		config: config,
	}
}

// ValidatePath validates a JSONPath expression for security issues.
func (v *SecurityValidator) ValidatePath(path string) error {
	// Check path complexity
	complexity := v.calculateComplexity(path)
	if complexity > v.config.MaxPathComplexity {
		return NewError(ErrInvalidPath,
			fmt.Sprintf("path complexity %d exceeds limit %d", complexity, v.config.MaxPathComplexity),
			path, -1)
	}

	// Check for blocked patterns
	for _, pattern := range v.config.BlockedPatterns {
		if pattern.MatchString(path) {
			return NewError(ErrInvalidExpression,
				fmt.Sprintf("path contains blocked pattern: %s", pattern.String()),
				path, -1)
		}
	}

	// Validate filter expressions
	if err := v.validateFilterExpressions(path); err != nil {
		return err
	}

	return nil
}

// calculateComplexity calculates the complexity score of a JSONPath expression.
func (v *SecurityValidator) calculateComplexity(path string) int {
	complexity := 0

	// Base complexity
	complexity += len(path) / ComplexityPerChar

	// Operators add complexity
	complexity += strings.Count(path, "..") * ComplexityRecursive // Recursive descent
	complexity += strings.Count(path, "*") * 2                    // Wildcards
	complexity += strings.Count(path, "?") * ComplexityFilter     // Filters
	complexity += strings.Count(path, "[") * 1                    // Brackets
	complexity += strings.Count(path, ",") * 1                    // Union

	// Nested expressions add complexity
	depth := 0
	maxDepth := 0
	for _, char := range path {
		switch char {
		case '[', '(':
			depth++
			if depth > maxDepth {
				maxDepth = depth
			}
		case ']', ')':
			depth--
		}
	}
	complexity += maxDepth * 2

	return complexity
}

// validateFilterExpressions validates filter expressions in the path.
func (v *SecurityValidator) validateFilterExpressions(path string) error {
	// Find all filter expressions
	filterPattern := regexp.MustCompile(`\?\([^)]+\)`)
	filters := filterPattern.FindAllString(path, -1)

	for _, filter := range filters {
		if err := v.validateSingleFilter(filter, path); err != nil {
			return err
		}
	}

	return nil
}

// validateSingleFilter validates a single filter expression.
func (v *SecurityValidator) validateSingleFilter(filter, path string) error {
	// Remove the wrapper
	filter = strings.TrimPrefix(filter, "?(")
	filter = strings.TrimSuffix(filter, ")")

	// Check for function calls
	functionPattern := regexp.MustCompile(`(\w+)\s*\(`)
	functions := functionPattern.FindAllStringSubmatch(filter, -1)

	for _, match := range functions {
		if len(match) > 1 {
			funcName := match[1]
			if !v.isFunctionAllowed(funcName) {
				return NewError(ErrInvalidExpression,
					fmt.Sprintf("function '%s' is not allowed", funcName),
					path, -1)
			}
		}
	}

	// Check for network/file access attempts
	if !v.config.AllowNetworkAccess {
		networkPatterns := []string{"http://", "https://", "ws://", "wss://"}
		for _, pattern := range networkPatterns {
			if strings.Contains(strings.ToLower(filter), pattern) {
				return NewError(ErrInvalidExpression,
					"network access not allowed in filter expressions",
					path, -1)
			}
		}
	}

	if !v.config.AllowFileAccess {
		filePatterns := []string{"file://", "../", "./", "/etc/", "/var/", "c:\\", "\\\\"}
		for _, pattern := range filePatterns {
			if strings.Contains(strings.ToLower(filter), pattern) {
				return NewError(ErrInvalidExpression,
					"file access not allowed in filter expressions",
					path, -1)
			}
		}
	}

	return nil
}

// isFunctionAllowed checks if a function is in the allowed list.
func (v *SecurityValidator) isFunctionAllowed(funcName string) bool {
	for _, allowed := range v.config.AllowedFunctions {
		if strings.EqualFold(funcName, allowed) {
			return true
		}
	}
	return false
}

// SanitizePath sanitizes a JSONPath expression by removing dangerous elements.
func (v *SecurityValidator) SanitizePath(path string) string {
	// Remove blocked patterns
	for _, pattern := range v.config.BlockedPatterns {
		path = pattern.ReplaceAllString(path, "")
	}

	// Remove excessive whitespace
	path = regexp.MustCompile(`\s+`).ReplaceAllString(path, " ")
	path = strings.TrimSpace(path)

	return path
}

// RateLimiter implements a simple rate limiter for queries.
type RateLimiter struct {
	requests    map[string]*requestTracker
	maxRequests int
	window      time.Duration
	mu          sync.RWMutex
}

type requestTracker struct {
	count  int
	window time.Time
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:    make(map[string]*requestTracker),
		maxRequests: maxRequests,
		window:      window,
	}
}

// Allow checks if a request is allowed for the given identifier.
func (rl *RateLimiter) Allow(identifier string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	tracker, exists := rl.requests[identifier]
	if !exists {
		rl.requests[identifier] = &requestTracker{
			count:  1,
			window: now,
		}
		return true
	}

	// Check if we're in a new window
	if now.Sub(tracker.window) >= rl.window {
		tracker.count = 1
		tracker.window = now
		return true
	}

	// Check if we're under the limit
	if tracker.count < rl.maxRequests {
		tracker.count++
		return true
	}

	return false
}

// Cleanup removes old entries from the rate limiter.
func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for id, tracker := range rl.requests {
		if now.Sub(tracker.window) >= rl.window*2 {
			delete(rl.requests, id)
		}
	}
}
