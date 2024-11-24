package utils

import (
	"encoding/json"
	"io"
)

// EncodeJSON encodes the given data to JSON and writes it to the provided writer.
func EncodeJSON(w io.Writer, data interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// DecodeJSON decodes JSON data from the provided reader into the given target.
func DecodeJSON(r io.Reader, target interface{}) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(target)
}
