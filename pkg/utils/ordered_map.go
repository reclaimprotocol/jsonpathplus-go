package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// OrderedMap represents a map that preserves insertion order
type OrderedMap struct {
	keys   []string
	values map[string]interface{}
}

// NewOrderedMap creates a new ordered map
func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		keys:   make([]string, 0),
		values: make(map[string]interface{}),
	}
}

// Set sets a key-value pair
func (om *OrderedMap) Set(key string, value interface{}) {
	if _, exists := om.values[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.values[key] = value
}

// Get gets a value by key
func (om *OrderedMap) Get(key string) (interface{}, bool) {
	value, exists := om.values[key]
	return value, exists
}

// Delete deletes a key-value pair
func (om *OrderedMap) Delete(key string) {
	if _, exists := om.values[key]; exists {
		delete(om.values, key)
		// Remove from keys slice
		for i, k := range om.keys {
			if k == key {
				om.keys = append(om.keys[:i], om.keys[i+1:]...)
				break
			}
		}
	}
}

// Keys returns all keys in insertion order
func (om *OrderedMap) Keys() []string {
	return append([]string(nil), om.keys...)
}

// Values returns all values in insertion order
func (om *OrderedMap) Values() []interface{} {
	values := make([]interface{}, len(om.keys))
	for i, key := range om.keys {
		values[i] = om.values[key]
	}
	return values
}

// Len returns the number of key-value pairs
func (om *OrderedMap) Len() int {
	return len(om.keys)
}

// Range iterates over key-value pairs in insertion order
func (om *OrderedMap) Range(fn func(key string, value interface{}) bool) {
	for _, key := range om.keys {
		if !fn(key, om.values[key]) {
			break
		}
	}
}

// ToMap converts to a regular map[string]interface{}
func (om *OrderedMap) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	for _, key := range om.keys {
		result[key] = om.values[key]
	}
	return result
}

// MarshalJSON implements json.Marshaler to preserve order during JSON serialization
func (om *OrderedMap) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	
	for i, key := range om.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		
		// Marshal key
		keyBytes, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)
		buf.WriteByte(':')
		
		// Marshal value
		valueBytes, err := json.Marshal(om.values[key])
		if err != nil {
			return nil, err
		}
		buf.Write(valueBytes)
	}
	
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshaler to preserve order during JSON deserialization
func (om *OrderedMap) UnmarshalJSON(data []byte) error {
	// Reset the map
	om.keys = om.keys[:0]
	for k := range om.values {
		delete(om.values, k)
	}
	
	// Parse JSON to preserve order
	decoder := json.NewDecoder(bytes.NewReader(data))
	
	// Expect '{'
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	if delim, ok := token.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expected '{', got %v", token)
	}
	
	// Parse key-value pairs
	for decoder.More() {
		// Read key
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		key, ok := token.(string)
		if !ok {
			return fmt.Errorf("expected string key, got %v", token)
		}
		
		// Read value
		var value interface{}
		if err := decoder.Decode(&value); err != nil {
			return err
		}
		
		om.Set(key, value)
	}
	
	// Expect '}'
	token, err = decoder.Token()
	if err != nil {
		return err
	}
	if delim, ok := token.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("expected '}', got %v", token)
	}
	
	return nil
}

// ParseOrderedJSON parses JSON while preserving object property order
func ParseOrderedJSON(data []byte) (interface{}, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	return parseOrderedValue(decoder)
}

func parseOrderedValue(decoder *json.Decoder) (interface{}, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	
	switch token := token.(type) {
	case json.Delim:
		if token == '{' {
			// Parse object with preserved order
			om := NewOrderedMap()
			for decoder.More() {
				// Read key
				keyToken, err := decoder.Token()
				if err != nil {
					return nil, err
				}
				key, ok := keyToken.(string)
				if !ok {
					return nil, fmt.Errorf("expected string key, got %v", keyToken)
				}
				
				// Read value recursively
				value, err := parseOrderedValue(decoder)
				if err != nil {
					return nil, err
				}
				
				om.Set(key, value)
			}
			
			// Expect '}'
			if _, err := decoder.Token(); err != nil {
				return nil, err
			}
			
			return om, nil
		} else if token == '[' {
			// Parse array
			var arr []interface{}
			for decoder.More() {
				value, err := parseOrderedValue(decoder)
				if err != nil {
					return nil, err
				}
				arr = append(arr, value)
			}
			
			// Expect ']'
			if _, err := decoder.Token(); err != nil {
				return nil, err
			}
			
			return arr, nil
		}
	default:
		// Primitive value
		return token, nil
	}
	
	return nil, fmt.Errorf("unexpected token: %v", token)
}

// Convert a regular interface{} to ordered structure
func ConvertToOrdered(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		om := NewOrderedMap()
		// Sort keys to get consistent order (since Go map iteration is random)
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		
		for _, key := range keys {
			om.Set(key, ConvertToOrdered(v[key]))
		}
		return om
	case []interface{}:
		for i, item := range v {
			v[i] = ConvertToOrdered(item)
		}
		return v
	default:
		return v
	}
}