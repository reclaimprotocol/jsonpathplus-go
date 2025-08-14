const { JSONPath } = require('jsonpath-plus');

const data = {
    "orders": [
        {"id": "ORD001", "items": [{"product": "laptop", "price": 999.99}, {"product": "mouse", "price": 29.99}]},
        {"id": "ORD002", "items": [{"product": "keyboard", "price": 79.99}]},
        {"id": "ORD003", "items": [{"product": "laptop", "price": 1299.99}]}
    ]
};
const query = '$.orders[?(@.items[*].product === \'laptop\')]';

console.log(`Testing: ${query}`);

try {
    const results = JSONPath({path: query, json: data});
    console.log(`Count: ${results.length}`);
    results.forEach((result, i) => {
        console.log(`[${i}] Value: ${JSON.stringify(result)}`);
    });
} catch (error) {
    console.log(`Error: ${error.message}`);
}