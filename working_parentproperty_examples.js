#!/usr/bin/env node

const { JSONPath } = require('jsonpath-plus');

// Test working examples of @parentProperty
const testData = {
  "company": {
    "departments": {
      "engineering": {
        "employees": [
          {"name": "Alice", "level": "senior"},
          {"name": "Bob", "level": "junior"}
        ]
      },
      "sales": {
        "employees": [
          {"name": "Charlie", "level": "senior"}
        ]
      }
    },
    "metadata": {
      "founded": 2010
    }
  }
};

console.log("=== Working @parentProperty Examples ===\n");

// Let's trace through what @parentProperty should be at each level
console.log("Data structure analysis:");
console.log("- company (root)");  
console.log("  - departments (parentProperty = 'company')");
console.log("    - engineering (parentProperty = 'departments')"); 
console.log("      - employees (parentProperty = 'engineering')");
console.log("        - [0] Alice (parentProperty = 'employees')");
console.log("        - [1] Bob (parentProperty = 'employees')");
console.log("    - sales (parentProperty = 'departments')");
console.log("      - employees (parentProperty = 'sales')");
console.log("        - [0] Charlie (parentProperty = 'employees')");
console.log("  - metadata (parentProperty = 'company')");
console.log("");

const workingExamples = [
  {
    name: "Array elements with parentProperty = 'employees'",
    query: "$.company.departments.*.employees[?(@parentProperty === 'employees')]",
    description: "Employees where their parent array property is 'employees'"
  },
  {
    name: "Array elements using numeric property", 
    query: "$.company.departments.*.employees[?(@property === 0)]",
    description: "First employee in each department"
  },
  {
    name: "Departments with parentProperty = 'departments'",
    query: "$.company.departments[?(@parentProperty === 'company')]", 
    description: "Departments object where parent is company"
  },
  {
    name: "All company children",
    query: "$.company.*[?(@parentProperty === 'company')]",
    description: "Direct children of company"
  },
  {
    name: "Test what @parentProperty returns for departments.*",
    query: "$.company.departments.*",
    description: "Just get departments to understand structure"
  }
];

workingExamples.forEach((example, index) => {
  console.log(`${index + 1}. ${example.name}`);
  console.log(`   Query: ${example.query}`);
  
  try {
    const result = JSONPath({
      path: example.query,
      json: testData, 
      resultType: 'value'
    });
    
    const paths = JSONPath({
      path: example.query,
      json: testData,
      resultType: 'path'  
    });
    
    console.log(`   Results: ${result.length} items`);
    result.forEach((val, i) => {
      console.log(`     [${i}] Path: ${paths[i]}`);
      console.log(`         Value: ${JSON.stringify(val)}`);
    });
    
  } catch (error) {
    console.log(`   ❌ ERROR: ${error.message}`);
  }
  
  console.log(`   Description: ${example.description}`);
  console.log("");
});

// Now test the original failing query with a simpler structure
console.log("=== Testing Original Failing Query Logic ===");

console.log("Original failing query: $.company.departments.*[?(@parentProperty === 'departments')]");
console.log("This should find: engineering and sales objects");
console.log("Because: both are properties of the 'departments' object");
console.log("");

// Let's manually trace what @parentProperty should be:
console.log("Manual trace:");
console.log("- $.company.departments.* gets engineering and sales objects");  
console.log("- For engineering object: @parentProperty should be ???");
console.log("- For sales object: @parentProperty should be ???");
console.log("");

// Test different interpretations
const interpretations = [
  "$.company.departments.*[?(@parentProperty === 'company')]",  // Parent of departments
  "$.company.*[?(@parentProperty === 'company')]",              // Children of company  
  "$.company.departments.*[?(@property === 'engineering')]",    // Filter by property name
  "$.company.departments[?(@parentProperty === 'company')]"     // Departments object itself
];

interpretations.forEach((query, i) => {
  console.log(`Interpretation ${i + 1}: ${query}`);
  try {
    const result = JSONPath({ path: query, json: testData, resultType: 'value' });
    console.log(`  Results: ${result.length} items`);
    result.forEach((val, j) => console.log(`    [${j}] ${JSON.stringify(val).substring(0, 50)}...`));
  } catch (e) {
    console.log(`  ERROR: ${e.message}`);
  }
  console.log("");
});