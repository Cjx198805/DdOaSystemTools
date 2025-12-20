package utils

import (
	"encoding/json"
	"io"
)

// ParseJSON 解析JSON响应
func ParseJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
