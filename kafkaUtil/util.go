package kafkaUtil

import (
	"encoding/json"
	"strconv"
)

// decodeBase64OrRaw tries to decode the input as base64, returns raw bytes if not base64
func decodeBase64OrRaw(data []byte) []byte {
	decoded, err := strconv.Unquote("\"" + string(data) + "\"")
	if err == nil {
		return []byte(decoded)
	}
	return data
}

// decodeJSONOrRaw tries to decode the input as JSON, returns raw bytes if not JSON
func decodeJSONOrRaw(data []byte) interface{} {
	var v interface{}
	if err := json.Unmarshal(data, &v); err == nil {
		return v
	}
	return data
}
