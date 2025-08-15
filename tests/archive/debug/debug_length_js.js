const { JSONPath } = require('jsonpath-plus');

const data = {"data":[42,"hello",true,null,{"key":"value"},[1,2,3]]};
const query = '$.data[?(@.length > 3)]';

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