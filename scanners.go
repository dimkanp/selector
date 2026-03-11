package selector

import (
	"encoding/json"
	"fmt"
	"strings"
)

// IgnoreNullSlice used in JsonScanner in case when database returned "[null]"
// and unmarshalling should return empty slice instead of slice with one default element.
// [null] will be treated like []
var IgnoreNullSlice = true

// SetIgnoreNullSlice change JsonScanner.Scan behavior in case when
// database returned "[null]" but this value should be treated like "[]".
//
// By default, ignoring enabled.
//
// Reason of implementing such functionality lies in
// postgres json_agg() function which returns "[null]" when
// no rows in set are present.
func SetIgnoreNullSlice(ignore bool) {
	IgnoreNullSlice = ignore
}

// JsonScanner made for fields that selected as a JSON
// to scan JSON data and parse it into Value
type JsonScanner[T any] struct {
	Value *T
}

func (r *JsonScanner[T]) Scan(src any) error {
	if src == nil {
		src = []byte("null") // leave all handling to the encoding/json package
	}

	data, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan Row[%T]", src)
	}

	data = checkNullSlice(data)

	err := json.Unmarshal(data, r.Value)
	if err != nil {
		return fmt.Errorf("unmarshal scanned JSON data into %T: %w", r.Value, err)
	}

	return nil
}

// replace "[null]" with "[]" when ignoring null slice enabled
func checkNullSlice(input []byte) []byte {
	if !IgnoreNullSlice {
		return input
	}

	// when input equal to "[null]" with removed all trailing spaces and ignored letters case
	trimmed := strings.TrimSpace(string(input))
	if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
		trimmed = trimmed[1 : len(trimmed)-1] // remove '[' at the start and ']' at the end
		if strings.ToLower(strings.TrimSpace(trimmed)) == "null" {
			return []byte("[]")
		}
	}

	return input
}

// ScanJson returns structure ready to scan JSON data from database and parse it in the object passed as parameter p.
// Returned structure are affected by IgnoreNullSlice variable.
func ScanJson[T any](p *T) *JsonScanner[T] {
	return &JsonScanner[T]{
		Value: p,
	}
}

type FirstScanner[T any] struct {
	Value *T
}

// ScanFirst is like ScanJson but expects slice to be scanned
// after scanning it set first element from slice to passed value by pointer p.
//
// Expected to use in cases where query cant return JSON object of single structure
// but can return slice of single element, for example instead
// Postgres function row_to_json() used json_agg() for some reason
func ScanFirst[T any](p *T) *FirstScanner[T] {
	return &FirstScanner[T]{
		Value: p,
	}
}

func (s *FirstScanner[T]) Scan(src any) error {
	data, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan Row[%T]", src)
	}

	var slice []T
	err := json.Unmarshal(data, &slice)
	if err != nil {
		return fmt.Errorf("unmarshal scanned JSON data into %T: %w", s.Value, err)
	}

	if len(slice) == 0 {
		// set default value
		slice = make([]T, 1)
	}

	*s.Value = slice[0]
	return nil
}
