const { JSONPath } = require('jsonpath-plus');
const data = {users: {'1': {name: 'Alice'}, '2': {name: 'Bob'}}};

console.log('Testing path filters:');
const query = "$.users[?(@path === \"$['users']['1']\")]";
console.log('Query:', query);
const result1 = JSONPath({ path: query, json: data });
console.log('Result:', result1);
console.log('Count:', result1.length);

console.log('\nFor comparison, get path of user 1:');
const result2 = JSONPath({ path: '$.users.1', json: data, resultType: 'all' });
console.log('Paths:', result2.map(r => r.path));

console.log('\nTesting different path formats:');
const result3 = JSONPath({ path: '$.users[*]', json: data, resultType: 'all' });
result3.forEach(r => {
    console.log('Value:', r.value, 'Path:', r.path);
});