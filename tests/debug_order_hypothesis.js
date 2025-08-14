const JSONPath = require('jsonpath-plus');
const fs = require('fs');

// Test data with clear property ordering
const testData = {
    store: {
        book: [
            { price: 12.99 },
            { price: 8.99 },
            { price: 22.99 }
        ],
        bicycle: {
            price: 19.95
        }
    }
};

console.log("=== HYPOTHESIS TEST: Object Property Processing Order ===");
console.log("Test data keys order:", Object.keys(testData.store));
console.log("Should be: ['book', 'bicycle']");

console.log("\n=== Test 1: Plain recursive descent $..*  ===");
const plainResults = JSONPath.JSONPath({
    path: "$..*",
    json: testData,
    resultType: 'all'
});

console.log("Plain $.* results (just price properties):");
plainResults.forEach((result, i) => {
    if (result.parentProperty === 'price') {
        console.log(`${i+1}. Path: ${result.path}, Value: ${result.value}`);
    }
});

console.log("\n=== Test 2: Recursive descent with filter $..*[?(@property === 'price')] ===");
const filterResults = JSONPath.JSONPath({
    path: "$..*[?(@property === 'price')]",
    json: testData,
    resultType: 'all'
});

console.log("Filter results:");
filterResults.forEach((result, i) => {
    console.log(`${i+1}. Path: ${result.path}, Value: ${result.value}`);
});

console.log("\n=== Test 3: Check Object.keys() order ===");
console.log("Object.keys(testData.store):", Object.keys(testData.store));

console.log("\n=== Test 4: Reverse order test ===");
const reverseData = {
    store: {
        bicycle: { price: 19.95 },
        book: [
            { price: 12.99 },
            { price: 8.99 }, 
            { price: 22.99 }
        ]
    }
};

console.log("Reverse data keys:", Object.keys(reverseData.store));

const reverseFilterResults = JSONPath.JSONPath({
    path: "$..*[?(@property === 'price')]",
    json: reverseData,
    resultType: 'all'
});

console.log("Reverse order filter results:");
reverseFilterResults.forEach((result, i) => {
    console.log(`${i+1}. Path: ${result.path}, Value: ${result.value}`);
});