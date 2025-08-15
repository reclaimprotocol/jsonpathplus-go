package utils

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ParseValue parses a string value into the appropriate type
func ParseValue(s string) interface{} {
	s = strings.TrimSpace(s)

	// Remove quotes if present
	if (strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'")) ||
		(strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) {
		return s[1 : len(s)-1]
	}

	// Try boolean
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}

	// Try null
	if s == "null" {
		return nil
	}

	// Try number
	if num, err := strconv.ParseFloat(s, 64); err == nil {
		return num
	}

	return s
}

// CompareValues compares two values using the given operator
func CompareValues(left interface{}, operator string, right interface{}) bool {
	// Handle nil values
	if left == nil && right == nil {
		return operator == "==" || operator == "==="
	}
	if left == nil || right == nil {
		return operator == "!=" || operator == "!=="
	}

	// Handle string comparisons
	leftStr, leftIsStr := left.(string)
	rightStr, rightIsStr := right.(string)

	if leftIsStr && rightIsStr {
		switch operator {
		case "==", "===":
			return leftStr == rightStr
		case "!=", "!==":
			return leftStr != rightStr
		case "<":
			return leftStr < rightStr
		case "<=":
			return leftStr <= rightStr
		case ">":
			return leftStr > rightStr
		case ">=":
			return leftStr >= rightStr
		}
	}

	// Handle numeric comparisons
	leftNum := toFloat64(left)
	rightNum := toFloat64(right)

	if isNumeric(left) && isNumeric(right) {
		switch operator {
		case "==", "===":
			return leftNum == rightNum
		case "!=", "!==":
			return leftNum != rightNum
		case "<":
			return leftNum < rightNum
		case "<=":
			return leftNum <= rightNum
		case ">":
			return leftNum > rightNum
		case ">=":
			return leftNum >= rightNum
		}
	}

	// Default to equality check for other types
	switch operator {
	case "==", "===":
		return reflect.DeepEqual(left, right)
	case "!=", "!==":
		return !reflect.DeepEqual(left, right)
	default:
		return false
	}
}

// GetPropertyValue gets a property value from an object, supporting nested access
func GetPropertyValue(obj interface{}, property string) interface{} {
	if obj == nil {
		return nil
	}

	// Handle nested property access (e.g., "profile.bio")
	if strings.Contains(property, ".") {
		return getNestedProperty(obj, property)
	}

	switch v := obj.(type) {
	case *OrderedMap:
		value, _ := v.Get(property)
		return value
	case map[string]interface{}:
		return v[property]
	case map[interface{}]interface{}:
		return v[property]
	default:
		// Use reflection for struct fields
		rv := reflect.ValueOf(obj)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		if rv.Kind() == reflect.Struct {
			field := rv.FieldByName(property)
			if field.IsValid() {
				return field.Interface()
			}
		}
		return nil
	}
}

// getNestedProperty gets a nested property value
func getNestedProperty(obj interface{}, path string) interface{} {
	if obj == nil {
		return nil
	}

	parts := strings.Split(path, ".")
	current := obj

	for _, part := range parts {
		current = GetPropertyValue(current, part)
		if current == nil {
			return nil
		}
	}

	return current
}

// isNumeric checks if a value is numeric
func isNumeric(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

// toFloat64 converts a value to float64
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case uint8:
		return float64(val)
	case uint16:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	case float32:
		return float64(val)
	case float64:
		return val
	default:
		return 0
	}
}

// Function predicate implementations

// TryMatchFunction handles regex matching
func TryMatchFunction(expr string, current interface{}) (bool, bool) {
	// Try regex literal format with optional flags: .match(/pattern/flags)
	reRegexLiteral := regexp.MustCompile(`\.match\(/(.+?)/([gimsx]*)\)`)
	matches := reRegexLiteral.FindStringSubmatch(expr)

	if len(matches) >= 2 {
		pattern := matches[1]
		flags := ""
		if len(matches) == 3 {
			flags = matches[2]
		}

		// Get the actual string value
		str := getStringValue(expr, current)
		if str == "" {
			return false, true
		}

		// Apply flags to pattern
		if strings.Contains(flags, "i") {
			pattern = "(?i)" + pattern
		}
		if strings.Contains(flags, "m") {
			pattern = "(?m)" + pattern
		}
		if strings.Contains(flags, "s") {
			pattern = "(?s)" + pattern
		}

		// Compile and match the regex
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return false, true
		}

		return regex.MatchString(str), true
	}

	// Try quoted string format: .match("pattern") or .match('pattern')
	reQuoted := regexp.MustCompile(`\.match\(['"](.+?)['"]\)`)
	matches = reQuoted.FindStringSubmatch(expr)
	if len(matches) != 2 {
		return false, false
	}

	pattern := matches[1]

	// Get the actual string value
	str := getStringValue(expr, current)
	if str == "" {
		return false, true
	}

	// Compile and match the regex
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false, true
	}

	return regex.MatchString(str), true
}

// TryContainsFunction handles contains checking for strings and arrays
func TryContainsFunction(expr string, current interface{}) (bool, bool) {
	re := regexp.MustCompile(`\.contains\(['"](.+?)['"]\)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 2 {
		return false, false
	}

	searchTerm := matches[1]

	// Get the base property
	baseProperty := getBaseProperty(expr)
	value := getPropertyValueForFunction(current, baseProperty)

	if value == nil {
		return false, true
	}

	switch v := value.(type) {
	case string:
		return strings.Contains(v, searchTerm), true
	case []interface{}:
		// For arrays, check if any item contains the search term as a string
		for _, item := range v {
			if str, ok := item.(string); ok && strings.Contains(str, searchTerm) {
				return true, true
			}
			// Also check string representation for exact match
			if fmt.Sprintf("%v", item) == searchTerm {
				return true, true
			}
		}
		return false, true
	default:
		str := fmt.Sprintf("%v", v)
		return strings.Contains(str, searchTerm), true
	}
}

// TryStartsWithFunction handles startsWith checking
func TryStartsWithFunction(expr string, current interface{}) (bool, bool) {
	re := regexp.MustCompile(`\.startsWith\(['"](.+?)['"]\)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 2 {
		return false, false
	}

	prefix := matches[1]
	baseProperty := getBaseProperty(expr)
	value := getPropertyValueForFunction(current, baseProperty)

	if value == nil {
		return false, true
	}

	str, ok := value.(string)
	if !ok {
		str = fmt.Sprintf("%v", value)
	}

	return strings.HasPrefix(str, prefix), true
}

// TryEndsWithFunction handles endsWith checking
func TryEndsWithFunction(expr string, current interface{}) (bool, bool) {
	re := regexp.MustCompile(`\.endsWith\(['"](.+?)['"]\)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 2 {
		return false, false
	}

	suffix := matches[1]
	baseProperty := getBaseProperty(expr)
	value := getPropertyValueForFunction(current, baseProperty)

	if value == nil {
		return false, true
	}

	str, ok := value.(string)
	if !ok {
		str = fmt.Sprintf("%v", value)
	}

	return strings.HasSuffix(str, suffix), true
}

// TryLengthFunction handles length comparisons
func TryLengthFunction(expr string, current interface{}) (bool, bool) {
	re := regexp.MustCompile(`\.length\s*(===|!==|<=|>=|==|!=|<|>)\s*(\d+)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 3 {
		return false, false
	}

	operator := matches[1]
	expectedLength, err := strconv.Atoi(matches[2])
	if err != nil {
		return false, false
	}

	baseProperty := getBaseProperty(expr)
	value := getPropertyValueForFunction(current, baseProperty)

	if value == nil {
		// JavaScript throws "Cannot read properties of null (reading 'length')"
		// We return false, false to indicate this is an invalid operation that should fail the entire query
		return false, false
	}

	var actualLength int
	switch v := value.(type) {
	case string:
		actualLength = len(v)
	case []interface{}:
		actualLength = len(v)
	case map[string]interface{}:
		actualLength = len(v)
	default:
		return false, true
	}

	switch operator {
	case "===", "==":
		return actualLength == expectedLength, true
	case "!==", "!=":
		return actualLength != expectedLength, true
	case "<":
		return actualLength < expectedLength, true
	case "<=":
		return actualLength <= expectedLength, true
	case ">":
		return actualLength > expectedLength, true
	case ">=":
		return actualLength >= expectedLength, true
	}

	return false, true
}

// TryCaseFunction handles case conversion functions
func TryCaseFunction(expr string, current interface{}) (bool, bool) {
	// Handle .toLowerCase() === 'value' or .toLowerCase() == 'value'
	lowerRe := regexp.MustCompile(`\.toLowerCase\(\)\s*(===|!==|==|!=)\s*['"](.+?)['"]`)
	if matches := lowerRe.FindStringSubmatch(expr); len(matches) == 3 {
		operator := matches[1]
		expectedValue := matches[2]

		baseProperty := getBaseProperty(expr)
		value := getPropertyValueForFunction(current, baseProperty)

		if value == nil {
			return false, true
		}

		str, ok := value.(string)
		if !ok {
			str = fmt.Sprintf("%v", value)
		}

		lowerStr := strings.ToLower(str)
		switch operator {
		case "===", "==":
			return lowerStr == expectedValue, true
		case "!==", "!=":
			return lowerStr != expectedValue, true
		}
	}

	// Handle .toUpperCase() === 'value' or .toUpperCase() == 'value'
	upperRe := regexp.MustCompile(`\.toUpperCase\(\)\s*(===|!==|==|!=)\s*['"](.+?)['"]`)
	if matches := upperRe.FindStringSubmatch(expr); len(matches) == 3 {
		operator := matches[1]
		expectedValue := matches[2]

		baseProperty := getBaseProperty(expr)
		value := getPropertyValueForFunction(current, baseProperty)

		if value == nil {
			return false, true
		}

		str, ok := value.(string)
		if !ok {
			str = fmt.Sprintf("%v", value)
		}

		upperStr := strings.ToUpper(str)
		switch operator {
		case "===", "==":
			return upperStr == expectedValue, true
		case "!==", "!=":
			return upperStr != expectedValue, true
		}
	}

	return false, false
}

// TryTypeofFunction handles type checking
func TryTypeofFunction(expr string, current interface{}) (bool, bool) {
	re := regexp.MustCompile(`\.typeof\(\)\s*(===|!==|==|!=)\s*['"](.+?)['"]`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 3 {
		return false, false
	}

	operator := matches[1]
	expectedType := matches[2]

	baseProperty := getBaseProperty(expr)
	value := getPropertyValueForFunction(current, baseProperty)

	var actualType string
	switch value.(type) {
	case string:
		actualType = "string"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		actualType = "number"
	case bool:
		actualType = "boolean"
	case []interface{}:
		actualType = "array"
	case map[string]interface{}:
		actualType = "object"
	case nil:
		actualType = "null"
	default:
		actualType = "unknown"
	}

	switch operator {
	case "===", "==":
		return actualType == expectedType, true
	case "!==", "!=":
		return actualType != expectedType, true
	}

	return false, true
}

// TryChainedOperations handles method chaining like .toLowerCase().contains()
func TryChainedOperations(expr string, current interface{}) (bool, bool) {
	// Handle .toLowerCase().contains()
	lowerContainsRe := regexp.MustCompile(`\.toLowerCase\(\)\.contains\(['"](.+?)['"]\)`)
	if matches := lowerContainsRe.FindStringSubmatch(expr); len(matches) == 2 {
		searchTerm := matches[1]
		str := getStringValue(expr, current)
		if str == "" {
			return false, true
		}

		lowerStr := strings.ToLower(str)
		return strings.Contains(lowerStr, searchTerm), true
	}

	// Handle .toUpperCase().contains()
	upperContainsRe := regexp.MustCompile(`\.toUpperCase\(\)\.contains\(['"](.+?)['"]\)`)
	if matches := upperContainsRe.FindStringSubmatch(expr); len(matches) == 2 {
		searchTerm := matches[1]
		str := getStringValue(expr, current)
		if str == "" {
			return false, true
		}

		upperStr := strings.ToUpper(str)
		return strings.Contains(upperStr, searchTerm), true
	}

	return false, false
}

// Helper functions for function predicates

func getStringValue(expr string, current interface{}) string {
	baseProperty := getBaseProperty(expr)
	value := getPropertyValueForFunction(current, baseProperty)

	if value == nil {
		return ""
	}

	if str, ok := value.(string); ok {
		return str
	}

	return fmt.Sprintf("%v", value)
}

func getBaseProperty(expr string) string {
	// Extract the base property from expressions like ".name.contains()" or ".name.startsWith()"
	if strings.HasPrefix(expr, ".") {
		// Find the function call first
		for _, suffix := range []string{".contains(", ".startsWith(", ".endsWith(", ".match(", ".length", ".toLowerCase(", ".toUpperCase(", ".typeof(", ".floor(", ".round(", ".ceil("} {
			if idx := strings.Index(expr, suffix); idx != -1 {
				if idx == 0 {
					// Expression like ".contains(...)" means current value (@)
					return ""
				}
				// Extract property from ".property.function(" -> "property"
				property := expr[1:idx]
				return property
			}
		}

		// If no function found, look for nested property access
		dotIndex := strings.Index(expr[1:], ".")
		if dotIndex == -1 {
			// Single property like ".name"
			return strings.TrimPrefix(expr, ".")
		}
		// Nested property like ".user.name" -> "user.name"
		return expr[1:]
	}
	return ""
}

func getPropertyValueForFunction(current interface{}, property string) interface{} {
	if property == "" {
		return current
	}

	return GetPropertyValue(current, property)
}

// Math function implementations

// TryFloorFunction handles floor() math function
func TryFloorFunction(expr string, current interface{}) (bool, bool) {
	re := regexp.MustCompile(`\.floor\(\)\s*(===|!==|<=|>=|==|!=|<|>)\s*(\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 3 {
		return false, false
	}

	operator := matches[1]
	expectedValue, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return false, false
	}

	baseProperty := getBaseProperty(expr)
	value := getPropertyValueForFunction(current, baseProperty)

	if value == nil {
		return false, true
	}

	if !isNumeric(value) {
		return false, true
	}

	floatValue := toFloat64(value)
	floorValue := math.Floor(floatValue)

	switch operator {
	case "===", "==":
		return floorValue == expectedValue, true
	case "!==", "!=":
		return floorValue != expectedValue, true
	case "<":
		return floorValue < expectedValue, true
	case "<=":
		return floorValue <= expectedValue, true
	case ">":
		return floorValue > expectedValue, true
	case ">=":
		return floorValue >= expectedValue, true
	}

	return false, true
}

// TryRoundFunction handles round() math function
func TryRoundFunction(expr string, current interface{}) (bool, bool) {
	re := regexp.MustCompile(`\.round\(\)\s*(===|!==|<=|>=|==|!=|<|>)\s*(\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 3 {
		return false, false
	}

	operator := matches[1]
	expectedValue, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return false, false
	}

	baseProperty := getBaseProperty(expr)
	value := getPropertyValueForFunction(current, baseProperty)

	if value == nil {
		return false, true
	}

	if !isNumeric(value) {
		return false, true
	}

	floatValue := toFloat64(value)
	roundValue := math.Round(floatValue)

	switch operator {
	case "===", "==":
		return roundValue == expectedValue, true
	case "!==", "!=":
		return roundValue != expectedValue, true
	case "<":
		return roundValue < expectedValue, true
	case "<=":
		return roundValue <= expectedValue, true
	case ">":
		return roundValue > expectedValue, true
	case ">=":
		return roundValue >= expectedValue, true
	}

	return false, true
}

// TryCeilFunction handles ceil() math function
func TryCeilFunction(expr string, current interface{}) (bool, bool) {
	re := regexp.MustCompile(`\.ceil\(\)\s*(===|!==|<=|>=|==|!=|<|>)\s*(\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 3 {
		return false, false
	}

	operator := matches[1]
	expectedValue, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return false, false
	}

	baseProperty := getBaseProperty(expr)
	value := getPropertyValueForFunction(current, baseProperty)

	if value == nil {
		return false, true
	}

	if !isNumeric(value) {
		return false, true
	}

	floatValue := toFloat64(value)
	ceilValue := math.Ceil(floatValue)

	switch operator {
	case "===", "==":
		return ceilValue == expectedValue, true
	case "!==", "!=":
		return ceilValue != expectedValue, true
	case "<":
		return ceilValue < expectedValue, true
	case "<=":
		return ceilValue <= expectedValue, true
	case ">":
		return ceilValue > expectedValue, true
	case ">=":
		return ceilValue >= expectedValue, true
	}

	return false, true
}
