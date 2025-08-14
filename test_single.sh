#!/bin/bash

# Test a single JSONPath query with both implementations
# Usage: ./test_single.sh "<jsonpath>" "<data_key>"

if [ $# -ne 2 ]; then
    echo "Usage: $0 '<jsonpath>' '<data_key>'"
    echo "Example: $0 '\$..book[?(@property !== 0)]' 'goessner_spec_data'"
    exit 1
fi

JSONPATH="$1"
DATA_KEY="$2"

# Get the data
DATA=$(node -e "
const registry = JSON.parse(require('fs').readFileSync('./test_data_registry.json', 'utf8'));
console.log(JSON.stringify(registry['$DATA_KEY']));
")

echo "=============================================================="
echo "Testing JSONPath: $JSONPATH"
echo "Data: $DATA_KEY"
echo "=============================================================="
echo

echo "--- Go Implementation ---"
(cd cmd/test_go && go run main.go "$JSONPATH" "$DATA") | jq '.'

echo
echo "--- JS Implementation ---"
node test_js.js "$JSONPATH" "$DATA" | jq '.'

echo
echo "--- Direct Comparison ---"
GO_RESULT=$(cd cmd/test_go && go run main.go "$JSONPATH" "$DATA")
JS_RESULT=$(node test_js.js "$JSONPATH" "$DATA")

GO_COUNT=$(echo "$GO_RESULT" | jq '.count')
JS_COUNT=$(echo "$JS_RESULT" | jq '.count')

if [ "$GO_COUNT" == "$JS_COUNT" ]; then
    echo "✅ Count Match: $GO_COUNT"
else
    echo "❌ Count Mismatch: Go=$GO_COUNT, JS=$JS_COUNT"
fi

GO_VALUES=$(echo "$GO_RESULT" | jq '.values')
JS_VALUES=$(echo "$JS_RESULT" | jq '.values')

if [ "$GO_VALUES" == "$JS_VALUES" ]; then
    echo "✅ Values Match"
else
    echo "❌ Values Mismatch"
    echo "Go Values: $GO_VALUES"
    echo "JS Values: $JS_VALUES"
fi