package operators

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/reclaimprotocol/jsonpathplus-go/pkg/types"
	"github.com/reclaimprotocol/jsonpathplus-go/pkg/utils"
)

// OperatorEvaluator handles special JSONPath-Plus operators
type OperatorEvaluator struct{}

// NewOperatorEvaluator creates a new operator evaluator
func NewOperatorEvaluator() *OperatorEvaluator {
	return &OperatorEvaluator{}
}

// EvaluatePropertyNames handles the property names operator (~)
func (o *OperatorEvaluator) EvaluatePropertyNames(ctx types.Result, options *types.Options) []types.Result {
	var results []types.Result

	switch v := ctx.Value.(type) {
	case *utils.OrderedMap:
		index := 0
		v.Range(func(key string, value interface{}) bool {
			results = append(results, types.Result{
				Value:          key,
				Path:           fmt.Sprintf("%s~[%d]", ctx.Path, index),
				Parent:         ctx.Value,
				ParentProperty: strconv.Itoa(index),
				Index:          index,
				OriginalIndex:  index,
			})
			index++
			return true
		})
	case map[string]interface{}:
		index := 0
		for key := range v {
			results = append(results, types.Result{
				Value:          key,
				Path:           fmt.Sprintf("%s~[%d]", ctx.Path, index),
				Parent:         ctx.Value,
				ParentProperty: strconv.Itoa(index),
				Index:          index,
				OriginalIndex:  index,
			})
			index++
		}
	case map[interface{}]interface{}:
		index := 0
		for key := range v {
			keyStr := fmt.Sprintf("%v", key)
			results = append(results, types.Result{
				Value:          keyStr,
				Path:           fmt.Sprintf("%s~[%d]", ctx.Path, index),
				Parent:         ctx.Value,
				ParentProperty: strconv.Itoa(index),
				Index:          index,
				OriginalIndex:  index,
			})
			index++
		}
	}

	return results
}

// EvaluateParent handles the parent operator (^)
func (o *OperatorEvaluator) EvaluateParent(ctx types.Result, options *types.Options) []types.Result {
	var results []types.Result

	// Return the parent object
	if ctx.Parent != nil {
		// Calculate parent path by removing the last segment
		parentPath := o.calculateParentPath(ctx.Path)

		results = append(results, types.Result{
			Value:          ctx.Parent,
			Path:           parentPath,
			Parent:         nil, // We don't track grandparents for now
			ParentProperty: "",
			Index:          0,
			OriginalIndex:  0,
		})
	}

	return results
}

// EvaluateParentWithDeduplication handles the parent operator with deduplication
func (o *OperatorEvaluator) EvaluateParentWithDeduplication(contexts []types.Result, options *types.Options) []types.Result {
	seen := make(map[string]bool)
	var results []types.Result

	for _, ctx := range contexts {
		if ctx.Parent != nil {
			parentPath := o.calculateParentPath(ctx.Path)

			// Only add if we haven't seen this parent path before
			if !seen[parentPath] {
				seen[parentPath] = true
				results = append(results, types.Result{
					Value:          ctx.Parent,
					Path:           parentPath,
					Parent:         nil,
					ParentProperty: "",
					Index:          0,
					OriginalIndex:  0,
				})
			}
		}
	}

	return results
}

// calculateParentPath calculates the parent path from a child path
func (o *OperatorEvaluator) calculateParentPath(childPath string) string {
	// Remove the last segment from the path
	// Examples:
	// "$.store.book[0]" -> "$.store.book"
	// "$.store.book[0].title" -> "$.store.book[0]"
	// "$.users[1].profile.bio" -> "$.users[1].profile"

	// Find the last meaningful separator
	lastDot := -1
	lastBracket := -1

	for i := len(childPath) - 1; i >= 0; i-- {
		if childPath[i] == '.' && lastDot == -1 {
			lastDot = i
		}
		if childPath[i] == ']' && lastBracket == -1 {
			// Find the matching opening bracket
			bracketDepth := 1
			for j := i - 1; j >= 0; j-- {
				if childPath[j] == ']' {
					bracketDepth++
				} else if childPath[j] == '[' {
					bracketDepth--
					if bracketDepth == 0 {
						lastBracket = j
						break
					}
				}
			}
			break
		}
	}

	// Determine which separator is more recent
	cutPoint := -1
	if lastDot > lastBracket {
		cutPoint = lastDot
	} else if lastBracket > -1 {
		cutPoint = lastBracket
	}

	if cutPoint > 0 {
		return childPath[:cutPoint]
	}

	// If we can't find a separator, return root
	return "$"
}

// EvaluateChainedOperations handles chained bracket operations
func (o *OperatorEvaluator) EvaluateChainedOperations(results []types.Result, chainNodes []*types.AstNode, evaluateNodeFunc func(*types.AstNode, []types.Result, *types.Options) []types.Result, options *types.Options) []types.Result {
	// Start with the initial results
	currentResults := results

	// Apply each chained operation in sequence
	for _, chainNode := range chainNodes {
		// Special handling for slice operations applied to multiple results
		if chainNode.Type == "slice" && len(currentResults) > 1 {
			// Apply slice to the collection of results, not individually
			nextResults := o.applySliceToResults(chainNode, currentResults)
			currentResults = nextResults
		} else {
			var nextResults []types.Result

			// Apply the chain node to each current result
			for _, result := range currentResults {
				nodeResults := evaluateNodeFunc(chainNode, []types.Result{result}, options)
				nextResults = append(nextResults, nodeResults...)
			}

			currentResults = nextResults
		}
	}

	return currentResults
}

// applySliceToResults applies a slice operation to a collection of results
func (o *OperatorEvaluator) applySliceToResults(sliceNode *types.AstNode, results []types.Result) []types.Result {
	slice := sliceNode.Value
	start, end, step := o.parseSliceParams(slice, len(results))

	var slicedResults []types.Result

	// Handle forward and reverse iteration
	if step > 0 {
		for i := start; i < end && i < len(results); i += step {
			if i >= 0 {
				slicedResults = append(slicedResults, results[i])
			}
		}
	} else {
		for i := start; i > end && i >= 0; i += step {
			if i < len(results) {
				slicedResults = append(slicedResults, results[i])
			}
		}
	}

	return slicedResults
}

// parseSliceParams parses slice parameters (similar to evaluator's parseSliceParams)
func (o *OperatorEvaluator) parseSliceParams(slice string, arrLen int) (start, end, step int) {
	slice = strings.TrimSpace(slice)
	parts := strings.Split(slice, ":")

	step = 1

	if len(parts) > 2 && parts[2] != "" {
		step, _ = strconv.Atoi(strings.TrimSpace(parts[2]))
	}

	if step == 0 {
		step = 1
	}

	// Set defaults based on step direction
	if step > 0 {
		start = 0
		end = arrLen
	} else {
		start = arrLen - 1
		end = -1
	}

	// Parse start if provided
	if len(parts) > 0 && parts[0] != "" {
		start, _ = strconv.Atoi(strings.TrimSpace(parts[0]))
	}

	// Parse end if provided
	if len(parts) > 1 && parts[1] != "" {
		end, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
	}

	// Handle negative indices
	if start < 0 {
		start = arrLen + start
	}
	if end < 0 && step > 0 {
		end = arrLen + end
	}

	// Clamp to valid range for forward iteration
	if step > 0 {
		if start < 0 {
			start = 0
		}
		if end > arrLen {
			end = arrLen
		}
	} else {
		// For reverse iteration, clamp differently
		if start >= arrLen {
			start = arrLen - 1
		}
		if start < 0 {
			start = 0
		}
	}

	return start, end, step
}

// ContextualEvaluator handles evaluations that need enhanced context
type ContextualEvaluator struct {
	operatorEval *OperatorEvaluator
}

// NewContextualEvaluator creates a new contextual evaluator
func NewContextualEvaluator() *ContextualEvaluator {
	return &ContextualEvaluator{
		operatorEval: NewOperatorEvaluator(),
	}
}

// CreateContext creates an enhanced context for evaluation
func (c *ContextualEvaluator) CreateContext(result types.Result, root interface{}) *types.Context {
	return types.NewContext(
		root,
		result.Value,
		result.Parent,
		result.ParentProperty,
		result.Path,
		result.Index,
	)
}

// CreateArrayElementContext creates a context for array elements with proper parent tracking
func (c *ContextualEvaluator) CreateArrayElementContext(result types.Result, root interface{}, actualArray interface{}) *types.Context {
	return types.NewArrayElementContext(
		root,
		result.Value,
		result.Parent,
		result.ParentProperty,
		result.Path,
		result.Index,
		actualArray, // The actual array containing this element
	)
}

// CreateChildContext creates a child context for nested evaluation
func (c *ContextualEvaluator) CreateChildContext(parent *types.Context, value interface{}, property string, index int) *types.Context {
	var newPath string
	if parent.Path == "$" {
		if property != "" {
			newPath = fmt.Sprintf("$.%s", property)
		} else {
			newPath = fmt.Sprintf("$[%d]", index)
		}
	} else {
		if property != "" {
			newPath = fmt.Sprintf("%s.%s", parent.Path, property)
		} else {
			newPath = fmt.Sprintf("%s[%d]", parent.Path, index)
		}
	}

	return types.NewContext(
		parent.Root,
		value,
		parent.Current,
		property,
		newPath,
		index,
	)
}

// EvaluateWithContext evaluates a node with enhanced context support
func (c *ContextualEvaluator) EvaluateWithContext(ctx *types.Context, evaluateFunc func(interface{}) []types.Result) []types.Result {
	// This is a helper that can be used by evaluators that need context
	return evaluateFunc(ctx.Current)
}
