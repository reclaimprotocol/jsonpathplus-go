# String Index Preservation Report - v2.0.0

## Status: 🟢 Fully Supporting Value-Centric Ranges

### Summary
In v2.0.0, JSONPath-Plus Go has been promoted to a major version to introduce **Value-Centric Indexing** by default. Unlike traditional JSONPath libraries, we now track the entire character range of object properties (key-to-value).

### Comparison: v1 vs v2

| Feature | v1 Behavior | v2 Behavior (Default) |
|---------|-------------|-----------------------|
| Object Property | Start of key | **Key to End of Value** |
| Array Element | Start of element | Start of element |
| Whitespace | Preserved | Preserved |
| Complexity | Key-only | **Range (Start, End, Length)** |

### Example Mapping

JSON: `{"id": 123, "name": "test"}`

| Path | v1 Index | v2 Range [Start:End] | v2 Length |
|------|----------|----------------------|-----------|
| `$.id` | 1 | `[1:10]` | 9 |
| `$.name` | 12 | `[12:26]` | 14 |

### Implementation Details
- **`Start`**: Marks the beginning of the property key.
- **`End`**: Marks the end of the JSON value (exclusive).
- **`Length`**: Calculates `End - Start`, representing the total segment in the original JSON.

---
*Verified as of v2.0.0*