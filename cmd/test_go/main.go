package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <jsonpath> <json-data>\n", os.Args[0])
		os.Exit(1)
	}

	jsonpath := os.Args[1]
	jsonData := os.Args[2]

	results, err := jp.Query(jsonpath, jsonData)
	if err != nil {
		output := map[string]interface{}{
			"error":  err.Error(),
			"count":  0,
			"values": []interface{}{},
			"paths":  []string{},
		}
		jsonBytes, _ := json.Marshal(output)
		fmt.Println(string(jsonBytes))
		return
	}

	values := make([]interface{}, len(results))
	paths := make([]string, len(results))
	for i, result := range results {
		values[i] = result.Value
		paths[i] = result.Path
	}

	output := map[string]interface{}{
		"count":  len(results),
		"values": values,
		"paths":  paths,
	}

	jsonBytes, err := json.Marshal(output)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonBytes))
}
