package jsonpathplus

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	nodeTypeRecursive = "recursive"
)

// Result represents a JSONPath query result with position information.
type Result struct {
	Value          interface{} // The actual value
	Path           string      // JSONPath to this element
	Parent         interface{} // Reference to parent object/array
	ParentProperty string      // Property name or array index in parent
	Index          int         // Position in result set
	OriginalIndex  int         // Character position in original JSON string
	Length         int         // Length of the element in the JSON string
}

// Options configures JSONPath evaluation behavior.
type Options struct {
	ResultType string
	Flatten    bool
	Wrap       bool
}

// JSONPath represents a compiled JSONPath expression.
type JSONPath struct {
	path    string
	options *Options
}

// New creates a new JSONPath instance with the given path and options.
func New(path string, options *Options) *JSONPath {
	if options == nil {
		options = &Options{
			ResultType: "value",
			Flatten:    false,
			Wrap:       true,
		}
	}
	return &JSONPath{
		path:    path,
		options: options,
	}
}

// Query executes a JSONPath query with string character position tracking.
func Query(path string, jsonStr string) ([]Result, error) {
	return QueryWithStringIndex(path, jsonStr)
}

// QueryData executes a JSONPath query on parsed data (legacy support).
func QueryData(path string, data interface{}) ([]Result, error) {
	jp := New(path, nil)
	return jp.Execute(data)
}

// Execute evaluates the JSONPath against the provided data.
func (jp *JSONPath) Execute(data interface{}) ([]Result, error) {
	tokens, err := tokenize(jp.path)
	if err != nil {
		return nil, err
	}

	ast, err := parse(tokens)
	if err != nil {
		return nil, err
	}

	return evaluate(ast, data, jp.options)
}

type tokenType int

const (
	tokenRoot tokenType = iota
	tokenCurrent
	tokenDot
	tokenDoubleDot
	tokenBracketOpen
	tokenBracketClose
	tokenIdentifier
	tokenNumber
	tokenString
	tokenWildcard
	tokenFilter
	tokenSlice
	tokenComma
	tokenUnion
)

type token struct {
	Type  tokenType
	Value string
}

func tokenize(path string) ([]token, error) {
	var tokens []token
	i := 0
	inBracket := false
	bracketStart := 0

	for i < len(path) {
		switch path[i] {
		case '$':
			tokens = append(tokens, token{Type: tokenRoot, Value: "$"})
			i++
		case '@':
			tokens = append(tokens, token{Type: tokenCurrent, Value: "@"})
			i++
		case '.':
			if i+1 < len(path) && path[i+1] == '.' {
				tokens = append(tokens, token{Type: tokenDoubleDot, Value: ".."})
				i += 2
			} else {
				tokens = append(tokens, token{Type: tokenDot, Value: "."})
				i++
			}
		case '[':
			tokens = append(tokens, token{Type: tokenBracketOpen, Value: "["})
			inBracket = true
			bracketStart = i + 1
			i++
		case ']':
			tokens = append(tokens, token{Type: tokenBracketClose, Value: "]"})
			inBracket = false
			i++
		case '*':
			tokens = append(tokens, token{Type: tokenWildcard, Value: "*"})
			i++
		case ',':
			tokens = append(tokens, token{Type: tokenComma, Value: ","})
			i++
		case '?':
			filterEnd := findFilterEnd(path, i)
			if filterEnd == -1 {
				return nil, fmt.Errorf("unclosed filter expression")
			}
			tokens = append(tokens, token{Type: tokenFilter, Value: path[i:filterEnd]})
			i = filterEnd
		case '\'', '"':
			strEnd := findStringEnd(path, i)
			if strEnd == -1 {
				return nil, fmt.Errorf("unclosed string")
			}
			tokens = append(tokens, token{Type: tokenString, Value: path[i+1 : strEnd]})
			i = strEnd + 1
		case ':':
			if inBracket {
				// We need to handle the slice token properly
				// Remove the last token if it was a number (part of the slice)
				if len(tokens) > 0 && tokens[len(tokens)-1].Type == tokenNumber {
					tokens = tokens[:len(tokens)-1]
				}
				sliceEnd := findSliceEnd(path, bracketStart)
				tokens = append(tokens, token{Type: tokenSlice, Value: path[bracketStart:sliceEnd]})
				i = sliceEnd
				if i < len(path) && path[i] == ']' {
					continue
				}
			} else {
				i++
			}
		case ' ', '\t', '\n', '\r':
			i++
		default:
			switch {
			case isDigit(path[i]) || path[i] == '-':
				numEnd := findNumberEnd(path, i)
				tokens = append(tokens, token{Type: tokenNumber, Value: path[i:numEnd]})
				i = numEnd
			case isIdentifierStart(path[i]):
				idEnd := findIdentifierEnd(path, i)
				tokens = append(tokens, token{Type: tokenIdentifier, Value: path[i:idEnd]})
				i = idEnd
			default:
				return nil, fmt.Errorf("unexpected character: %c at position %d", path[i], i)
			}
		}
	}

	return tokens, nil
}

func findFilterEnd(path string, start int) int {
	depth := 0
	inString := false
	stringChar := byte(0)

	for i := start; i < len(path); i++ {
		if inString {
			if path[i] == stringChar && (i == 0 || path[i-1] != '\\') {
				inString = false
			}
			continue
		}

		switch path[i] {
		case '\'', '"':
			inString = true
			stringChar = path[i]
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i + 1
			}
		}
	}
	return -1
}

func findStringEnd(path string, start int) int {
	quote := path[start]
	for i := start + 1; i < len(path); i++ {
		if path[i] == quote && (i == start+1 || path[i-1] != '\\') {
			return i
		}
	}
	return -1
}

func findSliceEnd(path string, start int) int {
	for i := start; i < len(path); i++ {
		if path[i] == ']' {
			return i
		}
	}
	return len(path)
}

func findNumberEnd(path string, start int) int {
	i := start
	if i < len(path) && path[i] == '-' {
		i++
	}
	for i < len(path) && isDigit(path[i]) {
		i++
	}
	return i
}

func findIdentifierEnd(path string, start int) int {
	i := start
	for i < len(path) && (isIdentifierChar(path[i])) {
		i++
	}
	return i
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isIdentifierStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isIdentifierChar(ch byte) bool {
	return isIdentifierStart(ch) || isDigit(ch)
}

type astNode struct {
	Type     string
	Value    string
	Children []*astNode
}

func parse(tokens []token) (*astNode, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty path")
	}

	root := &astNode{Type: "path"}
	current := root
	i := 0

	for i < len(tokens) {
		switch tokens[i].Type {
		case tokenRoot:
			node := &astNode{Type: "root", Value: "$"}
			current.Children = append(current.Children, node)
			i++
		case tokenCurrent:
			node := &astNode{Type: "current", Value: "@"}
			current.Children = append(current.Children, node)
			i++
		case tokenDot:
			i++
			if i < len(tokens) {
				switch tokens[i].Type {
				case tokenIdentifier:
					node := &astNode{Type: "property", Value: tokens[i].Value}
					current.Children = append(current.Children, node)
					i++
				case tokenWildcard:
					node := &astNode{Type: "wildcard", Value: "*"}
					current.Children = append(current.Children, node)
					i++
				case tokenRoot, tokenCurrent, tokenDot, tokenDoubleDot,
					tokenBracketOpen, tokenBracketClose, tokenNumber,
					tokenString, tokenFilter, tokenSlice, tokenComma, tokenUnion:
					// Skip other token types after dot
				}
			}
		case tokenDoubleDot:
			node := &astNode{Type: "recursive", Value: ".."}
			current.Children = append(current.Children, node)
			i++
			if i < len(tokens) {
				switch tokens[i].Type {
				case tokenIdentifier:
					childNode := &astNode{Type: "property", Value: tokens[i].Value}
					node.Children = append(node.Children, childNode)
					i++
				case tokenWildcard:
					childNode := &astNode{Type: "wildcard", Value: "*"}
					node.Children = append(node.Children, childNode)
					i++
				case tokenRoot, tokenCurrent, tokenDot, tokenDoubleDot,
					tokenBracketOpen, tokenBracketClose, tokenNumber,
					tokenString, tokenFilter, tokenSlice, tokenComma, tokenUnion:
					// Skip other token types after double dot
				}
			}
		case tokenBracketOpen:
			i++
			bracketNode := &astNode{Type: "bracket"}
			for i < len(tokens) && tokens[i].Type != tokenBracketClose {
				switch tokens[i].Type {
				case tokenNumber:
					indexNode := &astNode{Type: "index", Value: tokens[i].Value}
					bracketNode.Children = append(bracketNode.Children, indexNode)
					i++
				case tokenString:
					propNode := &astNode{Type: "property", Value: tokens[i].Value}
					bracketNode.Children = append(bracketNode.Children, propNode)
					i++
				case tokenWildcard:
					wildcardNode := &astNode{Type: "wildcard", Value: "*"}
					bracketNode.Children = append(bracketNode.Children, wildcardNode)
					i++
				case tokenFilter:
					filterNode := &astNode{Type: "filter", Value: tokens[i].Value}
					bracketNode.Children = append(bracketNode.Children, filterNode)
					i++
				case tokenSlice:
					sliceNode := &astNode{Type: "slice", Value: tokens[i].Value}
					bracketNode.Children = append(bracketNode.Children, sliceNode)
					i++
				case tokenComma:
					i++
				case tokenRoot, tokenCurrent, tokenDot, tokenDoubleDot,
					tokenBracketOpen, tokenBracketClose, tokenIdentifier, tokenUnion:
					i++
				}
			}
			if i < len(tokens) && tokens[i].Type == tokenBracketClose {
				i++
			}
			current.Children = append(current.Children, bracketNode)
		case tokenBracketClose, tokenIdentifier, tokenNumber, tokenString,
			tokenWildcard, tokenFilter, tokenSlice, tokenComma, tokenUnion:
			i++
		}
	}

	return root, nil
}

func evaluate(ast *astNode, data interface{}, options *Options) ([]Result, error) {
	results := []Result{{
		Value:         data,
		Path:          "$",
		Index:         0,
		OriginalIndex: 0,
	}}

	for _, child := range ast.Children {
		results = evaluateNode(child, results, options)
	}

	if options.ResultType == "path" {
		for i := range results {
			results[i].Value = results[i].Path
		}
	}

	return results, nil
}

func evaluateNode(node *astNode, contexts []Result, options *Options) []Result {
	var results []Result

	for _, ctx := range contexts {
		switch node.Type {
		case "root", "current":
			results = append(results, ctx)
		case "property":
			results = append(results, evaluateProperty(node, ctx)...)
		case "wildcard":
			results = append(results, evaluateWildcard(ctx)...)
		case "index":
			results = append(results, evaluateIndex(node, ctx)...)
		case "bracket":
			for _, child := range node.Children {
				results = append(results, evaluateNode(child, []Result{ctx}, options)...)
			}
		case "filter":
			results = append(results, evaluateFilter(node.Value, ctx, options)...)
		case "slice":
			results = append(results, evaluateSlice(node.Value, ctx, options)...)
		case nodeTypeRecursive:
			results = append(results, evaluateRecursive(node, ctx, options)...)
		default:
			// Unknown node type, skip
		}
	}

	return results
}

func evaluateProperty(node *astNode, ctx Result) []Result {
	var results []Result
	obj, ok := ctx.Value.(map[string]interface{})
	if !ok {
		return results
	}
	
	val, exists := obj[node.Value]
	if !exists {
		return results
	}

	results = append(results, Result{
		Value:          val,
		Path:           ctx.Path + "." + node.Value,
		Parent:         ctx.Value,
		ParentProperty: node.Value,
		Index:          len(results),
		OriginalIndex:  getOriginalIndex(ctx.Value, node.Value),
	})
	return results
}

func evaluateWildcard(ctx Result) []Result {
	var results []Result
	
	switch v := ctx.Value.(type) {
	case map[string]interface{}:
		idx := 0
		for key, val := range v {
			results = append(results, Result{
				Value:          val,
				Path:           ctx.Path + "." + key,
				Parent:         ctx.Value,
				ParentProperty: key,
				Index:          idx,
				OriginalIndex:  getOriginalIndex(ctx.Value, key),
			})
			idx++
		}
	case []interface{}:
		for i, val := range v {
			results = append(results, Result{
				Value:          val,
				Path:           fmt.Sprintf("%s[%d]", ctx.Path, i),
				Parent:         ctx.Value,
				ParentProperty: strconv.Itoa(i),
				Index:          i,
				OriginalIndex:  i,
			})
		}
	default:
		// Other types don't support wildcard
	}
	return results
}

func evaluateIndex(node *astNode, ctx Result) []Result {
	var results []Result
	arr, ok := ctx.Value.([]interface{})
	if !ok {
		return results
	}
	
	idx, _ := strconv.Atoi(node.Value)
	if idx < 0 {
		idx = len(arr) + idx
	}
	if idx < 0 || idx >= len(arr) {
		return results
	}
	
	results = append(results, Result{
		Value:          arr[idx],
		Path:           fmt.Sprintf("%s[%d]", ctx.Path, idx),
		Parent:         ctx.Value,
		ParentProperty: strconv.Itoa(idx),
		Index:          idx,
		OriginalIndex:  idx,
	})
	return results
}

func evaluateFilter(filter string, ctx Result, _ *Options) []Result {
	var results []Result

	filter = strings.TrimPrefix(filter, "?(")
	filter = strings.TrimSuffix(filter, ")")

	v, ok := ctx.Value.([]interface{})
	if !ok {
		return results
	}
	
	for i, item := range v {
		if evaluateFilterExpression(filter, item, ctx.Value) {
			results = append(results, Result{
				Value:          item,
				Path:           fmt.Sprintf("%s[%d]", ctx.Path, i),
				Parent:         ctx.Value,
				ParentProperty: strconv.Itoa(i),
				Index:          i,
				OriginalIndex:  i,
			})
		}
	}

	return results
}

func evaluateFilterExpression(expr string, current, root interface{}) bool {
	expr = cleanFilterExpression(expr)
	
	// Try comparison expression first
	if result, ok := tryComparisonFilter(expr, current); ok {
		return result
	}
	
	// Try existence check
	if result, ok := tryExistenceFilter(expr, current); ok {
		return result
	}
	
	return false
}

func cleanFilterExpression(expr string) string {
	expr = strings.TrimSpace(expr)
	expr = strings.TrimPrefix(expr, "@")
	return strings.TrimSpace(expr)
}

func tryComparisonFilter(expr string, current interface{}) (bool, bool) {
	re := regexp.MustCompile(`\.(\w+)\s*(<=|>=|==|!=|<|>)\s*(.+)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 4 { // Full match + 3 captured groups
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
		return operator == "!=", true
	}
	
	parsedValue := parseValue(valueStr)
	return compareValues(propValue, operator, parsedValue), true
}

func tryExistenceFilter(expr string, current interface{}) (bool, bool) {
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
	
	return propValue != nil && propValue != false, true
}

func compareValues(left interface{}, operator string, right interface{}) bool {
	switch operator {
	case "==":
		return reflect.DeepEqual(left, right)
	case "!=":
		return !reflect.DeepEqual(left, right)
	case "<", ">", "<=", ">=":
		leftNum, leftOk := toNumber(left)
		rightNum, rightOk := toNumber(right)
		if leftOk && rightOk {
			switch operator {
			case "<":
				return leftNum < rightNum
			case ">":
				return leftNum > rightNum
			case "<=":
				return leftNum <= rightNum
			case ">=":
				return leftNum >= rightNum
			}
		}
	}
	return false
}

func toNumber(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case string:
		if num, err := strconv.ParseFloat(val, 64); err == nil {
			return num, true
		}
	}
	return 0, false
}

func parseValue(s string) interface{} {
	s = strings.TrimSpace(s)

	if strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'") {
		return s[1 : len(s)-1]
	}
	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		return s[1 : len(s)-1]
	}

	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}
	if s == "null" {
		return nil
	}

	if num, err := strconv.ParseFloat(s, 64); err == nil {
		return num
	}

	return s
}

func evaluateSlice(slice string, ctx Result, _ *Options) []Result {
	var results []Result
	
	arr, ok := ctx.Value.([]interface{})
	if !ok {
		return results
	}

	start, end, step := parseSliceParams(slice, len(arr))
	
	for i := start; i < end && i < len(arr); i += step {
		if i >= 0 {
			results = append(results, Result{
				Value:          arr[i],
				Path:           fmt.Sprintf("%s[%d]", ctx.Path, i),
				Parent:         ctx.Value,
				ParentProperty: strconv.Itoa(i),
				Index:          len(results),
				OriginalIndex:  i,
			})
		}
	}

	return results
}

func parseSliceParams(slice string, arrLen int) (start, end, step int) {
	slice = strings.TrimSpace(slice)
	parts := strings.Split(slice, ":")
	
	start = 0
	end = arrLen
	step = 1

	if len(parts) > 0 && parts[0] != "" {
		start, _ = strconv.Atoi(strings.TrimSpace(parts[0]))
	}

	if len(parts) > 1 && parts[1] != "" {
		end, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
	}

	if len(parts) > 2 && parts[2] != "" {
		step, _ = strconv.Atoi(strings.TrimSpace(parts[2]))
	}

	if step == 0 {
		step = 1
	}

	// Handle negative indices
	if start < 0 {
		start = arrLen + start
	}
	if end < 0 {
		end = arrLen + end
	}

	// Clamp to valid range
	if start < 0 {
		start = 0
	}
	if end > arrLen {
		end = arrLen
	}
	
	return start, end, step
}

func evaluateRecursive(node *astNode, ctx Result, options *Options) []Result {
	var results []Result
	visited := make(map[string]bool)

	var traverse func(current Result)
	traverse = func(current Result) {
		if visited[current.Path] {
			return
		}
		visited[current.Path] = true

		if len(node.Children) > 0 {
			childResults := evaluateNode(node.Children[0], []Result{current}, options)
			results = append(results, childResults...)
		} else {
			switch v := current.Value.(type) {
			case map[string]interface{}:
				for key, val := range v {
					results = append(results, Result{
						Value:          val,
						Path:           current.Path + "." + key,
						Parent:         current.Value,
						ParentProperty: key,
						Index:          len(results),
						OriginalIndex:  getOriginalIndex(current.Value, key),
					})
				}
			case []interface{}:
				for i, val := range v {
					results = append(results, Result{
						Value:          val,
						Path:           fmt.Sprintf("%s[%d]", current.Path, i),
						Parent:         current.Value,
						ParentProperty: strconv.Itoa(i),
						Index:          i,
						OriginalIndex:  i,
					})
				}
			default:
				if len(node.Children) == 0 {
					results = append(results, current)
				}
			}
		}

		switch v := current.Value.(type) {
		case map[string]interface{}:
			for key, val := range v {
				traverse(Result{
					Value:          val,
					Path:           current.Path + "." + key,
					Parent:         current.Value,
					ParentProperty: key,
					Index:          len(results),
					OriginalIndex:  getOriginalIndex(current.Value, key),
				})
			}
		case []interface{}:
			for i, val := range v {
				traverse(Result{
					Value:          val,
					Path:           fmt.Sprintf("%s[%d]", current.Path, i),
					Parent:         current.Value,
					ParentProperty: strconv.Itoa(i),
					Index:          i,
					OriginalIndex:  i,
				})
			}
		default:
			// Leaf node or unsupported type, no traversal needed
		}
	}

	traverse(ctx)
	return results
}

func getOriginalIndex(data interface{}, key string) int {
	switch v := data.(type) {
	case map[string]interface{}:
		idx := 0
		for k := range v {
			if k == key {
				return idx
			}
			idx++
		}
	case []interface{}:
		if i, err := strconv.Atoi(key); err == nil {
			return i
		}
	}
	return 0
}

// JSONParse parses a JSON string into a Go data structure.
func JSONParse(jsonStr string) (interface{}, error) {
	decoder := json.NewDecoder(strings.NewReader(jsonStr))
	decoder.UseNumber()

	var result interface{}
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}

	return convertNumbers(result), nil
}

func convertNumbers(v interface{}) interface{} {
	switch val := v.(type) {
	case json.Number:
		if i, err := val.Int64(); err == nil {
			if i >= -2147483648 && i <= 2147483647 {
				return int(i)
			}
			return i
		}
		if f, err := val.Float64(); err == nil {
			return f
		}
		return val.String()
	case map[string]interface{}:
		for k, v := range val {
			val[k] = convertNumbers(v)
		}
		return val
	case []interface{}:
		for i, v := range val {
			val[i] = convertNumbers(v)
		}
		return val
	default:
		return val
	}
}
