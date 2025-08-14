#!/usr/bin/env node

const { JSONPath } = require('jsonpath-plus');

if (process.argv.length !== 4) {
  console.error(`Usage: ${process.argv[1]} <jsonpath> <json-data>`);
  process.exit(1);
}

const jsonpath = process.argv[2];
const jsonData = process.argv[3];

try {
  const data = JSON.parse(jsonData);
  
  const values = JSONPath({
    path: jsonpath,
    json: data,
    resultType: 'value'
  });
  
  const paths = JSONPath({
    path: jsonpath, 
    json: data,
    resultType: 'path'
  });
  
  const output = {
    count: values.length,
    values: values,
    paths: paths
  };
  
  console.log(JSON.stringify(output));
} catch (error) {
  const output = {
    error: error.message,
    count: 0,
    values: [],
    paths: []
  };
  
  console.log(JSON.stringify(output));
}