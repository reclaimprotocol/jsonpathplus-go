const JSONPath = require('jsonpath-plus');
const fs = require('fs');

const data = JSON.parse(fs.readFileSync('data/goessner_spec_data.json', 'utf8'));
const query = "$..*[?(@property === 'price' && @ !== 8.95)]";

console.log("JS Results for query:", query);
console.log("Data keys:", Object.keys(data));
console.log("Store keys:", Object.keys(data.store));

try {
    const results = JSONPath.JSONPath({
        path: query,
        json: data,
        resultType: 'all'
    });
    
    console.log("\nJS Results:", results.length);
    results.forEach((result, i) => {
        console.log(`${i+1}. Path: ${result.path}, Value: ${result.value}, Property: ${result.parentProperty}`);
    });
    
    // Also test with individual price checks
    console.log("\n=== Testing individual price values ===");
    const allPrices = JSONPath.JSONPath({
        path: "$..price", 
        json: data,
        resultType: 'all'
    });
    console.log("All price paths:");
    allPrices.forEach((result, i) => {
        console.log(`${i+1}. Path: ${result.path}, Value: ${result.value}, NotEqual8.95: ${result.value !== 8.95}`);
    });
} catch (error) {
    console.log("JS Error:", error.message);
}