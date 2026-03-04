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

// SetAlias fieldsTree parameter used to set Alias for nested structure.
// Can be omitted to set Alias for current structure Selector or
// passed few values for several levels of nesting
func (s *Selector) SetAlias(alias string, fieldsTree ...string) error {
	if len(fieldsTree) == 0 {
		s.Alias = alias
		return nil
	}

	name := fieldsTree[0]

	for _, field := range s.Fields {
		if field.Name == name {
			err := field.SetAlias(alias, fieldsTree[1:]...)
			if err != nil {
				return fmt.Errorf("%s.%w", name, err)
			}
		}
	}

	return fmt.Errorf("%s field not found", name)
}

// GetAlias fieldsTree parameter used to get Alias of nested structure.
// Can be omitted to get Alias of current structure Selector or
// passed few values for several levels of nesting
func (s *Selector) GetAlias(fieldsTree ...string) (string, error) {
	if len(fieldsTree) == 0 {
		return s.Alias, nil
	}

	name := fieldsTree[0]

	for _, field := range s.Fields {
		if field.Name == name {
			alias, err := field.GetAlias(fieldsTree[1:]...)
			if err != nil {
				return "", fmt.Errorf("%s.%w", name, err)
			}

			return alias, nil
		}
	}

	return "", fmt.Errorf("%s field not found", name)
}
