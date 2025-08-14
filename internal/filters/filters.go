package filters

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/reclaimprotocol/jsonpathplus-go/pkg/types"
	"github.com/reclaimprotocol/jsonpathplus-go/pkg/utils"
)

// FilterEvaluator handles filter expression evaluation
type FilterEvaluator struct{}

// NewFilterEvaluator creates a new filter evaluator
func NewFilterEvaluator() *FilterEvaluator {
	return &FilterEvaluator{}
}

// EvaluateFilter evaluates a filter expression with enhanced context support
func (f *FilterEvaluator) EvaluateFilter(filter string, ctx *types.Context) bool {
	filter = strings.TrimPrefix(filter, "?(")
	filter = strings.TrimSuffix(filter, ")")

	return f.evaluateFilterExpression(filter, ctx)
}

// evaluateFilterExpression evaluates the main filter logic with context
func (f *FilterEvaluator) evaluateFilterExpression(expr string, ctx *types.Context) bool {
	expr = f.cleanFilterExpression(expr)

	// Handle logical operators (&&, ||)
	if result, ok := f.tryLogicalFilter(expr, ctx); ok {
		return result
	}

	// Handle negation (!@.field)
	if result, ok := f.tryNegationFilter(expr, ctx); ok {
		return result
	}

	// Handle root references ($.field)
	if result, ok := f.tryRootComparisonFilter(expr, ctx); ok {
		return result
	}

	// Handle context-based filters (@property, @parent, @parentProperty)
	if result, ok := f.tryContextFilter(expr, ctx); ok {
		return result
	}

	// Handle function predicates (.match(), .contains(), .startsWith(), .endsWith())
	if result, ok := f.tryFunctionPredicateFilter(expr, ctx.Current); ok {
		return result
	}

	// Try direct value comparison first (@ > 5)
	if result, ok := f.tryDirectComparisonFilter(expr, ctx.Current); ok {
		return result
	}

	// Try comparison expression (property-based)
	if result, ok := f.tryComparisonFilter(expr, ctx.Current); ok {
		return result
	}

	// Try existence check
	if result, ok := f.tryExistenceFilter(expr, ctx.Current); ok {
		return result
	}

	return false
}

// tryContextFilter handles context-based filter expressions
func (f *FilterEvaluator) tryContextFilter(expr string, ctx *types.Context) (bool, bool) {
	// Handle @parentProperty expressions FIRST (before @parent)
	if strings.Contains(expr, "@parentProperty") {
		return f.handleParentPropertyFilter(expr, ctx)
	}

	// Handle @property expressions
	if strings.Contains(expr, "@property") {
		return f.handlePropertyFilter(expr, ctx)
	}

	// Handle @parent expressions
	if strings.Contains(expr, "@parent") {
		return f.handleParentFilter(expr, ctx)
	}

	// Handle @path expressions (check for @path but not @parentProperty)
	if strings.Contains(expr, "@path") && !strings.Contains(expr, "@parentProperty") {
		return f.handlePathFilter(expr, ctx)
	}

	return false, false
}

// handlePropertyFilter handles @property-based filters
func (f *FilterEvaluator) handlePropertyFilter(expr string, ctx *types.Context) (bool, bool) {
	// Pattern: @property === 'value' or @property !== 'value'
	re := regexp.MustCompile(`@property\s*(===|!==|==|!=)\s*['"](.*?)['"]`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 3 {
		return false, false
	}

	operator := matches[1]
	expectedValue := matches[2]
	actualValue := ctx.GetPropertyName()

	switch operator {
	case "===", "==":
		return actualValue == expectedValue, true
	case "!==", "!=":
		return actualValue != expectedValue, true
	}

	return false, true
}

// handleParentFilter handles @parent-based filters
func (f *FilterEvaluator) handleParentFilter(expr string, ctx *types.Context) (bool, bool) {
	// Handle @parent.property.nested.path expressions
	re := regexp.MustCompile(`@parent\.([a-zA-Z_]\w*(?:\.[a-zA-Z_]\w*)*)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) < 2 {
		// Simple @parent existence check
		if strings.TrimSpace(expr) == "@parent" {
			return ctx.GetParent() != nil, true
		}
		return false, false
	}

	propertyPath := matches[1]
	parent := ctx.GetParent()

	if parent == nil {
		return false, true
	}

	// Get the property value from parent (supports nested paths like "bicycle.color")
	parentValue := utils.GetPropertyValue(parent, propertyPath)

	// Check if this is just an existence check
	if strings.TrimSpace(expr) == fmt.Sprintf("@parent.%s", propertyPath) {
		return parentValue != nil, true
	}

	// Handle comparison expressions with parent properties
	// Pattern: @parent.property.nested == 'value'
	compRe := regexp.MustCompile(fmt.Sprintf(`@parent\.%s\s*(===|!==|==|!=|<=|>=|<|>)\s*(.+)`, regexp.QuoteMeta(propertyPath)))
	compMatches := compRe.FindStringSubmatch(expr)
	if len(compMatches) != 3 {
		return false, false
	}

	operator := compMatches[1]
	valueStr := strings.TrimSpace(compMatches[2])
	expectedValue := utils.ParseValue(valueStr)

	return utils.CompareValues(parentValue, operator, expectedValue), true
}

// handlePathFilter handles @path-based filters
func (f *FilterEvaluator) handlePathFilter(expr string, ctx *types.Context) (bool, bool) {
	// Handle simple @path existence check
	if strings.TrimSpace(expr) == "@path" {
		return ctx.Path != "", true
	}

	// Pattern: @path === 'value' or @path !== 'value'
	// Handle paths with nested quotes by finding the outer quote boundaries
	re := regexp.MustCompile(`@path\s*(===|!==|==|!=)\s*['\"](.*)['\"]`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 3 {
		return false, false
	}

	operator := matches[1]
	expectedPath := matches[2]
	actualPath := ctx.GetBracketPath() // Use bracket notation

	switch operator {
	case "===", "==":
		return actualPath == expectedPath, true
	case "!==", "!=":
		return actualPath != expectedPath, true
	}

	return false, true
}

// handleParentPropertyFilter handles @parentProperty-based filters
func (f *FilterEvaluator) handleParentPropertyFilter(expr string, ctx *types.Context) (bool, bool) {
	// Remove outer quotes if present
	if (strings.HasPrefix(expr, "'") && strings.HasSuffix(expr, "'")) ||
		(strings.HasPrefix(expr, "\"") && strings.HasSuffix(expr, "\"")) {
		expr = expr[1 : len(expr)-1]
	}

	// Pattern: @parentProperty === 'value' or @parentProperty !== 'value'
	re := regexp.MustCompile(`@parentProperty\s*(===|!==|==|!=)\s*['"](.*?)['"]`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 3 {
		// Handle numeric comparisons for array indices
		numRe := regexp.MustCompile(`@parentProperty\s*(===|!==|==|!=|<=|>=|<|>)\s*(\d+)`)
		numMatches := numRe.FindStringSubmatch(expr)
		if len(numMatches) != 3 {
			return false, false
		}

		operator := numMatches[1]
		expectedIndex, err := strconv.Atoi(numMatches[2])
		if err != nil {
			return false, false
		}

		actualIndex, err := strconv.Atoi(ctx.GetPropertyName())
		if err != nil {
			return false, true // Not a numeric property
		}

		switch operator {
		case "===", "==":
			return actualIndex == expectedIndex, true
		case "!==", "!=":
			return actualIndex != expectedIndex, true
		case "<":
			return actualIndex < expectedIndex, true
		case "<=":
			return actualIndex <= expectedIndex, true
		case ">":
			return actualIndex > expectedIndex, true
		case ">=":
			return actualIndex >= expectedIndex, true
		}

		return false, true
	}

	operator := matches[1]
	expectedValue := matches[2]
	actualValue := ctx.GetPropertyName()

	switch operator {
	case "===", "==":
		return actualValue == expectedValue, true
	case "!==", "!=":
		return actualValue != expectedValue, true
	}

	return false, true
}

// Legacy filter methods (refactored from original code)

func (f *FilterEvaluator) cleanFilterExpression(expr string) string {
	expr = strings.TrimSpace(expr)
	// Don't remove @ symbols for context-based expressions
	// Only remove leading @ for simple property references like "@.field"
	if strings.HasPrefix(expr, "@.") {
		expr = strings.TrimPrefix(expr, "@")
	}
	// Handle direct @ comparisons by replacing @ with empty to indicate current value
	// Simplify by unconditionally trimming the "@ " prefix if present
	expr = strings.TrimPrefix(expr, "@ ")
	return strings.TrimSpace(expr)
}

func (f *FilterEvaluator) tryLogicalFilter(expr string, ctx *types.Context) (bool, bool) {
	// Handle && operator
	if andPos := f.findLogicalOperator(expr, "&&"); andPos != -1 {
		left := strings.TrimSpace(expr[:andPos])
		right := strings.TrimSpace(expr[andPos+2:])

		// Strip outer parentheses if present
		left = f.stripOuterParentheses(left)
		right = f.stripOuterParentheses(right)

		leftResult := f.evaluateFilterExpression(left, ctx)
		if !leftResult {
			return false, true // Short-circuit evaluation
		}

		rightResult := f.evaluateFilterExpression(right, ctx)
		return rightResult, true
	}

	// Handle || operator
	if orPos := f.findLogicalOperator(expr, "||"); orPos != -1 {
		left := strings.TrimSpace(expr[:orPos])
		right := strings.TrimSpace(expr[orPos+2:])

		// Strip outer parentheses if present
		left = f.stripOuterParentheses(left)
		right = f.stripOuterParentheses(right)

		leftResult := f.evaluateFilterExpression(left, ctx)
		if leftResult {
			return true, true // Short-circuit evaluation
		}

		rightResult := f.evaluateFilterExpression(right, ctx)
		return rightResult, true
	}

	return false, false
}

func (f *FilterEvaluator) tryNegationFilter(expr string, ctx *types.Context) (bool, bool) {
	if !strings.HasPrefix(expr, "!") {
		return false, false
	}

	innerExpr := strings.TrimSpace(expr[1:])
	if strings.HasPrefix(innerExpr, "(") && strings.HasSuffix(innerExpr, ")") {
		innerExpr = innerExpr[1 : len(innerExpr)-1]
	}

	result := f.evaluateFilterExpression(innerExpr, ctx)
	return !result, true
}

func (f *FilterEvaluator) tryRootComparisonFilter(expr string, ctx *types.Context) (bool, bool) {
	re := regexp.MustCompile(`\.(\w+)\s*(===|!==|<=|>=|==|!=|<|>)\s*\$\.(\w+)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 4 {
		return false, false
	}

	property := matches[1]
	operator := matches[2]
	rootProperty := matches[3]

	obj, ok := ctx.Current.(map[string]interface{})
	if !ok {
		return false, true
	}

	propValue, exists := obj[property]
	if !exists {
		return operator == "!=" || operator == "!==", true
	}

	rootValue := utils.GetPropertyValue(ctx.Root, rootProperty)
	return utils.CompareValues(propValue, operator, rootValue), true
}

func (f *FilterEvaluator) tryDirectComparisonFilter(expr string, current interface{}) (bool, bool) {
	re := regexp.MustCompile(`^\s*(===|!==|<=|>=|==|!=|<|>)\s*(.+)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 3 {
		return false, false
	}

	operator := matches[1]
	valueStr := strings.TrimSpace(matches[2])
	parsedValue := utils.ParseValue(valueStr)

	return utils.CompareValues(current, operator, parsedValue), true
}

func (f *FilterEvaluator) tryComparisonFilter(expr string, current interface{}) (bool, bool) {
	re := regexp.MustCompile(`\.(\w+)\s*(===|!==|<=|>=|==|!=|<|>)\s*(.+)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 4 {
		return false, false
	}

	property := matches[1]
	operator := matches[2]
	valueStr := strings.TrimSpace(matches[3])

	obj, ok := current.(map[string]interface{})
	if !ok {
		return false, true
	}

	propValue, exists := obj[property]
	if !exists {
		return operator == "!=" || operator == "!==", true
	}

	parsedValue := utils.ParseValue(valueStr)
	return utils.CompareValues(propValue, operator, parsedValue), true
}

func (f *FilterEvaluator) tryExistenceFilter(expr string, current interface{}) (bool, bool) {
	reExists := regexp.MustCompile(`^\.(\w+)$`)
	matches := reExists.FindStringSubmatch(expr)
	if len(matches) != 2 {
		return false, false
	}

	property := matches[1]
	obj, ok := current.(map[string]interface{})
	if !ok {
		return false, true
	}

	propValue, exists := obj[property]
	if !exists {
		return false, true
	}

	if propValue == nil {
		return false, true
	}

	// Check for empty arrays/objects and boolean values
	switch v := propValue.(type) {
	case []interface{}:
		return len(v) > 0, true
	case map[string]interface{}:
		return len(v) > 0, true
	case string:
		return v != "", true
	case bool:
		return v, true
	default:
		return true, true
	}
}

func (f *FilterEvaluator) tryFunctionPredicateFilter(expr string, current interface{}) (bool, bool) {
	// Handle .match(pattern)
	if result, ok := f.tryMatchFunction(expr, current); ok {
		return result, true
	}

	// Handle .contains(substring)
	if result, ok := f.tryContainsFunction(expr, current); ok {
		return result, true
	}

	// Handle .startsWith(prefix)
	if result, ok := f.tryStartsWithFunction(expr, current); ok {
		return result, true
	}

	// Handle .endsWith(suffix)
	if result, ok := f.tryEndsWithFunction(expr, current); ok {
		return result, true
	}

	// Handle .length comparisons
	if result, ok := f.tryLengthFunction(expr, current); ok {
		return result, true
	}

	// Handle .toLowerCase() and .toUpperCase()
	if result, ok := f.tryCaseFunction(expr, current); ok {
		return result, true
	}

	// Handle chained operations (e.g., .toLowerCase().contains())
	if result, ok := f.tryChainedOperations(expr, current); ok {
		return result, true
	}

	// Handle .typeof() comparisons
	if result, ok := f.tryTypeofFunction(expr, current); ok {
		return result, true
	}

	return false, false
}

// Function predicate implementations (delegated to utils for reusability)

func (f *FilterEvaluator) tryMatchFunction(expr string, current interface{}) (bool, bool) {
	return utils.TryMatchFunction(expr, current)
}

func (f *FilterEvaluator) tryContainsFunction(expr string, current interface{}) (bool, bool) {
	return utils.TryContainsFunction(expr, current)
}

func (f *FilterEvaluator) tryStartsWithFunction(expr string, current interface{}) (bool, bool) {
	return utils.TryStartsWithFunction(expr, current)
}

func (f *FilterEvaluator) tryEndsWithFunction(expr string, current interface{}) (bool, bool) {
	return utils.TryEndsWithFunction(expr, current)
}

func (f *FilterEvaluator) tryLengthFunction(expr string, current interface{}) (bool, bool) {
	return utils.TryLengthFunction(expr, current)
}

func (f *FilterEvaluator) tryCaseFunction(expr string, current interface{}) (bool, bool) {
	return utils.TryCaseFunction(expr, current)
}

func (f *FilterEvaluator) tryChainedOperations(expr string, current interface{}) (bool, bool) {
	return utils.TryChainedOperations(expr, current)
}

func (f *FilterEvaluator) tryTypeofFunction(expr string, current interface{}) (bool, bool) {
	return utils.TryTypeofFunction(expr, current)
}

// Helper functions

func (f *FilterEvaluator) findLogicalOperator(expr string, op string) int {
	inQuotes := false
	quoteChar := byte(0)
	parenDepth := 0

	for i := 0; i < len(expr)-len(op)+1; i++ {
		ch := expr[i]

		if !inQuotes {
			if ch == '\'' || ch == '"' {
				inQuotes = true
				quoteChar = ch
			} else if ch == '(' {
				parenDepth++
			} else if ch == ')' {
				parenDepth--
			} else if parenDepth == 0 && expr[i:i+len(op)] == op {
				// Make sure it's not part of a longer operator
				if i+len(op) < len(expr) {
					nextChar := expr[i+len(op)]
					if nextChar == '=' || nextChar == '&' || nextChar == '|' {
						continue
					}
				}
				return i
			}
		} else {
			if ch == quoteChar && (i == 0 || expr[i-1] != '\\') {
				inQuotes = false
				quoteChar = 0
			}
		}
	}

	return -1
}

// stripOuterParentheses removes outer parentheses if they wrap the entire expression
func (f *FilterEvaluator) stripOuterParentheses(expr string) string {
	expr = strings.TrimSpace(expr)

	if len(expr) < 2 || expr[0] != '(' || expr[len(expr)-1] != ')' {
		return expr
	}

	// Check if the parentheses actually wrap the entire expression
	// by tracking depth and ensuring we never close all parentheses before the end
	depth := 0
	for i, ch := range expr {
		if ch == '(' {
			depth++
		} else if ch == ')' {
			depth--
			// If depth hits 0 before the last character, the outer parens don't wrap everything
			if depth == 0 && i < len(expr)-1 {
				return expr
			}
		}
	}

	// If we get here, the outer parentheses wrap the entire expression
	return strings.TrimSpace(expr[1 : len(expr)-1])
}
