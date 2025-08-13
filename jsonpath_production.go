package jsonpathplus

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

// JSONPathEngine is the main production-ready JSONPath engine.
type JSONPathEngine struct {
	config  *Config
	logger  Logger
	metrics *MetricsCollector
	mu      sync.RWMutex
}

// NewEngine creates a new JSONPath engine with the given configuration.
func NewEngine(config *Config) (*JSONPathEngine, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	var logger Logger
	if config.EnableLogging {
		logger = NewDefaultLogger(LogLevelInfo)
	} else {
		logger = &NoOpLogger{}
	}

	return &JSONPathEngine{
		config:  config.Clone(),
		logger:  logger,
		metrics: NewMetricsCollector(config.EnableMetrics),
	}, nil
}

// SetLogger sets a custom logger.
func (e *JSONPathEngine) SetLogger(logger Logger) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.logger = logger
}

// Query executes a JSONPath query with string character position tracking.
func (e *JSONPathEngine) Query(path string, jsonStr string) ([]Result, error) {
	return QueryWithStringIndex(path, jsonStr)
}

// QueryData executes a JSONPath query on parsed data with timeout and resource limits.
func (e *JSONPathEngine) QueryData(path string, data interface{}) ([]Result, error) {
	return e.QueryDataWithContext(context.Background(), path, data)
}

// QueryDataWithContext executes a JSONPath query on parsed data with context for cancellation.
func (e *JSONPathEngine) QueryDataWithContext(ctx context.Context, path string, data interface{}) ([]Result, error) {
	start := time.Now()
	var err error

	defer func() {
		duration := time.Since(start)
		e.metrics.RecordQuery(duration, err)

		if err != nil {
			e.logger.Error("Query failed",
				String("path", path),
				Duration("duration", duration),
				Error("error", err),
			)
		} else {
			e.logger.Debug("Query completed",
				String("path", path),
				Duration("duration", duration),
			)
		}
	}()

	// Validate input
	if err = e.validateInput(path, data); err != nil {
		return nil, err
	}

	// Check timeout
	if e.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.config.Timeout)
		defer cancel()
	}

	// Create execution context
	execCtx := &executionContext{
		ctx:            ctx,
		config:         e.config,
		logger:         e.logger,
		metrics:        e.metrics,
		recursionDepth: 0,
		resultCount:    0,
	}

	// Compile AST
	ast, err := e.compileAST(path)
	if err != nil {
		return nil, fmt.Errorf("failed to compile path: %w", err)
	}

	// Execute query with resource limits
	results, err := e.executeWithLimits(execCtx, ast, data)
	if err != nil {
		return nil, err
	}

	// Apply result limits
	if len(results) > e.config.MaxResultCount {
		e.logger.Warn("Result count exceeded limit",
			Int("count", len(results)),
			Int("limit", e.config.MaxResultCount),
		)
		results = results[:e.config.MaxResultCount]
	}

	return results, nil
}

// executionContext holds state for query execution.
type executionContext struct {
	ctx            context.Context
	config         *Config
	logger         Logger
	metrics        *MetricsCollector
	recursionDepth int
	resultCount    int
}

// validateInput validates the input parameters.
func (e *JSONPathEngine) validateInput(path string, data interface{}) error {
	if path == "" {
		return NewError(ErrInvalidPath, "empty path", "", -1)
	}

	if len(path) > e.config.MaxPathLength {
		return &PathLengthError{
			Length: len(path),
			Limit:  e.config.MaxPathLength,
		}
	}

	if data == nil {
		return NewError(ErrInvalidJSON, "nil data", path, -1)
	}

	// Basic path validation
	if !strings.HasPrefix(path, "$") {
		return NewError(ErrInvalidPath, "path must start with '$'", path, 0)
	}

	return nil
}

// compileAST compiles the JSONPath expression into an AST.
func (e *JSONPathEngine) compileAST(path string) (*astNode, error) {
	// Compile AST
	tokens, err := e.tokenizeWithValidation(path)
	if err != nil {
		return nil, err
	}

	ast, err := e.parseWithValidation(tokens, path)
	if err != nil {
		return nil, err
	}

	return ast, nil
}

// tokenizeWithValidation tokenizes with enhanced validation.
func (e *JSONPathEngine) tokenizeWithValidation(path string) ([]token, error) {
	tokens, err := tokenize(path)
	if err != nil {
		return nil, WrapError(ErrParseError, err, path, -1)
	}

	// Additional validation in strict mode
	if e.config.StrictMode {
		if err := e.validateTokens(tokens, path); err != nil {
			return nil, err
		}
	}

	return tokens, nil
}

// validateTokens performs additional token validation.
func (e *JSONPathEngine) validateTokens(tokens []token, path string) error {
	if len(tokens) == 0 {
		return NewError(ErrInvalidPath, "no tokens found", path, -1)
	}

	if tokens[0].Type != tokenRoot {
		return NewError(ErrInvalidPath, "path must start with root token", path, 0)
	}

	// Check for suspicious patterns
	for i, tok := range tokens {
		switch tok.Type {
		case tokenFilter:
			if !e.config.AllowUnsafeOperations {
				// Basic filter validation
				if strings.Contains(tok.Value, "eval") || strings.Contains(tok.Value, "exec") {
					return NewError(ErrInvalidExpression, "unsafe filter expression", path, i)
				}
			}
		case tokenRoot, tokenCurrent, tokenDot, tokenDoubleDot,
			tokenBracketOpen, tokenBracketClose, tokenIdentifier,
			tokenNumber, tokenString, tokenWildcard, tokenSlice,
			tokenComma, tokenUnion:
			// These tokens are safe, no special validation needed
		}
	}

	return nil
}

// parseWithValidation parses with enhanced validation.
func (e *JSONPathEngine) parseWithValidation(tokens []token, path string) (*astNode, error) {
	ast, err := parse(tokens)
	if err != nil {
		return nil, WrapError(ErrParseError, err, path, -1)
	}

	// Validate AST structure
	if err := e.validateAST(ast, path); err != nil {
		return nil, err
	}

	return ast, nil
}

// validateAST validates the AST structure.
func (e *JSONPathEngine) validateAST(ast *astNode, path string) error {
	if ast == nil {
		return NewError(ErrParseError, "nil AST", path, -1)
	}

	// Count recursive nodes to prevent excessive recursion
	recursiveCount := 0
	var validateNode func(*astNode) error
	validateNode = func(node *astNode) error {
		if node.Type == "recursive" {
			recursiveCount++
			if recursiveCount > MaxRecursionDepthLimit {
				return NewError(ErrRecursionLimit, "too many recursive operators", path, -1)
			}
		}

		for _, child := range node.Children {
			if err := validateNode(child); err != nil {
				return err
			}
		}
		return nil
	}

	return validateNode(ast)
}

// executeWithLimits executes the query with resource monitoring.
func (e *JSONPathEngine) executeWithLimits(
	execCtx *executionContext, ast *astNode, data interface{},
) ([]Result, error) {
	// Check memory usage periodically
	if execCtx.config.MaxMemoryUsage > 0 {
		var m runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m)

		currentUsage := int64(m.Alloc) // #nosec G115 - Memory stats are safe to convert
		execCtx.metrics.UpdateMemoryUsage(currentUsage)

		if currentUsage > execCtx.config.MaxMemoryUsage {
			return nil, NewError(ErrEvaluationError,
				fmt.Sprintf("memory usage %d exceeds limit %d", currentUsage, execCtx.config.MaxMemoryUsage),
				"", -1)
		}
	}

	// Check context cancellation
	select {
	case <-execCtx.ctx.Done():
		return nil, execCtx.ctx.Err()
	default:
	}

	// Execute the evaluation
	results, err := e.evaluateWithContext(execCtx, ast, data)
	if err != nil {
		return nil, WrapError(ErrEvaluationError, err, "", -1)
	}

	return results, nil
}

// evaluateWithContext is a wrapper around the original evaluate function with context.
func (e *JSONPathEngine) evaluateWithContext(
	execCtx *executionContext, ast *astNode, data interface{},
) ([]Result, error) {
	// Use the original evaluate function but with context checking
	options := &Options{
		ResultType: "value",
		Flatten:    false,
		Wrap:       true,
	}

	return e.evaluateWithResourceLimits(execCtx, ast, data, options)
}

// evaluateWithResourceLimits performs evaluation with resource limits.
func (e *JSONPathEngine) evaluateWithResourceLimits(
	execCtx *executionContext, ast *astNode, data interface{}, options *Options,
) ([]Result, error) {
	results := []Result{{
		Value:         data,
		Path:          "$",
		Index:         0,
		OriginalIndex: 0,
	}}

	for _, child := range ast.Children {
		// Check context cancellation
		select {
		case <-execCtx.ctx.Done():
			return nil, execCtx.ctx.Err()
		default:
		}

		// Check result count limit
		if len(results) > execCtx.config.MaxResultCount {
			break
		}

		newResults, err := e.evaluateNodeWithLimits(execCtx, child, results, options)
		if err != nil {
			return nil, err
		}
		results = newResults
	}

	if options.ResultType == "path" {
		for i := range results {
			results[i].Value = results[i].Path
		}
	}

	return results, nil
}

// evaluateNodeWithLimits evaluates a single node with resource limits.
func (e *JSONPathEngine) evaluateNodeWithLimits(
	execCtx *executionContext, node *astNode, contexts []Result, options *Options,
) ([]Result, error) {
	// Check recursion depth
	if node.Type == "recursive" {
		if execCtx.recursionDepth >= execCtx.config.MaxRecursionDepth {
			return nil, &RecursionLimitError{
				Depth: execCtx.recursionDepth,
				Limit: execCtx.config.MaxRecursionDepth,
			}
		}
		execCtx.recursionDepth++
		defer func() { execCtx.recursionDepth-- }()
	}

	// Use original evaluation logic but with context checking
	return evaluateNode(node, contexts, options), nil
}

// GetMetrics returns current performance metrics.
func (e *JSONPathEngine) GetMetrics() Metrics {
	return e.metrics.GetMetrics()
}

// GetConfig returns a copy of the current configuration.
func (e *JSONPathEngine) GetConfig() *Config {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.config.Clone()
}

// UpdateConfig updates the engine configuration.
func (e *JSONPathEngine) UpdateConfig(newConfig *Config) error {
	if err := newConfig.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.config = newConfig.Clone()

	return nil
}

// Close performs cleanup and releases resources.
func (e *JSONPathEngine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.metrics != nil {
		e.metrics.Reset()
	}

	e.logger.Info("JSONPath engine closed")
	return nil
}
