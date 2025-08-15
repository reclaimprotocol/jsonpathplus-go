const { JSONPath } = require('jsonpath-plus');

const data = {"users": ["alice", "bob", "charlie", "alice"]};
const query = '$.users[-1]';

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