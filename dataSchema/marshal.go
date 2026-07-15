package dataSchema

import "bytes"

// MergeJSONObjects flat-merges several JSON-encoded objects into a single JSON
// object, preserving the order of the given parts.
//
// It replaces the fragile hand-rolled "strip the braces and splice a comma"
// approach used in the custom MarshalJSON implementations. That approach
// produced invalid JSON (e.g. `,,` or `{,`) whenever one of the parts was an
// empty object `{}` or a non-object value such as `null`, which surfaced as:
//
//	json: error calling MarshalJSON for type *interaction.Property:
//	invalid character ',' looking for beginning of object key string
//
// Parts that are not JSON objects (empty input, `null`, `{}`, arrays, ...) are
// skipped. The caller remains responsible for ensuring the merged objects do
// not declare duplicate keys.
func MergeJSONObjects(parts ...[]byte) ([]byte, error) {
	inners := make([][]byte, 0, len(parts))
	for _, part := range parts {
		part = bytes.TrimSpace(part)
		// Only merge actual JSON objects ("{...}"). Skip anything else
		// (empty, "null", arrays, ...).
		if len(part) < 2 || part[0] != '{' || part[len(part)-1] != '}' {
			continue
		}
		inner := bytes.TrimSpace(part[1 : len(part)-1])
		if len(inner) == 0 { // empty object "{}"
			continue
		}
		inners = append(inners, inner)
	}

	result := make([]byte, 0)
	result = append(result, '{')
	result = append(result, bytes.Join(inners, []byte{','})...)
	result = append(result, '}')
	return result, nil
}
