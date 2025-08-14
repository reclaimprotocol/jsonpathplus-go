const JSONPath = require('jsonpath-plus');
const fs = require('fs');

const data = JSON.parse(fs.readFileSync('data/mixed_types.json', 'utf8'));
const query = "$.data[?(@.length > 3)]";

console.log("Testing query:", query);
console.log("Data:", JSON.stringify(data, null, 2));

try {
    const results = JSONPath.JSONPath({
        path: query,
        json: data,
        resultType: 'all'
    });
    console.log("JS Results:", results.length);
    results.forEach((result, i) => {
        console.log(`${i+1}. Value: ${JSON.stringify(result.value)}`);
    });
} catch (error) {
    console.log("JS Error:", error.message);
    console.log("Error occurred when processing query");
}

// Let's test each element individually
console.log("\n=== Testing each element's length property ===");
data.data.forEach((item, i) => {
    console.log(`Element ${i} (${typeof item}):`, JSON.stringify(item));
    try {
        const length = item?.length;
        console.log(`  .length = ${length}`);
        console.log(`  .length > 3 = ${length > 3}`);
    } catch (e) {
        console.log(`  Error accessing .length: ${e.message}`);
    }
});