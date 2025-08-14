package main

import (
	"encoding/json"
	"fmt"
	"os"
	
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

// normalizeObjectOrdering ensures consistent property ordering to match JavaScript behavior
// JavaScript preserves insertion order, but since Go maps are unordered, we'll sort alphabetically
func normalizeObjectOrdering(value interface{}) interface{} {
	return value // For now, disable normalization to test if the issue is elsewhere
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <jsonpath-or-file> <json-data-or-file> [--query-file] [--data-file]\n", os.Args[0])
		os.Exit(1)
	}

	jsonpathOrFile := os.Args[1]
	jsonDataOrFile := os.Args[2]
	var jsonpath, jsonData string
	
	// Check flags
	queryFromFile := false
	dataFromFile := false
	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--query-file":
			queryFromFile = true
		case "--data-file", "--file":
			dataFromFile = true
		}
	}
	
	// Read JSONPath query
	if queryFromFile {
		queryData, err := os.ReadFile(jsonpathOrFile)
		if err != nil {
			output := map[string]interface{}{
				"error": fmt.Sprintf("Failed to read query file: %v", err),
				"count": 0,
				"values": []interface{}{},
				"paths": []string{},
			}
			jsonBytes, _ := json.Marshal(output)
			fmt.Println(string(jsonBytes))
			return
		}
		jsonpath = string(queryData)
	} else {
		jsonpath = jsonpathOrFile
	}
	
	// Read JSON data
	if dataFromFile {
		fileData, err := os.ReadFile(jsonDataOrFile)
		if err != nil {
			output := map[string]interface{}{
				"error": fmt.Sprintf("Failed to read data file: %v", err),
				"count": 0,
				"values": []interface{}{},
				"paths": []string{},
			}
			jsonBytes, _ := json.Marshal(output)
			fmt.Println(string(jsonBytes))
			return
		}
		jsonData = string(fileData)
	} else {
		jsonData = jsonDataOrFile
	}

	results, err := jp.Query(jsonpath, jsonData)
	if err != nil {
		output := map[string]interface{}{
			"error": err.Error(),
			"count": 0,
			"values": []interface{}{},
			"paths": []string{},
		}
		jsonBytes, _ := json.Marshal(output)
		fmt.Println(string(jsonBytes))
		return
	}

	values := make([]interface{}, len(results))
	paths := make([]string, len(results))
	for i, result := range results {
		values[i] = normalizeObjectOrdering(result.Value)
		paths[i] = result.Path
	}

	output := map[string]interface{}{
		"count": len(results),
		"values": values,
		"paths": paths,
	}

	jsonBytes, err := json.Marshal(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON marshal error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonBytes))
}