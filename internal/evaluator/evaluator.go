package evaluator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/reclaimprotocol/jsonpathplus-go/internal/filters"
	"github.com/reclaimprotocol/jsonpathplus-go/internal/operators"
	"github.com/reclaimprotocol/jsonpathplus-go/pkg/types"
)

// Evaluator handles JSONPath expression evaluation
type Evaluator struct {
	filterEval     *filters.FilterEvaluator
	operatorEval   *operators.OperatorEvaluator
	contextualEval *operators.ContextualEvaluator
}

// NewEvaluator creates a new evaluator
func NewEvaluator() *Evaluator {
	return &Evaluator{
		filterEval:     filters.NewFilterEvaluator(),
		operatorEval:   operators.NewOperatorEvaluator(),
		contextualEval: operators.NewContextualEvaluator(),
	}
}

// Evaluate evaluates an AST against data
func (e *Evaluator) Evaluate(ast *types.AstNode, data interface{}, options *types.Options) []types.Result {
	if options == nil {
		options = &types.Options{}
	}

	// Set root for $ references
	if options.Root == nil {
		options.Root = data
	}

	rootResult := types.Result{
		Value:          data,
		Path:           "$",
		Parent:         nil,
		ParentProperty: "",
		Index:          0,
		OriginalIndex:  0,
	}

	return e.evaluateNode(ast, []types.Result{rootResult}, options)
}

// evaluateNode evaluates a single AST node
func (e *Evaluator) evaluateNode(node *types.AstNode, contexts []types.Result, options *types.Options) []types.Result {
	var results []types.Result

	// Special handling for filters applied to multiple contexts (e.g., after wildcard)
	if node.Type == "filter" && len(contexts) > 1 {
		return e.evaluateFilterOnResults(node, contexts, options)
	}

	for _, ctx := range contexts {
		nodeResults := e.evaluateSingleNode(node, ctx, options)
		results = append(results, nodeResults...)
	}

	return results
}

// evaluateFilterOnResults applies a filter to a collection of results
func (e *Evaluator) evaluateFilterOnResults(node *types.AstNode, contexts []types.Result, options *types.Options) []types.Result {
	var results []types.Result

	for _, ctx := range contexts {
		// Create enhanced context for filter evaluation
		itemContext := e.contextualEval.CreateContext(ctx, options.Root)

		if e.filterEval.EvaluateFilter(node.Value, itemContext) {
			if len(node.Children) > 0 {
				childResults := e.evaluateNode(node.Children[0], []types.Result{ctx}, options)
				results = append(results, childResults...)
			} else {
				results = append(results, ctx)
			}
		}
	}

	return results
}

// evaluateSingleNode evaluates a node against a single context
func (e *Evaluator) evaluateSingleNode(node *types.AstNode, ctx types.Result, options *types.Options) []types.Result {
	switch node.Type {
	case "root":
		return e.evaluateRoot(node, ctx, options)
	case "property":
		return e.evaluateProperty(node, ctx, options)
	case "wildcard":
		return e.evaluateWildcard(node, ctx, options)
	case "index_wildcard":
		return e.evaluateIndexWildcard(node, ctx, options)
	case "index":
		return e.evaluateIndex(node, ctx, options)
	case "slice":
		return e.evaluateSlice(node, ctx, options)
	case "filter":
		return e.evaluateFilter(node, ctx, options)
	case "recursive":
		return e.evaluateRecursive(node, ctx, options)
	case "union":
		return e.evaluateUnion(node, ctx, options)
	case "chain":
		return e.evaluateChain(node, ctx, options)
	case "property_names":
		return e.evaluatePropertyNames(node, ctx, options)
	case "parent":
		return e.evaluateParent(node, ctx, options)
	default:
		return nil
	}
}

// Node type evaluators

func (e *Evaluator) evaluateRoot(node *types.AstNode, ctx types.Result, options *types.Options) []types.Result {
	if len(node.Children) == 0 {
		return []types.Result{ctx}
	}

	return e.evaluateNode(node.Children[0], []types.Result{ctx}, options)
}

func (e *Evaluator) evaluateProperty(node *types.AstNode, ctx types.Result, options *types.Options) []types.Result {
	var results []types.Result
	property := node.Value

	switch v := ctx.Value.(type) {
	case map[string]interface{}:
		if value, exists := v[property]; exists {
			result := types.Result{
				Value:          value,
				Path:           fmt.Sprintf("%s.%s", ctx.Path, property),
				Parent:         ctx.Value,
				ParentProperty: property,
				Index:          0,
				OriginalIndex:  0,
			}

			if len(node.Children) > 0 {
				return e.evaluateNode(node.Children[0], []types.Result{result}, options)
			}
			results = append(results, result)
		}
	case []interface{}:
		// For arrays, treat property as index if it's numeric
		if idx, err := strconv.Atoi(property); err == nil && idx >= 0 && idx < len(v) {
			result := types.Result{
				Value:          v[idx],
				Path:           fmt.Sprintf("%s[%d]", ctx.Path, idx),
				Parent:         ctx.Value,
				ParentProperty: strconv.Itoa(idx),
				Index:          idx,
				OriginalIndex:  idx,
			}

			if len(node.Children) > 0 {
				return e.evaluateNode(node.Children[0], []types.Result{result}, options)
			}
			results = append(results, result)
		}
	}

	return results
}

func (e *Evaluator) evaluateWildcard(node *types.AstNode, ctx types.Result, options *types.Options) []types.Result {
	var results []types.Result

	switch v := ctx.Value.(type) {
	case map[string]interface{}:
		index := 0
		for key, value := range v {
			result := types.Result{
				Value:          value,
				Path:           fmt.Sprintf("%s.%s", ctx.Path, key),
				Parent:         ctx.Value,
				ParentProperty: key,
				Index:          index,
				OriginalIndex:  index,
			}
			results = append(results, result)
			index++
		}
	case []interface{}:
		// Check if this is a property wildcard (from dot notation) vs index wildcard (from bracket notation)
		// Property wildcards on arrays should return properties of array elements
		// Index wildcards on arrays should return the array elements themselves
		isPropertyWildcard := ctx.Path != "" && !strings.HasSuffix(ctx.Path, "]")

		if isPropertyWildcard {
			// Property wildcard: $.store.book.* should return all properties of all books
			for i, value := range v {
				if valueMap, ok := value.(map[string]interface{}); ok {
					// For each object in the array, add all its properties
					propIndex := 0
					for key, propValue := range valueMap {
						result := types.Result{
							Value:          propValue,
							Path:           fmt.Sprintf("%s[%d].%s", ctx.Path, i, key),
							Parent:         value,
							ParentProperty: key,
							Index:          propIndex,
							OriginalIndex:  propIndex,
						}
						results = append(results, result)
						propIndex++
					}
				} else {
					// For non-object array elements, return the element itself
					result := types.Result{
						Value:          value,
						Path:           fmt.Sprintf("%s[%d]", ctx.Path, i),
						Parent:         ctx.Value,
						ParentProperty: strconv.Itoa(i),
						Index:          i,
						OriginalIndex:  i,
					}
					results = append(results, result)
				}
			}
		} else {
			// Index wildcard: $.store.book[*] should return array elements
			for i, value := range v {
				result := types.Result{
					Value:          value,
					Path:           fmt.Sprintf("%s[%d]", ctx.Path, i),
					Parent:         ctx.Value,
					ParentProperty: strconv.Itoa(i),
					Index:          i,
					OriginalIndex:  i,
				}
				results = append(results, result)
			}
		}
	}

	// Apply children to all results at once (important for filters)
	if len(node.Children) > 0 {
		return e.evaluateNode(node.Children[0], results, options)
	}

	return results
}

func (e *Evaluator) evaluateIndexWildcard(node *types.AstNode, ctx types.Result, options *types.Options) []types.Result {
	var results []types.Result

	switch v := ctx.Value.(type) {
	case map[string]interface{}:
		// For objects, index wildcard behaves like property wildcard
		index := 0
		for key, value := range v {
			result := types.Result{
				Value:          value,
				Path:           fmt.Sprintf("%s.%s", ctx.Path, key),
				Parent:         ctx.Value,
				ParentProperty: key,
				Index:          index,
				OriginalIndex:  index,
			}
			results = append(results, result)
			index++
		}
	case []interface{}:
		// For arrays, index wildcard returns array elements themselves
		for i, value := range v {
			result := types.Result{
				Value:          value,
				Path:           fmt.Sprintf("%s[%d]", ctx.Path, i),
				Parent:         ctx.Value,
				ParentProperty: strconv.Itoa(i),
				Index:          i,
				OriginalIndex:  i,
			}
			results = append(results, result)
		}
	}

	// Apply children to all results at once (important for filters)
	if len(node.Children) > 0 {
		return e.evaluateNode(node.Children[0], results, options)
	}

	return results
}

func (e *Evaluator) evaluateIndex(node *types.AstNode, ctx types.Result, options *types.Options) []types.Result {
	var results []types.Result

	arr, ok := ctx.Value.([]interface{})
	if !ok {
		return results
	}

	idx, err := strconv.Atoi(node.Value)
	if err != nil {
		return results
	}

	// Handle negative indices
	if idx < 0 {
		idx = len(arr) + idx
	}

	if idx >= 0 && idx < len(arr) {
		result := types.Result{
			Value:          arr[idx],
			Path:           fmt.Sprintf("%s[%d]", ctx.Path, idx),
			Parent:         ctx.Value,
			ParentProperty: strconv.Itoa(idx),
			Index:          idx,
			OriginalIndex:  idx,
		}

		if len(node.Children) > 0 {
			return e.evaluateNode(node.Children[0], []types.Result{result}, options)
		}
		results = append(results, result)
	}

	return results
}

func (e *Evaluator) evaluateSlice(node *types.AstNode, ctx types.Result, options *types.Options) []types.Result {
	var results []types.Result

	arr, ok := ctx.Value.([]interface{})
	if !ok {
		return results
	}

	start, end, step := e.parseSliceParams(node.Value, len(arr))

	// Handle forward and reverse iteration
	if step > 0 {
		for i := start; i < end && i < len(arr); i += step {
			if i >= 0 {
				result := types.Result{
					Value:          arr[i],
					Path:           fmt.Sprintf("%s[%d]", ctx.Path, i),
					Parent:         ctx.Value,
					ParentProperty: strconv.Itoa(i),
					Index:          len(results),
					OriginalIndex:  i,
				}

				if len(node.Children) > 0 {
					childResults := e.evaluateNode(node.Children[0], []types.Result{result}, options)
					results = append(results, childResults...)
				} else {
					results = append(results, result)
				}
			}
		}
	} else {
		for i := start; i > end && i >= 0; i += step {
			if i < len(arr) {
				result := types.Result{
					Value:          arr[i],
					Path:           fmt.Sprintf("%s[%d]", ctx.Path, i),
					Parent:         ctx.Value,
					ParentProperty: strconv.Itoa(i),
					Index:          len(results),
					OriginalIndex:  i,
				}

				if len(node.Children) > 0 {
					childResults := e.evaluateNode(node.Children[0], []types.Result{result}, options)
					results = append(results, childResults...)
				} else {
					results = append(results, result)
				}
			}
		}
	}

	return results
}

func (e *Evaluator) evaluateFilter(node *types.AstNode, ctx types.Result, options *types.Options) []types.Result {
	var results []types.Result

	// Handle array filtering
	if arr, ok := ctx.Value.([]interface{}); ok {
		for i, item := range arr {
			// Create enhanced context for filter evaluation
			// For @parentProperty, use the property name that led to this array
			parentProp := ctx.ParentProperty
			if parentProp == "" {
				parentProp = strconv.Itoa(i)
			}

			itemResult := types.Result{
				Value:          item,
				Path:           fmt.Sprintf("%s[%d]", ctx.Path, i),
				Parent:         ctx.Parent, // Use the parent of the array, not the array itself
				ParentProperty: parentProp, // Use the property that contains the array
				Index:          i,
				OriginalIndex:  i,
			}

			itemContext := e.contextualEval.CreateContext(itemResult, options.Root)

			if e.filterEval.EvaluateFilter(node.Value, itemContext) {
				if len(node.Children) > 0 {
					childResults := e.evaluateNode(node.Children[0], []types.Result{itemResult}, options)
					results = append(results, childResults...)
				} else {
					results = append(results, itemResult)
				}
			}
		}
	} else {
		// Handle single item filtering (for chained operations)
		// Create context for the single item
		itemContext := e.contextualEval.CreateContext(ctx, options.Root)

		if e.filterEval.EvaluateFilter(node.Value, itemContext) {
			if len(node.Children) > 0 {
				childResults := e.evaluateNode(node.Children[0], []types.Result{ctx}, options)
				results = append(results, childResults...)
			} else {
				results = append(results, ctx)
			}
		}
	}

	return results
}

func (e *Evaluator) evaluateRecursive(node *types.AstNode, ctx types.Result, options *types.Options) []types.Result {
	var results []types.Result
	visited := make(map[string]bool)

	// If we have children, we need to find all nodes that match the child criteria
	if len(node.Children) > 0 {
		var allNodes []types.Result

		// First, collect all nodes recursively
		var traverse func(current types.Result)
		traverse = func(current types.Result) {
			if visited[current.Path] {
				return
			}
			visited[current.Path] = true

			// Add current node to potential matches
			allNodes = append(allNodes, current)

			// Recursively traverse children
			switch v := current.Value.(type) {
			case map[string]interface{}:
				for key, val := range v {
					childResult := types.Result{
						Value:          val,
						Path:           fmt.Sprintf("%s.%s", current.Path, key),
						Parent:         current.Value,
						ParentProperty: key,
						Index:          0,
						OriginalIndex:  0,
					}
					traverse(childResult)
				}
			case []interface{}:
				for i, val := range v {
					childResult := types.Result{
						Value:          val,
						Path:           fmt.Sprintf("%s[%d]", current.Path, i),
						Parent:         current.Value,
						ParentProperty: strconv.Itoa(i),
						Index:          i,
						OriginalIndex:  i,
					}
					traverse(childResult)
				}
			}
		}

		traverse(ctx)

		// Now apply the child node to each collected node
		for _, nodeResult := range allNodes {
			childResults := e.evaluateNode(node.Children[0], []types.Result{nodeResult}, options)
			results = append(results, childResults...)
		}
	} else {
		// No children - return all nodes at all levels (this case shouldn't happen with ..)
		var traverse func(current types.Result)
		traverse = func(current types.Result) {
			if visited[current.Path] {
				return
			}
			visited[current.Path] = true

			switch v := current.Value.(type) {
			case map[string]interface{}:
				for key, val := range v {
					childResult := types.Result{
						Value:          val,
						Path:           fmt.Sprintf("%s.%s", current.Path, key),
						Parent:         current.Value,
						ParentProperty: key,
						Index:          0,
						OriginalIndex:  0,
					}
					results = append(results, childResult)
					traverse(childResult)
				}
			case []interface{}:
				for i, val := range v {
					childResult := types.Result{
						Value:          val,
						Path:           fmt.Sprintf("%s[%d]", current.Path, i),
						Parent:         current.Value,
						ParentProperty: strconv.Itoa(i),
						Index:          i,
						OriginalIndex:  i,
					}
					results = append(results, childResult)
					traverse(childResult)
				}
			}
		}

		traverse(ctx)
	}

	return results
}

func (e *Evaluator) evaluateUnion(node *types.AstNode, ctx types.Result, options *types.Options) []types.Result {
	var results []types.Result

	for _, child := range node.Children {
		childResults := e.evaluateSingleNode(child, ctx, options)
		results = append(results, childResults...)
	}

	return results
}

func (e *Evaluator) evaluateChain(node *types.AstNode, ctx types.Result, options *types.Options) []types.Result {
	if len(node.Children) == 0 {
		return []types.Result{ctx}
	}

	// Start with the first operation
	currentResults := e.evaluateSingleNode(node.Children[0], ctx, options)

	// Apply subsequent operations to the results
	return e.operatorEval.EvaluateChainedOperations(
		currentResults,
		node.Children[1:],
		func(n *types.AstNode, results []types.Result, opts *types.Options) []types.Result {
			return e.evaluateNode(n, results, opts)
		},
		options,
	)
}

func (e *Evaluator) evaluatePropertyNames(node *types.AstNode, ctx types.Result, options *types.Options) []types.Result {
	results := e.operatorEval.EvaluatePropertyNames(ctx, options)

	if len(node.Children) > 0 {
		return e.evaluateNode(node.Children[0], results, options)
	}

	return results
}

func (e *Evaluator) evaluateParent(node *types.AstNode, ctx types.Result, options *types.Options) []types.Result {
	// If there are children, evaluate them first to get the target results,
	// then apply the parent operator to those results
	if len(node.Children) > 0 {
		childResults := e.evaluateNode(node.Children[0], []types.Result{ctx}, options)

		// Use deduplication when we have multiple child results
		if len(childResults) > 1 {
			return e.operatorEval.EvaluateParentWithDeduplication(childResults, options)
		}

		// Single result - no need for deduplication
		var results []types.Result
		for _, childResult := range childResults {
			parentResults := e.operatorEval.EvaluateParent(childResult, options)
			results = append(results, parentResults...)
		}

		return results
	}

	// No children - apply parent operator directly to current context
	return e.operatorEval.EvaluateParent(ctx, options)
}

// parseSliceParams parses slice parameters with support for reverse iteration
func (e *Evaluator) parseSliceParams(slice string, arrLen int) (start, end, step int) {
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
