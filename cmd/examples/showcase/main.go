package main

import (
	"fmt"
	"log"

	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

// Comprehensive showcase of JSONPath Plus Go features
func main() {
	fmt.Println("🚀 JSONPath Plus Go - Comprehensive Showcase")
	fmt.Println("============================================")

	// Create production engine
	engine := createEngine()
	defer engine.Close()

	// Sample data
	data := createSampleData()

	// Demonstrate features
	basicQueries(engine, data)
	indexPreservation(engine, data)
	filterExpressions(engine, data)
	performanceFeatures(engine, data)

	fmt.Println("\n✨ All features demonstrated successfully!")
	fmt.Println("📚 Check the documentation for more details.")
}

func createEngine() *jp.JSONPathEngine {
	fmt.Println("\n1. Creating Production Engine")

	config := jp.ProductionConfig()
	config.EnableLogging = true
	config.EnableMetrics = true

	engine, err := jp.NewEngine(config)
	if err != nil {
		log.Fatalf("❌ Failed to create engine: %v", err)
	}

	fmt.Println("✅ Engine created with production configuration")
	return engine
}

func createSampleData() interface{} {
	jsonData := `{
		"company": {
			"name": "TechCorp International",
			"founded": 2010,
			"locations": ["San Francisco", "New York", "London", "Tokyo"],
			"departments": {
				"engineering": {
					"name": "Engineering",
					"budget": 5000000,
					"teams": [
						{
							"name": "Backend",
							"lead": "Alice Johnson",
							"members": [
								{"name": "Bob Smith", "role": "Senior Engineer", "salary": 120000, "skills": ["Go", "Docker", "Kubernetes"]},
								{"name": "Carol White", "role": "Engineer", "salary": 95000, "skills": ["Python", "Redis", "PostgreSQL"]},
								{"name": "Dave Brown", "role": "Junior Engineer", "salary": 75000, "skills": ["JavaScript", "Node.js", "MongoDB"]}
							]
						},
						{
							"name": "Frontend",
							"lead": "Eve Davis",
							"members": [
								{"name": "Frank Wilson", "role": "Senior Engineer", "salary": 115000, "skills": ["React", "TypeScript", "CSS"]},
								{"name": "Grace Lee", "role": "Engineer", "salary": 90000, "skills": ["Vue.js", "JavaScript", "SASS"]}
							]
						}
					]
				},
				"sales": {
					"name": "Sales",
					"budget": 2000000,
					"regions": [
						{"name": "North America", "revenue": 15000000, "manager": "Henry Garcia"},
						{"name": "Europe", "revenue": 8000000, "manager": "Ivy Martinez"},
						{"name": "Asia", "revenue": 12000000, "manager": "Jack Anderson"}
					]
				}
			},
			"products": [
				{"id": "P001", "name": "CloudSync Pro", "price": 299.99, "category": "Software", "active": true},
				{"id": "P002", "name": "DataFlow Enterprise", "price": 1999.99, "category": "Software", "active": true},
				{"id": "P003", "name": "SecureVault", "price": 149.99, "category": "Security", "active": false},
				{"id": "P004", "name": "AnalyticsDash", "price": 599.99, "category": "Analytics", "active": true}
			]
		},
		"metadata": {
			"version": "2.1",
			"last_updated": "2024-01-15",
			"total_employees": 45,
			"annual_revenue": 35000000
		}
	}`

	data, err := jp.JSONParse(jsonData)
	if err != nil {
		log.Fatalf("❌ Failed to parse JSON: %v", err)
	}

	return data
}

func basicQueries(engine *jp.JSONPathEngine, data interface{}) {
	fmt.Println("\n2. Basic JSONPath Queries")

	queries := []struct {
		path        string
		description string
	}{
		{"$.company.name", "Company name"},
		{"$.company.locations[*]", "All locations"},
		{"$.company.locations[0:2]", "First two locations"},
		{"$.company.locations[-1]", "Last location"},
		{"$..name", "All names recursively"},
		{"$.company.products[*].id", "All product IDs"},
		{"$.company.departments.*", "All departments"},
	}

	for i, query := range queries {
		fmt.Printf("\n  Query %d: %s\n", i+1, query.description)
		fmt.Printf("  Path: %s\n", query.path)

		results, err := engine.QueryData(query.path, data)
		if err != nil {
			fmt.Printf("  ❌ Error: %v\n", err)
			continue
		}

		fmt.Printf("  ✅ Results (%d):\n", len(results))
		for j, result := range results {
			if j < 3 { // Show first 3 results
				fmt.Printf("    [%d] %v (path: %s)\n", j, result.Value, result.Path)
			} else if j == 3 {
				fmt.Printf("    ... and %d more\n", len(results)-3)
				break
			}
		}
	}
}

func indexPreservation(engine *jp.JSONPathEngine, data interface{}) {
	fmt.Println("\n3. Original Index Preservation")

	queries := []struct {
		path        string
		description string
	}{
		{"$.company.locations[*]", "Array elements with original indices"},
		{"$.company.departments.engineering.teams[*].members[*]", "Nested array elements"},
		{"$..skills[*]", "Deeply nested arrays"},
	}

	for i, query := range queries {
		fmt.Printf("\n  Demo %d: %s\n", i+1, query.description)
		fmt.Printf("  Path: %s\n", query.path)

		results, err := engine.QueryData(query.path, data)
		if err != nil {
			fmt.Printf("  ❌ Error: %v\n", err)
			continue
		}

		fmt.Printf("  ✅ Index preservation:\n")
		for j, result := range results {
			if j < 5 { // Show first 5 results
				fmt.Printf("    Value: %-20v | OriginalIndex: %d | Path: %s\n",
					result.Value, result.OriginalIndex, result.Path)
			} else if j == 5 {
				fmt.Printf("    ... and %d more with preserved indices\n", len(results)-5)
				break
			}
		}
	}
}

func filterExpressions(engine *jp.JSONPathEngine, data interface{}) {
	fmt.Println("\n4. Advanced Filter Expressions")

	filters := []struct {
		path        string
		description string
	}{
		{"$.company.products[?(@.active == true)]", "Active products"},
		{"$.company.products[?(@.price > 500)]", "Expensive products"},
		{"$..members[?(@.salary >= 100000)]", "High earners"},
		{"$..members[?(@.role == 'Senior Engineer')]", "Senior engineers"},
		{"$.company.departments.sales.regions[?(@.revenue > 10000000)]", "High revenue regions"},
		{"$..skills[?(@)]", "All skills (existence check)"},
	}

	for i, filter := range filters {
		fmt.Printf("\n  Filter %d: %s\n", i+1, filter.description)
		fmt.Printf("  Expression: %s\n", filter.path)

		results, err := engine.QueryData(filter.path, data)
		if err != nil {
			fmt.Printf("  ❌ Error: %v\n", err)
			continue
		}

		fmt.Printf("  ✅ Matching items (%d):\n", len(results))
		for j, result := range results {
			if j < 3 { // Show first 3 results
				switch v := result.Value.(type) {
				case map[string]interface{}:
					if name, ok := v["name"]; ok {
						fmt.Printf("    [%d] %s\n", j, name)
					} else if id, ok := v["id"]; ok {
						fmt.Printf("    [%d] Product %s\n", j, id)
					} else {
						fmt.Printf("    [%d] %v\n", j, v)
					}
				default:
					fmt.Printf("    [%d] %v\n", j, v)
				}
			} else if j == 3 {
				fmt.Printf("    ... and %d more matches\n", len(results)-3)
				break
			}
		}
	}
}

func performanceFeatures(engine *jp.JSONPathEngine, data interface{}) {
	fmt.Println("\n5. Performance & Production Features")

	// Execute several queries to generate metrics
	queries := []string{
		"$.company.name",
		"$..members[*].name",
		"$.company.products[?(@.active == true)]",
		"$..revenue",
		"$.metadata.*",
	}

	fmt.Println("\n  Executing queries for metrics...")
	for i, query := range queries {
		for j := 0; j < 5; j++ { // Run each query 5 times
			_, err := engine.QueryData(query, data)
			if err != nil {
				fmt.Printf("  ❌ Query %d.%d failed: %v\n", i+1, j+1, err)
			}
		}
	}

	// Show metrics
	metrics := engine.GetMetrics()
	fmt.Printf("\n  📊 Performance Metrics:\n")
	fmt.Printf("    Queries executed: %d\n", metrics.QueriesExecuted)
	fmt.Printf("    Average execution time: %v\n", metrics.AverageExecutionTime)
	fmt.Printf("    Error count: %d\n", metrics.ErrorCount)
	fmt.Printf("    Memory usage: %d bytes\n", metrics.MemoryUsage)

	// Cache functionality has been removed for simplicity

	// Configuration info
	config := engine.GetConfig()
	fmt.Printf("\n  ⚙️  Configuration:\n")
	fmt.Printf("    Max path length: %d\n", config.MaxPathLength)
	fmt.Printf("    Max recursion depth: %d\n", config.MaxRecursionDepth)
	fmt.Printf("    Max result count: %d\n", config.MaxResultCount)
	fmt.Printf("    Timeout: %v\n", config.Timeout)
	fmt.Printf("    Strict mode: %v\n", config.StrictMode)

	// Security demo
	fmt.Printf("\n  🔒 Security Validation:\n")
	validator := jp.NewSecurityValidator(jp.DefaultSecurityConfig())

	testPaths := []string{
		"$.company.name",                // Safe
		"$.users[?(eval('malicious'))]", // Unsafe
	}

	for _, path := range testPaths {
		err := validator.ValidatePath(path)
		if err != nil {
			fmt.Printf("    ❌ BLOCKED: %s\n", path)
		} else {
			fmt.Printf("    ✅ ALLOWED: %s\n", path)
		}
	}
}
