package internal

import (
	"encoding/json"
)

// Marshal marshals a value to JSON.
func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal unmarshals JSON data into a value.
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
