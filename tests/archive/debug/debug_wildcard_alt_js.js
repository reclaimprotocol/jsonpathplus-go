const { JSONPath } = require('jsonpath-plus');

const data = {
    "orders": [
        {"id": "ORD001", "items": [{"product": "laptop", "price": 999.99}, {"product": "mouse", "price": 29.99}]},
        {"id": "ORD002", "items": [{"product": "keyboard", "price": 79.99}]},
        {"id": "ORD003", "items": [{"product": "laptop", "price": 1299.99}]}
    ]
};

// Try alternative approaches
const queries = [
    '$.orders[?(@.items[*].product === \'laptop\')]', // Original failing
    '$.orders[?(@.items[0].product === \'laptop\')]', // Specific index
    '$.orders[?(@.items[?(@.product === \'laptop\')])]', // Nested filter
    '$..orders[?(@.items[?(@.product === \'laptop\')])]' // With recursive descent
];

queries.forEach((query, index) => {
    console.log(`\nTesting query ${index + 1}: ${query}`);
    try {
        const results = JSONPath({path: query, json: data});
        console.log(`Count: ${results.length}`);
        results.forEach((result, i) => {
            console.log(`[${i}] ID: ${result.id}`);
        });
    } catch (error) {
        console.log(`Error: ${error.message}`);
    }
});