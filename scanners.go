package selector

import (
	"encoding/json"
	"fmt"
)

// JsonScanner made for fields that selected as a JSON
// to scan JSON data and parse it into Value
type JsonScanner[T any] struct {
	Value *T
}

func (r *JsonScanner[T]) Scan(src any) error {
	data, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan Row[%T]", src)
	}

	err := json.Unmarshal(data, r.Value)
	if err != nil {
		return fmt.Errorf("unmarshal scanned JSON data into %T: %w", r.Value, err)
	}

	return nil
}

// ScanJson returns structure ready to scan JSON data from database and parse it in the object passed as parameter p
func ScanJson[T any](p *T) *JsonScanner[T] {
	return &JsonScanner[T]{
		Value: p,
	}
}
