package selector

import "fmt"

const defaultName = "default"

type Selector struct {
	// Name of the field to select
	Name string
	// Alist to use on selecting field
	// or to set generated value
	Alias string
	// Any parameters to filter instances
	Params map[string]any
	Fields []*Selector
}

func Select(input ...any) []*Selector {
	var results []*Selector
	for _, in := range input {
		results = append(results, normalize(in)...)
	}
	return results
}

func normalize(input any) []*Selector {
	switch in := input.(type) {
	case *Selector:
		return []*Selector{in}
	case []*Selector:
		return in
	case string:
		return []*Selector{{Name: in}}
	case []string:
		list := make([]*Selector, len(in))
		for i, s := range in {
			list[i] = &Selector{Name: s}
		}
		return list
	default:
		panic(fmt.Sprintf("type %T not supported as Selector", input))
	}
}

func DefaultSelector() *Selector {
	return &Selector{Name: defaultName}
}

func GetParameter[T any](s *Selector, name string) (res T, ok bool) {
	p, ok := s.Params[name]
	if !ok {
		return res, false
	}

	res, ok = p.(T)
	return res, ok
}

func (s *Selector) IsDefault() bool {
	return s.Name == defaultName && len(s.Params) == 0 && len(s.Fields) == 0
}
