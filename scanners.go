package selector

import (
	"encoding/json"
	"fmt"
)

type RowScanner[T any] struct {
	Value *T
}

func (r *RowScanner[T]) Scan(src any) error {
	data, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan Row[%T]", src)
	}

	err := json.Unmarshal(data, r.Value)
	if err != nil {
		return err
	}

	return nil
}

func Scan[T any](p *T) *RowScanner[T] {
	return &RowScanner[T]{
		Value: p,
	}
}
